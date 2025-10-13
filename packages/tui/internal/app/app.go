package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"log/slog"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/aaronmrosenthal/rycode-sdk-go"
	"github.com/aaronmrosenthal/rycode/internal/auth"
	"github.com/aaronmrosenthal/rycode/internal/clipboard"
	"github.com/aaronmrosenthal/rycode/internal/commands"
	"github.com/aaronmrosenthal/rycode/internal/components/toast"
	"github.com/aaronmrosenthal/rycode/internal/id"
	"github.com/aaronmrosenthal/rycode/internal/styles"
	"github.com/aaronmrosenthal/rycode/internal/theme"
	"github.com/aaronmrosenthal/rycode/internal/util"
)

type Message struct {
	Info  opencode.MessageUnion
	Parts []opencode.PartUnion
}

type App struct {
	Project           opencode.Project
	Agents            []opencode.Agent
	Providers         []opencode.Provider
	Version           string
	StatePath         string
	Config            *opencode.Config
	Client            *opencode.Client
	State             *State
	AgentIndex        int
	Provider          *opencode.Provider
	Model             *opencode.Model
	Session           *opencode.Session
	Messages          []Message
	Permissions       []opencode.Permission
	CurrentPermission opencode.Permission
	Commands          commands.CommandRegistry
	InitialModel      *string
	InitialPrompt     *string
	InitialAgent      *string
	InitialSession    *string
	compactCancel     context.CancelFunc
	IsLeaderSequence  bool
	IsBashMode        bool
	ScrollSpeed       int
	AuthBridge        *auth.Bridge // Auth system bridge
	CurrentCost       float64      // Cached cost from auth system
	LastCostUpdate    time.Time    // When cost was last fetched
}

func (a *App) Agent() *opencode.Agent {
	return &a.Agents[a.AgentIndex]
}

type SessionCreatedMsg = struct {
	Session *opencode.Session
}
type SessionSelectedMsg = *opencode.Session
type MessageRevertedMsg struct {
	Session opencode.Session
	Message Message
}
type SessionUnrevertedMsg struct {
	Session opencode.Session
}
type SessionLoadedMsg struct{}
type ModelSelectedMsg struct {
	Provider opencode.Provider
	Model    opencode.Model
}

type AgentSelectedMsg struct {
	AgentName string
}

type SessionClearedMsg struct{}
type CompactSessionMsg struct{}

// CostUpdatedMsg is sent when cost summary is updated
type CostUpdatedMsg struct {
	Cost float64
}
type SendPrompt = Prompt
type SendShell = struct {
	Command string
}
type SendCommand = struct {
	Command string
	Args    string
}
type SetEditorContentMsg struct {
	Text string
}
type FileRenderedMsg struct {
	FilePath string
}
type PermissionRespondedToMsg struct {
	Response opencode.SessionPermissionRespondParamsResponse
}

func New(
	ctx context.Context,
	version string,
	project *opencode.Project,
	path *opencode.Path,
	agents []opencode.Agent,
	httpClient *opencode.Client,
	initialModel *string,
	initialPrompt *string,
	initialAgent *string,
	initialSession *string,
) (*App, error) {
	util.RootPath = project.Worktree
	util.CwdPath, _ = os.Getwd()

	configInfo, err := httpClient.Config.Get(ctx, opencode.ConfigGetParams{})
	if err != nil {
		return nil, err
	}

	if configInfo.Keybinds.Leader == "" {
		configInfo.Keybinds.Leader = "ctrl+x"
	}

	appStatePath := filepath.Join(path.State, "tui")
	appState, err := LoadState(appStatePath)
	if err != nil {
		appState = NewState()
		SaveState(appStatePath, appState)
	}

	if appState.AgentModel == nil {
		appState.AgentModel = make(map[string]AgentModel)
	}

	if configInfo.Theme != "" {
		appState.Theme = configInfo.Theme
	}

	themeEnv := os.Getenv("OPENCODE_THEME")
	if themeEnv != "" {
		appState.Theme = themeEnv
	}

	agentIndex := slices.IndexFunc(agents, func(a opencode.Agent) bool {
		return a.Mode != "subagent"
	})
	var agent *opencode.Agent
	modeName := "build"
	if appState.Agent != "" {
		modeName = appState.Agent
	}
	if initialAgent != nil && *initialAgent != "" {
		modeName = *initialAgent
	}
	for i, m := range agents {
		if m.Name == modeName {
			agentIndex = i
			break
		}
	}
	agent = &agents[agentIndex]

	if agent.Model.ModelID != "" {
		appState.AgentModel[agent.Name] = AgentModel{
			ProviderID: agent.Model.ProviderID,
			ModelID:    agent.Model.ModelID,
		}
	}

	if err := theme.LoadThemesFromDirectories(
		path.Config,
		util.RootPath,
		util.CwdPath,
	); err != nil {
		slog.Warn("Failed to load themes from directories", "error", err)
	}

	if appState.Theme != "" {
		if appState.Theme == "system" && styles.Terminal != nil {
			theme.UpdateSystemTheme(
				styles.Terminal.Background,
				styles.Terminal.BackgroundIsDark,
			)
		}
		theme.SetTheme(appState.Theme)
	}

	slog.Debug("Loaded config", "config", configInfo)

	customCommands, err := httpClient.Command.List(ctx, opencode.CommandListParams{})
	if err != nil {
		return nil, err
	}

	app := &App{
		Project:        *project,
		Agents:         agents,
		Version:        version,
		StatePath:      appStatePath,
		Config:         configInfo,
		State:          appState,
		Client:         httpClient,
		AgentIndex:     agentIndex,
		Session:        &opencode.Session{},
		Messages:       []Message{},
		Commands:       commands.LoadFromConfig(configInfo, *customCommands),
		InitialModel:   initialModel,
		InitialPrompt:  initialPrompt,
		InitialAgent:   initialAgent,
		InitialSession: initialSession,
		ScrollSpeed:    int(configInfo.Tui.ScrollSpeed),
		AuthBridge:     auth.NewBridge(project.Worktree),
		CurrentCost:    0.0,
		LastCostUpdate: time.Now(),
	}

	return app, nil
}

func (a *App) Keybind(commandName commands.CommandName) string {
	command := a.Commands[commandName]
	if len(command.Keybindings) == 0 {
		return ""
	}
	kb := command.Keybindings[0]
	key := kb.Key
	if kb.RequiresLeader {
		key = a.Config.Keybinds.Leader + " " + kb.Key
	}
	return key
}

func (a *App) Key(commandName commands.CommandName) string {
	t := theme.CurrentTheme()
	base := styles.NewStyle().Background(t.Background()).Foreground(t.Text()).Bold(true).Render
	muted := styles.NewStyle().
		Background(t.Background()).
		Foreground(t.TextMuted()).
		Faint(true).
		Render
	command := a.Commands[commandName]
	key := a.Keybind(commandName)
	return base(key) + muted(" "+command.Description)
}

func SetClipboard(text string) tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, func() tea.Msg {
		clipboard.Write(clipboard.FmtText, []byte(text))
		return nil
	})
	// try to set the clipboard using OSC52 for terminals that support it
	cmds = append(cmds, tea.SetClipboard(text))
	return tea.Sequence(cmds...)
}

func (a *App) cycleMode(forward bool) (*App, tea.Cmd) {
	if forward {
		a.AgentIndex++
		if a.AgentIndex >= len(a.Agents) {
			a.AgentIndex = 0
		}
	} else {
		a.AgentIndex--
		if a.AgentIndex < 0 {
			a.AgentIndex = len(a.Agents) - 1
		}
	}
	if a.Agent().Mode == "subagent" {
		return a.cycleMode(forward)
	}

	modelID := a.Agent().Model.ModelID
	providerID := a.Agent().Model.ProviderID
	if modelID == "" {
		if model, ok := a.State.AgentModel[a.Agent().Name]; ok {
			modelID = model.ModelID
			providerID = model.ProviderID
		}
	}

	if modelID != "" {
		for _, provider := range a.Providers {
			if provider.ID == providerID {
				a.Provider = &provider
				for _, model := range provider.Models {
					if model.ID == modelID {
						a.Model = &model
						break
					}
				}
				break
			}
		}
	}

	a.State.Agent = a.Agent().Name
	a.State.UpdateAgentUsage(a.Agent().Name)
	return a, a.SaveState()
}

func (a *App) SwitchAgent() (*App, tea.Cmd) {
	return a.cycleMode(true)
}

func (a *App) SwitchAgentReverse() (*App, tea.Cmd) {
	return a.cycleMode(false)
}

func (a *App) cycleRecentModel(forward bool) (*App, tea.Cmd) {
	recentModels := a.State.RecentlyUsedModels
	if len(recentModels) > 5 {
		recentModels = recentModels[:5]
	}
	if len(recentModels) < 2 {
		return a, toast.NewInfoToast("Need at least 2 recent models to cycle")
	}
	nextIndex := 0
	prevIndex := 0
	for i, recentModel := range recentModels {
		if a.Provider != nil && a.Model != nil && recentModel.ProviderID == a.Provider.ID &&
			recentModel.ModelID == a.Model.ID {
			nextIndex = (i + 1) % len(recentModels)
			prevIndex = (i - 1 + len(recentModels)) % len(recentModels)
			break
		}
	}
	targetIndex := nextIndex
	if !forward {
		targetIndex = prevIndex
	}
	for range recentModels {
		currentRecentModel := recentModels[targetIndex%len(recentModels)]
		provider, model := findModelByProviderAndModelID(
			a.Providers,
			currentRecentModel.ProviderID,
			currentRecentModel.ModelID,
		)
		if provider != nil && model != nil {
			a.Provider, a.Model = provider, model
			a.State.AgentModel[a.Agent().Name] = AgentModel{
				ProviderID: provider.ID,
				ModelID:    model.ID,
			}
			return a, tea.Sequence(
				a.SaveState(),
				toast.NewSuccessToast(
					fmt.Sprintf("Switched to %s (%s)", model.Name, provider.Name),
				),
			)
		}
		recentModels = append(
			recentModels[:targetIndex%len(recentModels)],
			recentModels[targetIndex%len(recentModels)+1:]...)
		if len(recentModels) < 2 {
			a.State.RecentlyUsedModels = recentModels
			return a, tea.Sequence(
				a.SaveState(),
				toast.NewInfoToast("Not enough valid recent models to cycle"),
			)
		}
	}
	a.State.RecentlyUsedModels = recentModels
	return a, toast.NewErrorToast("Recent model not found")
}

func (a *App) CycleRecentModel() (*App, tea.Cmd) {
	return a.cycleRecentModel(true)
}

func (a *App) CycleRecentModelReverse() (*App, tea.Cmd) {
	return a.cycleRecentModel(false)
}

// CycleAuthenticatedProviders cycles through all authenticated providers
func (a *App) CycleAuthenticatedProviders(forward bool) (*App, tea.Cmd) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get authenticated CLI providers
	cliProviders, err := a.AuthBridge.GetCLIProviders(ctx)
	if err != nil {
		return a, toast.NewErrorToast("Failed to get CLI providers")
	}

	// Filter to only authenticated providers
	authenticatedProviders := []string{}
	for _, cliProv := range cliProviders {
		authStatus, err := a.AuthBridge.CheckAuthStatus(ctx, cliProv.Provider)
		if err != nil || !authStatus.IsAuthenticated {
			continue
		}
		authenticatedProviders = append(authenticatedProviders, cliProv.Provider)
	}

	if len(authenticatedProviders) == 0 {
		return a, toast.NewInfoToast("No authenticated providers. Press 'd' in /model to auto-detect.")
	}

	if len(authenticatedProviders) == 1 {
		return a, toast.NewInfoToast("Only one provider authenticated")
	}

	// Find current provider index in authenticated list
	currentIndex := -1
	for i, provID := range authenticatedProviders {
		if a.Provider != nil && provID == a.Provider.ID {
			currentIndex = i
			break
		}
	}

	// Calculate next index
	nextIndex := 0
	if currentIndex != -1 {
		if forward {
			nextIndex = (currentIndex + 1) % len(authenticatedProviders)
		} else {
			nextIndex = (currentIndex - 1 + len(authenticatedProviders)) % len(authenticatedProviders)
		}
	}

	// Get next provider ID
	nextProviderID := authenticatedProviders[nextIndex]

	// Find the actual provider in a.Providers
	var nextProvider *opencode.Provider
	var nextModel *opencode.Model

	for i := range a.Providers {
		if a.Providers[i].ID == nextProviderID {
			nextProvider = &a.Providers[i]
			break
		}
	}

	if nextProvider == nil {
		return a, toast.NewErrorToast(fmt.Sprintf("Provider %s not found in providers list", nextProviderID))
	}

	// Try to find the most recently used model for this provider
	for _, recentModel := range a.State.RecentlyUsedModels {
		if recentModel.ProviderID == nextProvider.ID {
			// Find this model in the provider's models
			for _, model := range nextProvider.Models {
				if model.ID == recentModel.ModelID {
					nextModel = &model
					break
				}
			}
			if nextModel != nil {
				break
			}
		}
	}

	// If no recent model, use the first available model
	if nextModel == nil && len(nextProvider.Models) > 0 {
		for _, model := range nextProvider.Models {
			nextModel = &model
			break
		}
	}

	if nextModel == nil {
		return a, toast.NewErrorToast(fmt.Sprintf("No models found for %s", nextProvider.Name))
	}

	// Update app state
	a.Provider = nextProvider
	a.Model = nextModel
	a.State.AgentModel[a.Agent().Name] = AgentModel{
		ProviderID: nextProvider.ID,
		ModelID:    nextModel.ID,
	}
	a.State.UpdateModelUsage(nextProvider.ID, nextModel.ID)

	return a, tea.Sequence(
		a.SaveState(),
		toast.NewSuccessToast(
			fmt.Sprintf("→ %s: %s", nextProvider.Name, nextModel.Name),
		),
	)
}

func (a *App) CycleAuthenticatedProvider() (*App, tea.Cmd) {
	return a.CycleAuthenticatedProviders(true)
}

func (a *App) CycleAuthenticatedProviderReverse() (*App, tea.Cmd) {
	return a.CycleAuthenticatedProviders(false)
}

func (a *App) SwitchToAgent(agentName string) (*App, tea.Cmd) {
	// Find the agent index by name
	for i, agent := range a.Agents {
		if agent.Name == agentName {
			a.AgentIndex = i
			break
		}
	}

	// Set up model for the new agent
	modelID := a.Agent().Model.ModelID
	providerID := a.Agent().Model.ProviderID
	if modelID == "" {
		if model, ok := a.State.AgentModel[a.Agent().Name]; ok {
			modelID = model.ModelID
			providerID = model.ProviderID
		}
	}

	if modelID != "" {
		for _, provider := range a.Providers {
			if provider.ID == providerID {
				a.Provider = &provider
				for _, model := range provider.Models {
					if model.ID == modelID {
						a.Model = &model
						break
					}
				}
				break
			}
		}
	}

	a.State.Agent = a.Agent().Name
	a.State.UpdateAgentUsage(agentName)
	return a, a.SaveState()
}

// findModelByFullID finds a model by its full ID in the format "provider/model"
func findModelByFullID(
	providers []opencode.Provider,
	fullModelID string,
) (*opencode.Provider, *opencode.Model) {
	modelParts := strings.SplitN(fullModelID, "/", 2)
	if len(modelParts) < 2 {
		return nil, nil
	}

	providerID := modelParts[0]
	modelID := modelParts[1]

	return findModelByProviderAndModelID(providers, providerID, modelID)
}

// findModelByProviderAndModelID finds a model by provider ID and model ID
func findModelByProviderAndModelID(
	providers []opencode.Provider,
	providerID, modelID string,
) (*opencode.Provider, *opencode.Model) {
	for _, provider := range providers {
		if provider.ID != providerID {
			continue
		}

		for _, model := range provider.Models {
			if model.ID == modelID {
				return &provider, &model
			}
		}

		// Provider found but model not found
		return nil, nil
	}

	// Provider not found
	return nil, nil
}

// findProviderByID finds a provider by its ID
func findProviderByID(providers []opencode.Provider, providerID string) *opencode.Provider {
	for _, provider := range providers {
		if provider.ID == providerID {
			return &provider
		}
	}
	return nil
}

func (a *App) InitializeProvider() tea.Cmd {
	ctx := context.Background()

	// Get merged providers (HTTP API + CLI providers)
	providers, err := a.ListProviders(ctx)
	if err != nil {
		slog.Error("Failed to list providers", "error", err)
		// TODO: notify user
		return nil
	}
	if len(providers) == 0 {
		slog.Error("No providers configured")
		return nil
	}

	// Get the HTTP-only response for default model selection
	providersResponse, err := a.Client.App.Providers(ctx, opencode.AppProvidersParams{})
	if err != nil {
		slog.Warn("Failed to get provider response for defaults", "error", err)
	}

	a.Providers = providers

	// Auto-detect credentials on every startup (silent if none found)
	var autoDetectCmd tea.Cmd
	autoDetectCmd = a.autoDetectAllCredentialsQuiet()

	// retains backwards compatibility with old state format
	if model, ok := a.State.AgentModel[a.State.Agent]; ok {
		a.State.Provider = model.ProviderID
		a.State.Model = model.ModelID
	}

	var selectedProvider *opencode.Provider
	var selectedModel *opencode.Model

	// Priority 1: Command line --model flag (InitialModel)
	if a.InitialModel != nil && *a.InitialModel != "" {
		if provider, model := findModelByFullID(providers, *a.InitialModel); provider != nil &&
			model != nil {
			selectedProvider = provider
			selectedModel = model
			slog.Debug(
				"Selected model from command line",
				"provider",
				provider.ID,
				"model",
				model.ID,
			)
		} else {
			slog.Debug("Command line model not found", "model", *a.InitialModel)
		}
	}

	// Priority 2: Config file model setting
	if selectedProvider == nil && a.Config.Model != "" {
		if provider, model := findModelByFullID(providers, a.Config.Model); provider != nil &&
			model != nil {
			selectedProvider = provider
			selectedModel = model
			slog.Debug("Selected model from config", "provider", provider.ID, "model", model.ID)
		} else {
			slog.Debug("Config model not found", "model", a.Config.Model)
		}
	}

	// Priority 3: Current agent's preferred model
	if selectedProvider == nil && a.Agent().Model.ModelID != "" {
		if provider, model := findModelByProviderAndModelID(providers, a.Agent().Model.ProviderID, a.Agent().Model.ModelID); provider != nil &&
			model != nil {
			selectedProvider = provider
			selectedModel = model
			slog.Debug(
				"Selected model from current agent",
				"provider",
				provider.ID,
				"model",
				model.ID,
				"agent",
				a.Agent().Name,
			)
		} else {
			slog.Debug("Agent model not found", "provider", a.Agent().Model.ProviderID, "model", a.Agent().Model.ModelID, "agent", a.Agent().Name)
		}
	}

	// Priority 4: Recent model usage (most recently used model)
	if selectedProvider == nil && len(a.State.RecentlyUsedModels) > 0 {
		recentUsage := a.State.RecentlyUsedModels[0] // Most recent is first
		if provider, model := findModelByProviderAndModelID(providers, recentUsage.ProviderID, recentUsage.ModelID); provider != nil &&
			model != nil {
			selectedProvider = provider
			selectedModel = model
			slog.Debug(
				"Selected model from recent usage",
				"provider",
				provider.ID,
				"model",
				model.ID,
			)
		} else {
			slog.Debug("Recent model not found", "provider", recentUsage.ProviderID, "model", recentUsage.ModelID)
		}
	}

	// Priority 5: State-based model (backwards compatibility)
	if selectedProvider == nil && a.State.Provider != "" && a.State.Model != "" {
		if provider, model := findModelByProviderAndModelID(providers, a.State.Provider, a.State.Model); provider != nil &&
			model != nil {
			selectedProvider = provider
			selectedModel = model
			slog.Debug("Selected model from state", "provider", provider.ID, "model", model.ID)
		} else {
			slog.Debug("State model not found", "provider", a.State.Provider, "model", a.State.Model)
		}
	}

	// Priority 6: Internal priority fallback (Anthropic preferred, then first available)
	if selectedProvider == nil {
		// Try Anthropic first as internal priority
		if provider := findProviderByID(providers, "anthropic"); provider != nil {
			if model := getDefaultModel(providersResponse, *provider); model != nil {
				selectedProvider = provider
				selectedModel = model
				slog.Debug(
					"Selected model from internal priority (Anthropic)",
					"provider",
					provider.ID,
					"model",
					model.ID,
				)
			}
		}

		// If Anthropic not available, use first available provider
		if selectedProvider == nil && len(providers) > 0 {
			provider := &providers[0]
			if model := getDefaultModel(providersResponse, *provider); model != nil {
				selectedProvider = provider
				selectedModel = model
				slog.Debug(
					"Selected model from fallback (first available)",
					"provider",
					provider.ID,
					"model",
					model.ID,
				)
			}
		}
	}

	// Final safety check
	if selectedProvider == nil || selectedModel == nil {
		slog.Error("Failed to select any model")
		return nil
	}

	var cmds []tea.Cmd
	cmds = append(cmds, util.CmdHandler(ModelSelectedMsg{
		Provider: *selectedProvider,
		Model:    *selectedModel,
	}))

	// Add auto-detect command if this is first run
	if autoDetectCmd != nil {
		cmds = append(cmds, autoDetectCmd)
	}

	// Load initial session if provided
	if a.InitialSession != nil && *a.InitialSession != "" {
		cmds = append(cmds, func() tea.Msg {
			// Find the session by ID
			sessions, err := a.ListSessions(context.Background())
			if err != nil {
				slog.Error("Failed to list sessions for initial session", "error", err)
				return toast.NewErrorToast("Failed to load initial session")()
			}

			for _, session := range sessions {
				if session.ID == *a.InitialSession {
					return SessionSelectedMsg(&session)
				}
			}

			slog.Warn("Initial session not found", "sessionID", *a.InitialSession)
			return toast.NewErrorToast("Session not found: " + *a.InitialSession)()
		})
	}

	if a.InitialPrompt != nil && *a.InitialPrompt != "" {
		cmds = append(cmds, util.CmdHandler(SendPrompt{Text: *a.InitialPrompt}))
	}
	return tea.Batch(cmds...)
}

func getDefaultModel(
	response *opencode.AppProvidersResponse,
	provider opencode.Provider,
) *opencode.Model {
	if match, ok := response.Default[provider.ID]; ok {
		model := provider.Models[match]
		return &model
	} else {
		for _, model := range provider.Models {
			return &model
		}
	}
	return nil
}

func (a *App) IsBusy() bool {
	if len(a.Messages) == 0 {
		return false
	}
	if a.IsCompacting() {
		return true
	}
	lastMessage := a.Messages[len(a.Messages)-1]
	if casted, ok := lastMessage.Info.(opencode.AssistantMessage); ok {
		return casted.Time.Completed == 0
	}
	return false
}

func (a *App) IsCompacting() bool {
	if time.Since(time.UnixMilli(int64(a.Session.Time.Compacting))) < time.Second*30 {
		return true
	}
	return false
}

func (a *App) HasAnimatingWork() bool {
	for _, msg := range a.Messages {
		switch casted := msg.Info.(type) {
		case opencode.AssistantMessage:
			if casted.Time.Completed == 0 {
				return true
			}
		}
		for _, p := range msg.Parts {
			if tp, ok := p.(opencode.ToolPart); ok {
				if tp.State.Status == opencode.ToolPartStateStatusPending {
					return true
				}
			}
		}
	}
	return false
}

func (a *App) SaveState() tea.Cmd {
	return func() tea.Msg {
		err := SaveState(a.StatePath, a.State)
		if err != nil {
			slog.Error("Failed to save state", "error", err)
		}
		return nil
	}
}

// UpdateCost fetches the latest cost from the auth bridge
func (a *App) UpdateCost() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		summary, err := a.AuthBridge.GetCostSummary(ctx)
		if err != nil {
			slog.Debug("Failed to get cost summary", "error", err)
			return nil
		}

		return CostUpdatedMsg{Cost: summary.TodayCost}
	}
}

func (a *App) InitializeProject(ctx context.Context) tea.Cmd {
	cmds := []tea.Cmd{}

	session, err := a.CreateSession(ctx)
	if err != nil {
		// status.Error(err.Error())
		return nil
	}

	a.Session = session
	cmds = append(cmds, util.CmdHandler(SessionCreatedMsg{Session: session}))

	go func() {
		_, err := a.Client.Session.Init(ctx, a.Session.ID, opencode.SessionInitParams{
			MessageID:  opencode.F(id.Ascending(id.Message)),
			ProviderID: opencode.F(a.Provider.ID),
			ModelID:    opencode.F(a.Model.ID),
		})
		if err != nil {
			slog.Error("Failed to initialize project", "error", err)
			// status.Error(err.Error())
		}
	}()

	return tea.Batch(cmds...)
}

func (a *App) CompactSession(ctx context.Context) tea.Cmd {
	if a.compactCancel != nil {
		a.compactCancel()
	}

	compactCtx, cancel := context.WithCancel(ctx)
	a.compactCancel = cancel

	go func() {
		defer func() {
			a.compactCancel = nil
		}()

		_, err := a.Client.Session.Summarize(
			compactCtx,
			a.Session.ID,
			opencode.SessionSummarizeParams{
				ProviderID: opencode.F(a.Provider.ID),
				ModelID:    opencode.F(a.Model.ID),
			},
		)
		if err != nil {
			if compactCtx.Err() != context.Canceled {
				slog.Error("Failed to compact session", "error", err)
			}
		}
	}()
	return nil
}

func (a *App) MarkProjectInitialized(ctx context.Context) error {
	return nil
	/*
		_, err := a.Client.App.Init(ctx)
		if err != nil {
			slog.Error("Failed to mark project as initialized", "error", err)
			return err
		}
		return nil
	*/
}

func (a *App) CreateSession(ctx context.Context) (*opencode.Session, error) {
	session, err := a.Client.Session.New(ctx, opencode.SessionNewParams{})
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (a *App) SendPrompt(ctx context.Context, prompt Prompt) (*App, tea.Cmd) {
	var cmds []tea.Cmd
	if a.Session.ID == "" {
		session, err := a.CreateSession(ctx)
		if err != nil {
			return a, toast.NewErrorToast(err.Error())
		}
		a.Session = session
		cmds = append(cmds, util.CmdHandler(SessionCreatedMsg{Session: session}))
	}

	messageID := id.Ascending(id.Message)
	message := prompt.ToMessage(messageID, a.Session.ID)

	a.Messages = append(a.Messages, message)

	cmds = append(cmds, func() tea.Msg {
		_, err := a.Client.Session.Prompt(ctx, a.Session.ID, opencode.SessionPromptParams{
			Model: opencode.F(opencode.SessionPromptParamsModel{
				ProviderID: opencode.F(a.Provider.ID),
				ModelID:    opencode.F(a.Model.ID),
			}),
			Agent:     opencode.F(a.Agent().Name),
			MessageID: opencode.F(messageID),
			Parts:     opencode.F(message.ToSessionChatParams()),
		})
		if err != nil {
			errormsg := fmt.Sprintf("failed to send message: %v", err)
			slog.Error(errormsg)
			return toast.NewErrorToast(errormsg)()
		}
		return nil
	})

	// The actual response will come through SSE
	// For now, just return success
	return a, tea.Batch(cmds...)
}

func (a *App) SendCommand(ctx context.Context, command string, args string) (*App, tea.Cmd) {
	var cmds []tea.Cmd
	if a.Session.ID == "" {
		session, err := a.CreateSession(ctx)
		if err != nil {
			return a, toast.NewErrorToast(err.Error())
		}
		a.Session = session
		cmds = append(cmds, util.CmdHandler(SessionCreatedMsg{Session: session}))
	}

	cmds = append(cmds, func() tea.Msg {
		params := opencode.SessionCommandParams{
			Command:   opencode.F(command),
			Arguments: opencode.F(args),
			Agent:     opencode.F(a.Agents[a.AgentIndex].Name),
		}
		if a.Provider != nil && a.Model != nil {
			params.Model = opencode.F(a.Provider.ID + "/" + a.Model.ID)
		}
		_, err := a.Client.Session.Command(
			context.Background(),
			a.Session.ID,
			params,
		)
		if err != nil {
			slog.Error("Failed to execute command", "error", err)
			return toast.NewErrorToast(fmt.Sprintf("Failed to execute command: %v", err))()
		}
		return nil
	})

	// The actual response will come through SSE
	// For now, just return success
	return a, tea.Batch(cmds...)
}

func (a *App) SendShell(ctx context.Context, command string) (*App, tea.Cmd) {
	var cmds []tea.Cmd
	if a.Session.ID == "" {
		session, err := a.CreateSession(ctx)
		if err != nil {
			return a, toast.NewErrorToast(err.Error())
		}
		a.Session = session
		cmds = append(cmds, util.CmdHandler(SessionCreatedMsg{Session: session}))
	}

	cmds = append(cmds, func() tea.Msg {
		_, err := a.Client.Session.Shell(
			context.Background(),
			a.Session.ID,
			opencode.SessionShellParams{
				Agent:   opencode.F(a.Agent().Name),
				Command: opencode.F(command),
			},
		)
		if err != nil {
			slog.Error("Failed to submit shell command", "error", err)
			return toast.NewErrorToast(fmt.Sprintf("Failed to submit shell command: %v", err))()
		}
		return nil
	})

	// The actual response will come through SSE
	// For now, just return success
	return a, tea.Batch(cmds...)
}

func (a *App) Cancel(ctx context.Context, sessionID string) error {
	// Cancel any running compact operation
	if a.compactCancel != nil {
		a.compactCancel()
		a.compactCancel = nil
	}

	_, err := a.Client.Session.Abort(ctx, sessionID, opencode.SessionAbortParams{})
	if err != nil {
		slog.Error("Failed to cancel session", "error", err)
		return err
	}
	return nil
}

func (a *App) ListSessions(ctx context.Context) ([]opencode.Session, error) {
	response, err := a.Client.Session.List(ctx, opencode.SessionListParams{})
	if err != nil {
		return nil, err
	}
	if response == nil {
		return []opencode.Session{}, nil
	}
	sessions := *response
	return sessions, nil
}

func (a *App) DeleteSession(ctx context.Context, sessionID string) error {
	_, err := a.Client.Session.Delete(ctx, sessionID, opencode.SessionDeleteParams{})
	if err != nil {
		slog.Error("Failed to delete session", "error", err)
		return err
	}
	return nil
}

func (a *App) UpdateSession(ctx context.Context, sessionID string, title string) error {
	_, err := a.Client.Session.Update(ctx, sessionID, opencode.SessionUpdateParams{
		Title: opencode.F(title),
	})
	if err != nil {
		slog.Error("Failed to update session", "error", err)
		return err
	}
	return nil
}

// isFirstRun checks if this is the first run (no authenticated providers)
func (a *App) isFirstRun() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	status, err := a.AuthBridge.GetAuthStatus(ctx)
	if err != nil {
		// If we can't check, assume it's not first run to avoid disruption
		slog.Debug("Failed to check first run status", "error", err)
		return false
	}

	return len(status.Authenticated) == 0
}

// autoDetectAllCredentials attempts to auto-detect credentials on first run
func (a *App) autoDetectAllCredentials() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		slog.Debug("Auto-detecting credentials on first run")
		result, err := a.AuthBridge.AutoDetect(ctx)
		if err != nil {
			slog.Debug("Auto-detect failed on first run", "error", err)
			return nil // Silent fail on first run
		}

		if result.Found > 0 {
			slog.Info("Auto-detected credentials", "count", result.Found)
			return toast.NewSuccessToast(
				fmt.Sprintf("Found %d provider(s). Ready to code!", result.Found),
			)()
		}

		slog.Debug("No credentials auto-detected on first run")
		return nil
	}
}

// autoDetectAllCredentialsQuiet runs auto-detect silently on every startup
func (a *App) autoDetectAllCredentialsQuiet() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		slog.Debug("Auto-detecting credentials")
		result, err := a.AuthBridge.AutoDetect(ctx)
		if err != nil {
			slog.Debug("Auto-detect failed", "error", err)
			return nil // Silent fail
		}

		// Get authenticated providers status
		status, err := a.AuthBridge.GetAuthStatus(ctx)
		if err != nil {
			slog.Debug("Failed to get auth status", "error", err)
			// Still show basic message if we found credentials
			if result.Found > 0 {
				return toast.NewSuccessToast(
					fmt.Sprintf("Found %d provider(s). Ready to code!", result.Found),
				)()
			}
			return nil
		}

		if len(status.Authenticated) == 0 {
			// No providers authenticated
			return nil
		}

		// Build provider names list
		providerNames := make([]string, 0, len(status.Authenticated))
		for _, p := range status.Authenticated {
			providerNames = append(providerNames, p.Name)
		}

		// Generate friendly message
		var msg string
		if len(providerNames) == 1 {
			msg = fmt.Sprintf("Ready: %s ✓", providerNames[0])
		} else if len(providerNames) <= 3 {
			msg = fmt.Sprintf("Ready: %s ✓", strings.Join(providerNames, ", "))
		} else {
			// Show first 3 and count
			msg = fmt.Sprintf("Ready: %s, +%d more ✓",
				strings.Join(providerNames[:3], ", "),
				len(providerNames)-3)
		}

		// Check if all providers are authenticated
		totalProviders := len(a.Providers)
		if len(status.Authenticated) == totalProviders && totalProviders > 1 {
			msg = fmt.Sprintf("All providers ready: %s ✓", strings.Join(providerNames, ", "))
		}

		slog.Info("Authenticated providers ready", "count", len(status.Authenticated), "providers", providerNames)
		return toast.NewSuccessToast(msg)()
	}
}

// AnalyzePromptAndRecommendModel analyzes a prompt and recommends a better model if available
func (a *App) AnalyzePromptAndRecommendModel(prompt string) tea.Cmd {
	return func() tea.Msg {
		// Detect task type from prompt
		taskType := detectTaskType(prompt)
		if taskType == "general" {
			// Don't recommend for general tasks
			return nil
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		// Get AI recommendations for this task
		recommendations, err := a.AuthBridge.GetRecommendations(ctx, taskType)
		if err != nil {
			slog.Debug("Failed to get recommendations", "error", err)
			return nil
		}

		if len(recommendations) == 0 {
			return nil
		}

		// Check if current model is already optimal
		bestRec := recommendations[0]
		if a.Model != nil && a.Provider != nil {
			currentModelID := a.Provider.ID + "/" + a.Model.ID
			bestModelID := bestRec.Provider + "/" + bestRec.Model

			if currentModelID == bestModelID {
				// Already using the best model
				return nil
			}

			// Only recommend if confidence is high enough
			if bestRec.Score < 0.7 {
				return nil
			}

			// Find the recommended model in providers
			recommendedProvider, recommendedModel := findModelByProviderAndModelID(
				a.Providers,
				bestRec.Provider,
				bestRec.Model,
			)

			if recommendedProvider == nil || recommendedModel == nil {
				return nil
			}

			// Create toast with action to switch
			return toast.NewInfoToast(
				fmt.Sprintf("%s might be better for %s tasks", recommendedModel.Name, taskType),
			)()
		}

		return nil
	}
}

// detectTaskType detects the task type from a prompt
func detectTaskType(prompt string) string {
	lower := strings.ToLower(prompt)

	// Check for debugging/testing keywords
	if strings.Contains(lower, "test") || strings.Contains(lower, "bug") ||
		strings.Contains(lower, "debug") || strings.Contains(lower, "fix") ||
		strings.Contains(lower, "error") || strings.Contains(lower, "issue") {
		return "debugging"
	}

	// Check for refactoring keywords
	if strings.Contains(lower, "refactor") || strings.Contains(lower, "clean") ||
		strings.Contains(lower, "improve") || strings.Contains(lower, "optimize") {
		return "refactoring"
	}

	// Check for code generation keywords
	if strings.Contains(lower, "build") || strings.Contains(lower, "create") ||
		strings.Contains(lower, "implement") || strings.Contains(lower, "add") ||
		strings.Contains(lower, "write") || strings.Contains(lower, "generate") {
		return "code_generation"
	}

	// Check for code review keywords
	if strings.Contains(lower, "review") || strings.Contains(lower, "analyze") ||
		strings.Contains(lower, "explain") || strings.Contains(lower, "understand") {
		return "code_review"
	}

	// Check for quick questions
	if strings.Contains(lower, "quick") || strings.Contains(lower, "?") ||
		strings.Contains(lower, "how") || strings.Contains(lower, "what") ||
		strings.Contains(lower, "why") {
		return "quick_question"
	}

	return "general"
}

func (a *App) ListMessages(ctx context.Context, sessionId string) ([]Message, error) {
	response, err := a.Client.Session.Messages(ctx, sessionId, opencode.SessionMessagesParams{})
	if err != nil {
		return nil, err
	}
	if response == nil {
		return []Message{}, nil
	}
	messages := []Message{}
	for _, message := range *response {
		msg := Message{
			Info:  message.Info.AsUnion(),
			Parts: []opencode.PartUnion{},
		}
		for _, part := range message.Parts {
			msg.Parts = append(msg.Parts, part.AsUnion())
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

func (a *App) ListProviders(ctx context.Context) ([]opencode.Provider, error) {
	// Get providers from API
	response, err := a.Client.App.Providers(ctx, opencode.AppProvidersParams{})
	apiProviders := []opencode.Provider{}
	if err == nil && response != nil {
		providers := *response
		apiProviders = providers.Providers
	}

	// ALWAYS try to get CLI providers and merge them
	cliProviders, err := a.AuthBridge.GetCLIProviders(ctx)
	if err != nil {
		// CLI detection failed - just use API providers
		if len(apiProviders) > 0 {
			return apiProviders, nil
		}
		return nil, fmt.Errorf("failed to get providers: API returned %d, CLI error: %w", len(apiProviders), err)
	}

	// Convert CLI providers to opencode.Provider format
	providerMap := make(map[string]opencode.Provider)

	// Add API providers to map first
	for _, p := range apiProviders {
		providerMap[p.ID] = p
	}

	// Add or merge CLI providers
	for _, cliProv := range cliProviders {
		models := make(map[string]opencode.Model)
		for _, modelID := range cliProv.Models {
			models[modelID] = opencode.Model{
				ID:   modelID,
				Name: formatModelName(modelID),
			}
		}

		provider := opencode.Provider{
			ID:     cliProv.Provider,
			Name:   formatProviderName(cliProv.Provider),
			Models: models,
		}

		// If provider already exists from API, merge models
		if existing, exists := providerMap[cliProv.Provider]; exists {
			// Merge models from CLI into existing provider
			for modelID, model := range models {
				if _, hasModel := existing.Models[modelID]; !hasModel {
					existing.Models[modelID] = model
				}
			}
			providerMap[cliProv.Provider] = existing
		} else {
			// Add new CLI provider
			providerMap[cliProv.Provider] = provider
		}
	}

	// Convert map back to slice
	result := make([]opencode.Provider, 0, len(providerMap))
	for _, provider := range providerMap {
		result = append(result, provider)
	}

	return result, nil
}

// formatModelName formats a model ID into a human-readable name
func formatModelName(modelID string) string {
	// Simple formatting: replace hyphens with spaces and title case
	parts := strings.Split(modelID, "-")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, " ")
}

// formatProviderName formats a provider ID into a human-readable name
func formatProviderName(providerID string) string {
	names := map[string]string{
		"claude":    "Anthropic",
		"qwen":      "Alibaba",
		"codex":     "OpenAI",
		"gemini":    "Google",
		"anthropic": "Anthropic",
		"openai":    "OpenAI",
		"google":    "Google",
	}
	if name, ok := names[providerID]; ok {
		return name
	}
	return strings.ToUpper(providerID[:1]) + providerID[1:]
}

// func (a *App) loadCustomKeybinds() {
//
// }

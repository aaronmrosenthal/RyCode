package dialog

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/aaronmrosenthal/rycode-sdk-go"
	"github.com/aaronmrosenthal/rycode/internal/app"
	"github.com/aaronmrosenthal/rycode/internal/components/list"
	"github.com/aaronmrosenthal/rycode/internal/components/modal"
	"github.com/aaronmrosenthal/rycode/internal/components/toast"
	"github.com/aaronmrosenthal/rycode/internal/layout"
	"github.com/aaronmrosenthal/rycode/internal/styles"
	"github.com/aaronmrosenthal/rycode/internal/theme"
	"github.com/aaronmrosenthal/rycode/internal/util"
)

const (
	numVisibleModels = 10
	minDialogWidth   = 40
	maxDialogWidth   = 80
	maxRecentModels  = 5
)

// ModelDialog interface for the model selection dialog
type ModelDialog interface {
	layout.Modal
}

type modelDialog struct {
	app                   *app.App
	allModels             []ModelWithProvider
	width                 int
	height                int
	modal                 *modal.Modal
	searchDialog          *SearchDialog
	dialogWidth           int
	providerAuthStatus    map[string]*ProviderAuthStatus // Cached auth status per provider
	authPrompt            *AuthPromptDialog              // Auth prompt dialog
	showingAuthPrompt     bool                           // Whether auth prompt is visible
	authingProvider       string                         // Provider being authenticated
	recommendationPanel   *ModelRecommendationPanel      // AI recommendation panel
	showRecommendations   bool                           // Whether to show recommendations
}

// ProviderAuthStatus holds authentication and health information for a provider
type ProviderAuthStatus struct {
	IsAuthenticated bool
	Health          string // "healthy", "degraded", "down", "unknown"
	ModelsCount     int
	LastChecked     time.Time
}

type ModelWithProvider struct {
	Model    opencode.Model
	Provider opencode.Provider
}

// modelItem is a custom list item for model selections
type modelItem struct {
	model           ModelWithProvider
	isAuthenticated bool // Whether the provider is authenticated
}

func (m modelItem) Render(
	selected bool,
	width int,
	baseStyle styles.Style,
) string {
	t := theme.CurrentTheme()

	itemStyle := baseStyle.
		Background(t.BackgroundPanel()).
		Foreground(t.Text())

	// Gray out locked models
	if !m.isAuthenticated {
		itemStyle = itemStyle.Foreground(t.TextMuted()).Faint(true)
	} else if selected {
		itemStyle = itemStyle.Foreground(t.Primary())
	}

	providerStyle := baseStyle.
		Foreground(t.TextMuted()).
		Background(t.BackgroundPanel())

	modelPart := itemStyle.Render(m.model.Model.Name)
	providerPart := providerStyle.Render(fmt.Sprintf(" %s", m.model.Provider.Name))

	// Add lock indicator for unauthenticated providers
	lockPart := ""
	if !m.isAuthenticated {
		lockPart = providerStyle.Render(" [locked]")
	}

	combinedText := modelPart + providerPart + lockPart
	return baseStyle.
		Background(t.BackgroundPanel()).
		PaddingLeft(1).
		Render(combinedText)
}

func (m modelItem) Selectable() bool {
	// Only selectable if provider is authenticated
	return m.isAuthenticated
}

type modelKeyMap struct {
	Enter  key.Binding
	Escape key.Binding
}

var modelKeys = modelKeyMap{
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select model"),
	),
	Escape: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "close"),
	),
}

func (m *modelDialog) Init() tea.Cmd {
	m.setupAllModels()

	// Generate initial recommendations
	if m.recommendationPanel != nil {
		ctx := GetDefaultContext()
		m.recommendationPanel.GenerateRecommendations(ctx)
	}

	return m.searchDialog.Init()
}

func (m *modelDialog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// If showing auth prompt, route messages there
	if m.showingAuthPrompt {
		return m.handleAuthPromptUpdate(msg)
	}

	switch msg := msg.(type) {
	case SearchSelectionMsg:
		// Handle selection from search dialog
		if item, ok := msg.Item.(modelItem); ok {
			// If model is locked, show auth prompt
			if !item.isAuthenticated {
				m.showAuthPrompt(item.model.Provider.ID, item.model.Provider.Name)
				return m, nil
			}

			return m, tea.Sequence(
				util.CmdHandler(modal.CloseModalMsg{}),
				util.CmdHandler(
					app.ModelSelectedMsg{
						Provider: item.model.Provider,
						Model:    item.model.Model,
					}),
			)
		}
		return m, util.CmdHandler(modal.CloseModalMsg{})
	case SearchCancelledMsg:
		return m, util.CmdHandler(modal.CloseModalMsg{})

	case SearchRemoveItemMsg:
		if item, ok := msg.Item.(modelItem); ok {
			if m.isModelInRecentSection(item.model, msg.Index) {
				m.app.State.RemoveModelFromRecentlyUsed(item.model.Provider.ID, item.model.Model.ID)
				items := m.buildDisplayList(m.searchDialog.GetQuery())
				m.searchDialog.SetItems(items)
				return m, m.app.SaveState()
			}
		}
		return m, nil

	case SearchQueryChangedMsg:
		// Update the list based on search query
		items := m.buildDisplayList(msg.Query)
		m.searchDialog.SetItems(items)
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.searchDialog.SetWidth(m.dialogWidth)
		m.searchDialog.SetHeight(msg.Height)
		if m.authPrompt != nil {
			m.authPrompt.SetSize(msg.Width, msg.Height)
		}
		if m.recommendationPanel != nil {
			m.recommendationPanel.SetWidth(m.dialogWidth - 4)
		}

	case tea.KeyPressMsg:
		// Handle 'a' key to start authentication on focused provider
		if msg.String() == "a" {
			providerID, providerName := m.getFocusedProvider()
			if providerID != "" {
				authStatus := m.checkProviderAuth(providerID)
				if !authStatus.IsAuthenticated {
					m.showAuthPrompt(providerID, providerName)
					return m, nil
				}
			}
		}

		// Handle 'd' key for auto-detect
		if msg.String() == "d" {
			return m, m.performAutoDetect()
		}

		// Handle 'i' key to toggle recommendations (AI insights)
		if msg.String() == "i" {
			m.showRecommendations = !m.showRecommendations
			if m.recommendationPanel != nil {
				m.recommendationPanel.SetVisible(m.showRecommendations)
			}
			return m, nil
		}

	case AuthSuccessMsg:
		// Authentication succeeded
		m.showingAuthPrompt = false
		m.authPrompt = nil

		// Invalidate cache for this provider
		delete(m.providerAuthStatus, msg.Provider)

		// Refresh display
		items := m.buildDisplayList(m.searchDialog.GetQuery())
		m.searchDialog.SetItems(items)

		// Show success toast
		return m, toast.NewSuccessToast(
			fmt.Sprintf("✓ Authenticated with %s (%d models)", msg.Provider, msg.ModelsCount),
		)

	case AuthFailureMsg:
		// Authentication failed - show error in prompt
		if m.authPrompt != nil {
			m.authPrompt.SetError(msg.Error)
		}
		return m, nil

	case AuthStatusRefreshMsg:
		// Refresh auth status for all providers
		m.providerAuthStatus = make(map[string]*ProviderAuthStatus)
		items := m.buildDisplayList(m.searchDialog.GetQuery())
		m.searchDialog.SetItems(items)
		return m, nil
	}

	updatedDialog, cmd := m.searchDialog.Update(msg)
	m.searchDialog = updatedDialog.(*SearchDialog)
	return m, cmd
}

func (m *modelDialog) View() string {
	if m.showingAuthPrompt && m.authPrompt != nil {
		return m.authPrompt.View()
	}

	// Show model list
	listView := m.searchDialog.View()

	// Add recommendations if enabled and no search query
	if m.showRecommendations && m.recommendationPanel != nil && m.searchDialog.GetQuery() == "" {
		recommendationsView := m.recommendationPanel.View()
		if recommendationsView != "" {
			// Stack vertically
			return listView + "\n\n" + recommendationsView
		}
	}

	return listView
}

// showAuthPrompt displays the authentication prompt for a provider
func (m *modelDialog) showAuthPrompt(providerID, providerName string) {
	m.authPrompt = NewAuthPromptDialog(providerName)
	m.authPrompt.SetSize(m.width, m.height)
	m.showingAuthPrompt = true
	m.authingProvider = providerID
}

// handleAuthPromptUpdate handles messages when auth prompt is visible
func (m *modelDialog) handleAuthPromptUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			// Submit authentication
			apiKey := m.authPrompt.GetValue()
			if apiKey != "" {
				return m, m.performAuthentication(m.authingProvider, apiKey)
			}
			return m, nil

		case "ctrl+d":
			// Auto-detect credentials
			m.showingAuthPrompt = false
			m.authPrompt = nil
			return m, m.performAutoDetect()

		case "esc":
			// Cancel authentication
			m.showingAuthPrompt = false
			m.authPrompt = nil
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.authPrompt != nil {
			m.authPrompt.SetSize(msg.Width, msg.Height)
		}
		return m, nil
	}

	// Pass message to auth prompt
	var cmd tea.Cmd
	m.authPrompt, cmd = m.authPrompt.Update(msg)
	return m, cmd
}

// performAuthentication authenticates with a provider using an API key
func (m *modelDialog) performAuthentication(providerID, apiKey string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		result, err := m.app.AuthBridge.Authenticate(ctx, providerID, apiKey)
		if err != nil {
			return AuthFailureMsg{
				Provider: providerID,
				Error:    err.Error(),
			}
		}

		return AuthSuccessMsg{
			Provider:    result.Provider,
			ModelsCount: result.ModelsCount,
		}
	}
}

// performAutoDetect attempts to auto-detect credentials
func (m *modelDialog) performAutoDetect() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		result, err := m.app.AuthBridge.AutoDetect(ctx)
		if err != nil {
			return toast.NewErrorToast("Auto-detect failed: " + err.Error())()
		}

		if result.Found == 0 {
			return toast.NewInfoToast("No credentials found. Please enter manually.")()
		}

		// Refresh auth status after auto-detect
		return tea.Batch(
			util.CmdHandler(AuthStatusRefreshMsg{}),
			toast.NewSuccessToast(fmt.Sprintf("✓ Auto-detected %d credential(s)", result.Found)),
		)
	}
}

// getFocusedProvider returns the provider ID and name of the currently focused item
func (m *modelDialog) getFocusedProvider() (string, string) {
	// Get the selected item from search dialog's list
	selectedItem, selectedIndex := m.searchDialog.list.GetSelectedItem()

	if selectedIndex == -1 {
		return "", ""
	}

	// Check if it's a model item
	if item, ok := selectedItem.(modelItem); ok {
		return item.model.Provider.ID, item.model.Provider.Name
	}

	// If it's a header, try to get the next item
	items := m.buildDisplayList(m.searchDialog.GetQuery())
	if selectedIndex+1 < len(items) {
		if item, ok := items[selectedIndex+1].(modelItem); ok {
			return item.model.Provider.ID, item.model.Provider.Name
		}
	}

	return "", ""
}

func (m *modelDialog) calculateOptimalWidth(models []ModelWithProvider) int {
	maxWidth := minDialogWidth

	for _, model := range models {
		// Calculate the width needed for this item: "ModelName (ProviderName)"
		// Add 4 for the parentheses, space, and some padding
		itemWidth := len(model.Model.Name) + len(model.Provider.Name) + 4
		if itemWidth > maxWidth {
			maxWidth = itemWidth
		}
	}

	if maxWidth > maxDialogWidth {
		maxWidth = maxDialogWidth
	}

	return maxWidth
}

// buildProviderHeader creates a header string with authentication status indicators
func (m *modelDialog) buildProviderHeader(providerName string, status *ProviderAuthStatus) string {
	if status.IsAuthenticated {
		// Authenticated provider - check health
		switch status.Health {
		case "healthy":
			return providerName + " ✓"
		case "degraded":
			return providerName + " ⚠"
		case "down":
			return providerName + " ✗"
		default:
			return providerName + " ✓"
		}
	} else {
		// Not authenticated
		return providerName + " 🔒"
	}
}

// checkProviderAuth checks and caches the authentication status for a provider
func (m *modelDialog) checkProviderAuth(providerID string) *ProviderAuthStatus {
	// Return cached status if available and fresh (< 30 seconds old)
	if cached, ok := m.providerAuthStatus[providerID]; ok {
		if time.Since(cached.LastChecked) < 30*time.Second {
			return cached
		}
	}

	// Create context with timeout to avoid blocking UI
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	status := &ProviderAuthStatus{
		IsAuthenticated: false,
		Health:          "unknown",
		ModelsCount:     0,
		LastChecked:     time.Now(),
	}

	// Check authentication status via bridge
	authStatus, err := m.app.AuthBridge.CheckAuthStatus(ctx, providerID)
	if err == nil {
		status.IsAuthenticated = authStatus.IsAuthenticated
		status.ModelsCount = authStatus.ModelsCount
	}

	// Check provider health
	health, err := m.app.AuthBridge.GetProviderHealth(ctx, providerID)
	if err == nil {
		status.Health = health.Status
	}

	// Cache the result
	m.providerAuthStatus[providerID] = status
	return status
}

func (m *modelDialog) setupAllModels() {
	providers, _ := m.app.ListProviders(context.Background())

	m.allModels = make([]ModelWithProvider, 0)
	for _, provider := range providers {
		for _, model := range provider.Models {
			m.allModels = append(m.allModels, ModelWithProvider{
				Model:    model,
				Provider: provider,
			})
		}
	}

	m.sortModels()

	// Calculate optimal width based on all models
	m.dialogWidth = m.calculateOptimalWidth(m.allModels)

	// Initialize search dialog
	m.searchDialog = NewSearchDialog("Search models...", numVisibleModels)
	m.searchDialog.SetWidth(m.dialogWidth)

	// Build initial display list (empty query shows grouped view)
	items := m.buildDisplayList("")
	m.searchDialog.SetItems(items)
}

func (m *modelDialog) sortModels() {
	sort.Slice(m.allModels, func(i, j int) bool {
		modelA := m.allModels[i]
		modelB := m.allModels[j]

		usageA := m.getModelUsageTime(modelA.Provider.ID, modelA.Model.ID)
		usageB := m.getModelUsageTime(modelB.Provider.ID, modelB.Model.ID)

		// If both have usage times, sort by most recent first
		if !usageA.IsZero() && !usageB.IsZero() {
			return usageA.After(usageB)
		}

		// If only one has usage time, it goes first
		if !usageA.IsZero() && usageB.IsZero() {
			return true
		}
		if usageA.IsZero() && !usageB.IsZero() {
			return false
		}

		// If neither has usage time, sort by release date desc if available
		if modelA.Model.ReleaseDate != "" && modelB.Model.ReleaseDate != "" {
			dateA := m.parseReleaseDate(modelA.Model.ReleaseDate)
			dateB := m.parseReleaseDate(modelB.Model.ReleaseDate)
			if !dateA.IsZero() && !dateB.IsZero() {
				return dateA.After(dateB)
			}
		}

		// If only one has release date, it goes first
		if modelA.Model.ReleaseDate != "" && modelB.Model.ReleaseDate == "" {
			return true
		}
		if modelA.Model.ReleaseDate == "" && modelB.Model.ReleaseDate != "" {
			return false
		}

		// If neither has usage time nor release date, fall back to alphabetical sorting
		return modelA.Model.Name < modelB.Model.Name
	})
}

func (m *modelDialog) parseReleaseDate(dateStr string) time.Time {
	if parsed, err := time.Parse("2006-01-02", dateStr); err == nil {
		return parsed
	}

	return time.Time{}
}

func (m *modelDialog) getModelUsageTime(providerID, modelID string) time.Time {
	for _, usage := range m.app.State.RecentlyUsedModels {
		if usage.ProviderID == providerID && usage.ModelID == modelID {
			return usage.LastUsed
		}
	}
	return time.Time{}
}

// buildDisplayList creates the list items based on search query
func (m *modelDialog) buildDisplayList(query string) []list.Item {
	if query != "" {
		// Search mode: use fuzzy matching
		return m.buildSearchResults(query)
	} else {
		// Grouped mode: show Recent section and provider groups
		return m.buildGroupedResults()
	}
}

// buildSearchResults creates a flat list of search results using fuzzy matching
func (m *modelDialog) buildSearchResults(query string) []list.Item {
	type modelMatch struct {
		model ModelWithProvider
		score int
	}

	modelNames := []string{}
	modelMap := make(map[string]ModelWithProvider)

	// Create search strings and perform fuzzy matching
	for _, model := range m.allModels {
		searchStr := fmt.Sprintf("%s %s", model.Model.Name, model.Provider.Name)
		modelNames = append(modelNames, searchStr)
		modelMap[searchStr] = model

		searchStr = fmt.Sprintf("%s %s", model.Provider.Name, model.Model.Name)
		modelNames = append(modelNames, searchStr)
		modelMap[searchStr] = model
	}

	matches := fuzzy.RankFindFold(query, modelNames)
	sort.Sort(matches)

	items := []list.Item{}
	seenModels := make(map[string]bool)

	for _, match := range matches {
		model := modelMap[match.Target]
		// Create a unique key to avoid duplicates
		// Include name to handle custom models with same ID but different names
		key := fmt.Sprintf("%s:%s:%s", model.Provider.ID, model.Model.ID, model.Model.Name)
		if seenModels[key] {
			continue
		}
		seenModels[key] = true

		// Check auth status for search results
		authStatus := m.checkProviderAuth(model.Provider.ID)
		items = append(items, modelItem{
			model:           model,
			isAuthenticated: authStatus.IsAuthenticated,
		})
	}

	return items
}

// buildGroupedResults creates a grouped list with Recent section and provider groups
func (m *modelDialog) buildGroupedResults() []list.Item {
	var items []list.Item

	// Add Recent section
	recentModels := m.getRecentModels(maxRecentModels)
	if len(recentModels) > 0 {
		items = append(items, list.HeaderItem("Recent"))
		for _, model := range recentModels {
			// Check auth status for recent models
			authStatus := m.checkProviderAuth(model.Provider.ID)
			items = append(items, modelItem{
				model:           model,
				isAuthenticated: authStatus.IsAuthenticated,
			})
		}
	}

	// Group models by provider
	providerGroups := make(map[string][]ModelWithProvider)
	for _, model := range m.allModels {
		providerName := model.Provider.Name
		providerGroups[providerName] = append(providerGroups[providerName], model)
	}

	// Get sorted provider names for consistent order
	var providerNames []string
	for name := range providerGroups {
		providerNames = append(providerNames, name)
	}
	sort.Strings(providerNames)

	// Add provider groups
	for _, providerName := range providerNames {
		models := providerGroups[providerName]
		providerID := models[0].Provider.ID

		// Check auth status for this provider
		authStatus := m.checkProviderAuth(providerID)

		// Sort models within provider group
		sort.Slice(models, func(i, j int) bool {
			modelA := models[i]
			modelB := models[j]

			usageA := m.getModelUsageTime(modelA.Provider.ID, modelA.Model.ID)
			usageB := m.getModelUsageTime(modelB.Provider.ID, modelB.Model.ID)

			// Sort by usage time first, then by release date, then alphabetically
			if !usageA.IsZero() && !usageB.IsZero() {
				return usageA.After(usageB)
			}
			if !usageA.IsZero() && usageB.IsZero() {
				return true
			}
			if usageA.IsZero() && !usageB.IsZero() {
				return false
			}

			// Sort by release date if available
			if modelA.Model.ReleaseDate != "" && modelB.Model.ReleaseDate != "" {
				dateA := m.parseReleaseDate(modelA.Model.ReleaseDate)
				dateB := m.parseReleaseDate(modelB.Model.ReleaseDate)
				if !dateA.IsZero() && !dateB.IsZero() {
					return dateA.After(dateB)
				}
			}

			return modelA.Model.Name < modelB.Model.Name
		})

		// Add provider header with auth status indicator
		headerText := m.buildProviderHeader(providerName, authStatus)
		items = append(items, list.HeaderItem(headerText))

		// Add models in this provider group
		for _, model := range models {
			items = append(items, modelItem{
				model:           model,
				isAuthenticated: authStatus.IsAuthenticated,
			})
		}
	}

	return items
}

// getRecentModels returns the most recently used models
func (m *modelDialog) getRecentModels(limit int) []ModelWithProvider {
	var recentModels []ModelWithProvider

	// Get recent models from app state
	for _, usage := range m.app.State.RecentlyUsedModels {
		if len(recentModels) >= limit {
			break
		}

		// Find the corresponding model
		for _, model := range m.allModels {
			if model.Provider.ID == usage.ProviderID && model.Model.ID == usage.ModelID {
				recentModels = append(recentModels, model)
				break
			}
		}
	}

	return recentModels
}

func (m *modelDialog) isModelInRecentSection(model ModelWithProvider, index int) bool {
	// Only check if we're in grouped mode (no search query)
	if m.searchDialog.GetQuery() != "" {
		return false
	}

	recentModels := m.getRecentModels(maxRecentModels)
	if len(recentModels) == 0 {
		return false
	}

	// Index 0 is the "Recent" header, so recent models are at indices 1 to len(recentModels)
	if index >= 1 && index <= len(recentModels) {
		if index-1 < len(recentModels) {
			recentModel := recentModels[index-1]
			return recentModel.Provider.ID == model.Provider.ID &&
				recentModel.Model.ID == model.Model.ID
		}
	}

	return false
}

func (m *modelDialog) Render(background string) string {
	return m.modal.Render(m.View(), background)
}

func (s *modelDialog) Close() tea.Cmd {
	return nil
}

func NewModelDialog(app *app.App) ModelDialog {
	dialog := &modelDialog{
		app:                 app,
		providerAuthStatus:  make(map[string]*ProviderAuthStatus),
		recommendationPanel: NewModelRecommendationPanel(),
		showRecommendations: true, // Show by default
	}

	dialog.setupAllModels()

	dialog.modal = modal.New(
		modal.WithTitle("Select Model"),
		modal.WithMaxWidth(dialog.dialogWidth+4),
	)

	return dialog
}

package dialog

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/lipgloss/v2/compat"

	"github.com/aaronmrosenthal/rycode-sdk-go"
	"github.com/aaronmrosenthal/rycode/internal/app"
	"github.com/aaronmrosenthal/rycode/internal/auth"
	"github.com/aaronmrosenthal/rycode/internal/components/modal"
	"github.com/aaronmrosenthal/rycode/internal/components/splash"
	"github.com/aaronmrosenthal/rycode/internal/theme"
	"github.com/aaronmrosenthal/rycode/internal/util"
)

// providersLoadedMsg is sent when providers finish loading
type providersLoadedMsg struct {
	providers []opencode.Provider
	err       error
}

// loadingTickMsg is sent for loading animation updates
type loadingTickMsg time.Time

// SimpleProviderToggle is a minimal provider selector that just cycles between CLI providers
type SimpleProviderToggle struct {
	app             *app.App
	providers       []opencode.Provider
	selectedIndex   int
	width           int
	height          int
	isLoading       bool
	loadError       error
	cortexRenderer  *splash.CortexRenderer
	isSwitching     bool      // True when switching providers (show cortex)
	switchStartTime time.Time // When the switch animation started
	fadeOpacity     float64   // Current fade opacity (0.0 to 1.0)
}

// NewSimpleProviderToggle creates a new simple provider toggle
func NewSimpleProviderToggle(app *app.App) *SimpleProviderToggle {
	// Create a compact 3D torus for loading animation (40x12)
	cortexRenderer := splash.NewCortexRenderer(40, 12)

	return &SimpleProviderToggle{
		app:            app,
		providers:      []opencode.Provider{},
		selectedIndex:  0,
		isLoading:      true,
		cortexRenderer: cortexRenderer,
	}
}

func (s *SimpleProviderToggle) Init() tea.Cmd {
	// Start loading animation and load providers asynchronously
	return tea.Batch(
		s.loadProvidersAsync(),
		s.tickLoadingAnimation(),
	)
}

func (s *SimpleProviderToggle) tickLoadingAnimation() tea.Cmd {
	return tea.Tick(splash.CortexAnimationTickInterval, func(t time.Time) tea.Msg {
		return loadingTickMsg(t)
	})
}

func (s *SimpleProviderToggle) loadProvidersAsync() tea.Cmd {
	return func() tea.Msg {
		providers, err := s.loadAuthenticatedProvidersSync()
		return providersLoadedMsg{providers: providers, err: err}
	}
}

func (s *SimpleProviderToggle) loadAuthenticatedProvidersSync() ([]opencode.Provider, error) {
	// Use timeout context to prevent hanging
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	slog.Debug("loading authenticated providers")

	// Get CLI providers directly from AuthBridge (doesn't require full client)
	cliProviders, err := s.app.AuthBridge.GetCLIProviders(ctx)
	if err != nil {
		slog.Error("failed to load CLI providers", "error", err)
		return nil, fmt.Errorf("failed to load CLI providers: %w", err)
	}

	slog.Debug("loaded CLI providers", "count", len(cliProviders))
	for i, p := range cliProviders {
		slog.Debug("CLI provider details", "index", i, "id", p.Provider, "models", len(p.Models))
	}

	// Check authentication in parallel for faster loading (reduces from ~4-5s to ~1s)
	type authCheckResult struct {
		providerID  string
		models      []string
		authStatus  *auth.AuthStatus
		err         error
	}

	resultChan := make(chan authCheckResult, len(cliProviders))
	for _, cliProv := range cliProviders {
		go func(providerID string, models []string) {
			slog.Debug("checking provider auth", "provider", providerID)
			authStatus, err := s.app.AuthBridge.CheckAuthStatus(ctx, providerID)

			resultChan <- authCheckResult{
				providerID: providerID,
				models:     models,
				authStatus: authStatus,
				err:        err,
			}
		}(cliProv.Provider, cliProv.Models)
	}

	// Collect results
	providers := []opencode.Provider{}
	for range cliProviders {
		result := <-resultChan

		if result.err != nil {
			slog.Warn("auth check failed", "provider", result.providerID, "error", result.err)
			continue
		}

		slog.Debug("auth status checked",
			"provider", result.providerID,
			"authenticated", result.authStatus.IsAuthenticated,
			"models_count", result.authStatus.ModelsCount)

		if !result.authStatus.IsAuthenticated {
			slog.Debug("skipping unauthenticated provider", "provider", result.providerID)
			continue
		}

		// Convert to opencode.Provider format
		models := make(map[string]opencode.Model)
		for _, modelID := range result.models {
			models[modelID] = opencode.Model{
				ID:   modelID,
				Name: formatModelName(modelID),
			}
		}

		provider := opencode.Provider{
			ID:     result.providerID,
			Name:   formatProviderName(result.providerID),
			Models: models,
		}

		slog.Debug("added provider to list", "provider", result.providerID, "models_count", len(models))
		providers = append(providers, provider)
	}

	slog.Debug("providers loading complete", "total_providers", len(providers))
	for i, p := range providers {
		slog.Debug("final provider", "index", i, "id", p.ID, "name", p.Name, "models_count", len(p.Models))
	}

	// Sort providers by priority: Claude, Codex, Gemini, Grok, Qwen
	providerOrder := map[string]int{
		"claude":  1,
		"codex":   2,
		"gemini":  3,
		"grok":    4,
		"qwen":    5,
	}
	sort.Slice(providers, func(i, j int) bool {
		orderI, okI := providerOrder[providers[i].ID]
		orderJ, okJ := providerOrder[providers[j].ID]
		if okI && okJ {
			return orderI < orderJ
		}
		if okI {
			return true
		}
		if okJ {
			return false
		}
		return providers[i].ID < providers[j].ID
	})

	return providers, nil
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

// formatProviderName formats a provider ID into a human-readable company name
func formatProviderName(providerID string) string {
	names := map[string]string{
		"anthropic": "Anthropic",
		"claude":    "Anthropic",
		"openai":    "OpenAI",
		"codex":     "OpenAI",
		"google":    "Google",
		"gemini":    "Google",
		"xai":       "xAI",
		"grok":      "xAI",
		"qwen":      "Alibaba",
	}
	if name, ok := names[providerID]; ok {
		return name
	}
	return strings.ToUpper(providerID[:1]) + providerID[1:]
}

// getProviderDisplayName returns the short display name for chips/UI
func getProviderDisplayName(providerID string) string {
	names := map[string]string{
		"anthropic": "Claude",
		"claude":    "Claude",
		"google":    "Gemini",
		"gemini":    "Gemini",
		"openai":    "Codex",
		"codex":     "Codex",
		"xai":       "Grok",
		"grok":      "Grok",
		"qwen":      "Qwen",
	}
	if name, ok := names[providerID]; ok {
		return name
	}
	return providerID
}

func (s *SimpleProviderToggle) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case providersLoadedMsg:
		s.isLoading = false
		if msg.err != nil {
			s.loadError = msg.err
			return s, nil
		}
		s.providers = msg.providers
		s.loadError = nil

		// Log dynamic width calculation for visibility
		requiredWidth := s.CalculateRequiredWidth()
		slog.Debug("providers loaded, calculated required width",
			"providers", len(s.providers),
			"width", requiredWidth)

		return s, nil

	case loadingTickMsg:
		if s.isLoading {
			// Torus animation advances automatically in Render()
			return s, s.tickLoadingAnimation()
		}

		// Fade animation when switching providers
		if s.isSwitching {
			elapsed := time.Since(s.switchStartTime)
			opacity, finished := splash.CalculateAnimationOpacity(elapsed, 1.0)

			if finished {
				// Animation complete
				s.isSwitching = false
				s.fadeOpacity = 0.0
				slog.Debug("modal provider switch animation finished", "duration", elapsed)
				return s, nil
			}

			s.fadeOpacity = opacity
			return s, s.tickLoadingAnimation()
		}

		return s, nil

	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		return s, nil

	case tea.KeyPressMsg:
		slog.Debug("modal key press", "key", msg.String(), "selected_index", s.selectedIndex, "providers_count", len(s.providers))
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("tab"))):
			// Tab: cycle to next provider with cortex fade animation
			if len(s.providers) > 0 {
				oldIndex := s.selectedIndex
				s.selectedIndex = (s.selectedIndex + 1) % len(s.providers)
				selectedProvider := s.providers[s.selectedIndex]
				slog.Debug("modal tab cycle", "from", oldIndex, "to", s.selectedIndex, "provider", selectedProvider.ID)

				// Set cortex to provider's brand color
				brandColor := splash.GetProviderBrandColor(selectedProvider.ID)
				s.cortexRenderer.SetBrandColor(brandColor)
				slog.Debug("cortex brand color set",
					"provider", selectedProvider.ID,
					"color", fmt.Sprintf("#%02X%02X%02X", brandColor.R, brandColor.G, brandColor.B))

				// Start cortex fade animation (1.2s total: fade-in, hold, fade-out)
				s.isSwitching = true
				s.switchStartTime = time.Now()
				s.fadeOpacity = 1.0 // Start at FULL visibility for instant feedback
				slog.Debug("modal switch animation started", "opacity", s.fadeOpacity)
				return s, s.tickLoadingAnimation()
			} else {
				slog.Debug("tab ignored - no providers loaded")
			}
			// Consume the Tab key even if providers aren't loaded
			// This prevents the main TUI from cycling providers in the background
			return s, nil

		case key.Matches(msg, key.NewBinding(key.WithKeys("shift+tab"))):
			// Shift+Tab: cycle to previous provider with cortex fade animation
			if len(s.providers) > 0 {
				oldIndex := s.selectedIndex
				s.selectedIndex--
				if s.selectedIndex < 0 {
					s.selectedIndex = len(s.providers) - 1
				}
				selectedProvider := s.providers[s.selectedIndex]
				slog.Debug("modal shift+tab cycle", "from", oldIndex, "to", s.selectedIndex, "provider", selectedProvider.ID)

				// Set cortex to provider's brand color
				brandColor := splash.GetProviderBrandColor(selectedProvider.ID)
				s.cortexRenderer.SetBrandColor(brandColor)
				slog.Debug("cortex brand color set",
					"provider", selectedProvider.ID,
					"color", fmt.Sprintf("#%02X%02X%02X", brandColor.R, brandColor.G, brandColor.B))

				// Start cortex fade animation (1.2s total: fade-in, hold, fade-out)
				s.isSwitching = true
				s.switchStartTime = time.Now()
				s.fadeOpacity = 1.0 // Start at FULL visibility for instant feedback
				slog.Debug("modal switch animation started", "opacity", s.fadeOpacity)
				return s, s.tickLoadingAnimation()
			} else {
				slog.Debug("shift+tab ignored - no providers loaded")
			}
			// Consume the Shift+Tab key even if providers aren't loaded
			// This prevents the main TUI from cycling providers in the background
			return s, nil

		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			// Enter: select current provider
			if len(s.providers) > 0 && s.selectedIndex < len(s.providers) {
				provider := s.providers[s.selectedIndex]

				// Get the best/default model for this provider
				selectedModel := s.getDefaultModelForProvider(provider)

				return s, tea.Sequence(
					util.CmdHandler(modal.CloseModalMsg{}),
					util.CmdHandler(app.ModelSelectedMsg{
						Provider: provider,
						Model:    selectedModel,
					}),
				)
			}
			return s, nil

		case key.Matches(msg, key.NewBinding(key.WithKeys("esc"))):
			// Esc: close without selecting
			slog.Debug("modal closing via esc key")
			return s, util.CmdHandler(modal.CloseModalMsg{})

		case key.Matches(msg, key.NewBinding(key.WithKeys("r"))):
			// R: retry loading if there was an error
			if s.loadError != nil {
				s.isLoading = true
				s.loadError = nil
				return s, tea.Batch(
					s.loadProvidersAsync(),
					s.tickLoadingAnimation(),
				)
			}
			return s, nil

		case key.Matches(msg, key.NewBinding(key.WithKeys("1", "2", "3", "4", "5"))):
			// Number keys: direct provider selection
			index := int(msg.String()[0] - '1') // Convert '1'-'5' to 0-4
			if index >= 0 && index < len(s.providers) {
				provider := s.providers[index]
				selectedModel := s.getDefaultModelForProvider(provider)

				slog.Debug("modal direct selection", "index", index, "provider", provider.ID)
				return s, tea.Sequence(
					util.CmdHandler(modal.CloseModalMsg{}),
					util.CmdHandler(app.ModelSelectedMsg{
						Provider: provider,
						Model:    selectedModel,
					}),
				)
			}
			return s, nil

		default:
			slog.Debug("modal key not handled", "key", msg.String())
		}
	}

	return s, nil
}

func (s *SimpleProviderToggle) View() string {
	t := theme.CurrentTheme()

	if s.isLoading {
		return s.renderLoading(t)
	}

	// Show cortex when switching providers
	if s.isSwitching {
		return s.renderSwitching(t)
	}

	if s.loadError != nil {
		return s.renderError(t)
	}

	if len(s.providers) == 0 {
		return s.renderEmpty(t)
	}

	var b strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().
		Foreground(t.Primary()).
		Bold(true).
		Padding(1, 0)
	b.WriteString(titleStyle.Render("Select Provider"))
	b.WriteString("\n\n")

	// Provider chips (vertical stacked layout for reliability)
	for i, provider := range s.providers {
		isSelected := i == s.selectedIndex
		chip := s.renderProviderChip(provider, isSelected, i+1, t)
		chipContainerStyle := lipgloss.NewStyle().Padding(0, 2)
		b.WriteString(chipContainerStyle.Render(chip))
		b.WriteString("\n")
	}
	b.WriteString("\n")

	// Show models for selected provider
	if s.selectedIndex < len(s.providers) {
		selectedProvider := s.providers[s.selectedIndex]
		b.WriteString(s.renderProviderModels(selectedProvider, t))
	}

	b.WriteString("\n\n")

	// Footer with shortcuts
	footerStyle := lipgloss.NewStyle().
		Foreground(t.TextMuted()).
		Padding(1, 2)
	footer := "1-5: Quick Select | Tab: Next | Shift+Tab: Previous | Enter: Select | Esc: Cancel"
	b.WriteString(footerStyle.Render(footer))

	return b.String()
}

func (s *SimpleProviderToggle) renderLoading(t theme.Theme) string {
	textStyle := lipgloss.NewStyle().
		Foreground(t.TextMuted()).
		Align(lipgloss.Center)

	var builder strings.Builder
	builder.WriteString("\n")

	// Render the 3D spinning torus
	torusOutput := s.cortexRenderer.Render()

	// Center the torus (use default width if not set yet)
	width := s.width
	if width == 0 {
		width = 80 // Default terminal width
	}
	torusStyle := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(width)
	builder.WriteString(torusStyle.Render(torusOutput))

	builder.WriteString("\n")
	builder.WriteString(textStyle.Render("Loading providers..."))
	builder.WriteString("\n")

	return builder.String()
}

func (s *SimpleProviderToggle) renderSwitching(t theme.Theme) string {
	// Apply opacity to the cortex and text based on fade state
	// Opacity ranges from 0.0 (transparent) to 1.0 (full opacity)
	opacity := s.fadeOpacity

	// Calculate color with opacity (blend with background)
	// Simple opacity simulation by reducing brightness (terminals don't support true alpha)
	var textColor compat.AdaptiveColor
	if opacity >= 0.9 {
		textColor = t.Text()
	} else if opacity >= 0.7 {
		textColor = t.TextMuted()
	} else if opacity >= 0.4 {
		textColor = compat.AdaptiveColor{
			Light: lipgloss.Color("#666666"),
			Dark:  lipgloss.Color("#444444"),
		}
	} else if opacity >= 0.2 {
		textColor = compat.AdaptiveColor{
			Light: lipgloss.Color("#333333"),
			Dark:  lipgloss.Color("#222222"),
		}
	} else {
		textColor = compat.AdaptiveColor{
			Light: lipgloss.Color("#222222"),
			Dark:  lipgloss.Color("#111111"),
		}
	}

	textStyle := lipgloss.NewStyle().
		Foreground(textColor).
		Align(lipgloss.Center)

	var builder strings.Builder
	builder.WriteString("\n")

	// Only render cortex if opacity > 0.1 (optimization)
	if opacity > 0.1 {
		// Render cortex frame (updates screen buffer)
		s.cortexRenderer.RenderFrame()

		// Manually build the cortex output with opacity applied to colors
		var cortexBuilder strings.Builder
		cortexWidth := s.cortexRenderer.Width()
		cortexHeight := s.cortexRenderer.Height()
		for y := 0; y < cortexHeight; y++ {
			for x := 0; x < cortexWidth; x++ {
				idx := y*cortexWidth + x
				char := s.cortexRenderer.Screen(idx)

				if char != ' ' {
					// Get color with opacity applied
					rgb := s.cortexRenderer.GetColorAtWithOpacity(x, y, opacity)
					cortexBuilder.WriteString(splash.Colorize(string(char), rgb))
				} else {
					cortexBuilder.WriteRune(' ')
				}
			}
			if y < cortexHeight-1 {
				cortexBuilder.WriteRune('\n')
			}
		}

		// Center the torus (use default width if not set yet)
		width := s.width
		if width == 0 {
			width = 80 // Default terminal width
		}
		torusStyle := lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(width)
		builder.WriteString(torusStyle.Render(cortexBuilder.String()))

		builder.WriteString("\n")

		// Show which provider we're switching to
		if s.selectedIndex < len(s.providers) {
			providerName := getProviderDisplayName(s.providers[s.selectedIndex].ID)
			builder.WriteString(textStyle.Render(fmt.Sprintf("Switching to %s...", providerName)))
		} else {
			builder.WriteString(textStyle.Render("Switching provider..."))
		}

		builder.WriteString("\n")
	}

	return builder.String()
}

func (s *SimpleProviderToggle) renderEmpty(t theme.Theme) string {
	emptyStyle := lipgloss.NewStyle().
		Foreground(t.TextMuted()).
		Padding(2)

	msg := "No authenticated CLI providers found.\n\n"
	msg += "Run /auth to authenticate with providers."

	return emptyStyle.Render(msg)
}

func (s *SimpleProviderToggle) renderError(t theme.Theme) string {
	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF4444")).
		Bold(true).
		Padding(1, 2)

	helpStyle := lipgloss.NewStyle().
		Foreground(t.TextMuted()).
		Padding(0, 2)

	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(errorStyle.Render("⚠ Failed to load providers"))
	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render(fmt.Sprintf("Error: %s", s.loadError.Error())))
	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("Press R to retry | Esc to cancel"))
	b.WriteString("\n")

	return b.String()
}

func (s *SimpleProviderToggle) renderProviderChip(provider opencode.Provider, isSelected bool, number int, t theme.Theme) string {
	// Get provider color
	color := s.getProviderColor(provider.ID)

	chipStyle := lipgloss.NewStyle().
		Padding(0, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(color)

	if isSelected {
		chipStyle = chipStyle.
			Background(color).
			Foreground(t.Background()).
			Bold(true)
	} else {
		chipStyle = chipStyle.
			Foreground(color)
	}

	// Show: [number] Name (modelCount)
	displayName := fmt.Sprintf("[%d] %s (%d)", number, getProviderDisplayName(provider.ID), len(provider.Models))
	return chipStyle.Render(displayName)
}

func (s *SimpleProviderToggle) getProviderColor(providerID string) compat.AdaptiveColor {
	// Use the same brand colors as RGB definitions for consistency
	switch providerID {
	case "anthropic", "claude":
		// Claude brand: warm orange/peach #E07856
		return compat.AdaptiveColor{Light: lipgloss.Color("#E07856"), Dark: lipgloss.Color("#E07856")}
	case "google", "gemini":
		// Gemini brand: blue-to-purple gradient (mid-purple) #8B7FD8
		return compat.AdaptiveColor{Light: lipgloss.Color("#8B7FD8"), Dark: lipgloss.Color("#8B7FD8")}
	case "openai", "codex":
		// OpenAI/Codex brand: teal/cyan #10A37F
		return compat.AdaptiveColor{Light: lipgloss.Color("#10A37F"), Dark: lipgloss.Color("#10A37F")}
	case "xai", "grok":
		// Grok/xAI brand: red #FF4444
		return compat.AdaptiveColor{Light: lipgloss.Color("#FF4444"), Dark: lipgloss.Color("#FF4444")}
	case "qwen":
		// Qwen brand: golden orange (from badge) #FFA726
		return compat.AdaptiveColor{Light: lipgloss.Color("#FFA726"), Dark: lipgloss.Color("#FFA726")}
	default:
		return theme.CurrentTheme().Primary()
	}
}

// getDefaultModelForProvider returns the best/default model for a provider (deterministic)
func (s *SimpleProviderToggle) getDefaultModelForProvider(provider opencode.Provider) opencode.Model {
	// Priority order for each provider's models (latest SOTA models as of 2025)
	priorities := map[string][]string{
		"claude": {
			"claude-sonnet-4-5",              // Latest Sonnet 4.5 (main coding model)
			"claude-opus-4-1",                // Opus 4.1 (complex reasoning)
			"claude-sonnet-4",                // Sonnet 4
			"claude-3-7-sonnet",              // Sonnet 3.7
			"claude-3-5-sonnet-20241022",     // Sonnet 3.5 (dated version)
			"claude-3-5-haiku-20241022",      // Haiku 3.5 (fast/lightweight)
		},
		"codex": {
			"gpt-5",                          // GPT-5 base model
			"o3",                             // O3 reasoning model
			"gpt-5-mini",                     // GPT-5 mini
			"o3-mini",                        // O3 mini
			"gpt-5-nano",                     // GPT-5 nano
			"gpt-4-5",                        // GPT-4.5
			"gpt-4o",                         // GPT-4o
			"gpt-4o-mini",                    // GPT-4o mini
		},
		"gemini": {
			"gemini-2.5-pro",                 // Gemini Pro 2.5 (latest)
			"gemini-2.5-flash",               // Gemini 2.5 Flash
			"gemini-2.5-flash-lite",          // Gemini 2.5 Flash Lite
			"gemini-2.5-flash-image",         // Gemini 2.5 Flash Image
			"gemini-2.5-computer-use",        // Gemini 2.5 Computer Use
			"gemini-2.5-deep-think",          // Gemini 2.5 Deep Think
			"gemini-exp-1206",                // Experimental release
			"gemini-2.0-flash-exp",           // Gemini 2.0 Flash (legacy)
		},
		"grok": {
			"grok-beta",                      // Grok beta (latest)
			"grok-2-1212",                    // Grok 2 dated release
		},
		"qwen": {
			"qwen3-max",                      // Qwen 3 Max (best general purpose)
			"qwen3-thinking-2507",            // Qwen 3 Thinking (reasoning model)
			"qwen3-next",                     // Qwen 3 Next
			"qwen3-omni",                     // Qwen 3 Omni (multimodal)
			"qwen3-instruct-2507",            // Qwen 3 Instruct
			"qwen3-235b",                     // Qwen 3 235B
			"qwen3-32b",                      // Qwen 3 32B
		},
	}

	// Check if we have priorities for this provider
	if prefs, ok := priorities[provider.ID]; ok {
		for _, modelID := range prefs {
			if model, exists := provider.Models[modelID]; exists {
				return model
			}
		}
	}

	// Fallback: return first model (with sorted keys for determinism)
	if len(provider.Models) > 0 {
		keys := make([]string, 0, len(provider.Models))
		for k := range provider.Models {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		return provider.Models[keys[0]]
	}

	// Should never happen, but return empty model as last resort
	return opencode.Model{}
}

func (s *SimpleProviderToggle) renderProviderModels(provider opencode.Provider, t theme.Theme) string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(t.TextMuted()).
		Padding(0, 2)
	b.WriteString(headerStyle.Render(fmt.Sprintf("%d models available:", len(provider.Models))))
	b.WriteString("\n\n")

	modelStyle := lipgloss.NewStyle().
		Foreground(t.Text()).
		Padding(0, 3)

	// Get models sorted by priority (best models first)
	modelIDs := make([]string, 0, len(provider.Models))
	for id := range provider.Models {
		modelIDs = append(modelIDs, id)
	}

	// Sort models using the same priority system as getDefaultModelForProvider (latest SOTA models as of 2025)
	priorities := map[string][]string{
		"claude":  {"claude-sonnet-4-5", "claude-opus-4-1", "claude-sonnet-4", "claude-3-7-sonnet", "claude-3-5-sonnet-20241022", "claude-3-5-haiku-20241022"},
		"codex":   {"gpt-5", "o3", "gpt-5-mini", "o3-mini", "gpt-5-nano", "gpt-4-5", "gpt-4o", "gpt-4o-mini"},
		"gemini":  {"gemini-2.5-pro", "gemini-2.5-flash", "gemini-2.5-flash-lite", "gemini-2.5-flash-image", "gemini-2.5-computer-use", "gemini-2.5-deep-think", "gemini-exp-1206", "gemini-2.0-flash-exp"},
		"grok":    {"grok-beta", "grok-2-1212"},
		"qwen":    {"qwen3-max", "qwen3-thinking-2507", "qwen3-next", "qwen3-omni", "qwen3-instruct-2507", "qwen3-235b", "qwen3-32b"},
	}

	priorityList := priorities[provider.ID]
	sort.Slice(modelIDs, func(i, j int) bool {
		// Find positions in priority list
		posI, posJ := -1, -1
		for idx, pID := range priorityList {
			if pID == modelIDs[i] {
				posI = idx
			}
			if pID == modelIDs[j] {
				posJ = idx
			}
		}

		// If both in priority list, compare positions
		if posI >= 0 && posJ >= 0 {
			return posI < posJ
		}
		// Priority models come first
		if posI >= 0 {
			return true
		}
		if posJ >= 0 {
			return false
		}
		// Otherwise alphabetical
		return modelIDs[i] < modelIDs[j]
	})

	count := 0
	for _, modelID := range modelIDs {
		if count >= 5 { // Show max 5 models
			remaining := len(provider.Models) - count
			moreStyle := lipgloss.NewStyle().
				Foreground(t.TextMuted()).
				Padding(0, 3)
			b.WriteString(moreStyle.Render(fmt.Sprintf("... and %d more", remaining)))
			break
		}
		model := provider.Models[modelID]
		b.WriteString(modelStyle.Render("• " + model.Name))
		b.WriteString("\n")
		count++
	}

	return b.String()
}

func (s *SimpleProviderToggle) SetSize(width, height int) {
	s.width = width
	s.height = height
}

// IsLoading returns whether the toggle is currently loading providers
func (s *SimpleProviderToggle) IsLoading() bool {
	return s.isLoading
}

// IsSwitching returns whether the toggle is currently showing the switching animation
func (s *SimpleProviderToggle) IsSwitching() bool {
	return s.isSwitching
}

// CalculateRequiredWidth calculates the minimum modal width needed to fit all provider chips
func (s *SimpleProviderToggle) CalculateRequiredWidth() int {
	if len(s.providers) == 0 {
		return 60 // Default minimum width
	}

	totalWidth := 0
	for i, provider := range s.providers {
		// Calculate chip content: [number] DisplayName (modelCount)
		// Example: "[1] Claude (6)" = 14 chars
		displayName := getProviderDisplayName(provider.ID)
		chipText := fmt.Sprintf("[%d] %s (%d)", i+1, displayName, len(provider.Models))

		// Add border (2 chars) + padding (4 chars) = 6 extra chars
		chipWidth := len(chipText) + 6
		totalWidth += chipWidth

		// Add spacing between chips (2 spaces)
		if i < len(s.providers)-1 {
			totalWidth += 2
		}
	}

	// Add container padding (4 chars on sides) + title padding (4 chars)
	totalWidth += 8

	// Ensure minimum width of 60 and maximum of 120
	if totalWidth < 60 {
		totalWidth = 60
	}
	if totalWidth > 120 {
		totalWidth = 120
	}

	slog.Debug("calculated modal width",
		"providers", len(s.providers),
		"required_width", totalWidth)

	return totalWidth
}

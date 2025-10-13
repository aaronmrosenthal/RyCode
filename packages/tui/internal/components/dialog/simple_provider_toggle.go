package dialog

import (
	"context"
	"fmt"
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
	"github.com/aaronmrosenthal/rycode/internal/splash"
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
	return tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
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

	logModelsDebug("=== SimpleProviderToggle.loadAuthenticatedProviders() START ===")

	// Get CLI providers directly from AuthBridge (doesn't require full client)
	cliProviders, err := s.app.AuthBridge.GetCLIProviders(ctx)
	if err != nil {
		logModelsDebug("ERROR: Failed to load CLI providers: %v", err)
		return nil, fmt.Errorf("failed to load CLI providers: %w", err)
	}

	logModelsDebug("Got %d CLI providers from GetCLIProviders()", len(cliProviders))
	for i, p := range cliProviders {
		logModelsDebug("  CLI Provider %d: ID=%s, Models=%d", i, p.Provider, len(p.Models))
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
			logModelsDebug("Checking provider: %s", providerID)
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
			logModelsDebug("  Auth check ERROR for %s: %v", result.providerID, result.err)
			continue
		}

		logModelsDebug("  Auth status for %s: IsAuthenticated=%v, ModelsCount=%d",
			result.providerID, result.authStatus.IsAuthenticated, result.authStatus.ModelsCount)

		if !result.authStatus.IsAuthenticated {
			logModelsDebug("  SKIPPED %s: not authenticated", result.providerID)
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

		logModelsDebug("  ADDED %s to providers list with %d models", result.providerID, len(models))
		providers = append(providers, provider)
	}

	logModelsDebug("=== SimpleProviderToggle.loadAuthenticatedProviders() END: %d providers loaded ===", len(providers))
	for i, p := range providers {
		logModelsDebug("  Final Provider %d: ID=%s, Name=%s, ModelsCount=%d", i, p.ID, p.Name, len(p.Models))
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
		return s, nil

	case loadingTickMsg:
		if s.isLoading {
			// Torus animation advances automatically in Render()
			return s, s.tickLoadingAnimation()
		}

		// Fade animation when switching providers
		if s.isSwitching {
			elapsed := time.Since(s.switchStartTime)
			totalDuration := 600 * time.Millisecond // Total animation: 200ms fade-in + 200ms hold + 200ms fade-out
			fadeInDuration := 200 * time.Millisecond
			fadeOutStart := 400 * time.Millisecond

			if elapsed < fadeInDuration {
				// Fade in (0.0 -> 1.0 over 200ms)
				s.fadeOpacity = float64(elapsed) / float64(fadeInDuration)
				logModelsDebug("Fade IN: opacity=%.2f", s.fadeOpacity)
			} else if elapsed < fadeOutStart {
				// Hold at full opacity
				s.fadeOpacity = 1.0
				logModelsDebug("Hold: opacity=%.2f", s.fadeOpacity)
			} else if elapsed < totalDuration {
				// Fade out (1.0 -> 0.0 over 200ms)
				fadeOutElapsed := elapsed - fadeOutStart
				s.fadeOpacity = 1.0 - (float64(fadeOutElapsed) / float64(totalDuration-fadeOutStart))
				logModelsDebug("Fade OUT: opacity=%.2f", s.fadeOpacity)
			} else {
				// Animation complete
				s.isSwitching = false
				s.fadeOpacity = 0.0
				logModelsDebug("Animation COMPLETE")
				return s, nil
			}

			// Continue animation
			return s, s.tickLoadingAnimation()
		}

		return s, nil

	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		return s, nil

	case tea.KeyPressMsg:
		logModelsDebug("=== SimpleProviderToggle KeyPress: %s ===", msg.String())
		logModelsDebug("Current selectedIndex: %d, Total providers: %d", s.selectedIndex, len(s.providers))
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("tab"))):
			logModelsDebug("✓ Tab key MATCHED!")
			// Tab: cycle to next provider with cortex fade animation
			if len(s.providers) > 0 {
				oldIndex := s.selectedIndex
				s.selectedIndex = (s.selectedIndex + 1) % len(s.providers)
				logModelsDebug("Tab cycling: %d -> %d (provider: %s)", oldIndex, s.selectedIndex, s.providers[s.selectedIndex].ID)

				// Start cortex fade animation (600ms total: fade-in, hold, fade-out)
				s.isSwitching = true
				s.switchStartTime = time.Now()
				s.fadeOpacity = 0.0 // Start from transparent
				logModelsDebug("✓ SWITCHING ANIMATION ACTIVATED - starting fade-in")
				return s, s.tickLoadingAnimation()
			} else {
				logModelsDebug("Cannot cycle: no providers loaded")
			}
			return s, nil

		case key.Matches(msg, key.NewBinding(key.WithKeys("shift+tab"))):
			logModelsDebug("✓ Shift+Tab key MATCHED!")
			// Shift+Tab: cycle to previous provider with cortex fade animation
			if len(s.providers) > 0 {
				oldIndex := s.selectedIndex
				s.selectedIndex--
				if s.selectedIndex < 0 {
					s.selectedIndex = len(s.providers) - 1
				}
				logModelsDebug("Shift+Tab cycling: %d -> %d (provider: %s)", oldIndex, s.selectedIndex, s.providers[s.selectedIndex].ID)

				// Start cortex fade animation (600ms total: fade-in, hold, fade-out)
				s.isSwitching = true
				s.switchStartTime = time.Now()
				s.fadeOpacity = 0.0 // Start from transparent
				logModelsDebug("✓ SWITCHING ANIMATION ACTIVATED - starting fade-in")
				return s, s.tickLoadingAnimation()
			} else {
				logModelsDebug("Cannot cycle: no providers loaded")
			}
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
			logModelsDebug("✓ Esc key MATCHED - closing modal")
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

				logModelsDebug("Direct selection: index=%d, provider=%s", index, provider.ID)
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
			logModelsDebug("Key not matched by any case: %s", msg.String())
		}
	}

	return s, nil
}

func (s *SimpleProviderToggle) View() string {
	t := theme.CurrentTheme()

	logModelsDebug("View() called - isLoading=%v, isSwitching=%v", s.isLoading, s.isSwitching)

	if s.isLoading {
		return s.renderLoading(t)
	}

	// Show cortex when switching providers
	if s.isSwitching {
		logModelsDebug("✓ Rendering SWITCHING view with cortex!")
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

	// Provider chips (horizontal layout with numbers)
	var chipsWithSpacing []string
	for i, provider := range s.providers {
		isSelected := i == s.selectedIndex
		chip := s.renderProviderChip(provider, isSelected, i+1, t)
		chipsWithSpacing = append(chipsWithSpacing, chip)
		// Add spacing between chips (but not after the last one)
		if i < len(s.providers)-1 {
			chipsWithSpacing = append(chipsWithSpacing, "  ")
		}
	}

	// Join chips horizontally
	chipsRow := lipgloss.JoinHorizontal(lipgloss.Left, chipsWithSpacing...)
	chipContainerStyle := lipgloss.NewStyle().Padding(0, 2)
	b.WriteString(chipContainerStyle.Render(chipsRow))
	b.WriteString("\n\n")

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
	switch providerID {
	case "anthropic", "claude":
		// Claude brand: warm orange/peach (from screenshot)
		return compat.AdaptiveColor{Light: lipgloss.Color("#E07856"), Dark: lipgloss.Color("#E07856")}
	case "google", "gemini":
		// Gemini brand: light blue (from screenshot - matches the blue gradient)
		return compat.AdaptiveColor{Light: lipgloss.Color("#4A90E2"), Dark: lipgloss.Color("#4A90E2")}
	case "openai", "codex":
		// OpenAI/Codex brand: teal/cyan (keeping OpenAI's signature color)
		return compat.AdaptiveColor{Light: lipgloss.Color("#10A37F"), Dark: lipgloss.Color("#10A37F")}
	case "xai", "grok":
		// Grok/xAI brand: red
		return compat.AdaptiveColor{Light: lipgloss.Color("#FF4444"), Dark: lipgloss.Color("#FF4444")}
	case "qwen":
		// Qwen brand: orange/amber (from screenshot - matches the orange gradient)
		return compat.AdaptiveColor{Light: lipgloss.Color("#FFA500"), Dark: lipgloss.Color("#FFA500")}
	default:
		return theme.CurrentTheme().Primary()
	}
}

// getDefaultModelForProvider returns the best/default model for a provider (deterministic)
func (s *SimpleProviderToggle) getDefaultModelForProvider(provider opencode.Provider) opencode.Model {
	// Priority order for each provider's models (latest SOTA models as of 2025)
	priorities := map[string][]string{
		"claude": {
			"claude-opus-4-5",                // Latest Opus 4.5 (plan/complex reasoning)
			"claude-sonnet-4-5-20250929",     // Latest Sonnet 4.5 dated version
			"claude-sonnet-4-5",              // Latest Sonnet 4.5 (main coding model)
			"claude-sonnet-4-5-1m",           // Sonnet 4.5 with 1M context (failover)
			"claude-sonnet-3-5-20241022",     // Previous Sonnet 3.5
			"claude-sonnet-3-5",              // Legacy Sonnet 3.5
		},
		"codex": {
			"gpt-5-codex",                    // GPT-5 specialized codex model
			"codex-gpt-5",                    // Alternative GPT-5 codex naming
			"gpt-5",                          // GPT-5 base model
			"gpt-4o",                         // GPT-4o (current best non-5)
			"gpt-4o-mini",                    // GPT-4o mini
			"gpt-4-turbo",                    // GPT-4 turbo
			"gpt-4",                          // GPT-4 base
		},
		"gemini": {
			"gemini-pro-2.5",                 // Gemini Pro 2.5 (latest)
			"gemini-2.5-flash",               // Gemini 2.5 Flash (failover)
			"gemini-flash-2.5",               // Alternative Flash 2.5 naming
			"gemini-2.0-flash-exp",           // Gemini 2.0 Flash (current)
			"gemini-exp-1206",                // Experimental release
			"gemini-pro-1.5",                 // Pro 1.5
			"gemini-pro",                     // Legacy Pro
		},
		"grok": {
			"grok-beta",                      // Grok beta (latest)
			"grok-2-1212",                    // Grok 2 dated release
		},
		"qwen": {
			"qwen-coder-3",                   // Qwen Coder 3 (specialized coding model)
			"qwen-coder-3.0",                 // Alternative Qwen Coder 3 naming
			"qwen3-coder",                    // Alternative naming format
			"qwen-max",                       // Qwen Max (general purpose)
			"qwen-plus",                      // Qwen Plus
			"qwen-turbo",                     // Qwen Turbo
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
		"claude":  {"claude-opus-4-5", "claude-sonnet-4-5-20250929", "claude-sonnet-4-5", "claude-sonnet-4-5-1m", "claude-sonnet-3-5-20241022", "claude-sonnet-3-5"},
		"codex":   {"gpt-5-codex", "codex-gpt-5", "gpt-5", "gpt-4o", "gpt-4o-mini", "gpt-4-turbo", "gpt-4"},
		"gemini":  {"gemini-pro-2.5", "gemini-2.5-flash", "gemini-flash-2.5", "gemini-2.0-flash-exp", "gemini-exp-1206", "gemini-pro-1.5", "gemini-pro"},
		"grok":    {"grok-beta", "grok-2-1212"},
		"qwen":    {"qwen-coder-3", "qwen-coder-3.0", "qwen3-coder", "qwen-max", "qwen-plus", "qwen-turbo"},
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

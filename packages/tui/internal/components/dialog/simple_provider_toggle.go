package dialog

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/lipgloss/v2/compat"

	"github.com/aaronmrosenthal/rycode-sdk-go"
	"github.com/aaronmrosenthal/rycode/internal/app"
	"github.com/aaronmrosenthal/rycode/internal/components/modal"
	"github.com/aaronmrosenthal/rycode/internal/theme"
	"github.com/aaronmrosenthal/rycode/internal/util"
)

// SimpleProviderToggle is a minimal provider selector that just cycles between CLI providers
type SimpleProviderToggle struct {
	app             *app.App
	providers       []opencode.Provider
	selectedIndex   int
	width           int
	height          int
}

// NewSimpleProviderToggle creates a new simple provider toggle
func NewSimpleProviderToggle(app *app.App) *SimpleProviderToggle {
	return &SimpleProviderToggle{
		app:           app,
		providers:     []opencode.Provider{},
		selectedIndex: 0,
	}
}

func (s *SimpleProviderToggle) Init() tea.Cmd {
	// Load only authenticated CLI providers
	s.loadAuthenticatedProviders()
	return nil
}

func (s *SimpleProviderToggle) loadAuthenticatedProviders() {
	ctx := context.Background()

	logModelsDebug("=== SimpleProviderToggle.loadAuthenticatedProviders() START ===")

	// Get CLI providers directly from AuthBridge (doesn't require full client)
	cliProviders, err := s.app.AuthBridge.GetCLIProviders(ctx)
	if err != nil {
		logModelsDebug("ERROR: Failed to load CLI providers: %v", err)
		return
	}

	logModelsDebug("Got %d CLI providers from GetCLIProviders()", len(cliProviders))
	for i, p := range cliProviders {
		logModelsDebug("  CLI Provider %d: ID=%s, Models=%d", i, p.Provider, len(p.Models))
	}

	// Convert CLI providers to opencode.Provider format and check authentication
	s.providers = []opencode.Provider{}
	for _, cliProv := range cliProviders {
		providerID := cliProv.Provider
		logModelsDebug("Checking provider: %s", providerID)

		// Check if authenticated
		authStatus, err := s.app.AuthBridge.CheckAuthStatus(ctx, providerID)
		if err != nil {
			logModelsDebug("  Auth check ERROR for %s: %v", providerID, err)
			continue
		}
		logModelsDebug("  Auth status for %s: IsAuthenticated=%v, ModelsCount=%d",
			providerID, authStatus.IsAuthenticated, authStatus.ModelsCount)

		if !authStatus.IsAuthenticated {
			logModelsDebug("  SKIPPED %s: not authenticated", providerID)
			continue
		}

		// Convert to opencode.Provider format
		models := make(map[string]opencode.Model)
		for _, modelID := range cliProv.Models {
			models[modelID] = opencode.Model{
				ID:   modelID,
				Name: formatModelName(modelID),
			}
		}

		provider := opencode.Provider{
			ID:     providerID,
			Name:   formatProviderName(providerID),
			Models: models,
		}

		logModelsDebug("  ADDED %s to providers list with %d models", providerID, len(models))
		s.providers = append(s.providers, provider)
	}

	logModelsDebug("=== SimpleProviderToggle.loadAuthenticatedProviders() END: %d providers loaded ===", len(s.providers))
	for i, p := range s.providers {
		logModelsDebug("  Final Provider %d: ID=%s, Name=%s, ModelsCount=%d", i, p.ID, p.Name, len(p.Models))
	}
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
		"anthropic": "Anthropic",
		"openai":    "OpenAI",
		"google":    "Google",
		"xai":       "xAI",
		"qwen":      "Alibaba",
	}
	if name, ok := names[providerID]; ok {
		return name
	}
	return strings.ToUpper(providerID[:1]) + providerID[1:]
}


func (s *SimpleProviderToggle) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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
			// Tab: cycle to next provider
			if len(s.providers) > 0 {
				oldIndex := s.selectedIndex
				s.selectedIndex = (s.selectedIndex + 1) % len(s.providers)
				logModelsDebug("Tab cycling: %d -> %d (provider: %s)", oldIndex, s.selectedIndex, s.providers[s.selectedIndex].ID)
			} else {
				logModelsDebug("Cannot cycle: no providers loaded")
			}
			return s, nil

		case key.Matches(msg, key.NewBinding(key.WithKeys("shift+tab"))):
			logModelsDebug("✓ Shift+Tab key MATCHED!")
			// Shift+Tab: cycle to previous provider
			if len(s.providers) > 0 {
				oldIndex := s.selectedIndex
				s.selectedIndex--
				if s.selectedIndex < 0 {
					s.selectedIndex = len(s.providers) - 1
				}
				logModelsDebug("Shift+Tab cycling: %d -> %d (provider: %s)", oldIndex, s.selectedIndex, s.providers[s.selectedIndex].ID)
			} else {
				logModelsDebug("Cannot cycle: no providers loaded")
			}
			return s, nil

		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			// Enter: select current provider
			if len(s.providers) > 0 && s.selectedIndex < len(s.providers) {
				provider := s.providers[s.selectedIndex]

				// Get the first model from this provider
				var selectedModel opencode.Model
				for _, model := range provider.Models {
					selectedModel = model
					break
				}

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
		default:
			logModelsDebug("Key not matched by any case: %s", msg.String())
		}
	}

	return s, nil
}

func (s *SimpleProviderToggle) View() string {
	t := theme.CurrentTheme()

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

	// Provider chips (horizontal layout)
	var chipsWithSpacing []string
	for i, provider := range s.providers {
		isSelected := i == s.selectedIndex
		chip := s.renderProviderChip(provider, isSelected, t)
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
	footer := "Tab: Next | Shift+Tab: Previous | Enter: Select | Esc: Cancel"
	b.WriteString(footerStyle.Render(footer))

	return b.String()
}

func (s *SimpleProviderToggle) renderEmpty(t theme.Theme) string {
	emptyStyle := lipgloss.NewStyle().
		Foreground(t.TextMuted()).
		Padding(2)

	msg := "No authenticated CLI providers found.\n\n"
	msg += "Run /auth to authenticate with providers."

	return emptyStyle.Render(msg)
}

func (s *SimpleProviderToggle) renderProviderChip(provider opencode.Provider, isSelected bool, t theme.Theme) string {
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

	return chipStyle.Render(s.getProviderDisplayName(provider.ID))
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

func (s *SimpleProviderToggle) getProviderDisplayName(providerID string) string {
	switch providerID {
	case "anthropic", "claude":
		return "Claude"
	case "google", "gemini":
		return "Gemini"
	case "openai", "codex":
		return "Codex"
	case "xai", "grok":
		return "Grok"
	case "qwen":
		return "Qwen"
	default:
		return providerID
	}
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

	count := 0
	for _, model := range provider.Models {
		if count >= 5 { // Show max 5 models
			remaining := len(provider.Models) - count
			moreStyle := lipgloss.NewStyle().
				Foreground(t.TextMuted()).
				Padding(0, 3)
			b.WriteString(moreStyle.Render(fmt.Sprintf("... and %d more", remaining)))
			break
		}
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

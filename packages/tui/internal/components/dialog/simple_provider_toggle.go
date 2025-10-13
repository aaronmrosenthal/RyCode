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

	// Get all providers
	allProviders, err := s.app.ListProviders(ctx)
	if err != nil {
		logModelsDebug("Failed to load providers: %v", err)
		return
	}

	// Filter to only authenticated CLI providers
	s.providers = []opencode.Provider{}
	for _, provider := range allProviders {
		// Check if authenticated
		authStatus, err := s.app.AuthBridge.CheckAuthStatus(ctx, provider.ID)
		if err != nil || !authStatus.IsAuthenticated {
			continue
		}

		// Only include CLI providers (claude, qwen, codex, gemini, grok)
		if s.isCLIProvider(provider.ID) {
			s.providers = append(s.providers, provider)
		}
	}

	logModelsDebug("Loaded %d authenticated CLI providers", len(s.providers))
}

func (s *SimpleProviderToggle) isCLIProvider(providerID string) bool {
	cliProviders := []string{"anthropic", "qwen", "codex", "google", "xai"}
	for _, id := range cliProviders {
		if providerID == id {
			return true
		}
	}
	return false
}

func (s *SimpleProviderToggle) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		return s, nil

	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("tab"))):
			// Tab: cycle to next provider
			if len(s.providers) > 0 {
				s.selectedIndex = (s.selectedIndex + 1) % len(s.providers)
			}
			return s, nil

		case key.Matches(msg, key.NewBinding(key.WithKeys("shift+tab"))):
			// Shift+Tab: cycle to previous provider
			if len(s.providers) > 0 {
				s.selectedIndex--
				if s.selectedIndex < 0 {
					s.selectedIndex = len(s.providers) - 1
				}
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
			return s, util.CmdHandler(modal.CloseModalMsg{})
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

	// Provider chips (horizontal layout like landing page)
	chipContainerStyle := lipgloss.NewStyle().Padding(0, 2)
	var chips []string

	for i, provider := range s.providers {
		isSelected := i == s.selectedIndex
		chip := s.renderProviderChip(provider, isSelected, t)
		chips = append(chips, chip)
	}

	b.WriteString(chipContainerStyle.Render(strings.Join(chips, "  ")))
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
	case "anthropic":
		return compat.AdaptiveColor{Light: lipgloss.Color("#7aa2f7"), Dark: lipgloss.Color("#7aa2f7")} // Claude blue
	case "google":
		return compat.AdaptiveColor{Light: lipgloss.Color("#ea4aaa"), Dark: lipgloss.Color("#ea4aaa")} // Gemini pink
	case "openai":
		return compat.AdaptiveColor{Light: lipgloss.Color("#ff6b35"), Dark: lipgloss.Color("#ff6b35")} // OpenAI orange
	case "xai":
		return compat.AdaptiveColor{Light: lipgloss.Color("#00ffff"), Dark: lipgloss.Color("#00ffff")} // Grok cyan
	case "qwen":
		return compat.AdaptiveColor{Light: lipgloss.Color("#ff00ff"), Dark: lipgloss.Color("#ff00ff")} // Qwen magenta
	default:
		return theme.CurrentTheme().Primary()
	}
}

func (s *SimpleProviderToggle) getProviderDisplayName(providerID string) string {
	switch providerID {
	case "anthropic":
		return "Claude"
	case "google":
		return "Gemini"
	case "openai":
		return "GPT-5"
	case "xai":
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

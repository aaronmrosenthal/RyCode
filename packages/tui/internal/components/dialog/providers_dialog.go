package dialog

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aaronmrosenthal/rycode/internal/app"
	"github.com/aaronmrosenthal/rycode/internal/layout"
	"github.com/aaronmrosenthal/rycode/internal/styles"
	"github.com/aaronmrosenthal/rycode/internal/theme"
	"github.com/aaronmrosenthal/rycode/internal/typography"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

// ProviderStatus represents the status of a single provider
type ProviderStatus struct {
	ID              string
	Name            string
	IsAuthenticated bool
	Health          string // "healthy", "degraded", "down", "unknown"
	ModelsCount     int
	LastChecked     time.Time
	APIKeyMasked    string // Last 4 chars only
}

// ProvidersDialog displays provider management dashboard
type ProvidersDialog interface {
	layout.Modal
}

type providersDialog struct {
	app       *app.App
	providers []ProviderStatus
	width     int
	height    int
	focused   int // Currently focused provider
}

// NewProvidersDialog creates a new providers management dialog
func NewProvidersDialog(app *app.App) ProvidersDialog {
	dialog := &providersDialog{
		app:       app,
		providers: make([]ProviderStatus, 0),
		focused:   0,
	}

	// Load provider statuses
	dialog.loadProviders()

	return dialog
}

// loadProviders fetches current provider statuses
func (d *providersDialog) loadProviders() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Known providers
	providerIDs := []string{"anthropic", "openai", "google", "xai", "qwen"}
	providerNames := map[string]string{
		"anthropic": "Anthropic (Claude)",
		"openai":    "OpenAI (GPT)",
		"google":    "Google (Gemini)",
		"xai":       "X.AI (Grok)",
		"qwen":      "Alibaba (Qwen)",
	}

	for _, id := range providerIDs {
		status := ProviderStatus{
			ID:          id,
			Name:        providerNames[id],
			LastChecked: time.Now(),
		}

		// Check authentication status
		if d.app.AuthBridge != nil {
			authStatus, err := d.app.AuthBridge.CheckAuthStatus(ctx, id)
			if err == nil {
				status.IsAuthenticated = authStatus.IsAuthenticated
				status.ModelsCount = authStatus.ModelsCount
			}

			// Check health
			health, err := d.app.AuthBridge.GetProviderHealth(ctx, id)
			if err == nil {
				status.Health = health.Status
			} else {
				status.Health = "unknown"
			}
		}

		// Mask API key (would come from secure storage in production)
		if status.IsAuthenticated {
			status.APIKeyMasked = "sk-...xxxx" // Placeholder
		}

		d.providers = append(d.providers, status)
	}
}

func (d *providersDialog) Init() tea.Cmd {
	return nil
}

func (d *providersDialog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		d.width = msg.Width
		d.height = msg.Height

	case tea.KeyPressMsg:
		switch msg.String() {
		case "up", "k":
			if d.focused > 0 {
				d.focused--
			}
		case "down", "j":
			if d.focused < len(d.providers)-1 {
				d.focused++
			}
		case "r":
			// Refresh provider statuses
			d.loadProviders()
		case "a":
			// Authenticate focused provider
			if d.focused < len(d.providers) {
				provider := d.providers[d.focused]
				if !provider.IsAuthenticated {
					// Would open auth prompt here
					// For now, just refresh
					d.loadProviders()
				}
			}
		}
	}

	return d, nil
}

func (d *providersDialog) View() string {
	t := theme.CurrentTheme()
	typo := typography.New()

	var sections []string

	// Header
	header := styles.NewStyle().
		Foreground(t.Primary()).
		Bold(true).
		Render("🔐 Provider Management")

	sections = append(sections, header)
	sections = append(sections, "")

	// Summary stats
	authenticated := 0
	healthy := 0
	totalModels := 0

	for _, p := range d.providers {
		if p.IsAuthenticated {
			authenticated++
			totalModels += p.ModelsCount
		}
		if p.Health == "healthy" {
			healthy++
		}
	}

	summary := fmt.Sprintf("%d/%d authenticated • %d models available • %d healthy",
		authenticated, len(d.providers), totalModels, healthy)

	summaryStyle := styles.NewStyle().
		Foreground(t.TextMuted()).
		Faint(true)

	sections = append(sections, summaryStyle.Render(summary))
	sections = append(sections, "")
	sections = append(sections, typo.Subheading.Render("Providers"))
	sections = append(sections, "")

	// Provider cards
	for i, provider := range d.providers {
		card := d.renderProviderCard(provider, i == d.focused)
		sections = append(sections, card)
		if i < len(d.providers)-1 {
			sections = append(sections, "")
		}
	}

	// Footer with help
	sections = append(sections, "")
	sections = append(sections, "")

	helpStyle := styles.NewStyle().
		Foreground(t.TextMuted()).
		Faint(true)

	help := helpStyle.Render("↑/↓: Navigate  [a] Authenticate  [r] Refresh  [ESC] Close")
	sections = append(sections, help)

	content := strings.Join(sections, "\n")

	// Wrap in bordered box
	boxStyle := styles.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Border()).
		Padding(2, 3).
		Width(d.width - 4)

	return boxStyle.Render(content)
}

// renderProviderCard creates a card for a single provider
func (d *providersDialog) renderProviderCard(provider ProviderStatus, focused bool) string {
	t := theme.CurrentTheme()

	// Border color based on focus
	borderColor := t.Border()
	if focused {
		borderColor = t.Primary()
	}

	// Provider name and status
	nameStyle := styles.NewStyle().
		Foreground(t.Text()).
		Bold(true)

	name := nameStyle.Render(provider.Name)

	// Status badge
	var statusBadge string
	if provider.IsAuthenticated {
		statusBadge = styles.NewStyle().
			Foreground(t.Background()).
			Background(t.Success()).
			Bold(true).
			Padding(0, 1).
			Render("✓ Authenticated")
	} else {
		statusBadge = styles.NewStyle().
			Foreground(t.Background()).
			Background(t.TextMuted()).
			Bold(true).
			Padding(0, 1).
			Render("🔒 Locked")
	}

	// Health indicator
	healthIcon := ""
	healthColor := t.TextMuted()

	switch provider.Health {
	case "healthy":
		healthIcon = "●"
		healthColor = t.Success()
	case "degraded":
		healthIcon = "●"
		healthColor = t.Warning()
	case "down":
		healthIcon = "●"
		healthColor = t.Error()
	default:
		healthIcon = "○"
	}

	healthStyle := styles.NewStyle().
		Foreground(healthColor)

	health := healthStyle.Render(healthIcon + " " + provider.Health)

	// Models count
	modelsStyle := styles.NewStyle().
		Foreground(t.TextMuted()).
		Faint(true)

	models := ""
	if provider.IsAuthenticated {
		models = modelsStyle.Render(fmt.Sprintf("%d models", provider.ModelsCount))
	} else {
		models = modelsStyle.Render("No access")
	}

	// API key (masked)
	keyStyle := styles.NewStyle().
		Foreground(t.TextMuted()).
		Faint(true)

	key := ""
	if provider.IsAuthenticated {
		key = keyStyle.Render("Key: " + provider.APIKeyMasked)
	}

	// Layout
	line1 := name + "  " + statusBadge
	line2 := health + "  " + models
	line3 := key

	var lines []string
	lines = append(lines, line1)
	lines = append(lines, line2)
	if key != "" {
		lines = append(lines, line3)
	}

	content := strings.Join(lines, "\n")

	// Card style
	cardStyle := styles.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1, 2).
		Width(60)

	return cardStyle.Render(content)
}

func (d *providersDialog) Render(background string) string {
	return d.View()
}

func (d *providersDialog) Close() tea.Cmd {
	return nil
}

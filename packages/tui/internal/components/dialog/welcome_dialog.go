package dialog

import (
	"fmt"
	"strings"

	"github.com/aaronmrosenthal/rycode/internal/app"
	"github.com/aaronmrosenthal/rycode/internal/layout"
	"github.com/aaronmrosenthal/rycode/internal/styles"
	"github.com/aaronmrosenthal/rycode/internal/theme"
	"github.com/aaronmrosenthal/rycode/internal/typography"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

// WelcomeStep represents a step in the onboarding flow
type WelcomeStep struct {
	Title       string
	Content     string
	Bullets     []string
	Action      string // What user should do
	KeyHint     string // Keyboard hint
	Celebration string // Emoji/icon for this step
}

// WelcomeDialog displays the first-time user experience
type WelcomeDialog interface {
	layout.Modal
}

type welcomeDialog struct {
	app         *app.App
	steps       []WelcomeStep
	currentStep int
	width       int
	height      int
}

// NewWelcomeDialog creates a new welcome/onboarding dialog
func NewWelcomeDialog(app *app.App) WelcomeDialog {
	dialog := &welcomeDialog{
		app:         app,
		currentStep: 0,
	}

	dialog.setupSteps()

	return dialog
}

// setupSteps defines the onboarding flow
func (w *welcomeDialog) setupSteps() {
	w.steps = []WelcomeStep{
		{
			Title:       "Welcome to RyCode! 🚀",
			Celebration: "🎉",
			Content:     "RyCode is your AI-powered development assistant, built by Claude to demonstrate what's possible when AI designs tools.",
			Bullets: []string{
				"Intelligent model recommendations",
				"Real-time cost tracking & budgeting",
				"Multi-provider support (Claude, GPT, Gemini, Grok, Qwen)",
				"Beautiful TUI with keyboard shortcuts",
				"Usage analytics & optimization insights",
			},
			Action:  "Let's get you set up!",
			KeyHint: "Press ENTER to continue",
		},
		{
			Title:       "Choose Your AI Provider",
			Celebration: "🔐",
			Content:     "RyCode supports multiple AI providers. You'll need at least one API key to get started.",
			Bullets: []string{
				"Anthropic (Claude): Best for coding & reasoning",
				"OpenAI (GPT): Wide range of models",
				"Google (Gemini): Large context windows",
				"X.AI (Grok): Fast responses",
				"Alibaba (Qwen): Multilingual support",
			},
			Action:  "We'll help you authenticate in the next step",
			KeyHint: "Press ENTER to continue",
		},
		{
			Title:       "Quick Setup",
			Celebration: "⚡",
			Content:     "You have two options to get started:",
			Bullets: []string{
				"[A] Auto-detect credentials from environment variables",
				"[M] Manually enter API key for a specific provider",
				"",
				"Auto-detect looks for:",
				"  • ANTHROPIC_API_KEY",
				"  • OPENAI_API_KEY",
				"  • GOOGLE_API_KEY / GEMINI_API_KEY",
				"  • XAI_API_KEY / GROK_API_KEY",
				"  • QWEN_API_KEY",
			},
			Action:  "Choose your setup method",
			KeyHint: "Press [A] for auto-detect or [M] for manual",
		},
		{
			Title:       "Keyboard Shortcuts",
			Celebration: "⌨️",
			Content:     "RyCode is designed for keyboard warriors. Here are the essential shortcuts:",
			Bullets: []string{
				"Tab: Cycle through available models",
				"Ctrl+M: Open model selector",
				"Ctrl+P: Provider management dashboard",
				"Ctrl+I: Usage insights & analytics",
				"Ctrl+B: Budget forecast",
				"Ctrl+?: Show all keyboard shortcuts",
				"Ctrl+C: Exit",
			},
			Action:  "You'll learn more as you go!",
			KeyHint: "Press ENTER to continue",
		},
		{
			Title:       "Smart Features",
			Celebration: "🧠",
			Content:     "RyCode includes AI-powered features to help you optimize usage:",
			Bullets: []string{
				"💰 Cost alerts: Get notified before exceeding budget",
				"🔮 Predictive budgeting: Forecast month-end spending",
				"💡 Model recommendations: AI suggests the best model for each task",
				"📊 Usage insights: Beautiful charts showing your patterns",
				"⚡ Auto-optimization: Learn from your usage over time",
			},
			Action:  "All of this runs locally - no data sent to third parties",
			KeyHint: "Press ENTER to continue",
		},
		{
			Title:       "You're All Set! 🎊",
			Celebration: "✨",
			Content:     "RyCode is ready to use. Here's what to do next:",
			Bullets: []string{
				"1. Select a model (Tab or Ctrl+M)",
				"2. Start a conversation",
				"3. Check your usage insights (Ctrl+I)",
				"4. Explore keyboard shortcuts (Ctrl+?)",
				"5. Have fun building! 🚀",
			},
			Action:  "Remember: RyCode learns from your usage to provide better recommendations over time.",
			KeyHint: "Press ENTER to start using RyCode!",
		},
	}
}

func (w *welcomeDialog) Init() tea.Cmd {
	return nil
}

func (w *welcomeDialog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		w.width = msg.Width
		w.height = msg.Height

	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			// Next step
			if w.currentStep < len(w.steps)-1 {
				w.currentStep++
			} else {
				// Completed onboarding - would close dialog here
				// Mark as completed in user settings
			}

		case "left", "h":
			// Previous step
			if w.currentStep > 0 {
				w.currentStep--
			}

		case "right", "l":
			// Next step (alternative to enter)
			if w.currentStep < len(w.steps)-1 {
				w.currentStep++
			}

		case "a":
			// Auto-detect (on step 2)
			if w.currentStep == 2 {
				// Trigger auto-detect
				// Would call d.app.AuthBridge.AutoDetect() here
			}

		case "m":
			// Manual auth (on step 2)
			if w.currentStep == 2 {
				// Open manual auth dialog
				// Would open auth prompt here
			}
		}
	}

	return w, nil
}

func (w *welcomeDialog) View() string {
	if w.currentStep >= len(w.steps) {
		return ""
	}

	step := w.steps[w.currentStep]
	t := theme.CurrentTheme()
	typo := typography.New()

	var sections []string

	// Progress indicator
	progress := fmt.Sprintf("Step %d of %d", w.currentStep+1, len(w.steps))
	progressStyle := styles.NewStyle().
		Foreground(t.TextMuted()).
		Faint(true)

	sections = append(sections, progressStyle.Render(progress))
	sections = append(sections, "")

	// Celebration icon (big and centered)
	celebrationStyle := styles.NewStyle().
		Foreground(t.Primary()).
		Bold(true)

	celebration := celebrationStyle.Render(step.Celebration + "  " + step.Title)
	sections = append(sections, celebration)
	sections = append(sections, "")

	// Content
	contentStyle := typo.Body.
		Foreground(t.Text()).
		Width(60)

	content := contentStyle.Render(step.Content)
	sections = append(sections, content)
	sections = append(sections, "")

	// Bullets
	if len(step.Bullets) > 0 {
		for _, bullet := range step.Bullets {
			if bullet == "" {
				sections = append(sections, "")
				continue
			}

			bulletStyle := typo.Body.
				Foreground(t.Text()).
				Width(58)

			// Check for special formatting
			if strings.HasPrefix(bullet, "[") {
				// Action item - highlight the key
				bulletStyle = bulletStyle.Foreground(t.Primary())
			} else if strings.HasPrefix(bullet, "  •") || strings.HasPrefix(bullet, "  ") {
				// Indented item
				bulletStyle = bulletStyle.Foreground(t.TextMuted())
			}

			formatted := bulletStyle.Render("  • " + bullet)
			sections = append(sections, formatted)
		}
		sections = append(sections, "")
	}

	// Action
	if step.Action != "" {
		actionStyle := styles.NewStyle().
			Foreground(t.Info()).
			Italic(true)

		action := actionStyle.Render(step.Action)
		sections = append(sections, action)
		sections = append(sections, "")
	}

	// Navigation hint
	navStyle := styles.NewStyle().
		Foreground(t.TextMuted()).
		Faint(true)

	// Build nav based on current step
	var navParts []string

	if w.currentStep > 0 {
		navParts = append(navParts, "← Previous")
	}

	navParts = append(navParts, step.KeyHint)

	if w.currentStep < len(w.steps)-1 {
		navParts = append(navParts, "→ Next")
	}

	nav := navStyle.Render(strings.Join(navParts, "  |  "))
	sections = append(sections, nav)

	content = strings.Join(sections, "\n")

	// Progress bar
	progressBar := w.renderProgressBar()
	content = content + "\n\n" + progressBar

	// Wrap in bordered box
	boxStyle := styles.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Primary()).
		Padding(2, 3).
		Width(70)

	return boxStyle.Render(content)
}

// renderProgressBar creates a visual progress indicator
func (w *welcomeDialog) renderProgressBar() string {
	t := theme.CurrentTheme()

	bar := ""
	for i := 0; i < len(w.steps); i++ {
		if i == w.currentStep {
			// Current step - solid
			bar += styles.NewStyle().
				Foreground(t.Primary()).
				Render("●")
		} else if i < w.currentStep {
			// Completed step - solid but muted
			bar += styles.NewStyle().
				Foreground(t.Success()).
				Render("●")
		} else {
			// Future step - outlined
			bar += styles.NewStyle().
				Foreground(t.TextMuted()).
				Faint(true).
				Render("○")
		}

		// Add connector
		if i < len(w.steps)-1 {
			if i < w.currentStep {
				bar += styles.NewStyle().
					Foreground(t.Success()).
					Render("─")
			} else {
				bar += styles.NewStyle().
					Foreground(t.TextMuted()).
					Faint(true).
					Render("─")
			}
		}
	}

	// Center the bar
	return strings.Repeat(" ", 15) + bar
}

func (w *welcomeDialog) Render(background string) string {
	return w.View()
}

func (w *welcomeDialog) Close() tea.Cmd {
	return nil
}

// ShouldShowWelcome checks if welcome flow should be shown
func ShouldShowWelcome(app *app.App) bool {
	// Check if this is first run
	// In production, would check:
	// - User settings for completed_welcome flag
	// - Existence of any authenticated providers
	// - Previous usage history

	// For now, return false (would be true on actual first run)
	return false
}

package dialog

import (
	"strings"

	"github.com/aaronmrosenthal/rycode/internal/accessibility"
	"github.com/aaronmrosenthal/rycode/internal/app"
	"github.com/aaronmrosenthal/rycode/internal/layout"
	"github.com/aaronmrosenthal/rycode/internal/styles"
	"github.com/aaronmrosenthal/rycode/internal/theme"
	"github.com/aaronmrosenthal/rycode/internal/typography"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

// AccessibilityDialog displays accessibility settings
type AccessibilityDialog interface {
	layout.Modal
}

type accessibilityDialog struct {
	app      *app.App
	settings *accessibility.AccessibilitySettings
	width    int
	height   int
	focused  int // Currently focused setting
}

// Setting represents a toggleable accessibility setting
type Setting struct {
	Name        string
	Description string
	Enabled     *bool
	OnToggle    func()
}

// NewAccessibilityDialog creates a new accessibility settings dialog
func NewAccessibilityDialog(app *app.App) AccessibilityDialog {
	return &accessibilityDialog{
		app:      app,
		settings: accessibility.GetSettings(),
		focused:  0,
	}
}

func (d *accessibilityDialog) Init() tea.Cmd {
	return nil
}

func (d *accessibilityDialog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		d.width = msg.Width
		d.height = msg.Height

	case tea.KeyPressMsg:
		settings := d.getSettings()
		switch msg.String() {
		case "up", "k":
			if d.focused > 0 {
				d.focused--
			}
		case "down", "j":
			if d.focused < len(settings)-1 {
				d.focused++
			}
		case "enter", " ":
			// Toggle focused setting
			if d.focused < len(settings) {
				setting := settings[d.focused]
				*setting.Enabled = !*setting.Enabled
				if setting.OnToggle != nil {
					setting.OnToggle()
				}

				// Announce change for screen readers
				status := "disabled"
				if *setting.Enabled {
					status = "enabled"
				}
				accessibility.AnnounceAction("Toggled "+setting.Name, status)
			}
		}
	}

	return d, nil
}

func (d *accessibilityDialog) View() string {
	t := theme.CurrentTheme()
	typo := typography.New()

	var sections []string

	// Header
	header := styles.NewStyle().
		Foreground(t.Primary()).
		Bold(true).
		Render("♿ Accessibility Settings")

	sections = append(sections, header)
	sections = append(sections, "")

	// Subtitle
	subtitle := styles.NewStyle().
		Foreground(t.TextMuted()).
		Faint(true).
		Render("Customize RyCode for your needs")

	sections = append(sections, subtitle)
	sections = append(sections, "")

	// Settings categories
	sections = append(sections, typo.Subheading.Render("Visual"))
	sections = append(sections, "")

	settings := d.getSettings()
	for i, setting := range settings {
		card := d.renderSetting(setting, i == d.focused)
		sections = append(sections, card)
	}

	sections = append(sections, "")

	// Help footer
	helpStyle := styles.NewStyle().
		Foreground(t.TextMuted()).
		Faint(true)

	help := helpStyle.Render("↑/↓: Navigate  [Space/Enter] Toggle  [ESC] Close")
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

// getSettings returns all available settings
func (d *accessibilityDialog) getSettings() []Setting {
	return []Setting{
		{
			Name:        "High Contrast Mode",
			Description: "Increase color contrast for better visibility",
			Enabled:     &d.settings.HighContrast,
			OnToggle: func() {
				if d.settings.HighContrast {
					d.settings.EnableHighContrast()
				} else {
					d.settings.DisableHighContrast()
				}
			},
		},
		{
			Name:        "Reduced Motion",
			Description: "Minimize animations and motion effects",
			Enabled:     &d.settings.ReducedMotion,
			OnToggle: func() {
				if d.settings.ReducedMotion {
					d.settings.EnableReducedMotion()
				} else {
					d.settings.DisableReducedMotion()
				}
			},
		},
		{
			Name:        "Large Text",
			Description: "Increase text size for better readability",
			Enabled:     &d.settings.LargeText,
		},
		{
			Name:        "Increased Spacing",
			Description: "Add more space between UI elements",
			Enabled:     &d.settings.IncreasedSpacing,
		},
		{
			Name:        "Screen Reader Mode",
			Description: "Optimize for screen reader usage with verbose labels",
			Enabled:     &d.settings.ScreenReaderMode,
			OnToggle: func() {
				if d.settings.ScreenReaderMode {
					d.settings.EnableScreenReader()
				}
			},
		},
		{
			Name:        "Keyboard-Only Mode",
			Description: "Enhance keyboard navigation with visual focus indicators",
			Enabled:     &d.settings.KeyboardOnly,
			OnToggle: func() {
				if d.settings.KeyboardOnly {
					d.settings.EnableKeyboardOnly()
				}
			},
		},
		{
			Name:        "Show Keyboard Hints",
			Description: "Display keyboard shortcuts in dialogs and menus",
			Enabled:     &d.settings.ShowKeyboardHints,
		},
		{
			Name:        "Verbose Labels",
			Description: "Show detailed labels and descriptions everywhere",
			Enabled:     &d.settings.VerboseLabels,
		},
		{
			Name:        "Enhanced Focus Indicators",
			Description: "Make focus indicators larger and more visible",
			Enabled:     &d.settings.EnhancedFocus,
		},
	}
}

// renderSetting renders a single setting card
func (d *accessibilityDialog) renderSetting(setting Setting, focused bool) string {
	t := theme.CurrentTheme()

	// Border color based on focus
	borderColor := t.Border()
	if focused {
		borderColor = t.Primary()
	}

	// Toggle indicator
	toggleStyle := styles.NewStyle().
		Foreground(t.Background()).
		Bold(true).
		Padding(0, 1)

	var toggle string
	if *setting.Enabled {
		toggle = toggleStyle.
			Background(t.Success()).
			Render("✓ ON")
	} else {
		toggle = toggleStyle.
			Background(t.TextMuted()).
			Render("○ OFF")
	}

	// Name
	nameStyle := styles.NewStyle().
		Foreground(t.Text()).
		Bold(true)

	name := nameStyle.Render(setting.Name)

	// Description
	descStyle := styles.NewStyle().
		Foreground(t.TextMuted()).
		Faint(true)

	desc := descStyle.Render(setting.Description)

	// Layout
	line1 := name + "  " + toggle
	line2 := desc

	content := line1 + "\n" + line2

	// Card style
	cardStyle := styles.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(0, 2).
		Width(60)

	return cardStyle.Render(content)
}

func (d *accessibilityDialog) Render(background string) string {
	return d.View()
}

func (d *accessibilityDialog) Close() tea.Cmd {
	return nil
}

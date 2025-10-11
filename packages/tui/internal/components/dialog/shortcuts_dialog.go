package dialog

import (
	"strings"

	"github.com/aaronmrosenthal/rycode/internal/app"
	"github.com/aaronmrosenthal/rycode/internal/layout"
	"github.com/aaronmrosenthal/rycode/internal/styles"
	"github.com/aaronmrosenthal/rycode/internal/theme"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

// ShortcutCategory represents a group of related shortcuts
type ShortcutCategory struct {
	Name      string
	Icon      string
	Shortcuts []Shortcut
}

// Shortcut represents a single keyboard shortcut
type Shortcut struct {
	Keys        string
	Description string
	IsImportant bool // Highlight as essential shortcut
}

// ShortcutsDialog displays comprehensive keyboard shortcuts
type ShortcutsDialog interface {
	layout.Modal
}

type shortcutsDialog struct {
	app        *app.App
	categories []ShortcutCategory
	width      int
	height     int
}

// NewShortcutsDialog creates a new keyboard shortcuts reference
func NewShortcutsDialog(app *app.App) ShortcutsDialog {
	dialog := &shortcutsDialog{
		app: app,
	}

	dialog.setupCategories()

	return dialog
}

// setupCategories defines all keyboard shortcuts by category
func (s *shortcutsDialog) setupCategories() {
	s.categories = []ShortcutCategory{
		{
			Name: "Essential",
			Icon: "⭐",
			Shortcuts: []Shortcut{
				{"Tab", "Cycle through available models", true},
				{"Ctrl+M", "Open model selector", true},
				{"Ctrl+?", "Show this shortcuts guide", true},
				{"Ctrl+C", "Exit RyCode", true},
				{"ESC", "Close current dialog", true},
			},
		},
		{
			Name: "Navigation",
			Icon: "🧭",
			Shortcuts: []Shortcut{
				{"↑/↓ or j/k", "Navigate lists", false},
				{"←/→ or h/l", "Navigate steps", false},
				{"Enter", "Select/Confirm", false},
				{"/", "Search/Filter", false},
				{"Home/End", "Jump to first/last", false},
			},
		},
		{
			Name: "Models & Providers",
			Icon: "🤖",
			Shortcuts: []Shortcut{
				{"Tab", "Quick model switch", true},
				{"Ctrl+M", "Model selector dialog", false},
				{"Ctrl+P", "Provider management", false},
				{"a", "Authenticate provider (in provider list)", false},
				{"r", "Refresh provider status", false},
				{"d", "Auto-detect credentials", false},
			},
		},
		{
			Name: "Analytics & Insights",
			Icon: "📊",
			Shortcuts: []Shortcut{
				{"Ctrl+I", "Usage insights dashboard", true},
				{"Ctrl+B", "Budget forecast", true},
				{"Ctrl+$", "Cost summary", false},
				{"i", "Toggle AI recommendations (in model selector)", false},
			},
		},
		{
			Name: "Editing",
			Icon: "✏️",
			Shortcuts: []Shortcut{
				{"Ctrl+A", "Select all", false},
				{"Ctrl+X", "Cut", false},
				{"Ctrl+V", "Paste", false},
				{"Ctrl+Z", "Undo", false},
				{"Ctrl+Y", "Redo", false},
			},
		},
		{
			Name: "Advanced",
			Icon: "🔧",
			Shortcuts: []Shortcut{
				{"Ctrl+R", "Refresh all data", false},
				{"Ctrl+T", "Toggle theme", false},
				{"Ctrl+L", "Clear screen", false},
				{"Ctrl+D", "Developer tools", false},
			},
		},
	}
}

func (s *shortcutsDialog) Init() tea.Cmd {
	return nil
}

func (s *shortcutsDialog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
	}

	return s, nil
}

func (s *shortcutsDialog) View() string {
	t := theme.CurrentTheme()

	var sections []string

	// Header
	header := styles.NewStyle().
		Foreground(t.Primary()).
		Bold(true).
		Render("⌨️  Keyboard Shortcuts")

	sections = append(sections, header)
	sections = append(sections, "")

	// Subtitle
	subtitle := styles.NewStyle().
		Foreground(t.TextMuted()).
		Faint(true).
		Render("Master RyCode with these keyboard shortcuts")

	sections = append(sections, subtitle)
	sections = append(sections, "")

	// Render categories in two columns
	leftColumn := []ShortcutCategory{}
	rightColumn := []ShortcutCategory{}

	for i, cat := range s.categories {
		if i%2 == 0 {
			leftColumn = append(leftColumn, cat)
		} else {
			rightColumn = append(rightColumn, cat)
		}
	}

	// Render left column
	leftContent := s.renderColumn(leftColumn)
	rightContent := s.renderColumn(rightColumn)

	// Combine columns
	combined := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftContent,
		strings.Repeat(" ", 4),
		rightContent,
	)

	sections = append(sections, combined)
	sections = append(sections, "")

	// Footer tip
	tipStyle := styles.NewStyle().
		Foreground(t.Info()).
		Italic(true)

	tip := tipStyle.Render("💡 Tip: Essential shortcuts are marked with ⭐")
	sections = append(sections, tip)
	sections = append(sections, "")

	// Help text
	helpStyle := styles.NewStyle().
		Foreground(t.TextMuted()).
		Faint(true)

	help := helpStyle.Render("Press ESC to close this guide")
	sections = append(sections, help)

	content := strings.Join(sections, "\n")

	// Wrap in bordered box
	boxStyle := styles.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Border()).
		Padding(2, 3).
		Width(s.width - 4)

	return boxStyle.Render(content)
}

// renderColumn renders a column of shortcut categories
func (s *shortcutsDialog) renderColumn(categories []ShortcutCategory) string {
	t := theme.CurrentTheme()

	var sections []string

	for _, category := range categories {
		// Category header
		categoryHeader := styles.NewStyle().
			Foreground(t.Primary()).
			Bold(true).
			Render(category.Icon + " " + category.Name)

		sections = append(sections, categoryHeader)
		sections = append(sections, "")

		// Shortcuts
		for _, shortcut := range category.Shortcuts {
			line := s.renderShortcut(shortcut)
			sections = append(sections, line)
		}

		sections = append(sections, "")
	}

	return strings.Join(sections, "\n")
}

// renderShortcut renders a single shortcut line
func (s *shortcutsDialog) renderShortcut(shortcut Shortcut) string {
	t := theme.CurrentTheme()

	// Key style
	keyStyle := styles.NewStyle().
		Foreground(t.Background()).
		Background(t.Primary()).
		Bold(true).
		Padding(0, 1)

	if !shortcut.IsImportant {
		// Non-essential shortcuts get muted background
		keyStyle = keyStyle.
			Background(t.TextMuted())
	}

	key := keyStyle.Render(shortcut.Keys)

	// Description style
	descStyle := styles.NewStyle().
		Foreground(t.Text())

	if !shortcut.IsImportant {
		descStyle = descStyle.Foreground(t.TextMuted())
	}

	desc := descStyle.Render(shortcut.Description)

	// Star for important shortcuts
	star := ""
	if shortcut.IsImportant {
		star = styles.NewStyle().
			Foreground(t.Warning()).
			Render(" ⭐")
	}

	return key + "  " + desc + star
}

func (s *shortcutsDialog) Render(background string) string {
	return s.View()
}

func (s *shortcutsDialog) Close() tea.Cmd {
	return nil
}

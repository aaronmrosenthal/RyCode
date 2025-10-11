package dialog

import (
	"time"

	"github.com/aaronmrosenthal/rycode/internal/app"
	"github.com/aaronmrosenthal/rycode/internal/intelligence"
	"github.com/aaronmrosenthal/rycode/internal/layout"
	"github.com/aaronmrosenthal/rycode/internal/styles"
	"github.com/aaronmrosenthal/rycode/internal/theme"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

// InsightsDialog displays usage analytics and insights
type InsightsDialog interface {
	layout.Modal
}

type insightsDialog struct {
	app     *app.App
	insights *intelligence.UsageInsights
	width   int
	height  int
}

// NewInsightsDialog creates a new usage insights dialog
func NewInsightsDialog(app *app.App) InsightsDialog {
	dialog := &insightsDialog{
		app:     app,
		insights: intelligence.NewUsageInsights(),
	}

	// TODO: Load actual usage data from app state
	// For now, creating sample data for demonstration
	dialog.loadSampleData()

	return dialog
}

// loadSampleData creates sample usage data for demonstration
func (d *insightsDialog) loadSampleData() {
	// Sample data showing realistic usage patterns
	// In production, this would load from persistent storage

	// Simulated 7 days of usage
	// Day 1: Light usage
	d.insights.AddUsage(d.daysAgo(6), 0.45, 15, 12000, "claude-3-5-haiku", "anthropic")
	d.insights.AddUsage(d.daysAgo(6), 0.25, 8, 8000, "gpt-3.5-turbo", "openai")

	// Day 2: Moderate usage
	d.insights.AddUsage(d.daysAgo(5), 1.20, 35, 45000, "claude-3-5-sonnet", "anthropic")
	d.insights.AddUsage(d.daysAgo(5), 0.60, 20, 18000, "claude-3-5-haiku", "anthropic")

	// Day 3: Heavy usage
	d.insights.AddUsage(d.daysAgo(4), 2.50, 50, 80000, "claude-3-5-sonnet", "anthropic")
	d.insights.AddUsage(d.daysAgo(4), 1.10, 40, 35000, "gpt-4-turbo", "openai")

	// Day 4: Light usage
	d.insights.AddUsage(d.daysAgo(3), 0.35, 12, 10000, "claude-3-5-haiku", "anthropic")

	// Day 5: Moderate usage
	d.insights.AddUsage(d.daysAgo(2), 1.80, 45, 55000, "claude-3-5-sonnet", "anthropic")
	d.insights.AddUsage(d.daysAgo(2), 0.40, 15, 12000, "gemini-1.5-flash", "google")

	// Day 6: Heavy usage
	d.insights.AddUsage(d.daysAgo(1), 3.20, 65, 95000, "claude-3-5-sonnet", "anthropic")
	d.insights.AddUsage(d.daysAgo(1), 0.80, 25, 22000, "claude-3-5-haiku", "anthropic")

	// Day 7 (today): Moderate usage so far
	d.insights.AddUsage(d.daysAgo(0), 1.50, 38, 48000, "claude-3-5-sonnet", "anthropic")
}

// daysAgo returns a time N days in the past
func (d *insightsDialog) daysAgo(days int) time.Time {
	return time.Now().AddDate(0, 0, -days)
}

func (i *insightsDialog) Init() tea.Cmd {
	return nil
}

func (i *insightsDialog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		i.width = msg.Width
		i.height = msg.Height
	}
	return i, nil
}

func (i *insightsDialog) View() string {
	t := theme.CurrentTheme()

	// Render dashboard
	dashboard := i.insights.RenderDashboard(i.width - 8)

	// Add footer with help text
	footer := styles.NewStyle().
		Foreground(t.TextMuted()).
		Faint(true).
		MarginTop(2).
		Render("Press ESC to close")

	content := dashboard + "\n" + footer

	// Wrap in bordered box
	boxStyle := styles.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Border()).
		Padding(2, 3).
		Width(i.width - 4)

	return boxStyle.Render(content)
}

func (i *insightsDialog) Render(background string) string {
	// Center the dialog
	dialog := i.View()

	// Simple centering
	return dialog
}

func (i *insightsDialog) Close() tea.Cmd {
	return nil
}

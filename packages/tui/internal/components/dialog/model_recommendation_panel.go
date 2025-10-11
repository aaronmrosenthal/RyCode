package dialog

import (
	"fmt"
	"strings"
	"time"

	"github.com/aaronmrosenthal/rycode/internal/intelligence"
	"github.com/aaronmrosenthal/rycode/internal/styles"
	"github.com/aaronmrosenthal/rycode/internal/theme"
	"github.com/aaronmrosenthal/rycode/internal/typography"
	"github.com/charmbracelet/lipgloss/v2"
)

// ModelRecommendationPanel shows intelligent model recommendations
type ModelRecommendationPanel struct {
	engine          *intelligence.RecommendationEngine
	recommendations []intelligence.ModelRecommendation
	width           int
	visible         bool
}

// NewModelRecommendationPanel creates a new recommendation panel
func NewModelRecommendationPanel() *ModelRecommendationPanel {
	return &ModelRecommendationPanel{
		engine:  intelligence.NewRecommendationEngine(),
		visible: true,
	}
}

// SetWidth sets the panel width
func (p *ModelRecommendationPanel) SetWidth(width int) {
	p.width = width
}

// SetVisible controls panel visibility
func (p *ModelRecommendationPanel) SetVisible(visible bool) {
	p.visible = visible
}

// GenerateRecommendations generates recommendations based on context
func (p *ModelRecommendationPanel) GenerateRecommendations(ctx intelligence.TaskContext) {
	p.recommendations = p.engine.GetRecommendations(ctx)
}

// View renders the recommendation panel
func (p *ModelRecommendationPanel) View() string {
	if !p.visible || len(p.recommendations) == 0 {
		return ""
	}

	t := theme.CurrentTheme()

	// Panel title
	titleStyle := styles.NewStyle().
		Foreground(t.Primary()).
		Bold(true).
		Padding(0, 1)

	title := titleStyle.Render("💡 AI Recommendations")

	// Build recommendation cards
	var cards []string
	for i, rec := range p.recommendations {
		if i >= 3 {
			break // Show top 3
		}
		cards = append(cards, p.renderCompactRecommendation(rec, i == 0))
	}

	content := strings.Join(cards, "\n")

	// Wrap in panel
	panelStyle := styles.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Border()).
		Padding(1, 2).
		Width(p.width)

	return panelStyle.Render(title + "\n\n" + content)
}

// renderCompactRecommendation renders a compact recommendation card
func (p *ModelRecommendationPanel) renderCompactRecommendation(rec intelligence.ModelRecommendation, highlighted bool) string {
	t := theme.CurrentTheme()
	typo := typography.New()

	// Model name and score
	modelStyle := styles.NewStyle()
	if highlighted {
		modelStyle = modelStyle.
			Foreground(t.Background()).
			Background(t.Primary()).
			Bold(true)
	} else {
		modelStyle = modelStyle.
			Foreground(t.Text()).
			Bold(true)
	}

	modelName := modelStyle.Render(fmt.Sprintf("%s - %s", rec.Provider, rec.Model))

	// Score badge
	scoreColor := t.Success()
	if rec.Score < 80 {
		scoreColor = t.Warning()
	}
	if rec.Score < 60 {
		scoreColor = t.Error()
	}

	scoreBadge := styles.NewStyle().
		Foreground(t.Background()).
		Background(scoreColor).
		Bold(true).
		Padding(0, 1).
		Render(fmt.Sprintf("%.0f%%", rec.Score))

	// Metadata badges
	costBadge := styles.NewStyle().
		Foreground(t.TextMuted()).
		Render(fmt.Sprintf("$%.3f", rec.CostPerUse))

	speedIcon := "⚡"
	if rec.Speed == "medium" {
		speedIcon = "→"
	} else if rec.Speed == "slow" {
		speedIcon = "🐌"
	}
	speedBadge := styles.NewStyle().
		Foreground(t.TextMuted()).
		Render(fmt.Sprintf("%s %s", speedIcon, rec.Speed))

	qualityIcon := "⭐"
	if rec.Quality == "medium" {
		qualityIcon = "★"
	} else if rec.Quality == "basic" {
		qualityIcon = "☆"
	}
	qualityBadge := styles.NewStyle().
		Foreground(t.TextMuted()).
		Render(fmt.Sprintf("%s %s", qualityIcon, rec.Quality))

	metadata := typography.Inline([]string{costBadge, speedBadge, qualityBadge}, " • ")

	// Reasoning
	reasoning := typo.Caption.
		Width(p.width - 8).
		Render("💭 " + rec.Reasoning)

	// Top pros (show 2)
	pros := ""
	for i, pro := range rec.Pros {
		if i >= 2 {
			break
		}
		pros += typo.Caption.
			Foreground(t.Success()).
			Render(fmt.Sprintf("  ✓ %s", pro)) + "\n"
	}

	// Build card
	var parts []string
	parts = append(parts, modelName+" "+scoreBadge)
	parts = append(parts, metadata)
	parts = append(parts, reasoning)
	if pros != "" {
		parts = append(parts, strings.TrimSuffix(pros, "\n"))
	}

	card := strings.Join(parts, "\n")

	// Add separator if not highlighted
	if !highlighted {
		separator := styles.NewStyle().
			Foreground(t.Border()).
			Render("─────────────────────")
		card = card + "\n" + separator
	}

	return card
}

// GetDefaultContext creates a default task context based on time of day
func GetDefaultContext() intelligence.TaskContext {
	now := time.Now()
	hour := now.Hour()

	// Default to balanced approach
	priority := "balanced"
	complexity := "medium"
	budget := 0.10 // $0.10 per task

	// Work hours (9am-5pm): prioritize quality
	if hour >= 9 && hour <= 17 {
		priority = "quality"
		budget = 0.20
	}

	// After hours: prioritize cost
	if hour < 9 || hour > 17 {
		priority = "cost"
		budget = 0.05
	}

	return intelligence.TaskContext{
		Description: "General coding task",
		Complexity:  complexity,
		Priority:    priority,
		Budget:      budget,
		TimeOfDay:   now,
		RecentUsage: []intelligence.ModelUsage{},
	}
}

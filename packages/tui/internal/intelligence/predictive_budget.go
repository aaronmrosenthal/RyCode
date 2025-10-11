package intelligence

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/aaronmrosenthal/rycode/internal/styles"
	"github.com/aaronmrosenthal/rycode/internal/theme"
	"github.com/aaronmrosenthal/rycode/internal/typography"
	"github.com/charmbracelet/lipgloss/v2/compat"
)

// BudgetForecast represents predicted spending
type BudgetForecast struct {
	CurrentSpend      float64
	ProjectedMonthEnd float64
	Confidence        float64 // 0-100
	DaysRemaining     int
	DaysElapsed       int
	AverageDailySpend float64
	Trend             string // "increasing", "stable", "decreasing"
	WillExceedBudget  bool
	ExcessAmount      float64
}

// BudgetRecommendation provides actionable budget guidance
type BudgetRecommendation struct {
	Type        string // "stay_course", "reduce_usage", "upgrade_budget", "critical"
	Title       string
	Message     string
	Suggestions []string
	DailyTarget float64 // Suggested daily spend to stay on budget
}

// PredictiveBudget provides ML-style budget forecasting
type PredictiveBudget struct {
	monthlyLimit      float64
	dailyLimit        float64
	currentMonthSpend float64
	dailySpendHistory []float64
	usageInsights     *UsageInsights
}

// NewPredictiveBudget creates a new predictive budget system
func NewPredictiveBudget(monthlyLimit, dailyLimit float64, insights *UsageInsights) *PredictiveBudget {
	return &PredictiveBudget{
		monthlyLimit:      monthlyLimit,
		dailyLimit:        dailyLimit,
		usageInsights:     insights,
		dailySpendHistory: make([]float64, 0),
	}
}

// RecordDailySpend records spending for a day
func (p *PredictiveBudget) RecordDailySpend(amount float64) {
	p.dailySpendHistory = append(p.dailySpendHistory, amount)
	p.currentMonthSpend += amount
}

// GetForecast generates a budget forecast for the current month
func (p *PredictiveBudget) GetForecast() BudgetForecast {
	now := time.Now()
	daysElapsed := now.Day()
	daysInMonth := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, now.Location()).Day()
	daysRemaining := daysInMonth - daysElapsed

	// Calculate average daily spend
	avgDailySpend := 0.0
	if daysElapsed > 0 {
		avgDailySpend = p.currentMonthSpend / float64(daysElapsed)
	}

	// Simple linear projection
	projectedMonthEnd := p.currentMonthSpend + (avgDailySpend * float64(daysRemaining))

	// Calculate trend with weighted recent days
	trend := p.calculateTrend()

	// Adjust projection based on trend
	if trend == "increasing" {
		projectedMonthEnd *= 1.15 // 15% increase
	} else if trend == "decreasing" {
		projectedMonthEnd *= 0.90 // 10% decrease
	}

	// Calculate confidence based on data points
	confidence := math.Min(100, float64(daysElapsed)*10) // Max 100% confidence

	// Check if will exceed budget
	willExceed := projectedMonthEnd > p.monthlyLimit
	excessAmount := 0.0
	if willExceed {
		excessAmount = projectedMonthEnd - p.monthlyLimit
	}

	return BudgetForecast{
		CurrentSpend:      p.currentMonthSpend,
		ProjectedMonthEnd: projectedMonthEnd,
		Confidence:        confidence,
		DaysRemaining:     daysRemaining,
		DaysElapsed:       daysElapsed,
		AverageDailySpend: avgDailySpend,
		Trend:             trend,
		WillExceedBudget:  willExceed,
		ExcessAmount:      excessAmount,
	}
}

// calculateTrend analyzes spending trend over recent days
func (p *PredictiveBudget) calculateTrend() string {
	if len(p.dailySpendHistory) < 3 {
		return "stable"
	}

	// Compare last 3 days with previous 3 days
	recentDays := 3
	if len(p.dailySpendHistory) < 6 {
		return "stable"
	}

	recentAvg := 0.0
	previousAvg := 0.0

	for i := 0; i < recentDays; i++ {
		recentAvg += p.dailySpendHistory[len(p.dailySpendHistory)-1-i]
		previousAvg += p.dailySpendHistory[len(p.dailySpendHistory)-1-recentDays-i]
	}

	recentAvg /= float64(recentDays)
	previousAvg /= float64(recentDays)

	// Calculate percentage change
	change := (recentAvg - previousAvg) / previousAvg

	if change > 0.15 {
		return "increasing"
	} else if change < -0.15 {
		return "decreasing"
	}

	return "stable"
}

// GetRecommendation generates budget recommendations based on forecast
func (p *PredictiveBudget) GetRecommendation() BudgetRecommendation {
	forecast := p.GetForecast()

	// Calculate daily target to stay on budget
	remaining := p.monthlyLimit - forecast.CurrentSpend
	dailyTarget := 0.0
	if forecast.DaysRemaining > 0 {
		dailyTarget = remaining / float64(forecast.DaysRemaining)
	}

	// Critical: Already exceeded budget
	if forecast.CurrentSpend >= p.monthlyLimit {
		return BudgetRecommendation{
			Type:        "critical",
			Title:       "🚨 Budget Exceeded",
			Message:     fmt.Sprintf("You've spent $%.2f of your $%.2f monthly budget.", forecast.CurrentSpend, p.monthlyLimit),
			Suggestions: []string{
				"Stop using AI until next month",
				"Switch to free tier models only",
				"Increase your monthly budget",
				"Review and delete unnecessary usage",
			},
			DailyTarget: 0,
		}
	}

	// Projected to exceed by >20%
	if forecast.WillExceedBudget && (forecast.ExcessAmount/p.monthlyLimit) > 0.20 {
		return BudgetRecommendation{
			Type:  "reduce_usage",
			Title: "⚠️ Significant Overrun Projected",
			Message: fmt.Sprintf("At current pace, you'll spend $%.2f this month ($%.2f over budget).",
				forecast.ProjectedMonthEnd, forecast.ExcessAmount),
			Suggestions: []string{
				fmt.Sprintf("Reduce daily spend to $%.2f to stay on budget", dailyTarget),
				"Use Claude Haiku instead of Sonnet for simple tasks",
				"Use GPT-3.5 Turbo instead of GPT-4",
				"Batch requests to reduce API calls",
				fmt.Sprintf("Consider increasing monthly budget to $%.2f", forecast.ProjectedMonthEnd*1.1),
			},
			DailyTarget: dailyTarget,
		}
	}

	// Projected to slightly exceed
	if forecast.WillExceedBudget {
		return BudgetRecommendation{
			Type:  "reduce_usage",
			Title: "⚡ Budget Overrun Possible",
			Message: fmt.Sprintf("You may exceed budget by $%.2f this month.",
				forecast.ExcessAmount),
			Suggestions: []string{
				fmt.Sprintf("Target $%.2f/day to stay within budget", dailyTarget),
				"Switch to cheaper models for routine tasks",
				"Monitor usage more closely",
			},
			DailyTarget: dailyTarget,
		}
	}

	// Using less than 50% of budget with <5 days left
	percentUsed := (forecast.CurrentSpend / p.monthlyLimit) * 100
	if percentUsed < 50 && forecast.DaysRemaining < 5 {
		return BudgetRecommendation{
			Type:  "upgrade_budget",
			Title: "💡 Budget Underutilized",
			Message: fmt.Sprintf("You've only used %.0f%% of your budget with %d days left.",
				percentUsed, forecast.DaysRemaining),
			Suggestions: []string{
				"You can safely use premium models more often",
				"Consider reducing monthly budget to save money",
				fmt.Sprintf("Your current pace: $%.2f/month (vs $%.2f budget)", forecast.ProjectedMonthEnd, p.monthlyLimit),
			},
			DailyTarget: dailyTarget,
		}
	}

	// On track
	return BudgetRecommendation{
		Type:  "stay_course",
		Title: "✓ Budget On Track",
		Message: fmt.Sprintf("Projected: $%.2f of $%.2f budget (%.0f%% confidence).",
			forecast.ProjectedMonthEnd, p.monthlyLimit, forecast.Confidence),
		Suggestions: []string{
			fmt.Sprintf("Continue current usage pattern ($%.2f/day)", forecast.AverageDailySpend),
			"Budget management looks good",
		},
		DailyTarget: dailyTarget,
	}
}

// RenderForecast creates a beautiful forecast visualization
func (p *PredictiveBudget) RenderForecast(width int) string {
	forecast := p.GetForecast()
	recommendation := p.GetRecommendation()

	t := theme.CurrentTheme()

	var sections []string

	// Header
	header := styles.NewStyle().
		Foreground(t.Primary()).
		Bold(true).
		Render("🔮 Budget Forecast")

	sections = append(sections, header)
	sections = append(sections, "")

	// Current status
	currentStatus := p.renderCurrentStatus(forecast)
	sections = append(sections, currentStatus)
	sections = append(sections, "")

	// Projection visualization
	projectionViz := p.renderProjectionBar(forecast)
	sections = append(sections, projectionViz)
	sections = append(sections, "")

	// Trend indicator
	trendLine := p.renderTrend(forecast)
	sections = append(sections, trendLine)
	sections = append(sections, "")

	// Recommendation card
	recCard := p.renderRecommendation(recommendation)
	sections = append(sections, recCard)

	return strings.Join(sections, "\n")
}

// renderCurrentStatus displays current month status
func (p *PredictiveBudget) renderCurrentStatus(forecast BudgetForecast) string {
	t := theme.CurrentTheme()

	// Current spend
	currentStyle := styles.NewStyle().
		Foreground(t.Info()).
		Bold(true)

	current := currentStyle.Render(fmt.Sprintf("$%.2f", forecast.CurrentSpend))

	// Budget
	budgetStyle := styles.NewStyle().
		Foreground(t.TextMuted())

	budget := budgetStyle.Render(fmt.Sprintf("/ $%.2f", p.monthlyLimit))

	// Percentage
	percentUsed := (forecast.CurrentSpend / p.monthlyLimit) * 100
	percentColor := t.Success()
	if percentUsed > 80 {
		percentColor = t.Warning()
	}
	if percentUsed > 95 {
		percentColor = t.Error()
	}

	percentStyle := styles.NewStyle().
		Foreground(percentColor).
		Bold(true)

	percent := percentStyle.Render(fmt.Sprintf("(%.0f%%)", percentUsed))

	// Days info
	daysStyle := styles.NewStyle().
		Foreground(t.TextMuted()).
		Faint(true)

	days := daysStyle.Render(fmt.Sprintf("Day %d of ~%d", forecast.DaysElapsed, forecast.DaysElapsed+forecast.DaysRemaining))

	return current + " " + budget + " " + percent + "     " + days
}

// renderProjectionBar shows projected spending as a visual bar
func (p *PredictiveBudget) renderProjectionBar(forecast BudgetForecast) string {
	t := theme.CurrentTheme()

	barWidth := 40
	currentBar := int((forecast.CurrentSpend / p.monthlyLimit) * float64(barWidth))
	projectedBar := int((forecast.ProjectedMonthEnd / p.monthlyLimit) * float64(barWidth))

	// Clamp to bar width
	if currentBar > barWidth {
		currentBar = barWidth
	}
	if projectedBar > barWidth {
		projectedBar = barWidth
	}

	// Build bar
	bar := ""

	// Current spend (solid)
	currentColor := t.Success()
	if forecast.CurrentSpend/p.monthlyLimit > 0.80 {
		currentColor = t.Warning()
	}
	if forecast.CurrentSpend/p.monthlyLimit > 0.95 {
		currentColor = t.Error()
	}

	for i := 0; i < currentBar; i++ {
		bar += styles.NewStyle().Foreground(currentColor).Render("█")
	}

	// Projected spend (outlined)
	projectedColor := t.Info()
	if forecast.WillExceedBudget {
		projectedColor = t.Error()
	}

	for i := currentBar; i < projectedBar; i++ {
		bar += styles.NewStyle().Foreground(projectedColor).Render("░")
	}

	// Empty space
	for i := projectedBar; i < barWidth; i++ {
		bar += styles.NewStyle().Foreground(t.TextMuted()).Faint(true).Render("·")
	}

	// Labels
	labelStyle := styles.NewStyle().
		Foreground(t.TextMuted()).
		Faint(true)

	currentLabel := labelStyle.Render("Current")
	projectedLabel := labelStyle.Render(fmt.Sprintf("Projected: $%.2f", forecast.ProjectedMonthEnd))

	return bar + "\n" + currentLabel + "     " + projectedLabel
}

// renderTrend shows spending trend indicator
func (p *PredictiveBudget) renderTrend(forecast BudgetForecast) string {
	t := theme.CurrentTheme()

	var icon string
	var text string
	var color compat.AdaptiveColor

	switch forecast.Trend {
	case "increasing":
		icon = "📈"
		text = "Spending is increasing"
		color = t.Warning()
	case "decreasing":
		icon = "📉"
		text = "Spending is decreasing"
		color = t.Success()
	default:
		icon = "➡️"
		text = "Spending is stable"
		color = t.Info()
	}

	trendStyle := styles.NewStyle().
		Foreground(color)

	confidence := styles.NewStyle().
		Foreground(t.TextMuted()).
		Faint(true).
		Render(fmt.Sprintf("(%.0f%% confidence)", forecast.Confidence))

	return icon + " " + trendStyle.Render(text) + " " + confidence
}

// renderRecommendation displays the budget recommendation card
func (p *PredictiveBudget) renderRecommendation(rec BudgetRecommendation) string {
	t := theme.CurrentTheme()
	typo := typography.New()

	// Title with icon
	titleColor := t.Info()
	if rec.Type == "critical" {
		titleColor = t.Error()
	} else if rec.Type == "reduce_usage" {
		titleColor = t.Warning()
	} else if rec.Type == "stay_course" {
		titleColor = t.Success()
	}

	title := styles.NewStyle().
		Foreground(titleColor).
		Bold(true).
		Render(rec.Title)

	// Message
	message := typo.Body.
		Foreground(t.Text()).
		Render(rec.Message)

	// Suggestions
	var suggestions []string
	for _, suggestion := range rec.Suggestions {
		bullet := typo.Body.
			Foreground(t.TextMuted()).
			Render("  • " + suggestion)
		suggestions = append(suggestions, bullet)
	}

	// Build card
	parts := []string{title, "", message}
	if len(suggestions) > 0 {
		parts = append(parts, "")
		parts = append(parts, suggestions...)
	}

	return strings.Join(parts, "\n")
}

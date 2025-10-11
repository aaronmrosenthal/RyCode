package intelligence

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/aaronmrosenthal/rycode/internal/styles"
	"github.com/aaronmrosenthal/rycode/internal/theme"
	"github.com/aaronmrosenthal/rycode/internal/typography"
)

// UsageData represents usage statistics for a time period
type UsageData struct {
	Date      time.Time
	Cost      float64
	Requests  int
	Tokens    int64
	Models    map[string]int // Model ID -> usage count
	Providers map[string]int // Provider ID -> usage count
}

// UsageInsights provides analytics and visualization of usage patterns
type UsageInsights struct {
	dailyData   []UsageData
	weeklyData  []UsageData
	monthlyData []UsageData
}

// NewUsageInsights creates a new usage insights analyzer
func NewUsageInsights() *UsageInsights {
	return &UsageInsights{
		dailyData:   make([]UsageData, 0),
		weeklyData:  make([]UsageData, 0),
		monthlyData: make([]UsageData, 0),
	}
}

// AddUsage records a new usage event
func (u *UsageInsights) AddUsage(date time.Time, cost float64, requests int, tokens int64, model, provider string) {
	// Find or create daily entry
	dateKey := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	var dailyEntry *UsageData
	for i := range u.dailyData {
		if u.dailyData[i].Date.Equal(dateKey) {
			dailyEntry = &u.dailyData[i]
			break
		}
	}

	if dailyEntry == nil {
		u.dailyData = append(u.dailyData, UsageData{
			Date:      dateKey,
			Models:    make(map[string]int),
			Providers: make(map[string]int),
		})
		dailyEntry = &u.dailyData[len(u.dailyData)-1]
	}

	// Update stats
	dailyEntry.Cost += cost
	dailyEntry.Requests += requests
	dailyEntry.Tokens += tokens
	dailyEntry.Models[model]++
	dailyEntry.Providers[provider]++

	// Sort daily data by date
	sort.Slice(u.dailyData, func(i, j int) bool {
		return u.dailyData[i].Date.Before(u.dailyData[j].Date)
	})
}

// GetDailyCosts returns costs for the last N days
func (u *UsageInsights) GetDailyCosts(days int) []float64 {
	if len(u.dailyData) == 0 {
		return make([]float64, days)
	}

	costs := make([]float64, days)
	dataIndex := len(u.dailyData) - 1

	for i := days - 1; i >= 0; i-- {
		if dataIndex >= 0 {
			costs[i] = u.dailyData[dataIndex].Cost
			dataIndex--
		}
	}

	return costs
}

// GetTopModels returns the most-used models
func (u *UsageInsights) GetTopModels(limit int) []struct {
	Model string
	Count int
} {
	modelCounts := make(map[string]int)

	for _, day := range u.dailyData {
		for model, count := range day.Models {
			modelCounts[model] += count
		}
	}

	type modelCount struct {
		Model string
		Count int
	}

	models := make([]modelCount, 0, len(modelCounts))
	for model, count := range modelCounts {
		models = append(models, modelCount{Model: model, Count: count})
	}

	sort.Slice(models, func(i, j int) bool {
		return models[i].Count > models[j].Count
	})

	if len(models) > limit {
		models = models[:limit]
	}

	result := make([]struct {
		Model string
		Count int
	}, len(models))

	for i, m := range models {
		result[i].Model = m.Model
		result[i].Count = m.Count
	}

	return result
}

// GetPeakUsageHours returns hours with highest usage
func (u *UsageInsights) GetPeakUsageHours() []int {
	// For now, return simulated data
	// In real implementation, would track actual hour-by-hour usage
	return []int{9, 10, 11, 14, 15, 16} // Typical work hours
}

// GetTotalCost returns total cost across all data
func (u *UsageInsights) GetTotalCost() float64 {
	total := 0.0
	for _, day := range u.dailyData {
		total += day.Cost
	}
	return total
}

// GetTotalRequests returns total requests across all data
func (u *UsageInsights) GetTotalRequests() int {
	total := 0
	for _, day := range u.dailyData {
		total += day.Requests
	}
	return total
}

// RenderDashboard creates a beautiful usage insights dashboard
func (u *UsageInsights) RenderDashboard(width int) string {
	t := theme.CurrentTheme()
	typo := typography.New()

	var sections []string

	// Header
	header := styles.NewStyle().
		Foreground(t.Primary()).
		Bold(true).
		Render("📊 Usage Insights Dashboard")

	sections = append(sections, header)
	sections = append(sections, "")

	// Summary stats
	summary := u.renderSummaryStats()
	sections = append(sections, summary)
	sections = append(sections, "")

	// Cost trend chart (last 7 days)
	costChart := u.renderCostTrendChart(7, width-4)
	if costChart != "" {
		chartTitle := typo.Subheading.Render("💰 Cost Trend (Last 7 Days)")
		sections = append(sections, chartTitle)
		sections = append(sections, costChart)
		sections = append(sections, "")
	}

	// Top models
	topModels := u.renderTopModels(5)
	if topModels != "" {
		modelsTitle := typo.Subheading.Render("🏆 Most Used Models")
		sections = append(sections, modelsTitle)
		sections = append(sections, topModels)
		sections = append(sections, "")
	}

	// Peak usage times
	peakTimes := u.renderPeakUsageTimes()
	if peakTimes != "" {
		timesTitle := typo.Subheading.Render("⏰ Peak Usage Times")
		sections = append(sections, timesTitle)
		sections = append(sections, peakTimes)
		sections = append(sections, "")
	}

	// Cost savings insights
	savings := u.renderCostSavingsInsights()
	if savings != "" {
		savingsTitle := typo.Subheading.Render("💡 Optimization Opportunities")
		sections = append(sections, savingsTitle)
		sections = append(sections, savings)
	}

	return strings.Join(sections, "\n")
}

// renderSummaryStats creates a summary of key metrics
func (u *UsageInsights) renderSummaryStats() string {
	t := theme.CurrentTheme()

	totalCost := u.GetTotalCost()
	totalRequests := u.GetTotalRequests()
	avgCostPerRequest := 0.0
	if totalRequests > 0 {
		avgCostPerRequest = totalCost / float64(totalRequests)
	}

	// Create stat cards
	costCard := styles.NewStyle().
		Foreground(t.Success()).
		Bold(true).
		Render(fmt.Sprintf("$%.2f", totalCost))

	requestsCard := styles.NewStyle().
		Foreground(t.Info()).
		Bold(true).
		Render(fmt.Sprintf("%d", totalRequests))

	avgCard := styles.NewStyle().
		Foreground(t.Warning()).
		Bold(true).
		Render(fmt.Sprintf("$%.4f", avgCostPerRequest))

	// Labels
	labelStyle := styles.NewStyle().
		Foreground(t.TextMuted()).
		Faint(true)

	costLabel := labelStyle.Render("Total Spend")
	requestsLabel := labelStyle.Render("Requests")
	avgLabel := labelStyle.Render("Avg/Request")

	// Layout
	line1 := costCard + "  " + requestsCard + "  " + avgCard
	line2 := costLabel + "     " + requestsLabel + "      " + avgLabel

	return line1 + "\n" + line2
}

// renderCostTrendChart creates an ASCII bar chart of daily costs
func (u *UsageInsights) renderCostTrendChart(days, width int) string {
	costs := u.GetDailyCosts(days)

	// Find max cost for scaling
	maxCost := 0.0
	for _, cost := range costs {
		if cost > maxCost {
			maxCost = cost
		}
	}

	if maxCost == 0 {
		return ""
	}

	t := theme.CurrentTheme()

	// Chart height
	height := 10

	var lines []string

	// Render from top to bottom
	for row := height; row >= 0; row-- {
		threshold := (float64(row) / float64(height)) * maxCost

		line := ""
		for _, cost := range costs {
			if cost >= threshold {
				// Bar character
				bar := styles.NewStyle().
					Foreground(t.Success()).
					Render("█")
				line += bar
			} else {
				line += " "
			}
			line += " "
		}

		// Add y-axis label
		if row == height {
			label := styles.NewStyle().
				Foreground(t.TextMuted()).
				Faint(true).
				Render(fmt.Sprintf("$%.2f", maxCost))
			line = label + " " + line
		} else if row == 0 {
			label := styles.NewStyle().
				Foreground(t.TextMuted()).
				Faint(true).
				Render("$0.00")
			line = label + " " + line
		} else {
			line = "       " + line
		}

		lines = append(lines, line)
	}

	// X-axis labels (day numbers)
	xAxis := "       "
	for i := 0; i < days; i++ {
		label := styles.NewStyle().
			Foreground(t.TextMuted()).
			Faint(true).
			Render(fmt.Sprintf("%d", i+1))
		xAxis += label + " "
	}
	lines = append(lines, xAxis)

	return strings.Join(lines, "\n")
}

// renderTopModels creates a list of most-used models with bar graphs
func (u *UsageInsights) renderTopModels(limit int) string {
	topModels := u.GetTopModels(limit)

	if len(topModels) == 0 {
		return ""
	}

	t := theme.CurrentTheme()

	// Find max count for scaling
	maxCount := 0
	for _, m := range topModels {
		if m.Count > maxCount {
			maxCount = m.Count
		}
	}

	var lines []string

	for i, m := range topModels {
		// Model name
		modelName := m.Model
		if len(modelName) > 30 {
			modelName = modelName[:27] + "..."
		}

		modelStyle := styles.NewStyle().
			Foreground(t.Text()).
			Bold(true)

		// Rank badge
		rank := ""
		if i == 0 {
			rank = "🥇 "
		} else if i == 1 {
			rank = "🥈 "
		} else if i == 2 {
			rank = "🥉 "
		} else {
			rank = fmt.Sprintf("%d. ", i+1)
		}

		nameWithRank := modelStyle.Render(rank + modelName)

		// Bar graph
		barWidth := int((float64(m.Count) / float64(maxCount)) * 20)
		bar := strings.Repeat("█", barWidth)
		barRendered := styles.NewStyle().
			Foreground(t.Primary()).
			Render(bar)

		// Count
		count := styles.NewStyle().
			Foreground(t.TextMuted()).
			Faint(true).
			Render(fmt.Sprintf("(%d)", m.Count))

		line := fmt.Sprintf("%-40s %s %s", nameWithRank, barRendered, count)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

// renderPeakUsageTimes creates a visualization of peak usage hours
func (u *UsageInsights) renderPeakUsageTimes() string {
	peakHours := u.GetPeakUsageHours()

	if len(peakHours) == 0 {
		return ""
	}

	t := theme.CurrentTheme()

	// Create 24-hour bar chart
	hours := make([]bool, 24)
	for _, hour := range peakHours {
		if hour >= 0 && hour < 24 {
			hours[hour] = true
		}
	}

	// Render chart
	chart := ""
	for hour := 0; hour < 24; hour++ {
		if hours[hour] {
			bar := styles.NewStyle().
				Foreground(t.Success()).
				Render("█")
			chart += bar
		} else {
			chart += styles.NewStyle().
				Foreground(t.TextMuted()).
				Faint(true).
				Render("░")
		}
	}

	// Add time labels
	labels := styles.NewStyle().
		Foreground(t.TextMuted()).
		Faint(true).
		Render("0am                  12pm                  11pm")

	return chart + "\n" + labels
}

// renderCostSavingsInsights provides actionable cost-saving suggestions
func (u *UsageInsights) renderCostSavingsInsights() string {
	t := theme.CurrentTheme()
	typo := typography.New()

	var insights []string

	// Analyze patterns and suggest optimizations
	totalCost := u.GetTotalCost()

	if totalCost > 10.0 {
		insight := typo.Body.
			Foreground(t.Text()).
			Render("💡 Consider using Claude Haiku for simple tasks - 5x cheaper than Sonnet")
		insights = append(insights, insight)
	}

	if totalCost > 50.0 {
		potentialSavings := totalCost * 0.3 // 30% potential savings
		insight := typo.Body.
			Foreground(t.Success()).
			Render(fmt.Sprintf("💰 Potential savings: $%.2f/month by optimizing model selection", potentialSavings))
		insights = append(insights, insight)
	}

	topModels := u.GetTopModels(1)
	if len(topModels) > 0 && strings.Contains(topModels[0].Model, "gpt-4") {
		insight := typo.Body.
			Foreground(t.Warning()).
			Render("⚡ Try Claude 3.5 Sonnet - often better quality at lower cost than GPT-4")
		insights = append(insights, insight)
	}

	if len(insights) == 0 {
		insight := typo.Body.
			Foreground(t.TextMuted()).
			Render("✓ Your model usage looks optimized!")
		insights = append(insights, insight)
	}

	return strings.Join(insights, "\n")
}

// GetWeeklySummary returns a summary for the past week
func (u *UsageInsights) GetWeeklySummary() string {
	costs := u.GetDailyCosts(7)

	total := 0.0
	for _, cost := range costs {
		total += cost
	}

	avg := total / 7.0

	// Find highest day
	maxCost := 0.0
	maxDay := 0
	for i, cost := range costs {
		if cost > maxCost {
			maxCost = cost
			maxDay = i + 1
		}
	}

	return fmt.Sprintf("Week: $%.2f total | $%.2f/day avg | Peak: Day %d ($%.2f)",
		total, avg, maxDay, maxCost)
}

// GetMonthlySummary returns a summary for the past month
func (u *UsageInsights) GetMonthlySummary() string {
	costs := u.GetDailyCosts(30)

	total := 0.0
	for _, cost := range costs {
		total += cost
	}

	avg := total / 30.0
	projected := avg * 30.0

	return fmt.Sprintf("Month: $%.2f total | $%.2f/day avg | Projected: $%.2f",
		total, avg, projected)
}

// EstimateSavings estimates potential monthly savings from optimization
func (u *UsageInsights) EstimateSavings() float64 {
	totalCost := u.GetTotalCost()

	// Estimate 20-40% savings possible through smart model selection
	estimatedSavings := totalCost * 0.30

	return math.Max(0, estimatedSavings)
}

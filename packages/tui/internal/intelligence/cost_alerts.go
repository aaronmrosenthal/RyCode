package intelligence

import (
	"fmt"
	"time"

	"github.com/aaronmrosenthal/rycode/internal/components/toast"
	tea "github.com/charmbracelet/bubbletea/v2"
)

// AlertLevel represents the severity of a cost alert
type AlertLevel int

const (
	AlertLevelInfo AlertLevel = iota
	AlertLevelWarning
	AlertLevelCritical
	AlertLevelUrgent
)

// CostAlert represents a cost-related alert
type CostAlert struct {
	Level       AlertLevel
	Title       string
	Message     string
	Suggestion  string
	Threshold   float64
	CurrentCost float64
	Timestamp   time.Time
}

// BudgetConfig holds budget settings
type BudgetConfig struct {
	DailyLimit   float64
	MonthlyLimit float64
	WarningAt    float64 // percentage (e.g., 0.8 for 80%)
	CriticalAt   float64 // percentage (e.g., 0.95 for 95%)
}

// DefaultBudget provides sensible defaults
var DefaultBudget = BudgetConfig{
	DailyLimit:   10.0,  // $10/day
	MonthlyLimit: 50.0,  // $50/month
	WarningAt:    0.80,  // 80% threshold
	CriticalAt:   0.95,  // 95% threshold
}

// CostAlertSystem monitors spending and generates alerts
type CostAlertSystem struct {
	budget         BudgetConfig
	lastDailyCost  float64
	lastMonthlyCost float64
	lastAlert      time.Time
	alertCooldown  time.Duration
}

// NewCostAlertSystem creates a new cost alert system
func NewCostAlertSystem(budget BudgetConfig) *CostAlertSystem {
	return &CostAlertSystem{
		budget:        budget,
		alertCooldown: 15 * time.Minute, // Don't spam alerts
	}
}

// CheckCosts analyzes costs and returns alerts if needed
func (c *CostAlertSystem) CheckCosts(dailyCost, monthlyCost, projectedMonthly float64) []CostAlert {
	alerts := []CostAlert{}

	// Check if we're in cooldown
	if time.Since(c.lastAlert) < c.alertCooldown {
		return alerts
	}

	// Daily budget alerts
	if dailyAlert := c.checkDailyBudget(dailyCost); dailyAlert != nil {
		alerts = append(alerts, *dailyAlert)
	}

	// Monthly budget alerts
	if monthlyAlert := c.checkMonthlyBudget(monthlyCost); monthlyAlert != nil {
		alerts = append(alerts, *monthlyAlert)
	}

	// Projection alerts
	if projectionAlert := c.checkProjection(monthlyCost, projectedMonthly); projectionAlert != nil {
		alerts = append(alerts, *projectionAlert)
	}

	// Spike detection
	if spikeAlert := c.checkSpike(dailyCost, monthlyCost); spikeAlert != nil {
		alerts = append(alerts, *spikeAlert)
	}

	// Update state
	c.lastDailyCost = dailyCost
	c.lastMonthlyCost = monthlyCost
	if len(alerts) > 0 {
		c.lastAlert = time.Now()
	}

	return alerts
}

// checkDailyBudget checks daily spending against limits
func (c *CostAlertSystem) checkDailyBudget(dailyCost float64) *CostAlert {
	if c.budget.DailyLimit <= 0 {
		return nil
	}

	percentage := dailyCost / c.budget.DailyLimit

	if percentage >= 1.0 {
		return &CostAlert{
			Level:       AlertLevelUrgent,
			Title:       "Daily Budget Exceeded!",
			Message:     fmt.Sprintf("You've spent $%.2f today (limit: $%.2f)", dailyCost, c.budget.DailyLimit),
			Suggestion:  "Consider stopping for today or increasing your budget",
			Threshold:   c.budget.DailyLimit,
			CurrentCost: dailyCost,
			Timestamp:   time.Now(),
		}
	}

	if percentage >= c.budget.CriticalAt {
		return &CostAlert{
			Level:       AlertLevelCritical,
			Title:       "Daily Budget Nearly Exceeded",
			Message:     fmt.Sprintf("You've used %.0f%% of today's budget ($%.2f/$%.2f)", percentage*100, dailyCost, c.budget.DailyLimit),
			Suggestion:  "Switch to cheaper models or pause usage",
			Threshold:   c.budget.DailyLimit * c.budget.CriticalAt,
			CurrentCost: dailyCost,
			Timestamp:   time.Now(),
		}
	}

	if percentage >= c.budget.WarningAt {
		return &CostAlert{
			Level:       AlertLevelWarning,
			Title:       "Daily Budget Warning",
			Message:     fmt.Sprintf("You've used %.0f%% of today's budget ($%.2f/$%.2f)", percentage*100, dailyCost, c.budget.DailyLimit),
			Suggestion:  "Consider using Claude Haiku or GPT-3.5 for remaining tasks",
			Threshold:   c.budget.DailyLimit * c.budget.WarningAt,
			CurrentCost: dailyCost,
			Timestamp:   time.Now(),
		}
	}

	return nil
}

// checkMonthlyBudget checks monthly spending against limits
func (c *CostAlertSystem) checkMonthlyBudget(monthlyCost float64) *CostAlert {
	if c.budget.MonthlyLimit <= 0 {
		return nil
	}

	percentage := monthlyCost / c.budget.MonthlyLimit

	if percentage >= 1.0 {
		return &CostAlert{
			Level:       AlertLevelUrgent,
			Title:       "Monthly Budget Exceeded!",
			Message:     fmt.Sprintf("You've spent $%.2f this month (limit: $%.2f)", monthlyCost, c.budget.MonthlyLimit),
			Suggestion:  "Budget exceeded. Consider upgrading or pausing usage.",
			Threshold:   c.budget.MonthlyLimit,
			CurrentCost: monthlyCost,
			Timestamp:   time.Now(),
		}
	}

	if percentage >= c.budget.CriticalAt {
		return &CostAlert{
			Level:       AlertLevelCritical,
			Title:       "Monthly Budget Nearly Exceeded",
			Message:     fmt.Sprintf("You've used %.0f%% of this month's budget ($%.2f/$%.2f)", percentage*100, monthlyCost, c.budget.MonthlyLimit),
			Suggestion:  "Optimize usage to stay within budget",
			Threshold:   c.budget.MonthlyLimit * c.budget.CriticalAt,
			CurrentCost: monthlyCost,
			Timestamp:   time.Now(),
		}
	}

	if percentage >= c.budget.WarningAt {
		return &CostAlert{
			Level:       AlertLevelWarning,
			Title:       "Monthly Budget Warning",
			Message:     fmt.Sprintf("You've used %.0f%% of this month's budget ($%.2f/$%.2f)", percentage*100, monthlyCost, c.budget.MonthlyLimit),
			Suggestion:  "Monitor usage closely to avoid exceeding budget",
			Threshold:   c.budget.MonthlyLimit * c.budget.WarningAt,
			CurrentCost: monthlyCost,
			Timestamp:   time.Now(),
		}
	}

	return nil
}

// checkProjection checks if projected spending will exceed budget
func (c *CostAlertSystem) checkProjection(monthlyCost, projectedMonthly float64) *CostAlert {
	if c.budget.MonthlyLimit <= 0 || projectedMonthly <= c.budget.MonthlyLimit {
		return nil
	}

	// Already over budget, don't send projection alert
	if monthlyCost >= c.budget.MonthlyLimit {
		return nil
	}

	overage := projectedMonthly - c.budget.MonthlyLimit
	percentage := (overage / c.budget.MonthlyLimit) * 100

	if percentage > 20 {
		return &CostAlert{
			Level:       AlertLevelCritical,
			Title:       "Budget Overrun Projected",
			Message:     fmt.Sprintf("At current pace, you'll spend $%.2f this month ($%.2f over budget)", projectedMonthly, overage),
			Suggestion:  "Reduce usage or increase budget to avoid overrun",
			Threshold:   c.budget.MonthlyLimit,
			CurrentCost: projectedMonthly,
			Timestamp:   time.Now(),
		}
	}

	return &CostAlert{
		Level:       AlertLevelWarning,
		Title:       "Projected to Exceed Budget",
		Message:     fmt.Sprintf("At current pace, you'll spend $%.2f this month ($%.2f over budget)", projectedMonthly, overage),
		Suggestion:  "Consider optimizing model selection",
		Threshold:   c.budget.MonthlyLimit,
		CurrentCost: projectedMonthly,
		Timestamp:   time.Now(),
	}
}

// checkSpike detects unusual spending spikes
func (c *CostAlertSystem) checkSpike(dailyCost, monthlyCost float64) *CostAlert {
	// Need historical data to detect spike
	if c.lastDailyCost == 0 || c.lastMonthlyCost == 0 {
		return nil
	}

	// Calculate average daily spend
	daysInMonth := time.Now().Day()
	if daysInMonth == 0 {
		return nil
	}

	avgDaily := monthlyCost / float64(daysInMonth)

	// Spike if today is 3x average
	if dailyCost > avgDaily*3 && dailyCost > 1.0 {
		return &CostAlert{
			Level:       AlertLevelWarning,
			Title:       "Unusual Spending Spike Detected",
			Message:     fmt.Sprintf("Today's cost ($%.2f) is %.1fx your daily average ($%.2f)", dailyCost, dailyCost/avgDaily, avgDaily),
			Suggestion:  "Was this intentional? Consider reviewing recent usage.",
			Threshold:   avgDaily * 3,
			CurrentCost: dailyCost,
			Timestamp:   time.Now(),
		}
	}

	return nil
}

// ToastCmd converts an alert to a toast command
func (a *CostAlert) ToastCmd() tea.Cmd {
	message := a.Title + ": " + a.Message

	switch a.Level {
	case AlertLevelUrgent:
		return toast.NewErrorToast(message)
	case AlertLevelCritical:
		return toast.NewWarningToast(message)
	case AlertLevelWarning:
		return toast.NewWarningToast(message)
	default:
		return toast.NewInfoToast(message)
	}
}

// CostCheckMsg is sent periodically to check costs
type CostCheckMsg struct{}

// CostAlertMsg is sent when an alert is triggered
type CostAlertMsg struct {
	Alerts []CostAlert
}

// SmartCostSuggestion provides intelligent cost-saving suggestions
type SmartCostSuggestion struct {
	Title          string
	CurrentCost    float64
	PotentialCost  float64
	Savings        float64
	Recommendation string
}

// AnalyzeModelCost suggests cheaper alternatives for a task
func AnalyzeModelCost(provider, model string, taskComplexity string) *SmartCostSuggestion {
	// Simplified cost analysis - in real implementation would use actual pricing
	expensiveModels := map[string]bool{
		"gpt-4":                  true,
		"claude-3-opus":          true,
		"claude-3-5-sonnet":      false, // Good balance
		"claude-3-5-haiku":       false, // Cheap
		"gpt-3.5-turbo":          false, // Cheap
	}

	cheapAlternatives := map[string]string{
		"gpt-4":             "gpt-3.5-turbo or claude-3-5-haiku",
		"claude-3-opus":     "claude-3-5-sonnet or claude-3-5-haiku",
		"claude-3-5-sonnet": "claude-3-5-haiku for simple tasks",
	}

	// For simple tasks with expensive models, suggest alternatives
	if taskComplexity == "simple" && expensiveModels[model] {
		alternative := cheapAlternatives[model]
		if alternative != "" {
			return &SmartCostSuggestion{
				Title:          "Cost Optimization Available",
				CurrentCost:    0.03, // Example cost
				PotentialCost:  0.005, // Example cheaper cost
				Savings:        0.025,
				Recommendation: fmt.Sprintf("For simple tasks, try %s instead of %s", alternative, model),
			}
		}
	}

	return nil
}

package intelligence

import (
	"fmt"
	"time"

	"github.com/aaronmrosenthal/rycode/internal/styles"
	"github.com/aaronmrosenthal/rycode/internal/theme"
	"github.com/aaronmrosenthal/rycode/internal/typography"
	"github.com/charmbracelet/lipgloss/v2"
)

// ModelRecommendation represents an AI-powered model suggestion
type ModelRecommendation struct {
	Provider   string
	Model      string
	Score      float64  // 0-100 confidence score
	Reasoning  string   // Why this model is recommended
	Pros       []string // Advantages
	Cons       []string // Disadvantages
	CostPerUse float64  // Estimated cost for this task
	Speed      string   // "fast", "medium", "slow"
	Quality    string   // "high", "medium", "basic"
}

// TaskContext provides context for recommendations
type TaskContext struct {
	Description string        // What the user wants to do
	Complexity  string        // "simple", "medium", "complex"
	Priority    string        // "cost", "quality", "speed"
	Budget      float64       // Available budget
	TimeOfDay   time.Time     // When the task is being done
	RecentUsage []ModelUsage  // Recent model usage history
}

// ModelUsage tracks model usage patterns
type ModelUsage struct {
	Provider  string
	Model     string
	UsedAt    time.Time
	TaskType  string
	Satisfied bool // Was user satisfied with result
}

// RecommendationEngine generates intelligent model recommendations
type RecommendationEngine struct {
	usageHistory []ModelUsage
}

// NewRecommendationEngine creates a new recommendation engine
func NewRecommendationEngine() *RecommendationEngine {
	return &RecommendationEngine{
		usageHistory: make([]ModelUsage, 0),
	}
}

// GetRecommendations returns top 3 model recommendations for a task
func (r *RecommendationEngine) GetRecommendations(ctx TaskContext) []ModelRecommendation {
	recommendations := []ModelRecommendation{}

	// Analyze context and generate recommendations
	if ctx.Priority == "cost" {
		recommendations = append(recommendations, r.getCostOptimizedRecommendations(ctx)...)
	} else if ctx.Priority == "quality" {
		recommendations = append(recommendations, r.getQualityOptimizedRecommendations(ctx)...)
	} else if ctx.Priority == "speed" {
		recommendations = append(recommendations, r.getSpeedOptimizedRecommendations(ctx)...)
	} else {
		// Balanced recommendations
		recommendations = append(recommendations, r.getBalancedRecommendations(ctx)...)
	}

	// Apply learning from usage history
	recommendations = r.adjustForUserPreferences(recommendations, ctx)

	// Sort by score and return top 3
	return r.getTopRecommendations(recommendations, 3)
}

// getCostOptimizedRecommendations suggests cheapest viable options
func (r *RecommendationEngine) getCostOptimizedRecommendations(ctx TaskContext) []ModelRecommendation {
	recommendations := []ModelRecommendation{}

	// Claude 3.5 Haiku - Best value
	recommendations = append(recommendations, ModelRecommendation{
		Provider: "Anthropic",
		Model:    "claude-3-5-haiku-20241022",
		Score:    95,
		Reasoning: "Best cost-to-quality ratio for most tasks",
		Pros: []string{
			"Extremely cost-effective ($0.001/1K input tokens)",
			"Fast response times",
			"200K context window",
			"Good quality for simple-medium tasks",
		},
		Cons: []string{
			"May struggle with highly complex tasks",
			"Not as nuanced as premium models",
		},
		CostPerUse: 0.005,
		Speed:      "fast",
		Quality:    "medium",
	})

	// GPT-3.5 Turbo - Ultra cheap
	recommendations = append(recommendations, ModelRecommendation{
		Provider: "OpenAI",
		Model:    "gpt-3.5-turbo",
		Score:    85,
		Reasoning: "Cheapest option for simple tasks",
		Pros: []string{
			"Very low cost ($0.0005/1K input tokens)",
			"Fast responses",
			"Reliable for simple queries",
		},
		Cons: []string{
			"Limited for complex reasoning",
			"Smaller context window (16K)",
		},
		CostPerUse: 0.002,
		Speed:      "fast",
		Quality:    "basic",
	})

	// Gemini Flash - Fast and cheap
	if ctx.Complexity == "simple" {
		recommendations = append(recommendations, ModelRecommendation{
			Provider: "Google",
			Model:    "gemini-1.5-flash",
			Score:    80,
			Reasoning: "Ultra-fast and ultra-cheap for simple tasks",
			Pros: []string{
				"Extremely low cost ($0.000075/1K tokens)",
				"Lightning fast",
				"1M context window",
			},
			Cons: []string{
				"Basic quality only",
				"Not suitable for complex tasks",
			},
			CostPerUse: 0.001,
			Speed:      "fast",
			Quality:    "basic",
		})
	}

	return recommendations
}

// getQualityOptimizedRecommendations suggests highest quality options
func (r *RecommendationEngine) getQualityOptimizedRecommendations(ctx TaskContext) []ModelRecommendation {
	recommendations := []ModelRecommendation{}

	// Claude 3.5 Sonnet - Best overall quality
	recommendations = append(recommendations, ModelRecommendation{
		Provider: "Anthropic",
		Model:    "claude-3-5-sonnet-20241022",
		Score:    98,
		Reasoning: "Industry-leading quality and reasoning",
		Pros: []string{
			"Best-in-class reasoning and coding",
			"200K context window",
			"Excellent at complex tasks",
			"Balanced cost ($0.003/1K input)",
		},
		Cons: []string{
			"Slightly more expensive than Haiku",
		},
		CostPerUse: 0.015,
		Speed:      "medium",
		Quality:    "high",
	})

	// GPT-4 Turbo - Strong alternative
	recommendations = append(recommendations, ModelRecommendation{
		Provider: "OpenAI",
		Model:    "gpt-4-turbo-preview",
		Score:    90,
		Reasoning: "Excellent for creative and analytical tasks",
		Pros: []string{
			"Strong creative capabilities",
			"128K context window",
			"Good for diverse tasks",
		},
		Cons: []string{
			"Higher cost ($0.01/1K input)",
			"Can be slower than Claude",
		},
		CostPerUse: 0.05,
		Speed:      "medium",
		Quality:    "high",
	})

	// Claude 3 Opus - Premium quality
	if ctx.Budget > 0.10 {
		recommendations = append(recommendations, ModelRecommendation{
			Provider: "Anthropic",
			Model:    "claude-3-opus-20240229",
			Score:    95,
			Reasoning: "Premium quality for most demanding tasks",
			Pros: []string{
				"Highest quality reasoning",
				"Best for very complex tasks",
				"200K context window",
			},
			Cons: []string{
				"Most expensive ($0.015/1K input)",
				"Slower responses",
			},
			CostPerUse: 0.075,
			Speed:      "slow",
			Quality:    "high",
		})
	}

	return recommendations
}

// getSpeedOptimizedRecommendations suggests fastest options
func (r *RecommendationEngine) getSpeedOptimizedRecommendations(ctx TaskContext) []ModelRecommendation {
	recommendations := []ModelRecommendation{}

	// Gemini Flash - Fastest
	recommendations = append(recommendations, ModelRecommendation{
		Provider: "Google",
		Model:    "gemini-1.5-flash",
		Score:    95,
		Reasoning: "Lightning-fast responses",
		Pros: []string{
			"Extremely fast",
			"Very low cost",
			"1M context window",
		},
		Cons: []string{
			"Basic quality",
		},
		CostPerUse: 0.001,
		Speed:      "fast",
		Quality:    "basic",
	})

	// Claude Haiku - Fast and good quality
	recommendations = append(recommendations, ModelRecommendation{
		Provider: "Anthropic",
		Model:    "claude-3-5-haiku-20241022",
		Score:    90,
		Reasoning: "Fast with better quality than ultra-cheap models",
		Pros: []string{
			"Very fast",
			"Good quality",
			"Low cost",
		},
		Cons: []string{
			"Not as fast as Gemini Flash",
		},
		CostPerUse: 0.005,
		Speed:      "fast",
		Quality:    "medium",
	})

	// Grok 2 Mini - Fast alternative
	recommendations = append(recommendations, ModelRecommendation{
		Provider: "Grok",
		Model:    "grok-2-mini",
		Score:    75,
		Reasoning: "Fast and efficient",
		Pros: []string{
			"Quick responses",
			"Low cost",
		},
		Cons: []string{
			"Newer, less proven",
		},
		CostPerUse: 0.003,
		Speed:      "fast",
		Quality:    "medium",
	})

	return recommendations
}

// getBalancedRecommendations suggests best overall options
func (r *RecommendationEngine) getBalancedRecommendations(ctx TaskContext) []ModelRecommendation {
	recommendations := []ModelRecommendation{}

	// Default to Claude 3.5 Sonnet - best balance
	recommendations = append(recommendations, ModelRecommendation{
		Provider: "Anthropic",
		Model:    "claude-3-5-sonnet-20241022",
		Score:    95,
		Reasoning: "Perfect balance of cost, speed, and quality",
		Pros: []string{
			"Excellent quality",
			"Reasonable cost",
			"Good speed",
			"200K context",
		},
		Cons: []string{
			"None significant",
		},
		CostPerUse: 0.015,
		Speed:      "medium",
		Quality:    "high",
	})

	// For simpler tasks, suggest Haiku
	if ctx.Complexity == "simple" || ctx.Complexity == "medium" {
		recommendations = append(recommendations, ModelRecommendation{
			Provider: "Anthropic",
			Model:    "claude-3-5-haiku-20241022",
			Score:    90,
			Reasoning: "Great value for simple-medium tasks",
			Pros: []string{
				"Very cost-effective",
				"Fast",
				"Good quality",
			},
			Cons: []string{
				"May struggle with complex tasks",
			},
			CostPerUse: 0.005,
			Speed:      "fast",
			Quality:    "medium",
		})
	}

	return recommendations
}

// adjustForUserPreferences learns from usage history
func (r *RecommendationEngine) adjustForUserPreferences(recs []ModelRecommendation, ctx TaskContext) []ModelRecommendation {
	// Boost score for models user has been satisfied with
	for i, rec := range recs {
		for _, usage := range ctx.RecentUsage {
			if usage.Provider == rec.Provider && usage.Model == rec.Model && usage.Satisfied {
				recs[i].Score += 5
				recs[i].Reasoning += " (You've been satisfied with this model before)"
			}
		}
	}

	// Time-of-day preferences
	hour := ctx.TimeOfDay.Hour()
	if hour >= 9 && hour <= 17 {
		// Work hours - boost quality
		for i, rec := range recs {
			if rec.Quality == "high" {
				recs[i].Score += 3
			}
		}
	} else {
		// After hours - boost cost savings
		for i, rec := range recs {
			if rec.CostPerUse < 0.01 {
				recs[i].Score += 3
			}
		}
	}

	return recs
}

// getTopRecommendations returns top N recommendations by score
func (r *RecommendationEngine) getTopRecommendations(recs []ModelRecommendation, n int) []ModelRecommendation {
	// Sort by score descending
	for i := 0; i < len(recs)-1; i++ {
		for j := i + 1; j < len(recs); j++ {
			if recs[j].Score > recs[i].Score {
				recs[i], recs[j] = recs[j], recs[i]
			}
		}
	}

	if len(recs) <= n {
		return recs
	}
	return recs[:n]
}

// RenderRecommendation renders a recommendation card
func RenderRecommendation(rec ModelRecommendation, highlighted bool) string {
	t := theme.CurrentTheme()
	typo := typography.New()

	// Header with model name and score
	header := ""
	if highlighted {
		header = styles.NewStyle().
			Background(t.Primary()).
			Foreground(t.Background()).
			Bold(true).
			Padding(0, 1).
			Render(fmt.Sprintf("⭐ %s - %s", rec.Provider, rec.Model))
	} else {
		header = typo.Subheading.Render(fmt.Sprintf("%s - %s", rec.Provider, rec.Model))
	}

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

	// Metadata
	costBadge := typography.Badge(fmt.Sprintf("$%.3f/use", rec.CostPerUse))
	speedBadge := typography.StatusBadge(rec.Speed, "info")
	qualityBadge := typography.StatusBadge(rec.Quality, "success")

	metadata := typography.Inline([]string{costBadge, speedBadge, qualityBadge}, "  ")

	// Reasoning
	reasoning := typo.Body.
		Foreground(t.Text()).
		MarginTop(1).
		Render(rec.Reasoning)

	// Pros/Cons
	prosTitle := typo.Label.Render("PROS:")
	prosList := ""
	for _, pro := range rec.Pros {
		prosList += typo.Body.Foreground(t.Success()).Render("  ✓ " + pro) + "\n"
	}

	consTitle := typo.Label.Render("CONS:")
	consList := ""
	for _, con := range rec.Cons {
		consList += typo.Body.Foreground(t.TextMuted()).Render("  • " + con) + "\n"
	}

	// Combine all parts
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header+" "+scoreBadge,
		"",
		metadata,
		"",
		reasoning,
		"",
		prosTitle,
		prosList,
		consTitle,
		consList,
	)

	// Wrap in card
	borderColor := t.Border()
	if highlighted {
		borderColor = t.Primary()
	}

	return styles.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1, 2).
		Width(60).
		Render(content)
}

// RecommendationPanelMsg requests model recommendations
type RecommendationPanelMsg struct {
	Context TaskContext
}

// RecommendationsReadyMsg contains generated recommendations
type RecommendationsReadyMsg struct {
	Recommendations []ModelRecommendation
}

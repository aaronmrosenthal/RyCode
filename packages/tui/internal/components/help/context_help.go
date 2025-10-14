package help

import (
	"strings"

	"github.com/aaronmrosenthal/rycode/internal/styles"
	"github.com/aaronmrosenthal/rycode/internal/theme"
)

// HelpContext represents different areas of the app where help can be shown
type HelpContext string

const (
	ContextModelSelector   HelpContext = "model_selector"
	ContextProviderList    HelpContext = "provider_list"
	ContextInsightsDash    HelpContext = "insights_dashboard"
	ContextBudgetForecast  HelpContext = "budget_forecast"
	ContextChat            HelpContext = "chat"
	ContextEmpty           HelpContext = "empty_state"
	ContextFirstRun        HelpContext = "first_run"
	ContextAuthentication  HelpContext = "authentication"
	ContextModelRecommend  HelpContext = "model_recommendations"
)

// ContextualHint represents a helpful tip for a specific context
type ContextualHint struct {
	Icon        string
	Title       string
	Message     string
	Action      string // Suggested action
	Shortcut    string // Keyboard shortcut if applicable
	Dismissable bool
}

// ContextHelpProvider manages contextual help throughout the app
type ContextHelpProvider struct {
	dismissedHints map[HelpContext]bool
	hintDatabase   map[HelpContext][]ContextualHint
}

// NewContextHelpProvider creates a new contextual help provider
func NewContextHelpProvider() *ContextHelpProvider {
	provider := &ContextHelpProvider{
		dismissedHints: make(map[HelpContext]bool),
		hintDatabase:   make(map[HelpContext][]ContextualHint),
	}

	provider.setupHints()
	return provider
}

// setupHints populates the hint database
func (p *ContextHelpProvider) setupHints() {
	p.hintDatabase = map[HelpContext][]ContextualHint{
		ContextModelSelector: {
			{
				Icon:        "💡",
				Title:       "Quick Tip",
				Message:     "Press Tab to quickly cycle through your most-used models",
				Action:      "Try pressing Tab now",
				Shortcut:    "Tab",
				Dismissable: true,
			},
			{
				Icon:        "🎯",
				Title:       "Smart Recommendations",
				Message:     "Toggle 'i' to see AI-powered model recommendations based on your task",
				Action:      "Press 'i' to toggle recommendations",
				Shortcut:    "i",
				Dismissable: true,
			},
		},
		ContextProviderList: {
			{
				Icon:        "🔐",
				Title:       "No Providers Authenticated",
				Message:     "You need to authenticate at least one provider to use RyCode",
				Action:      "Press 'a' to authenticate, or 'd' to auto-detect credentials",
				Shortcut:    "a / d",
				Dismissable: false,
			},
			{
				Icon:        "⚡",
				Title:       "Quick Auth",
				Message:     "Auto-detect can find API keys from your environment variables",
				Action:      "Press 'd' to scan for credentials automatically",
				Shortcut:    "d",
				Dismissable: true,
			},
		},
		ContextInsightsDash: {
			{
				Icon:        "📊",
				Title:       "No Usage Data Yet",
				Message:     "Start using RyCode to see beautiful analytics and insights here",
				Action:      "Make your first API call to see data",
				Shortcut:    "",
				Dismissable: false,
			},
			{
				Icon:        "💰",
				Title:       "Cost Optimization",
				Message:     "Check the 'Optimization Opportunities' section for ways to save money",
				Action:      "Scroll down to see suggestions",
				Shortcut:    "",
				Dismissable: true,
			},
		},
		ContextBudgetForecast: {
			{
				Icon:        "🔮",
				Title:       "Predictive Budgeting",
				Message:     "RyCode learns from your spending patterns to forecast month-end costs",
				Action:      "Use RyCode for a few days to get accurate predictions",
				Shortcut:    "",
				Dismissable: true,
			},
			{
				Icon:        "⚠️",
				Title:       "Budget Alert",
				Message:     "You're projected to exceed your budget this month",
				Action:      "Check recommendations below for cost-saving tips",
				Shortcut:    "",
				Dismissable: false,
			},
		},
		ContextChat: {
			{
				Icon:        "⌨️",
				Title:       "Keyboard Shortcuts",
				Message:     "Press Ctrl+? to see all available keyboard shortcuts",
				Action:      "Master shortcuts to boost productivity",
				Shortcut:    "Ctrl+?",
				Dismissable: true,
			},
			{
				Icon:        "🔄",
				Title:       "Switch Models Anytime",
				Message:     "Press Tab to quickly switch between models mid-conversation",
				Action:      "Try switching to a cheaper model for simple questions",
				Shortcut:    "Tab",
				Dismissable: true,
			},
		},
		ContextEmpty: {
			{
				Icon:        "👋",
				Title:       "Welcome to RyCode!",
				Message:     "Start by selecting a model and asking a question",
				Action:      "Press Ctrl+M to open the model selector",
				Shortcut:    "Ctrl+M",
				Dismissable: false,
			},
		},
		ContextFirstRun: {
			{
				Icon:        "🚀",
				Title:       "First Time Here?",
				Message:     "Let's get you set up with a quick walkthrough",
				Action:      "Press Enter to start the welcome guide",
				Shortcut:    "Enter",
				Dismissable: false,
			},
		},
		ContextAuthentication: {
			{
				Icon:        "🔑",
				Title:       "API Key Required",
				Message:     "Enter your API key to authenticate with this provider",
				Action:      "Paste your key and press Enter",
				Shortcut:    "",
				Dismissable: false,
			},
			{
				Icon:        "🔒",
				Title:       "Your Keys Are Safe",
				Message:     "API keys are stored securely in your system keychain",
				Action:      "RyCode never sends your keys to third parties",
				Shortcut:    "",
				Dismissable: true,
			},
		},
		ContextModelRecommend: {
			{
				Icon:        "🤖",
				Title:       "AI Recommendations",
				Message:     "RyCode analyzes your task to suggest the best model for quality, cost, and speed",
				Action:      "Check the 'Reasoning' section to understand each recommendation",
				Shortcut:    "",
				Dismissable: true,
			},
		},
	}
}

// GetHint returns the primary hint for a given context
func (p *ContextHelpProvider) GetHint(ctx HelpContext) *ContextualHint {
	hints := p.hintDatabase[ctx]
	if len(hints) == 0 {
		return nil
	}

	// Return first non-dismissed hint
	for _, hint := range hints {
		if hint.Dismissable && p.dismissedHints[ctx] {
			continue
		}
		return &hint
	}

	return nil
}

// GetAllHints returns all hints for a context
func (p *ContextHelpProvider) GetAllHints(ctx HelpContext) []ContextualHint {
	hints := p.hintDatabase[ctx]

	// Filter out dismissed hints
	var active []ContextualHint
	for _, hint := range hints {
		if hint.Dismissable && p.dismissedHints[ctx] {
			continue
		}
		active = append(active, hint)
	}

	return active
}

// DismissHint marks a hint as dismissed for a context
func (p *ContextHelpProvider) DismissHint(ctx HelpContext) {
	p.dismissedHints[ctx] = true
}

// ResetHints clears all dismissed hints
func (p *ContextHelpProvider) ResetHints() {
	p.dismissedHints = make(map[HelpContext]bool)
}

// RenderHint creates a beautiful hint card
func (p *ContextHelpProvider) RenderHint(hint *ContextualHint, width int) string {
	if hint == nil {
		return ""
	}

	t := theme.CurrentTheme()

	var lines []string

	// Icon and title
	titleStyle := styles.NewStyle().
		Foreground(t.Info()).
		Bold(true)

	title := titleStyle.Render(hint.Icon + " " + hint.Title)
	lines = append(lines, title)

	// Message
	messageStyle := styles.NewStyle().
		Foreground(t.Text()).
		Width(width - 4)

	message := messageStyle.Render(hint.Message)
	lines = append(lines, message)

	// Action
	if hint.Action != "" {
		actionStyle := styles.NewStyle().
			Foreground(t.Primary()).
			Italic(true)

		action := actionStyle.Render("→ " + hint.Action)
		lines = append(lines, action)
	}

	// Shortcut badge
	if hint.Shortcut != "" {
		shortcutStyle := styles.NewStyle().
			Foreground(t.Background()).
			Background(t.Primary()).
			Bold(true).
			Padding(0, 1)

		shortcut := shortcutStyle.Render(hint.Shortcut)
		lines = append(lines, shortcut)
	}

	content := strings.Join(lines, "\n")

	// Wrap in subtle box
	boxStyle := styles.NewStyle().
		Border(styles.NewStyle().GetBorder()).
		BorderForeground(t.Info()).
		Padding(1, 2).
		Width(width)

	return boxStyle.Render(content)
}

// RenderCompactHint creates a single-line hint for status bar
func (p *ContextHelpProvider) RenderCompactHint(hint *ContextualHint) string {
	if hint == nil {
		return ""
	}

	t := theme.CurrentTheme()

	// Icon
	icon := styles.NewStyle().
		Foreground(t.Info()).
		Render(hint.Icon)

	// Message
	message := styles.NewStyle().
		Foreground(t.TextMuted()).
		Render(hint.Message)

	// Shortcut
	shortcut := ""
	if hint.Shortcut != "" {
		shortcut = styles.NewStyle().
			Foreground(t.Primary()).
			Bold(true).
			Render(" [" + hint.Shortcut + "]")
	}

	return icon + " " + message + shortcut
}

// GetStatusBarHint returns a context-appropriate hint for the status bar
func (p *ContextHelpProvider) GetStatusBarHint(ctx HelpContext, viewContext string) string {
	// Dynamic hints based on current view state
	hints := map[string]string{
		"model_selector_empty":       "Press Ctrl+P to authenticate a provider first",
		"model_selector_loaded":      "Tab: Quick switch | i: Toggle recommendations | Ctrl+?: All shortcuts",
		"chat_empty":                 "Select a model with Ctrl+M, then start chatting",
		"chat_active":                "Tab: Switch model | Ctrl+I: Usage insights | Ctrl+B: Budget forecast",
		"provider_list_none":         "Press 'a' to authenticate or 'd' to auto-detect credentials",
		"provider_list_partial":      "Press 'a' to add more providers",
		"insights_no_data":           "Make API calls to see analytics here",
		"insights_with_data":         "Ctrl+B: Budget forecast | Scroll for optimization tips",
		"budget_under":               "You're under budget - consider using premium models",
		"budget_over":                "Budget overrun likely - check recommendations",
		"welcome_flow":               "Press Enter to continue | ← → to navigate steps",
		"shortcuts_guide":            "ESC: Close | Scroll to see all categories",
	}

	t := theme.CurrentTheme()

	key := string(ctx) + "_" + viewContext
	if msg, ok := hints[key]; ok {
		return styles.NewStyle().
			Foreground(t.TextMuted()).
			Faint(true).
			Render("💡 " + msg)
	}

	// Fallback to generic hint
	hint := p.GetHint(ctx)
	if hint != nil {
		return p.RenderCompactHint(hint)
	}

	return ""
}

// GetEmptyStateMessage returns a helpful empty state message
func GetEmptyStateMessage(ctx HelpContext) string {
	messages := map[HelpContext]string{
		ContextModelSelector:  "No models available yet\n\nAuthenticate a provider with Ctrl+P to get started",
		ContextProviderList:   "No providers authenticated\n\nPress 'a' to add a provider or 'd' to auto-detect",
		ContextInsightsDash:   "No usage data yet\n\nStart using RyCode to see beautiful analytics here",
		ContextBudgetForecast: "Building your forecast...\n\nUse RyCode for a few days to get predictions",
		ContextChat:           "Ready to chat!\n\nSelect a model with Ctrl+M to begin",
	}

	t := theme.CurrentTheme()

	if msg, ok := messages[ctx]; ok {
		return styles.NewStyle().
			Foreground(t.TextMuted()).
			Faint(true).
			Italic(true).
			Render(msg)
	}

	return ""
}

// GetTooltip returns a tooltip for UI elements
func GetTooltip(element string) string {
	tooltips := map[string]string{
		"model_name":           "The AI model that will handle your request",
		"cost_estimate":        "Estimated cost per 1K tokens (input + output)",
		"provider_badge":       "Provider and authentication status",
		"budget_bar":           "Your spending vs monthly budget",
		"recommendation_score": "AI confidence score (0-100)",
		"trend_indicator":      "Spending trend compared to previous period",
		"health_status":        "Provider API health status",
		"model_speed":          "Response time: fast (<1s), medium (1-3s), slow (>3s)",
		"model_quality":        "Output quality tier: high, medium, basic",
		"confidence_badge":     "Forecast confidence based on data points",
		"api_key_mask":         "Last 4 characters of your API key (full key hidden)",
	}

	t := theme.CurrentTheme()

	if tip, ok := tooltips[element]; ok {
		return styles.NewStyle().
			Foreground(t.Background()).
			Background(t.Text()).
			Padding(0, 1).
			Render("ℹ️  " + tip)
	}

	return ""
}

// ProgressiveTips provides tips that appear based on user progression
type ProgressiveTips struct {
	tipsShown      map[string]bool
	interactionCount map[string]int
}

// NewProgressiveTips creates a new progressive tips tracker
func NewProgressiveTips() *ProgressiveTips {
	return &ProgressiveTips{
		tipsShown:      make(map[string]bool),
		interactionCount: make(map[string]int),
	}
}

// RecordInteraction tracks user interactions for tip timing
func (pt *ProgressiveTips) RecordInteraction(action string) {
	pt.interactionCount[action]++
}

// ShouldShowTip determines if a progressive tip should appear
func (pt *ProgressiveTips) ShouldShowTip(tipID string, minInteractions int) bool {
	// Don't show if already shown
	if pt.tipsShown[tipID] {
		return false
	}

	// Show after minimum interactions
	count := pt.interactionCount[tipID]
	return count >= minInteractions
}

// MarkTipShown marks a tip as displayed
func (pt *ProgressiveTips) MarkTipShown(tipID string) {
	pt.tipsShown[tipID] = true
}

// GetProgressiveTip returns a tip based on user progression
func (pt *ProgressiveTips) GetProgressiveTip() *ContextualHint {
	// After 5 model switches, suggest Tab shortcut
	if pt.ShouldShowTip("model_switch", 5) {
		pt.MarkTipShown("model_switch")
		return &ContextualHint{
			Icon:        "⚡",
			Title:       "Power User Tip",
			Message:     "You're switching models often - try Tab for instant cycling",
			Action:      "Press Tab to quickly switch between your favorite models",
			Shortcut:    "Tab",
			Dismissable: true,
		}
	}

	// After 10 requests, suggest insights
	if pt.ShouldShowTip("request", 10) {
		pt.MarkTipShown("request")
		return &ContextualHint{
			Icon:        "📊",
			Title:       "Check Your Insights",
			Message:     "You've made 10+ requests - see your usage patterns and savings opportunities",
			Action:      "Press Ctrl+I to view your analytics dashboard",
			Shortcut:    "Ctrl+I",
			Dismissable: true,
		}
	}

	// After spending $5, suggest budget forecast
	if pt.ShouldShowTip("spending", 1) {
		pt.MarkTipShown("spending")
		return &ContextualHint{
			Icon:        "🔮",
			Title:       "Budget Management",
			Message:     "You're spending regularly - let RyCode predict your month-end costs",
			Action:      "Press Ctrl+B to see your budget forecast",
			Shortcut:    "Ctrl+B",
			Dismissable: true,
		}
	}

	return nil
}

// FormatHelpFooter creates a contextual help footer for dialogs
func FormatHelpFooter(primaryAction, primaryKey, secondaryAction, secondaryKey string) string {
	t := theme.CurrentTheme()

	// Primary action (highlighted)
	primaryKeyStyle := styles.NewStyle().
		Foreground(t.Background()).
		Background(t.Primary()).
		Bold(true).
		Padding(0, 1)

	primaryTextStyle := styles.NewStyle().
		Foreground(t.Text())

	primary := primaryKeyStyle.Render(primaryKey) + " " + primaryTextStyle.Render(primaryAction)

	// Secondary action (muted)
	secondaryKeyStyle := styles.NewStyle().
		Foreground(t.Background()).
		Background(t.TextMuted()).
		Padding(0, 1)

	secondaryTextStyle := styles.NewStyle().
		Foreground(t.TextMuted())

	secondary := secondaryKeyStyle.Render(secondaryKey) + " " + secondaryTextStyle.Render(secondaryAction)

	return primary + "  " + secondary
}

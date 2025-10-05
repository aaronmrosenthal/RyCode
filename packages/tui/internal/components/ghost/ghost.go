package ghost

import (
	"strings"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/sst/opencode/internal/theme"
)

// Suggestion represents a ghost text prediction
type Suggestion struct {
	Text       string
	Confidence float64 // 0-1
	Trigger    string  // "tab", "enter", "auto"
	Source     string  // "pattern", "history", "ml"
}

// Predictor interface for different prediction strategies
type Predictor interface {
	Predict(input string, context map[string]interface{}) (*Suggestion, error)
	Learn(accepted bool, suggestion *Suggestion)
}

// PatternPredictor learns from user's accepted suggestions
type PatternPredictor struct {
	patterns map[string]*patternStats
}

type patternStats struct {
	count    int
	accepted int
}

// NewPatternPredictor creates a new pattern-based predictor
func NewPatternPredictor() *PatternPredictor {
	return &PatternPredictor{
		patterns: make(map[string]*patternStats),
	}
}

// Predict generates a suggestion based on patterns
func (p *PatternPredictor) Predict(input string, context map[string]interface{}) (*Suggestion, error) {
	// Common completions for partial commands
	completions := map[string]string{
		"how do i test":        "how do i test this component?",
		"fix the":              "fix the bug in",
		"add a":                "add a feature to",
		"explain":              "explain how this works",
		"refactor":             "refactor this code",
		"debug":                "debug the issue",
		"why is":               "why is this happening?",
		"what does":            "what does this do?",
		"show me":              "show me the implementation",
		"implement":            "implement the feature",
		"/t":                   "/test",
		"/d":                   "/debug",
		"/f":                   "/fix",
		"/e":                   "/explain",
		"/p":                   "/preview",
		"/c":                   "/commit",
		"/r":                   "/review",
	}

	input = strings.ToLower(strings.TrimSpace(input))

	for prefix, completion := range completions {
		if strings.HasPrefix(input, prefix) && input != completion {
			confidence := 0.8
			if len(input) > len(prefix)*2/3 {
				confidence = 0.9
			}

			return &Suggestion{
				Text:       completion,
				Confidence: confidence,
				Trigger:    "tab",
				Source:     "pattern",
			}, nil
		}
	}

	return nil, nil
}

// Learn updates pattern statistics
func (p *PatternPredictor) Learn(accepted bool, suggestion *Suggestion) {
	key := suggestion.Text
	stats, exists := p.patterns[key]
	if !exists {
		stats = &patternStats{count: 0, accepted: 0}
		p.patterns[key] = stats
	}

	stats.count++
	if accepted {
		stats.accepted++
	}
}

// Render ghost text with styling
func Render(suggestion *Suggestion, theme *theme.Theme) string {
	if suggestion == nil {
		return ""
	}

	// Different styling based on confidence
	var style lipgloss.Style
	if suggestion.Confidence > 0.8 {
		style = lipgloss.NewStyle().
			Foreground(theme.GhostTextHigh).
			Faint(true)
	} else {
		style = lipgloss.NewStyle().
			Foreground(theme.GhostTextLow).
			Faint(true)
	}

	// Add trigger hint
	triggerHint := ""
	switch suggestion.Trigger {
	case "tab":
		triggerHint = " ⇥"
	case "enter":
		triggerHint = " ↵"
	}

	return style.Render(suggestion.Text + triggerHint)
}

// RenderInline renders ghost text inline with current input
func RenderInline(currentInput string, suggestion *Suggestion, theme *theme.Theme) string {
	if suggestion == nil || suggestion.Text == "" {
		return currentInput
	}

	// Only show the completion part (what comes after current input)
	if !strings.HasPrefix(strings.ToLower(suggestion.Text), strings.ToLower(currentInput)) {
		return currentInput
	}

	completion := suggestion.Text[len(currentInput):]

	ghostStyle := lipgloss.NewStyle().
		Foreground(theme.GhostTextHigh).
		Faint(true)

	return currentInput + ghostStyle.Render(completion)
}

// ContextualSuggestions provides context-aware command suggestions
type ContextualSuggestions struct {
	CurrentFile       string
	HasErrors         bool
	HasUncommitted    bool
	RecentCommands    []string
	IsInTestFile      bool
	IsInReactFile     bool
}

// GetSuggestions returns contextual command suggestions
func (c *ContextualSuggestions) GetSuggestions() []Suggestion {
	suggestions := []Suggestion{}

	// Test file context
	if c.IsInTestFile {
		suggestions = append(suggestions, Suggestion{
			Text:       "/test",
			Confidence: 0.95,
			Trigger:    "tab",
			Source:     "context",
		})
		suggestions = append(suggestions, Suggestion{
			Text:       "/coverage",
			Confidence: 0.8,
			Trigger:    "tab",
			Source:     "context",
		})
	}

	// React file context
	if c.IsInReactFile {
		suggestions = append(suggestions, Suggestion{
			Text:       "/preview",
			Confidence: 0.9,
			Trigger:    "tab",
			Source:     "context",
		})
		suggestions = append(suggestions, Suggestion{
			Text:       "/test components",
			Confidence: 0.8,
			Trigger:    "tab",
			Source:     "context",
		})
	}

	// Error context
	if c.HasErrors {
		suggestions = append(suggestions, Suggestion{
			Text:       "/debug",
			Confidence: 0.95,
			Trigger:    "tab",
			Source:     "context",
		})
		suggestions = append(suggestions, Suggestion{
			Text:       "/fix",
			Confidence: 0.9,
			Trigger:    "tab",
			Source:     "context",
		})
	}

	// Uncommitted changes
	if c.HasUncommitted {
		suggestions = append(suggestions, Suggestion{
			Text:       "/commit",
			Confidence: 0.85,
			Trigger:    "tab",
			Source:     "context",
		})
		suggestions = append(suggestions, Suggestion{
			Text:       "/review",
			Confidence: 0.75,
			Trigger:    "tab",
			Source:     "context",
		})
	}

	return suggestions
}

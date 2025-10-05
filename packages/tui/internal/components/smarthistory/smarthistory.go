package smarthistory

import (
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/aaronmrosenthal/rycode/internal/theme"
)

// HistoryItem represents a command in history with context
type HistoryItem struct {
	Command   string
	Timestamp time.Time
	Context   Context
	Success   bool
	Duration  time.Duration
}

// Context captures the state when command was run
type Context struct {
	CurrentFile   string
	FileType      string // "test", "component", "config", etc.
	HasErrors     bool
	Branch        string
	LastError     string
	RecentChanges []string
}

// SmartHistory manages contextual command history
type SmartHistory struct {
	items        []HistoryItem
	maxItems     int
	contextIndex map[string][]int // context type -> indices
}

// NewSmartHistory creates a new smart history manager
func NewSmartHistory() *SmartHistory {
	return &SmartHistory{
		items:        []HistoryItem{},
		maxItems:     100,
		contextIndex: make(map[string][]int),
	}
}

// Add adds a command to history
func (sh *SmartHistory) Add(item HistoryItem) {
	sh.items = append([]HistoryItem{item}, sh.items...)

	if len(sh.items) > sh.maxItems {
		sh.items = sh.items[:sh.maxItems]
	}

	// Index by context
	contextKey := sh.getContextKey(item.Context)
	sh.contextIndex[contextKey] = append(sh.contextIndex[contextKey], 0)

	// Update all indices
	for key := range sh.contextIndex {
		for i := range sh.contextIndex[key] {
			sh.contextIndex[key][i]++
		}
	}
}

// GetContextual returns commands relevant to current context
func (sh *SmartHistory) GetContextual(currentContext Context) []HistoryItem {
	contextKey := sh.getContextKey(currentContext)
	indices, exists := sh.contextIndex[contextKey]

	if !exists || len(indices) == 0 {
		return sh.GetRecent(5)
	}

	// Get commands from this context
	contextualItems := []HistoryItem{}
	for _, idx := range indices {
		if idx < len(sh.items) {
			contextualItems = append(contextualItems, sh.items[idx])
		}
		if len(contextualItems) >= 5 {
			break
		}
	}

	return contextualItems
}

// GetRecent returns the most recent commands
func (sh *SmartHistory) GetRecent(n int) []HistoryItem {
	if n > len(sh.items) {
		n = len(sh.items)
	}
	return sh.items[:n]
}

// GetSuccessful returns commands that succeeded
func (sh *SmartHistory) GetSuccessful(n int) []HistoryItem {
	successful := []HistoryItem{}
	for _, item := range sh.items {
		if item.Success {
			successful = append(successful, item)
			if len(successful) >= n {
				break
			}
		}
	}
	return successful
}

// Search searches history by text
func (sh *SmartHistory) Search(query string) []HistoryItem {
	query = strings.ToLower(query)
	results := []HistoryItem{}

	for _, item := range sh.items {
		if strings.Contains(strings.ToLower(item.Command), query) {
			results = append(results, item)
		}
	}

	return results
}

// GetFrequent returns most frequently used commands in context
func (sh *SmartHistory) GetFrequent(context Context, n int) []HistoryItem {
	frequency := make(map[string]int)
	commandToItem := make(map[string]HistoryItem)

	contextKey := sh.getContextKey(context)
	indices := sh.contextIndex[contextKey]

	for _, idx := range indices {
		if idx < len(sh.items) {
			item := sh.items[idx]
			frequency[item.Command]++
			if _, exists := commandToItem[item.Command]; !exists {
				commandToItem[item.Command] = item
			}
		}
	}

	// Sort by frequency
	type freqItem struct {
		command string
		count   int
	}

	freqItems := []freqItem{}
	for cmd, count := range frequency {
		freqItems = append(freqItems, freqItem{cmd, count})
	}

	// Simple sort (bubble sort for small n)
	for i := 0; i < len(freqItems); i++ {
		for j := i + 1; j < len(freqItems); j++ {
			if freqItems[j].count > freqItems[i].count {
				freqItems[i], freqItems[j] = freqItems[j], freqItems[i]
			}
		}
	}

	// Get top n
	result := []HistoryItem{}
	for i := 0; i < n && i < len(freqItems); i++ {
		result = append(result, commandToItem[freqItems[i].command])
	}

	return result
}

// Render renders contextual history
func (sh *SmartHistory) Render(context Context, theme *theme.Theme) string {
	contextual := sh.GetContextual(context)

	if len(contextual) == 0 {
		return ""
	}

	// Title
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.AccentPrimary).
		Bold(true).
		MarginBottom(1)

	contextDesc := sh.getContextDescription(context)
	title := titleStyle.Render("Recent commands" + contextDesc)

	// Items
	items := []string{}
	for i, item := range contextual {
		if i >= 5 {
			break
		}

		// Number
		numberStyle := lipgloss.NewStyle().
			Foreground(theme.AccentSecondary).
			PaddingRight(1)

		// Command
		commandStyle := lipgloss.NewStyle().
			Foreground(theme.TextPrimary)

		// Time ago
		timeStyle := lipgloss.NewStyle().
			Foreground(theme.TextDim).
			PaddingLeft(1)

		// Success indicator
		successStyle := lipgloss.NewStyle().
			Foreground(theme.Success)
		failStyle := lipgloss.NewStyle().
			Foreground(theme.Error)

		indicator := ""
		if item.Success {
			indicator = successStyle.Render("✓")
		} else {
			indicator = failStyle.Render("✗")
		}

		timeAgo := formatTimeAgo(item.Timestamp)

		itemStr := lipgloss.JoinHorizontal(
			lipgloss.Left,
			numberStyle.Render(string(rune('1'+i))),
			indicator,
			commandStyle.Render(" "+item.Command),
			timeStyle.Render(" • "+timeAgo),
		)

		items = append(items, itemStr)
	}

	// Instruction
	instructionStyle := lipgloss.NewStyle().
		Foreground(theme.TextDim).
		MarginTop(1)

	instruction := instructionStyle.Render("Press 1-5 to reuse, / to search")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		lipgloss.JoinVertical(lipgloss.Left, items...),
		instruction,
	)
}

// Helper functions
func (sh *SmartHistory) getContextKey(context Context) string {
	parts := []string{}

	if context.FileType != "" {
		parts = append(parts, "type:"+context.FileType)
	}

	if context.HasErrors {
		parts = append(parts, "errors:true")
	}

	if len(parts) == 0 {
		return "general"
	}

	return strings.Join(parts, "|")
}

func (sh *SmartHistory) getContextDescription(context Context) string {
	if context.FileType == "test" {
		return " (in test files)"
	}
	if context.FileType == "component" {
		return " (in components)"
	}
	if context.HasErrors {
		return " (with errors)"
	}
	return ""
}

func formatTimeAgo(t time.Time) string {
	duration := time.Since(t)

	if duration < time.Minute {
		return "just now"
	}
	if duration < time.Hour {
		mins := int(duration.Minutes())
		return formatPlural(mins, "min")
	}
	if duration < 24*time.Hour {
		hours := int(duration.Hours())
		return formatPlural(hours, "hour")
	}
	days := int(duration.Hours() / 24)
	return formatPlural(days, "day")
}

func formatPlural(n int, unit string) string {
	if n == 1 {
		return "1 " + unit + " ago"
	}
	return string(rune('0'+n)) + " " + unit + "s ago"
}

// GetPatterns analyzes patterns in command usage
func (sh *SmartHistory) GetPatterns() map[string][]string {
	patterns := make(map[string][]string)

	// Detect common sequences
	for i := 0; i < len(sh.items)-1; i++ {
		curr := sh.items[i].Command
		next := sh.items[i+1].Command

		key := curr
		if _, exists := patterns[key]; !exists {
			patterns[key] = []string{}
		}

		// Add if not already there
		found := false
		for _, cmd := range patterns[key] {
			if cmd == next {
				found = true
				break
			}
		}
		if !found {
			patterns[key] = append(patterns[key], next)
		}
	}

	return patterns
}

// SuggestNext suggests next command based on current
func (sh *SmartHistory) SuggestNext(currentCommand string) []string {
	patterns := sh.GetPatterns()

	if suggestions, exists := patterns[currentCommand]; exists {
		return suggestions
	}

	return []string{}
}

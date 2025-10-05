package timeline

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/sst/opencode/internal/theme"
)

// Event represents a timeline event
type Event struct {
	Type        EventType
	Timestamp   time.Time
	Label       string
	Significance float64 // 0-1, importance of event
	Data        map[string]interface{}
}

// EventType categorizes timeline events
type EventType string

const (
	EventMessage  EventType = "message"
	EventBranch   EventType = "branch"
	EventSnapshot EventType = "snapshot"
	EventError    EventType = "error"
	EventSuccess  EventType = "success"
	EventEdit     EventType = "edit"
	EventTest     EventType = "test"
	EventCommit   EventType = "commit"
)

// Timeline manages conversation timeline
type Timeline struct {
	Events          []Event
	CurrentPosition int
	Width           int
}

// NewTimeline creates a new timeline
func NewTimeline(width int) *Timeline {
	return &Timeline{
		Events:          []Event{},
		CurrentPosition: 0,
		Width:           width,
	}
}

// AddEvent adds an event to timeline
func (t *Timeline) AddEvent(event Event) {
	t.Events = append(t.Events, event)
	t.CurrentPosition = len(t.Events) - 1
}

// ScrubTo moves to a specific position
func (t *Timeline) ScrubTo(position int) {
	if position >= 0 && position < len(t.Events) {
		t.CurrentPosition = position
	}
}

// GetCurrentEvent returns current event
func (t *Timeline) GetCurrentEvent() *Event {
	if t.CurrentPosition >= 0 && t.CurrentPosition < len(t.Events) {
		return &t.Events[t.CurrentPosition]
	}
	return nil
}

// Render renders the timeline bar
func (t *Timeline) Render(theme *theme.Theme) string {
	if len(t.Events) == 0 {
		return ""
	}

	// Calculate step size
	width := t.Width
	if width > 80 {
		width = 80
	}

	// Build timeline bar
	bar := ""
	step := float64(len(t.Events)) / float64(width)

	for i := 0; i < width; i++ {
		eventIndex := int(float64(i) * step)
		if eventIndex >= len(t.Events) {
			eventIndex = len(t.Events) - 1
		}

		event := t.Events[eventIndex]

		// Check if this is current position
		currentPos := int(float64(t.CurrentPosition) / step)
		if i == currentPos {
			bar += theme.TimelineCurrent + "●" + theme.Reset
		} else {
			// Use different symbols for different event types
			symbol, color := t.getEventSymbol(event, theme)
			bar += color + symbol + theme.Reset
		}
	}

	// Build top border
	borderStyle := lipgloss.NewStyle().
		Foreground(theme.Border)

	topBorder := borderStyle.Render("┌" + lipgloss.NewStyle().Render(bar) + "┐")

	// Build labels for major events
	labels := t.renderLabels(theme)

	// Build bottom border with legend
	bottomBorder := borderStyle.Render("└" + "─" + "┘")

	// Legend
	legendStyle := lipgloss.NewStyle().
		Foreground(theme.TextDim).
		MarginTop(1)

	legend := legendStyle.Render(
		theme.TimelineCurrent + "●" + theme.Reset + " Current  " +
			theme.Success + "✓" + theme.Reset + " Success  " +
			theme.Error + "✖" + theme.Reset + " Error  " +
			theme.AccentSecondary + "⎇" + theme.Reset + " Branch  " +
			theme.AccentPrimary + "◆" + theme.Reset + " Snapshot",
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		topBorder,
		labels,
		bottomBorder,
		legend,
	)
}

// RenderCompact renders a compact timeline
func (t *Timeline) RenderCompact(theme *theme.Theme) string {
	if len(t.Events) == 0 {
		return ""
	}

	// Show last 5 events
	recentEvents := t.Events
	if len(recentEvents) > 5 {
		recentEvents = recentEvents[len(recentEvents)-5:]
	}

	items := []string{}
	for i, event := range recentEvents {
		symbol, color := t.getEventSymbol(event, theme)

		// Check if current
		isCurrent := (len(t.Events) - len(recentEvents) + i) == t.CurrentPosition

		style := lipgloss.NewStyle().Foreground(color)
		if isCurrent {
			style = style.Bold(true)
		}

		timeStr := event.Timestamp.Format("15:04")
		label := event.Label
		if label == "" {
			label = string(event.Type)
		}

		item := lipgloss.JoinHorizontal(
			lipgloss.Left,
			style.Render(symbol),
			lipgloss.NewStyle().
				Foreground(theme.TextDim).
				PaddingLeft(1).
				Render(timeStr),
			lipgloss.NewStyle().
				Foreground(theme.TextPrimary).
				PaddingLeft(1).
				Render(label),
		)

		items = append(items, item)
	}

	return lipgloss.JoinVertical(lipgloss.Left, items...)
}

// RenderProgress renders a progress indicator
func (t *Timeline) RenderProgress(theme *theme.Theme) string {
	if len(t.Events) == 0 {
		return ""
	}

	progress := float64(t.CurrentPosition) / float64(len(t.Events)-1)
	if progress > 1.0 {
		progress = 1.0
	}

	width := 40
	filled := int(progress * float64(width))
	empty := width - filled

	barStyle := lipgloss.NewStyle().
		Foreground(theme.AccentPrimary)

	emptyStyle := lipgloss.NewStyle().
		Foreground(theme.Border)

	bar := "[" +
		barStyle.Render(lipgloss.NewStyle().Render(lipgloss.PlaceHorizontal(filled, lipgloss.Left, "█", lipgloss.WithWhitespaceChars("█")))) +
		emptyStyle.Render(lipgloss.NewStyle().Render(lipgloss.PlaceHorizontal(empty, lipgloss.Left, "░", lipgloss.WithWhitespaceChars("░")))) +
		"]"

	percentage := fmt.Sprintf(" %.0f%%", progress*100)

	return bar + percentage
}

// Helper functions
func (t *Timeline) getEventSymbol(event Event, theme *theme.Theme) (string, lipgloss.TerminalColor) {
	switch event.Type {
	case EventError:
		return "✖", theme.Error
	case EventSuccess:
		return "✓", theme.Success
	case EventBranch:
		return "⎇", theme.AccentSecondary
	case EventSnapshot:
		return "◆", theme.AccentPrimary
	case EventTest:
		return "🧪", theme.Info
	case EventCommit:
		return "📝", theme.Success
	case EventEdit:
		return "✎", theme.Warning
	default:
		return "═", theme.Border
	}
}

func (t *Timeline) renderLabels(theme *theme.Theme) string {
	if len(t.Events) == 0 {
		return ""
	}

	// Find major events (high significance)
	majorEvents := []Event{}
	for _, event := range t.Events {
		if event.Significance >= 0.7 {
			majorEvents = append(majorEvents, event)
		}
	}

	if len(majorEvents) == 0 {
		return ""
	}

	// Limit to 5 labels to avoid crowding
	if len(majorEvents) > 5 {
		majorEvents = majorEvents[len(majorEvents)-5:]
	}

	labelStyle := lipgloss.NewStyle().
		Foreground(theme.TextDim)

	labels := []string{}
	for _, event := range majorEvents {
		label := event.Label
		if label == "" {
			label = string(event.Type)
		}
		labels = append(labels, labelStyle.Render(label))
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, labels...)
}

// GetEventsBetween returns events between two positions
func (t *Timeline) GetEventsBetween(start, end int) []Event {
	if start < 0 {
		start = 0
	}
	if end >= len(t.Events) {
		end = len(t.Events) - 1
	}
	if start > end {
		start, end = end, start
	}

	return t.Events[start : end+1]
}

// GetStats returns timeline statistics
func (t *Timeline) GetStats() map[string]int {
	stats := make(map[string]int)

	for _, event := range t.Events {
		stats[string(event.Type)]++
	}

	return stats
}

// Export exports timeline to a structured format
func (t *Timeline) Export() []map[string]interface{} {
	exported := []map[string]interface{}{}

	for _, event := range t.Events {
		item := map[string]interface{}{
			"type":        string(event.Type),
			"timestamp":   event.Timestamp.Format(time.RFC3339),
			"label":       event.Label,
			"significance": event.Significance,
		}

		if event.Data != nil {
			item["data"] = event.Data
		}

		exported = append(exported, item)
	}

	return exported
}

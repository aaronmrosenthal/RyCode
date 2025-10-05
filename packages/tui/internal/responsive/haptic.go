package responsive

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/sst/opencode/internal/theme"
)

// HapticType represents different haptic feedback patterns
type HapticType string

const (
	HapticLight       HapticType = "light"        // Light tap
	HapticMedium      HapticType = "medium"       // Medium tap
	HapticHeavy       HapticType = "heavy"        // Heavy tap
	HapticSuccess     HapticType = "success"      // Success pattern
	HapticWarning     HapticType = "warning"      // Warning pattern
	HapticError       HapticType = "error"        // Error pattern
	HapticSelection   HapticType = "selection"    // Selection feedback
	HapticImpact      HapticType = "impact"       // Impact feedback
	HapticNotification HapticType = "notification" // Notification
)

// HapticEvent represents a haptic feedback event
type HapticEvent struct {
	Type      HapticType
	Intensity float64 // 0.0 - 1.0
	Duration  time.Duration
	Pattern   []int // Pattern of vibrations in ms
}

// HapticMsg is sent when haptic feedback should occur
type HapticMsg struct {
	Event  HapticEvent
	Visual string // Visual representation for terminal
}

// HapticEngine manages haptic feedback simulation
type HapticEngine struct {
	enabled       bool
	lastHaptic    time.Time
	throttlePeriod time.Duration
	visualMode    bool // Show visual haptic indicators
}

// NewHapticEngine creates a new haptic engine
func NewHapticEngine(enabled bool) *HapticEngine {
	return &HapticEngine{
		enabled:       enabled,
		throttlePeriod: 50 * time.Millisecond, // Prevent haptic spam
		visualMode:    true, // Always show visual for terminal
	}
}

// Enable enables haptic feedback
func (he *HapticEngine) Enable() {
	he.enabled = true
}

// Disable disables haptic feedback
func (he *HapticEngine) Disable() {
	he.enabled = false
}

// IsEnabled returns whether haptic is enabled
func (he *HapticEngine) IsEnabled() bool {
	return he.enabled
}

// Trigger sends a haptic event
func (he *HapticEngine) Trigger(hapticType HapticType) tea.Cmd {
	if !he.enabled {
		return nil
	}

	// Throttle rapid haptic events
	now := time.Now()
	if now.Sub(he.lastHaptic) < he.throttlePeriod {
		return nil
	}
	he.lastHaptic = now

	event := he.createEvent(hapticType)

	return func() tea.Msg {
		return HapticMsg{
			Event:  event,
			Visual: he.createVisual(event),
		}
	}
}

// createEvent creates a haptic event for a given type
func (he *HapticEngine) createEvent(hapticType HapticType) HapticEvent {
	switch hapticType {
	case HapticLight:
		return HapticEvent{
			Type:      HapticLight,
			Intensity: 0.3,
			Duration:  10 * time.Millisecond,
			Pattern:   []int{10},
		}

	case HapticMedium:
		return HapticEvent{
			Type:      HapticMedium,
			Intensity: 0.6,
			Duration:  20 * time.Millisecond,
			Pattern:   []int{20},
		}

	case HapticHeavy:
		return HapticEvent{
			Type:      HapticHeavy,
			Intensity: 1.0,
			Duration:  30 * time.Millisecond,
			Pattern:   []int{30},
		}

	case HapticSuccess:
		return HapticEvent{
			Type:      HapticSuccess,
			Intensity: 0.7,
			Duration:  50 * time.Millisecond,
			Pattern:   []int{10, 10, 30}, // Short-short-long
		}

	case HapticWarning:
		return HapticEvent{
			Type:      HapticWarning,
			Intensity: 0.8,
			Duration:  60 * time.Millisecond,
			Pattern:   []int{20, 10, 20, 10}, // Double pulse
		}

	case HapticError:
		return HapticEvent{
			Type:      HapticError,
			Intensity: 1.0,
			Duration:  80 * time.Millisecond,
			Pattern:   []int{30, 20, 30}, // Triple thud
		}

	case HapticSelection:
		return HapticEvent{
			Type:      HapticSelection,
			Intensity: 0.4,
			Duration:  15 * time.Millisecond,
			Pattern:   []int{15},
		}

	case HapticImpact:
		return HapticEvent{
			Type:      HapticImpact,
			Intensity: 0.9,
			Duration:  25 * time.Millisecond,
			Pattern:   []int{25},
		}

	case HapticNotification:
		return HapticEvent{
			Type:      HapticNotification,
			Intensity: 0.6,
			Duration:  40 * time.Millisecond,
			Pattern:   []int{15, 10, 15}, // Three gentle taps
		}

	default:
		return HapticEvent{
			Type:      HapticLight,
			Intensity: 0.3,
			Duration:  10 * time.Millisecond,
			Pattern:   []int{10},
		}
	}
}

// createVisual creates a visual representation of haptic feedback
func (he *HapticEngine) createVisual(event HapticEvent) string {
	if !he.visualMode {
		return ""
	}

	switch event.Type {
	case HapticLight:
		return "〰️"
	case HapticMedium:
		return "〰️〰️"
	case HapticHeavy:
		return "〰️〰️〰️"
	case HapticSuccess:
		return "✨"
	case HapticWarning:
		return "⚠️"
	case HapticError:
		return "💥"
	case HapticSelection:
		return "👆"
	case HapticImpact:
		return "💫"
	case HapticNotification:
		return "🔔"
	default:
		return "〰️"
	}
}

// RenderVisualFeedback renders visual haptic feedback in terminal
func RenderVisualFeedback(msg HapticMsg, theme *theme.Theme, x, y int) string {
	if msg.Visual == "" {
		return ""
	}

	style := lipgloss.NewStyle().
		Foreground(theme.AccentPrimary).
		Bold(true)

	// Position the visual indicator
	visual := style.Render(msg.Visual)

	// Add some animation-like effect by varying the display
	intensity := int(msg.Event.Intensity * 3)
	if intensity > 0 {
		visual = visual + " " + renderIntensityBar(intensity, theme)
	}

	return visual
}

// renderIntensityBar renders a visual intensity indicator
func renderIntensityBar(intensity int, theme *theme.Theme) string {
	bar := ""
	for i := 0; i < intensity; i++ {
		bar += "▂"
	}

	style := lipgloss.NewStyle().
		Foreground(theme.AccentPrimary)

	return style.Render(bar)
}

// HapticOverlay shows a temporary haptic feedback overlay
type HapticOverlay struct {
	active    bool
	message   string
	startTime time.Time
	duration  time.Duration
	theme     *theme.Theme
}

// NewHapticOverlay creates a haptic feedback overlay
func NewHapticOverlay(theme *theme.Theme) *HapticOverlay {
	return &HapticOverlay{
		theme:    theme,
		duration: 300 * time.Millisecond,
	}
}

// Show displays the haptic overlay
func (ho *HapticOverlay) Show(msg HapticMsg) tea.Cmd {
	ho.active = true
	ho.message = msg.Visual
	ho.startTime = time.Now()

	// Auto-hide after duration
	return tea.Tick(ho.duration, func(t time.Time) tea.Msg {
		return HapticHideMsg{}
	})
}

// Hide hides the overlay
func (ho *HapticOverlay) Hide() {
	ho.active = false
	ho.message = ""
}

// IsActive returns whether overlay is active
func (ho *HapticOverlay) IsActive() bool {
	return ho.active && time.Since(ho.startTime) < ho.duration
}

// Render renders the overlay
func (ho *HapticOverlay) Render(width, height int) string {
	if !ho.IsActive() {
		return ""
	}

	// Center overlay
	style := lipgloss.NewStyle().
		Foreground(ho.theme.AccentPrimary).
		Background(ho.theme.BackgroundSecondary).
		Padding(1, 3).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(ho.theme.AccentPrimary).
		Bold(true)

	overlay := style.Render(ho.message)

	// Center it
	positioned := lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		overlay,
	)

	return positioned
}

// HapticHideMsg signals to hide the haptic overlay
type HapticHideMsg struct{}

// HapticPatterns provides common haptic patterns for different interactions
var HapticPatterns = map[string]HapticType{
	"swipe":         HapticLight,
	"tap":           HapticSelection,
	"long_press":    HapticMedium,
	"button_press":  HapticSelection,
	"scroll":        HapticLight,
	"menu_open":     HapticMedium,
	"menu_close":    HapticLight,
	"message_sent":  HapticSuccess,
	"message_received": HapticNotification,
	"error":         HapticError,
	"warning":       HapticWarning,
	"reaction_add":  HapticSuccess,
	"ai_switch":     HapticImpact,
	"voice_start":   HapticMedium,
	"voice_stop":    HapticLight,
}

// GetPatternForGesture returns appropriate haptic for a gesture
func GetPatternForGesture(gesture GestureType) HapticType {
	switch gesture {
	case GestureSwipeLeft, GestureSwipeRight:
		return HapticLight
	case GestureSwipeUp, GestureSwipeDown:
		return HapticLight
	case GestureTap:
		return HapticSelection
	case GestureDoubleTap:
		return HapticMedium
	case GestureLongPress:
		return HapticHeavy
	default:
		return HapticLight
	}
}

// GetPatternForAction returns appropriate haptic for an action
func GetPatternForAction(action GestureAction) HapticType {
	switch action {
	case ActionReact:
		return HapticSuccess
	case ActionVoiceInput:
		return HapticMedium
	case ActionSwitchAI:
		return HapticImpact
	case ActionDeleteMessage:
		return HapticWarning
	case ActionCopy:
		return HapticSuccess
	case ActionShare:
		return HapticSuccess
	default:
		return HapticSelection
	}
}

// DebugHapticInfo returns debug info about haptic event
func DebugHapticInfo(event HapticEvent) string {
	pattern := ""
	for i, p := range event.Pattern {
		if i > 0 {
			pattern += "-"
		}
		pattern += fmt.Sprintf("%dms", p)
	}

	return fmt.Sprintf(
		"Haptic: %s (%.1f%%, %s, pattern: %s)",
		event.Type,
		event.Intensity*100,
		event.Duration,
		pattern,
	)
}

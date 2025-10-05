package responsive

import (
	"time"

	tea "github.com/charmbracelet/bubbletea/v2"
)

// GestureType represents different gesture types
type GestureType string

const (
	GestureSwipeLeft  GestureType = "swipe_left"
	GestureSwipeRight GestureType = "swipe_right"
	GestureSwipeUp    GestureType = "swipe_up"
	GestureSwipeDown  GestureType = "swipe_down"
	GestureTap        GestureType = "tap"
	GestureDoubleTap  GestureType = "double_tap"
	GestureLongPress  GestureType = "long_press"
	GesturePinch      GestureType = "pinch"
)

// Gesture represents a touch gesture
type Gesture struct {
	Type      GestureType
	StartTime time.Time
	EndTime   time.Time
	StartX    int
	StartY    int
	EndX      int
	EndY      int
	Duration  time.Duration
	Velocity  float64
}

// GestureMsg is sent when a gesture is recognized
type GestureMsg struct {
	Gesture Gesture
	Action  GestureAction
}

// GestureAction defines what action a gesture performs
type GestureAction string

const (
	ActionNextMessage     GestureAction = "next_message"
	ActionPrevMessage     GestureAction = "prev_message"
	ActionOpenMenu        GestureAction = "open_menu"
	ActionCloseMenu       GestureAction = "close_menu"
	ActionReact           GestureAction = "react"
	ActionCopy            GestureAction = "copy"
	ActionShare           GestureAction = "share"
	ActionDeleteMessage   GestureAction = "delete_message"
	ActionShowHistory     GestureAction = "show_history"
	ActionVoiceInput      GestureAction = "voice_input"
	ActionSwitchAI        GestureAction = "switch_ai"
	ActionScrollUp        GestureAction = "scroll_up"
	ActionScrollDown      GestureAction = "scroll_down"
	ActionQuickCommand    GestureAction = "quick_command"
)

// GestureRecognizer tracks and recognizes gestures
type GestureRecognizer struct {
	enabled        bool
	lastTap        time.Time
	lastGesture    Gesture
	gestureStart   time.Time
	startX         int
	startY         int
	currentX       int
	currentY       int
	isTracking     bool
	swipeThreshold int // Minimum distance for swipe
	tapTimeout     time.Duration
}

// NewGestureRecognizer creates a new gesture recognizer
func NewGestureRecognizer() *GestureRecognizer {
	return &GestureRecognizer{
		enabled:        true,
		swipeThreshold: 3, // 3 character widths
		tapTimeout:     300 * time.Millisecond,
	}
}

// Enable enables gesture recognition
func (gr *GestureRecognizer) Enable() {
	gr.enabled = true
}

// Disable disables gesture recognition
func (gr *GestureRecognizer) Disable() {
	gr.enabled = false
}

// IsEnabled returns whether gestures are enabled
func (gr *GestureRecognizer) IsEnabled() bool {
	return gr.enabled
}

// StartTracking begins tracking a gesture
func (gr *GestureRecognizer) StartTracking(x, y int) {
	if !gr.enabled {
		return
	}

	gr.isTracking = true
	gr.gestureStart = time.Now()
	gr.startX = x
	gr.startY = y
	gr.currentX = x
	gr.currentY = y
}

// UpdateTracking updates gesture position
func (gr *GestureRecognizer) UpdateTracking(x, y int) {
	if !gr.enabled || !gr.isTracking {
		return
	}

	gr.currentX = x
	gr.currentY = y
}

// EndTracking completes gesture and returns recognized gesture
func (gr *GestureRecognizer) EndTracking() *Gesture {
	if !gr.enabled || !gr.isTracking {
		return nil
	}

	gr.isTracking = false
	endTime := time.Now()
	duration := endTime.Sub(gr.gestureStart)

	deltaX := gr.currentX - gr.startX
	deltaY := gr.currentY - gr.startY

	gesture := &Gesture{
		StartTime: gr.gestureStart,
		EndTime:   endTime,
		StartX:    gr.startX,
		StartY:    gr.startY,
		EndX:      gr.currentX,
		EndY:      gr.currentY,
		Duration:  duration,
		Velocity:  gr.calculateVelocity(deltaX, deltaY, duration),
	}

	// Recognize gesture type
	gesture.Type = gr.recognizeType(deltaX, deltaY, duration)
	gr.lastGesture = *gesture

	return gesture
}

// recognizeType determines gesture type from movement
func (gr *GestureRecognizer) recognizeType(deltaX, deltaY int, duration time.Duration) GestureType {
	absX := abs(deltaX)
	absY := abs(deltaY)

	// Tap or long press (minimal movement)
	if absX < 2 && absY < 2 {
		// Long press
		if duration > 500*time.Millisecond {
			return GestureLongPress
		}

		// Double tap
		if time.Since(gr.lastTap) < gr.tapTimeout {
			gr.lastTap = time.Time{} // Reset
			return GestureDoubleTap
		}

		gr.lastTap = time.Now()
		return GestureTap
	}

	// Swipe gestures
	if absX > absY {
		// Horizontal swipe
		if absX >= gr.swipeThreshold {
			if deltaX > 0 {
				return GestureSwipeRight
			}
			return GestureSwipeLeft
		}
	} else {
		// Vertical swipe
		if absY >= gr.swipeThreshold {
			if deltaY > 0 {
				return GestureSwipeDown
			}
			return GestureSwipeUp
		}
	}

	return GestureTap
}

// calculateVelocity calculates gesture velocity
func (gr *GestureRecognizer) calculateVelocity(deltaX, deltaY int, duration time.Duration) float64 {
	if duration == 0 {
		return 0
	}

	distance := sqrt(float64(deltaX*deltaX + deltaY*deltaY))
	seconds := duration.Seconds()
	return distance / seconds
}

// MapGestureToAction maps a gesture to an action based on context
func MapGestureToAction(gesture Gesture, context GestureContext) GestureAction {
	switch gesture.Type {
	case GestureSwipeLeft:
		if context.InMessageView {
			return ActionNextMessage
		}
		return ActionOpenMenu

	case GestureSwipeRight:
		if context.InMessageView {
			return ActionPrevMessage
		}
		if context.MenuOpen {
			return ActionCloseMenu
		}
		return ActionShowHistory

	case GestureSwipeUp:
		if context.InMessageView {
			return ActionScrollUp
		}
		return ActionQuickCommand

	case GestureSwipeDown:
		return ActionScrollDown

	case GestureTap:
		// Context-dependent tap
		return ActionNextMessage

	case GestureDoubleTap:
		return ActionReact

	case GestureLongPress:
		if context.OnMessage {
			return ActionCopy
		}
		return ActionVoiceInput

	default:
		return ""
	}
}

// GestureContext provides context for gesture interpretation
type GestureContext struct {
	InMessageView bool
	OnMessage     bool
	MenuOpen      bool
	InputActive   bool
	MessageID     string
}

// Phone-optimized gesture mappings
var PhoneGestures = map[GestureType]GestureAction{
	GestureSwipeLeft:  ActionNextMessage,
	GestureSwipeRight: ActionPrevMessage,
	GestureSwipeUp:    ActionShowHistory,
	GestureSwipeDown:  ActionCloseMenu,
	GestureDoubleTap:  ActionReact,
	GestureLongPress:  ActionVoiceInput,
}

// Tablet-optimized gesture mappings
var TabletGestures = map[GestureType]GestureAction{
	GestureSwipeLeft:  ActionOpenMenu,
	GestureSwipeRight: ActionCloseMenu,
	GestureSwipeUp:    ActionScrollUp,
	GestureSwipeDown:  ActionScrollDown,
	GestureDoubleTap:  ActionReact,
	GestureLongPress:  ActionCopy,
}

// KeyToGesture converts keyboard input to gesture (for desktop testing)
func KeyToGesture(key string) *Gesture {
	now := time.Now()

	switch key {
	case "left", "h":
		return &Gesture{
			Type:      GestureSwipeLeft,
			StartTime: now,
			EndTime:   now,
		}

	case "right", "l":
		return &Gesture{
			Type:      GestureSwipeRight,
			StartTime: now,
			EndTime:   now,
		}

	case "up", "k":
		return &Gesture{
			Type:      GestureSwipeUp,
			StartTime: now,
			EndTime:   now,
		}

	case "down", "j":
		return &Gesture{
			Type:      GestureSwipeDown,
			StartTime: now,
			EndTime:   now,
		}

	case " ", "enter":
		return &Gesture{
			Type:      GestureTap,
			StartTime: now,
			EndTime:   now,
		}

	default:
		return nil
	}
}

// Helper functions
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func sqrt(x float64) float64 {
	// Simple approximation for small numbers
	if x == 0 {
		return 0
	}

	z := x
	for i := 0; i < 10; i++ {
		z = (z + x/z) / 2
	}
	return z
}

// GestureHelpText returns help text for gestures
func GestureHelpText(device DeviceType) string {
	switch device {
	case DevicePhone:
		return `📱 Phone Gestures:
  ← Swipe left: Next message
  → Swipe right: Previous message
  ↑ Swipe up: Show history
  ↓ Swipe down: Close menu
  👆 Double tap: React to message
  🤚 Long press: Voice input`

	case DeviceTablet:
		return `📱 Tablet Gestures:
  ← Swipe left: Open menu
  → Swipe right: Close menu
  ↑ Swipe up: Scroll up
  ↓ Swipe down: Scroll down
  👆 Double tap: React
  🤚 Long press: Copy message`

	default:
		return `⌨️  Keyboard Shortcuts:
  ←/h: Previous
  →/l: Next
  ↑/k: Up
  ↓/j: Down
  Space: Select
  r: React`
	}
}

// GestureUpdate is a Bubble Tea update for gestures
func GestureUpdate(msg tea.Msg, gr *GestureRecognizer, context GestureContext) (*GestureMsg, tea.Cmd) {
	if !gr.IsEnabled() {
		return nil, nil
	}

	// Handle keyboard as gestures for testing on desktop
	if keyMsg, ok := msg.(tea.KeyPressMsg); ok {
		gesture := KeyToGesture(keyMsg.String())
		if gesture != nil {
			action := MapGestureToAction(*gesture, context)
			return &GestureMsg{
				Gesture: *gesture,
				Action:  action,
			}, nil
		}
	}

	// TODO: Handle actual mouse/touch events when available
	// This would require extended terminal capabilities

	return nil, nil
}

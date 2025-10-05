package responsive

import (
	"testing"
	"time"
)

func TestGestureRecognizer_SwipeDetection(t *testing.T) {
	gr := NewGestureRecognizer()

	tests := []struct {
		name     string
		startX   int
		startY   int
		endX     int
		endY     int
		expected GestureType
	}{
		{
			name:     "Swipe Right",
			startX:   0,
			startY:   0,
			endX:     5,
			endY:     0,
			expected: GestureSwipeRight,
		},
		{
			name:     "Swipe Left",
			startX:   5,
			startY:   0,
			endX:     0,
			endY:     0,
			expected: GestureSwipeLeft,
		},
		{
			name:     "Swipe Down",
			startX:   0,
			startY:   0,
			endX:     0,
			endY:     5,
			expected: GestureSwipeDown,
		},
		{
			name:     "Swipe Up",
			startX:   0,
			startY:   5,
			endX:     0,
			endY:     0,
			expected: GestureSwipeUp,
		},
		{
			name:     "Tap (no movement)",
			startX:   5,
			startY:   5,
			endX:     5,
			endY:     5,
			expected: GestureTap,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gr.StartTracking(tt.startX, tt.startY)
			gr.UpdateTracking(tt.endX, tt.endY)
			gesture := gr.EndTracking()

			if gesture == nil {
				t.Fatal("Expected gesture, got nil")
			}

			if gesture.Type != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, gesture.Type)
			}
		})
	}
}

func TestGestureRecognizer_DoubleTap(t *testing.T) {
	gr := NewGestureRecognizer()
	gr.tapTimeout = 300 * time.Millisecond

	// First tap
	gr.StartTracking(5, 5)
	gr.UpdateTracking(5, 5)
	gesture1 := gr.EndTracking()

	if gesture1.Type != GestureTap {
		t.Error("Expected first tap to be GestureTap")
	}

	// Second tap (within timeout)
	time.Sleep(100 * time.Millisecond)
	gr.StartTracking(5, 5)
	gr.UpdateTracking(5, 5)
	gesture2 := gr.EndTracking()

	if gesture2.Type != GestureDoubleTap {
		t.Error("Expected second tap to be GestureDoubleTap")
	}
}

func TestGestureRecognizer_LongPress(t *testing.T) {
	gr := NewGestureRecognizer()

	gr.StartTracking(5, 5)
	time.Sleep(600 * time.Millisecond) // > 500ms threshold
	gr.UpdateTracking(5, 5)
	gesture := gr.EndTracking()

	if gesture == nil {
		t.Fatal("Expected gesture, got nil")
	}

	if gesture.Type != GestureLongPress {
		t.Errorf("Expected GestureLongPress, got %v", gesture.Type)
	}
}

func TestGestureRecognizer_EnableDisable(t *testing.T) {
	gr := NewGestureRecognizer()

	if !gr.IsEnabled() {
		t.Error("Expected gesture recognizer to be enabled by default")
	}

	gr.Disable()

	if gr.IsEnabled() {
		t.Error("Expected gesture recognizer to be disabled")
	}

	gr.StartTracking(0, 0)
	gr.UpdateTracking(5, 0)
	gesture := gr.EndTracking()

	if gesture != nil {
		t.Error("Expected no gesture when disabled")
	}

	gr.Enable()

	if !gr.IsEnabled() {
		t.Error("Expected gesture recognizer to be enabled")
	}
}

func TestMapGestureToAction(t *testing.T) {
	tests := []struct {
		name     string
		gesture  GestureType
		context  GestureContext
		expected GestureAction
	}{
		{
			name:    "Swipe Left in Message View",
			gesture: GestureSwipeLeft,
			context: GestureContext{InMessageView: true},
			expected: ActionNextMessage,
		},
		{
			name:    "Swipe Right in Message View",
			gesture: GestureSwipeRight,
			context: GestureContext{InMessageView: true},
			expected: ActionPrevMessage,
		},
		{
			name:    "Double Tap",
			gesture: GestureDoubleTap,
			context: GestureContext{},
			expected: ActionReact,
		},
		{
			name:    "Long Press",
			gesture: GestureLongPress,
			context: GestureContext{},
			expected: ActionVoiceInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gesture := Gesture{Type: tt.gesture}
			action := MapGestureToAction(gesture, tt.context)

			if action != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, action)
			}
		})
	}
}

func TestKeyToGesture(t *testing.T) {
	tests := []struct {
		key      string
		expected GestureType
	}{
		{"left", GestureSwipeLeft},
		{"right", GestureSwipeRight},
		{"up", GestureSwipeUp},
		{"down", GestureSwipeDown},
		{"h", GestureSwipeLeft},
		{"l", GestureSwipeRight},
		{"k", GestureSwipeUp},
		{"j", GestureSwipeDown},
		{" ", GestureTap},
		{"enter", GestureTap},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			gesture := KeyToGesture(tt.key)

			if gesture == nil {
				t.Fatal("Expected gesture, got nil")
			}

			if gesture.Type != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, gesture.Type)
			}
		})
	}
}

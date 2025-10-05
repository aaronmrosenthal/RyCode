package responsive

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/sst/opencode/internal/theme"
)

// Mock focusable element for testing
type mockFocusable struct {
	id      string
	focused bool
}

func (m *mockFocusable) ID() string                               { return m.id }
func (m *mockFocusable) IsFocused() bool                          { return m.focused }
func (m *mockFocusable) Focus()                                   { m.focused = true }
func (m *mockFocusable) Blur()                                    { m.focused = false }
func (m *mockFocusable) HandleKey(key string) tea.Cmd             { return nil }
func (m *mockFocusable) Render(theme *theme.Theme) string         { return m.id }

func TestFocusManager_Navigation(t *testing.T) {
	fm := NewFocusManager()

	elem1 := &mockFocusable{id: "elem1"}
	elem2 := &mockFocusable{id: "elem2"}
	elem3 := &mockFocusable{id: "elem3"}

	elements := []FocusableElement{elem1, elem2, elem3}
	fm.RegisterZone(ZoneInput, elements)
	fm.SetZone(ZoneInput)

	// Test Next
	fm.Next()
	if !elem2.focused {
		t.Error("Expected elem2 to be focused after Next()")
	}

	// Test Previous
	fm.Previous()
	if !elem1.focused {
		t.Error("Expected elem1 to be focused after Previous()")
	}

	// Test wrap around
	fm.Previous()
	if !elem3.focused {
		t.Error("Expected elem3 to be focused (wrap around)")
	}
}

func TestFocusManager_ZoneSwitching(t *testing.T) {
	fm := NewFocusManager()

	zone1elem := &mockFocusable{id: "zone1"}
	zone2elem := &mockFocusable{id: "zone2"}

	fm.RegisterZone(ZoneInput, []FocusableElement{zone1elem})
	fm.RegisterZone(ZoneMessages, []FocusableElement{zone2elem})

	fm.SetZone(ZoneInput)

	if fm.GetCurrentZone() != ZoneInput {
		t.Error("Expected current zone to be ZoneInput")
	}

	fm.NextZone()

	if fm.GetCurrentZone() != ZoneMessages {
		t.Error("Expected current zone to be ZoneMessages after NextZone()")
	}
}

func TestFocusManager_FocusElement(t *testing.T) {
	fm := NewFocusManager()

	elem1 := &mockFocusable{id: "elem1"}
	elem2 := &mockFocusable{id: "elem2"}

	fm.RegisterZone(ZoneInput, []FocusableElement{elem1, elem2})
	fm.SetZone(ZoneInput)

	fm.FocusElement("elem2")

	if !elem2.focused {
		t.Error("Expected elem2 to be focused")
	}
	if elem1.focused {
		t.Error("Expected elem1 to not be focused")
	}
}

func TestFocusManager_KeyboardMode(t *testing.T) {
	fm := NewFocusManager()

	if fm.IsKeyboardMode() {
		t.Error("Expected keyboard mode to be false initially")
	}

	elem := &mockFocusable{id: "elem"}
	fm.RegisterZone(ZoneInput, []FocusableElement{elem})
	fm.SetZone(ZoneInput)

	fm.Next()

	if !fm.IsKeyboardMode() {
		t.Error("Expected keyboard mode to be true after Next()")
	}

	fm.SetMouseMode()

	if fm.IsKeyboardMode() {
		t.Error("Expected keyboard mode to be false after SetMouseMode()")
	}
}

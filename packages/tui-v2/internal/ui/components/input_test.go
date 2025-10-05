package components

import (
	"strings"
	"testing"
)

func TestNewInputBar(t *testing.T) {
	input := NewInputBar(80)

	if input.Width != 80 {
		t.Errorf("Expected width 80, got %d", input.Width)
	}

	if input.Value != "" {
		t.Error("Expected empty value")
	}

	if input.Cursor != 0 {
		t.Error("Expected cursor at 0")
	}
}

func TestInputBar_Render(t *testing.T) {
	input := NewInputBar(80)
	rendered := input.Render()

	if rendered == "" {
		t.Error("Expected non-empty render output")
	}
}

func TestInputBar_RenderWithValue(t *testing.T) {
	input := NewInputBar(80)
	input.SetValue("Hello, world!")
	rendered := input.Render()

	if !strings.Contains(rendered, "Hello, world!") {
		t.Error("Expected rendered output to contain input value")
	}
}

func TestInputBar_RenderWithGhostText(t *testing.T) {
	input := NewInputBar(80)
	input.SetFocus(true)
	input.SetGhostText(" (suggestion)")
	rendered := input.Render()

	if rendered == "" {
		t.Error("Expected non-empty render output")
	}
}

func TestInputBar_SetValue(t *testing.T) {
	input := NewInputBar(80)
	input.SetValue("Test")

	if input.Value != "Test" {
		t.Errorf("Expected value 'Test', got %s", input.Value)
	}

	if input.Cursor != 4 {
		t.Errorf("Expected cursor at 4, got %d", input.Cursor)
	}
}

func TestInputBar_InsertRune(t *testing.T) {
	input := NewInputBar(80)
	input.InsertRune('H')
	input.InsertRune('i')

	if input.Value != "Hi" {
		t.Errorf("Expected value 'Hi', got %s", input.Value)
	}

	if input.Cursor != 2 {
		t.Errorf("Expected cursor at 2, got %d", input.Cursor)
	}
}

func TestInputBar_DeleteCharBefore(t *testing.T) {
	input := NewInputBar(80)
	input.SetValue("Hello")
	input.DeleteCharBefore()

	if input.Value != "Hell" {
		t.Errorf("Expected value 'Hell', got %s", input.Value)
	}

	if input.Cursor != 4 {
		t.Errorf("Expected cursor at 4, got %d", input.Cursor)
	}
}

func TestInputBar_DeleteCharAfter(t *testing.T) {
	input := NewInputBar(80)
	input.SetValue("Hello")
	input.MoveCursorToStart()
	input.DeleteCharAfter()

	if input.Value != "ello" {
		t.Errorf("Expected value 'ello', got %s", input.Value)
	}

	if input.Cursor != 0 {
		t.Errorf("Expected cursor at 0, got %d", input.Cursor)
	}
}

func TestInputBar_MoveCursor(t *testing.T) {
	input := NewInputBar(80)
	input.SetValue("Test")

	// Test move left
	input.MoveCursorLeft()
	if input.Cursor != 3 {
		t.Errorf("Expected cursor at 3, got %d", input.Cursor)
	}

	// Test move right
	input.MoveCursorRight()
	if input.Cursor != 4 {
		t.Errorf("Expected cursor at 4, got %d", input.Cursor)
	}

	// Test move to start
	input.MoveCursorToStart()
	if input.Cursor != 0 {
		t.Error("Expected cursor at 0")
	}

	// Test move to end
	input.MoveCursorToEnd()
	if input.Cursor != 4 {
		t.Errorf("Expected cursor at 4, got %d", input.Cursor)
	}
}

func TestInputBar_Clear(t *testing.T) {
	input := NewInputBar(80)
	input.SetValue("Test")
	input.SetGhostText("suggestion")
	input.Clear()

	if input.Value != "" {
		t.Error("Expected empty value after clear")
	}

	if input.Cursor != 0 {
		t.Error("Expected cursor at 0 after clear")
	}

	if input.GhostText != "" {
		t.Error("Expected empty ghost text after clear")
	}
}

func TestInputBar_GhostText(t *testing.T) {
	input := NewInputBar(80)
	input.SetValue("Hello")
	input.SetGhostText(" world")

	if input.GhostText != " world" {
		t.Errorf("Expected ghost text ' world', got %s", input.GhostText)
	}

	// Test accepting ghost text
	input.AcceptGhostText()

	if input.Value != "Hello world" {
		t.Errorf("Expected value 'Hello world', got %s", input.Value)
	}

	if input.GhostText != "" {
		t.Error("Expected empty ghost text after accept")
	}

	if input.Cursor != len("Hello world") {
		t.Errorf("Expected cursor at end, got %d", input.Cursor)
	}
}

func TestInputBar_IsEmpty(t *testing.T) {
	input := NewInputBar(80)

	if !input.IsEmpty() {
		t.Error("Expected input to be empty")
	}

	input.SetValue("Test")

	if input.IsEmpty() {
		t.Error("Expected input to not be empty")
	}

	input.SetValue("   ")

	if !input.IsEmpty() {
		t.Error("Expected input with only whitespace to be empty")
	}
}

func TestInputBar_SetFocus(t *testing.T) {
	input := NewInputBar(80)

	if input.Focused {
		t.Error("Expected input to not be focused initially")
	}

	input.SetFocus(true)

	if !input.Focused {
		t.Error("Expected input to be focused")
	}

	input.SetFocus(false)

	if input.Focused {
		t.Error("Expected input to not be focused")
	}
}

func TestInputBar_SetWidth(t *testing.T) {
	input := NewInputBar(80)
	input.SetWidth(100)

	if input.Width != 100 {
		t.Errorf("Expected width 100, got %d", input.Width)
	}
}

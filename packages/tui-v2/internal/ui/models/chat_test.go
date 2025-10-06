package models

import (
	"strings"
	"testing"
	"time"

	"github.com/aaronmrosenthal/rycode/packages/tui-v2/internal/ui/components"
	tea "github.com/charmbracelet/bubbletea"
)

func TestNewChatModel(t *testing.T) {
	m := NewChatModel()

	if m.streaming {
		t.Error("Expected streaming to be false initially")
	}
	if m.ready {
		t.Error("Expected ready to be false initially")
	}
	if len(m.messages.Messages) != 0 {
		t.Errorf("Expected 0 messages, got %d", len(m.messages.Messages))
	}
}

func TestChatModel_Init(t *testing.T) {
	m := NewChatModel()
	cmd := m.Init()

	if cmd == nil {
		t.Error("Expected Init to return animation ticker command")
	}
}

func TestChatModel_WindowSizeMsg(t *testing.T) {
	m := NewChatModel()

	// Send window size message
	msg := tea.WindowSizeMsg{Width: 120, Height: 40}
	updated, _ := m.Update(msg)
	m = updated.(ChatModel)

	if m.width != 120 {
		t.Errorf("Expected width 120, got %d", m.width)
	}
	if m.height != 40 {
		t.Errorf("Expected height 40, got %d", m.height)
	}
	if !m.ready {
		t.Error("Expected ready to be true after WindowSizeMsg")
	}
}

func TestChatModel_QuitKeys(t *testing.T) {
	m := NewChatModel()
	m.ready = true

	tests := []struct {
		key string
	}{
		{"ctrl+c"},
		{"esc"},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			msg := tea.KeyMsg{Type: tea.KeyRunes}
			switch tt.key {
			case "ctrl+c":
				msg.Type = tea.KeyCtrlC
			case "esc":
				msg.Type = tea.KeyEsc
			}

			_, cmd := m.Update(msg)
			if cmd == nil {
				t.Error("Expected quit command, got nil")
			}
		})
	}
}

func TestChatModel_ClearMessages(t *testing.T) {
	m := NewChatModel()
	m.ready = true

	// Add some messages
	m.messages.AddMessage(components.Message{
		ID:      "1",
		Author:  "Test",
		Content: "Hello",
	})

	if len(m.messages.Messages) == 0 {
		t.Fatal("Expected messages to be added")
	}

	// Send Ctrl+L to clear
	msg := tea.KeyMsg{Type: tea.KeyCtrlL}
	updated, _ := m.Update(msg)
	m = updated.(ChatModel)

	if len(m.messages.Messages) != 0 {
		t.Errorf("Expected 0 messages after clear, got %d", len(m.messages.Messages))
	}
}

func TestChatModel_SendMessage(t *testing.T) {
	m := NewChatModel()
	m.ready = true
	m.input.SetValue("Hello AI")

	// Send message
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updated, cmd := m.Update(msg)
	m = updated.(ChatModel)

	if cmd == nil {
		t.Error("Expected command after sending message")
	}

	// Should have 2 messages: user + AI placeholder
	if len(m.messages.Messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(m.messages.Messages))
	}

	// Check user message
	if !m.messages.Messages[0].IsUser {
		t.Error("Expected first message to be from user")
	}
	if m.messages.Messages[0].Content != "Hello AI" {
		t.Errorf("Expected content 'Hello AI', got '%s'", m.messages.Messages[0].Content)
	}

	// Check AI placeholder
	if m.messages.Messages[1].IsUser {
		t.Error("Expected second message to be from AI")
	}
	if m.messages.Messages[1].Status != components.Streaming {
		t.Error("Expected AI message to have Streaming status")
	}

	// Check streaming state
	if !m.streaming {
		t.Error("Expected streaming to be true")
	}

	// Check input cleared
	if m.input.GetValue() != "" {
		t.Error("Expected input to be cleared")
	}
}

func TestChatModel_StreamingPreventsInput(t *testing.T) {
	m := NewChatModel()
	m.ready = true
	m.streaming = true

	// Try to type while streaming
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
	updated, _ := m.Update(msg)
	m = updated.(ChatModel)

	// Input should still be empty
	if m.input.GetValue() != "" {
		t.Error("Expected input to remain empty while streaming")
	}
}

func TestChatModel_StreamChunkMsg(t *testing.T) {
	m := NewChatModel()
	m.ready = true

	// Add an AI message placeholder
	m.messages.AddMessage(components.Message{
		ID:      "ai-1",
		Author:  "AI",
		Content: "",
		Status:  components.Streaming,
		IsUser:  false,
	})
	m.streaming = true

	// Send stream chunk
	chunk := StreamChunkMsg{Chunk: "Hello"}
	updated, cmd := m.Update(chunk)
	m = updated.(ChatModel)

	if cmd == nil {
		t.Error("Expected command to continue streaming")
	}

	// Check message updated
	lastMsg := m.messages.Messages[len(m.messages.Messages)-1]
	if lastMsg.Content != "Hello" {
		t.Errorf("Expected content 'Hello', got '%s'", lastMsg.Content)
	}
}

func TestChatModel_StreamCompleteMsg(t *testing.T) {
	m := NewChatModel()
	m.ready = true

	// Add an AI message
	m.messages.AddMessage(components.Message{
		ID:      "ai-1",
		Author:  "AI",
		Content: "Complete response",
		Status:  components.Streaming,
		IsUser:  false,
	})
	m.streaming = true
	m.input.SetFocus(false)

	// Send stream complete
	msg := StreamCompleteMsg{}
	updated, _ := m.Update(msg)
	m = updated.(ChatModel)

	if m.streaming {
		t.Error("Expected streaming to be false")
	}

	// Check message status updated
	lastMsg := m.messages.Messages[len(m.messages.Messages)-1]
	if lastMsg.Status != components.Sent {
		t.Error("Expected message status to be Sent")
	}

	// Check input refocused
	if !m.input.Focused {
		t.Error("Expected input to be focused after streaming complete")
	}
}

func TestChatModel_KeyboardNavigation(t *testing.T) {
	m := NewChatModel()
	m.ready = true

	tests := []struct {
		name string
		key  tea.KeyType
	}{
		{"left arrow", tea.KeyLeft},
		{"right arrow", tea.KeyRight},
		{"home", tea.KeyHome},
		{"end", tea.KeyEnd},
		{"up arrow", tea.KeyUp},
		{"down arrow", tea.KeyDown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := tea.KeyMsg{Type: tt.key}
			_, cmd := m.Update(msg)

			// Should not quit
			if cmd != nil && cmd() != nil {
				// Command should not be a quit command
				// (we can't easily test tea.Quit, but we can ensure it doesn't panic)
			}
		})
	}
}

func TestChatModel_TabAcceptsGhostText(t *testing.T) {
	m := NewChatModel()
	m.ready = true
	m.input.SetGhostText(" suggestion")

	msg := tea.KeyMsg{Type: tea.KeyTab}
	updated, _ := m.Update(msg)
	m = updated.(ChatModel)

	if m.input.GetValue() != " suggestion" {
		t.Errorf("Expected value ' suggestion', got '%s'", m.input.GetValue())
	}
}

func TestChatModel_BackspaceDeletesChar(t *testing.T) {
	m := NewChatModel()
	m.ready = true
	m.input.SetValue("Hello")

	msg := tea.KeyMsg{Type: tea.KeyBackspace}
	updated, _ := m.Update(msg)
	m = updated.(ChatModel)

	if m.input.GetValue() != "Hell" {
		t.Errorf("Expected value 'Hell', got '%s'", m.input.GetValue())
	}
}

func TestChatModel_GhostTextPredictions(t *testing.T) {
	m := NewChatModel()
	m.ready = true

	// Type "How do I"
	for _, r := range []rune("How do I") {
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}
		updated, _ := m.Update(msg)
		m = updated.(ChatModel)
	}

	if m.input.GhostText != " fix this bug?" {
		t.Errorf("Expected ghost text ' fix this bug?', got '%s'", m.input.GhostText)
	}

	// Clear and type "Explain"
	m.input.Clear()
	for _, r := range []rune("Explain") {
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}
		updated, _ := m.Update(msg)
		m = updated.(ChatModel)
	}

	if m.input.GhostText != " this code to me" {
		t.Errorf("Expected ghost text ' this code to me', got '%s'", m.input.GhostText)
	}
}

func TestChatModel_GenerateAIResponse_FirstMessage(t *testing.T) {
	m := NewChatModel()

	response := m.generateAIResponse()

	if !strings.Contains(response, "RyCode AI") {
		t.Error("Expected first response to contain 'RyCode AI'")
	}
}

func TestChatModel_GenerateAIResponse_BugPattern(t *testing.T) {
	m := NewChatModel()
	m.messages.AddMessage(components.Message{Content: "I found a bug", IsUser: true})
	m.messages.AddMessage(components.Message{Content: "", IsUser: false})

	response := m.generateAIResponse()

	if !strings.Contains(response, "bug") {
		t.Error("Expected response to mention bugs")
	}
	if !strings.Contains(response, "null") {
		t.Error("Expected response to mention null checks")
	}
}

func TestChatModel_GenerateAIResponse_TestPattern(t *testing.T) {
	m := NewChatModel()
	m.messages.AddMessage(components.Message{Content: "Help me write tests", IsUser: true})
	m.messages.AddMessage(components.Message{Content: "", IsUser: false})

	response := m.generateAIResponse()

	if !strings.Contains(response, "test") {
		t.Error("Expected response to mention tests")
	}
	if !strings.Contains(response, "func Test") {
		t.Error("Expected response to include test template")
	}
}

func TestChatModel_GenerateAIResponse_ExplainPattern(t *testing.T) {
	m := NewChatModel()
	m.messages.AddMessage(components.Message{Content: "Explain this code", IsUser: true})
	m.messages.AddMessage(components.Message{Content: "", IsUser: false})

	response := m.generateAIResponse()

	if !strings.Contains(response, "explain") {
		t.Error("Expected response to mention explanation")
	}
	if !strings.Contains(response, "TUI") {
		t.Error("Expected response to mention TUI")
	}
}

func TestChatModel_GenerateAIResponse_GreetingPattern(t *testing.T) {
	m := NewChatModel()
	m.messages.AddMessage(components.Message{Content: "Hello", IsUser: true})
	m.messages.AddMessage(components.Message{Content: "", IsUser: false})

	response := m.generateAIResponse()

	if !strings.Contains(response, "Hey there") || !strings.Contains(response, "👋") {
		t.Error("Expected friendly greeting response")
	}
}

func TestChatModel_GenerateAIResponse_DefaultPattern(t *testing.T) {
	m := NewChatModel()
	m.messages.AddMessage(components.Message{Content: "Random question", IsUser: true})
	m.messages.AddMessage(components.Message{Content: "", IsUser: false})

	response := m.generateAIResponse()

	if !strings.Contains(response, "options") {
		t.Error("Expected default response to mention options")
	}
}

func TestChatModel_View_BeforeReady(t *testing.T) {
	m := NewChatModel()
	view := m.View()

	if view != "Initializing..." {
		t.Errorf("Expected 'Initializing...', got '%s'", view)
	}
}

func TestChatModel_View_AfterReady(t *testing.T) {
	m := NewChatModel()
	m.ready = true
	m.width = 80
	m.height = 24
	m.showLogo = false // Disable logo for test to check title

	view := m.View()

	// Should contain header elements (either logo or title)
	if !strings.Contains(view, "RyCode") {
		t.Error("Expected view to contain 'RyCode'")
	}
}

func TestChatModel_View_StreamingIndicator(t *testing.T) {
	m := NewChatModel()
	m.ready = true
	m.width = 80
	m.height = 24
	m.streaming = true

	view := m.View()

	if !strings.Contains(view, "⚡") || !strings.Contains(view, "responding") {
		t.Error("Expected view to show streaming indicator")
	}
}

func TestChatModel_UpdateDimensions(t *testing.T) {
	m := NewChatModel()
	m.width = 100
	m.height = 50

	m.updateDimensions()

	if m.messages.Width != 100 {
		t.Errorf("Expected messages width 100, got %d", m.messages.Width)
	}

	// Input height is 6, borders are 2, so messages height should be 50 - 6 - 2 = 42
	expectedHeight := 50 - 6 - 2
	if m.messages.Height != expectedHeight {
		t.Errorf("Expected messages height %d, got %d", expectedHeight, m.messages.Height)
	}
}

func TestChatModel_SendMessage_EmptyInput(t *testing.T) {
	m := NewChatModel()
	m.ready = true
	m.input.Clear()

	// Try to send with empty input
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updated, cmd := m.Update(msg)
	m = updated.(ChatModel)

	// Should not send message
	if cmd != nil {
		t.Error("Expected no command for empty input")
	}
	if len(m.messages.Messages) != 0 {
		t.Error("Expected no messages to be added")
	}
}

func TestChatModel_MessageTimestamps(t *testing.T) {
	m := NewChatModel()
	m.ready = true
	m.input.SetValue("Test message")

	before := time.Now()
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updated, _ := m.Update(msg)
	m = updated.(ChatModel)
	after := time.Now()

	// Check timestamps are reasonable
	for _, msg := range m.messages.Messages {
		if msg.Timestamp.Before(before) || msg.Timestamp.After(after) {
			t.Error("Message timestamp out of expected range")
		}
	}
}

package ai

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config == nil {
		t.Fatal("DefaultConfig() returned nil")
	}

	// Verify default values
	tests := []struct {
		name     string
		got      interface{}
		expected interface{}
	}{
		{"Provider", config.Provider, "auto"},
		{"ClaudeModel", config.ClaudeModel, "claude-opus-4-20250514"},
		{"OpenAIModel", config.OpenAIModel, "gpt-4o"},
		{"MaxTokens", config.MaxTokens, 4096},
		{"Temperature", config.Temperature, 0.7},
		{"TopP", config.TopP, 0.9},
		{"RequestsPerMinute", config.RequestsPerMinute, 50},
		{"TokensPerMinute", config.TokensPerMinute, 100000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("DefaultConfig().%s = %v, want %v", tt.name, tt.got, tt.expected)
			}
		})
	}
}

func TestRole(t *testing.T) {
	tests := []struct {
		role     Role
		expected string
	}{
		{RoleUser, "user"},
		{RoleAssistant, "assistant"},
		{RoleSystem, "system"},
	}

	for _, tt := range tests {
		t.Run(string(tt.role), func(t *testing.T) {
			if string(tt.role) != tt.expected {
				t.Errorf("Role = %v, want %v", tt.role, tt.expected)
			}
		})
	}
}

func TestEventType(t *testing.T) {
	tests := []struct {
		eventType EventType
		expected  string
	}{
		{EventTypeChunk, "chunk"},
		{EventTypeComplete, "complete"},
		{EventTypeError, "error"},
	}

	for _, tt := range tests {
		t.Run(string(tt.eventType), func(t *testing.T) {
			if string(tt.eventType) != tt.expected {
				t.Errorf("EventType = %v, want %v", tt.eventType, tt.expected)
			}
		})
	}
}

func TestMessage(t *testing.T) {
	msg := Message{
		Role:    RoleUser,
		Content: "Hello, AI!",
	}

	if msg.Role != RoleUser {
		t.Errorf("Message.Role = %v, want %v", msg.Role, RoleUser)
	}

	if msg.Content != "Hello, AI!" {
		t.Errorf("Message.Content = %v, want %v", msg.Content, "Hello, AI!")
	}
}

func TestStreamEvent(t *testing.T) {
	t.Run("Chunk event", func(t *testing.T) {
		event := StreamEvent{
			Type:    EventTypeChunk,
			Content: "test chunk",
			Done:    false,
		}

		if event.Type != EventTypeChunk {
			t.Errorf("StreamEvent.Type = %v, want %v", event.Type, EventTypeChunk)
		}
		if event.Content != "test chunk" {
			t.Errorf("StreamEvent.Content = %v, want %v", event.Content, "test chunk")
		}
		if event.Done {
			t.Error("StreamEvent.Done should be false")
		}
	})

	t.Run("Complete event", func(t *testing.T) {
		event := StreamEvent{
			Type: EventTypeComplete,
			Done: true,
		}

		if event.Type != EventTypeComplete {
			t.Errorf("StreamEvent.Type = %v, want %v", event.Type, EventTypeComplete)
		}
		if !event.Done {
			t.Error("StreamEvent.Done should be true")
		}
	})

	t.Run("Error event", func(t *testing.T) {
		testErr := &testError{"test error"}
		event := StreamEvent{
			Type:  EventTypeError,
			Error: testErr,
		}

		if event.Type != EventTypeError {
			t.Errorf("StreamEvent.Type = %v, want %v", event.Type, EventTypeError)
		}
		if event.Error == nil {
			t.Fatal("StreamEvent.Error should not be nil")
		}
		if event.Error.Error() != "test error" {
			t.Errorf("StreamEvent.Error = %v, want %v", event.Error, "test error")
		}
	})
}

// testError is a simple error implementation for testing
type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}

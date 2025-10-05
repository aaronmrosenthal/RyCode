package components

import (
	"strings"
	"testing"
	"time"
)

func TestMessageBubble_Render(t *testing.T) {
	msg := Message{
		ID:        "1",
		Author:    "You",
		Content:   "Hello, AI!",
		Timestamp: time.Now(),
		Status:    Sent,
		IsUser:    true,
	}

	bubble := NewMessageBubble(msg, 80)
	rendered := bubble.Render()

	if rendered == "" {
		t.Error("Expected non-empty render output")
	}

	if !strings.Contains(rendered, "You") {
		t.Error("Expected rendered output to contain author name")
	}

	if !strings.Contains(rendered, "Hello, AI!") {
		t.Error("Expected rendered output to contain message content")
	}
}

func TestMessageBubble_RenderMarkdown(t *testing.T) {
	msg := Message{
		ID:        "2",
		Author:    "AI",
		Content:   "# Heading\n\nSome **bold** text.",
		Timestamp: time.Now(),
		Status:    Sent,
		IsUser:    false,
	}

	bubble := NewMessageBubble(msg, 80)
	rendered := bubble.Render()

	if rendered == "" {
		t.Error("Expected non-empty render output")
	}
}

func TestMessageBubble_RenderCodeBlock(t *testing.T) {
	msg := Message{
		ID:      "3",
		Author:  "AI",
		Content: "```go\nfunc main() {\n  println(\"hello\")\n}\n```",
		Timestamp: time.Now(),
		Status:    Sent,
		IsUser:    false,
	}

	bubble := NewMessageBubble(msg, 80)
	rendered := bubble.Render()

	if rendered == "" {
		t.Error("Expected non-empty render output")
	}
}

func TestMessageBubble_RenderReactions(t *testing.T) {
	msg := Message{
		ID:        "4",
		Author:    "You",
		Content:   "Great!",
		Timestamp: time.Now(),
		Status:    Sent,
		Reactions: []string{"👍", "🎉"},
		IsUser:    true,
	}

	bubble := NewMessageBubble(msg, 80)
	rendered := bubble.Render()

	if rendered == "" {
		t.Error("Expected non-empty render output")
	}
}

func TestMessageBubble_RenderStatus(t *testing.T) {
	tests := []struct {
		name   string
		status MessageStatus
	}{
		{"Sending", Sending},
		{"Sent", Sent},
		{"Error", Error},
		{"Streaming", Streaming},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := Message{
				ID:        "5",
				Author:    "You",
				Content:   "Test",
				Timestamp: time.Now(),
				Status:    tt.status,
				IsUser:    true,
			}

			bubble := NewMessageBubble(msg, 80)
			rendered := bubble.Render()

			if rendered == "" {
				t.Error("Expected non-empty render output")
			}
		})
	}
}

func TestFormatTimestamp(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name string
		time time.Time
		want string
	}{
		{"Just now", now.Add(-30 * time.Second), "just now"},
		{"1 minute ago", now.Add(-1 * time.Minute), "1 minute ago"},
		{"5 minutes ago", now.Add(-5 * time.Minute), "5 minutes ago"},
		{"1 hour ago", now.Add(-1 * time.Hour), "1 hour ago"},
		{"2 hours ago", now.Add(-2 * time.Hour), "2 hours ago"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatTimestamp(tt.time)
			if got != tt.want {
				t.Errorf("formatTimestamp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessageList_Render(t *testing.T) {
	messages := []Message{
		{
			ID:        "1",
			Author:    "You",
			Content:   "First message",
			Timestamp: time.Now(),
			Status:    Sent,
			IsUser:    true,
		},
		{
			ID:        "2",
			Author:    "AI",
			Content:   "Second message",
			Timestamp: time.Now(),
			Status:    Sent,
			IsUser:    false,
		},
	}

	list := NewMessageList(messages, 80, 24)
	rendered := list.Render()

	if rendered == "" {
		t.Error("Expected non-empty render output")
	}

	if !strings.Contains(rendered, "First message") {
		t.Error("Expected rendered output to contain first message")
	}

	if !strings.Contains(rendered, "Second message") {
		t.Error("Expected rendered output to contain second message")
	}
}

func TestMessageList_RenderEmpty(t *testing.T) {
	list := NewMessageList([]Message{}, 80, 24)
	rendered := list.Render()

	if rendered == "" {
		t.Error("Expected non-empty render output for empty list")
	}

	if !strings.Contains(rendered, "No messages") {
		t.Error("Expected empty state message")
	}
}

func TestMessageList_AddMessage(t *testing.T) {
	list := NewMessageList([]Message{}, 80, 24)

	if len(list.Messages) != 0 {
		t.Error("Expected empty message list")
	}

	msg := Message{
		ID:        "1",
		Author:    "You",
		Content:   "Test",
		Timestamp: time.Now(),
		Status:    Sent,
		IsUser:    true,
	}

	list.AddMessage(msg)

	if len(list.Messages) != 1 {
		t.Error("Expected 1 message in list")
	}

	if list.Messages[0].Content != "Test" {
		t.Error("Expected message content to be 'Test'")
	}
}

func TestMessageList_UpdateLastMessage(t *testing.T) {
	msg := Message{
		ID:        "1",
		Author:    "AI",
		Content:   "Initial content",
		Timestamp: time.Now(),
		Status:    Streaming,
		IsUser:    false,
	}

	list := NewMessageList([]Message{msg}, 80, 24)
	list.UpdateLastMessage("Updated content")

	if list.Messages[0].Content != "Updated content" {
		t.Errorf("Expected content to be 'Updated content', got %s", list.Messages[0].Content)
	}
}

func TestMessageList_SetLastMessageStatus(t *testing.T) {
	msg := Message{
		ID:        "1",
		Author:    "You",
		Content:   "Test",
		Timestamp: time.Now(),
		Status:    Sending,
		IsUser:    true,
	}

	list := NewMessageList([]Message{msg}, 80, 24)
	list.SetLastMessageStatus(Sent)

	if list.Messages[0].Status != Sent {
		t.Errorf("Expected status to be Sent, got %v", list.Messages[0].Status)
	}
}

func TestMessageList_Scroll(t *testing.T) {
	messages := make([]Message, 10)
	for i := 0; i < 10; i++ {
		messages[i] = Message{
			ID:        string(rune(i)),
			Author:    "Test",
			Content:   "Message",
			Timestamp: time.Now(),
			Status:    Sent,
			IsUser:    i%2 == 0,
		}
	}

	list := NewMessageList(messages, 80, 24)

	// Test scroll up
	list.ScrollUp()
	if list.Offset != 1 {
		t.Errorf("Expected offset 1, got %d", list.Offset)
	}

	// Test scroll down
	list.ScrollDown()
	if list.Offset != 0 {
		t.Errorf("Expected offset 0, got %d", list.Offset)
	}

	// Test scroll to bottom
	list.Offset = 5
	list.ScrollToBottom()
	if list.Offset != 0 {
		t.Errorf("Expected offset 0 after ScrollToBottom, got %d", list.Offset)
	}
}

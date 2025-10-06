package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/aaronmrosenthal/rycode/packages/tui-v2/internal/theme"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

// MessageStatus represents the state of a message
type MessageStatus int

const (
	Sending MessageStatus = iota
	Sent
	Error
	Streaming
)

// Message represents a chat message
type Message struct {
	ID        string
	Author    string // "user" or "ai" or model name
	Content   string
	Timestamp time.Time
	Status    MessageStatus
	Reactions []string
	IsUser    bool
}

// MessageBubble renders a message with markdown and syntax highlighting
type MessageBubble struct {
	Message      Message
	Width        int
	Theme        theme.Theme
	AnimFrame    int    // Animation frame for effects
	ProviderName string // Provider name for branding (e.g., "claude", "gpt-4o")
}

// NewMessageBubble creates a new message bubble
func NewMessageBubble(msg Message, width int) MessageBubble {
	return MessageBubble{
		Message: msg,
		Width:   width,
		Theme:   theme.MatrixTheme,
	}
}

// Render renders the message bubble
func (mb MessageBubble) Render() string {
	// Build header (author + timestamp)
	header := mb.renderHeader()

	// Render content with markdown
	content := mb.renderContent()

	// Reactions
	reactions := mb.renderReactions()

	// Status indicator
	status := mb.renderStatus()

	// Compose all parts
	var parts []string
	parts = append(parts, header)
	parts = append(parts, content)
	if reactions != "" {
		parts = append(parts, reactions)
	}
	if status != "" {
		parts = append(parts, status)
	}

	body := strings.Join(parts, "\n")

	// Apply border based on message type
	var style lipgloss.Style
	if mb.Message.IsUser {
		style = mb.Theme.MessageUser.Width(mb.Width - 4)
	} else {
		style = mb.Theme.MessageAI.Width(mb.Width - 4)
	}

	return style.Render(body)
}

// renderHeader renders the message header (author + timestamp)
func (mb MessageBubble) renderHeader() string {
	author := mb.Message.Author
	timestamp := formatTimestamp(mb.Message.Timestamp)

	// Color author based on type
	var authorText string
	if mb.Message.IsUser {
		// User: cyan with user icon
		authorStyle := lipgloss.NewStyle().
			Foreground(theme.NeonCyan).
			Bold(true)
		authorText = authorStyle.Render("👤 " + author)
	} else {
		// AI: provider-branded with gradient
		providerColor := theme.GetProviderColor(mb.ProviderName)
		icon := theme.GetProviderIcon(mb.ProviderName)

		// Create gradient from provider color to matrix green
		gradientText := theme.GradientText(author, providerColor, theme.MatrixGreen)
		iconStyle := lipgloss.NewStyle().Foreground(providerColor)
		authorText = iconStyle.Render(icon+" ") + gradientText
	}

	timestampStyle := lipgloss.NewStyle().
		Foreground(theme.MatrixGreenDark)

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		authorText,
		timestampStyle.Render(" • "+timestamp),
	)
}

// renderContent renders the message content with markdown
func (mb MessageBubble) renderContent() string {
	if mb.Message.Content == "" {
		return mb.Theme.Dim.Render("(empty message)")
	}

	// Render markdown with Glamour
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(mb.Width-8), // Account for padding and border
	)

	if err != nil {
		// Fallback to plain text if markdown rendering fails
		return mb.Theme.Primary.Render(mb.Message.Content)
	}

	rendered, err := renderer.Render(mb.Message.Content)
	if err != nil {
		return mb.Theme.Primary.Render(mb.Message.Content)
	}

	return strings.TrimSpace(rendered)
}

// renderReactions renders emoji reactions
func (mb MessageBubble) renderReactions() string {
	if len(mb.Message.Reactions) == 0 {
		return ""
	}

	style := lipgloss.NewStyle().
		Foreground(theme.NeonYellow).
		MarginTop(1)

	return style.Render(strings.Join(mb.Message.Reactions, " "))
}

// renderStatus renders the message status indicator
func (mb MessageBubble) renderStatus() string {
	var icon string
	var color lipgloss.Color

	switch mb.Message.Status {
	case Sending:
		icon = "⏳"
		color = theme.NeonYellow
	case Sent:
		icon = "✓"
		color = theme.MatrixGreenDim
	case Error:
		icon = "✗"
		color = theme.NeonPink
	case Streaming:
		// Advanced streaming visualization with animated spinner
		return mb.renderStreamingIndicator()
	default:
		return ""
	}

	style := lipgloss.NewStyle().
		Foreground(color).
		Align(lipgloss.Right)

	return style.Render(icon)
}

// renderStreamingIndicator renders an advanced streaming indicator with animation
func (mb MessageBubble) renderStreamingIndicator() string {
	// Animated spinner frames (braille dots)
	spinners := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	spinnerFrame := mb.AnimFrame % len(spinners)
	spinner := spinners[spinnerFrame]

	// Pulsing "AI is thinking" text
	thinkingText := "AI is thinking"
	intensity := 0.5 + 0.5*(float64(mb.AnimFrame%30)/30.0)
	thinkingColor := theme.InterpolateBrightness(theme.NeonCyan, intensity)
	thinkingStyle := lipgloss.NewStyle().Foreground(thinkingColor)

	// Spinner in cyan
	spinnerStyle := lipgloss.NewStyle().Foreground(theme.NeonCyan).Bold(true)

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		spinnerStyle.Render(spinner+" "),
		thinkingStyle.Render(thinkingText),
	)
}

// formatTimestamp formats a timestamp for display
func formatTimestamp(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	switch {
	case diff < time.Minute:
		return "just now"
	case diff < time.Hour:
		minutes := int(diff.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	case diff < 24*time.Hour:
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	default:
		return t.Format("Jan 2, 15:04")
	}
}

// MessageList renders a list of messages
type MessageList struct {
	Messages []Message
	Width    int
	Height   int
	Offset   int // Scroll offset
}

// NewMessageList creates a new message list
func NewMessageList(messages []Message, width, height int) MessageList {
	return MessageList{
		Messages: messages,
		Width:    width,
		Height:   height,
		Offset:   0,
	}
}

// Render renders the message list
func (ml MessageList) Render() string {
	if len(ml.Messages) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(theme.MatrixGreenDark).
			Italic(true).
			Width(ml.Width).
			Height(ml.Height).
			Align(lipgloss.Center, lipgloss.Center)

		return emptyStyle.Render("No messages yet.\nStart a conversation!")
	}

	var bubbles []string
	for _, msg := range ml.Messages {
		bubble := NewMessageBubble(msg, ml.Width)
		bubbles = append(bubbles, bubble.Render())
	}

	// Join messages with spacing
	content := strings.Join(bubbles, "\n\n")

	// Apply height constraint with scrolling
	style := lipgloss.NewStyle().
		Width(ml.Width).
		MaxHeight(ml.Height)

	return style.Render(content)
}

// ScrollDown scrolls the message list down
func (ml *MessageList) ScrollDown() {
	if ml.Offset > 0 {
		ml.Offset--
	}
}

// ScrollUp scrolls the message list up
func (ml *MessageList) ScrollUp() {
	maxOffset := len(ml.Messages) - 1
	if ml.Offset < maxOffset {
		ml.Offset++
	}
}

// ScrollToBottom scrolls to the most recent message
func (ml *MessageList) ScrollToBottom() {
	ml.Offset = 0
}

// AddMessage adds a new message to the list
func (ml *MessageList) AddMessage(msg Message) {
	ml.Messages = append(ml.Messages, msg)
	ml.ScrollToBottom()
}

// UpdateLastMessage updates the last message in the list
func (ml *MessageList) UpdateLastMessage(content string) {
	if len(ml.Messages) > 0 {
		ml.Messages[len(ml.Messages)-1].Content = content
	}
}

// SetLastMessageStatus sets the status of the last message
func (ml *MessageList) SetLastMessageStatus(status MessageStatus) {
	if len(ml.Messages) > 0 {
		ml.Messages[len(ml.Messages)-1].Status = status
	}
}

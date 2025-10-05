package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/aaronmrosenthal/rycode/packages/tui-v2/internal/layout"
	"github.com/aaronmrosenthal/rycode/packages/tui-v2/internal/theme"
	"github.com/aaronmrosenthal/rycode/packages/tui-v2/internal/ui/components"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// StreamChunkMsg is sent when a new chunk of streaming text arrives
type StreamChunkMsg struct {
	Chunk string
}

// StreamCompleteMsg is sent when streaming is complete
type StreamCompleteMsg struct{}

// ChatModel represents the chat interface
type ChatModel struct {
	messages    components.MessageList
	input       components.InputBar
	width       int
	height      int
	layoutMgr   *layout.LayoutManager
	streaming   bool
	theme       theme.Theme
	ready       bool
}

// NewChatModel creates a new chat model
func NewChatModel() ChatModel {
	return ChatModel{
		messages:  components.NewMessageList([]components.Message{}, 80, 20),
		input:     components.NewInputBar(80),
		layoutMgr: layout.NewLayoutManager(80, 24),
		streaming: false,
		theme:     theme.MatrixTheme,
		ready:     false,
	}
}

// Init initializes the chat model
func (m ChatModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.layoutMgr.Update(msg.Width, msg.Height)
		m.updateDimensions()
		m.ready = true
		return m, nil

	case StreamChunkMsg:
		// Update the last message with new chunk
		m.messages.UpdateLastMessage(m.messages.Messages[len(m.messages.Messages)-1].Content + msg.Chunk)
		return m, m.streamNextChunk()

	case StreamCompleteMsg:
		// Mark streaming as complete
		m.streaming = false
		m.messages.SetLastMessageStatus(components.Sent)
		m.input.SetFocus(true)
		return m, nil
	}

	return m, nil
}

// handleKeyPress handles keyboard input
func (m ChatModel) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global shortcuts
	switch msg.String() {
	case "ctrl+c", "esc":
		return m, tea.Quit

	case "ctrl+l":
		// Clear messages
		m.messages = components.NewMessageList([]components.Message{}, m.messages.Width, m.messages.Height)
		m.input.Clear()
		return m, nil
	}

	// If streaming, don't process input
	if m.streaming {
		return m, nil
	}

	// Handle input bar shortcuts
	switch msg.String() {
	case "enter":
		if !m.input.IsEmpty() {
			return m, m.sendMessage()
		}
		return m, nil

	case "tab":
		// Accept ghost text
		m.input.AcceptGhostText()
		return m, nil

	case "backspace":
		m.input.DeleteCharBefore()
		return m, nil

	case "delete":
		m.input.DeleteCharAfter()
		return m, nil

	case "left":
		m.input.MoveCursorLeft()
		return m, nil

	case "right":
		m.input.MoveCursorRight()
		return m, nil

	case "home", "ctrl+a":
		m.input.MoveCursorToStart()
		return m, nil

	case "end", "ctrl+e":
		m.input.MoveCursorToEnd()
		return m, nil

	case "up":
		// Scroll messages up
		m.messages.ScrollUp()
		return m, nil

	case "down":
		// Scroll messages down
		m.messages.ScrollDown()
		return m, nil

	case "ctrl+d":
		// Scroll to bottom
		m.messages.ScrollToBottom()
		return m, nil

	default:
		// Insert character if it's a single rune
		if len(msg.Runes) == 1 {
			m.input.InsertRune(msg.Runes[0])

			// Simulate ghost text prediction
			if m.input.GetValue() == "How do I" {
				m.input.SetGhostText(" fix this bug?")
			} else if strings.HasPrefix(m.input.GetValue(), "Explain") {
				m.input.SetGhostText(" this code to me")
			} else {
				m.input.SetGhostText("")
			}
		}
		return m, nil
	}
}

// sendMessage sends a message and triggers AI response
func (m *ChatModel) sendMessage() tea.Cmd {
	// Create user message
	userMsg := components.Message{
		ID:        fmt.Sprintf("user-%d", time.Now().Unix()),
		Author:    "You",
		Content:   m.input.GetValue(),
		Timestamp: time.Now(),
		Status:    components.Sent,
		IsUser:    true,
	}

	// Add to message list
	m.messages.AddMessage(userMsg)

	// Clear input
	m.input.Clear()
	m.input.SetFocus(false)

	// Create AI response placeholder
	aiMsg := components.Message{
		ID:        fmt.Sprintf("ai-%d", time.Now().Unix()),
		Author:    "AI",
		Content:   "",
		Timestamp: time.Now(),
		Status:    components.Streaming,
		IsUser:    false,
	}

	m.messages.AddMessage(aiMsg)
	m.streaming = true

	// Start streaming response
	return m.streamNextChunk()
}

// streamNextChunk simulates streaming AI response
func (m *ChatModel) streamNextChunk() tea.Cmd {
	// Simulate AI response based on user input
	response := m.generateAIResponse()
	words := strings.Fields(response)

	currentContent := ""
	if len(m.messages.Messages) > 0 {
		currentContent = m.messages.Messages[len(m.messages.Messages)-1].Content
	}

	// Get next word to stream
	currentWords := strings.Fields(currentContent)
	if len(currentWords) < len(words) {
		nextWord := words[len(currentWords)]

		return tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
			if len(currentWords) == 0 {
				return StreamChunkMsg{Chunk: nextWord}
			}
			return StreamChunkMsg{Chunk: " " + nextWord}
		})
	}

	// Streaming complete
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return StreamCompleteMsg{}
	})
}

// generateAIResponse generates a simulated AI response
func (m *ChatModel) generateAIResponse() string {
	if len(m.messages.Messages) < 2 {
		return "Hello! I'm RyCode AI, powered by the Matrix TUI. How can I help you today?"
	}

	lastUserMsg := m.messages.Messages[len(m.messages.Messages)-2].Content

	// Pattern-based responses
	if strings.Contains(strings.ToLower(lastUserMsg), "bug") {
		return "I'll analyze the code for bugs. Based on the context, I recommend:\n\n1. Check for null/undefined values\n2. Add error handling\n3. Validate input parameters\n\nWould you like me to show specific examples?"
	}

	if strings.Contains(strings.ToLower(lastUserMsg), "test") {
		return "I can help you write tests! Here's a test template:\n\n```go\nfunc TestExample(t *testing.T) {\n  // Arrange\n  input := \"test\"\n  \n  // Act\n  result := YourFunction(input)\n  \n  // Assert\n  if result != expected {\n    t.Errorf(\"got %v, want %v\", result, expected)\n  }\n}\n```"
	}

	if strings.Contains(strings.ToLower(lastUserMsg), "explain") {
		return "Let me explain this code:\n\nThis implements a **responsive TUI framework** with:\n- Device detection (phone/tablet/desktop)\n- Dynamic layout switching\n- Theme system with Matrix aesthetics\n\nThe key insight is using terminal dimensions to adapt the UI automatically!"
	}

	if strings.Contains(strings.ToLower(lastUserMsg), "hello") || strings.Contains(strings.ToLower(lastUserMsg), "hi") {
		return "Hey there! 👋 I'm here to help with coding, debugging, and explanations. What are you working on?"
	}

	// Default response
	return "Interesting! I can help with that. Some options:\n\n1. 🔧 Fix bugs\n2. 🧪 Write tests\n3. 📝 Add documentation\n4. ⚡ Optimize performance\n\nWhat would you like to focus on?"
}

// updateDimensions updates component dimensions based on layout
func (m *ChatModel) updateDimensions() {
	// Calculate available space
	inputHeight := 6 // Input bar + buttons + actions
	messagesHeight := m.height - inputHeight - 2 // -2 for borders

	// Update message list dimensions
	m.messages.Width = m.width
	m.messages.Height = messagesHeight

	// Update input bar width
	m.input.SetWidth(m.width)
}

// View renders the chat interface
func (m ChatModel) View() string {
	if !m.ready {
		return "Initializing..."
	}

	// Header
	deviceClass := m.layoutMgr.GetDeviceClass()
	header := m.renderHeader(deviceClass)

	// Messages area
	messagesView := m.messages.Render()

	// Separator
	separator := strings.Repeat("─", m.width)
	separatorStyle := lipgloss.NewStyle().Foreground(theme.MatrixGreenDark)

	// Input area
	inputView := m.input.Render()

	// Status bar
	statusBar := m.renderStatusBar()

	// Compose layout
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		messagesView,
		separatorStyle.Render(separator),
		inputView,
		statusBar,
	)

	return content
}

// renderHeader renders the chat header
func (m ChatModel) renderHeader(deviceClass layout.DeviceClass) string {
	title := theme.GradientTextPreset("RyCode Matrix TUI", theme.GradientMatrix)
	subtitle := m.theme.Subtitle.Render(fmt.Sprintf("Device: %s • %dx%d", deviceClass, m.width, m.height))

	headerStyle := lipgloss.NewStyle().
		Padding(0, 1).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderBottom(true).
		BorderForeground(theme.MatrixGreen)

	return headerStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			subtitle,
		),
	)
}

// renderStatusBar renders the status bar
func (m ChatModel) renderStatusBar() string {
	var status string

	if m.streaming {
		status = m.theme.Info.Render("⚡ AI is responding...")
	} else if m.input.Focused {
		status = m.theme.Hint.Render("Press Enter to send • Tab to accept suggestion • Ctrl+L to clear • Esc to quit")
	} else {
		status = m.theme.Hint.Render("Type to start • Ctrl+C to quit")
	}

	messageCount := fmt.Sprintf("%d messages", len(m.messages.Messages))

	statusStyle := lipgloss.NewStyle().
		Foreground(theme.MatrixGreenDim).
		Padding(0, 1)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		statusStyle.Render(status),
		lipgloss.NewStyle().Foreground(theme.MatrixGreenDark).Render(" │ "),
		statusStyle.Render(messageCount),
	)
}

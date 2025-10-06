package models

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aaronmrosenthal/rycode/packages/tui-v2/internal/ai"
	_ "github.com/aaronmrosenthal/rycode/packages/tui-v2/internal/ai/providers" // Register providers
	"github.com/aaronmrosenthal/rycode/packages/tui-v2/internal/layout"
	"github.com/aaronmrosenthal/rycode/packages/tui-v2/internal/theme"
	"github.com/aaronmrosenthal/rycode/packages/tui-v2/internal/ui/components"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	// StreamDelayMs is the delay between streaming chunks in milliseconds
	StreamDelayMs = 50
	// StreamCompleteDelayMs is the delay before marking stream complete
	StreamCompleteDelayMs = 100
	// InputBarHeight is the height of the input bar area
	InputBarHeight = 6
	// BorderHeight is the height of borders
	BorderHeight = 2
)

// StreamChunkMsg is sent when a new chunk of streaming text arrives
type StreamChunkMsg struct {
	Chunk        string
	TokensUsed   int // Tokens in this chunk (0 if unknown)
	PromptTokens int // Prompt tokens (set on first chunk)
}

// StreamCompleteMsg is sent when streaming is complete
type StreamCompleteMsg struct{}

// TokenUpdateMsg is sent to update token counters (thread-safe)
type TokenUpdateMsg struct {
	PromptTokens   int
	ResponseTokens int
}

// ChatModel represents the chat interface
type ChatModel struct {
	messages           components.MessageList
	input              components.InputBar
	width              int
	height             int
	layoutMgr          *layout.LayoutManager
	streaming          bool
	theme              theme.Theme
	ready              bool
	aiProvider         ai.Provider
	aiEnabled          bool
	aiError            error
	streamChan         <-chan ai.StreamEvent
	streamActive       bool
	sessionTokens      int                // Total tokens used this session
	lastPromptTokens   int                // Tokens in last prompt
	lastResponseTokens int                // Tokens in last response
	activeCtx          context.Context    // Context for active AI request
	cancelRequest      context.CancelFunc // Cancel function for active request
}

// NewChatModel creates a new chat model
func NewChatModel() ChatModel {
	// Try to initialize AI provider
	provider, err := ai.NewProvider(nil) // nil = use env config

	aiEnabled := err == nil
	if !aiEnabled {
		// Fall back to mock responses
		provider = nil
	}

	return ChatModel{
		messages:   components.NewMessageList([]components.Message{}, 80, 20),
		input:      components.NewInputBar(80),
		layoutMgr:  layout.NewLayoutManager(80, 24),
		streaming:  false,
		theme:      theme.MatrixTheme,
		ready:      false,
		aiProvider: provider,
		aiEnabled:  aiEnabled,
		aiError:    err,
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
		if len(m.messages.Messages) == 0 {
			return m, nil
		}
		lastMsg := m.messages.Messages[len(m.messages.Messages)-1]
		m.messages.UpdateLastMessage(lastMsg.Content + msg.Chunk)

		// Update token counters (thread-safe)
		if msg.PromptTokens > 0 {
			m.lastPromptTokens = msg.PromptTokens
		}
		if msg.TokensUsed > 0 {
			m.sessionTokens += msg.TokensUsed
			m.lastResponseTokens += msg.TokensUsed
		}

		// Continue streaming from active channel or fall back to mock
		if m.streamActive && m.streamChan != nil {
			return m, m.waitForNextStreamEvent()
		}
		return m, m.streamNextChunk()

	case TokenUpdateMsg:
		// Update token counters (thread-safe message-based update)
		if msg.PromptTokens > 0 {
			m.lastPromptTokens = msg.PromptTokens
		}
		if msg.ResponseTokens > 0 {
			m.sessionTokens += msg.ResponseTokens
			m.lastResponseTokens += msg.ResponseTokens
		}
		return m, nil

	case StreamCompleteMsg:
		// Mark streaming as complete
		m.streaming = false
		m.streamActive = false
		m.streamChan = nil

		// Clean up context and cancel function
		if m.cancelRequest != nil {
			m.cancelRequest()
			m.cancelRequest = nil
		}
		m.activeCtx = nil

		if len(m.messages.Messages) > 0 {
			m.messages.SetLastMessageStatus(components.Sent)
		}
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
		// Cancel active AI request if streaming
		if m.cancelRequest != nil {
			m.cancelRequest()
		}
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

	// Store prompt for AI
	prompt := m.input.GetValue()

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

	// Use real AI if available, otherwise fall back to mock
	if m.aiEnabled && m.aiProvider != nil {
		return m.streamRealAI(prompt)
	}

	// Fall back to mock streaming
	return m.streamNextChunk()
}

// streamRealAI handles real AI provider streaming
func (m *ChatModel) streamRealAI(prompt string) tea.Cmd {
	return func() tea.Msg {
		// Build conversation history
		history := make([]ai.Message, 0, len(m.messages.Messages)-1)
		for i := 0; i < len(m.messages.Messages)-1; i++ {
			msg := m.messages.Messages[i]
			role := ai.RoleUser
			if !msg.IsUser {
				role = ai.RoleAssistant
			}
			history = append(history, ai.Message{
				Role:    role,
				Content: msg.Content,
			})
		}

		// Estimate prompt tokens
		estimatedPromptTokens := ai.EstimateConversationTokens(history) + ai.EstimateTokens(prompt)

		// Create cancellable context with 2 minute timeout
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		m.activeCtx = ctx
		m.cancelRequest = cancel

		// Start streaming from AI provider
		eventCh, err := m.aiProvider.Stream(ctx, prompt, history)
		if err != nil {
			return StreamChunkMsg{Chunk: fmt.Sprintf("❌ Error: %v", err)}
		}

		// Store channel and wait for first event
		m.streamChan = eventCh
		m.streamActive = true

		// Wait for first event
		event, ok := <-eventCh
		if !ok {
			return StreamCompleteMsg{}
		}

		// Prepare token counts for first chunk
		promptTokens := event.PromptTokens
		if promptTokens == 0 {
			promptTokens = estimatedPromptTokens
		}

		tokensUsed := event.TokensUsed
		if event.Content != "" && tokensUsed == 0 {
			// Estimate if provider doesn't give us exact count
			tokensUsed = ai.EstimateTokens(event.Content)
		}

		switch event.Type {
		case ai.EventTypeChunk:
			return StreamChunkMsg{
				Chunk:        event.Content,
				TokensUsed:   tokensUsed,
				PromptTokens: promptTokens,
			}
		case ai.EventTypeComplete:
			return StreamCompleteMsg{}
		case ai.EventTypeError:
			return StreamChunkMsg{Chunk: fmt.Sprintf("\n\n❌ Error: %v", event.Error)}
		default:
			return StreamCompleteMsg{}
		}
	}
}

// streamNextChunk simulates streaming AI response (fallback when no AI provider)
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

		return tea.Tick(StreamDelayMs*time.Millisecond, func(t time.Time) tea.Msg {
			if len(currentWords) == 0 {
				return StreamChunkMsg{Chunk: nextWord}
			}
			return StreamChunkMsg{Chunk: " " + nextWord}
		})
	}

	// Streaming complete
	return tea.Tick(StreamCompleteDelayMs*time.Millisecond, func(t time.Time) tea.Msg {
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
	messagesHeight := m.height - InputBarHeight - BorderHeight

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
		provider := "Mock"
		if m.aiEnabled && m.aiProvider != nil {
			provider = fmt.Sprintf("%s (%s)", m.aiProvider.Name(), m.aiProvider.Model())
		}
		status = m.theme.Info.Render(fmt.Sprintf("⚡ %s is responding...", provider))
	} else if m.input.Focused {
		status = m.theme.Hint.Render("Press Enter to send • Tab to accept suggestion • Ctrl+L to clear • Esc to quit")
	} else {
		status = m.theme.Hint.Render("Type to start • Ctrl+C to quit")
	}

	// Add AI provider info with token usage
	aiInfo := "Mock AI"
	if m.aiEnabled && m.aiProvider != nil {
		providerName := fmt.Sprintf("%s (%s)", m.aiProvider.Name(), m.aiProvider.Model())
		if m.sessionTokens > 0 {
			aiInfo = fmt.Sprintf("%s • %d tokens", providerName, m.sessionTokens)
		} else {
			aiInfo = providerName
		}
	} else if m.aiError != nil {
		aiInfo = "⚠️  No AI (set ANTHROPIC_API_KEY or OPENAI_API_KEY)"
	}

	messageCount := fmt.Sprintf("%d messages", len(m.messages.Messages))

	statusStyle := lipgloss.NewStyle().
		Foreground(theme.MatrixGreenDim).
		Padding(0, 1)

	separator := lipgloss.NewStyle().Foreground(theme.MatrixGreenDark).Render(" │ ")

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		statusStyle.Render(status),
		separator,
		statusStyle.Render(aiInfo),
		separator,
		statusStyle.Render(messageCount),
	)
}

// waitForNextStreamEvent waits for the next stream event from active AI channel
func (m *ChatModel) waitForNextStreamEvent() tea.Cmd {
	return func() tea.Msg {
		event, ok := <-m.streamChan
		if !ok {
			// Channel closed
			return StreamCompleteMsg{}
		}

		// Prepare token info (don't mutate state here - that happens in Update())
		promptTokens := event.PromptTokens
		tokensUsed := event.TokensUsed
		if event.Content != "" && tokensUsed == 0 {
			// Estimate if provider doesn't give exact count
			tokensUsed = ai.EstimateTokens(event.Content)
		}

		switch event.Type {
		case ai.EventTypeChunk:
			return StreamChunkMsg{
				Chunk:        event.Content,
				TokensUsed:   tokensUsed,
				PromptTokens: promptTokens,
			}
		case ai.EventTypeComplete:
			return StreamCompleteMsg{}
		case ai.EventTypeError:
			return StreamChunkMsg{Chunk: fmt.Sprintf("\n\n❌ Error: %v", event.Error)}
		default:
			return StreamCompleteMsg{}
		}
	}
}

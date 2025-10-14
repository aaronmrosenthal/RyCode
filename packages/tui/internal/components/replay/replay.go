package replay

import (
	"time"

	"github.com/charmbracelet/lipgloss/v2"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/aaronmrosenthal/rycode/internal/theme"
)

// Message represents a conversation message
type Message struct {
	Role      string
	Content   string
	Timestamp time.Time
	Tools     []string
}

// ReplayState tracks replay status
type ReplayState struct {
	Messages      []Message
	CurrentIndex  int
	Playing       bool
	Speed         float64 // 1.0 = normal, 0.5 = slow-mo, 2.0 = fast
	ShowThinking  bool
	ExplainMode   bool
}

// ReplayModel handles instant replay functionality
type ReplayModel struct {
	state  *ReplayState
	width  int
	height int
	theme  *theme.Theme
}

// NewReplayModel creates a new replay model
func NewReplayModel(messages []Message, theme *theme.Theme) *ReplayModel {
	return &ReplayModel{
		state: &ReplayState{
			Messages:     messages,
			CurrentIndex: 0,
			Playing:      false,
			Speed:        1.0,
			ShowThinking: false,
			ExplainMode:  false,
		},
		theme: theme,
	}
}

// TickMsg for replay animation
type TickMsg struct{}

// Init initializes the replay
func (m *ReplayModel) Init() tea.Cmd {
	return nil
}

// Update handles replay updates
func (m *ReplayModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "space":
			// Toggle play/pause
			m.state.Playing = !m.state.Playing
			if m.state.Playing {
				return m, m.tick()
			}
			return m, nil

		case "left", "h":
			// Previous message
			if m.state.CurrentIndex > 0 {
				m.state.CurrentIndex--
			}
			return m, nil

		case "right", "l":
			// Next message
			if m.state.CurrentIndex < len(m.state.Messages)-1 {
				m.state.CurrentIndex++
			}
			return m, nil

		case "1":
			// Slow motion (0.5x)
			m.state.Speed = 0.5
			return m, nil

		case "2":
			// Normal speed (1.0x)
			m.state.Speed = 1.0
			return m, nil

		case "3":
			// Fast forward (2.0x)
			m.state.Speed = 2.0
			return m, nil

		case "t":
			// Toggle thinking mode
			m.state.ShowThinking = !m.state.ShowThinking
			return m, nil

		case "e":
			// Toggle explain mode
			m.state.ExplainMode = !m.state.ExplainMode
			return m, nil

		case "r":
			// Restart from beginning
			m.state.CurrentIndex = 0
			m.state.Playing = true
			return m, m.tick()

		case "esc", "q":
			// Exit replay
			return m, tea.Quit
		}

	case TickMsg:
		if m.state.Playing && m.state.CurrentIndex < len(m.state.Messages)-1 {
			m.state.CurrentIndex++
			return m, m.tick()
		}
		if m.state.CurrentIndex >= len(m.state.Messages)-1 {
			m.state.Playing = false
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	return m, nil
}

// View renders the replay
func (m *ReplayModel) View() string {
	if len(m.state.Messages) == 0 {
		return "No messages to replay"
	}

	sections := []string{}

	// Header
	sections = append(sections, m.renderHeader())

	// Current message
	sections = append(sections, m.renderCurrentMessage())

	// Explanation (if in explain mode)
	if m.state.ExplainMode {
		sections = append(sections, m.renderExplanation())
	}

	// Controls
	sections = append(sections, m.renderControls())

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// Helper functions
func (m *ReplayModel) tick() tea.Cmd {
	delay := time.Duration(1000.0 / m.state.Speed)
	return tea.Tick(delay*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg{}
	})
}

func (m *ReplayModel) renderHeader() string {
	titleStyle := lipgloss.NewStyle().
		Foreground((*m.theme).Primary()).
		Bold(true)

	title := titleStyle.Render("🎬 Instant Replay")

	// Progress
	progress := float64(m.state.CurrentIndex) / float64(len(m.state.Messages)-1)
	if progress > 1.0 {
		progress = 1.0
	}

	progressWidth := 40
	filled := int(progress * float64(progressWidth))
	empty := progressWidth - filled

	progressStyle := lipgloss.NewStyle().
		Foreground((*m.theme).Primary())

	emptyStyle := lipgloss.NewStyle().
		Foreground((*m.theme).Border())

	progressBar := "[" +
		progressStyle.Render(lipgloss.NewStyle().Render(lipgloss.PlaceHorizontal(filled, lipgloss.Left, "█", lipgloss.WithWhitespaceChars("█")))) +
		emptyStyle.Render(lipgloss.NewStyle().Render(lipgloss.PlaceHorizontal(empty, lipgloss.Left, "░", lipgloss.WithWhitespaceChars("░")))) +
		"]"

	// Position info
	positionStyle := lipgloss.NewStyle().
		Foreground((*m.theme).TextMuted())

	position := positionStyle.Render(
		" " + string(rune('0'+m.state.CurrentIndex+1)) + "/" + string(rune('0'+len(m.state.Messages))),
	)

	// Speed indicator
	speedStyle := lipgloss.NewStyle().
		Foreground((*m.theme).Info())

	speedText := ""
	if m.state.Speed == 0.5 {
		speedText = " 🐌 0.5x"
	} else if m.state.Speed == 2.0 {
		speedText = " ⚡ 2.0x"
	} else {
		speedText = " ▶ 1.0x"
	}
	speed := speedStyle.Render(speedText)

	header := lipgloss.JoinHorizontal(
		lipgloss.Left,
		title,
		"  ",
		progressBar,
		position,
		speed,
	)

	return lipgloss.NewStyle().
		MarginBottom(1).
		Render(header)
}

func (m *ReplayModel) renderCurrentMessage() string {
	if m.state.CurrentIndex >= len(m.state.Messages) {
		return ""
	}

	msg := m.state.Messages[m.state.CurrentIndex]

	// Role indicator
	roleStyle := lipgloss.NewStyle().
		Bold(true)

	var role string
	if msg.Role == "user" {
		role = roleStyle.
			Foreground((*m.theme).Secondary()).
			Render("You:")
	} else {
		role = roleStyle.
			Foreground((*m.theme).Primary()).
			Render("AI:")
	}

	// Timestamp
	timeStyle := lipgloss.NewStyle().
		Foreground((*m.theme).TextMuted())

	timestamp := timeStyle.Render(msg.Timestamp.Format("15:04:05"))

	// Content
	contentStyle := lipgloss.NewStyle().
		Foreground((*m.theme).Text()).
		Width(m.width - 4).
		MarginTop(1).
		MarginBottom(1)

	content := contentStyle.Render(msg.Content)

	// Tools used (if any)
	toolsSection := ""
	if len(msg.Tools) > 0 && m.state.ShowThinking {
		toolStyle := lipgloss.NewStyle().
			Foreground((*m.theme).Info()).
			Italic(true)

		tools := "🛠️  Tools used: " + lipgloss.JoinHorizontal(lipgloss.Left, msg.Tools...)
		toolsSection = toolStyle.Render(tools)
	}

	header := lipgloss.JoinHorizontal(
		lipgloss.Left,
		role,
		"  ",
		timestamp,
	)

	sections := []string{header, content}
	if toolsSection != "" {
		sections = append(sections, toolsSection)
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m *ReplayModel) renderExplanation() string {
	if m.state.CurrentIndex >= len(m.state.Messages) {
		return ""
	}

	msg := m.state.Messages[m.state.CurrentIndex]

	// Generate explanation based on message
	explanation := m.generateExplanation(msg)

	explainStyle := lipgloss.NewStyle().
		Foreground((*m.theme).Warning()).
		Background((*m.theme).BackgroundPanel()).
		Padding(1).
		MarginTop(1).
		Width(m.width - 4)

	return explainStyle.Render("💡 " + explanation)
}

func (m *ReplayModel) renderControls() string {
	controlStyle := lipgloss.NewStyle().
		Foreground((*m.theme).TextMuted()).
		MarginTop(1)

	controls := []string{
		"Space: Play/Pause",
		"←/→: Navigate",
		"1/2/3: Speed",
		"T: Thinking",
		"E: Explain",
		"R: Restart",
		"Q: Exit",
	}

	return controlStyle.Render(lipgloss.JoinHorizontal(lipgloss.Left, controls...))
}

func (m *ReplayModel) generateExplanation(msg Message) string {
	// Simple explanation generator
	// In production, this could use AI to generate explanations

	if msg.Role == "user" {
		return "User is asking: " + msg.Content
	}

	explanations := []string{
		"AI analyzed the request and generated a response",
		"AI used its knowledge to provide this answer",
		"AI considered the context before responding",
	}

	if len(msg.Tools) > 0 {
		return "AI used tools (" + lipgloss.JoinHorizontal(lipgloss.Left, msg.Tools...) + ") to generate this response"
	}

	return explanations[m.state.CurrentIndex%len(explanations)]
}

// GetState returns current replay state
func (m *ReplayModel) GetState() *ReplayState {
	return m.state
}

// SetSpeed sets replay speed
func (m *ReplayModel) SetSpeed(speed float64) {
	m.state.Speed = speed
}

// ToggleExplain toggles explanation mode
func (m *ReplayModel) ToggleExplain() {
	m.state.ExplainMode = !m.state.ExplainMode
}

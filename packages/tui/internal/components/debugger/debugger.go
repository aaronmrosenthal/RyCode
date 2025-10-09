package debugger

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/aaronmrosenthal/rycode-sdk-go"
	"github.com/aaronmrosenthal/rycode/internal/styles"
	"github.com/aaronmrosenthal/rycode/internal/theme"
)

// DebuggerMsg is sent when the debugger is activated
type DebuggerMsg struct {
	SessionID string
	Program   string
	Language  string
}

// DebuggerStoppedMsg is sent when execution is paused
type DebuggerStoppedMsg struct {
	File   string
	Line   int
	Reason string
}

// DebuggerState represents the current debugging state
type DebuggerState int

const (
	StateInactive DebuggerState = iota
	StateInitializing
	StateRunning
	StatePaused
	StateStopped
)

type Model struct {
	width  int
	height int

	// Debugger state
	state     DebuggerState
	sessionID string
	program   string
	language  string

	// Current execution point
	currentFile string
	currentLine int
	stopReason  string

	// Panels
	sourceView    SourceViewModel
	variablesView VariablesViewModel
	callStackView CallStackViewModel

	// Active panel
	activePanel int // 0=source, 1=variables, 2=callstack

	// Theme
	theme theme.Theme

	// HTTP client for API calls
	client *opencode.Client
}

func New(width, height int, client *opencode.Client) Model {
	return Model{
		width:         width,
		height:        height,
		state:         StateInactive,
		activePanel:   0,
		theme:         theme.CurrentTheme(),
		sourceView:    NewSourceView(width/2, height-4),
		variablesView: NewVariablesView(width/2, height/2-2),
		callStackView: NewCallStackView(width/2, height/2-2),
		client:        client,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateLayout()
		return m, nil

	case DebuggerMsg:
		m.state = StateInitializing
		m.sessionID = msg.SessionID
		m.program = msg.Program
		m.language = msg.Language
		return m, nil

	case DebuggerStoppedMsg:
		m.state = StatePaused
		m.currentFile = msg.File
		m.currentLine = msg.Line
		m.stopReason = msg.Reason
		m.sourceView = m.sourceView.UpdateCurrentLine(msg.File, msg.Line)
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}

	return m, nil
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (Model, tea.Cmd) {
	// Only handle keys when debugger is active
	if m.state == StateInactive {
		return m, nil
	}

	switch msg.String() {
	case "tab":
		// Cycle through panels
		m.activePanel = (m.activePanel + 1) % 3
		return m, nil

	case "c":
		// Continue execution
		if m.state == StatePaused {
			m.state = StateRunning
			return m, m.sendDebugCommand("continue")
		}
		return m, nil

	case "s":
		// Step over
		if m.state == StatePaused {
			return m, m.sendDebugCommand("step-over")
		}
		return m, nil

	case "i":
		// Step into
		if m.state == StatePaused {
			return m, m.sendDebugCommand("step-into")
		}
		return m, nil

	case "o":
		// Step out
		if m.state == StatePaused {
			return m, m.sendDebugCommand("step-out")
		}
		return m, nil

	case "q":
		// Quit debugger
		m.state = StateInactive
		return m, m.sendDebugCommand("disconnect")
	}

	return m, nil
}

func (m Model) View() string {
	if m.state == StateInactive {
		return ""
	}

	// Build the debugger layout
	header := m.renderHeader()
	content := m.renderContent()
	footer := m.renderFooter()

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		content,
		footer,
	)
}

func (m Model) renderHeader() string {
	t := m.theme

	// Status indicator
	var statusText string
	var statusColor lipgloss.AdaptiveColor

	switch m.state {
	case StateInitializing:
		statusText = "INITIALIZING"
		statusColor = t.Warning()
	case StateRunning:
		statusText = "RUNNING"
		statusColor = t.Success()
	case StatePaused:
		statusText = "PAUSED"
		statusColor = t.Error()
	case StateStopped:
		statusText = "STOPPED"
		statusColor = t.TextMuted()
	}

	statusStyle := styles.NewStyle().
		Foreground(statusColor).
		Bold(true).
		Padding(0, 1)

	// Program info
	programStyle := styles.NewStyle().
		Foreground(t.Text()).
		Padding(0, 1)

	// Build header
	headerContent := lipgloss.JoinHorizontal(lipgloss.Top,
		statusStyle.Render(fmt.Sprintf("🐛 %s", statusText)),
		programStyle.Render(fmt.Sprintf("│ %s", m.program)),
	)

	if m.currentFile != "" && m.currentLine > 0 {
		locationStyle := styles.NewStyle().
			Foreground(t.Primary()).
			Padding(0, 1)
		headerContent += locationStyle.Render(fmt.Sprintf("│ %s:%d", m.currentFile, m.currentLine))
	}

	return styles.NewStyle().
		Width(m.width).
		Background(t.BackgroundPanel()).
		Foreground(t.Text()).
		Render(headerContent)
}

func (m Model) renderContent() string {
	leftWidth := m.width / 2
	rightWidth := m.width - leftWidth

	// Left panel: Source code
	leftPanel := m.renderLeftPanel(leftWidth, m.height-4)

	// Right panel: Variables and Call Stack
	rightPanel := m.renderRightPanel(rightWidth, m.height-4)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)
}

func (m Model) renderLeftPanel(width, height int) string {
	t := m.theme

	title := "SOURCE CODE"
	if m.activePanel == 0 {
		title = "► " + title
	}

	titleStyle := styles.NewStyle().
		Foreground(t.Primary()).
		Bold(true).
		Padding(0, 1)

	border := styles.NewStyle().
		Width(width).
		Height(height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Border()).
		Padding(1)

	content := titleStyle.Render(title) + "\n\n"

	if m.state == StatePaused {
		// Render source code with current line highlighted
		content += m.sourceView.View()
	} else if m.state == StateRunning {
		content += styles.NewStyle().
			Foreground(t.TextMuted()).
			Render("Program is running...")
	} else {
		content += styles.NewStyle().
			Foreground(t.TextMuted()).
			Render("No source to display")
	}

	return border.Render(content)
}

func (m Model) renderRightPanel(width, height int) string {
	topHeight := height / 2
	bottomHeight := height - topHeight

	// Top: Variables
	variablesPanel := m.renderVariablesPanel(width, topHeight)

	// Bottom: Call Stack
	callStackPanel := m.renderCallStackPanel(width, bottomHeight)

	return lipgloss.JoinVertical(lipgloss.Left, variablesPanel, callStackPanel)
}

func (m Model) renderVariablesPanel(width, height int) string {
	t := m.theme

	title := "VARIABLES"
	if m.activePanel == 1 {
		title = "► " + title
	}

	titleStyle := styles.NewStyle().
		Foreground(t.Primary()).
		Bold(true).
		Padding(0, 1)

	border := styles.NewStyle().
		Width(width).
		Height(height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Border()).
		Padding(1)

	content := titleStyle.Render(title) + "\n\n"

	if m.state == StatePaused {
		content += m.variablesView.View()
	} else {
		content += styles.NewStyle().
			Foreground(t.TextMuted()).
			Render("No variables to display")
	}

	return border.Render(content)
}

func (m Model) renderCallStackPanel(width, height int) string {
	t := m.theme

	title := "CALL STACK"
	if m.activePanel == 2 {
		title = "► " + title
	}

	titleStyle := styles.NewStyle().
		Foreground(t.Primary()).
		Bold(true).
		Padding(0, 1)

	border := styles.NewStyle().
		Width(width).
		Height(height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Border()).
		Padding(1)

	content := titleStyle.Render(title) + "\n\n"

	if m.state == StatePaused {
		content += m.callStackView.View()
	} else {
		content += styles.NewStyle().
			Foreground(t.TextMuted()).
			Render("No call stack to display")
	}

	return border.Render(content)
}

func (m Model) renderFooter() string {
	t := m.theme

	if m.state != StatePaused {
		return ""
	}

	shortcuts := []string{
		"[c]ontinue",
		"[s]tep over",
		"[i]nto",
		"[o]ut",
		"[tab] switch panel",
		"[q]uit",
	}

	footerStyle := styles.NewStyle().
		Width(m.width).
		Background(t.BackgroundPanel()).
		Foreground(t.TextMuted()).
		Padding(0, 1)

	return footerStyle.Render(strings.Join(shortcuts, " • "))
}

func (m Model) updateLayout() {
	leftWidth := m.width / 2
	m.sourceView = m.sourceView.UpdateSize(leftWidth, m.height-4)
	m.variablesView = m.variablesView.UpdateSize(leftWidth, m.height/2-2)
	m.callStackView = m.callStackView.UpdateSize(leftWidth, m.height/2-2)
}

// GetState returns the current debugger state
func (m Model) GetState() DebuggerState {
	return m.state
}

// IsActive returns whether the debugger is currently active
func (m Model) IsActive() bool {
	return m.state != StateInactive
}

// sendDebugCommand sends a debug control command to the backend
func (m Model) sendDebugCommand(command string) tea.Cmd {
	return func() tea.Msg {
		if m.client == nil {
			return nil
		}

		path := fmt.Sprintf("/debug/%s/%s", m.sessionID, command)
		err := m.client.Post(context.Background(), path, map[string]interface{}{}, nil)
		if err != nil {
			// Log error but don't crash the TUI
			return nil
		}
		return nil
	}
}

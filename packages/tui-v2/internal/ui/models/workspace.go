package models

import (
	"strings"

	"github.com/aaronmrosenthal/rycode/packages/tui-v2/internal/layout"
	"github.com/aaronmrosenthal/rycode/packages/tui-v2/internal/theme"
	"github.com/aaronmrosenthal/rycode/packages/tui-v2/internal/ui/components"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// FocusPane represents which pane is currently focused
type FocusPane int

const (
	FocusFileTree FocusPane = iota
	FocusChat
)

// WorkspaceModel represents the full IDE workspace with FileTree and Chat
type WorkspaceModel struct {
	fileTree    *components.FileTree
	chat        ChatModel
	focus       FocusPane
	width       int
	height      int
	layoutMgr   *layout.LayoutManager
	ready       bool
	fileTreeWidth int // Width of file tree pane
}

// NewWorkspaceModel creates a new workspace model
func NewWorkspaceModel(rootPath string) WorkspaceModel {
	return WorkspaceModel{
		fileTree:      nil, // Will be initialized after first WindowSizeMsg
		chat:          NewChatModel(),
		focus:         FocusChat, // Start with chat focused
		width:         80,
		height:        24,
		layoutMgr:     layout.NewLayoutManager(80, 24),
		ready:         false,
		fileTreeWidth: 30, // Default file tree width
	}
}

// Init initializes the workspace model
func (m WorkspaceModel) Init() tea.Cmd {
	return m.chat.Init()
}

// Update handles messages and updates the model
func (m WorkspaceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.layoutMgr.Update(msg.Width, msg.Height)
		m.updateDimensions()
		m.ready = true

		// Also update chat
		chatUpdated, cmd := m.chat.Update(msg)
		m.chat = chatUpdated.(ChatModel)
		return m, cmd

	case StreamChunkMsg, StreamCompleteMsg:
		// Forward to chat
		chatUpdated, cmd := m.chat.Update(msg)
		m.chat = chatUpdated.(ChatModel)
		return m, cmd
	}

	return m, nil
}

// handleKeyPress handles keyboard input
func (m WorkspaceModel) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global shortcuts
	switch msg.String() {
	case "ctrl+c", "esc":
		return m, tea.Quit

	case "ctrl+b":
		// Toggle focus between FileTree and Chat
		if m.focus == FocusFileTree {
			m.focus = FocusChat
		} else {
			m.focus = FocusFileTree
		}
		return m, nil

	case "ctrl+t":
		// Toggle FileTree visibility
		if m.fileTreeWidth > 0 {
			m.fileTreeWidth = 0
		} else {
			m.fileTreeWidth = 30
		}
		m.updateDimensions()
		return m, nil
	}

	// Route keys based on focus
	if m.focus == FocusFileTree && m.fileTree != nil {
		return m.handleFileTreeKeys(msg)
	} else {
		// Forward to chat
		chatUpdated, cmd := m.chat.Update(msg)
		m.chat = chatUpdated.(ChatModel)
		return m, cmd
	}
}

// handleFileTreeKeys handles FileTree-specific keyboard shortcuts
func (m WorkspaceModel) handleFileTreeKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "j", "down":
		m.fileTree.SelectNext()
	case "k", "up":
		m.fileTree.SelectPrev()
	case "g":
		m.fileTree.SelectFirst()
	case "G":
		m.fileTree.SelectLast()
	case "h", "left", "backspace":
		m.fileTree.GoToParent()
	case "l", "right", "enter":
		m.fileTree.ToggleExpanded()
	case ".":
		m.fileTree.ToggleHidden()
	case "r":
		m.fileTree.Refresh()
	case "o":
		// Open selected file in chat
		selected := m.fileTree.GetSelected()
		if selected != nil && !selected.IsDir {
			// TODO: Send file path to chat
			// For now, just switch focus to chat
			m.focus = FocusChat
		}
	}
	return m, nil
}

// updateDimensions updates component dimensions based on layout
func (m *WorkspaceModel) updateDimensions() {
	// Initialize FileTree if needed
	if m.fileTree == nil {
		// Get current working directory for now
		// TODO: Make this configurable
		m.fileTree = components.NewFileTree(".", m.fileTreeWidth, m.height-2)
	}

	// Calculate dimensions based on device class
	deviceClass := m.layoutMgr.GetDeviceClass()

	// On mobile, hide file tree by default
	if deviceClass.IsMobile() && m.fileTreeWidth > 0 {
		m.fileTreeWidth = 0
	}

	// Update FileTree dimensions
	if m.fileTreeWidth > 0 {
		m.fileTree.SetWidth(m.fileTreeWidth)
		m.fileTree.SetHeight(m.height - 2) // Account for header
	}

	// Update Chat dimensions
	chatWidth := m.width
	if m.fileTreeWidth > 0 {
		chatWidth = m.width - m.fileTreeWidth - 1 // -1 for separator
	}

	m.chat.width = chatWidth
	m.chat.height = m.height
	m.chat.updateDimensions()
}

// View renders the workspace
func (m WorkspaceModel) View() string {
	if !m.ready {
		return "Initializing workspace..."
	}

	// Render header
	header := m.renderHeader()

	// Render main content
	var mainContent string
	if m.fileTreeWidth > 0 && m.fileTree != nil {
		// Split view: FileTree + Chat
		fileTreeView := m.fileTree.Render()
		chatView := m.chat.View()

		// Add focus indicators
		if m.focus == FocusFileTree {
			fileTreeView = m.addFocusIndicator(fileTreeView, true)
			chatView = m.addFocusIndicator(chatView, false)
		} else {
			fileTreeView = m.addFocusIndicator(fileTreeView, false)
			chatView = m.addFocusIndicator(chatView, true)
		}

		// Combine side by side
		mainContent = lipgloss.JoinHorizontal(
			lipgloss.Top,
			fileTreeView,
			chatView,
		)
	} else {
		// Chat only
		mainContent = m.chat.View()
	}

	// Render status bar
	statusBar := m.renderStatusBar()

	// Combine all parts
	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		mainContent,
		statusBar,
	)
}

// renderHeader renders the workspace header
func (m WorkspaceModel) renderHeader() string {
	deviceClass := m.layoutMgr.GetDeviceClass()

	title := theme.GradientTextPreset("RyCode Workspace", theme.GradientMatrix)
	subtitle := theme.MatrixTheme.Subtitle.Render(
		deviceClass.String() + " • " +
		m.layoutMgr.GetDeviceClass().String(),
	)

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

// renderStatusBar renders the workspace status bar
func (m WorkspaceModel) renderStatusBar() string {
	hints := []string{}

	if m.fileTreeWidth > 0 {
		hints = append(hints, "Ctrl+B: Switch pane")
	}
	hints = append(hints, "Ctrl+T: Toggle tree")

	if m.focus == FocusFileTree {
		hints = append(hints, "j/k: Navigate")
		hints = append(hints, "Enter: Expand")
		hints = append(hints, "o: Open")
	} else {
		hints = append(hints, "Enter: Send")
		hints = append(hints, "Tab: Accept")
	}

	hintsText := strings.Join(hints, " • ")

	// File info
	var fileInfo string
	if m.fileTree != nil {
		selected := m.fileTree.GetSelected()
		if selected != nil {
			fileInfo = selected.Name
		}
	}

	statusStyle := lipgloss.NewStyle().
		Foreground(theme.MatrixGreenDim).
		Padding(0, 1)

	leftStatus := statusStyle.Render(hintsText)
	rightStatus := statusStyle.Render(fileInfo)

	// Calculate spacing
	totalWidth := lipgloss.Width(leftStatus) + lipgloss.Width(rightStatus)
	spacing := m.width - totalWidth
	if spacing < 0 {
		spacing = 0
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftStatus,
		strings.Repeat(" ", spacing),
		rightStatus,
	)
}

// addFocusIndicator adds a visual indicator for focused/unfocused panes
func (m WorkspaceModel) addFocusIndicator(content string, focused bool) string {
	if focused {
		// Add bright border for focused pane
		style := lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(theme.MatrixGreen)
		return style.Render(content)
	}
	// Dim border for unfocused pane
	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(theme.MatrixGreenDark)
	return style.Render(content)
}

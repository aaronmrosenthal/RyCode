package components

import (
	"strings"

	"github.com/aaronmrosenthal/rycode/packages/tui-v2/internal/theme"
	"github.com/charmbracelet/lipgloss"
)

// InputBar represents the message input field
type InputBar struct {
	Value          string
	Cursor         int
	Placeholder    string
	MaxLines       int
	Width          int
	GhostText      string
	ShowVoiceButton bool
	ShowActions    bool
	Focused        bool
	Theme          theme.Theme
}

// NewInputBar creates a new input bar
func NewInputBar(width int) InputBar {
	return InputBar{
		Value:          "",
		Cursor:         0,
		Placeholder:    "Type a message or press 🎤 to speak...",
		MaxLines:       10,
		Width:          width,
		GhostText:      "",
		ShowVoiceButton: true,
		ShowActions:    true,
		Focused:        false,
		Theme:          theme.MatrixTheme,
	}
}

// Render renders the input bar
func (ib InputBar) Render() string {
	// Main input area
	inputContent := ib.renderInput()

	// Action buttons
	buttons := ib.renderButtons()

	// Quick actions (if enabled)
	var quickActions string
	if ib.ShowActions {
		quickActions = ib.renderQuickActions()
	}

	// Compose all parts
	var parts []string
	parts = append(parts, inputContent)
	if buttons != "" {
		parts = append(parts, buttons)
	}
	if quickActions != "" {
		parts = append(parts, quickActions)
	}

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

// renderInput renders the main input field
func (ib InputBar) renderInput() string {
	displayText := ib.Value

	// Show placeholder if empty
	if displayText == "" && !ib.Focused {
		placeholderStyle := lipgloss.NewStyle().
			Foreground(theme.MatrixGreenDark).
			Italic(true)
		displayText = placeholderStyle.Render(ib.Placeholder)
	}

	// Add ghost text if present and focused
	if ib.GhostText != "" && ib.Focused {
		ghostStyle := lipgloss.NewStyle().
			Foreground(theme.MatrixGreenDim)
		displayText += ghostStyle.Render(ib.GhostText)
	}

	// Add cursor if focused
	if ib.Focused && ib.Value != "" {
		// Insert cursor at position
		if ib.Cursor <= len(ib.Value) {
			before := ib.Value[:ib.Cursor]
			after := ""
			if ib.Cursor < len(ib.Value) {
				after = ib.Value[ib.Cursor:]
			}

			cursorStyle := lipgloss.NewStyle().
				Background(theme.MatrixGreen).
				Foreground(theme.Black)

			cursorChar := "█"
			if ib.Cursor < len(ib.Value) {
				cursorChar = string(ib.Value[ib.Cursor])
			}

			displayText = before + cursorStyle.Render(cursorChar) + after

			// Add ghost text after cursor
			if ib.GhostText != "" {
				ghostStyle := lipgloss.NewStyle().Foreground(theme.MatrixGreenDim)
				displayText += ghostStyle.Render(ib.GhostText)
			}
		}
	}

	// Apply input style
	var style lipgloss.Style
	if ib.Focused {
		style = ib.Theme.InputFocused.Width(ib.Width - 4)
	} else {
		style = ib.Theme.Input.Width(ib.Width - 4)
	}

	// Limit height to MaxLines
	lines := strings.Split(displayText, "\n")
	if len(lines) > ib.MaxLines {
		lines = lines[len(lines)-ib.MaxLines:]
		displayText = strings.Join(lines, "\n")
	}

	return style.Render(displayText)
}

// renderButtons renders action buttons
func (ib InputBar) renderButtons() string {
	var buttons []string

	// Voice button
	if ib.ShowVoiceButton {
		voiceStyle := lipgloss.NewStyle().
			Foreground(theme.NeonCyan)
		buttons = append(buttons, voiceStyle.Render("🎤 Voice"))
	}

	// Send button
	sendStyle := ib.Theme.Button
	if ib.Focused && ib.Value != "" {
		sendStyle = ib.Theme.ButtonActive
	}
	buttons = append(buttons, sendStyle.Render("Send ↵"))

	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		strings.Join(buttons, "  "),
	)
}

// renderQuickActions renders quick action buttons
func (ib InputBar) renderQuickActions() string {
	actions := []string{
		"Fix",
		"Test",
		"Explain",
		"Refactor",
		"Run",
	}

	buttonStyle := lipgloss.NewStyle().
		Foreground(theme.MatrixGreenDim).
		Padding(0, 1)

	var renderedActions []string
	for _, action := range actions {
		renderedActions = append(renderedActions, buttonStyle.Render(action))
	}

	labelStyle := lipgloss.NewStyle().
		Foreground(theme.MatrixGreenDark)

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		labelStyle.Render("Quick: "),
		strings.Join(renderedActions, " │ "),
	)
}

// SetValue sets the input value
func (ib *InputBar) SetValue(value string) {
	ib.Value = value
	ib.Cursor = len(value)
}

// InsertRune inserts a character at cursor position
func (ib *InputBar) InsertRune(r rune) {
	before := ib.Value[:ib.Cursor]
	after := ib.Value[ib.Cursor:]
	ib.Value = before + string(r) + after
	ib.Cursor++
}

// DeleteCharBefore deletes the character before the cursor (backspace)
func (ib *InputBar) DeleteCharBefore() {
	if ib.Cursor > 0 {
		before := ib.Value[:ib.Cursor-1]
		after := ib.Value[ib.Cursor:]
		ib.Value = before + after
		ib.Cursor--
	}
}

// DeleteCharAfter deletes the character after the cursor (delete key)
func (ib *InputBar) DeleteCharAfter() {
	if ib.Cursor < len(ib.Value) {
		before := ib.Value[:ib.Cursor]
		after := ib.Value[ib.Cursor+1:]
		ib.Value = before + after
	}
}

// MoveCursorLeft moves cursor left
func (ib *InputBar) MoveCursorLeft() {
	if ib.Cursor > 0 {
		ib.Cursor--
	}
}

// MoveCursorRight moves cursor right
func (ib *InputBar) MoveCursorRight() {
	if ib.Cursor < len(ib.Value) {
		ib.Cursor++
	}
}

// MoveCursorToStart moves cursor to start
func (ib *InputBar) MoveCursorToStart() {
	ib.Cursor = 0
}

// MoveCursorToEnd moves cursor to end
func (ib *InputBar) MoveCursorToEnd() {
	ib.Cursor = len(ib.Value)
}

// Clear clears the input
func (ib *InputBar) Clear() {
	ib.Value = ""
	ib.Cursor = 0
	ib.GhostText = ""
}

// SetGhostText sets the ghost text suggestion
func (ib *InputBar) SetGhostText(text string) {
	ib.GhostText = text
}

// AcceptGhostText accepts the ghost text suggestion
func (ib *InputBar) AcceptGhostText() {
	if ib.GhostText != "" {
		ib.Value += ib.GhostText
		ib.Cursor = len(ib.Value)
		ib.GhostText = ""
	}
}

// GetValue returns the current input value
func (ib InputBar) GetValue() string {
	return ib.Value
}

// IsEmpty returns true if input is empty
func (ib InputBar) IsEmpty() bool {
	return strings.TrimSpace(ib.Value) == ""
}

// SetFocus sets the focus state
func (ib *InputBar) SetFocus(focused bool) {
	ib.Focused = focused
}

// SetWidth sets the width
func (ib *InputBar) SetWidth(width int) {
	ib.Width = width
}

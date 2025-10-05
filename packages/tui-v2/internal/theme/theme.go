package theme

import "github.com/charmbracelet/lipgloss"

// Theme represents a complete UI theme
type Theme struct {
	Name string

	// Text styles
	Primary   lipgloss.Style
	Secondary lipgloss.Style
	Dim       lipgloss.Style
	Error     lipgloss.Style
	Success   lipgloss.Style
	Warning   lipgloss.Style
	Info      lipgloss.Style

	// UI element styles
	Border    lipgloss.Style
	Highlight lipgloss.Style
	Selected  lipgloss.Style

	// Component styles
	Button       lipgloss.Style
	ButtonActive lipgloss.Style
	Input        lipgloss.Style
	InputFocused lipgloss.Style
	CodeBlock    lipgloss.Style
	Quote        lipgloss.Style
	Link         lipgloss.Style

	// Message styles
	MessageUser lipgloss.Style
	MessageAI   lipgloss.Style

	// Status styles
	StatusBar lipgloss.Style
	Title     lipgloss.Style
	Subtitle  lipgloss.Style
	Hint      lipgloss.Style

	// Special effects
	Glow     lipgloss.Style
	Gradient lipgloss.Style
}

// MatrixTheme is the default Matrix-themed style
var MatrixTheme = Theme{
	Name: "Matrix",

	// Text styles
	Primary: lipgloss.NewStyle().
		Foreground(MatrixGreen),

	Secondary: lipgloss.NewStyle().
		Foreground(MatrixGreenDim),

	Dim: lipgloss.NewStyle().
		Foreground(MatrixGreenDark),

	Error: lipgloss.NewStyle().
		Foreground(ColorError).
		Bold(true),

	Success: lipgloss.NewStyle().
		Foreground(ColorSuccess).
		Bold(true),

	Warning: lipgloss.NewStyle().
		Foreground(ColorWarning),

	Info: lipgloss.NewStyle().
		Foreground(ColorInfo),

	// UI elements
	Border: lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(MatrixGreen),

	Highlight: lipgloss.NewStyle().
		Background(DarkGreen),

	Selected: lipgloss.NewStyle().
		Background(DarkGreen).
		Foreground(MatrixGreen).
		Bold(true),

	// Components
	Button: lipgloss.NewStyle().
		Foreground(Black).
		Background(MatrixGreen).
		Padding(0, 2).
		Bold(true),

	ButtonActive: lipgloss.NewStyle().
		Foreground(Black).
		Background(MatrixGreenBright).
		Padding(0, 2).
		Bold(true),

	Input: lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(MatrixGreenDim).
		Padding(0, 1),

	InputFocused: lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(MatrixGreen).
		Padding(0, 1),

	CodeBlock: lipgloss.NewStyle().
		Background(DarkerGreen).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(MatrixGreenDark).
		Padding(1).
		MarginTop(1).
		MarginBottom(1),

	Quote: lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderLeft(true).
		BorderForeground(NeonCyan).
		PaddingLeft(2).
		Foreground(MatrixGreenDim).
		Italic(true),

	Link: lipgloss.NewStyle().
		Foreground(NeonCyan).
		Underline(true),

	// Messages
	MessageUser: lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(NeonCyan).
		Padding(1).
		MarginBottom(1),

	MessageAI: lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(MatrixGreen).
		Padding(1).
		MarginBottom(1),

	// Status
	StatusBar: lipgloss.NewStyle().
		Background(DarkGreen).
		Foreground(MatrixGreen).
		Padding(0, 1),

	Title: lipgloss.NewStyle().
		Foreground(MatrixGreen).
		Bold(true).
		MarginTop(1).
		MarginBottom(1),

	Subtitle: lipgloss.NewStyle().
		Foreground(MatrixGreenDim),

	Hint: lipgloss.NewStyle().
		Foreground(MatrixGreenDark).
		Italic(true),

	// Effects
	Glow: lipgloss.NewStyle().
		Foreground(MatrixGreen).
		Bold(true),

	Gradient: lipgloss.NewStyle().
		Foreground(MatrixGreen),
}

// Helper methods for Theme

// RenderTitle renders a title with the theme
func (t Theme) RenderTitle(text string) string {
	return t.Title.Render(text)
}

// RenderSubtitle renders a subtitle with the theme
func (t Theme) RenderSubtitle(text string) string {
	return t.Subtitle.Render(text)
}

// RenderHint renders a hint with the theme
func (t Theme) RenderHint(text string) string {
	return t.Hint.Render(text)
}

// RenderError renders an error message with the theme
func (t Theme) RenderError(text string) string {
	return t.Error.Render("✗ " + text)
}

// RenderSuccess renders a success message with the theme
func (t Theme) RenderSuccess(text string) string {
	return t.Success.Render("✓ " + text)
}

// RenderWarning renders a warning message with the theme
func (t Theme) RenderWarning(text string) string {
	return t.Warning.Render("⚠ " + text)
}

// RenderInfo renders an info message with the theme
func (t Theme) RenderInfo(text string) string {
	return t.Info.Render("ℹ " + text)
}

// RenderButton renders a button with the theme
func (t Theme) RenderButton(text string, active bool) string {
	if active {
		return t.ButtonActive.Render(text)
	}
	return t.Button.Render(text)
}

// RenderInput renders an input field with the theme
func (t Theme) RenderInput(text string, focused bool) string {
	if focused {
		return t.InputFocused.Render(text)
	}
	return t.Input.Render(text)
}

// RenderCodeBlock renders a code block with the theme
func (t Theme) RenderCodeBlock(code string) string {
	return t.CodeBlock.Render(code)
}

// RenderQuote renders a quote with the theme
func (t Theme) RenderQuote(text string) string {
	return t.Quote.Render(text)
}

// RenderLink renders a link with the theme
func (t Theme) RenderLink(text string) string {
	return t.Link.Render(text)
}

// RenderMessage renders a message with appropriate styling
func (t Theme) RenderMessage(text string, isUser bool) string {
	if isUser {
		return t.MessageUser.Render(text)
	}
	return t.MessageAI.Render(text)
}

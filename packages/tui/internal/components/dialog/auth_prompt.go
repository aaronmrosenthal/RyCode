package dialog

import (
	"fmt"

	"github.com/aaronmrosenthal/rycode/internal/styles"
	"github.com/aaronmrosenthal/rycode/internal/theme"
	"github.com/charmbracelet/bubbles/v2/textinput"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

// AuthPromptDialog is a dialog for entering provider API keys
type AuthPromptDialog struct {
	provider       string
	input          textinput.Model
	error          string
	showAutoDetect bool
	width          int
	height         int
}

// NewAuthPromptDialog creates a new authentication prompt dialog
func NewAuthPromptDialog(provider string) *AuthPromptDialog {
	ti := textinput.New()
	ti.Placeholder = "Enter API key..."
	ti.Focus()
	ti.CharLimit = 256
	ti.SetWidth(60)
	ti.EchoMode = textinput.EchoPassword // Hide API key

	return &AuthPromptDialog{
		provider:       provider,
		input:          ti,
		showAutoDetect: true,
		error:          "",
	}
}

// SetSize sets the dialog dimensions
func (a *AuthPromptDialog) SetSize(width, height int) {
	a.width = width
	a.height = height

	// Adjust input width based on dialog width
	if width > 80 {
		a.input.SetWidth(60)
	} else if width > 60 {
		a.input.SetWidth(50)
	} else {
		a.input.SetWidth(40)
	}
}

// SetError sets an error message to display
func (a *AuthPromptDialog) SetError(err string) {
	a.error = err
}

// GetValue returns the current input value
func (a *AuthPromptDialog) GetValue() string {
	return a.input.Value()
}

// Update handles messages for the auth prompt
func (a *AuthPromptDialog) Update(msg tea.Msg) (*AuthPromptDialog, tea.Cmd) {
	var cmd tea.Cmd
	a.input, cmd = a.input.Update(msg)
	return a, cmd
}

// View renders the auth prompt dialog
func (a *AuthPromptDialog) View() string {
	t := theme.CurrentTheme()

	titleStyle := styles.NewStyle().
		Foreground(t.Primary()).
		Bold(true).
		MarginBottom(1)

	inputStyle := styles.NewStyle().
		Foreground(t.Text()).
		Background(t.BackgroundPanel()).
		MarginBottom(1)

	hintStyle := styles.NewStyle().
		Foreground(t.TextMuted()).
		Faint(true).
		MarginBottom(1)

	errorStyle := styles.NewStyle().
		Foreground(t.Error()).
		Bold(true)

	// Title
	title := titleStyle.Render(fmt.Sprintf("Authenticate with %s", a.provider))

	// Input field
	inputView := inputStyle.Render(a.input.View())

	// Hints
	hints := ""
	if a.showAutoDetect {
		hints = hintStyle.Render("Press Enter to submit | Ctrl+D for auto-detect | Esc to cancel")
	} else {
		hints = hintStyle.Render("Press Enter to submit | Esc to cancel")
	}

	// Error message
	errorView := ""
	if a.error != "" {
		errorView = "\n" + errorStyle.Render("✗ " + a.error)
	}

	// Combine all parts
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		inputView,
		"",
		hints,
		errorView,
	)

	// Center the content
	dialogStyle := styles.NewStyle().
		Padding(2, 4).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Primary()).
		Background(t.BackgroundPanel())

	return dialogStyle.Render(content)
}

// Message types for auth prompt

// AuthSubmitMsg is sent when user submits API key
type AuthSubmitMsg struct {
	Provider string
	APIKey   string
}

// AuthCancelMsg is sent when user cancels auth
type AuthCancelMsg struct{}

// AuthAutoDetectMsg is sent when user requests auto-detect
type AuthAutoDetectMsg struct {
	Provider string
}

// AuthSuccessMsg is sent when authentication succeeds
type AuthSuccessMsg struct {
	Provider    string
	ModelsCount int
}

// AuthFailureMsg is sent when authentication fails
type AuthFailureMsg struct {
	Provider string
	Error    string
}

// AuthStatusRefreshMsg is sent to refresh auth status
type AuthStatusRefreshMsg struct{}

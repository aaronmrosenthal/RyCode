package dialog

import (
	"fmt"
	"strings"

	"github.com/aaronmrosenthal/rycode/internal/styles"
	"github.com/aaronmrosenthal/rycode/internal/theme"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

// ErrorType represents different types of errors
type ErrorType int

const (
	ErrorTypeGeneric ErrorType = iota
	ErrorTypeAuth
	ErrorTypeNetwork
	ErrorTypeRateLimit
	ErrorTypePermission
	ErrorTypeValidation
)

// ErrorAction represents an action the user can take
type ErrorAction struct {
	Key         string
	Label       string
	Description string
	Action      func() tea.Cmd
}

// ErrorDialog displays enhanced error messages with helpful guidance
type ErrorDialog struct {
	errorType   ErrorType
	title       string
	message     string
	details     string
	suggestions []string
	actions     []ErrorAction
	docsLink    string
	width       int
	height      int
	focused     int
}

// NewErrorDialog creates a new enhanced error dialog
func NewErrorDialog(errorType ErrorType, title, message string) *ErrorDialog {
	return &ErrorDialog{
		errorType:   errorType,
		title:       title,
		message:     message,
		suggestions: make([]string, 0),
		actions:     make([]ErrorAction, 0),
		focused:     0,
	}
}

// NewAuthError creates an authentication error dialog
func NewAuthError(provider, errorMsg string) *ErrorDialog {
	dialog := NewErrorDialog(
		ErrorTypeAuth,
		fmt.Sprintf("Authentication Failed: %s", provider),
		errorMsg,
	)

	// Add provider-specific suggestions
	switch strings.ToLower(provider) {
	case "anthropic", "claude":
		dialog.WithSuggestions(
			"Check your API key for typos",
			"Verify your key at console.anthropic.com",
			"Ensure your API key has the correct permissions",
			"Try generating a new API key",
		).WithDocsLink("https://docs.anthropic.com/claude/reference/getting-started-with-the-api")

	case "openai", "gpt":
		dialog.WithSuggestions(
			"Check your API key format (should start with 'sk-')",
			"Verify your key at platform.openai.com/api-keys",
			"Ensure you have available credits",
			"Check if your organization ID is correct",
		).WithDocsLink("https://platform.openai.com/docs/quickstart")

	case "google", "gemini":
		dialog.WithSuggestions(
			"Verify your API key at console.cloud.google.com",
			"Ensure the Generative AI API is enabled",
			"Check your project quota limits",
			"Confirm your billing is active",
		).WithDocsLink("https://ai.google.dev/docs")

	case "grok", "x.ai":
		dialog.WithSuggestions(
			"Verify your API key at x.ai/api",
			"Ensure you have access to the Grok API",
			"Check your usage limits",
		).WithDocsLink("https://x.ai/api")

	case "qwen", "alibaba":
		dialog.WithSuggestions(
			"Verify your API key at dashscope.aliyun.com",
			"Ensure your account has sufficient balance",
			"Check regional availability",
		).WithDocsLink("https://help.aliyun.com/zh/dashscope")

	default:
		dialog.WithSuggestions(
			"Check your API key for typos",
			"Verify the key is valid and active",
			"Ensure you have the necessary permissions",
			"Try generating a new API key",
		)
	}

	// Add common actions
	dialog.WithActions(
		ErrorAction{
			Key:         "r",
			Label:       "Retry",
			Description: "Try authenticating again",
		},
		ErrorAction{
			Key:         "d",
			Label:       "Docs",
			Description: "Open documentation",
		},
	)

	return dialog
}

// NewNetworkError creates a network error dialog
func NewNetworkError(operation, errorMsg string) *ErrorDialog {
	dialog := NewErrorDialog(
		ErrorTypeNetwork,
		fmt.Sprintf("Network Error: %s", operation),
		errorMsg,
	)

	dialog.WithSuggestions(
		"Check your internet connection",
		"Verify the service is not down (check status page)",
		"Try again in a few moments",
		"Check if you need to configure a proxy",
	)

	dialog.WithActions(
		ErrorAction{
			Key:         "r",
			Label:       "Retry",
			Description: "Try the operation again",
		},
	)

	return dialog
}

// NewRateLimitError creates a rate limit error dialog
func NewRateLimitError(provider string, retryAfter string) *ErrorDialog {
	dialog := NewErrorDialog(
		ErrorTypeRateLimit,
		"Rate Limit Exceeded",
		fmt.Sprintf("You've exceeded the rate limit for %s", provider),
	)

	suggestions := []string{
		"Wait a moment before trying again",
		"Consider upgrading your API tier for higher limits",
		"Implement exponential backoff in your requests",
	}

	if retryAfter != "" {
		suggestions = append([]string{
			fmt.Sprintf("Retry after: %s", retryAfter),
		}, suggestions...)
	}

	dialog.WithSuggestions(suggestions...)

	dialog.WithActions(
		ErrorAction{
			Key:         "w",
			Label:       "Wait & Retry",
			Description: "Wait and automatically retry",
		},
	)

	return dialog
}

// WithSuggestions adds helpful suggestions
func (e *ErrorDialog) WithSuggestions(suggestions ...string) *ErrorDialog {
	e.suggestions = suggestions
	return e
}

// WithActions adds action buttons
func (e *ErrorDialog) WithActions(actions ...ErrorAction) *ErrorDialog {
	e.actions = actions
	return e
}

// WithDetails adds additional error details
func (e *ErrorDialog) WithDetails(details string) *ErrorDialog {
	e.details = details
	return e
}

// WithDocsLink adds a documentation link
func (e *ErrorDialog) WithDocsLink(link string) *ErrorDialog {
	e.docsLink = link
	return e
}

// SetSize sets the dialog dimensions
func (e *ErrorDialog) SetSize(width, height int) {
	e.width = width
	e.height = height
}

// Init implements tea.Model
func (e *ErrorDialog) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (e *ErrorDialog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		// Check action keys
		for _, action := range e.actions {
			if msg.String() == action.Key {
				if action.Action != nil {
					return e, action.Action()
				}
				return e, nil
			}
		}

		// Navigation
		switch msg.String() {
		case "up", "k":
			if e.focused > 0 {
				e.focused--
			}
		case "down", "j":
			if e.focused < len(e.actions)-1 {
				e.focused++
			}
		case "enter":
			if e.focused < len(e.actions) && e.actions[e.focused].Action != nil {
				return e, e.actions[e.focused].Action()
			}
		}
	}

	return e, nil
}

// View implements tea.Model
func (e *ErrorDialog) View() string {
	t := theme.CurrentTheme()

	// Icon based on error type
	icon := "✗"
	iconColor := t.Error()

	switch e.errorType {
	case ErrorTypeAuth:
		icon = "🔒"
	case ErrorTypeNetwork:
		icon = "🌐"
	case ErrorTypeRateLimit:
		icon = "⏱"
	case ErrorTypePermission:
		icon = "🚫"
	case ErrorTypeValidation:
		icon = "⚠"
	}

	// Styles
	iconStyle := styles.NewStyle().
		Foreground(iconColor).
		Bold(true)

	titleStyle := styles.NewStyle().
		Foreground(iconColor).
		Bold(true).
		MarginBottom(1)

	messageStyle := styles.NewStyle().
		Foreground(t.Text()).
		Width(60).
		MarginBottom(1)

	detailsStyle := styles.NewStyle().
		Foreground(t.TextMuted()).
		Faint(true).
		Width(60).
		MarginBottom(1)

	sectionTitleStyle := styles.NewStyle().
		Foreground(t.Primary()).
		Bold(true).
		MarginTop(1).
		MarginBottom(1)

	suggestionStyle := styles.NewStyle().
		Foreground(t.Text()).
		Width(58).
		PaddingLeft(2)

	actionKeyStyle := styles.NewStyle().
		Foreground(t.Primary()).
		Bold(true).
		Background(t.BackgroundElement()).
		Padding(0, 1)

	actionLabelStyle := styles.NewStyle().
		Foreground(t.Text()).
		Bold(true)

	actionDescStyle := styles.NewStyle().
		Foreground(t.TextMuted()).
		Faint(true)

	docsStyle := styles.NewStyle().
		Foreground(t.TextMuted()).
		Faint(true).
		MarginTop(1)

	// Build content
	var content strings.Builder

	// Header: icon + title
	header := iconStyle.Render(icon) + "  " + titleStyle.Render(e.title)
	content.WriteString(header)
	content.WriteString("\n\n")

	// Message
	content.WriteString(messageStyle.Render(e.message))
	content.WriteString("\n")

	// Details (if any)
	if e.details != "" {
		content.WriteString(detailsStyle.Render(e.details))
		content.WriteString("\n")
	}

	// Suggestions
	if len(e.suggestions) > 0 {
		content.WriteString(sectionTitleStyle.Render("What to do:"))
		content.WriteString("\n")

		for i, suggestion := range e.suggestions {
			bullet := fmt.Sprintf("%d. ", i+1)
			line := bullet + suggestion
			content.WriteString(suggestionStyle.Render(line))
			content.WriteString("\n")
		}
	}

	// Actions
	if len(e.actions) > 0 {
		content.WriteString("\n")
		actionLine := ""

		for i, action := range e.actions {
			if i > 0 {
				actionLine += "    "
			}

			keyPart := actionKeyStyle.Render(action.Key)
			labelPart := actionLabelStyle.Render(action.Label)

			if i == e.focused {
				// Highlight focused action
				keyPart = styles.NewStyle().
					Foreground(t.Background()).
					Background(t.Primary()).
					Bold(true).
					Padding(0, 1).
					Render(action.Key)
			}

			actionLine += keyPart + " " + labelPart
		}

		content.WriteString(actionLine)
		content.WriteString("\n")

		// Action descriptions
		if e.focused < len(e.actions) {
			desc := e.actions[e.focused].Description
			content.WriteString("\n")
			content.WriteString(actionDescStyle.Render(desc))
		}
	}

	// Docs link
	if e.docsLink != "" {
		content.WriteString("\n")
		docsText := fmt.Sprintf("📚 Docs: %s", e.docsLink)
		content.WriteString(docsStyle.Render(docsText))
	}

	// Wrap in dialog box
	dialogStyle := styles.NewStyle().
		Padding(2, 4).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(iconColor).
		Background(t.BackgroundPanel()).
		Width(68)

	return dialogStyle.Render(content.String())
}

// ErrorDialogMsg is sent when the error dialog requests an action
type ErrorDialogMsg struct {
	Action string
}

// CloseErrorDialogMsg is sent when the error dialog should close
type CloseErrorDialogMsg struct{}

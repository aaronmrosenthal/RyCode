package help

import (
	"strings"

	"github.com/aaronmrosenthal/rycode/internal/styles"
	"github.com/aaronmrosenthal/rycode/internal/theme"
)

// EmptyState represents a beautiful empty state UI
type EmptyState struct {
	Icon        string
	Title       string
	Description string
	Actions     []EmptyStateAction
}

// EmptyStateAction represents an actionable step in empty state
type EmptyStateAction struct {
	Label    string
	Shortcut string
	Primary  bool
}

// RenderEmptyState creates a beautiful empty state view
func RenderEmptyState(state EmptyState, width, height int) string {
	t := theme.CurrentTheme()

	var lines []string

	// Add vertical padding
	verticalPadding := (height - 10) / 2
	if verticalPadding > 0 {
		for i := 0; i < verticalPadding; i++ {
			lines = append(lines, "")
		}
	}

	// Large icon
	iconStyle := styles.NewStyle().
		Foreground(t.TextMuted()).
		Faint(true)

	// Make icon bigger
	bigIcon := state.Icon + state.Icon + state.Icon
	icon := iconStyle.Render(bigIcon)
	lines = append(lines, centerText(icon, width))
	lines = append(lines, "")

	// Title
	titleStyle := styles.NewStyle().
		Foreground(t.Text()).
		Bold(true)

	title := titleStyle.Render(state.Title)
	lines = append(lines, centerText(title, width))
	lines = append(lines, "")

	// Description
	descStyle := styles.NewStyle().
		Foreground(t.TextMuted()).
		Faint(true).
		Width(60)

	desc := descStyle.Render(state.Description)
	lines = append(lines, centerText(desc, width))
	lines = append(lines, "")
	lines = append(lines, "")

	// Actions
	if len(state.Actions) > 0 {
		actionsLine := renderActions(state.Actions)
		lines = append(lines, centerText(actionsLine, width))
	}

	return strings.Join(lines, "\n")
}

// renderActions creates action buttons
func renderActions(actions []EmptyStateAction) string {
	t := theme.CurrentTheme()

	var actionParts []string

	for _, action := range actions {
		// Button style
		buttonStyle := styles.NewStyle().
			Foreground(t.Background()).
			Background(t.TextMuted()).
			Padding(0, 2)

		if action.Primary {
			buttonStyle = buttonStyle.Background(t.Primary())
		}

		button := buttonStyle.Render(action.Label)

		// Shortcut hint
		shortcut := ""
		if action.Shortcut != "" {
			shortcutStyle := styles.NewStyle().
				Foreground(t.TextMuted()).
				Faint(true)

			shortcut = " " + shortcutStyle.Render("["+action.Shortcut+"]")
		}

		actionParts = append(actionParts, button+shortcut)
	}

	return strings.Join(actionParts, "  ")
}

// centerText centers text within a given width
func centerText(text string, width int) string {
	// Simple centering - doesn't account for ANSI codes perfectly
	textWidth := len(text)
	if textWidth >= width {
		return text
	}

	padding := (width - textWidth) / 2
	return strings.Repeat(" ", padding) + text
}

// Common empty states

// GetNoProvidersEmptyState returns empty state for no authenticated providers
func GetNoProvidersEmptyState() EmptyState {
	return EmptyState{
		Icon:        "🔐",
		Title:       "No Providers Authenticated",
		Description: "You need to authenticate at least one AI provider to use RyCode.\nDon't worry - it only takes a moment!",
		Actions: []EmptyStateAction{
			{
				Label:    "Auto-Detect Credentials",
				Shortcut: "d",
				Primary:  true,
			},
			{
				Label:    "Authenticate Manually",
				Shortcut: "a",
				Primary:  false,
			},
		},
	}
}

// GetNoModelsEmptyState returns empty state for no available models
func GetNoModelsEmptyState() EmptyState {
	return EmptyState{
		Icon:        "🤖",
		Title:       "No Models Available",
		Description: "Authenticate a provider to see available AI models.",
		Actions: []EmptyStateAction{
			{
				Label:    "Open Provider Management",
				Shortcut: "Ctrl+P",
				Primary:  true,
			},
		},
	}
}

// GetNoDataEmptyState returns empty state for no usage data
func GetNoDataEmptyState() EmptyState {
	return EmptyState{
		Icon:        "📊",
		Title:       "No Usage Data Yet",
		Description: "Start using RyCode to see beautiful analytics, insights, and cost forecasts here.\nYour first API call will begin tracking.",
		Actions: []EmptyStateAction{
			{
				Label:    "Select a Model",
				Shortcut: "Ctrl+M",
				Primary:  true,
			},
		},
	}
}

// GetWelcomeEmptyState returns empty state for first-time users
// Uses provider-specific welcome message if a ProviderTheme is active
func GetWelcomeEmptyState() EmptyState {
	t := theme.CurrentTheme()

	// Default welcome message
	welcomeMsg := "Your AI-powered development assistant is ready.\nLet's get started with a quick setup."

	// Check if current theme is a provider theme with custom welcome message
	if providerTheme, ok := t.(*theme.ProviderTheme); ok {
		if providerTheme.WelcomeMessage != "" {
			welcomeMsg = providerTheme.WelcomeMessage
		}
	}

	return EmptyState{
		Icon:        "👋",
		Title:       "Welcome to RyCode!",
		Description: welcomeMsg,
		Actions: []EmptyStateAction{
			{
				Label:    "Start Welcome Guide",
				Shortcut: "Enter",
				Primary:  true,
			},
			{
				Label:    "Skip to Provider Setup",
				Shortcut: "Ctrl+P",
				Primary:  false,
			},
		},
	}
}

// GetChatEmptyState returns empty state for empty chat
func GetChatEmptyState() EmptyState {
	return EmptyState{
		Icon:        "💬",
		Title:       "Ready to Chat!",
		Description: "Select an AI model and start a conversation.\nRyCode will help you choose the best model for your needs.",
		Actions: []EmptyStateAction{
			{
				Label:    "Select Model",
				Shortcut: "Ctrl+M",
				Primary:  true,
			},
			{
				Label:    "View Shortcuts",
				Shortcut: "Ctrl+?",
				Primary:  false,
			},
		},
	}
}

// GetErrorEmptyState returns empty state for errors
func GetErrorEmptyState(errorMsg string) EmptyState {
	return EmptyState{
		Icon:        "⚠️",
		Title:       "Something Went Wrong",
		Description: errorMsg + "\n\nTry refreshing or check your provider authentication.",
		Actions: []EmptyStateAction{
			{
				Label:    "Refresh",
				Shortcut: "Ctrl+R",
				Primary:  true,
			},
			{
				Label:    "Provider Management",
				Shortcut: "Ctrl+P",
				Primary:  false,
			},
		},
	}
}

// GetLoadingState returns a loading state view
func GetLoadingState(message string) string {
	t := theme.CurrentTheme()

	spinner := styles.NewStyle().
		Foreground(t.Primary()).
		Bold(true).
		Render("⠋ ")

	text := styles.NewStyle().
		Foreground(t.Text()).
		Render(message)

	return spinner + text
}

// GetSuccessState returns a success message view
func GetSuccessState(message string) string {
	t := theme.CurrentTheme()

	icon := styles.NewStyle().
		Foreground(t.Success()).
		Bold(true).
		Render("✓ ")

	text := styles.NewStyle().
		Foreground(t.Text()).
		Render(message)

	return icon + text
}

// GetWarningState returns a warning message view
func GetWarningState(message string) string {
	t := theme.CurrentTheme()

	icon := styles.NewStyle().
		Foreground(t.Warning()).
		Bold(true).
		Render("⚠ ")

	text := styles.NewStyle().
		Foreground(t.Text()).
		Render(message)

	return icon + text
}

// GetInfoState returns an info message view
func GetInfoState(message string) string {
	t := theme.CurrentTheme()

	icon := styles.NewStyle().
		Foreground(t.Info()).
		Bold(true).
		Render("ℹ ")

	text := styles.NewStyle().
		Foreground(t.Text()).
		Render(message)

	return icon + text
}

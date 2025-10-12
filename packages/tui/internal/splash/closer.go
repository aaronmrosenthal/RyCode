package splash

import (
	"strings"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/lipgloss/v2/compat"
)

// Closer manages the final "closer" screen
type Closer struct {
	width  int
	height int
}

// NewCloser creates a new closer screen
func NewCloser(width, height int) *Closer {
	return &Closer{
		width:  width,
		height: height,
	}
}

// Render renders the closer screen with proper borders using lipgloss
func (c *Closer) Render() string {
	// Matrix green/cyan colors matching terminal theme
	brightCyan := compat.AdaptiveColor{
		Dark:  lipgloss.Color("#00FFAA"),
		Light: lipgloss.Color("#00CC88"),
	}

	mediumGreen := compat.AdaptiveColor{
		Dark:  lipgloss.Color("#00CC88"),
		Light: lipgloss.Color("#008866"),
	}

	// Title style
	titleStyle := lipgloss.NewStyle().
		Foreground(brightCyan).
		Bold(true).
		Align(lipgloss.Center)

	// Calculate responsive width (max 70% of screen, min 50 chars)
	contentWidth := c.width * 7 / 10
	if contentWidth > 70 {
		contentWidth = 70
	}
	if contentWidth < 50 {
		contentWidth = 50
	}

	// Message style - responsive width
	messageStyle := lipgloss.NewStyle().
		Foreground(mediumGreen).
		Align(lipgloss.Center).
		Width(contentWidth - 8) // Account for padding and borders

	// Prompt style
	promptStyle := lipgloss.NewStyle().
		Foreground(brightCyan).
		Align(lipgloss.Center).
		Italic(true)

	// Build content
	var content strings.Builder
	content.WriteString(titleStyle.Render("🌀 RYCODE NEURAL CORTEX ACTIVE 🌀"))
	content.WriteString("\n\n")
	content.WriteString(messageStyle.Render(`"Every LLM fused. Every edge case covered.
You're not just coding—you're orchestrating
intelligence."`))
	content.WriteString("\n\n")
	content.WriteString(promptStyle.Render("Press any key to begin..."))

	// Create box with proper lipgloss borders - tighter padding
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(brightCyan).
		Padding(1, 2).
		Align(lipgloss.Center).
		MaxWidth(contentWidth)

	box := boxStyle.Render(content.String())

	// Center the box vertically and horizontally
	return lipgloss.Place(
		c.width,
		c.height,
		lipgloss.Center,
		lipgloss.Center,
		box,
	)
}

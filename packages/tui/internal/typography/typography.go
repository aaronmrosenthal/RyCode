package typography

import (
	"fmt"

	"github.com/aaronmrosenthal/rycode/internal/styles"
	"github.com/aaronmrosenthal/rycode/internal/theme"
	"github.com/charmbracelet/lipgloss/v2"
)

// Typography defines the type system for the TUI
type Typography struct {
	// Display - Largest text (hero sections, welcome screens)
	Display styles.Style

	// Title - Section titles, dialog headers
	Title styles.Style

	// Heading - Subsection headers
	Heading styles.Style

	// Subheading - Tertiary headers
	Subheading styles.Style

	// Body - Default text
	Body styles.Style

	// BodyEmphasis - Emphasized body text
	BodyEmphasis styles.Style

	// Caption - Small supporting text
	Caption styles.Style

	// Label - Input labels, tags
	Label styles.Style

	// Code - Inline code, technical text
	Code styles.Style

	// Link - Hyperlinks, clickable text
	Link styles.Style

	// Error - Error messages
	Error styles.Style

	// Success - Success messages
	Success styles.Style

	// Warning - Warning messages
	Warning styles.Style

	// Info - Info messages
	Info styles.Style
}

// New creates a new typography system based on current theme
func New() *Typography {
	t := theme.CurrentTheme()

	return &Typography{
		// Display: 24px equivalent, extra bold, high contrast
		Display: styles.NewStyle().
			Foreground(t.Primary()).
			Bold(true).
			MarginBottom(2),

		// Title: 20px equivalent, bold, primary color
		Title: styles.NewStyle().
			Foreground(t.Primary()).
			Bold(true).
			MarginBottom(1),

		// Heading: 18px equivalent, bold, text color
		Heading: styles.NewStyle().
			Foreground(t.Text()).
			Bold(true).
			MarginTop(1).
			MarginBottom(1),

		// Subheading: 16px equivalent, semibold, text color
		Subheading: styles.NewStyle().
			Foreground(t.Text()).
			Bold(true).
			MarginBottom(1),

		// Body: 14px equivalent, regular weight
		Body: styles.NewStyle().
			Foreground(t.Text()),

		// BodyEmphasis: 14px equivalent, bold
		BodyEmphasis: styles.NewStyle().
			Foreground(t.Text()).
			Bold(true),

		// Caption: 12px equivalent, muted
		Caption: styles.NewStyle().
			Foreground(t.TextMuted()).
			Faint(true),

		// Label: 12px equivalent, uppercase, tracked
		Label: styles.NewStyle().
			Foreground(t.TextMuted()).
			Bold(true).
			Transform(lipgloss.Uppercase),

		// Code: monospace, muted background
		Code: styles.NewStyle().
			Foreground(t.Primary()).
			Background(t.BackgroundElement()).
			Padding(0, 1),

		// Link: primary color, underline
		Link: styles.NewStyle().
			Foreground(t.Primary()).
			Underline(true),

		// Error: error color, bold
		Error: styles.NewStyle().
			Foreground(t.Error()).
			Bold(true),

		// Success: success color, bold
		Success: styles.NewStyle().
			Foreground(t.Success()).
			Bold(true),

		// Warning: warning color, bold
		Warning: styles.NewStyle().
			Foreground(t.Warning()).
			Bold(true),

		// Info: info color
		Info: styles.NewStyle().
			Foreground(t.Info()),
	}
}

// Spacing defines consistent spacing values
type Spacing struct {
	None   int
	XS     int // 0.25rem (4px)
	SM     int // 0.5rem (8px)
	MD     int // 1rem (16px)
	LG     int // 1.5rem (24px)
	XL     int // 2rem (32px)
	XXL    int // 3rem (48px)
	XXXL   int // 4rem (64px)
}

// DefaultSpacing provides the default spacing scale
var DefaultSpacing = Spacing{
	None: 0,
	XS:   1,
	SM:   1,
	MD:   2,
	LG:   3,
	XL:   4,
	XXL:  6,
	XXXL: 8,
}

// Container creates a container with consistent padding
func Container(content string) string {
	return styles.NewStyle().
		Padding(DefaultSpacing.MD, DefaultSpacing.LG).
		Render(content)
}

// Section creates a section with spacing
func Section(title, content string) string {
	typo := New()

	titleRendered := typo.Heading.Render(title)
	contentRendered := typo.Body.
		MarginBottom(DefaultSpacing.MD).
		Render(content)

	return titleRendered + "\n" + contentRendered
}

// List creates a bulleted list with proper spacing
func List(items []string) string {
	typo := New()
	result := ""

	for i, item := range items {
		bullet := "• "
		line := bullet + item

		style := typo.Body
		if i < len(items)-1 {
			style = style.MarginBottom(1)
		}

		result += style.Render(line)
		if i < len(items)-1 {
			result += "\n"
		}
	}

	return result
}

// NumberedList creates a numbered list with proper spacing
func NumberedList(items []string) string {
	typo := New()
	t := theme.CurrentTheme()
	result := ""

	for i, item := range items {
		numberStr := fmt.Sprintf("%2d.", i+1)

		number := styles.NewStyle().
			Foreground(t.Primary()).
			Bold(true).
			Render(numberStr)

		line := number + " " + item

		style := typo.Body
		if i < len(items)-1 {
			style = style.MarginBottom(1)
		}

		result += style.Render(line)
		if i < len(items)-1 {
			result += "\n"
		}
	}

	return result
}

// Card creates a card layout with consistent spacing and borders
func Card(title, content string) string {
	t := theme.CurrentTheme()
	typo := New()

	titleRendered := typo.Title.
		Foreground(t.Primary()).
		Render(title)

	contentRendered := typo.Body.
		MarginTop(DefaultSpacing.SM).
		Render(content)

	inner := titleRendered + "\n" + contentRendered

	return styles.NewStyle().
		Padding(DefaultSpacing.MD, DefaultSpacing.LG).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.BorderMuted()).
		Background(t.BackgroundPanel()).
		Render(inner)
}

// Divider creates a visual separator
func Divider(width int) string {
	t := theme.CurrentTheme()

	if width <= 0 {
		width = 60
	}

	line := ""
	for i := 0; i < width; i++ {
		line += "─"
	}

	return styles.NewStyle().
		Foreground(t.BorderMuted()).
		MarginTop(DefaultSpacing.SM).
		MarginBottom(DefaultSpacing.SM).
		Render(line)
}

// Hint creates hint text with proper styling
func Hint(text string) string {
	typo := New()
	return typo.Caption.Render("💡 " + text)
}

// ErrorMessage creates an error message with icon
func ErrorMessage(text string) string {
	typo := New()
	return typo.Error.Render("✗ " + text)
}

// SuccessMessage creates a success message with icon
func SuccessMessage(text string) string {
	typo := New()
	return typo.Success.Render("✓ " + text)
}

// WarningMessage creates a warning message with icon
func WarningMessage(text string) string {
	typo := New()
	return typo.Warning.Render("⚠ " + text)
}

// InfoMessage creates an info message with icon
func InfoMessage(text string) string {
	typo := New()
	return typo.Info.Render("ℹ " + text)
}

// KeyValue creates a key-value pair with proper alignment
func KeyValue(key, value string) string {
	typo := New()
	t := theme.CurrentTheme()

	keyRendered := styles.NewStyle().
		Foreground(t.TextMuted()).
		Bold(true).
		Width(20).
		Render(key + ":")

	valueRendered := typo.Body.Render(value)

	return keyRendered + " " + valueRendered
}

// Badge creates a badge/tag with background
func Badge(text string) string {
	t := theme.CurrentTheme()

	return styles.NewStyle().
		Foreground(t.Background()).
		Background(t.Primary()).
		Bold(true).
		Padding(0, 1).
		Render(text)
}

// StatusBadge creates a colored status badge
func StatusBadge(text string, status string) string {
	t := theme.CurrentTheme()

	var bg lipgloss.AdaptiveColor
	switch status {
	case "success":
		bg = t.Success()
	case "error":
		bg = t.Error()
	case "warning":
		bg = t.Warning()
	case "info":
		bg = t.Info()
	default:
		bg = t.TextMuted()
	}

	return styles.NewStyle().
		Foreground(t.Background()).
		Background(bg).
		Bold(true).
		Padding(0, 1).
		Render(text)
}

// Panel creates a panel with title and content
func Panel(title, content string) string {
	t := theme.CurrentTheme()
	typo := New()

	// Title bar
	titleBar := styles.NewStyle().
		Background(t.BackgroundElement()).
		Foreground(t.Primary()).
		Bold(true).
		Padding(0, DefaultSpacing.MD).
		Width(60).
		Render(title)

	// Content area
	contentArea := styles.NewStyle().
		Foreground(t.Text()).
		Padding(DefaultSpacing.MD).
		Width(60).
		Render(content)

	// Combine
	panel := titleBar + "\n" + contentArea

	return styles.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.BorderMuted()).
		Background(t.BackgroundPanel()).
		Render(panel)
}

// Stack vertically stacks elements with spacing
func Stack(elements []string, spacing int) string {
	result := ""
	for i, elem := range elements {
		result += elem
		if i < len(elements)-1 {
			for j := 0; j < spacing; j++ {
				result += "\n"
			}
		}
	}
	return result
}

// Inline horizontally arranges elements with spacing
func Inline(elements []string, spacing string) string {
	result := ""
	for i, elem := range elements {
		result += elem
		if i < len(elements)-1 {
			result += spacing
		}
	}
	return result
}

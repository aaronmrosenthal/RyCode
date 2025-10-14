package polish

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/lipgloss/v2/compat"
	"github.com/aaronmrosenthal/rycode/internal/accessibility"
	"github.com/aaronmrosenthal/rycode/internal/styles"
	"github.com/aaronmrosenthal/rycode/internal/theme"
)

// MicroInteraction represents a small delightful UI moment
type MicroInteraction struct {
	Type      string
	Message   string
	Icon      string
	Color     string
	Duration  time.Duration
	Animation bool
}

// InteractionType constants
const (
	InteractionSuccess     = "success"
	InteractionError       = "error"
	InteractionInfo        = "info"
	InteractionCelebration = "celebration"
	InteractionHover       = "hover"
	InteractionFocus       = "focus"
)

// ButtonHoverEffect creates a hover effect for buttons
func ButtonHoverEffect(text string, isHovered bool) string {
	t := theme.CurrentTheme()

	if !isHovered {
		return styles.NewStyle().
			Foreground(t.Text()).
			Render(text)
	}

	// Hovered state - add subtle glow effect
	return styles.NewStyle().
		Foreground(t.Primary()).
		Bold(true).
		Render("» " + text + " «")
}

// PulseEffect creates a pulsing animation effect
func PulseEffect(text string, frame int) string {
	if !accessibility.GetSettings().ShouldShowAnimations() {
		return text
	}

	t := theme.CurrentTheme()

	// Pulse between normal and bright every 30 frames
	intensity := 0.5 + 0.5*float64(frame%30)/30.0

	style := styles.NewStyle().
		Foreground(t.Primary())

	if intensity > 0.75 {
		style = style.Bold(true)
	}

	return style.Render(text)
}

// ShakeEffect creates a subtle shake animation (for errors)
func ShakeEffect(text string, frame int) string {
	if !accessibility.GetSettings().ShouldShowAnimations() {
		return text
	}

	// Shake pattern over 10 frames
	offsets := []int{0, 1, -1, 1, -1, 0, 0, 0, 0, 0}
	offset := offsets[frame%len(offsets)]

	if offset > 0 {
		return strings.Repeat(" ", offset) + text
	} else if offset < 0 {
		return text[1:]
	}

	return text
}

// TypewriterEffect simulates typing animation
func TypewriterEffect(text string, frame int) string {
	if !accessibility.GetSettings().ShouldShowAnimations() {
		return text
	}

	charsPerFrame := 2
	visibleChars := frame * charsPerFrame

	if visibleChars >= len(text) {
		return text
	}

	return text[:visibleChars] + "█"
}

// FadeInEffect creates a fade-in animation
func FadeInEffect(text string, frame int, totalFrames int) string {
	if !accessibility.GetSettings().ShouldShowAnimations() {
		return text
	}

	if frame >= totalFrames {
		return text
	}

	t := theme.CurrentTheme()

	// Fade from muted to normal
	if frame < totalFrames/2 {
		return styles.NewStyle().
			Foreground(t.TextMuted()).
			Faint(true).
			Render(text)
	}

	return text
}

// SuccessFlash creates a success flash effect
func SuccessFlash(message string, showCheck bool) string {
	t := theme.CurrentTheme()

	icon := ""
	if showCheck {
		icon = "✓ "
	}

	return styles.NewStyle().
		Foreground(t.Success()).
		Bold(true).
		Render(icon + message)
}

// LoadingSpinner returns animated loading spinner
func LoadingSpinner(frame int) string {
	if !accessibility.GetSettings().ShouldShowAnimations() {
		return "..."
	}

	spinners := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	return spinners[frame%len(spinners)]
}

// ProgressBar creates an animated progress bar
func ProgressBar(current, total int, width int) string {
	t := theme.CurrentTheme()

	if total == 0 {
		return ""
	}

	percentage := float64(current) / float64(total)
	filled := int(percentage * float64(width))

	bar := ""

	// Filled portion
	for i := 0; i < filled; i++ {
		bar += styles.NewStyle().
			Foreground(t.Success()).
			Render("█")
	}

	// Empty portion
	for i := filled; i < width; i++ {
		bar += styles.NewStyle().
			Foreground(t.TextMuted()).
			Faint(true).
			Render("░")
	}

	// Percentage label
	label := styles.NewStyle().
		Foreground(t.Text()).
		Render(fmt.Sprintf(" %d%%", int(percentage*100)))

	return bar + label
}

// Sparkle adds sparkle effect to text
func Sparkle(text string, frame int) string {
	if !accessibility.GetSettings().ShouldShowAnimations() {
		return text
	}

	sparkles := []string{"", "✨", "✨", ""}
	sparkle := sparkles[frame%len(sparkles)]

	return sparkle + " " + text + " " + sparkle
}

// Rainbow creates rainbow effect (for celebrations)
func Rainbow(text string, frame int) string {
	if !accessibility.GetSettings().ShouldShowAnimations() {
		return text
	}

	t := theme.CurrentTheme()
	colors := []compat.AdaptiveColor{
		{Light: lipgloss.Color("#FF0000"), Dark: lipgloss.Color("#FF0000")},
		{Light: lipgloss.Color("#FF7F00"), Dark: lipgloss.Color("#FF7F00")},
		{Light: lipgloss.Color("#FFFF00"), Dark: lipgloss.Color("#FFFF00")},
		{Light: lipgloss.Color("#00FF00"), Dark: lipgloss.Color("#00FF00")},
		{Light: lipgloss.Color("#0000FF"), Dark: lipgloss.Color("#0000FF")},
		{Light: lipgloss.Color("#4B0082"), Dark: lipgloss.Color("#4B0082")},
		{Light: lipgloss.Color("#9400D3"), Dark: lipgloss.Color("#9400D3")},
	}
	colorIndex := frame % len(colors)
	_ = t // Keep theme available for potential future use

	return styles.NewStyle().
		Foreground(colors[colorIndex]).
		Render(text)
}

// GlowEffect adds a glow effect to text
func GlowEffect(text string, frame int) string {
	if !accessibility.GetSettings().ShouldShowAnimations() {
		return text
	}

	t := theme.CurrentTheme()

	// Alternate between normal and glowing
	if frame%20 < 10 {
		return styles.NewStyle().
			Foreground(t.Primary()).
			Bold(true).
			Render("⚡ " + text + " ⚡")
	}

	return styles.NewStyle().
		Foreground(t.Primary()).
		Render(text)
}

// ConfettiEffect creates confetti for celebrations
func ConfettiEffect(width int, frame int) string {
	if !accessibility.GetSettings().ShouldShowAnimations() {
		return ""
	}

	confetti := []string{"🎉", "🎊", "✨", "⭐", "🌟", "💫"}

	// Random confetti positions
	line := strings.Repeat(" ", width)
	lineRunes := []rune(line)

	// Seed with frame for deterministic randomness
	r := rand.New(rand.NewSource(int64(frame)))

	// Add 3-5 confetti pieces
	for i := 0; i < 3+r.Intn(3); i++ {
		pos := r.Intn(width)
		if pos < len(lineRunes) {
			confettiChar := confetti[r.Intn(len(confetti))]
			// This is simplified - real implementation would need proper rune handling
			_ = confettiChar
		}
	}

	return string(lineRunes)
}

// StatusDot creates an animated status indicator
func StatusDot(status string, frame int) string {
	t := theme.CurrentTheme()

	var color = t.TextMuted()
	var icon = "○"

	switch status {
	case "active":
		color = t.Success()
		if accessibility.GetSettings().ShouldShowAnimations() && frame%20 < 10 {
			icon = "●"
		} else {
			icon = "◉"
		}
	case "warning":
		color = t.Warning()
		icon = "⚠"
	case "error":
		color = t.Error()
		icon = "✗"
	case "loading":
		if accessibility.GetSettings().ShouldShowAnimations() {
			spinners := []string{"◐", "◓", "◑", "◒"}
			icon = spinners[frame%len(spinners)]
		} else {
			icon = "○"
		}
		color = t.Info()
	}

	return styles.NewStyle().
		Foreground(color).
		Render(icon)
}

// CardHoverEffect adds elevation effect to cards
func CardHoverEffect(content string, isHovered bool) string {
	t := theme.CurrentTheme()

	borderColor := t.Border()
	if isHovered {
		borderColor = t.Primary()
	}

	// Would apply to lipgloss border
	_ = borderColor

	return content
}

// NumberCounter animates number changes
func NumberCounter(from, to int, progress float64) string {
	current := from + int(float64(to-from)*progress)
	return fmt.Sprintf("%d", current)
}

// SmoothScroll calculates smooth scroll offset
func SmoothScroll(target, current int, speed float64) int {
	diff := target - current
	step := int(float64(diff) * speed)

	if step == 0 && diff != 0 {
		if diff > 0 {
			step = 1
		} else {
			step = -1
		}
	}

	return current + step
}

// ElasticBounce creates elastic bounce effect
func ElasticBounce(frame, totalFrames int) float64 {
	if !accessibility.GetSettings().ShouldShowAnimations() {
		return 1.0
	}

	if frame >= totalFrames {
		return 1.0
	}

	progress := float64(frame) / float64(totalFrames)

	// Elastic easing out
	if progress == 0 || progress == 1 {
		return progress
	}

	p := 0.3
	s := p / 4

	// Calculate the exponent and ensure it's non-negative
	exponent := 10 * (progress - 1)
	if exponent < 0 {
		exponent = -exponent
	}

	// Bitshift must be done on integer, then convert to float64
	shiftAmount := uint(int(exponent))
	if shiftAmount > 31 {
		shiftAmount = 31 // Prevent overflow
	}

	return 1.0 + (-1.0)*float64(int(1)<<shiftAmount)*
		float64(float64(2)*float64(3.14159)*(progress-1-s)/p)
}

// PulseScale creates pulsing scale effect
func PulseScale(frame int) float64 {
	if !accessibility.GetSettings().ShouldShowAnimations() {
		return 1.0
	}

	// Pulse between 1.0 and 1.1 scale
	return 1.0 + 0.1*float64(frame%30)/30.0
}

// Notification creates a notification toast
type Notification struct {
	Message  string
	Type     string // success, error, info, warning
	Duration time.Duration
	Icon     string
}

// RenderNotification renders a notification with appropriate styling
func RenderNotification(notif Notification, frame int) string {
	t := theme.CurrentTheme()

	var color = t.Info()
	var icon = "ℹ"

	switch notif.Type {
	case "success":
		color = t.Success()
		icon = "✓"
	case "error":
		color = t.Error()
		icon = "✗"
	case "warning":
		color = t.Warning()
		icon = "⚠"
	}

	if notif.Icon != "" {
		icon = notif.Icon
	}

	iconStyle := styles.NewStyle().
		Foreground(color).
		Bold(true).
		Render(icon + " ")

	messageStyle := styles.NewStyle().
		Foreground(t.Text()).
		Render(notif.Message)

	notification := iconStyle + messageStyle

	// Fade in/out animation
	if accessibility.GetSettings().ShouldShowAnimations() {
		totalFrames := 60 // 1 second at 60fps
		if frame < 10 {
			// Fade in
			return FadeInEffect(notification, frame, 10)
		} else if frame > totalFrames-10 {
			// Fade out
			return FadeInEffect(notification, totalFrames-frame, 10)
		}
	}

	return notification
}

// Tooltip creates a tooltip with pointer
func Tooltip(text string, position string) string {
	t := theme.CurrentTheme()

	tooltipStyle := styles.NewStyle().
		Foreground(t.Background()).
		Background(t.Text()).
		Padding(0, 1)

	tooltip := tooltipStyle.Render(text)

	// Add pointer based on position
	pointer := ""
	switch position {
	case "top":
		pointer = "▼"
	case "bottom":
		pointer = "▲"
	case "left":
		pointer = "▶"
	case "right":
		pointer = "◀"
	}

	pointerStyle := styles.NewStyle().
		Foreground(t.Text())

	if position == "top" || position == "bottom" {
		return tooltip + "\n" + pointerStyle.Render(pointer)
	}

	return pointerStyle.Render(pointer) + " " + tooltip
}

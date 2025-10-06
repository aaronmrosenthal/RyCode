package components

import (
	"fmt"
	"strings"

	"github.com/aaronmrosenthal/rycode/packages/tui-v2/internal/theme"
	"github.com/charmbracelet/lipgloss"
)

// TokenMeter displays token usage with a visual progress bar
type TokenMeter struct {
	PromptTokens   int
	ResponseTokens int
	MaxTokens      int
	Width          int
	ShowBar        bool // Show visual bar (disable for small screens)
	Animated       bool // Enable pulsing animation
	Frame          int  // Animation frame
}

// NewTokenMeter creates a new token meter
func NewTokenMeter(promptTokens, responseTokens, maxTokens, width int) TokenMeter {
	return TokenMeter{
		PromptTokens:   promptTokens,
		ResponseTokens: responseTokens,
		MaxTokens:      maxTokens,
		Width:          width,
		ShowBar:        width >= 60, // Only show bar on larger screens
		Animated:       false,
		Frame:          0,
	}
}

// Render renders the token meter
func (tm TokenMeter) Render() string {
	total := tm.PromptTokens + tm.ResponseTokens
	percentage := float64(0)
	if tm.MaxTokens > 0 {
		percentage = float64(total) / float64(tm.MaxTokens)
	}

	// Choose color based on usage
	var barColor lipgloss.Color
	var icon string
	switch {
	case percentage > 0.95:
		barColor = theme.NeonPink // Critical (>95%)
		icon = "🔴"
	case percentage > 0.85:
		barColor = theme.NeonOrange // High (>85%)
		icon = "🟠"
	case percentage > 0.70:
		barColor = theme.NeonYellow // Warning (>70%)
		icon = "🟡"
	default:
		barColor = theme.MatrixGreen // Normal
		icon = "🟢"
	}

	// Add pulse animation if enabled and in warning state
	if tm.Animated && percentage > 0.85 {
		intensity := 0.6 + 0.4*(float64(tm.Frame%30)/30.0)
		barColor = theme.InterpolateBrightness(barColor, intensity)
	}

	// Render compact version for small screens
	if !tm.ShowBar {
		return tm.renderCompact(total, percentage, icon)
	}

	// Render full version with bar
	return tm.renderFull(total, percentage, barColor, icon)
}

// renderCompact renders a compact version without bar
func (tm TokenMeter) renderCompact(total int, percentage float64, icon string) string {
	label := fmt.Sprintf("%s Tokens: %d/%d (%.0f%%)",
		icon, total, tm.MaxTokens, percentage*100)

	style := lipgloss.NewStyle().Foreground(theme.MatrixGreenDim)
	return style.Render(label)
}

// renderFull renders full version with visual bar
func (tm TokenMeter) renderFull(total int, percentage float64, barColor lipgloss.Color, icon string) string {
	// Visual bar
	barWidth := tm.Width - 30 // Reserve space for label
	if barWidth < 10 {
		barWidth = 10
	}
	if barWidth > 50 {
		barWidth = 50
	}

	filled := int(float64(barWidth) * percentage)
	if filled > barWidth {
		filled = barWidth
	}

	bar := strings.Repeat("█", filled) +
		strings.Repeat("░", barWidth-filled)

	barStyle := lipgloss.NewStyle().Foreground(barColor)

	label := fmt.Sprintf(" %s %d/%d (%.0f%%)",
		icon, total, tm.MaxTokens, percentage*100)

	labelStyle := lipgloss.NewStyle().Foreground(theme.MatrixGreenDim)

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		barStyle.Render(bar),
		labelStyle.Render(label),
	)
}

// SetAnimated enables/disables animation
func (tm *TokenMeter) SetAnimated(animated bool) {
	tm.Animated = animated
}

// SetFrame updates the animation frame
func (tm *TokenMeter) SetFrame(frame int) {
	tm.Frame = frame
}

// UpdateTokens updates token counts
func (tm *TokenMeter) UpdateTokens(promptTokens, responseTokens int) {
	tm.PromptTokens = promptTokens
	tm.ResponseTokens = responseTokens
}

// GetPercentage returns current usage percentage
func (tm TokenMeter) GetPercentage() float64 {
	if tm.MaxTokens == 0 {
		return 0
	}
	total := tm.PromptTokens + tm.ResponseTokens
	return float64(total) / float64(tm.MaxTokens)
}

// IsWarning returns true if usage is in warning range
func (tm TokenMeter) IsWarning() bool {
	return tm.GetPercentage() > 0.70
}

// IsCritical returns true if usage is in critical range
func (tm TokenMeter) IsCritical() bool {
	return tm.GetPercentage() > 0.95
}

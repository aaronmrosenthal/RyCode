package theme

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// GradientText creates a horizontal gradient effect across text
func GradientText(text string, from, to lipgloss.Color) string {
	if len(text) == 0 {
		return ""
	}

	if len(text) == 1 {
		return lipgloss.NewStyle().Foreground(from).Render(text)
	}

	fromRGB := hexToRGB(string(from))
	toRGB := hexToRGB(string(to))

	var result strings.Builder
	runes := []rune(text)

	for i, char := range runes {
		progress := float64(i) / float64(len(runes)-1)
		color := interpolateRGB(fromRGB, toRGB, progress)
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(rgbToHex(color)))
		result.WriteString(style.Render(string(char)))
	}

	return result.String()
}

// GradientTextPreset applies a predefined gradient preset
func GradientTextPreset(text string, preset GradientPreset) string {
	return GradientText(text, preset.From, preset.To)
}

// GlowText simulates a glow effect (limited in terminal, uses bold + brightness)
func GlowText(text string, color lipgloss.Color, intensity float64) string {
	style := lipgloss.NewStyle().Foreground(color)

	// Intensity > 0.5 = bold
	if intensity > 0.5 {
		style = style.Bold(true)
	}

	// Intensity > 0.7 = use brighter color variant
	if intensity > 0.7 && color == MatrixGreen {
		style = style.Foreground(MatrixGreenBright)
	}

	return style.Render(text)
}

// PulseText creates a pulsing effect (for animations)
// frame: animation frame number (0-N)
// speed: how fast to pulse (higher = faster)
func PulseText(text string, color lipgloss.Color, frame int, speed float64) string {
	// Calculate intensity using sine wave
	intensity := (math.Sin(float64(frame)*speed) + 1.0) / 2.0 // 0.0 - 1.0

	return GlowText(text, color, intensity)
}

// RainbowText creates a rainbow effect using multiple colors
func RainbowText(text string) string {
	if len(text) == 0 {
		return ""
	}

	colors := []lipgloss.Color{
		NeonPink,
		NeonOrange,
		NeonYellow,
		MatrixGreen,
		NeonCyan,
		NeonBlue,
		NeonPurple,
	}

	var result strings.Builder
	runes := []rune(text)

	for i, char := range runes {
		colorIndex := i % len(colors)
		style := lipgloss.NewStyle().Foreground(colors[colorIndex])
		result.WriteString(style.Render(string(char)))
	}

	return result.String()
}

// BoxText wraps text in a Matrix-themed box
func BoxText(text string, width int) string {
	style := lipgloss.NewStyle().
		Width(width).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(MatrixGreen).
		Padding(1, 2).
		Align(lipgloss.Center)

	return style.Render(text)
}

// ShadowText adds a "shadow" effect (simulated with dark background)
func ShadowText(text string) string {
	fg := lipgloss.NewStyle().
		Foreground(MatrixGreen).
		Bold(true)

	shadow := lipgloss.NewStyle().
		Foreground(MatrixGreenDark).
		MarginLeft(1).
		MarginTop(1)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		fg.Render(text),
		shadow.Render(text),
	)
}

// RGB represents an RGB color
type RGB struct {
	R, G, B int
}

// hexToRGB converts hex color to RGB
func hexToRGB(hex string) RGB {
	var r, g, b int

	// Remove # if present
	hex = strings.TrimPrefix(hex, "#")

	// Parse hex string
	if len(hex) == 6 {
		fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	} else if len(hex) == 3 {
		// Handle short form (e.g., #0f0)
		fmt.Sscanf(hex, "%1x%1x%1x", &r, &g, &b)
		r, g, b = r*17, g*17, b*17
	}

	return RGB{R: r, G: g, B: b}
}

// rgbToHex converts RGB to hex color
func rgbToHex(rgb RGB) string {
	return fmt.Sprintf("#%02x%02x%02x", clamp(rgb.R), clamp(rgb.G), clamp(rgb.B))
}

// interpolateRGB interpolates between two RGB colors
func interpolateRGB(from, to RGB, progress float64) RGB {
	return RGB{
		R: int(float64(from.R) + (float64(to.R)-float64(from.R))*progress),
		G: int(float64(from.G) + (float64(to.G)-float64(from.G))*progress),
		B: int(float64(from.B) + (float64(to.B)-float64(from.B))*progress),
	}
}

// clamp ensures value is within 0-255 range
func clamp(val int) int {
	if val < 0 {
		return 0
	}
	if val > 255 {
		return 255
	}
	return val
}

// MatrixRain generates a single column of Matrix rain effect
// height: number of characters in the column
// offset: animation offset
func MatrixRain(height int, offset int) string {
	chars := "ｦｱｳｴｵｶｷｹｺｻｼｽｾｿﾀﾂﾃﾅﾆﾇﾈﾊﾋﾎﾏﾐﾑﾒﾓﾔﾕﾗﾘﾜ0123456789ZXCVBNM"
	runes := []rune(chars)

	var result strings.Builder

	for i := 0; i < height; i++ {
		// Calculate fade based on position in column
		fade := 1.0 - (float64(i) / float64(height))

		// Select character
		charIndex := (i + offset) % len(runes)
		char := runes[charIndex]

		// Apply color based on fade
		var color lipgloss.Color
		if fade > 0.8 {
			color = MatrixGreenBright
		} else if fade > 0.5 {
			color = MatrixGreen
		} else if fade > 0.2 {
			color = MatrixGreenDim
		} else {
			color = MatrixGreenDark
		}

		style := lipgloss.NewStyle().Foreground(color)
		result.WriteString(style.Render(string(char)) + "\n")
	}

	return result.String()
}

// ScanlineEffect adds horizontal scanline effect
func ScanlineEffect(text string, lineHeight int) string {
	lines := strings.Split(text, "\n")
	var result strings.Builder

	for i, line := range lines {
		if i%lineHeight == 0 {
			// Scanline - slightly dimmer
			style := lipgloss.NewStyle().Foreground(MatrixGreenDim)
			result.WriteString(style.Render(line) + "\n")
		} else {
			result.WriteString(line + "\n")
		}
	}

	return strings.TrimSuffix(result.String(), "\n")
}

// MatrixRainBackground represents an animated background with falling Matrix characters
type MatrixRainBackground struct {
	Width       int
	Height      int
	ColumnCount int
	Columns     []MatrixRainColumn
}

// MatrixRainColumn represents a single column of falling characters
type MatrixRainColumn struct {
	X      int // X position
	Y      int // Y position (top of trail)
	Speed  int // Fall speed (pixels per frame)
	Length int // Trail length
	Offset int // Character offset for variation
}

// NewMatrixRainBackground creates a new Matrix rain background
func NewMatrixRainBackground(width, height int) MatrixRainBackground {
	columnCount := width / 3 // Sparse columns to avoid clutter
	if columnCount < 1 {
		columnCount = 1
	}

	columns := make([]MatrixRainColumn, columnCount)
	for i := range columns {
		columns[i] = MatrixRainColumn{
			X:      (i * 3) + (i % 2), // Stagger positions
			Y:      -(i % height),     // Start at different heights
			Speed:  1 + (i % 2),       // Vary speeds
			Length: 5 + (i % 8),       // Vary trail lengths
			Offset: i * 7,             // Character variation
		}
	}

	return MatrixRainBackground{
		Width:       width,
		Height:      height,
		ColumnCount: columnCount,
		Columns:     columns,
	}
}

// Update updates the rain animation
func (mrb *MatrixRainBackground) Update() {
	for i := range mrb.Columns {
		// Move column down
		mrb.Columns[i].Y += mrb.Columns[i].Speed

		// Reset if off screen
		if mrb.Columns[i].Y > mrb.Height+mrb.Columns[i].Length {
			mrb.Columns[i].Y = -mrb.Columns[i].Length
			mrb.Columns[i].Offset += 13 // Change character pattern
		}
	}
}

// Render renders the Matrix rain (very dim to not distract)
func (mrb MatrixRainBackground) Render() string {
	// Create empty grid
	grid := make([][]rune, mrb.Height)
	for i := range grid {
		grid[i] = make([]rune, mrb.Width)
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}

	chars := "ｦｱｳｴｵｶｷｹｺｻｼｽｾｿﾀﾂﾃﾅﾆﾇﾈﾊﾋﾎﾏﾐﾑﾒﾓﾔﾕﾗﾘﾜ0123456789"
	runes := []rune(chars)

	// Draw columns
	for _, col := range mrb.Columns {
		for i := 0; i < col.Length; i++ {
			y := col.Y - i
			if y >= 0 && y < mrb.Height && col.X < mrb.Width {
				// Select character
				charIndex := (col.Offset + i) % len(runes)
				grid[y][col.X] = runes[charIndex]
			}
		}
	}

	// Render grid with very dim colors
	var result strings.Builder
	dimStyle := lipgloss.NewStyle().Foreground(MatrixGreenVDark)

	for _, row := range grid {
		result.WriteString(dimStyle.Render(string(row)) + "\n")
	}

	return strings.TrimSuffix(result.String(), "\n")
}

// AnimatedGradient represents a time-based color gradient animation
type AnimatedGradient struct {
	Colors    []lipgloss.Color // Multi-color palette to cycle through
	Duration  time.Duration    // Full animation cycle duration
	StartTime time.Time        // When animation started
}

// NewAnimatedGradient creates a new animated gradient
func NewAnimatedGradient(colors []lipgloss.Color, duration time.Duration) AnimatedGradient {
	return AnimatedGradient{
		Colors:    colors,
		Duration:  duration,
		StartTime: time.Now(),
	}
}

// ColorAt returns the interpolated color at a given position and current time
func (ag AnimatedGradient) ColorAt(position float64) lipgloss.Color {
	if len(ag.Colors) == 0 {
		return MatrixGreen
	}

	if len(ag.Colors) == 1 {
		return ag.Colors[0]
	}

	// Calculate animation progress (0.0 - 1.0)
	elapsed := time.Since(ag.StartTime)
	timeProgress := float64(elapsed%ag.Duration) / float64(ag.Duration)

	// Shift position based on time (creates wave effect)
	shiftedPosition := position + timeProgress
	if shiftedPosition > 1.0 {
		shiftedPosition -= 1.0
	}

	// Determine which two colors to interpolate between
	numSegments := len(ag.Colors) - 1
	segmentSize := 1.0 / float64(numSegments)
	segmentIndex := int(shiftedPosition / segmentSize)

	if segmentIndex >= numSegments {
		segmentIndex = numSegments - 1
	}

	// Calculate interpolation within segment
	segmentProgress := (shiftedPosition - float64(segmentIndex)*segmentSize) / segmentSize

	// Interpolate between colors
	fromRGB := hexToRGB(string(ag.Colors[segmentIndex]))
	toRGB := hexToRGB(string(ag.Colors[segmentIndex+1]))
	resultRGB := interpolateRGB(fromRGB, toRGB, segmentProgress)

	return lipgloss.Color(rgbToHex(resultRGB))
}

// AnimatedGradientText creates text with time-based animated gradient
func AnimatedGradientText(text string, ag AnimatedGradient) string {
	if len(text) == 0 {
		return ""
	}

	if len(text) == 1 {
		return lipgloss.NewStyle().Foreground(ag.ColorAt(0)).Render(text)
	}

	var result strings.Builder
	runes := []rune(text)

	for i, char := range runes {
		position := float64(i) / float64(len(runes)-1)
		color := ag.ColorAt(position)
		style := lipgloss.NewStyle().Foreground(color)
		result.WriteString(style.Render(string(char)))
	}

	return result.String()
}

// BreathingBorder creates a pulsing border effect
func BreathingBorder(content string, baseColor lipgloss.Color, frame int, width int) string {
	// Calculate intensity using sine wave (0.4 - 1.0 for subtle effect)
	intensity := 0.4 + 0.6*(math.Sin(float64(frame)*0.08)+1.0)/2.0

	// Adjust color brightness based on intensity
	rgb := hexToRGB(string(baseColor))
	adjustedRGB := RGB{
		R: int(float64(rgb.R) * intensity),
		G: int(float64(rgb.G) * intensity),
		B: int(float64(rgb.B) * intensity),
	}
	borderColor := lipgloss.Color(rgbToHex(adjustedRGB))

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1).
		Width(width - 4)

	return style.Render(content)
}

// InterpolateBrightness adjusts color brightness
func InterpolateBrightness(baseColor lipgloss.Color, intensity float64) lipgloss.Color {
	rgb := hexToRGB(string(baseColor))
	adjustedRGB := RGB{
		R: int(float64(rgb.R) * intensity),
		G: int(float64(rgb.G) * intensity),
		B: int(float64(rgb.B) * intensity),
	}
	return lipgloss.Color(rgbToHex(adjustedRGB))
}

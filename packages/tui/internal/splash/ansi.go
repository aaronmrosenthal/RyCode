package splash

import (
	"fmt"
	"math"
)

// RGB represents an RGB color
type RGB struct {
	R, G, B uint8
}

// ANSI returns the ANSI truecolor escape sequence for this color
func (c RGB) ANSI() string {
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", c.R, c.G, c.B)
}

// ResetColor returns the ANSI reset sequence
func ResetColor() string {
	return "\033[0m"
}

// Colorize wraps text in ANSI color codes
func Colorize(text string, color RGB) string {
	return color.ANSI() + text + ResetColor()
}

// lerp linearly interpolates between two values
func lerp(a, b uint8, t float64) uint8 {
	return uint8(float64(a)*(1.0-t) + float64(b)*t)
}

// lerpRGB linearly interpolates between two RGB colors
func lerpRGB(a, b RGB, t float64) RGB {
	return RGB{
		R: lerp(a.R, b.R, t),
		G: lerp(a.G, b.G, t),
		B: lerp(a.B, b.B, t),
	}
}

// GradientColor returns a color from the cyan-to-magenta gradient
// based on angle (0 to 2π)
func GradientColor(angle float64) RGB {
	cyan := RGB{0, 255, 255}    // #00FFFF
	magenta := RGB{255, 0, 255} // #FF00FF

	// Normalize angle to [0, 1]
	t := math.Mod(angle, 2*math.Pi) / (2 * math.Pi)
	if t < 0 {
		t += 1.0
	}

	return lerpRGB(cyan, magenta, t)
}

// ColorMode represents terminal color capabilities
type ColorMode int

const (
	Colors16 ColorMode = iota
	Colors256
	Truecolor
)

// String returns the name of the color mode
func (cm ColorMode) String() string {
	switch cm {
	case Colors16:
		return "16-color"
	case Colors256:
		return "256-color"
	case Truecolor:
		return "truecolor"
	default:
		return "unknown"
	}
}

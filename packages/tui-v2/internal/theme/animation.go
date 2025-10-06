package theme

import "time"

// Animation constants for consistent timing across all components
const (
	// AnimationFPS is the target frames per second for all animations
	AnimationFPS = 30

	// AnimationFrameDuration is the duration of each animation frame
	AnimationFrameDuration = time.Second / AnimationFPS

	// PulseFrameCycle is the number of frames in a complete pulse cycle
	// This creates a 1-second pulse when AnimationFPS = 30
	PulseFrameCycle = 30

	// MinLogoWidth is the minimum terminal width to show the full ASCII logo
	MinLogoWidth = 60
)

// CalculatePulseIntensity calculates a pulsing intensity value for animations
// that cycle over PulseFrameCycle frames.
//
// Parameters:
//   - frame: Current animation frame number
//   - baseIntensity: Minimum intensity value (e.g., 0.5 for 50-100% range)
//   - intensityRange: Range of intensity variation (e.g., 0.5 for 50% range)
//
// Returns intensity value between baseIntensity and (baseIntensity + intensityRange)
//
// Example:
//
//	CalculatePulseIntensity(15, 0.5, 0.5) -> returns value between 0.5 and 1.0
//	CalculatePulseIntensity(15, 0.6, 0.4) -> returns value between 0.6 and 1.0
func CalculatePulseIntensity(frame int, baseIntensity, intensityRange float64) float64 {
	// Calculate progress through the pulse cycle (0.0 to 1.0)
	progress := float64(frame%PulseFrameCycle) / float64(PulseFrameCycle)

	// Return intensity value in the specified range
	return baseIntensity + intensityRange*progress
}

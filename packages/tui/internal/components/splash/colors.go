package splash

import "time"

// Animation timing constants for cortex animations
const (
	CortexAnimationDuration     = 1200 * time.Millisecond // Total animation duration
	CortexFadeInDuration        = 300 * time.Millisecond  // Fade in phase
	CortexHoldDuration          = 600 * time.Millisecond  // Hold at full brightness
	CortexFadeOutDuration       = 300 * time.Millisecond  // Fade out phase
	CortexAnimationTickInterval = 50 * time.Millisecond   // 20 FPS for smooth animation
	CortexFadeOutStart          = 900 * time.Millisecond  // When fade out begins (fade-in + hold)
	CortexMinVisibleOpacity     = 0.05                    // Don't render below this opacity
)

// GetProviderBrandColor returns the brand color for a given provider ID
// This is the single source of truth for provider colors throughout the application
func GetProviderBrandColor(providerID string) RGB {
	switch providerID {
	case "anthropic", "claude":
		// Claude brand: warm orange/peach #E07856
		return RGB{R: 224, G: 120, B: 86}
	case "google", "gemini":
		// Gemini brand: blue-to-purple gradient (using mid-purple) #8B7FD8
		return RGB{R: 139, G: 127, B: 216}
	case "openai", "codex":
		// OpenAI/Codex brand: teal/cyan #10A37F
		return RGB{R: 16, G: 163, B: 127}
	case "xai", "grok":
		// Grok/xAI brand: red #FF4444
		return RGB{R: 255, G: 68, B: 68}
	case "qwen":
		// Qwen brand: golden orange (from badge) #FFA726
		return RGB{R: 255, G: 167, B: 38}
	default:
		// Default: cyan #00FFFF for unknown providers
		return RGB{R: 0, G: 255, B: 255}
	}
}

// CalculateAnimationOpacity calculates the current opacity based on animation timing
// Returns (opacity, finished) where finished=true means animation is complete
func CalculateAnimationOpacity(elapsed time.Duration, startOpacity float64) (opacity float64, finished bool) {
	// Animation complete
	if elapsed >= CortexAnimationDuration {
		return 0.0, true
	}

	// Fade in phase (0 -> startOpacity)
	if elapsed < CortexFadeInDuration {
		progress := float64(elapsed) / float64(CortexFadeInDuration)
		return startOpacity * progress, false
	}

	// Hold phase (full opacity)
	if elapsed < CortexFadeOutStart {
		return startOpacity, false
	}

	// Fade out phase (startOpacity -> 0)
	fadeOutElapsed := elapsed - CortexFadeOutStart
	progress := float64(fadeOutElapsed) / float64(CortexFadeOutDuration)
	return startOpacity * (1.0 - progress), false
}

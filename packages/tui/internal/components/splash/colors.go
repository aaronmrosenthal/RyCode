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

// providerBrandColors is a pre-computed map for O(1) color lookups
// This optimization reduces repeated switch statement overhead
var providerBrandColors = map[string]RGB{
	"anthropic": {R: 224, G: 120, B: 86},  // Claude brand: warm orange/peach #E07856
	"claude":    {R: 224, G: 120, B: 86},  // Claude brand: warm orange/peach #E07856
	"google":    {R: 139, G: 127, B: 216}, // Gemini brand: blue-to-purple gradient #8B7FD8
	"gemini":    {R: 139, G: 127, B: 216}, // Gemini brand: blue-to-purple gradient #8B7FD8
	"openai":    {R: 16, G: 163, B: 127},  // OpenAI/Codex brand: teal/cyan #10A37F
	"codex":     {R: 16, G: 163, B: 127},  // OpenAI/Codex brand: teal/cyan #10A37F
	"xai":       {R: 255, G: 68, B: 68},   // Grok/xAI brand: red #FF4444
	"grok":      {R: 255, G: 68, B: 68},   // Grok/xAI brand: red #FF4444
	"qwen":      {R: 255, G: 167, B: 38},  // Qwen brand: golden orange #FFA726
}

// defaultBrandColor is the fallback color for unknown providers
var defaultBrandColor = RGB{R: 0, G: 255, B: 255} // Cyan #00FFFF

// GetProviderBrandColor returns the brand color for a given provider ID
// This is the single source of truth for provider colors throughout the application
// Optimized with map lookup for O(1) performance
func GetProviderBrandColor(providerID string) RGB {
	if color, ok := providerBrandColors[providerID]; ok {
		return color
	}
	return defaultBrandColor
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

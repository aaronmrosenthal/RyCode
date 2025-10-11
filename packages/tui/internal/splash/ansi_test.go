package splash

import (
	"math"
	"testing"
)

// TestColorizeBasic tests basic colorization
func TestColorizeBasic(t *testing.T) {
	text := "Hello"
	color := RGB{255, 0, 0} // Red

	result := Colorize(text, color)

	// Should contain the text
	if result[len(result)-len(text)-len(ResetColor()):len(result)-len(ResetColor())] != text {
		t.Error("Colorized text doesn't contain original text")
	}

	// Should start with ANSI code
	if result[:2] != "\033[" {
		t.Error("Colorized text doesn't start with ANSI escape")
	}

	// Should end with reset
	if result[len(result)-len(ResetColor()):] != ResetColor() {
		t.Error("Colorized text doesn't end with reset code")
	}
}

// TestGradientColor tests gradient generation
func TestGradientColor(t *testing.T) {
	tests := []struct {
		angle    float64
		expected string // Rough color name
	}{
		{0.0, "cyan"},           // Start at cyan
		{math.Pi, "magenta"},    // Halfway to magenta
		{2 * math.Pi, "cyan"},   // Back to cyan
	}

	for _, tt := range tests {
		color := GradientColor(tt.angle)

		// Check it's a valid RGB
		if color.R == 0 && color.G == 0 && color.B == 0 {
			t.Errorf("Gradient at %.2f produced black", tt.angle)
		}

		// Cyan check (at 0)
		if tt.expected == "cyan" {
			if color.R > 50 || color.G < 200 || color.B < 200 {
				t.Errorf("Expected cyan-ish at %.2f, got RGB(%d,%d,%d)",
					tt.angle, color.R, color.G, color.B)
			}
		}

		// Magenta check (at π)
		if tt.expected == "magenta" {
			// At π, we should have high R and B, low G
			if color.R < 100 || color.G > 150 || color.B < 100 {
				t.Errorf("Expected magenta-ish at %.2f, got RGB(%d,%d,%d)",
					tt.angle, color.R, color.G, color.B)
			}
		}
	}
}

// TestLerpRGB tests RGB interpolation
func TestLerpRGB(t *testing.T) {
	red := RGB{255, 0, 0}
	blue := RGB{0, 0, 255}

	// At t=0, should be red
	result := lerpRGB(red, blue, 0.0)
	if result != red {
		t.Errorf("lerp at t=0 should be first color, got %+v", result)
	}

	// At t=1, should be blue
	result = lerpRGB(red, blue, 1.0)
	if result != blue {
		t.Errorf("lerp at t=1 should be second color, got %+v", result)
	}

	// At t=0.5, should be middle
	result = lerpRGB(red, blue, 0.5)
	expected := RGB{127, 0, 127}
	if result.R < 120 || result.R > 135 {
		t.Errorf("lerp at t=0.5 should be ~middle, got %+v (expected ~%+v)", result, expected)
	}
}

// TestANSIFormat tests ANSI code format
func TestANSIFormat(t *testing.T) {
	color := RGB{100, 150, 200}
	ansi := color.ANSI()

	// Should have truecolor format: \033[38;2;R;G;Bm
	expected := "\033[38;2;100;150;200m"
	if ansi != expected {
		t.Errorf("ANSI format incorrect: got %q, expected %q", ansi, expected)
	}
}

// TestResetColor tests reset code
func TestResetColor(t *testing.T) {
	reset := ResetColor()
	expected := "\033[0m"

	if reset != expected {
		t.Errorf("Reset color incorrect: got %q, expected %q", reset, expected)
	}
}

// BenchmarkGradientColor benchmarks gradient calculation
func BenchmarkGradientColor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = GradientColor(float64(i) * 0.1)
	}
}

// BenchmarkColorize benchmarks text colorization
func BenchmarkColorize(b *testing.B) {
	color := RGB{255, 100, 50}
	text := "Test"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Colorize(text, color)
	}
}

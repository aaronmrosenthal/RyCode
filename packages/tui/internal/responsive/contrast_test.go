package responsive

import (
	"math"
	"testing"
)

func TestContrastRatio(t *testing.T) {
	tests := []struct {
		name     string
		fg       RGB
		bg       RGB
		expected float64
	}{
		{
			name:     "White on Black",
			fg:       RGB{255, 255, 255},
			bg:       RGB{0, 0, 0},
			expected: 21.0,
		},
		{
			name:     "Black on White",
			fg:       RGB{0, 0, 0},
			bg:       RGB{255, 255, 255},
			expected: 21.0,
		},
		{
			name:     "Same Color",
			fg:       RGB{128, 128, 128},
			bg:       RGB{128, 128, 128},
			expected: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ratio := ContrastRatio(tt.fg, tt.bg)

			// Allow small floating point errors
			if math.Abs(ratio-tt.expected) > 0.01 {
				t.Errorf("Expected ratio %.2f, got %.2f", tt.expected, ratio)
			}
		})
	}
}

func TestMeetsWCAG_AA(t *testing.T) {
	tests := []struct {
		ratio       float64
		isLargeText bool
		expected    bool
	}{
		{4.5, false, true},   // Normal text, meets AA
		{4.4, false, false},  // Normal text, fails AA
		{3.0, true, true},    // Large text, meets AA
		{2.9, true, false},   // Large text, fails AA
		{7.0, false, true},   // Normal text, exceeds AA
		{4.5, true, true},    // Large text, exceeds AA
	}

	for _, tt := range tests {
		result := MeetsWCAG_AA(tt.ratio, tt.isLargeText)
		if result != tt.expected {
			t.Errorf("MeetsWCAG_AA(%.1f, %v) = %v, expected %v",
				tt.ratio, tt.isLargeText, result, tt.expected)
		}
	}
}

func TestParseHex(t *testing.T) {
	tests := []struct {
		hex      string
		expected RGB
		hasError bool
	}{
		{"#FF0000", RGB{255, 0, 0}, false},
		{"#00FF00", RGB{0, 255, 0}, false},
		{"#0000FF", RGB{0, 0, 255}, false},
		{"#FFF", RGB{255, 255, 255}, false}, // Shorthand
		{"#000", RGB{0, 0, 0}, false},       // Shorthand
		{"invalid", RGB{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.hex, func(t *testing.T) {
			rgb, err := parseHex(tt.hex)

			if tt.hasError {
				if err == nil {
					t.Error("Expected error, got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if rgb != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, rgb)
			}
		})
	}
}

func TestANSIToRGB(t *testing.T) {
	tests := []struct {
		ansi     int
		expected RGB
	}{
		{0, RGB{0, 0, 0}},           // Black
		{7, RGB{192, 192, 192}},     // White
		{9, RGB{255, 0, 0}},         // Bright Red
		{15, RGB{255, 255, 255}},    // Bright White
		{16, RGB{0, 0, 0}},          // Color cube start
		{232, RGB{8, 8, 8}},         // Grayscale start
		{255, RGB{238, 238, 238}},   // Grayscale end
	}

	for _, tt := range tests {
		t.Run(string(rune('0'+tt.ansi)), func(t *testing.T) {
			rgb := ansiToRGB(tt.ansi)

			// Allow some variance for grayscale/color cube
			if math.Abs(float64(rgb.R-tt.expected.R)) > 5 ||
				math.Abs(float64(rgb.G-tt.expected.G)) > 5 ||
				math.Abs(float64(rgb.B-tt.expected.B)) > 5 {
				t.Errorf("ANSI %d: expected %v, got %v", tt.ansi, tt.expected, rgb)
			}
		})
	}
}

func TestContrastChecker(t *testing.T) {
	checker := NewContrastChecker()

	// Good contrast (white on black)
	checker.Check("component1",
		RGB{255, 255, 255}.ToLipgloss(),
		RGB{0, 0, 0}.ToLipgloss(),
		false,
	)

	// Bad contrast (light gray on white)
	checker.Check("component2",
		RGB{200, 200, 200}.ToLipgloss(),
		RGB{255, 255, 255}.ToLipgloss(),
		false,
	)

	issues := checker.GetIssues()

	if len(issues) != 1 {
		t.Errorf("Expected 1 issue, got %d", len(issues))
	}

	if checker.HasIssues() == false {
		t.Error("Expected HasIssues() to be true")
	}
}

func TestLightenDarken(t *testing.T) {
	gray := RGB{128, 128, 128}

	lightened := lighten(gray, 0.5)
	if lightened.R <= gray.R {
		t.Error("Expected lighten to increase R value")
	}

	darkened := darken(gray, 0.5)
	if darkened.R >= gray.R {
		t.Error("Expected darken to decrease R value")
	}

	// Test bounds
	white := RGB{255, 255, 255}
	lightenedWhite := lighten(white, 0.5)
	if lightenedWhite.R != 255 {
		t.Error("Lightening white should stay at 255")
	}

	black := RGB{0, 0, 0}
	darkenedBlack := darken(black, 0.5)
	if darkenedBlack.R != 0 {
		t.Error("Darkening black should stay at 0")
	}
}

func TestToHex(t *testing.T) {
	tests := []struct {
		rgb      RGB
		expected string
	}{
		{RGB{255, 0, 0}, "#ff0000"},
		{RGB{0, 255, 0}, "#00ff00"},
		{RGB{0, 0, 255}, "#0000ff"},
		{RGB{255, 255, 255}, "#ffffff"},
		{RGB{0, 0, 0}, "#000000"},
	}

	for _, tt := range tests {
		result := tt.rgb.ToHex()
		if result != tt.expected {
			t.Errorf("Expected %s, got %s", tt.expected, result)
		}
	}
}

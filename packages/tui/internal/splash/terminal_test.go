package splash

import (
	"os"
	"testing"
)

func TestDetectColorMode(t *testing.T) {
	tests := []struct {
		name        string
		colorterm   string
		term        string
		noColor     string
		expected    ColorMode
		expectRange bool // Allow range check instead of exact match
	}{
		{
			name:      "Truecolor with COLORTERM",
			colorterm: "truecolor",
			term:      "xterm-256color",
			expected:  Truecolor,
		},
		{
			name:      "256 color terminal",
			colorterm: "",
			term:      "xterm-256color",
			expected:  Colors256,
		},
		{
			name:      "NO_COLOR with no COLORTERM",
			colorterm: "",
			term:      "xterm",
			noColor:   "1",
			expected:  Colors16,
		},
		{
			name:      "Basic terminal",
			colorterm: "",
			term:      "xterm",
			expected:  Colors16,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original env vars
			origColorterm := os.Getenv("COLORTERM")
			origTerm := os.Getenv("TERM")
			origNoColor := os.Getenv("NO_COLOR")

			// Set test env vars
			os.Setenv("COLORTERM", tt.colorterm)
			os.Setenv("TERM", tt.term)
			os.Setenv("NO_COLOR", tt.noColor)

			// Run test
			result := DetectColorMode()

			// Check result
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}

			// Restore env vars
			os.Setenv("COLORTERM", origColorterm)
			os.Setenv("TERM", origTerm)
			os.Setenv("NO_COLOR", origNoColor)
		})
	}
}

func TestSupportsUnicode(t *testing.T) {
	// Save original env
	origLang := os.Getenv("LANG")
	defer os.Setenv("LANG", origLang)

	tests := []struct {
		name     string
		lang     string
		expected bool
	}{
		{
			name:     "UTF-8 locale",
			lang:     "en_US.UTF-8",
			expected: true,
		},
		{
			name:     "Non-UTF-8 locale",
			lang:     "C",
			expected: true, // Still defaults to true
		},
		{
			name:     "Empty LANG",
			lang:     "",
			expected: true, // Default to true
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("LANG", tt.lang)
			result := SupportsUnicode()

			if result != tt.expected {
				t.Errorf("Expected %v for LANG=%s, got %v", tt.expected, tt.lang, result)
			}
		})
	}
}

func TestDetectTerminalCapabilities(t *testing.T) {
	// Just ensure it doesn't panic and returns reasonable values
	caps := DetectTerminalCapabilities()

	// Width and height should be > 0 (or at least defaults)
	if caps.Width <= 0 {
		t.Error("Width should be > 0")
	}

	if caps.Height <= 0 {
		t.Error("Height should be > 0")
	}

	// Colors should be valid enum
	validColors := map[ColorMode]bool{
		Truecolor: true,
		Colors256: true,
		Colors16:  true,
	}
	if !validColors[caps.Colors] {
		t.Errorf("Invalid color mode: %v", caps.Colors)
	}

	// Performance should be valid
	validPerf := map[string]bool{
		"fast":   true,
		"medium": true,
		"slow":   true,
	}
	if !validPerf[caps.Performance] {
		t.Errorf("Invalid performance: %s", caps.Performance)
	}
}

func TestTerminalCapabilities_ShouldSkipSplash(t *testing.T) {
	tests := []struct {
		name     string
		caps     TerminalCapabilities
		expected bool
	}{
		{
			name: "Normal terminal",
			caps: TerminalCapabilities{
				Width:    80,
				Height:   24,
				Colors:   Truecolor,
				Unicode:  true,
				TooSmall: false,
			},
			expected: false,
		},
		{
			name: "Too small terminal",
			caps: TerminalCapabilities{
				Width:    40,
				Height:   10,
				Colors:   Colors16,
				Unicode:  true,
				TooSmall: true,
			},
			expected: true,
		},
		{
			name: "Extremely small terminal",
			caps: TerminalCapabilities{
				Width:    30,
				Height:   10,
				Colors:   Colors16,
				Unicode:  false,
				TooSmall: true,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.caps.ShouldSkipSplash()
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestColorModeString(t *testing.T) {
	tests := []struct {
		mode     ColorMode
		expected string
	}{
		{Truecolor, "truecolor"},
		{Colors256, "256-color"},
		{Colors16, "16-color"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.mode.String()
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestEstimatePerformance(t *testing.T) {
	// Just ensure it doesn't panic and returns valid value
	perf := EstimatePerformance()

	validPerf := map[string]bool{
		"fast":   true,
		"medium": true,
		"slow":   true,
	}

	if !validPerf[perf] {
		t.Errorf("Invalid performance estimate: %s", perf)
	}
}

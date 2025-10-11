package splash

import (
	"strings"
	"testing"
)

func TestNewTextOnlySplash(t *testing.T) {
	splash := NewTextOnlySplash(80, 24)

	if splash.width != 80 {
		t.Errorf("Expected width 80, got %d", splash.width)
	}

	if splash.height != 24 {
		t.Errorf("Expected height 24, got %d", splash.height)
	}
}

func TestTextOnlySplashRender(t *testing.T) {
	splash := NewTextOnlySplash(80, 24)
	output := splash.Render()

	// Should contain key elements
	if !strings.Contains(output, "RYCODE NEURAL CORTEX") {
		t.Error("Output should contain 'RYCODE NEURAL CORTEX'")
	}

	if !strings.Contains(output, "Claude") {
		t.Error("Output should contain 'Claude'")
	}

	if !strings.Contains(output, "SIX MINDS") {
		t.Error("Output should contain 'SIX MINDS'")
	}

	if !strings.Contains(output, "Press any key") {
		t.Error("Output should contain 'Press any key'")
	}

	// Should have multiple lines
	lines := strings.Split(output, "\n")
	if len(lines) < 10 {
		t.Errorf("Expected at least 10 lines, got %d", len(lines))
	}
}

func TestCenterText(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		width    int
		expected string
	}{
		{
			name:     "Center short text",
			text:     "Hello",
			width:    15,
			expected: "     Hello",
		},
		{
			name:     "Text too long",
			text:     "Hello World This Is Long",
			width:    10,
			expected: "Hello World This Is Long",
		},
		{
			name:     "Exact fit",
			text:     "Hello",
			width:    5,
			expected: "Hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := centerText(tt.text, tt.width)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestStripANSI(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Plain text",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "ANSI colored text",
			input:    "\033[38;2;255;0;0mRed\033[0m",
			expected: "Red",
		},
		{
			name:     "Multiple ANSI codes",
			input:    "\033[1mBold\033[0m and \033[4mUnderline\033[0m",
			expected: "Bold and Underline",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stripANSI(tt.input)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestNewSimplifiedSplash(t *testing.T) {
	caps := TerminalCapabilities{
		Width:    60,
		Height:   20,
		Colors:   Colors16,
		Unicode:  false,
		TooSmall: true,
	}

	model := NewSimplifiedSplash(caps)

	// Should be in simplified mode (act 4)
	if model.act != 4 {
		t.Errorf("Expected act 4 (simplified), got %d", model.act)
	}

	// Should not show skip hint for too-small terminals
	if model.skipHint {
		t.Error("Skip hint should be false for too-small terminals")
	}
}

func TestRenderSimplified(t *testing.T) {
	model := NewSimplifiedSplash(TerminalCapabilities{
		Width:    60,
		Height:   20,
		Colors:   Colors16,
		Unicode:  false,
		TooSmall: true,
	})

	caps := TerminalCapabilities{
		Width:    60,
		Height:   20,
		Colors:   Colors16,
		Unicode:  false,
		TooSmall: true,
	}

	output := model.RenderSimplified(caps)

	// Should contain key elements
	if !strings.Contains(output, "RYCODE NEURAL CORTEX") {
		t.Error("Simplified output should contain 'RYCODE NEURAL CORTEX'")
	}

	if !strings.Contains(output, "Claude") {
		t.Error("Simplified output should contain 'Claude'")
	}
}

func TestShouldUseSimplifiedSplash(t *testing.T) {
	// This function relies on terminal detection, so we just ensure it doesn't panic
	result := ShouldUseSimplifiedSplash()

	// Result should be a boolean (either true or false)
	if result != true && result != false {
		t.Error("Should return a boolean value")
	}
}

func TestRenderStaticCloser(t *testing.T) {
	output := RenderStaticCloser(80, 24)

	// Should contain key elements
	if !strings.Contains(output, "CORTEX ACTIVE") {
		t.Error("Static closer should contain 'CORTEX ACTIVE'")
	}

	if !strings.Contains(output, "Six minds") {
		t.Error("Static closer should contain 'Six minds'")
	}

	// Should have box drawing characters
	if !strings.Contains(output, "╔") || !strings.Contains(output, "╗") {
		t.Error("Static closer should have box drawing characters")
	}
}

func TestCenterTextWithANSI(t *testing.T) {
	// Test centering text that contains ANSI codes
	cyan := RGB{0, 255, 255}
	coloredText := Colorize("RYCODE", cyan)

	result := centerText(coloredText, 20)

	// The visible length should be centered properly
	// RYCODE is 6 characters, so padding should be (20-6)/2 = 7 spaces
	if !strings.HasPrefix(result, "       \033") { // 7 spaces + ANSI start
		t.Errorf("Colored text not centered properly: %q", result)
	}
}

func TestTextOnlySplashSmallTerminal(t *testing.T) {
	// Test with very small terminal
	splash := NewTextOnlySplash(40, 12)
	output := splash.Render()

	// Should still render without panic
	if output == "" {
		t.Error("Should render something even for small terminal")
	}

	// Should still contain essential elements
	if !strings.Contains(output, "RYCODE") {
		t.Error("Should contain RYCODE even in small terminal")
	}
}

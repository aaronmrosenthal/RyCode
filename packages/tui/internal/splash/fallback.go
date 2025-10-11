package splash

import (
	"fmt"
	"strings"
)

// TextOnlySplash renders a simplified text-only splash for small terminals
type TextOnlySplash struct {
	width  int
	height int
}

// NewTextOnlySplash creates a text-only splash for limited terminals
func NewTextOnlySplash(width, height int) *TextOnlySplash {
	return &TextOnlySplash{
		width:  width,
		height: height,
	}
}

// Render returns the text-only splash screen
func (t *TextOnlySplash) Render() string {
	cyan := RGB{0, 255, 255}
	gold := RGB{255, 174, 0}
	green := RGB{10, 255, 10}

	var buf strings.Builder

	// Add vertical padding
	topPadding := (t.height - 15) / 2
	if topPadding < 0 {
		topPadding = 0
	}
	for i := 0; i < topPadding; i++ {
		buf.WriteRune('\n')
	}

	// Simplified banner
	buf.WriteString(centerText("═══════════════════════════════════", t.width))
	buf.WriteRune('\n')
	buf.WriteString(centerText(Colorize("RYCODE NEURAL CORTEX", cyan), t.width))
	buf.WriteRune('\n')
	buf.WriteString(centerText("═══════════════════════════════════", t.width))
	buf.WriteRune('\n')
	buf.WriteRune('\n')

	// AI models list (simplified)
	buf.WriteString(centerText(Colorize("🧩 Claude", green)+"  • Logical Reasoning", t.width))
	buf.WriteRune('\n')
	buf.WriteString(centerText(Colorize("⚙️  Gemini", green)+"  • System Architecture", t.width))
	buf.WriteRune('\n')
	buf.WriteString(centerText(Colorize("💻 Codex", green)+"   • Code Generation", t.width))
	buf.WriteRune('\n')
	buf.WriteString(centerText(Colorize("🔎 Qwen", green)+"   • Research Pipeline", t.width))
	buf.WriteRune('\n')
	buf.WriteString(centerText(Colorize("🤖 Grok", green)+"   • Humor & Chaos", t.width))
	buf.WriteRune('\n')
	buf.WriteString(centerText(Colorize("✅ GPT", green)+"    • Language Core", t.width))
	buf.WriteRune('\n')
	buf.WriteRune('\n')

	// Tagline
	buf.WriteString(centerText(Colorize("⚡ SIX MINDS. ONE COMMAND LINE.", gold), t.width))
	buf.WriteRune('\n')
	buf.WriteRune('\n')

	// Skip hint
	gray := RGB{100, 100, 100}
	buf.WriteString(centerText(Colorize("Press any key to continue...", gray), t.width))

	return buf.String()
}

// centerText centers a string within a given width
func centerText(text string, width int) string {
	// Calculate visible length (strip ANSI codes for counting)
	visibleLen := len(stripANSI(text))

	if visibleLen >= width {
		return text
	}

	padding := (width - visibleLen) / 2
	return strings.Repeat(" ", padding) + text
}

// stripANSI removes ANSI escape codes for length calculation
func stripANSI(text string) string {
	// Simple implementation: remove everything between \033[ and m
	result := ""
	inEscape := false
	for i := 0; i < len(text); i++ {
		if i < len(text)-1 && text[i] == '\033' && text[i+1] == '[' {
			inEscape = true
			i++ // Skip the '['
			continue
		}
		if inEscape {
			if text[i] == 'm' {
				inEscape = false
			}
			continue
		}
		result += string(text[i])
	}
	return result
}

// NewSimplifiedSplash creates a model with simplified splash for small terminals
func NewSimplifiedSplash(caps TerminalCapabilities) Model {
	m := New()
	m.act = 4 // Special "simplified" mode

	// Use text-only renderer
	if caps.TooSmall {
		m.skipHint = false
	}

	return m
}

// RenderSimplified renders the simplified splash
func (m Model) RenderSimplified(caps TerminalCapabilities) string {
	textSplash := NewTextOnlySplash(caps.Width, caps.Height)
	return textSplash.Render()
}

// ShouldUseSimplifiedSplash determines if simplified mode should be used
func ShouldUseSimplifiedSplash() bool {
	caps := DetectTerminalCapabilities()
	return caps.TooSmall || caps.ShouldSkipSplash()
}

// RenderStaticCloser renders a static version of the closer screen
func RenderStaticCloser(width, height int) string {
	cyan := RGB{0, 255, 255}
	gold := RGB{255, 174, 0}

	var buf strings.Builder

	// Vertical centering
	topPadding := (height - 8) / 2
	if topPadding < 0 {
		topPadding = 0
	}
	for i := 0; i < topPadding; i++ {
		buf.WriteRune('\n')
	}

	buf.WriteString(centerText("╔════════════════════════════════════╗", width))
	buf.WriteRune('\n')
	buf.WriteString(centerText("║                                    ║", width))
	buf.WriteRune('\n')
	buf.WriteString(centerText(fmt.Sprintf("║  %s  ║", Colorize("🌀 RYCODE CORTEX ACTIVE 🌀", cyan)), width))
	buf.WriteRune('\n')
	buf.WriteString(centerText("║                                    ║", width))
	buf.WriteRune('\n')
	buf.WriteString(centerText(fmt.Sprintf("║   %s  ║", Colorize("Six minds. One command line.", gold)), width))
	buf.WriteRune('\n')
	buf.WriteString(centerText("║                                    ║", width))
	buf.WriteRune('\n')
	buf.WriteString(centerText("╚════════════════════════════════════╝", width))

	return buf.String()
}

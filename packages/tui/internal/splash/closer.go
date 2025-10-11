package splash

import (
	"strings"
)

// Closer manages the final "closer" screen
type Closer struct {
	width  int
	height int
}

// NewCloser creates a new closer screen
func NewCloser(width, height int) *Closer {
	return &Closer{
		width:  width,
		height: height,
	}
}

// closerText is the final message shown to users
const closerText = `╔═══════════════════════════════════════════════════════════════════════╗
║                                                                       ║
║                  🌀 RYCODE NEURAL CORTEX ACTIVE 🌀                   ║
║                                                                       ║
║         "Every LLM fused. Every edge case covered.                   ║
║          You're not just coding—you're orchestrating                 ║
║          intelligence."                                              ║
║                                                                       ║
║                                                                       ║
║                   Press any key to begin...                          ║
║                                                                       ║
╚═══════════════════════════════════════════════════════════════════════╝`

// compactCloserText is shown on smaller terminals
const compactCloserText = `╔═══════════════════════════════════╗
║   🌀 RYCODE CORTEX ACTIVE 🌀     ║
║                                   ║
║   Six minds. One command line.    ║
║                                   ║
║   Press any key...                ║
╚═══════════════════════════════════╝`

// Render renders the closer screen
func (c *Closer) Render() string {
	// Choose version based on terminal size
	text := closerText
	if c.width < 80 || c.height < 24 {
		text = compactCloserText
	}

	lines := strings.Split(text, "\n")

	// Calculate vertical centering
	startY := (c.height - len(lines)) / 2
	if startY < 0 {
		startY = 0
	}

	var buf strings.Builder

	// Add top padding
	for i := 0; i < startY; i++ {
		buf.WriteRune('\n')
	}

	// Render centered lines
	cyan := RGB{0, 255, 255}
	for _, line := range lines {
		// Horizontal centering
		padding := (c.width - len(line)) / 2
		if padding > 0 {
			buf.WriteString(strings.Repeat(" ", padding))
		}

		// Colorize lines with emoji or "CORTEX"
		if strings.Contains(line, "🌀") || strings.Contains(line, "CORTEX") {
			buf.WriteString(Colorize(line, cyan))
		} else {
			buf.WriteString(line)
		}

		buf.WriteRune('\n')
	}

	return buf.String()
}

package splash

import (
	"fmt"
	"strings"
	"time"
)

// ModelInfo represents an AI model in the neural cortex
type ModelInfo struct {
	Name  string        // Model name (e.g., "Claude")
	Role  string        // Model's role (e.g., "Logical Reasoning")
	Icon  string        // Icon/emoji
	Color RGB           // Display color
	Delay time.Duration // Delay before showing
}

// AI models in the neural cortex
var models = []ModelInfo{
	{"Claude", "Logical Reasoning", "🧩", RGB{10, 255, 10}, 100 * time.Millisecond},
	{"Gemini", "System Architecture", "⚙️", RGB{10, 255, 10}, 100 * time.Millisecond},
	{"Codex", "Code Generation", "💻", RGB{10, 255, 10}, 100 * time.Millisecond},
	{"Qwen", "Research Pipeline", "🔎", RGB{10, 255, 10}, 100 * time.Millisecond},
	{"Grok", "Humor & Chaos Engine", "🤖", RGB{10, 255, 10}, 100 * time.Millisecond},
	{"GPT", "Language Core", "✅", RGB{10, 255, 10}, 100 * time.Millisecond},
}

// BootSequence manages the boot sequence animation
type BootSequence struct {
	frame      int // Current frame
	linesShown int // Number of lines currently shown
}

// NewBootSequence creates a new boot sequence
func NewBootSequence() *BootSequence {
	return &BootSequence{}
}

// Update updates the boot sequence state based on frame number
func (b *BootSequence) Update(frame int) {
	b.frame = frame
	// Show 1 line every 3 frames (100ms at 30 FPS)
	b.linesShown = frame / 3
	if b.linesShown > len(models) {
		b.linesShown = len(models)
	}
}

// Render renders the boot sequence
func (b *BootSequence) Render() string {
	var buf strings.Builder

	// Header
	cyan := RGB{0, 255, 255}
	buf.WriteString(Colorize("> [RYCODE NEURAL CORTEX v1.0.0]\n", cyan))
	buf.WriteString(">\n")

	// Model lines
	for i := 0; i < b.linesShown && i < len(models); i++ {
		model := models[i]

		// Tree structure
		prefix := "├─"
		if i == len(models)-1 {
			prefix = "└─"
		}

		line := fmt.Sprintf("> %s %s ▸ %s: ONLINE %s\n",
			prefix, model.Name, model.Role, model.Icon)

		buf.WriteString(Colorize(line, model.Color))
	}

	// Final message after all models loaded
	if b.linesShown >= len(models) {
		buf.WriteString(">\n")
		gold := RGB{255, 174, 0}
		buf.WriteString(Colorize("> ⚡ SIX MINDS. ONE COMMAND LINE.\n", gold))
	}

	return buf.String()
}

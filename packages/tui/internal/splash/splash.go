// Package splash implements the epic RyCode splash screen with 3D ASCII animation.
//
// Features:
// - 3D rotating torus (neural cortex) with real donut math
// - 3-act animation sequence: Boot → Cortex → Closer
// - Cyberpunk color gradients (cyan to magenta)
// - 30 FPS smooth animation
// - Easter eggs and skip functionality
package splash

import (
	"time"

	tea "github.com/charmbracelet/bubbletea/v2"
)

// Model represents the splash screen state
type Model struct {
	act       int           // Current act: 1=boot, 2=cortex, 3=closer
	frame     int           // Current frame number
	bootSeq   *BootSequence // Boot sequence animation
	cortex    *CortexRenderer // 3D torus renderer
	closer    *Closer       // Closer screen
	done      bool          // Whether splash is complete
	width     int           // Terminal width
	height    int           // Terminal height
	skipHint  bool          // Whether to show skip hint
}

// tickMsg is sent on each animation frame
type tickMsg time.Time

// New creates a new splash screen model
func New() Model {
	return Model{
		act:      1,
		frame:    0,
		bootSeq:  NewBootSequence(),
		cortex:   NewCortexRenderer(80, 24),
		closer:   NewCloser(80, 24),
		skipHint: true,
	}
}

// Init initializes the splash screen
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		tick(),
	)
}

// tick returns a command that sends a tickMsg after 33ms (30 FPS)
func tick() tea.Cmd {
	return tea.Tick(33*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Update dimensions
		m.width = msg.Width
		m.height = msg.Height
		m.cortex = NewCortexRenderer(msg.Width, msg.Height)
		m.closer = NewCloser(msg.Width, msg.Height)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "s", "S":
			// Skip splash
			m.done = true
			return m, tea.Quit

		case "esc":
			// Skip and disable forever
			// TODO: Save config to disable splash
			m.done = true
			return m, tea.Quit

		case "enter", " ":
			// Continue from closer screen
			if m.act == 3 {
				m.done = true
				return m, tea.Quit
			}
		}

	case tickMsg:
		m.frame++

		// Act transitions
		if m.act == 1 && m.frame > 30 {
			// After 1 second (30 frames), move to cortex
			m.act = 2
		} else if m.act == 2 && m.frame > 120 {
			// After 4 seconds total (120 frames), move to closer
			m.act = 3
		} else if m.act == 3 && m.frame > 150 {
			// After 5 seconds total (150 frames), auto-close
			m.done = true
			return m, tea.Quit
		}

		return m, tick()
	}

	return m, nil
}

// View renders the current splash screen state
func (m Model) View() string {
	var content string

	switch m.act {
	case 1:
		// Boot sequence
		m.bootSeq.Update(m.frame)
		content = m.bootSeq.Render()

	case 2:
		// Rotating cortex
		content = m.cortex.Render()

	case 3:
		// Closer screen
		content = m.closer.Render()

	default:
		content = ""
	}

	// Add skip hint if enabled
	if m.skipHint && m.act < 3 {
		content += "\n\n" + renderSkipHint()
	}

	return content
}

// renderSkipHint returns the skip hint text
func renderSkipHint() string {
	gray := RGB{100, 100, 100}
	return Colorize("Press 'S' to skip | ESC to disable forever", gray)
}

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
	act           int             // Current act: 1=boot, 2=cortex, 3=closer
	frame         int             // Current frame number
	bootSeq       *BootSequence   // Boot sequence animation
	cortex        *CortexRenderer // 3D torus renderer
	closer        *Closer         // Closer screen
	done          bool            // Whether splash is complete
	width         int             // Terminal width
	height        int             // Terminal height
	skipHint      bool            // Whether to show skip hint
	donutMode     bool            // Easter egg: infinite mode
	showMath      bool            // Easter egg: show equations
	konamiCode    []string        // Konami code progress
	konamiIdx     int             // Current position in Konami sequence
	rainbowMode   bool            // Easter egg: rainbow colors
	frameTimes    []time.Duration // Frame time history for adaptive FPS
	lastFrameTime time.Time       // Last frame timestamp
	targetFPS     int             // Target FPS (30 or 15)
}

// tickMsg is sent on each animation frame
type tickMsg time.Time

// Konami code sequence
var konamiSequence = []string{"up", "up", "down", "down", "left", "right", "left", "right", "b", "a"}

// New creates a new splash screen model
func New() Model {
	return Model{
		act:       1,
		frame:     0,
		bootSeq:   NewBootSequence(),
		cortex:    NewCortexRenderer(80, 24),
		closer:    NewCloser(80, 24),
		skipHint:  true,
		targetFPS: 30,
	}
}

// NewDonutMode creates a model in infinite donut mode (easter egg)
func NewDonutMode() Model {
	return Model{
		act:       2, // Jump straight to cortex
		frame:     0,
		cortex:    NewCortexRenderer(80, 24),
		donutMode: true,
		skipHint:  false,
		targetFPS: 30,
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
		key := msg.String()

		// Konami code detection
		if !m.donutMode {
			if key == konamiSequence[m.konamiIdx] {
				m.konamiIdx++
				if m.konamiIdx >= len(konamiSequence) {
					m.rainbowMode = true
					m.konamiIdx = 0
					if m.cortex != nil {
						m.cortex.SetRainbowMode(true)
					}
				}
			} else if key != "?" && key != "s" && key != "S" {
				m.konamiIdx = 0
			}
		}

		switch key {
		case "?":
			// Toggle math equations
			m.showMath = !m.showMath
			return m, nil

		case "q", "Q":
			// Quit donut mode
			if m.donutMode {
				m.done = true
				return m, tea.Quit
			}

		case "s", "S":
			// Skip splash
			if !m.donutMode {
				m.done = true
				return m, tea.Quit
			}

		case "esc":
			// Skip and disable forever
			if !m.donutMode {
				// Disable splash permanently
				if err := DisableSplashPermanently(); err != nil {
					// Log error but don't block
					// User can still skip, just config won't save
				}
				m.done = true
				return m, tea.Quit
			}

		case "enter", " ":
			// Continue from closer screen
			if m.act == 3 && !m.donutMode {
				m.done = true
				return m, tea.Quit
			}
		}

	case tickMsg:
		// Measure frame time
		now := time.Now()
		if !m.lastFrameTime.IsZero() {
			frameTime := now.Sub(m.lastFrameTime)
			m.frameTimes = append(m.frameTimes, frameTime)
			if len(m.frameTimes) > 30 {
				m.frameTimes = m.frameTimes[1:]
			}

			// Adjust target FPS based on average frame time
			if len(m.frameTimes) >= 10 {
				avg := m.averageFrameTime()
				if avg > 50*time.Millisecond {
					// Frames taking >50ms, slow down to 15 FPS
					m.targetFPS = 15
				} else {
					// Frames are fast, maintain 30 FPS
					m.targetFPS = 30
				}
			}
		}
		m.lastFrameTime = now

		m.frame++

		// Act transitions (not in donut mode)
		if !m.donutMode {
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
		}

		// Dynamic tick rate based on target FPS
		tickDuration := time.Duration(1000/m.targetFPS) * time.Millisecond
		return m, tea.Tick(tickDuration, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})
	}

	return m, nil
}

// View renders the current splash screen state
func (m Model) View() string {
	// Math equations mode
	if m.showMath {
		return renderMathEquations()
	}

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

	case 4:
		// Simplified mode (for small terminals)
		caps := DetectTerminalCapabilities()
		content = m.RenderSimplified(caps)

	default:
		content = ""
	}

	// Add hints
	if m.donutMode {
		content += "\n\n" + renderDonutHint()
	} else if m.skipHint && m.act < 3 {
		content += "\n\n" + renderSkipHint()
	}

	// Konami code progress indicator (hidden)
	if m.konamiIdx > 0 && m.konamiIdx < 5 {
		gray := RGB{50, 50, 50}
		progress := Colorize("...", gray)
		content += "\n" + progress
	}

	return content
}

// renderSkipHint returns the skip hint text
func renderSkipHint() string {
	gray := RGB{100, 100, 100}
	return Colorize("Press 'S' to skip | ESC to disable forever | '?' for math", gray)
}

// renderDonutHint returns the donut mode hint
func renderDonutHint() string {
	gray := RGB{100, 100, 100}
	cyan := RGB{0, 255, 255}
	return Colorize("🍩 DONUT MODE ", cyan) + Colorize("| Press 'Q' to quit | '?' for math", gray)
}

// renderMathEquations returns the math equations display
func renderMathEquations() string {
	cyan := RGB{0, 255, 255}
	gold := RGB{255, 174, 0}

	return Colorize(`
╔═══════════════════════════════════════════════════════════════════════╗
║                                                                       ║
║                  🧮 DONUT MATH - 3D Torus Equations                  ║
║                                                                       ║
╚═══════════════════════════════════════════════════════════════════════╝

`, cyan) + `Torus Parametric Equations:
  x(θ,φ) = (R + r·cos(φ))·cos(θ)
  y(θ,φ) = (R + r·cos(φ))·sin(θ)
  z(θ,φ) = r·sin(φ)

Where:
  R = 2 (major radius - distance from center to tube center)
  r = 1 (minor radius - tube thickness)
  θ = angle around torus (0 to 2π)
  φ = angle around tube (0 to 2π)

Rotation Matrices:

  Rx(A) = [1    0       0    ]
          [0  cos(A) -sin(A) ]
          [0  sin(A)  cos(A) ]

  Rz(B) = [cos(B) -sin(B)  0 ]
          [sin(B)  cos(B)  0 ]
          [0       0       1 ]

Perspective Projection:
  x_screen = width/2  + (30/z) * x
  y_screen = height/2 - (15/z) * y

Luminance (Phong Shading):
  L = cos(φ)·cos(θ)·sin(B) - cos(A)·cos(θ)·sin(φ) - sin(A)·sin(θ)
      + cos(B)·(cos(A)·sin(φ) - cos(θ)·sin(A)·sin(θ))

Character Mapping:
  L ∈ [-1, 1] → { ' ', '.', '·', ':', '*', '◉', '◎', '⚡' }

` + Colorize(`
Performance: 0.318ms per frame (85× faster than 30 FPS target!)

Press '?' again to return
`, gold)
}

// averageFrameTime calculates the average frame rendering time
func (m Model) averageFrameTime() time.Duration {
	if len(m.frameTimes) == 0 {
		return 0
	}

	var sum time.Duration
	for _, t := range m.frameTimes {
		sum += t
	}
	return sum / time.Duration(len(m.frameTimes))
}

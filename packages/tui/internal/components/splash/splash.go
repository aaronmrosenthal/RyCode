package splash

import (
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/lipgloss/v2/compat"
	"github.com/aaronmrosenthal/rycode/internal/styles"
	"github.com/aaronmrosenthal/rycode/internal/theme"
)

// SplashFinishedMsg is sent when the splash animation completes
type SplashFinishedMsg struct{}

// tickMsg is sent for animation updates
type tickMsg time.Time

const (
	splashDuration = 4500 * time.Millisecond // Extended for better viewing
	tickInterval   = 50 * time.Millisecond
	matrixChars    = "ﾊﾐﾋｰｳｼﾅﾓﾆｻﾜﾂｵﾘｱﾎﾃﾏｹﾒｴｶｷﾑﾕﾗｾﾈｽﾀﾇﾍ01"
)

type Model struct {
	width, height int
	startTime     time.Time
	rainColumns   []rainColumn
	logoVisible   bool
	fadeProgress  float64
}

type rainColumn struct {
	x        int
	y        int
	speed    int
	chars    []rune
	length   int
	brightness []float64
}

func New(width, height int) Model {
	// Initialize Matrix rain columns
	numColumns := width / 2
	columns := make([]rainColumn, numColumns)

	for i := range columns {
		columns[i] = rainColumn{
			x:      i * 2,
			y:      rand.Intn(height) - height,
			speed:  rand.Intn(3) + 1,
			length: rand.Intn(15) + 10,
			chars:  make([]rune, 0),
		}

		// Generate random Matrix characters
		for j := 0; j < columns[i].length; j++ {
			columns[i].chars = append(columns[i].chars, rune(matrixChars[rand.Intn(len(matrixChars))]))
			columns[i].brightness = append(columns[i].brightness, float64(columns[i].length-j)/float64(columns[i].length))
		}
	}

	return Model{
		width:       width,
		height:      height,
		rainColumns: columns,
		logoVisible: false,
		fadeProgress: 0,
	}
}

func (m Model) Init() tea.Cmd {
	m.startTime = time.Now()
	return tea.Batch(
		tickCmd(),
	)
}

func tickCmd() tea.Cmd {
	return tea.Tick(tickInterval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		elapsed := time.Since(m.startTime)

		// Check if animation should finish
		if elapsed >= splashDuration {
			return m, func() tea.Msg { return SplashFinishedMsg{} }
		}

		// Update fade progress with improved timing
		// Phase 1 (0-20%): Matrix only, no logo
		// Phase 2 (20-40%): Logo fades in with Matrix
		// Phase 3 (40-80%): Full visibility - logo and Matrix together
		// Phase 4 (80-100%): Fade out to chat view
		progress := float64(elapsed) / float64(splashDuration)

		if progress < 0.2 {
			// Matrix rain only
			m.fadeProgress = progress / 0.2
			m.logoVisible = false
		} else if progress < 0.4 {
			// Logo fading in with Matrix
			m.fadeProgress = 1.0
			m.logoVisible = true
		} else if progress < 0.8 {
			// Full visibility - enjoy the view
			m.fadeProgress = 1.0
			m.logoVisible = true
		} else {
			// Fade out
			m.fadeProgress = 1.0 - ((progress - 0.8) / 0.2)
			m.logoVisible = true
		}

		// Update Matrix rain
		for i := range m.rainColumns {
			m.rainColumns[i].y += m.rainColumns[i].speed

			// Reset column if it's off screen
			if m.rainColumns[i].y > m.height+m.rainColumns[i].length {
				m.rainColumns[i].y = -m.rainColumns[i].length
				m.rainColumns[i].x = rand.Intn(m.width/2) * 2

				// Regenerate characters
				for j := range m.rainColumns[i].chars {
					m.rainColumns[i].chars[j] = rune(matrixChars[rand.Intn(len(matrixChars))])
				}
			}
		}

		return m, tickCmd()

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m Model) View() string {
	t := theme.CurrentTheme()

	// Create empty canvas
	canvas := make([][]rune, m.height)
	colors := make([][]string, m.height)
	for i := range canvas {
		canvas[i] = make([]rune, m.width)
		colors[i] = make([]string, m.width)
		for j := range canvas[i] {
			canvas[i][j] = ' '
			colors[i][j] = "#000000"
		}
	}

	// Draw Matrix rain - EPIC neon green and cyan like the actual terminal
	// Match the terminal's blue prompt and green RyCode branding
	brightCyan := "#00FFAA"    // Brightest - cyan-green (matches terminal prompt)
	neonGreen := "#00FF00"     // Pure Matrix green
	mediumGreen := "#00CC88"   // Medium cyan-green (matches RyCode logo)
	darkGreen := "#008866"     // Darker green for depth

	for _, col := range m.rainColumns {
		for i, char := range col.chars {
			y := col.y + i
			if y >= 0 && y < m.height && col.x < m.width {
				canvas[y][col.x] = char

				// Brightest at the head, dimmer towards tail - EPIC color cascade
				brightness := col.brightness[i] * m.fadeProgress
				if brightness > 0.8 {
					// Head of rain: bright cyan (like terminal prompt)
					colors[y][col.x] = brightCyan
				} else if brightness > 0.6 {
					// Upper section: pure Matrix green
					colors[y][col.x] = neonGreen
				} else if brightness > 0.3 {
					// Middle section: medium green (RyCode brand color)
					colors[y][col.x] = mediumGreen
				} else {
					// Tail: dark green for depth
					colors[y][col.x] = darkGreen
				}
			}
		}
	}

	// Ry-Code ASCII art (toolkit-cli style - bright and readable)
	logo := []string{
		"",
		"  ________               _________     _________     ",
		"  ___  __ \\____  __      __  ____/___________  /____ ",
		"  __  /_/ /_  / / /_______  /    _  __ \\  __  /_  _ \\",
		"  _  _, _/_  /_/ /_/_____/ /___  / /_/ / /_/ / /  __/",
		"  /_/ |_| _\\__, /        \\____/  \\____/\\__,_/  \\___/ ",
		"          /____/                                     ",
		"",
	}

	tagline := "> Where Code Writes Itself"

	// Calculate center position for logo
	logoStartY := (m.height - len(logo) - 2) / 2

	// Overlay logo if visible
	if m.logoVisible && logoStartY >= 0 {
		for i, line := range logo {
			y := logoStartY + i
			if y >= 0 && y < m.height {
				// Center the line
				startX := (m.width - len(line)) / 2
				if startX < 0 {
					startX = 0
				}

				for j, char := range line {
					x := startX + j
					if x >= 0 && x < m.width {
						if char != ' ' {
							canvas[y][x] = char
							// Logo glows in Matrix green matching terminal branding
							colors[y][x] = brightCyan
						}
					}
				}
			}
		}

		// Add tagline below logo
		taglineY := logoStartY + len(logo)
		if taglineY < m.height {
			startX := (m.width - len(tagline)) / 2
			if startX < 0 {
				startX = 0
			}

			for j, char := range tagline {
				x := startX + j
				if x >= 0 && x < m.width {
					canvas[taglineY][x] = char
					// Tagline in medium green to match RyCode branding
					colors[taglineY][x] = mediumGreen
				}
			}
		}
	}

	// Render canvas with colors
	var output strings.Builder
	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			char := canvas[y][x]
			colorStr := colors[y][x]

			// Create adaptive color from string
			adaptiveColor := compat.AdaptiveColor{
				Dark:  lipgloss.Color(colorStr),
				Light: lipgloss.Color(colorStr),
			}

			style := styles.NewStyle().
				Foreground(adaptiveColor).
				Background(t.Background())

			output.WriteString(style.Render(string(char)))
		}
		if y < m.height-1 {
			output.WriteString("\n")
		}
	}

	return output.String()
}

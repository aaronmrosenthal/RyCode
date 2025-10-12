package splash

import (
	"math/rand"
	"strings"
	"time"
)

// Matrix rain animation constants
const (
	// Stream configuration
	streamDensityPercent = 60  // 60% of terminal width has active streams
	minStreamLength      = 5   // Minimum characters per stream
	maxStreamLength      = 20  // Maximum characters per stream
	minStreamSpeed       = 0.3 // Minimum fall speed (chars per frame)
	maxStreamSpeed       = 1.0 // Maximum fall speed (chars per frame)
	minStreamAge         = 60  // Minimum frames before respawn
	maxStreamAge         = 180 // Maximum frames before respawn

	// Animation timing
	logoFadeFrames      = 90  // Frames for full logo fade-in (3s at 30 FPS)
	logoRevealThreshold = 0.5 // When to start revealing logo (50% fade progress)
	charMutationChance  = 0.1 // Probability of character mutation per frame

	// Intensity thresholds for gradient
	intensityHeadMin   = 0.8 // Stream head (bright white)
	intensityBrightMin = 0.5 // Bright green section
	intensityMidMin    = 0.3 // Standard green section
)

// Matrix character set (Katakana, Latin, numbers, symbols)
var matrixChars = []rune{
	// Katakana
	'ア', 'イ', 'ウ', 'エ', 'オ', 'カ', 'キ', 'ク', 'ケ', 'コ',
	'サ', 'シ', 'ス', 'セ', 'ソ', 'タ', 'チ', 'ツ', 'テ', 'ト',
	'ナ', 'ニ', 'ヌ', 'ネ', 'ノ', 'ハ', 'ヒ', 'フ', 'ヘ', 'ホ',
	'マ', 'ミ', 'ム', 'メ', 'モ', 'ヤ', 'ユ', 'ヨ', 'ラ', 'リ',
	'ル', 'レ', 'ロ', 'ワ', 'ヲ', 'ン',
	// Numbers and symbols
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	':', '.', '=', '*', '+', '-', '<', '>', '¦', '|',
	'"', '\'', '^', '~', '`',
}

// MatrixRain renders falling Matrix-style character streams
type MatrixRain struct {
	width           int          // Screen width
	height          int          // Screen height
	streams         []RainStream // Array of character streams
	logo            string       // ASCII logo to reveal
	logoLines       []string     // Logo split into lines
	logoStartX      int          // Horizontal logo position
	logoStartY      int          // Vertical logo position
	logoMaxWidth    int          // Maximum logo line width
	frame           int          // Current frame number
	screenBuffer    [][]rune     // Reused screen buffer
	intensityBuffer [][]float64  // Reused intensity buffer
	logoMask        [][]bool     // Pre-calculated logo reveal mask
	rng             *rand.Rand   // Seeded random generator
}

// RainStream represents a single falling stream of characters
type RainStream struct {
	column int     // X position
	y      float64 // Y position (float for smooth animation)
	speed  float64 // Fall speed (characters per frame)
	length int     // Stream length
	chars  []rune  // Characters in the stream
	age    int     // How long this stream has existed
	maxAge int     // When to respawn
}

// NewMatrixRain creates a new Matrix rain renderer
func NewMatrixRain(width, height int, logo string) *MatrixRain {
	// Create seeded random generator for unique patterns each run
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Split logo into lines
	logoLines := strings.Split(logo, "\n")

	// Calculate logo dimensions and centering
	logoMaxWidth := 0
	for _, line := range logoLines {
		if len(line) > logoMaxWidth {
			logoMaxWidth = len(line)
		}
	}

	logoStartX := max(0, (width-logoMaxWidth)/2)
	logoStartY := max(0, (height-len(logoLines))/2)

	// Create initial streams (60% density)
	numStreams := width * streamDensityPercent / 100
	streams := make([]RainStream, numStreams)

	for i := range streams {
		streams[i] = createStreamWithRng(width, height, rng)
	}

	// Pre-allocate buffers (reused every frame)
	screenBuffer := make([][]rune, height)
	intensityBuffer := make([][]float64, height)
	logoMask := make([][]bool, height)

	for i := range screenBuffer {
		screenBuffer[i] = make([]rune, width)
		intensityBuffer[i] = make([]float64, width)
		logoMask[i] = make([]bool, width)
	}

	return &MatrixRain{
		width:           width,
		height:          height,
		streams:         streams,
		logo:            logo,
		logoLines:       logoLines,
		logoStartX:      logoStartX,
		logoStartY:      logoStartY,
		logoMaxWidth:    logoMaxWidth,
		frame:           0,
		screenBuffer:    screenBuffer,
		intensityBuffer: intensityBuffer,
		logoMask:        logoMask,
		rng:             rng,
	}
}

// createStreamWithRng creates a new rain stream with random properties
func createStreamWithRng(width, height int, rng *rand.Rand) RainStream {
	length := minStreamLength + rng.Intn(maxStreamLength-minStreamLength+1)
	chars := make([]rune, length)
	for j := range chars {
		chars[j] = matrixChars[rng.Intn(len(matrixChars))]
	}

	return RainStream{
		column: rng.Intn(width),
		y:      float64(-rng.Intn(height)), // Start above screen
		speed:  minStreamSpeed + rng.Float64()*(maxStreamSpeed-minStreamSpeed),
		length: length,
		chars:  chars,
		age:    0,
		maxAge: minStreamAge + rng.Intn(maxStreamAge-minStreamAge+1),
	}
}

// Update advances the animation by one frame
func (m *MatrixRain) Update() {
	m.frame++

	// Update all streams
	for i := range m.streams {
		stream := &m.streams[i]
		stream.y += stream.speed
		stream.age++

		// Mutate characters occasionally (Matrix effect)
		if m.rng.Float64() < charMutationChance {
			idx := m.rng.Intn(len(stream.chars))
			stream.chars[idx] = matrixChars[m.rng.Intn(len(matrixChars))]
		}

		// Respawn if stream is too old or off screen
		if stream.age >= stream.maxAge || stream.y > float64(m.height+stream.length) {
			m.streams[i] = createStreamWithRng(m.width, m.height, m.rng)
		}
	}

	// Pre-calculate logo reveal mask (deterministic for this frame)
	m.updateLogoMask()
}

// updateLogoMask pre-calculates which logo pixels should be visible this frame
func (m *MatrixRain) updateLogoMask() {
	fadeProgress := float64(m.frame) / logoFadeFrames
	if fadeProgress > 1.0 {
		fadeProgress = 1.0
	}

	// Clear mask
	for y := range m.logoMask {
		for x := range m.logoMask[y] {
			m.logoMask[y][x] = false
		}
	}

	// Calculate which logo pixels are revealed
	if fadeProgress > logoRevealThreshold {
		for ly := 0; ly < len(m.logoLines); ly++ {
			y := m.logoStartY + ly
			if y < 0 || y >= m.height {
				continue
			}

			logoLine := m.logoLines[ly]
			for lx := 0; lx < len(logoLine); lx++ {
				x := m.logoStartX + lx
				if x < 0 || x >= m.width {
					continue
				}

				logoChar := rune(logoLine[lx])
				if logoChar == ' ' || logoChar == 0 {
					continue
				}

				// Deterministically reveal based on position and fade progress
				// Use position-based hash for consistent but varied reveal
				posHash := float64((x*7 + y*13) % 100) / 100.0
				revealThreshold := (fadeProgress - logoRevealThreshold) / (1.0 - logoRevealThreshold)

				if posHash < revealThreshold {
					m.logoMask[y][x] = true
				}
			}
		}
	}
}

// Render renders the current frame with logo overlay
func (m *MatrixRain) Render() string {
	// Clear buffers (reuse existing allocations)
	for i := range m.screenBuffer {
		for j := range m.screenBuffer[i] {
			m.screenBuffer[i][j] = ' '
			m.intensityBuffer[i][j] = 0.0
		}
	}

	// Draw streams with intensity tracking
	for _, stream := range m.streams {
		for j := 0; j < stream.length; j++ {
			y := int(stream.y) - j
			if y < 0 || y >= m.height || stream.column >= m.width {
				continue
			}

			m.screenBuffer[y][stream.column] = stream.chars[j%len(stream.chars)]

			// Calculate intensity: 1.0 at head, fading to 0.2 at tail
			distFromHead := float64(j) / float64(stream.length)
			newIntensity := 1.0 - (distFromHead * 0.8)

			// Only update if this intensity is brighter (handles overlapping streams)
			if newIntensity > m.intensityBuffer[y][stream.column] {
				m.intensityBuffer[y][stream.column] = newIntensity
			}
		}
	}

	// Build output string
	var result strings.Builder
	// Pre-allocate approximate capacity (width * height * 20 bytes per colored char)
	result.Grow(m.width * m.height * 20)

	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			char := m.screenBuffer[y][x]
			charIntensity := m.intensityBuffer[y][x]

			// Check if logo should be shown here
			if m.logoMask[y][x] {
				// Logo revealed at this position
				logoY := y - m.logoStartY
				logoX := x - m.logoStartX

				if logoY >= 0 && logoY < len(m.logoLines) {
					logoLine := m.logoLines[logoY]
					if logoX >= 0 && logoX < len(logoLine) {
						logoChar := rune(logoLine[logoX])
						if logoChar != ' ' && logoChar != 0 {
							result.WriteString(colorizeLogoChar(logoChar))
							continue
						}
					}
				}
			}

			// Show rain character
			if char != ' ' {
				result.WriteString(colorizeMatrixCharWithIntensity(char, charIntensity))
			} else {
				result.WriteRune(' ')
			}
		}
		if y < m.height-1 {
			result.WriteRune('\n')
		}
	}

	return result.String()
}

// colorizeLogoChar applies bright cyan color to logo characters
func colorizeLogoChar(char rune) string {
	// Bright cyan/green for the logo (Matrix style)
	cyan := RGB{0, 255, 170}
	return Colorize(string(char), cyan)
}

// colorizeMatrixCharWithIntensity applies gradient based on position in stream
func colorizeMatrixCharWithIntensity(char rune, intensity float64) string {
	// Head: bright white-green (255, 255, 255)
	// Mid: standard green (0, 255, 100)
	// Tail: dark green (0, 100, 40)

	if intensity > intensityHeadMin {
		// Stream head - bright white
		white := RGB{220, 255, 220}
		return Colorize(string(char), white)
	} else if intensity > intensityBrightMin {
		// Upper section - bright green
		brightGreen := RGB{50, 255, 130}
		return Colorize(string(char), brightGreen)
	} else if intensity > intensityMidMin {
		// Middle section - standard green
		green := RGB{0, 255, 100}
		return Colorize(string(char), green)
	} else {
		// Tail - darker green
		darkGreen := RGB{0, uint8(100 + intensity*100), 40}
		return Colorize(string(char), darkGreen)
	}
}

// GetFrame returns the current frame number
func (m *MatrixRain) GetFrame() int {
	return m.frame
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

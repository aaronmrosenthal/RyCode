package responsive

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/aaronmrosenthal/rycode/internal/theme"
)

// VoiceState represents the current voice input state
type VoiceState string

const (
	VoiceIdle       VoiceState = "idle"
	VoiceListening  VoiceState = "listening"
	VoiceProcessing VoiceState = "processing"
	VoiceError      VoiceState = "error"
)

// VoiceInput manages voice input functionality
type VoiceInput struct {
	state         VoiceState
	transcript    string
	confidence    float64
	isRecording   bool
	startTime     time.Time
	duration      time.Duration
	errorMessage  string
	waveform      []int // Visual waveform
	enabled       bool
}

// NewVoiceInput creates a new voice input manager
func NewVoiceInput() *VoiceInput {
	return &VoiceInput{
		state:    VoiceIdle,
		waveform: make([]int, 20),
		enabled:  true,
	}
}

// VoiceStartMsg signals voice recording start
type VoiceStartMsg struct{}

// VoiceStopMsg signals voice recording stop
type VoiceStopMsg struct{}

// VoiceTranscriptMsg contains transcribed text
type VoiceTranscriptMsg struct {
	Text       string
	Confidence float64
}

// VoiceErrorMsg contains voice error
type VoiceErrorMsg struct {
	Error string
}

// Start begins voice recording
func (vi *VoiceInput) Start() tea.Cmd {
	if !vi.enabled {
		return nil
	}

	vi.state = VoiceListening
	vi.isRecording = true
	vi.startTime = time.Now()
	vi.transcript = ""
	vi.errorMessage = ""

	return tea.Batch(
		func() tea.Msg {
			return VoiceStartMsg{}
		},
		vi.animateWaveform(),
	)
}

// Stop ends voice recording
func (vi *VoiceInput) Stop() tea.Cmd {
	vi.state = VoiceProcessing
	vi.isRecording = false
	vi.duration = time.Since(vi.startTime)

	return func() tea.Msg {
		return VoiceStopMsg{}
	}
}

// Update handles voice input updates
func (vi *VoiceInput) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case VoiceStartMsg:
		return vi.simulateRecording()

	case VoiceTranscriptMsg:
		vi.state = VoiceIdle
		vi.transcript = msg.Text
		vi.confidence = msg.Confidence
		return nil

	case VoiceErrorMsg:
		vi.state = VoiceError
		vi.errorMessage = msg.Error
		vi.isRecording = false
		return nil

	case VoiceWaveformMsg:
		if vi.isRecording {
			vi.waveform = msg.Levels
			return vi.animateWaveform()
		}
	}

	return nil
}

// VoiceWaveformMsg updates visual waveform
type VoiceWaveformMsg struct {
	Levels []int
}

// animateWaveform creates animated waveform effect
func (vi *VoiceInput) animateWaveform() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		// Generate random waveform levels (in real app, would be from mic)
		levels := make([]int, 20)
		for i := range levels {
			// Simulate voice levels 0-5
			levels[i] = int(time.Now().UnixNano()/(int64(i+1)*1000000)) % 6
		}

		return VoiceWaveformMsg{Levels: levels}
	})
}

// simulateRecording simulates voice recording process
func (vi *VoiceInput) simulateRecording() tea.Cmd {
	// In real implementation, this would:
	// 1. Initialize Web Speech API or native speech recognition
	// 2. Start capturing microphone input
	// 3. Send audio to speech-to-text service
	// 4. Return transcribed text

	return tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
		// Simulate successful transcription
		return VoiceTranscriptMsg{
			Text:       "fix the bug in auth.go",
			Confidence: 0.92,
		}
	})
}

// GetTranscript returns the transcribed text
func (vi *VoiceInput) GetTranscript() string {
	return vi.transcript
}

// IsRecording returns whether currently recording
func (vi *VoiceInput) IsRecording() bool {
	return vi.isRecording
}

// GetState returns current state
func (vi *VoiceInput) GetState() VoiceState {
	return vi.state
}

// Reset resets voice input state
func (vi *VoiceInput) Reset() {
	vi.state = VoiceIdle
	vi.transcript = ""
	vi.confidence = 0
	vi.isRecording = false
	vi.errorMessage = ""
}

// Render renders voice input UI
func (vi *VoiceInput) Render(theme *theme.Theme, width int) string {
	switch vi.state {
	case VoiceListening:
		return vi.renderListening(theme, width)
	case VoiceProcessing:
		return vi.renderProcessing(theme, width)
	case VoiceError:
		return vi.renderError(theme, width)
	default:
		return ""
	}
}

// renderListening renders listening state
func (vi *VoiceInput) renderListening(theme *theme.Theme, width int) string {
	containerStyle := lipgloss.NewStyle().
		Width(width).
		Padding(2).
		Background(theme.BackgroundSecondary).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(theme.AccentPrimary)

	// Title
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.AccentPrimary).
		Bold(true).
		Width(width - 4).
		Align(lipgloss.Center)

	title := titleStyle.Render("🎤 Listening...")

	// Waveform visualization
	waveform := vi.renderWaveform(theme, width-4)

	// Duration
	duration := time.Since(vi.startTime)
	durationStyle := lipgloss.NewStyle().
		Foreground(theme.TextDim).
		Width(width - 4).
		Align(lipgloss.Center)

	durationText := durationStyle.Render(fmt.Sprintf("%.1fs", duration.Seconds()))

	// Hint
	hintStyle := lipgloss.NewStyle().
		Foreground(theme.TextDim).
		Width(width - 4).
		Align(lipgloss.Center).
		MarginTop(1)

	hint := hintStyle.Render("Press again to stop")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		waveform,
		"",
		durationText,
		hint,
	)

	return containerStyle.Render(content)
}

// renderWaveform renders audio waveform visualization
func (vi *VoiceInput) renderWaveform(theme *theme.Theme, width int) string {
	bars := []string{}

	for _, level := range vi.waveform {
		bar := vi.renderBar(level, theme)
		bars = append(bars, bar)
	}

	return lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Render(lipgloss.JoinHorizontal(lipgloss.Center, bars...))
}

// renderBar renders a single waveform bar
func (vi *VoiceInput) renderBar(level int, theme *theme.Theme) string {
	chars := []string{"▁", "▂", "▃", "▄", "▅", "▆"}

	if level >= len(chars) {
		level = len(chars) - 1
	}

	style := lipgloss.NewStyle().
		Foreground(theme.AccentPrimary)

	return style.Render(chars[level])
}

// renderProcessing renders processing state
func (vi *VoiceInput) renderProcessing(theme *theme.Theme, width int) string {
	style := lipgloss.NewStyle().
		Width(width).
		Padding(2).
		Background(theme.BackgroundSecondary).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(theme.Info).
		Align(lipgloss.Center)

	content := lipgloss.NewStyle().
		Foreground(theme.Info).
		Bold(true).
		Render("⏳ Processing...")

	return style.Render(content)
}

// renderError renders error state
func (vi *VoiceInput) renderError(theme *theme.Theme, width int) string {
	style := lipgloss.NewStyle().
		Width(width).
		Padding(2).
		Background(theme.BackgroundSecondary).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(theme.Error).
		Align(lipgloss.Center)

	title := lipgloss.NewStyle().
		Foreground(theme.Error).
		Bold(true).
		Render("❌ Error")

	message := lipgloss.NewStyle().
		Foreground(theme.TextDim).
		MarginTop(1).
		Render(vi.errorMessage)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		message,
	)

	return style.Render(content)
}

// VoiceQuickCommands provides voice-activated quick commands
type VoiceQuickCommands struct {
	commands map[string]string
}

// NewVoiceQuickCommands creates voice command mappings
func NewVoiceQuickCommands() *VoiceQuickCommands {
	return &VoiceQuickCommands{
		commands: map[string]string{
			// Debugging
			"debug":      "/debug",
			"fix":        "/fix",
			"find bug":   "/debug",
			"fix bug":    "/fix",

			// Testing
			"test":       "/test",
			"run tests":  "/test",
			"test this":  "/test",

			// Code review
			"review":     "/review",
			"check code": "/review",

			// Explain
			"explain":    "/explain",
			"what is":    "/explain",
			"how does":   "/explain",

			// Preview
			"preview":    "/preview",
			"show me":    "/preview",

			// Commit
			"commit":     "/commit",
			"save":       "/commit",

			// History
			"history":    "show history",
			"previous":   "show history",

			// AI switching
			"switch ai":     "switch ai",
			"use claude":    "switch to claude",
			"use codex":     "switch to codex",
			"use gemini":    "switch to gemini",

			// Navigation
			"next":       "next",
			"back":       "back",
			"scroll up":  "scroll up",
			"scroll down": "scroll down",
		},
	}
}

// ParseCommand converts voice transcript to command
func (vqc *VoiceQuickCommands) ParseCommand(transcript string) string {
	transcript = strings.ToLower(strings.TrimSpace(transcript))

	// Direct command match
	if cmd, exists := vqc.commands[transcript]; exists {
		return cmd
	}

	// Partial match
	for voice, cmd := range vqc.commands {
		if strings.Contains(transcript, voice) {
			return cmd
		}
	}

	// No match, return as-is (natural language query)
	return transcript
}

// GetSuggestions returns voice command suggestions
func (vqc *VoiceQuickCommands) GetSuggestions() []string {
	suggestions := []string{
		"Try: 'debug this'",
		"Try: 'run tests'",
		"Try: 'explain this code'",
		"Try: 'switch to Claude'",
		"Try: 'fix the bug'",
		"Or just ask naturally!",
	}

	return suggestions
}

// VoiceHelpOverlay renders voice help overlay
func VoiceHelpOverlay(theme *theme.Theme, width int) string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.AccentPrimary).
		Bold(true).
		Width(width - 4).
		Align(lipgloss.Center).
		MarginBottom(1)

	title := titleStyle.Render("🎤 Voice Commands")

	commands := []struct {
		category string
		examples []string
	}{
		{
			category: "🐛 Debugging",
			examples: []string{"debug this", "find bug", "fix bug"},
		},
		{
			category: "🧪 Testing",
			examples: []string{"run tests", "test this"},
		},
		{
			category: "👀 Review",
			examples: []string{"review code", "check this"},
		},
		{
			category: "💡 Explain",
			examples: []string{"explain this", "how does this work"},
		},
		{
			category: "🤖 AI Switch",
			examples: []string{"use Claude", "switch to Gemini"},
		},
	}

	sections := []string{}
	for _, cmd := range commands {
		categoryStyle := lipgloss.NewStyle().
			Foreground(theme.AccentSecondary).
			Bold(true)

		category := categoryStyle.Render(cmd.category)

		exampleStyle := lipgloss.NewStyle().
			Foreground(theme.TextDim).
			MarginLeft(2)

		examples := []string{}
		for _, ex := range cmd.examples {
			examples = append(examples, exampleStyle.Render("• "+ex))
		}

		section := lipgloss.JoinVertical(
			lipgloss.Left,
			category,
			lipgloss.JoinVertical(lipgloss.Left, examples...),
		)

		sections = append(sections, section)
	}

	hintStyle := lipgloss.NewStyle().
		Foreground(theme.Info).
		Width(width - 4).
		Align(lipgloss.Center).
		MarginTop(1)

	hint := hintStyle.Render("Or just speak naturally!")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		lipgloss.JoinVertical(lipgloss.Left, sections...),
		hint,
	)

	containerStyle := lipgloss.NewStyle().
		Width(width - 2).
		Padding(1, 2).
		Background(theme.BackgroundSecondary).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(theme.Border)

	return containerStyle.Render(content)
}

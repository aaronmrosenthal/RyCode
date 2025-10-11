package spinner

import (
	"time"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/aaronmrosenthal/rycode/internal/styles"
	"github.com/aaronmrosenthal/rycode/internal/theme"
)

// Spinner frames for different styles
var (
	// Dots - classic loading dots
	Dots = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

	// Line - rotating line
	Line = []string{"|", "/", "-", "\\"}

	// Bounce - bouncing ball
	Bounce = []string{"⠁", "⠂", "⠄", "⠂"}

	// Circle - circling dots
	Circle = []string{"◐", "◓", "◑", "◒"}

	// Square - rotating square
	Square = []string{"◰", "◳", "◲", "◱"}

	// Arrow - rotating arrow
	Arrow = []string{"←", "↖", "↑", "↗", "→", "↘", "↓", "↙"}

	// Pulse - pulsing circle
	Pulse = []string{"◯", "◉"}

	// Dots2 - three dots
	Dots2 = []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}

	// Dots3 - braille dots
	Dots3 = []string{"⠋", "⠙", "⠚", "⠞", "⠖", "⠦", "⠴", "⠲", "⠳", "⠓"}

	// Ellipsis - loading ellipsis
	Ellipsis = []string{"   ", ".  ", ".. ", "..."}

	// ProgressBar - progress bar animation
	ProgressBar = []string{
		"▱▱▱▱▱▱▱",
		"▰▱▱▱▱▱▱",
		"▰▰▱▱▱▱▱",
		"▰▰▰▱▱▱▱",
		"▰▰▰▰▱▱▱",
		"▰▰▰▰▰▱▱",
		"▰▰▰▰▰▰▱",
		"▰▰▰▰▰▰▰",
	}

	// Rainbow - colorful spinner (uses ANSI colors)
	Rainbow = []string{"●", "●", "●", "●", "●", "●", "●", "●"}
)

// TickMsg is sent on each spinner frame
type TickMsg time.Time

// Spinner represents a loading spinner
type Spinner struct {
	frames   []string
	frame    int
	message  string
	running  bool
	interval time.Duration
	style    styles.Style
}

// New creates a new spinner with the default style (Dots)
func New() *Spinner {
	t := theme.CurrentTheme()
	return &Spinner{
		frames:   Dots,
		frame:    0,
		running:  false,
		interval: 80 * time.Millisecond,
		style: styles.NewStyle().
			Foreground(t.Primary()),
	}
}

// WithFrames sets custom frames for the spinner
func (s *Spinner) WithFrames(frames []string) *Spinner {
	s.frames = frames
	return s
}

// WithMessage sets the spinner message
func (s *Spinner) WithMessage(message string) *Spinner {
	s.message = message
	return s
}

// WithInterval sets the frame interval
func (s *Spinner) WithInterval(interval time.Duration) *Spinner {
	s.interval = interval
	return s
}

// WithStyle sets the spinner style
func (s *Spinner) WithStyle(style styles.Style) *Spinner {
	s.style = style
	return s
}

// Start starts the spinner
func (s *Spinner) Start() {
	s.running = true
	s.frame = 0
}

// Stop stops the spinner
func (s *Spinner) Stop() {
	s.running = false
}

// IsRunning returns whether the spinner is running
func (s *Spinner) IsRunning() bool {
	return s.running
}

// Init implements tea.Model
func (s *Spinner) Init() tea.Cmd {
	return s.tick()
}

// Update implements tea.Model
func (s *Spinner) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case TickMsg:
		if s.running {
			s.frame = (s.frame + 1) % len(s.frames)
			return s, s.tick()
		}
	}
	return s, nil
}

// View implements tea.Model
func (s *Spinner) View() string {
	if !s.running {
		return ""
	}

	t := theme.CurrentTheme()
	frame := s.style.Render(s.frames[s.frame])

	if s.message != "" {
		messageStyle := styles.NewStyle().
			Foreground(t.Text()).
			MarginLeft(1)
		return frame + messageStyle.Render(s.message)
	}

	return frame
}

// tick generates a tick command
func (s *Spinner) tick() tea.Cmd {
	return tea.Tick(s.interval, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// LoadingState represents different loading states
type LoadingState int

const (
	LoadingStateIdle LoadingState = iota
	LoadingStateAuthenticating
	LoadingStateFetching
	LoadingStateProcessing
	LoadingStateComplete
	LoadingStateFailed
)

// LoadingSpinner is an enhanced spinner with state and progress
type LoadingSpinner struct {
	spinner  *Spinner
	state    LoadingState
	progress float64
	steps    []string
	current  int
}

// NewLoadingSpinner creates a new loading spinner with progress
func NewLoadingSpinner() *LoadingSpinner {
	return &LoadingSpinner{
		spinner: New(),
		state:   LoadingStateIdle,
	}
}

// WithSteps sets the loading steps to display
func (ls *LoadingSpinner) WithSteps(steps []string) *LoadingSpinner {
	ls.steps = steps
	ls.current = 0
	return ls
}

// Start starts the loading spinner
func (ls *LoadingSpinner) Start() {
	ls.spinner.Start()
	ls.state = LoadingStateAuthenticating
	ls.current = 0
}

// NextStep advances to the next step
func (ls *LoadingSpinner) NextStep() {
	if ls.current < len(ls.steps)-1 {
		ls.current++
	}
}

// Complete marks the loading as complete
func (ls *LoadingSpinner) Complete() {
	ls.state = LoadingStateComplete
	ls.spinner.Stop()
}

// Fail marks the loading as failed
func (ls *LoadingSpinner) Fail() {
	ls.state = LoadingStateFailed
	ls.spinner.Stop()
}

// Init implements tea.Model
func (ls *LoadingSpinner) Init() tea.Cmd {
	return ls.spinner.Init()
}

// Update implements tea.Model
func (ls *LoadingSpinner) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	model, cmd := ls.spinner.Update(msg)
	ls.spinner = model.(*Spinner)
	return ls, cmd
}

// View implements tea.Model
func (ls *LoadingSpinner) View() string {
	if ls.state == LoadingStateIdle {
		return ""
	}

	t := theme.CurrentTheme()

	// State icon
	var icon string
	var iconStyle styles.Style

	switch ls.state {
	case LoadingStateComplete:
		icon = "✓"
		iconStyle = styles.NewStyle().Foreground(t.Success())
	case LoadingStateFailed:
		icon = "✗"
		iconStyle = styles.NewStyle().Foreground(t.Error())
	default:
		icon = ls.spinner.View()
		iconStyle = styles.NewStyle()
	}

	// Current step
	var stepText string
	if len(ls.steps) > 0 && ls.current < len(ls.steps) {
		stepText = ls.steps[ls.current]
	}

	// Build output
	if stepText != "" {
		stepStyle := styles.NewStyle().
			Foreground(t.Text()).
			MarginLeft(1)
		return iconStyle.Render(icon) + stepStyle.Render(stepText)
	}

	return iconStyle.Render(icon)
}

// MultiStepLoading shows a multi-step loading process
type MultiStepLoading struct {
	steps    []LoadingStep
	current  int
	spinner  *Spinner
	complete bool
}

// LoadingStep represents a single loading step
type LoadingStep struct {
	Label    string
	State    StepState
	Duration time.Duration
}

// StepState represents the state of a loading step
type StepState int

const (
	StepStatePending StepState = iota
	StepStateRunning
	StepStateComplete
	StepStateFailed
)

// NewMultiStepLoading creates a new multi-step loading display
func NewMultiStepLoading(steps []string) *MultiStepLoading {
	loadingSteps := make([]LoadingStep, len(steps))
	for i, label := range steps {
		loadingSteps[i] = LoadingStep{
			Label: label,
			State: StepStatePending,
		}
	}

	return &MultiStepLoading{
		steps:   loadingSteps,
		current: 0,
		spinner: New(),
	}
}

// Start starts the multi-step loading
func (ml *MultiStepLoading) Start() {
	ml.spinner.Start()
	if len(ml.steps) > 0 {
		ml.steps[0].State = StepStateRunning
	}
}

// NextStep moves to the next step
func (ml *MultiStepLoading) NextStep() {
	if ml.current < len(ml.steps) {
		ml.steps[ml.current].State = StepStateComplete
		ml.current++
		if ml.current < len(ml.steps) {
			ml.steps[ml.current].State = StepStateRunning
		}
	}
}

// FailCurrentStep marks the current step as failed
func (ml *MultiStepLoading) FailCurrentStep() {
	if ml.current < len(ml.steps) {
		ml.steps[ml.current].State = StepStateFailed
	}
}

// Complete marks all steps as complete
func (ml *MultiStepLoading) Complete() {
	ml.complete = true
	ml.spinner.Stop()
	for i := range ml.steps {
		if ml.steps[i].State != StepStateFailed {
			ml.steps[i].State = StepStateComplete
		}
	}
}

// Init implements tea.Model
func (ml *MultiStepLoading) Init() tea.Cmd {
	return ml.spinner.Init()
}

// Update implements tea.Model
func (ml *MultiStepLoading) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	model, cmd := ml.spinner.Update(msg)
	ml.spinner = model.(*Spinner)
	return ml, cmd
}

// View implements tea.Model
func (ml *MultiStepLoading) View() string {
	t := theme.CurrentTheme()
	output := ""

	for i, step := range ml.steps {
		var icon string
		var iconStyle styles.Style
		var labelStyle styles.Style

		switch step.State {
		case StepStatePending:
			icon = "○"
			iconStyle = styles.NewStyle().Foreground(t.TextMuted())
			labelStyle = styles.NewStyle().Foreground(t.TextMuted())
		case StepStateRunning:
			icon = ml.spinner.frames[ml.spinner.frame]
			iconStyle = styles.NewStyle().Foreground(t.Primary())
			labelStyle = styles.NewStyle().Foreground(t.Text())
		case StepStateComplete:
			icon = "✓"
			iconStyle = styles.NewStyle().Foreground(t.Success())
			labelStyle = styles.NewStyle().Foreground(t.TextMuted())
		case StepStateFailed:
			icon = "✗"
			iconStyle = styles.NewStyle().Foreground(t.Error())
			labelStyle = styles.NewStyle().Foreground(t.Error())
		}

		line := iconStyle.Render(icon) + " " + labelStyle.Render(step.Label)
		if i > 0 {
			output += "\n"
		}
		output += line
	}

	return output
}

// IsComplete returns whether all steps are complete
func (ml *MultiStepLoading) IsComplete() bool {
	return ml.complete
}

// HasFailed returns whether any step has failed
func (ml *MultiStepLoading) HasFailed() bool {
	for _, step := range ml.steps {
		if step.State == StepStateFailed {
			return true
		}
	}
	return false
}

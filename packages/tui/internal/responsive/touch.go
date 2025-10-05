package responsive

import (
	"time"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/sst/opencode/internal/theme"
)

// TouchZone represents a touchable area in the UI
type TouchZone struct {
	ID       string
	X        int
	Y        int
	Width    int
	Height   int
	Action   func() tea.Cmd
	Enabled  bool
	MinSize  int // Minimum touch target size (44px iOS, 48px Android)
	Priority int // Higher priority zones are checked first
}

// TouchTarget represents a touch-optimized button/area
type TouchTarget struct {
	zone      TouchZone
	label     string
	icon      string
	focused   bool
	pressed   bool
	theme     *theme.Theme
	touchSize int // Actual touch target size
}

// NewTouchTarget creates a touch-optimized target
func NewTouchTarget(id, label, icon string, action func() tea.Cmd, theme *theme.Theme) *TouchTarget {
	return &TouchTarget{
		zone: TouchZone{
			ID:       id,
			Action:   action,
			Enabled:  true,
			MinSize:  48, // Material Design minimum (48dp)
			Priority: 0,
		},
		label:     label,
		icon:      icon,
		theme:     theme,
		touchSize: 48,
	}
}

// SetPosition sets the touch zone position
func (tt *TouchTarget) SetPosition(x, y, width, height int) {
	tt.zone.X = x
	tt.zone.Y = y
	tt.zone.Width = width
	tt.zone.Height = height

	// Ensure minimum touch target size
	if tt.zone.Width < tt.zone.MinSize {
		tt.zone.Width = tt.zone.MinSize
	}
	if tt.zone.Height < tt.zone.MinSize {
		tt.zone.Height = tt.zone.MinSize
	}
}

// Contains checks if coordinates are within touch zone
func (tt *TouchTarget) Contains(x, y int) bool {
	return x >= tt.zone.X &&
		x < tt.zone.X+tt.zone.Width &&
		y >= tt.zone.Y &&
		y < tt.zone.Y+tt.zone.Height
}

// Tap handles a tap event
func (tt *TouchTarget) Tap() tea.Cmd {
	if !tt.zone.Enabled {
		return nil
	}

	tt.pressed = true

	// Visual feedback + action
	return tea.Sequence(
		NewHapticEngine(true).Trigger(HapticSelection),
		tt.zone.Action,
		func() tea.Msg {
			time.Sleep(100 * time.Millisecond)
			return TouchReleaseMsg{ID: tt.zone.ID}
		},
	)
}

// Render renders the touch target
func (tt *TouchTarget) Render() string {
	style := lipgloss.NewStyle().
		Padding(1, 3).
		Background(tt.theme.BackgroundSecondary).
		Foreground(tt.theme.TextPrimary).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(tt.theme.Border)

	// Pressed state
	if tt.pressed {
		style = style.
			Background(tt.theme.AccentPrimary).
			Foreground(tt.theme.BackgroundPrimary).
			BorderForeground(tt.theme.AccentPrimary)
	}

	// Focused state
	if tt.focused {
		style = style.
			BorderForeground(tt.theme.AccentPrimary).
			BorderStyle(lipgloss.ThickBorder())
	}

	content := tt.icon
	if tt.label != "" {
		content += " " + tt.label
	}

	return style.Render(content)
}

// TouchReleaseMsg signals touch release
type TouchReleaseMsg struct {
	ID string
}

// TouchManager manages touch zones and hit detection
type TouchManager struct {
	zones       []*TouchZone
	activeTouch *TouchZone
	lastTap     time.Time
	doubleTapWindow time.Duration
}

// NewTouchManager creates a new touch manager
func NewTouchManager() *TouchManager {
	return &TouchManager{
		zones:           []*TouchZone{},
		doubleTapWindow: 300 * time.Millisecond,
	}
}

// RegisterZone registers a touch zone
func (tm *TouchManager) RegisterZone(zone *TouchZone) {
	tm.zones = append(tm.zones, zone)

	// Sort by priority (higher first)
	for i := len(tm.zones) - 1; i > 0; i-- {
		if tm.zones[i].Priority > tm.zones[i-1].Priority {
			tm.zones[i], tm.zones[i-1] = tm.zones[i-1], tm.zones[i]
		}
	}
}

// UnregisterZone removes a touch zone
func (tm *TouchManager) UnregisterZone(id string) {
	for i, zone := range tm.zones {
		if zone.ID == id {
			tm.zones = append(tm.zones[:i], tm.zones[i+1:]...)
			return
		}
	}
}

// ClearZones removes all touch zones
func (tm *TouchManager) ClearZones() {
	tm.zones = []*TouchZone{}
}

// HitTest finds the touch zone at coordinates
func (tm *TouchManager) HitTest(x, y int) *TouchZone {
	// Check zones in priority order
	for _, zone := range tm.zones {
		if !zone.Enabled {
			continue
		}

		if tm.contains(zone, x, y) {
			return zone
		}
	}

	return nil
}

// HandleTouch processes a touch event
func (tm *TouchManager) HandleTouch(x, y int) tea.Cmd {
	zone := tm.HitTest(x, y)
	if zone == nil {
		return nil
	}

	tm.activeTouch = zone

	// Check for double tap
	now := time.Now()
	isDoubleTap := now.Sub(tm.lastTap) < tm.doubleTapWindow
	tm.lastTap = now

	if isDoubleTap {
		// Double tap action could be different
		return tea.Sequence(
			NewHapticEngine(true).Trigger(HapticMedium),
			zone.Action,
		)
	}

	// Single tap
	return tea.Sequence(
		NewHapticEngine(true).Trigger(HapticSelection),
		zone.Action,
	)
}

// contains checks if point is within zone
func (tm *TouchManager) contains(zone *TouchZone, x, y int) bool {
	return x >= zone.X &&
		x < zone.X+zone.Width &&
		y >= zone.Y &&
		y < zone.Y+zone.Height
}

// TouchGrid creates a grid of touch targets
type TouchGrid struct {
	targets       []*TouchTarget
	cols          int
	rows          int
	spacing       int
	startX        int
	startY        int
	targetWidth   int
	targetHeight  int
	theme         *theme.Theme
}

// NewTouchGrid creates a grid layout for touch targets
func NewTouchGrid(cols, rows int, theme *theme.Theme) *TouchGrid {
	return &TouchGrid{
		cols:         cols,
		rows:         rows,
		spacing:      2,
		targetWidth:  20,
		targetHeight: 3,
		theme:        theme,
		targets:      []*TouchTarget{},
	}
}

// AddTarget adds a target to the grid
func (tg *TouchGrid) AddTarget(target *TouchTarget) {
	tg.targets = append(tg.targets, target)

	// Calculate position in grid
	index := len(tg.targets) - 1
	col := index % tg.cols
	row := index / tg.cols

	x := tg.startX + col*(tg.targetWidth+tg.spacing)
	y := tg.startY + row*(tg.targetHeight+tg.spacing)

	target.SetPosition(x, y, tg.targetWidth, tg.targetHeight)
}

// Render renders the grid
func (tg *TouchGrid) Render() string {
	if len(tg.targets) == 0 {
		return ""
	}

	// Build grid row by row
	rows := []string{}
	for r := 0; r < tg.rows; r++ {
		cols := []string{}
		for c := 0; c < tg.cols; c++ {
			idx := r*tg.cols + c
			if idx < len(tg.targets) {
				cols = append(cols, tg.targets[idx].Render())
			}
		}
		if len(cols) > 0 {
			row := lipgloss.JoinHorizontal(lipgloss.Top, cols...)
			rows = append(rows, row)
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

// PhoneTouchButtons creates phone-optimized button layout
func PhoneTouchButtons(actions []struct {
	ID    string
	Icon  string
	Label string
	Action func() tea.Cmd
}, theme *theme.Theme, width int) string {
	// Calculate button size based on screen width
	numButtons := len(actions)
	if numButtons == 0 {
		return ""
	}

	spacing := 2
	availableWidth := width - (numButtons-1)*spacing - 4
	buttonWidth := availableWidth / numButtons

	// Minimum touch target width
	minWidth := 12
	if buttonWidth < minWidth {
		// Stack vertically if too narrow
		return renderVerticalButtons(actions, theme, width)
	}

	// Render horizontally
	buttons := []string{}
	for _, action := range actions {
		buttonStyle := lipgloss.NewStyle().
			Width(buttonWidth).
			Height(3).
			Padding(1).
			Align(lipgloss.Center, lipgloss.Center).
			Background(theme.BackgroundSecondary).
			Foreground(theme.TextPrimary).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(theme.Border)

		content := action.Icon
		if buttonWidth > 15 && action.Label != "" {
			content += "\n" + action.Label
		}

		buttons = append(buttons, buttonStyle.Render(content))
	}

	return lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Render(lipgloss.JoinHorizontal(lipgloss.Top, buttons...))
}

func renderVerticalButtons(actions []struct {
	ID    string
	Icon  string
	Label string
	Action func() tea.Cmd
}, theme *theme.Theme, width int) string {
	buttons := []string{}

	for _, action := range actions {
		buttonStyle := lipgloss.NewStyle().
			Width(width - 4).
			Height(3).
			Padding(1, 2).
			Background(theme.BackgroundSecondary).
			Foreground(theme.TextPrimary).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(theme.Border).
			MarginBottom(1)

		content := lipgloss.JoinHorizontal(
			lipgloss.Left,
			action.Icon,
			" ",
			action.Label,
		)

		buttons = append(buttons, buttonStyle.Render(content))
	}

	return lipgloss.JoinVertical(lipgloss.Left, buttons...)
}

// TouchDebugOverlay shows touch zone boundaries (for debugging)
func TouchDebugOverlay(zones []*TouchZone, theme *theme.Theme, width, height int) string {
	// Draw touch zones as colored rectangles
	canvas := make([][]string, height)
	for i := range canvas {
		canvas[i] = make([]string, width)
		for j := range canvas[i] {
			canvas[i][j] = " "
		}
	}

	debugStyle := lipgloss.NewStyle().
		Foreground(theme.Error).
		Background(theme.BackgroundSecondary)

	// Draw each zone
	for _, zone := range zones {
		for y := zone.Y; y < zone.Y+zone.Height && y < height; y++ {
			for x := zone.X; x < zone.X+zone.Width && x < width; x++ {
				if y >= 0 && x >= 0 {
					canvas[y][x] = debugStyle.Render("▓")
				}
			}
		}

		// Draw zone label at top-left
		if zone.Y >= 0 && zone.Y < height && zone.X >= 0 && zone.X < width {
			labelStyle := lipgloss.NewStyle().
				Foreground(theme.Warning).
				Bold(true)

			label := labelStyle.Render(zone.ID)
			for i, ch := range label {
				if zone.X+i < width {
					canvas[zone.Y][zone.X+i] = string(ch)
				}
			}
		}
	}

	// Convert canvas to string
	lines := []string{}
	for _, row := range canvas {
		lines = append(lines, lipgloss.JoinHorizontal(lipgloss.Left, row...))
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

// TouchFeedbackOverlay shows visual feedback for touch
type TouchFeedbackOverlay struct {
	active   bool
	x        int
	y        int
	startTime time.Time
	duration time.Duration
	theme    *theme.Theme
}

// NewTouchFeedbackOverlay creates a touch feedback overlay
func NewTouchFeedbackOverlay(theme *theme.Theme) *TouchFeedbackOverlay {
	return &TouchFeedbackOverlay{
		theme:    theme,
		duration: 200 * time.Millisecond,
	}
}

// Show displays touch feedback at coordinates
func (tfo *TouchFeedbackOverlay) Show(x, y int) tea.Cmd {
	tfo.active = true
	tfo.x = x
	tfo.y = y
	tfo.startTime = time.Now()

	return tea.Tick(tfo.duration, func(t time.Time) tea.Msg {
		return TouchFeedbackHideMsg{}
	})
}

// Hide hides the feedback
func (tfo *TouchFeedbackOverlay) Hide() {
	tfo.active = false
}

// IsActive returns whether feedback is active
func (tfo *TouchFeedbackOverlay) IsActive() bool {
	return tfo.active && time.Since(tfo.startTime) < tfo.duration
}

// Render renders the touch feedback
func (tfo *TouchFeedbackOverlay) Render(width, height int) string {
	if !tfo.IsActive() {
		return ""
	}

	// Ripple effect
	elapsed := time.Since(tfo.startTime)
	progress := float64(elapsed) / float64(tfo.duration)

	// Animate size
	size := int(progress * 5)

	rippleStyle := lipgloss.NewStyle().
		Foreground(tfo.theme.AccentPrimary).
		Faint(true)

	ripple := ""
	for i := 0; i < size; i++ {
		ripple += "◯"
	}

	// Position at touch point
	positioned := lipgloss.Place(
		width,
		height,
		tfo.x,
		tfo.y,
		rippleStyle.Render(ripple),
	)

	return positioned
}

// TouchFeedbackHideMsg signals to hide touch feedback
type TouchFeedbackHideMsg struct{}

// Accessibility: Ensure touch targets meet minimum sizes
const (
	MinTouchTargetSize = 48 // Material Design guideline
	MinTouchTargetSizeiOS = 44 // iOS HIG guideline
	RecommendedTouchSize = 48 // Use 48dp for consistency
)

// ValidateTouchTarget checks if touch target meets accessibility guidelines
func ValidateTouchTarget(width, height int) bool {
	return width >= MinTouchTargetSize && height >= MinTouchTargetSize
}

// ExpandTouchTarget expands a target to minimum size if needed
func ExpandTouchTarget(zone *TouchZone) {
	if zone.Width < MinTouchTargetSize {
		// Center expansion
		diff := MinTouchTargetSize - zone.Width
		zone.X -= diff / 2
		zone.Width = MinTouchTargetSize
	}

	if zone.Height < MinTouchTargetSize {
		diff := MinTouchTargetSize - zone.Height
		zone.Y -= diff / 2
		zone.Height = MinTouchTargetSize
	}
}

package responsive

import (
	"github.com/charmbracelet/lipgloss/v2"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/sst/opencode/internal/theme"
)

// FocusableElement represents an element that can receive focus
type FocusableElement interface {
	ID() string
	IsFocused() bool
	Focus()
	Blur()
	HandleKey(key string) tea.Cmd
	Render(theme *theme.Theme) string
}

// FocusZone represents a group of focusable elements
type FocusZone string

const (
	ZoneInput       FocusZone = "input"
	ZoneMessages    FocusZone = "messages"
	ZoneSidebar     FocusZone = "sidebar"
	ZoneQuickActions FocusZone = "quick_actions"
	ZoneAIPicker    FocusZone = "ai_picker"
	ZoneReactions   FocusZone = "reactions"
	ZoneHistory     FocusZone = "history"
)

// FocusManager manages focus state and keyboard navigation
type FocusManager struct {
	zones          map[FocusZone][]FocusableElement
	currentZone    FocusZone
	focusIndex     map[FocusZone]int
	focusHistory   []FocusZone
	enabled        bool
	visualFocus    bool
	keyboardMode   bool // True when using keyboard (show focus rings)
}

// NewFocusManager creates a new focus manager
func NewFocusManager() *FocusManager {
	return &FocusManager{
		zones:        make(map[FocusZone][]FocusableElement),
		focusIndex:   make(map[FocusZone]int),
		focusHistory: []FocusZone{},
		enabled:      true,
		visualFocus:  true,
		keyboardMode: false,
	}
}

// RegisterZone registers a focus zone with elements
func (fm *FocusManager) RegisterZone(zone FocusZone, elements []FocusableElement) {
	fm.zones[zone] = elements
	if _, exists := fm.focusIndex[zone]; !exists {
		fm.focusIndex[zone] = 0
	}
}

// SetZone sets the current focus zone
func (fm *FocusManager) SetZone(zone FocusZone) {
	if _, exists := fm.zones[zone]; !exists {
		return
	}

	// Blur current zone
	fm.blurCurrentElement()

	// Save to history
	fm.focusHistory = append(fm.focusHistory, fm.currentZone)

	fm.currentZone = zone

	// Focus first element in new zone
	fm.focusCurrentElement()
}

// GetCurrentZone returns the current focus zone
func (fm *FocusManager) GetCurrentZone() FocusZone {
	return fm.currentZone
}

// Next moves focus to next element in current zone
func (fm *FocusManager) Next() tea.Cmd {
	if !fm.enabled {
		return nil
	}

	fm.keyboardMode = true
	elements := fm.zones[fm.currentZone]
	if len(elements) == 0 {
		return nil
	}

	fm.blurCurrentElement()

	// Move to next element (wrap around)
	idx := fm.focusIndex[fm.currentZone]
	idx = (idx + 1) % len(elements)
	fm.focusIndex[fm.currentZone] = idx

	fm.focusCurrentElement()

	// Haptic feedback for navigation
	return NewHapticEngine(true).Trigger(HapticSelection)
}

// Previous moves focus to previous element
func (fm *FocusManager) Previous() tea.Cmd {
	if !fm.enabled {
		return nil
	}

	fm.keyboardMode = true
	elements := fm.zones[fm.currentZone]
	if len(elements) == 0 {
		return nil
	}

	fm.blurCurrentElement()

	// Move to previous element (wrap around)
	idx := fm.focusIndex[fm.currentZone]
	idx = (idx - 1 + len(elements)) % len(elements)
	fm.focusIndex[fm.currentZone] = idx

	fm.focusCurrentElement()

	return NewHapticEngine(true).Trigger(HapticSelection)
}

// NextZone moves to next focus zone
func (fm *FocusManager) NextZone() tea.Cmd {
	if !fm.enabled {
		return nil
	}

	fm.keyboardMode = true
	zones := fm.getZoneOrder()
	if len(zones) == 0 {
		return nil
	}

	// Find current zone index
	currentIdx := 0
	for i, z := range zones {
		if z == fm.currentZone {
			currentIdx = i
			break
		}
	}

	// Move to next zone
	nextIdx := (currentIdx + 1) % len(zones)
	fm.SetZone(zones[nextIdx])

	return NewHapticEngine(true).Trigger(HapticMedium)
}

// PreviousZone moves to previous focus zone
func (fm *FocusManager) PreviousZone() tea.Cmd {
	if !fm.enabled {
		return nil
	}

	fm.keyboardMode = true
	zones := fm.getZoneOrder()
	if len(zones) == 0 {
		return nil
	}

	// Find current zone index
	currentIdx := 0
	for i, z := range zones {
		if z == fm.currentZone {
			currentIdx = i
			break
		}
	}

	// Move to previous zone
	prevIdx := (currentIdx - 1 + len(zones)) % len(zones)
	fm.SetZone(zones[prevIdx])

	return NewHapticEngine(true).Trigger(HapticMedium)
}

// Back returns to previous zone from history
func (fm *FocusManager) Back() tea.Cmd {
	if len(fm.focusHistory) == 0 {
		return nil
	}

	// Pop from history
	prevZone := fm.focusHistory[len(fm.focusHistory)-1]
	fm.focusHistory = fm.focusHistory[:len(fm.focusHistory)-1]

	fm.blurCurrentElement()
	fm.currentZone = prevZone
	fm.focusCurrentElement()

	return NewHapticEngine(true).Trigger(HapticLight)
}

// HandleKey handles keyboard input for focused element
func (fm *FocusManager) HandleKey(key string) tea.Cmd {
	fm.keyboardMode = true

	// Global navigation keys
	switch key {
	case "tab":
		return fm.Next()
	case "shift+tab":
		return fm.Previous()
	case "ctrl+tab":
		return fm.NextZone()
	case "ctrl+shift+tab":
		return fm.PreviousZone()
	case "esc":
		return fm.Back()
	}

	// Delegate to focused element
	element := fm.getCurrentElement()
	if element != nil {
		return element.HandleKey(key)
	}

	return nil
}

// FocusElement focuses a specific element by ID
func (fm *FocusManager) FocusElement(id string) {
	fm.keyboardMode = true

	// Search all zones for element
	for zone, elements := range fm.zones {
		for i, elem := range elements {
			if elem.ID() == id {
				fm.blurCurrentElement()
				fm.currentZone = zone
				fm.focusIndex[zone] = i
				fm.focusCurrentElement()
				return
			}
		}
	}
}

// GetFocusedElement returns currently focused element
func (fm *FocusManager) GetFocusedElement() FocusableElement {
	return fm.getCurrentElement()
}

// EnableVisualFocus shows/hides focus indicators
func (fm *FocusManager) EnableVisualFocus(enabled bool) {
	fm.visualFocus = enabled
}

// IsKeyboardMode returns whether in keyboard navigation mode
func (fm *FocusManager) IsKeyboardMode() bool {
	return fm.keyboardMode
}

// SetMouseMode switches to mouse/touch mode (hide focus rings)
func (fm *FocusManager) SetMouseMode() {
	fm.keyboardMode = false
}

// Helper methods
func (fm *FocusManager) getCurrentElement() FocusableElement {
	elements := fm.zones[fm.currentZone]
	if len(elements) == 0 {
		return nil
	}

	idx := fm.focusIndex[fm.currentZone]
	if idx >= len(elements) {
		idx = 0
		fm.focusIndex[fm.currentZone] = 0
	}

	return elements[idx]
}

func (fm *FocusManager) focusCurrentElement() {
	element := fm.getCurrentElement()
	if element != nil {
		element.Focus()
	}
}

func (fm *FocusManager) blurCurrentElement() {
	element := fm.getCurrentElement()
	if element != nil {
		element.Blur()
	}
}

func (fm *FocusManager) getZoneOrder() []FocusZone {
	// Define logical tab order
	order := []FocusZone{
		ZoneInput,
		ZoneQuickActions,
		ZoneMessages,
		ZoneSidebar,
		ZoneHistory,
		ZoneReactions,
		ZoneAIPicker,
	}

	// Filter to only registered zones
	available := []FocusZone{}
	for _, zone := range order {
		if _, exists := fm.zones[zone]; exists {
			available = append(available, zone)
		}
	}

	return available
}

// FocusRing renders a visual focus indicator
func FocusRing(focused bool, keyboardMode bool, theme *theme.Theme) lipgloss.Style {
	style := lipgloss.NewStyle()

	if focused && keyboardMode {
		// Show prominent focus ring for keyboard navigation
		style = style.
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(theme.AccentPrimary)
	} else if focused {
		// Subtle indicator for mouse/touch
		style = style.
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(theme.AccentSecondary)
	} else {
		// No focus
		style = style.
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(theme.Border)
	}

	return style
}

// FocusIndicator renders a simple focus indicator
func FocusIndicator(focused bool, theme *theme.Theme) string {
	if !focused {
		return "  "
	}

	return lipgloss.NewStyle().
		Foreground(theme.AccentPrimary).
		Bold(true).
		Render("▶ ")
}

// KeyboardHint shows keyboard navigation hints
type KeyboardHint struct {
	Key         string
	Description string
	Zone        FocusZone
}

// GlobalKeyboardHints returns global keyboard shortcuts
func GlobalKeyboardHints() []KeyboardHint {
	return []KeyboardHint{
		{"Tab", "Next element", ""},
		{"Shift+Tab", "Previous element", ""},
		{"Ctrl+Tab", "Next zone", ""},
		{"Ctrl+Shift+Tab", "Previous zone", ""},
		{"Esc", "Back / Cancel", ""},
		{"Enter", "Activate / Submit", ""},
		{"Space", "Select / Toggle", ""},
		{"↑↓", "Navigate list", ""},
		{"←→", "Navigate messages", ""},
		{"/", "Search / Command", ""},
		{"?", "Show help", ""},
		{"Ctrl+K", "Quick actions", ""},
		{"Ctrl+V", "Voice input", ""},
		{"Ctrl+R", "Instant replay", ""},
		{"Ctrl+H", "Show history", ""},
		{"Ctrl+,", "Settings", ""},
	}
}

// ZoneKeyboardHints returns zone-specific shortcuts
func ZoneKeyboardHints(zone FocusZone) []KeyboardHint {
	hints := map[FocusZone][]KeyboardHint{
		ZoneInput: {
			{"Enter", "Send message", ZoneInput},
			{"Ctrl+Enter", "New line", ZoneInput},
			{"↑", "Previous command", ZoneInput},
			{"↓", "Next command", ZoneInput},
		},
		ZoneMessages: {
			{"↑↓", "Navigate messages", ZoneMessages},
			{"r", "React to message", ZoneMessages},
			{"c", "Copy message", ZoneMessages},
			{"d", "Delete message", ZoneMessages},
			{"e", "Edit message", ZoneMessages},
		},
		ZoneQuickActions: {
			{"1-9", "Select action", ZoneQuickActions},
			{"Enter", "Activate action", ZoneQuickActions},
		},
		ZoneAIPicker: {
			{"1-3", "Select AI provider", ZoneAIPicker},
			{"Enter", "Confirm selection", ZoneAIPicker},
			{"Esc", "Cancel", ZoneAIPicker},
		},
		ZoneReactions: {
			{"1-7", "Select reaction", ZoneReactions},
			{"Esc", "Cancel", ZoneReactions},
		},
		ZoneHistory: {
			{"↑↓", "Navigate history", ZoneHistory},
			{"Enter", "Select command", ZoneHistory},
			{"/", "Search history", ZoneHistory},
		},
	}

	return hints[zone]
}

// RenderKeyboardHelp renders keyboard help overlay
func RenderKeyboardHelp(theme *theme.Theme, width int) string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.AccentPrimary).
		Bold(true).
		Width(width - 4).
		Align(lipgloss.Center).
		MarginBottom(1)

	title := titleStyle.Render("⌨️  Keyboard Shortcuts")

	// Global shortcuts
	globalSection := renderHintSection("Global", GlobalKeyboardHints(), theme, width)

	// Instructions
	instructionStyle := lipgloss.NewStyle().
		Foreground(theme.Info).
		Width(width - 4).
		Align(lipgloss.Center).
		MarginTop(1)

	instruction := instructionStyle.Render("Press ? to toggle • Esc to close")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		globalSection,
		instruction,
	)

	containerStyle := lipgloss.NewStyle().
		Width(width - 2).
		Padding(1, 2).
		Background(theme.BackgroundSecondary).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(theme.AccentPrimary)

	return containerStyle.Render(content)
}

func renderHintSection(title string, hints []KeyboardHint, theme *theme.Theme, width int) string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.AccentSecondary).
		Bold(true).
		MarginBottom(1)

	sectionTitle := titleStyle.Render(title)

	hintLines := []string{}
	for _, hint := range hints {
		keyStyle := lipgloss.NewStyle().
			Foreground(theme.AccentPrimary).
			Background(theme.BackgroundPrimary).
			Padding(0, 1).
			Bold(true)

		descStyle := lipgloss.NewStyle().
			Foreground(theme.TextPrimary).
			MarginLeft(2)

		line := lipgloss.JoinHorizontal(
			lipgloss.Left,
			keyStyle.Render(hint.Key),
			descStyle.Render(hint.Description),
		)

		hintLines = append(hintLines, line)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		sectionTitle,
		lipgloss.JoinVertical(lipgloss.Left, hintLines...),
	)
}

// FocusDebugInfo returns debug info about focus state
func (fm *FocusManager) FocusDebugInfo() string {
	element := fm.getCurrentElement()
	elementID := "none"
	if element != nil {
		elementID = element.ID()
	}

	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Render(
			"Focus: " + string(fm.currentZone) +
				" [" + elementID + "]" +
				" | Keyboard: " + boolToStr(fm.keyboardMode),
		)
}

func boolToStr(b bool) string {
	if b {
		return "YES"
	}
	return "NO"
}

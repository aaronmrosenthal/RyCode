package responsive

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/sst/opencode/internal/theme"
)

// AccessibilityLevel defines the level of accessibility features
type AccessibilityLevel string

const (
	A11yMinimal  AccessibilityLevel = "minimal"  // Basic features
	A11yStandard AccessibilityLevel = "standard" // Recommended
	A11yFull     AccessibilityLevel = "full"     // Maximum accessibility
)

// AccessibilityConfig configures accessibility features
type AccessibilityConfig struct {
	Level              AccessibilityLevel
	HighContrast       bool
	LargeText          bool
	ReducedMotion      bool
	ScreenReaderMode   bool
	KeyboardOnly       bool
	ShowFocusIndicators bool
	AnnounceChanges    bool
	ColorBlindMode     ColorBlindMode
}

// ColorBlindMode represents different color blindness types
type ColorBlindMode string

const (
	ColorBlindNone        ColorBlindMode = "none"
	ColorBlindProtanopia  ColorBlindMode = "protanopia"  // Red-blind
	ColorBlindDeuteranopia ColorBlindMode = "deuteranopia" // Green-blind
	ColorBlindTritanopia  ColorBlindMode = "tritanopia"   // Blue-blind
)

// NewAccessibilityConfig creates default accessibility config
func NewAccessibilityConfig() *AccessibilityConfig {
	return &AccessibilityConfig{
		Level:              A11yStandard,
		HighContrast:       false,
		LargeText:          false,
		ReducedMotion:      false,
		ScreenReaderMode:   false,
		KeyboardOnly:       false,
		ShowFocusIndicators: true,
		AnnounceChanges:    true,
		ColorBlindMode:     ColorBlindNone,
	}
}

// AccessibilityManager manages accessibility features
type AccessibilityManager struct {
	config     *AccessibilityConfig
	theme      *theme.Theme
	announcements []string
}

// NewAccessibilityManager creates an accessibility manager
func NewAccessibilityManager(config *AccessibilityConfig, theme *theme.Theme) *AccessibilityManager {
	return &AccessibilityManager{
		config:     config,
		theme:      theme,
		announcements: []string{},
	}
}

// Announce adds an announcement for screen readers
func (am *AccessibilityManager) Announce(message string) {
	if !am.config.AnnounceChanges {
		return
	}

	am.announcements = append(am.announcements, message)

	// Keep only last 10 announcements
	if len(am.announcements) > 10 {
		am.announcements = am.announcements[1:]
	}
}

// GetAnnouncements returns pending announcements
func (am *AccessibilityManager) GetAnnouncements() []string {
	announcements := am.announcements
	am.announcements = []string{} // Clear after reading
	return announcements
}

// GetTextScale returns text scale multiplier
func (am *AccessibilityManager) GetTextScale() float64 {
	if am.config.LargeText {
		return 1.5
	}
	return 1.0
}

// ShouldShowAnimation returns whether animations should play
func (am *AccessibilityManager) ShouldShowAnimation() bool {
	return !am.config.ReducedMotion
}

// GetFocusStyle returns focus style based on config
func (am *AccessibilityManager) GetFocusStyle(focused bool) lipgloss.Style {
	style := lipgloss.NewStyle()

	if !focused || !am.config.ShowFocusIndicators {
		return style
	}

	if am.config.HighContrast {
		// High contrast focus
		return style.
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(lipgloss.Color("15")). // White
			Background(lipgloss.Color("0")) // Black
	}

	// Standard focus
	return style.
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(am.theme.AccentPrimary)
}

// AdaptThemeForAccessibility modifies theme for accessibility
func (am *AccessibilityManager) AdaptThemeForAccessibility(baseTheme *theme.Theme) *theme.Theme {
	adapted := *baseTheme // Copy

	// High contrast mode
	if am.config.HighContrast {
		adapted.BackgroundPrimary = lipgloss.Color("0")   // Black
		adapted.BackgroundSecondary = lipgloss.Color("8") // Dark gray
		adapted.TextPrimary = lipgloss.Color("15")        // White
		adapted.TextSecondary = lipgloss.Color("15")
		adapted.Border = lipgloss.Color("15")
		adapted.AccentPrimary = lipgloss.Color("11")   // Bright yellow
		adapted.AccentSecondary = lipgloss.Color("14") // Bright cyan
	}

	// Color blind adaptations
	switch am.config.ColorBlindMode {
	case ColorBlindProtanopia, ColorBlindDeuteranopia:
		// Red-green blindness: use blue/yellow instead
		adapted.Success = lipgloss.Color("12")  // Bright blue
		adapted.Error = lipgloss.Color("11")    // Bright yellow
		adapted.Warning = lipgloss.Color("14")  // Bright cyan
		adapted.Info = lipgloss.Color("13")     // Bright magenta

	case ColorBlindTritanopia:
		// Blue-yellow blindness: use red/cyan
		adapted.Success = lipgloss.Color("10")  // Bright green
		adapted.Error = lipgloss.Color("9")     // Bright red
		adapted.Warning = lipgloss.Color("13")  // Bright magenta
		adapted.Info = lipgloss.Color("14")     // Bright cyan
	}

	return &adapted
}

// ARIALabel generates ARIA-like labels for screen readers
type ARIALabel struct {
	Label       string
	Role        string
	Description string
	State       string
	Level       int // Heading level
}

// RenderARIALabel renders an accessible label
func RenderARIALabel(label ARIALabel, theme *theme.Theme) string {
	if label.Label == "" {
		return ""
	}

	parts := []string{}

	// Role
	if label.Role != "" {
		parts = append(parts, fmt.Sprintf("[%s]", label.Role))
	}

	// Label
	parts = append(parts, label.Label)

	// State
	if label.State != "" {
		parts = append(parts, fmt.Sprintf("(%s)", label.State))
	}

	// Description
	if label.Description != "" {
		parts = append(parts, "-", label.Description)
	}

	text := strings.Join(parts, " ")

	style := lipgloss.NewStyle().
		Foreground(theme.TextPrimary)

	// Heading levels get bold
	if label.Level > 0 {
		style = style.Bold(true)
	}

	return style.Render(text)
}

// SkipLink allows keyboard users to skip navigation
type SkipLink struct {
	label  string
	target string
	visible bool
}

// NewSkipLink creates a skip link
func NewSkipLink(label, target string) *SkipLink {
	return &SkipLink{
		label:  label,
		target: target,
		visible: false,
	}
}

// Show makes skip link visible (on focus)
func (sl *SkipLink) Show() {
	sl.visible = true
}

// Hide hides skip link
func (sl *SkipLink) Hide() {
	sl.visible = false
}

// Render renders the skip link
func (sl *SkipLink) Render(theme *theme.Theme) string {
	if !sl.visible {
		return ""
	}

	style := lipgloss.NewStyle().
		Background(theme.AccentPrimary).
		Foreground(theme.BackgroundPrimary).
		Padding(0, 2).
		Bold(true)

	return style.Render("Skip to " + sl.label + " [Enter]")
}

// AccessibilityChecker validates UI for accessibility issues
type AccessibilityChecker struct {
	issues []AccessibilityIssue
}

// AccessibilityIssue represents an a11y problem
type AccessibilityIssue struct {
	Level       string // "error", "warning", "info"
	Component   string
	Description string
	Fix         string
}

// NewAccessibilityChecker creates a checker
func NewAccessibilityChecker() *AccessibilityChecker {
	return &AccessibilityChecker{
		issues: []AccessibilityIssue{},
	}
}

// CheckTouchTarget validates touch target size
func (ac *AccessibilityChecker) CheckTouchTarget(id string, width, height int) {
	if width < MinTouchTargetSize || height < MinTouchTargetSize {
		ac.issues = append(ac.issues, AccessibilityIssue{
			Level:     "error",
			Component: id,
			Description: fmt.Sprintf(
				"Touch target too small: %dx%d (minimum: %dx%d)",
				width, height,
				MinTouchTargetSize, MinTouchTargetSize,
			),
			Fix: "Increase target size to at least 48x48 pixels",
		})
	}
}

// CheckContrast validates color contrast
func (ac *AccessibilityChecker) CheckContrast(foreground, background lipgloss.Color) {
	// Simplified contrast check (real implementation would calculate actual contrast ratio)
	// WCAG AA requires 4.5:1 for normal text, 3:1 for large text

	// This is a placeholder - proper implementation would:
	// 1. Convert colors to RGB
	// 2. Calculate relative luminance
	// 3. Calculate contrast ratio
	// 4. Compare to WCAG standards

	// For now, just check if they're the same
	if foreground == background {
		ac.issues = append(ac.issues, AccessibilityIssue{
			Level:       "error",
			Component:   "color",
			Description: "Foreground and background colors are identical",
			Fix:         "Use contrasting colors (WCAG AA: 4.5:1 minimum)",
		})
	}
}

// CheckKeyboardAccess validates keyboard accessibility
func (ac *AccessibilityChecker) CheckKeyboardAccess(element string, hasKeyHandler bool, isFocusable bool) {
	if !hasKeyHandler && isFocusable {
		ac.issues = append(ac.issues, AccessibilityIssue{
			Level:       "warning",
			Component:   element,
			Description: "Element is focusable but has no keyboard handler",
			Fix:         "Add keyboard event handler for Enter/Space keys",
		})
	}
}

// CheckLabel validates that interactive elements have labels
func (ac *AccessibilityChecker) CheckLabel(element string, hasLabel bool) {
	if !hasLabel {
		ac.issues = append(ac.issues, AccessibilityIssue{
			Level:       "error",
			Component:   element,
			Description: "Interactive element has no accessible label",
			Fix:         "Add aria-label or visible text label",
		})
	}
}

// GetIssues returns all found issues
func (ac *AccessibilityChecker) GetIssues() []AccessibilityIssue {
	return ac.issues
}

// Report generates accessibility report
func (ac *AccessibilityChecker) Report(theme *theme.Theme) string {
	if len(ac.issues) == 0 {
		successStyle := lipgloss.NewStyle().
			Foreground(theme.Success).
			Bold(true)

		return successStyle.Render("✅ No accessibility issues found!")
	}

	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Warning).
		Bold(true).
		MarginBottom(1)

	title := titleStyle.Render(
		fmt.Sprintf("⚠️  Found %d accessibility issues:", len(ac.issues)),
	)

	issueLines := []string{title}

	for i, issue := range ac.issues {
		levelStyle := lipgloss.NewStyle()

		switch issue.Level {
		case "error":
			levelStyle = levelStyle.Foreground(theme.Error)
		case "warning":
			levelStyle = levelStyle.Foreground(theme.Warning)
		case "info":
			levelStyle = levelStyle.Foreground(theme.Info)
		}

		level := levelStyle.Render(strings.ToUpper(issue.Level))
		component := lipgloss.NewStyle().
			Foreground(theme.AccentSecondary).
			Render(issue.Component)

		desc := lipgloss.NewStyle().
			Foreground(theme.TextPrimary).
			Render(issue.Description)

		fix := lipgloss.NewStyle().
			Foreground(theme.TextDim).
			Render("→ " + issue.Fix)

		issueText := fmt.Sprintf(
			"%d. [%s] %s\n   %s\n   %s",
			i+1, level, component, desc, fix,
		)

		issueLines = append(issueLines, issueText)
	}

	return lipgloss.JoinVertical(lipgloss.Left, issueLines...)
}

// LiveRegion represents a dynamic content area for screen readers
type LiveRegion struct {
	content  string
	priority string // "polite" or "assertive"
	changed  bool
}

// NewLiveRegion creates a live region
func NewLiveRegion(priority string) *LiveRegion {
	return &LiveRegion{
		priority: priority,
	}
}

// Update updates the live region content
func (lr *LiveRegion) Update(content string) {
	if lr.content != content {
		lr.content = content
		lr.changed = true
	}
}

// GetUpdate returns update if content changed
func (lr *LiveRegion) GetUpdate() (string, bool) {
	if lr.changed {
		lr.changed = false
		return lr.content, true
	}
	return "", false
}

// Render renders the live region
func (lr *LiveRegion) Render(theme *theme.Theme) string {
	if lr.content == "" {
		return ""
	}

	style := lipgloss.NewStyle().
		Foreground(theme.Info).
		Padding(1)

	prefix := "[Live] "
	if lr.priority == "assertive" {
		prefix = "[Alert] "
	}

	return style.Render(prefix + lr.content)
}

// AccessibilitySettings renders accessibility settings UI
func AccessibilitySettings(config *AccessibilityConfig, theme *theme.Theme, width int) string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.AccentPrimary).
		Bold(true).
		Width(width - 4).
		Align(lipgloss.Center).
		MarginBottom(1)

	title := titleStyle.Render("♿ Accessibility Settings")

	settings := []struct {
		label   string
		value   string
		key     string
	}{
		{"High Contrast", boolToStr(config.HighContrast), "1"},
		{"Large Text", boolToStr(config.LargeText), "2"},
		{"Reduced Motion", boolToStr(config.ReducedMotion), "3"},
		{"Screen Reader Mode", boolToStr(config.ScreenReaderMode), "4"},
		{"Keyboard Only", boolToStr(config.KeyboardOnly), "5"},
		{"Show Focus Indicators", boolToStr(config.ShowFocusIndicators), "6"},
		{"Color Blind Mode", string(config.ColorBlindMode), "7"},
	}

	items := []string{}
	for _, setting := range settings {
		keyStyle := lipgloss.NewStyle().
			Foreground(theme.TextDim).
			Render(setting.key)

		labelStyle := lipgloss.NewStyle().
			Foreground(theme.TextPrimary).
			Render(setting.label)

		valueStyle := lipgloss.NewStyle().
			Foreground(theme.AccentSecondary).
			Bold(true).
			Render(setting.value)

		item := lipgloss.JoinHorizontal(
			lipgloss.Left,
			keyStyle, " ",
			labelStyle, ": ",
			valueStyle,
		)

		items = append(items, item)
	}

	hintStyle := lipgloss.NewStyle().
		Foreground(theme.TextDim).
		Width(width - 4).
		Align(lipgloss.Center).
		MarginTop(1)

	hint := hintStyle.Render("Press 1-7 to toggle • ESC to close")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		lipgloss.JoinVertical(lipgloss.Left, items...),
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

// ToggleAccessibilityFeature toggles an accessibility feature
func ToggleAccessibilityFeature(config *AccessibilityConfig, feature string) {
	switch feature {
	case "high_contrast":
		config.HighContrast = !config.HighContrast
	case "large_text":
		config.LargeText = !config.LargeText
	case "reduced_motion":
		config.ReducedMotion = !config.ReducedMotion
	case "screen_reader":
		config.ScreenReaderMode = !config.ScreenReaderMode
	case "keyboard_only":
		config.KeyboardOnly = !config.KeyboardOnly
	case "focus_indicators":
		config.ShowFocusIndicators = !config.ShowFocusIndicators
	}
}

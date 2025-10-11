package accessibility

import (
	"sync"
)

// AccessibilityMode represents different accessibility preferences
type AccessibilityMode int

const (
	ModeDefault AccessibilityMode = iota
	ModeHighContrast
	ModeReducedMotion
	ModeScreenReader
	ModeKeyboardOnly
)

// AccessibilitySettings holds all accessibility preferences
type AccessibilitySettings struct {
	mu sync.RWMutex

	// Visual settings
	HighContrast      bool
	ReducedMotion     bool
	LargeText         bool
	IncreasedSpacing  bool

	// Interaction settings
	KeyboardOnly      bool
	ScreenReaderMode  bool
	ShowKeyboardHints bool
	VerboseLabels     bool

	// Animation settings
	DisableAnimations bool
	SlowAnimations    bool
	AnimationSpeed    float64 // 0.5 = half speed, 1.0 = normal, 2.0 = double

	// Color settings
	ColorBlindMode    string // "", "protanopia", "deuteranopia", "tritanopia"
	ForceDarkMode     bool
	ForceLightMode    bool

	// Focus indicators
	EnhancedFocus     bool
	FocusIndicatorSize int // 1 = normal, 2 = large, 3 = extra large
}

// Global accessibility settings
var globalSettings = &AccessibilitySettings{
	HighContrast:       false,
	ReducedMotion:      false,
	LargeText:          false,
	IncreasedSpacing:   false,
	KeyboardOnly:       false,
	ScreenReaderMode:   false,
	ShowKeyboardHints:  true,
	VerboseLabels:      false,
	DisableAnimations:  false,
	SlowAnimations:     false,
	AnimationSpeed:     1.0,
	ColorBlindMode:     "",
	ForceDarkMode:      false,
	ForceLightMode:     false,
	EnhancedFocus:      false,
	FocusIndicatorSize: 1,
}

// GetSettings returns the global accessibility settings
func GetSettings() *AccessibilitySettings {
	return globalSettings
}

// EnableHighContrast enables high contrast mode
func (s *AccessibilitySettings) EnableHighContrast() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.HighContrast = true
}

// DisableHighContrast disables high contrast mode
func (s *AccessibilitySettings) DisableHighContrast() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.HighContrast = false
}

// IsHighContrast returns whether high contrast mode is enabled
func (s *AccessibilitySettings) IsHighContrast() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.HighContrast
}

// EnableReducedMotion enables reduced motion mode
func (s *AccessibilitySettings) EnableReducedMotion() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ReducedMotion = true
	s.DisableAnimations = true
}

// DisableReducedMotion disables reduced motion mode
func (s *AccessibilitySettings) DisableReducedMotion() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ReducedMotion = false
	s.DisableAnimations = false
}

// IsReducedMotion returns whether reduced motion is enabled
func (s *AccessibilitySettings) IsReducedMotion() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.ReducedMotion
}

// EnableScreenReader enables screen reader mode
func (s *AccessibilitySettings) EnableScreenReader() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ScreenReaderMode = true
	s.VerboseLabels = true
	s.ShowKeyboardHints = true
}

// IsScreenReaderMode returns whether screen reader mode is enabled
func (s *AccessibilitySettings) IsScreenReaderMode() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.ScreenReaderMode
}

// EnableKeyboardOnly enables keyboard-only mode
func (s *AccessibilitySettings) EnableKeyboardOnly() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.KeyboardOnly = true
	s.ShowKeyboardHints = true
	s.EnhancedFocus = true
}

// IsKeyboardOnly returns whether keyboard-only mode is enabled
func (s *AccessibilitySettings) IsKeyboardOnly() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.KeyboardOnly
}

// GetAnimationDuration returns animation duration based on settings
func (s *AccessibilitySettings) GetAnimationDuration(defaultDuration float64) float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.DisableAnimations {
		return 0
	}

	return defaultDuration / s.AnimationSpeed
}

// ShouldShowAnimations returns whether animations should be shown
func (s *AccessibilitySettings) ShouldShowAnimations() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return !s.DisableAnimations
}

// GetFocusIndicatorSize returns the focus indicator size
func (s *AccessibilitySettings) GetFocusIndicatorSize() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.FocusIndicatorSize
}

// HighContrastColors provides color names for high contrast mode
type HighContrastColors struct {
	Background     string
	Foreground     string
	Primary        string
	Secondary      string
	Success        string
	Warning        string
	Error          string
	Info           string
	Border         string
	Focus          string
}

// GetHighContrastColors returns high contrast color palette
func GetHighContrastColors() HighContrastColors {
	return HighContrastColors{
		Background: "#000000", // Black for maximum contrast
		Foreground: "#FFFFFF", // White for maximum contrast
		Primary:    "#00FFFF", // Cyan - bright and distinct
		Secondary:  "#FFFF00", // Yellow - highly visible
		Success:    "#00FF00", // Green - standard success
		Warning:    "#FFFF00", // Yellow - attention grabbing
		Error:      "#FF0000", // Red - clear error
		Info:       "#00FFFF", // Cyan - informative
		Border:     "#FFFFFF", // White for clear boundaries
		Focus:      "#FFFF00", // Yellow for clear focus indicator
	}
}

// ScreenReaderLabel represents a label for screen readers
type ScreenReaderLabel struct {
	Element     string
	Label       string
	Description string
	Hint        string
	State       string
	Role        string
}

// FormatScreenReaderLabel formats a label for screen readers
func FormatScreenReaderLabel(label ScreenReaderLabel) string {
	if !globalSettings.IsScreenReaderMode() {
		return label.Label
	}

	// Verbose format for screen readers
	text := label.Label

	if label.Role != "" {
		text = label.Role + ": " + text
	}

	if label.State != "" {
		text += " (" + label.State + ")"
	}

	if label.Description != "" {
		text += " - " + label.Description
	}

	if label.Hint != "" && globalSettings.ShowKeyboardHints {
		text += " [" + label.Hint + "]"
	}

	return text
}

// KeyboardNavigationHelper provides keyboard navigation assistance
type KeyboardNavigationHelper struct {
	FocusHistory []string
	CurrentFocus string
	FocusRing    []string
}

// NewKeyboardNavigationHelper creates a new keyboard navigation helper
func NewKeyboardNavigationHelper() *KeyboardNavigationHelper {
	return &KeyboardNavigationHelper{
		FocusHistory: make([]string, 0),
		FocusRing:    make([]string, 0),
	}
}

// Focus moves focus to an element
func (h *KeyboardNavigationHelper) Focus(elementID string) {
	h.FocusHistory = append(h.FocusHistory, h.CurrentFocus)
	h.CurrentFocus = elementID
}

// FocusNext moves focus to the next element in the ring
func (h *KeyboardNavigationHelper) FocusNext() string {
	if len(h.FocusRing) == 0 {
		return ""
	}

	currentIndex := -1
	for i, id := range h.FocusRing {
		if id == h.CurrentFocus {
			currentIndex = i
			break
		}
	}

	nextIndex := (currentIndex + 1) % len(h.FocusRing)
	h.Focus(h.FocusRing[nextIndex])
	return h.CurrentFocus
}

// FocusPrevious moves focus to the previous element in the ring
func (h *KeyboardNavigationHelper) FocusPrevious() string {
	if len(h.FocusRing) == 0 {
		return ""
	}

	currentIndex := -1
	for i, id := range h.FocusRing {
		if id == h.CurrentFocus {
			currentIndex = i
			break
		}
	}

	prevIndex := (currentIndex - 1 + len(h.FocusRing)) % len(h.FocusRing)
	h.Focus(h.FocusRing[prevIndex])
	return h.CurrentFocus
}

// FocusBack returns to the previous focus
func (h *KeyboardNavigationHelper) FocusBack() string {
	if len(h.FocusHistory) == 0 {
		return h.CurrentFocus
	}

	h.CurrentFocus = h.FocusHistory[len(h.FocusHistory)-1]
	h.FocusHistory = h.FocusHistory[:len(h.FocusHistory)-1]
	return h.CurrentFocus
}

// SetFocusRing sets the focus ring for tab navigation
func (h *KeyboardNavigationHelper) SetFocusRing(elements []string) {
	h.FocusRing = elements
}

// GetCurrentFocus returns the currently focused element
func (h *KeyboardNavigationHelper) GetCurrentFocus() string {
	return h.CurrentFocus
}

// AccessibilityAnnouncement represents an announcement for screen readers
type AccessibilityAnnouncement struct {
	Message  string
	Priority string // "low", "medium", "high", "assertive"
	Type     string // "info", "success", "warning", "error"
}

// AnnouncementQueue manages screen reader announcements
type AnnouncementQueue struct {
	mu            sync.Mutex
	announcements []AccessibilityAnnouncement
	maxSize       int
}

// NewAnnouncementQueue creates a new announcement queue
func NewAnnouncementQueue() *AnnouncementQueue {
	return &AnnouncementQueue{
		announcements: make([]AccessibilityAnnouncement, 0),
		maxSize:       10,
	}
}

// Announce adds an announcement to the queue
func (q *AnnouncementQueue) Announce(message, priority, typ string) {
	if !globalSettings.IsScreenReaderMode() {
		return
	}

	q.mu.Lock()
	defer q.mu.Unlock()

	announcement := AccessibilityAnnouncement{
		Message:  message,
		Priority: priority,
		Type:     typ,
	}

	q.announcements = append(q.announcements, announcement)
	if len(q.announcements) > q.maxSize {
		q.announcements = q.announcements[1:]
	}
}

// GetAnnouncements returns all pending announcements
func (q *AnnouncementQueue) GetAnnouncements() []AccessibilityAnnouncement {
	q.mu.Lock()
	defer q.mu.Unlock()

	announcements := make([]AccessibilityAnnouncement, len(q.announcements))
	copy(announcements, q.announcements)
	return announcements
}

// Clear clears all announcements
func (q *AnnouncementQueue) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.announcements = make([]AccessibilityAnnouncement, 0)
}

// Global announcement queue
var globalAnnouncementQueue = NewAnnouncementQueue()

// Announce adds a screen reader announcement
func Announce(message, priority, typ string) {
	globalAnnouncementQueue.Announce(message, priority, typ)
}

// GetAnnouncements returns pending screen reader announcements
func GetAnnouncements() []AccessibilityAnnouncement {
	return globalAnnouncementQueue.GetAnnouncements()
}

// ClearAnnouncements clears all announcements
func ClearAnnouncements() {
	globalAnnouncementQueue.Clear()
}

// Common accessibility helpers

// AnnounceSuccess announces a success message
func AnnounceSuccess(message string) {
	Announce(message, "medium", "success")
}

// AnnounceError announces an error message
func AnnounceError(message string) {
	Announce(message, "high", "error")
}

// AnnounceWarning announces a warning message
func AnnounceWarning(message string) {
	Announce(message, "medium", "warning")
}

// AnnounceInfo announces an info message
func AnnounceInfo(message string) {
	Announce(message, "low", "info")
}

// AnnounceNavigation announces navigation changes
func AnnounceNavigation(from, to string) {
	if from != "" {
		Announce("Navigated from "+from+" to "+to, "medium", "info")
	} else {
		Announce("Now at "+to, "medium", "info")
	}
}

// AnnounceFocus announces focus changes
func AnnounceFocus(element string) {
	Announce("Focus on "+element, "low", "info")
}

// AnnounceAction announces user actions
func AnnounceAction(action, result string) {
	Announce(action+" - "+result, "medium", "success")
}

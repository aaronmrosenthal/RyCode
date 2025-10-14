package theme

import (
	"sync"
	"sync/atomic"
	"time"
)

// ThemeTelemetry tracks usage statistics for the theming system.
// This data helps understand how users interact with provider themes.
type ThemeTelemetry struct {
	// Theme switch statistics
	TotalSwitches    atomic.Uint64
	SwitchesByTheme  map[string]*atomic.Uint64
	ActiveTime       map[string]time.Duration
	LastSwitchTime   time.Time
	SessionStartTime time.Time

	// Performance metrics
	AverageSwitchTime time.Duration
	FastestSwitch     time.Duration
	SlowestSwitch     time.Duration

	// Usage patterns
	TabCycles         atomic.Uint64 // User pressed Tab
	ModalSelections   atomic.Uint64 // User selected from modal
	ProgrammaticSwitches atomic.Uint64 // Code-triggered switches

	mu sync.RWMutex
}

// Global telemetry instance
var globalTelemetry *ThemeTelemetry

func init() {
	globalTelemetry = &ThemeTelemetry{
		SwitchesByTheme:  make(map[string]*atomic.Uint64),
		ActiveTime:       make(map[string]time.Duration),
		SessionStartTime: time.Now(),
		LastSwitchTime:   time.Now(),
		FastestSwitch:    time.Hour, // Initialize to large value
	}

	// Initialize counters for known providers
	for _, provider := range []string{"claude", "gemini", "codex", "qwen"} {
		globalTelemetry.SwitchesByTheme[provider] = &atomic.Uint64{}
	}
}

// RecordThemeSwitch records a theme switch event with timing information.
func RecordThemeSwitch(providerID string, switchType SwitchType, duration time.Duration) {
	if globalTelemetry == nil {
		return
	}

	globalTelemetry.mu.Lock()
	defer globalTelemetry.mu.Unlock()

	// Increment total switches
	globalTelemetry.TotalSwitches.Add(1)

	// Increment per-theme counter
	if counter, exists := globalTelemetry.SwitchesByTheme[providerID]; exists {
		counter.Add(1)
	} else {
		newCounter := &atomic.Uint64{}
		newCounter.Add(1)
		globalTelemetry.SwitchesByTheme[providerID] = newCounter
	}

	// Track switch type
	switch switchType {
	case SwitchTypeTab:
		globalTelemetry.TabCycles.Add(1)
	case SwitchTypeModal:
		globalTelemetry.ModalSelections.Add(1)
	case SwitchTypeProgrammatic:
		globalTelemetry.ProgrammaticSwitches.Add(1)
	}

	// Update performance metrics
	if duration < globalTelemetry.FastestSwitch {
		globalTelemetry.FastestSwitch = duration
	}
	if duration > globalTelemetry.SlowestSwitch {
		globalTelemetry.SlowestSwitch = duration
	}

	// Calculate running average switch time
	totalSwitches := globalTelemetry.TotalSwitches.Load()
	if totalSwitches > 0 {
		currentAvg := globalTelemetry.AverageSwitchTime
		globalTelemetry.AverageSwitchTime = ((currentAvg * time.Duration(totalSwitches-1)) + duration) / time.Duration(totalSwitches)
	}

	// Update active time for previous theme
	now := time.Now()
	timeSinceLastSwitch := now.Sub(globalTelemetry.LastSwitchTime)
	if prevTheme := getCurrentProviderID(); prevTheme != "" && prevTheme != providerID {
		globalTelemetry.ActiveTime[prevTheme] += timeSinceLastSwitch
	}

	globalTelemetry.LastSwitchTime = now
}

// SwitchType indicates how the theme switch was triggered
type SwitchType int

const (
	SwitchTypeTab SwitchType = iota
	SwitchTypeModal
	SwitchTypeProgrammatic
)

// GetTelemetryStats returns a snapshot of current telemetry statistics.
func GetTelemetryStats() TelemetryStats {
	if globalTelemetry == nil {
		return TelemetryStats{}
	}

	globalTelemetry.mu.RLock()
	defer globalTelemetry.mu.RUnlock()

	// Calculate session duration
	sessionDuration := time.Since(globalTelemetry.SessionStartTime)

	// Calculate time in current theme
	currentProvider := getCurrentProviderID()
	if currentProvider != "" {
		timeSinceLastSwitch := time.Since(globalTelemetry.LastSwitchTime)
		globalTelemetry.ActiveTime[currentProvider] += timeSinceLastSwitch
		globalTelemetry.LastSwitchTime = time.Now()
	}

	// Copy switch counts
	switchCounts := make(map[string]uint64)
	for provider, counter := range globalTelemetry.SwitchesByTheme {
		switchCounts[provider] = counter.Load()
	}

	// Copy active times
	activeTimes := make(map[string]time.Duration)
	for provider, duration := range globalTelemetry.ActiveTime {
		activeTimes[provider] = duration
	}

	return TelemetryStats{
		TotalSwitches:        globalTelemetry.TotalSwitches.Load(),
		SwitchesByTheme:      switchCounts,
		ActiveTimeByTheme:    activeTimes,
		SessionDuration:      sessionDuration,
		AverageSwitchTime:    globalTelemetry.AverageSwitchTime,
		FastestSwitch:        globalTelemetry.FastestSwitch,
		SlowestSwitch:        globalTelemetry.SlowestSwitch,
		TabCycles:            globalTelemetry.TabCycles.Load(),
		ModalSelections:      globalTelemetry.ModalSelections.Load(),
		ProgrammaticSwitches: globalTelemetry.ProgrammaticSwitches.Load(),
	}
}

// TelemetryStats provides a snapshot of theme usage statistics.
type TelemetryStats struct {
	TotalSwitches        uint64
	SwitchesByTheme      map[string]uint64
	ActiveTimeByTheme    map[string]time.Duration
	SessionDuration      time.Duration
	AverageSwitchTime    time.Duration
	FastestSwitch        time.Duration
	SlowestSwitch        time.Duration
	TabCycles            uint64
	ModalSelections      uint64
	ProgrammaticSwitches uint64
}

// MostUsedTheme returns the provider ID with the most time spent active.
func (s TelemetryStats) MostUsedTheme() string {
	var mostUsed string
	var maxDuration time.Duration

	for provider, duration := range s.ActiveTimeByTheme {
		if duration > maxDuration {
			maxDuration = duration
			mostUsed = provider
		}
	}

	return mostUsed
}

// LeastUsedTheme returns the provider ID with the least time spent active.
func (s TelemetryStats) LeastUsedTheme() string {
	var leastUsed string
	minDuration := time.Duration(1<<63 - 1) // Max duration

	for provider, duration := range s.ActiveTimeByTheme {
		if duration < minDuration && duration > 0 {
			minDuration = duration
			leastUsed = provider
		}
	}

	return leastUsed
}

// ThemePreference returns a score (0.0-1.0) indicating preference for a theme.
// Based on both time spent and number of switches to that theme.
func (s TelemetryStats) ThemePreference(providerID string) float64 {
	if s.SessionDuration == 0 {
		return 0.0
	}

	// Weight: 70% time spent, 30% selection frequency
	timeWeight := 0.7
	selectionWeight := 0.3

	// Time component
	timeScore := 0.0
	if activeTime, exists := s.ActiveTimeByTheme[providerID]; exists {
		timeScore = float64(activeTime) / float64(s.SessionDuration)
	}

	// Selection component
	selectionScore := 0.0
	if s.TotalSwitches > 0 {
		if switches, exists := s.SwitchesByTheme[providerID]; exists {
			selectionScore = float64(switches) / float64(s.TotalSwitches)
		}
	}

	return (timeScore * timeWeight) + (selectionScore * selectionWeight)
}

// ResetTelemetry clears all telemetry data (useful for testing).
func ResetTelemetry() {
	if globalTelemetry == nil {
		return
	}

	globalTelemetry.mu.Lock()
	defer globalTelemetry.mu.Unlock()

	globalTelemetry.TotalSwitches.Store(0)
	globalTelemetry.TabCycles.Store(0)
	globalTelemetry.ModalSelections.Store(0)
	globalTelemetry.ProgrammaticSwitches.Store(0)

	for _, counter := range globalTelemetry.SwitchesByTheme {
		counter.Store(0)
	}

	globalTelemetry.ActiveTime = make(map[string]time.Duration)
	globalTelemetry.SessionStartTime = time.Now()
	globalTelemetry.LastSwitchTime = time.Now()
	globalTelemetry.AverageSwitchTime = 0
	globalTelemetry.FastestSwitch = time.Hour
	globalTelemetry.SlowestSwitch = 0
}

// getCurrentProviderID is a helper to get the current provider ID.
// This would typically access the theme manager's current theme.
func getCurrentProviderID() string {
	theme := CurrentTheme()
	if providerTheme, ok := theme.(*ProviderTheme); ok {
		return providerTheme.ProviderID
	}
	return ""
}

// EnableTelemetry enables telemetry tracking (enabled by default).
var telemetryEnabled atomic.Bool

func init() {
	telemetryEnabled.Store(true)
}

// EnableTelemetryTracking enables telemetry data collection.
func EnableTelemetryTracking() {
	telemetryEnabled.Store(true)
}

// DisableTelemetryTracking disables telemetry data collection.
func DisableTelemetryTracking() {
	telemetryEnabled.Store(false)
}

// IsTelemetryEnabled returns whether telemetry tracking is currently enabled.
func IsTelemetryEnabled() bool {
	return telemetryEnabled.Load()
}

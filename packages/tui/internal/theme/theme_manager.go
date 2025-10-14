package theme

import (
	"sync"
)

// ThemeManager handles dynamic theme switching for provider-specific UI
type ThemeManager struct {
	mu sync.RWMutex

	// Available themes by provider ID
	themes map[string]*ProviderTheme

	// Current active theme
	current *ProviderTheme

	// Previous theme (for transitions)
	previous *ProviderTheme

	// Callbacks for theme change events
	changeCallbacks []func(*ProviderTheme)
}

// NewThemeManager creates a new theme manager with all provider themes loaded
func NewThemeManager() *ThemeManager {
	tm := &ThemeManager{
		themes:          make(map[string]*ProviderTheme),
		changeCallbacks: make([]func(*ProviderTheme), 0),
	}

	// Register all provider themes
	tm.registerTheme(NewClaudeTheme())
	tm.registerTheme(NewGeminiTheme())
	tm.registerTheme(NewCodexTheme())
	tm.registerTheme(NewQwenTheme())

	// Set Claude as default
	tm.current = tm.themes["claude"]

	return tm
}

// registerTheme adds a provider theme to the manager
func (tm *ThemeManager) registerTheme(theme *ProviderTheme) {
	tm.themes[theme.ProviderID] = theme
}

// Current returns the currently active theme
func (tm *ThemeManager) Current() *ProviderTheme {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.current
}

// CurrentTheme returns the current theme as a Theme interface
func (tm *ThemeManager) CurrentTheme() Theme {
	return tm.Current()
}

// Previous returns the previously active theme
func (tm *ThemeManager) Previous() *ProviderTheme {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.previous
}

// SwitchToProvider switches to the theme for the specified provider
// Returns true if theme was changed, false if already active
func (tm *ThemeManager) SwitchToProvider(providerID string) bool {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Get theme for provider
	newTheme, exists := tm.themes[providerID]
	if !exists {
		// Provider not found, keep current theme
		return false
	}

	// Check if already active
	if tm.current != nil && tm.current.ProviderID == providerID {
		return false
	}

	// Store previous theme
	tm.previous = tm.current

	// Switch to new theme
	tm.current = newTheme

	// Notify listeners (outside the lock to avoid deadlocks)
	tm.mu.Unlock()
	tm.notifyThemeChanged(newTheme)
	tm.mu.Lock()

	return true
}

// GetTheme returns the theme for a specific provider
func (tm *ThemeManager) GetTheme(providerID string) (*ProviderTheme, bool) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	theme, exists := tm.themes[providerID]
	return theme, exists
}

// AvailableProviders returns a list of all provider IDs with themes
func (tm *ThemeManager) AvailableProviders() []string {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	providers := make([]string, 0, len(tm.themes))
	for id := range tm.themes {
		providers = append(providers, id)
	}
	return providers
}

// OnThemeChange registers a callback that will be called when the theme changes
func (tm *ThemeManager) OnThemeChange(callback func(*ProviderTheme)) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tm.changeCallbacks = append(tm.changeCallbacks, callback)
}

// notifyThemeChanged calls all registered callbacks with the new theme
func (tm *ThemeManager) notifyThemeChanged(theme *ProviderTheme) {
	tm.mu.RLock()
	callbacks := make([]func(*ProviderTheme), len(tm.changeCallbacks))
	copy(callbacks, tm.changeCallbacks)
	tm.mu.RUnlock()

	// Call callbacks outside the lock
	for _, callback := range callbacks {
		callback(theme)
	}
}

// Reset resets the theme manager to default state (Claude theme)
func (tm *ThemeManager) Reset() {
	tm.SwitchToProvider("claude")
}

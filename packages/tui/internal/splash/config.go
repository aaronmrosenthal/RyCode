package splash

import (
	"encoding/json"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

func init() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())
}

// Config represents splash screen configuration
type Config struct {
	SplashEnabled   bool   `json:"splash_enabled"`
	SplashFrequency string `json:"splash_frequency"` // "always", "first", "random", "never"
	ReducedMotion   bool   `json:"reduced_motion"`
	ColorMode       string `json:"color_mode"` // "truecolor", "256", "16", "auto"
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		SplashEnabled:   true,
		SplashFrequency: "always", // THE AWAKENING - always show the epic splash
		ReducedMotion:   false,
		ColorMode:       "auto",
	}
}

// getConfigPath returns the path to the config file (variable for testing)
var getConfigPath = func() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".rycode", "config.json")
}

// getMarkerPath returns the path to the splash shown marker file (variable for testing)
var getMarkerPath = func() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".rycode", ".splash_shown")
}

// IsFirstRun checks if this is the first time RyCode is being run
func IsFirstRun() bool {
	_, err := os.Stat(getMarkerPath())
	return os.IsNotExist(err)
}

// MarkAsShown creates a marker file to indicate splash has been shown
func MarkAsShown() error {
	path := getMarkerPath()
	dir := filepath.Dir(path)

	// Ensure directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Create marker file
	return os.WriteFile(path, []byte("shown"), 0644)
}

// LoadConfig loads the configuration from disk
func LoadConfig() (*Config, error) {
	path := getConfigPath()

	// If config doesn't exist, return defaults
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	// Read config file
	data, err := os.ReadFile(path)
	if err != nil {
		return DefaultConfig(), nil
	}

	// Parse JSON
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		// If parse fails, return defaults (graceful degradation)
		return DefaultConfig(), nil
	}

	// Override with system preferences
	if os.Getenv("PREFERS_REDUCED_MOTION") == "1" {
		config.ReducedMotion = true
	}

	if os.Getenv("NO_COLOR") != "" {
		config.ColorMode = "16"
	}

	return &config, nil
}

// Save saves the configuration to disk
func (c *Config) Save() error {
	path := getConfigPath()
	dir := filepath.Dir(path)

	// Ensure directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(path, data, 0644)
}

// ShouldShowSplash determines whether to show the splash screen
func ShouldShowSplash(config *Config) bool {
	// User disabled splash
	if !config.SplashEnabled {
		return false
	}

	// Respect reduced motion preference
	if config.ReducedMotion {
		return false
	}

	// First run always shows (unless explicitly disabled)
	if IsFirstRun() {
		return true
	}

	// Check frequency setting
	switch config.SplashFrequency {
	case "always":
		return true
	case "never":
		return false
	case "first":
		return false // Already shown once
	case "random":
		// 10% random chance
		return rand.Float64() < 0.1
	default:
		// Default to "first" behavior
		return false
	}
}

// DisableSplashPermanently disables the splash screen in config
func DisableSplashPermanently() error {
	config, err := LoadConfig()
	if err != nil {
		config = DefaultConfig()
	}

	config.SplashEnabled = false
	return config.Save()
}

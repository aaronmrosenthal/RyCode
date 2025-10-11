package splash

import (
	"os"
	"runtime"
	"strings"

	"golang.org/x/term"
)

// TerminalCapabilities represents what the terminal can do
type TerminalCapabilities struct {
	Width       int
	Height      int
	Colors      ColorMode
	Unicode     bool
	TooSmall    bool
	Performance string // "fast", "medium", "slow"
}

// DetectTerminalCapabilities detects what the terminal can handle
func DetectTerminalCapabilities() TerminalCapabilities {
	caps := TerminalCapabilities{
		Width:       80,
		Height:      24,
		Colors:      Truecolor,
		Unicode:     true,
		Performance: "fast",
	}

	// Get terminal size
	if width, height, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		caps.Width = width
		caps.Height = height
	}

	// Check if terminal is too small
	if caps.Width < 80 || caps.Height < 24 {
		caps.TooSmall = true
	}

	// Detect color support
	caps.Colors = DetectColorMode()

	// Detect unicode support
	caps.Unicode = SupportsUnicode()

	// Estimate performance (conservative)
	caps.Performance = EstimatePerformance()

	return caps
}

// DetectColorMode detects terminal color capabilities
func DetectColorMode() ColorMode {
	colorterm := os.Getenv("COLORTERM")
	if colorterm == "truecolor" || colorterm == "24bit" {
		return Truecolor
	}

	term := os.Getenv("TERM")
	if strings.Contains(term, "256color") {
		return Colors256
	}

	// Check for NO_COLOR environment variable
	if os.Getenv("NO_COLOR") != "" {
		return Colors16
	}

	return Colors16 // Conservative default
}

// SupportsUnicode checks if the terminal supports unicode
func SupportsUnicode() bool {
	// Windows CMD has limited unicode support
	if runtime.GOOS == "windows" && os.Getenv("WT_SESSION") == "" {
		return false
	}

	// Check locale
	lang := os.Getenv("LANG")
	if lang != "" && (strings.Contains(strings.ToLower(lang), "utf") ||
		strings.Contains(strings.ToLower(lang), "utf-8")) {
		return true
	}

	// Default to true on modern systems
	return true
}

// EstimatePerformance estimates terminal rendering performance
func EstimatePerformance() string {
	// Check if running in remote session (likely slower)
	if os.Getenv("SSH_CONNECTION") != "" {
		return "medium"
	}

	// Windows Terminal and iTerm2 are generally fast
	if os.Getenv("WT_SESSION") != "" || os.Getenv("TERM_PROGRAM") == "iTerm.app" {
		return "fast"
	}

	// Conservative default
	return "medium"
}

// ShouldUseFallback determines if we should use a simplified mode
func (caps TerminalCapabilities) ShouldUseFallback() bool {
	return caps.TooSmall || caps.Colors == Colors16 || !caps.Unicode
}

// ShouldSkipSplash determines if splash should be skipped entirely
func (caps TerminalCapabilities) ShouldSkipSplash() bool {
	// Skip on extremely small terminals
	return caps.Width < 60 || caps.Height < 20
}

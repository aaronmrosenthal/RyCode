package responsive

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/term"
)

// TerminalCapabilities represents detected terminal features
type TerminalCapabilities struct {
	// Dimensions
	Width        int
	Height       int
	WidthPixels  int // If available
	HeightPixels int // If available

	// Colors
	SupportsTrueColor bool
	Supports256Color  bool
	Supports16Color   bool

	// Mouse/Touch
	SupportsMouseTracking bool
	SupportsSGRMouse      bool // SGR extended mouse mode
	SupportsPixelMouse    bool // Pixel-based mouse coordinates

	// Advanced features
	SupportsAltScreen   bool
	SupportsBracketPaste bool
	SupportsKittyGraphics bool
	SupportsSixel        bool

	// Terminal info
	TerminalType    string
	TerminalProgram string
	IsSSH           bool
	IsTmux          bool
	IsScreen        bool

	// Platform detection
	Platform    string // "ios", "android", "macos", "linux", "windows"
	IsPhone     bool
	IsTablet    bool
	IsMobile    bool
	IsTouchDevice bool
}

// DetectCapabilities detects terminal capabilities
func DetectCapabilities() *TerminalCapabilities {
	caps := &TerminalCapabilities{
		TerminalType:    os.Getenv("TERM"),
		TerminalProgram: detectTerminalProgram(),
	}

	// Detect dimensions
	caps.detectDimensions()

	// Detect color support
	caps.detectColorSupport()

	// Detect mouse support
	caps.detectMouseSupport()

	// Detect terminal multiplexers
	caps.IsTmux = os.Getenv("TMUX") != ""
	caps.IsScreen = strings.HasPrefix(caps.TerminalType, "screen")

	// Detect SSH
	caps.IsSSH = os.Getenv("SSH_CONNECTION") != "" || os.Getenv("SSH_CLIENT") != ""

	// Detect platform and device type
	caps.detectPlatform()

	// Detect advanced features
	caps.detectAdvancedFeatures()

	return caps
}

// detectDimensions detects terminal dimensions
func (tc *TerminalCapabilities) detectDimensions() {
	// Get character dimensions
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err == nil {
		tc.Width = width
		tc.Height = height
	} else {
		// Fallback to env vars
		tc.Width = getEnvInt("COLUMNS", 80)
		tc.Height = getEnvInt("LINES", 24)
	}

	// Try to detect pixel dimensions (kitty, iTerm2)
	tc.detectPixelDimensions()
}

// detectPixelDimensions attempts to get pixel dimensions
func (tc *TerminalCapabilities) detectPixelDimensions() {
	// kitty graphics protocol - query terminal size in pixels
	// CSI 14 t - Reports window size in pixels
	// Response: CSI 4 ; height ; width t

	// For now, estimate based on common font sizes
	// Most terminals use ~7-9 pixels per char width, ~14-18 per char height
	avgCharWidth := 8
	avgCharHeight := 16

	tc.WidthPixels = tc.Width * avgCharWidth
	tc.HeightPixels = tc.Height * avgCharHeight

	// TODO: Actually query terminal for real pixel dimensions
	// This requires sending escape sequences and reading responses
}

// detectColorSupport detects color capabilities
func (tc *TerminalCapabilities) detectColorSupport() {
	colorterm := os.Getenv("COLORTERM")

	// True color detection
	tc.SupportsTrueColor = colorterm == "truecolor" || colorterm == "24bit" ||
		strings.Contains(tc.TerminalType, "256color") ||
		tc.TerminalProgram == "iTerm.app" ||
		tc.TerminalProgram == "WezTerm" ||
		tc.TerminalProgram == "Alacritty" ||
		tc.TerminalProgram == "kitty"

	// 256 color support
	tc.Supports256Color = tc.SupportsTrueColor ||
		strings.Contains(tc.TerminalType, "256color") ||
		colorterm == "256"

	// 16 color support (almost all terminals)
	tc.Supports16Color = tc.TerminalType != "" && tc.TerminalType != "dumb"
}

// detectMouseSupport detects mouse capabilities
func (tc *TerminalCapabilities) detectMouseSupport() {
	// Most modern terminals support mouse tracking
	tc.SupportsMouseTracking = !tc.IsSSH && // SSH often blocks mouse
		tc.TerminalType != "dumb" &&
		!strings.HasPrefix(tc.TerminalType, "vt") // Old VT terminals

	// SGR extended mouse mode (1006) - modern standard
	tc.SupportsSGRMouse = tc.SupportsMouseTracking &&
		(tc.TerminalProgram == "iTerm.app" ||
			tc.TerminalProgram == "Alacritty" ||
			tc.TerminalProgram == "kitty" ||
			tc.TerminalProgram == "WezTerm" ||
			tc.TerminalProgram == "VSCode" ||
			strings.Contains(tc.TerminalType, "xterm"))

	// Pixel-based mouse coordinates (some terminals)
	tc.SupportsPixelMouse = tc.TerminalProgram == "kitty"
}

// detectPlatform detects OS platform and device type
func (tc *TerminalCapabilities) detectPlatform() {
	// Check environment variables for iOS/Android terminal apps
	termProgram := strings.ToLower(tc.TerminalProgram)
	termType := strings.ToLower(tc.TerminalType)

	// iOS terminal apps
	if strings.Contains(termProgram, "blink") ||
		strings.Contains(termProgram, "termius") ||
		strings.Contains(termProgram, "issh") ||
		os.Getenv("LC_TERMINAL") == "Blink" {
		tc.Platform = "ios"
		tc.IsMobile = true
		tc.IsPhone = tc.Width < 120 // Heuristic
		tc.IsTablet = !tc.IsPhone
		tc.IsTouchDevice = true
		return
	}

	// Android terminal apps
	if strings.Contains(termProgram, "termux") ||
		os.Getenv("TERMUX_VERSION") != "" ||
		strings.Contains(os.Getenv("PREFIX"), "termux") {
		tc.Platform = "android"
		tc.IsMobile = true
		tc.IsPhone = tc.Width < 120
		tc.IsTablet = !tc.IsPhone
		tc.IsTouchDevice = true
		return
	}

	// Desktop platforms
	switch {
	case strings.Contains(strings.ToLower(os.Getenv("OS")), "windows"):
		tc.Platform = "windows"
	case fileExists("/System/Library/CoreServices/SystemVersion.plist"):
		tc.Platform = "macos"
	default:
		tc.Platform = "linux"
	}

	// Desktop devices are not mobile
	tc.IsMobile = false
	tc.IsPhone = false
	tc.IsTablet = false

	// Some desktop devices have touch (Surface, MacBook with trackpad gestures)
	tc.IsTouchDevice = false
}

// detectAdvancedFeatures detects advanced terminal features
func (tc *TerminalCapabilities) detectAdvancedFeatures() {
	// Alt screen support (most modern terminals)
	tc.SupportsAltScreen = tc.TerminalType != "dumb" &&
		!strings.HasPrefix(tc.TerminalType, "vt")

	// Bracketed paste mode
	tc.SupportsBracketPaste = tc.SupportsAltScreen // Same capability set usually

	// Kitty graphics protocol
	tc.SupportsKittyGraphics = tc.TerminalProgram == "kitty" ||
		tc.TerminalProgram == "WezTerm"

	// Sixel graphics
	tc.SupportsSixel = tc.TerminalProgram == "WezTerm" ||
		tc.TerminalProgram == "mlterm" ||
		strings.Contains(termType, "sixel")
}

// detectTerminalProgram detects which terminal program is running
func detectTerminalProgram() string {
	// Check various environment variables
	if prog := os.Getenv("TERM_PROGRAM"); prog != "" {
		return prog
	}

	if prog := os.Getenv("LC_TERMINAL"); prog != "" {
		return prog
	}

	// Check for specific terminals
	if os.Getenv("KITTY_WINDOW_ID") != "" {
		return "kitty"
	}

	if os.Getenv("ALACRITTY_SOCKET") != "" || os.Getenv("ALACRITTY_LOG") != "" {
		return "Alacritty"
	}

	if os.Getenv("WEZTERM_EXECUTABLE") != "" {
		return "WezTerm"
	}

	if os.Getenv("VSCODE_INJECTION") != "" {
		return "VSCode"
	}

	if os.Getenv("TERMUX_VERSION") != "" {
		return "Termux"
	}

	// Default to TERM value
	return os.Getenv("TERM")
}

// EnableMouseTracking enables mouse tracking in terminal
func (tc *TerminalCapabilities) EnableMouseTracking() {
	if !tc.SupportsMouseTracking {
		return
	}

	if tc.SupportsSGRMouse {
		// Enable SGR extended mouse mode (1006)
		// Also enable mouse button tracking (1002) and any-event tracking (1003)
		fmt.Print("\x1b[?1002h") // Button event tracking
		fmt.Print("\x1b[?1006h") // SGR extended mode
	} else {
		// Fallback to basic mouse tracking (1000)
		fmt.Print("\x1b[?1000h")
	}
}

// DisableMouseTracking disables mouse tracking
func (tc *TerminalCapabilities) DisableMouseTracking() {
	if !tc.SupportsMouseTracking {
		return
	}

	fmt.Print("\x1b[?1002l")
	fmt.Print("\x1b[?1006l")
	fmt.Print("\x1b[?1000l")
}

// GetDeviceType returns simplified device type based on capabilities
func (tc *TerminalCapabilities) GetDeviceType() DeviceType {
	if tc.IsPhone {
		return DevicePhone
	}
	if tc.IsTablet {
		return DeviceTablet
	}
	return DeviceDesktop
}

// GetOrientation estimates orientation based on dimensions
func (tc *TerminalCapabilities) GetOrientation() Orientation {
	if tc.Width > tc.Height {
		return OrientationLandscape
	}
	return OrientationPortrait
}

// ShouldUseTouch returns whether touch interactions should be primary
func (tc *TerminalCapabilities) ShouldUseTouch() bool {
	return tc.IsTouchDevice && tc.IsMobile
}

// ShouldUseMouse returns whether mouse interactions are available
func (tc *TerminalCapabilities) ShouldUseMouse() bool {
	return tc.SupportsMouseTracking && !tc.IsMobile
}

// GetOptimalBreakpoint returns the best breakpoint for current terminal
func (tc *TerminalCapabilities) GetOptimalBreakpoint() Breakpoint {
	width := tc.Width
	orientation := tc.GetOrientation()

	switch {
	case width <= PhonePortrait.MaxWidth:
		return Breakpoint{
			Device:      DevicePhone,
			Orientation: orientation,
			MinWidth:    PhonePortrait.MinWidth,
			MaxWidth:    PhonePortrait.MaxWidth,
		}

	case width <= PhoneLandscape.MaxWidth:
		return Breakpoint{
			Device:      DevicePhone,
			Orientation: OrientationLandscape,
			MinWidth:    PhoneLandscape.MinWidth,
			MaxWidth:    PhoneLandscape.MaxWidth,
		}

	case width <= TabletPortrait.MaxWidth:
		return TabletPortrait

	case width <= TabletLandscape.MaxWidth:
		return TabletLandscape

	default:
		return Desktop
	}
}

// String returns human-readable description
func (tc *TerminalCapabilities) String() string {
	return fmt.Sprintf(
		"Terminal: %s (%s)\n"+
			"Dimensions: %dx%d chars (%dx%d px)\n"+
			"Platform: %s (mobile: %v, touch: %v)\n"+
			"Colors: true=%v, 256=%v\n"+
			"Mouse: tracking=%v, SGR=%v\n"+
			"Device: %s %s",
		tc.TerminalProgram,
		tc.TerminalType,
		tc.Width, tc.Height,
		tc.WidthPixels, tc.HeightPixels,
		tc.Platform, tc.IsMobile, tc.IsTouchDevice,
		tc.SupportsTrueColor, tc.Supports256Color,
		tc.SupportsMouseTracking, tc.SupportsSGRMouse,
		tc.GetDeviceType(), tc.GetOrientation(),
	)
}

// DebugReport generates detailed debug information
func (tc *TerminalCapabilities) DebugReport() string {
	lines := []string{
		"=== Terminal Capabilities Debug Report ===",
		"",
		"Environment:",
		fmt.Sprintf("  TERM=%s", tc.TerminalType),
		fmt.Sprintf("  TERM_PROGRAM=%s", tc.TerminalProgram),
		fmt.Sprintf("  COLORTERM=%s", os.Getenv("COLORTERM")),
		fmt.Sprintf("  SSH=%v", tc.IsSSH),
		fmt.Sprintf("  TMUX=%v", tc.IsTmux),
		"",
		"Dimensions:",
		fmt.Sprintf("  Size: %dx%d characters", tc.Width, tc.Height),
		fmt.Sprintf("  Pixels: %dx%d (estimated)", tc.WidthPixels, tc.HeightPixels),
		fmt.Sprintf("  Orientation: %s", tc.GetOrientation()),
		"",
		"Colors:",
		fmt.Sprintf("  True Color: %v", tc.SupportsTrueColor),
		fmt.Sprintf("  256 Color: %v", tc.Supports256Color),
		fmt.Sprintf("  16 Color: %v", tc.Supports16Color),
		"",
		"Input:",
		fmt.Sprintf("  Mouse Tracking: %v", tc.SupportsMouseTracking),
		fmt.Sprintf("  SGR Mouse: %v", tc.SupportsSGRMouse),
		fmt.Sprintf("  Pixel Mouse: %v", tc.SupportsPixelMouse),
		fmt.Sprintf("  Bracketed Paste: %v", tc.SupportsBracketPaste),
		"",
		"Platform:",
		fmt.Sprintf("  OS: %s", tc.Platform),
		fmt.Sprintf("  Mobile: %v", tc.IsMobile),
		fmt.Sprintf("  Phone: %v", tc.IsPhone),
		fmt.Sprintf("  Tablet: %v", tc.IsTablet),
		fmt.Sprintf("  Touch Device: %v", tc.IsTouchDevice),
		fmt.Sprintf("  Device Type: %s", tc.GetDeviceType()),
		"",
		"Advanced:",
		fmt.Sprintf("  Alt Screen: %v", tc.SupportsAltScreen),
		fmt.Sprintf("  Kitty Graphics: %v", tc.SupportsKittyGraphics),
		fmt.Sprintf("  Sixel: %v", tc.SupportsSixel),
		"",
		fmt.Sprintf("Recommended Breakpoint: %s %s",
			tc.GetOptimalBreakpoint().Device,
			tc.GetOptimalBreakpoint().Orientation,
		),
	}

	return strings.Join(lines, "\n")
}

// Helper functions

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// MouseEvent represents a parsed mouse event
type MouseEvent struct {
	X      int
	Y      int
	Button int // 0=left, 1=middle, 2=right
	Action MouseAction
	Mods   MouseMods
}

// MouseAction represents mouse event type
type MouseAction int

const (
	MousePress MouseAction = iota
	MouseRelease
	MouseMove
	MouseDrag
	MouseScrollUp
	MouseScrollDown
)

// MouseMods represents modifier keys
type MouseMods struct {
	Shift bool
	Alt   bool
	Ctrl  bool
}

// ParseMouseEvent parses SGR mouse event
// Format: CSI < Cb ; Cx ; Cy (M or m)
// M = press, m = release
func ParseMouseEvent(seq string) (*MouseEvent, error) {
	// Remove CSI prefix and parse
	// Example: "\x1b[<0;10;5M" -> button=0, x=10, y=5, press

	if !strings.HasPrefix(seq, "\x1b[<") {
		return nil, fmt.Errorf("not an SGR mouse event")
	}

	seq = strings.TrimPrefix(seq, "\x1b[<")
	isPress := strings.HasSuffix(seq, "M")
	seq = strings.TrimSuffix(strings.TrimSuffix(seq, "M"), "m")

	parts := strings.Split(seq, ";")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid mouse event format")
	}

	button, _ := strconv.Atoi(parts[0])
	x, _ := strconv.Atoi(parts[1])
	y, _ := strconv.Atoi(parts[2])

	event := &MouseEvent{
		X:      x - 1, // Convert to 0-based
		Y:      y - 1,
		Button: button & 3, // Lower 2 bits
	}

	// Determine action
	if isPress {
		event.Action = MousePress
	} else {
		event.Action = MouseRelease
	}

	// Parse modifiers (bits 2-4)
	event.Mods.Shift = (button & 4) != 0
	event.Mods.Alt = (button & 8) != 0
	event.Mods.Ctrl = (button & 16) != 0

	// Check for scroll
	if button&64 != 0 {
		if button&1 != 0 {
			event.Action = MouseScrollDown
		} else {
			event.Action = MouseScrollUp
		}
	}

	return event, nil
}

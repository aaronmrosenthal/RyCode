package responsive

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea/v2"
)

// DeviceType represents the type of device
type DeviceType string

const (
	DevicePhone   DeviceType = "phone"
	DeviceTablet  DeviceType = "tablet"
	DeviceDesktop DeviceType = "desktop"
)

// Orientation represents screen orientation
type Orientation string

const (
	OrientationPortrait  Orientation = "portrait"
	OrientationLandscape Orientation = "landscape"
)

// Breakpoint defines responsive breakpoints
type Breakpoint struct {
	MinWidth    int
	MaxWidth    int
	Device      DeviceType
	Orientation Orientation
}

// Standard breakpoints optimized for CLI usage
var (
	// Phone Portrait: 0-428px (iPhone 14 Pro Max max width)
	PhonePortrait = Breakpoint{
		MinWidth:    0,
		MaxWidth:    60, // ~60 chars typical terminal width on phone
		Device:      DevicePhone,
		Orientation: OrientationPortrait,
	}

	// Phone Landscape: 60-120 chars
	PhoneLandscape = Breakpoint{
		MinWidth:    61,
		MaxWidth:    120,
		Device:      DevicePhone,
		Orientation: OrientationLandscape,
	}

	// Tablet Portrait: 120-180 chars
	TabletPortrait = Breakpoint{
		MinWidth:    121,
		MaxWidth:    180,
		Device:      DeviceTablet,
		Orientation: OrientationPortrait,
	}

	// Tablet Landscape: 180-240 chars
	TabletLandscape = Breakpoint{
		MinWidth:    181,
		MaxWidth:    240,
		Device:      DeviceTablet,
		Orientation: OrientationLandscape,
	}

	// Desktop: 240+ chars
	Desktop = Breakpoint{
		MinWidth:    241,
		MaxWidth:    999999,
		Device:      DeviceDesktop,
		Orientation: OrientationLandscape,
	}
)

// LayoutConfig defines layout configuration for each breakpoint
type LayoutConfig struct {
	Device         DeviceType
	Orientation    Orientation
	Width          int
	Height         int
	ShowSidebar    bool
	ShowTimeline   bool
	ShowHistory    bool
	InputPosition  InputPosition
	MessageLayout  MessageLayoutType
	GesturesActive bool
	UseVoice       bool
	HapticFeedback bool
	ShowReactions  bool
}

// InputPosition defines where input appears
type InputPosition string

const (
	InputBottom InputPosition = "bottom" // Traditional (desktop)
	InputTop    InputPosition = "top"    // Phone-friendly (thumb zone)
	InputFloat  InputPosition = "float"  // Floating bottom-right (tablet)
)

// MessageLayoutType defines how messages are displayed
type MessageLayoutType string

const (
	LayoutCards      MessageLayoutType = "cards"      // Card-based (phone)
	LayoutList       MessageLayoutType = "list"       // Traditional list (desktop)
	LayoutSplit      MessageLayoutType = "split"      // Split view (tablet landscape)
	LayoutBubbles    MessageLayoutType = "bubbles"    // Chat bubbles (phone)
	LayoutTimeline   MessageLayoutType = "timeline"   // Timeline view (tablet)
)

// ViewportManager manages responsive layouts
type ViewportManager struct {
	width       int
	height      int
	config      LayoutConfig
	prevDevice  DeviceType
}

// NewViewportManager creates a responsive viewport manager
func NewViewportManager() *ViewportManager {
	return &ViewportManager{
		prevDevice: DeviceDesktop,
	}
}

// Update handles window size changes
func (vm *ViewportManager) Update(msg tea.WindowSizeMsg) LayoutConfig {
	vm.width = msg.Width
	vm.height = msg.Height
	vm.config = vm.determineLayout()
	return vm.config
}

// determineLayout calculates optimal layout for current dimensions
func (vm *ViewportManager) determineLayout() LayoutConfig {
	breakpoint := vm.getBreakpoint()

	config := LayoutConfig{
		Device:      breakpoint.Device,
		Orientation: breakpoint.Orientation,
		Width:       vm.width,
		Height:      vm.height,
	}

	switch breakpoint.Device {
	case DevicePhone:
		config = vm.phoneLayout(config, breakpoint.Orientation)
	case DeviceTablet:
		config = vm.tabletLayout(config, breakpoint.Orientation)
	case DeviceDesktop:
		config = vm.desktopLayout(config)
	}

	return config
}

// phoneLayout optimizes for phone screens (THE KILLER FEATURE)
func (vm *ViewportManager) phoneLayout(config LayoutConfig, orientation Orientation) LayoutConfig {
	if orientation == OrientationPortrait {
		// PORTRAIT: Thumb-zone optimized
		config.ShowSidebar = false
		config.ShowTimeline = false
		config.ShowHistory = false
		config.InputPosition = InputTop // TOP for thumb reach!
		config.MessageLayout = LayoutBubbles // Chat-style
		config.GesturesActive = true // Swipe gestures
		config.UseVoice = true // Voice input button
		config.HapticFeedback = true // Simulated haptic
		config.ShowReactions = true // Quick emoji reactions
	} else {
		// LANDSCAPE: More screen space
		config.ShowSidebar = false
		config.ShowTimeline = true // Compact timeline
		config.ShowHistory = false
		config.InputPosition = InputBottom
		config.MessageLayout = LayoutCards
		config.GesturesActive = true
		config.UseVoice = true
		config.HapticFeedback = true
		config.ShowReactions = true
	}

	return config
}

// tabletLayout optimizes for tablet screens
func (vm *ViewportManager) tabletLayout(config LayoutConfig, orientation Orientation) LayoutConfig {
	if orientation == OrientationPortrait {
		// PORTRAIT: Single column with context
		config.ShowSidebar = true // Collapsible sidebar
		config.ShowTimeline = true
		config.ShowHistory = true
		config.InputPosition = InputFloat // Floating input
		config.MessageLayout = LayoutTimeline
		config.GesturesActive = true
		config.UseVoice = false
		config.HapticFeedback = false
		config.ShowReactions = true
	} else {
		// LANDSCAPE: Split view power user
		config.ShowSidebar = true
		config.ShowTimeline = true
		config.ShowHistory = true
		config.InputPosition = InputBottom
		config.MessageLayout = LayoutSplit // Chat + code preview
		config.GesturesActive = true
		config.UseVoice = false
		config.HapticFeedback = false
		config.ShowReactions = true
	}

	return config
}

// desktopLayout optimizes for desktop screens
func (vm *ViewportManager) desktopLayout(config LayoutConfig) LayoutConfig {
	// DESKTOP: Full power, all features
	config.ShowSidebar = true
	config.ShowTimeline = true
	config.ShowHistory = true
	config.InputPosition = InputBottom
	config.MessageLayout = LayoutList
	config.GesturesActive = false // Keyboard-focused
	config.UseVoice = false
	config.HapticFeedback = false
	config.ShowReactions = true

	return config
}

// getBreakpoint determines current breakpoint
func (vm *ViewportManager) getBreakpoint() Breakpoint {
	width := vm.width
	height := vm.height

	// Determine orientation
	orientation := OrientationPortrait
	if width > height {
		orientation = OrientationLandscape
	}

	// Match breakpoint
	switch {
	case width <= PhonePortrait.MaxWidth:
		return Breakpoint{
			Device:      DevicePhone,
			Orientation: OrientationPortrait,
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
		if orientation == OrientationPortrait {
			return TabletPortrait
		}
		// Wide enough for landscape features
		return Breakpoint{
			Device:      DeviceTablet,
			Orientation: OrientationLandscape,
			MinWidth:    TabletPortrait.MinWidth,
			MaxWidth:    TabletPortrait.MaxWidth,
		}

	case width <= TabletLandscape.MaxWidth:
		return TabletLandscape

	default:
		return Desktop
	}
}

// GetConfig returns current layout config
func (vm *ViewportManager) GetConfig() LayoutConfig {
	return vm.config
}

// IsPhone returns true if current device is phone
func (vm *ViewportManager) IsPhone() bool {
	return vm.config.Device == DevicePhone
}

// IsTablet returns true if current device is tablet
func (vm *ViewportManager) IsTablet() bool {
	return vm.config.Device == DeviceTablet
}

// IsDesktop returns true if current device is desktop
func (vm *ViewportManager) IsDesktop() bool {
	return vm.config.Device == DeviceDesktop
}

// DeviceChanged returns true if device type changed
func (vm *ViewportManager) DeviceChanged() bool {
	changed := vm.config.Device != vm.prevDevice
	vm.prevDevice = vm.config.Device
	return changed
}

// String returns human-readable description
func (lc LayoutConfig) String() string {
	return fmt.Sprintf("%s (%s) %dx%d - Layout: %s, Input: %s",
		lc.Device,
		lc.Orientation,
		lc.Width,
		lc.Height,
		lc.MessageLayout,
		lc.InputPosition,
	)
}

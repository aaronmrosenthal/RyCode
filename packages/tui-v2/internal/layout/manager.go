package layout

import (
	"time"
)

// LayoutManager handles responsive layout detection and management
type LayoutManager struct {
	width      int
	height     int
	class      DeviceClass
	lastUpdate time.Time
	onChange   func(DeviceClass)
}

// NewLayoutManager creates a new layout manager
func NewLayoutManager(width, height int) *LayoutManager {
	lm := &LayoutManager{
		width:      width,
		height:     height,
		lastUpdate: time.Now(),
	}
	lm.detectDevice()
	return lm
}

// OnChange registers a callback for when device class changes
func (lm *LayoutManager) OnChange(callback func(DeviceClass)) {
	lm.onChange = callback
}

// detectDevice determines the device class based on terminal dimensions
// Optimized for education: fine-grained breakpoints for iPhone and iPad devices
func (lm *LayoutManager) detectDevice() {
	oldClass := lm.class

	switch {
	case lm.width >= 160:
		lm.class = DesktopLarge
	case lm.width >= 140:
		lm.class = LaptopStandard
	case lm.width >= 120:
		lm.class = LaptopSmall
	case lm.width >= 100:
		lm.class = TabletLarge // iPad landscape, iPad Pro portrait
	case lm.width >= 85:
		lm.class = TabletMedium // iPad portrait
	case lm.width >= 70:
		lm.class = TabletSmall // iPad Mini portrait
	case lm.width >= 55:
		lm.class = PhoneStandard // iPhone 12-14 portrait
	case lm.width >= 40:
		lm.class = PhoneCompact // iPhone SE, Mini portrait
	default:
		lm.class = PhoneTiny // Very small or large accessibility font
	}

	// Trigger callback if device class changed
	if oldClass != lm.class && lm.onChange != nil {
		lm.onChange(lm.class)
	}
}

// Update updates dimensions and re-detects device class
func (lm *LayoutManager) Update(width, height int) {
	lm.width = width
	lm.height = height
	lm.detectDevice()
	lm.lastUpdate = time.Now()
}

// GetDeviceClass returns the current device class
func (lm *LayoutManager) GetDeviceClass() DeviceClass {
	return lm.class
}

// GetDimensions returns current width and height
func (lm *LayoutManager) GetDimensions() (width, height int) {
	return lm.width, lm.height
}

// GetWidth returns current width
func (lm *LayoutManager) GetWidth() int {
	return lm.width
}

// GetHeight returns current height
func (lm *LayoutManager) GetHeight() int {
	return lm.height
}

// GetLastUpdate returns when dimensions were last updated
func (lm *LayoutManager) GetLastUpdate() time.Time {
	return lm.lastUpdate
}

// ShouldUseStackLayout returns true if stack layout is recommended
func (lm *LayoutManager) ShouldUseStackLayout() bool {
	return lm.class.IsMobile()
}

// ShouldUseSplitLayout returns true if split layout is recommended
func (lm *LayoutManager) ShouldUseSplitLayout() bool {
	return lm.class.IsTablet()
}

// ShouldUseMultiPaneLayout returns true if multi-pane layout is recommended
func (lm *LayoutManager) ShouldUseMultiPaneLayout() bool {
	return lm.class.IsDesktop()
}

// GetRecommendedSplitRatio returns the recommended split ratio for current device
// Returns the fraction of width for main pane (chat) vs file tree
func (lm *LayoutManager) GetRecommendedSplitRatio() float64 {
	switch lm.class {
	case PhoneTiny, PhoneCompact, PhoneStandard:
		return 1.0 // 100% chat (no file tree by default)
	case TabletSmall:
		return 0.75 // 75% chat / 25% file tree (iPad Mini)
	case TabletMedium:
		return 0.70 // 70% chat / 30% file tree (iPad)
	case TabletLarge:
		return 0.72 // 72% chat / 28% file tree (iPad landscape)
	case LaptopSmall:
		return 0.75 // 75% chat / 25% file tree (Chromebook)
	case LaptopStandard, DesktopLarge:
		return 0.75 // 75% chat / 25% file tree (Desktop)
	default:
		return 0.7
	}
}

// GetFileTreeWidthForDevice returns the recommended file tree width
func (lm *LayoutManager) GetFileTreeWidthForDevice() int {
	return lm.class.GetFileTreeWidth()
}

// ShouldShowFileTree returns true if file tree should be visible by default
func (lm *LayoutManager) ShouldShowFileTree() bool {
	return lm.class.SupportsFileTree()
}

// ShouldUseFileTreeOverlay returns true if file tree should be an overlay/drawer
func (lm *LayoutManager) ShouldUseFileTreeOverlay() bool {
	return lm.class.ShouldShowFileTreeOverlay()
}

// CanFitWidth returns true if the given width can fit in current dimensions
func (lm *LayoutManager) CanFitWidth(requiredWidth int) bool {
	return lm.width >= requiredWidth
}

// CanFitHeight returns true if the given height can fit in current dimensions
func (lm *LayoutManager) CanFitHeight(requiredHeight int) bool {
	return lm.height >= requiredHeight
}

// GetSafeWidth returns width minus padding
func (lm *LayoutManager) GetSafeWidth(padding int) int {
	safeWidth := lm.width - (padding * 2)
	if safeWidth < 0 {
		return 0
	}
	return safeWidth
}

// GetSafeHeight returns height minus padding
func (lm *LayoutManager) GetSafeHeight(padding int) int {
	safeHeight := lm.height - (padding * 2)
	if safeHeight < 0 {
		return 0
	}
	return safeHeight
}

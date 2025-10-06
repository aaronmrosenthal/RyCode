package layout

// DeviceClass represents different device categories based on terminal dimensions
// Optimized for education: supports students coding on iPhones and tablets
type DeviceClass int

const (
	PhoneTiny      DeviceClass = iota // 0-39: iPhone SE with large font
	PhoneCompact                      // 40-54: iPhone SE, Mini (portrait)
	PhoneStandard                     // 55-69: iPhone 12-14 (portrait)
	TabletSmall                       // 70-84: iPad Mini (portrait)
	TabletMedium                      // 85-99: iPad (portrait)
	TabletLarge                       // 100-119: iPad Pro / iPad landscape
	LaptopSmall                       // 120-139: Chromebook / small laptops
	LaptopStandard                    // 140-159: Standard laptops
	DesktopLarge                      // 160+: Large monitors
)

// String returns the string representation of DeviceClass
func (dc DeviceClass) String() string {
	return [...]string{
		"PhoneTiny",
		"PhoneCompact",
		"PhoneStandard",
		"TabletSmall",
		"TabletMedium",
		"TabletLarge",
		"LaptopSmall",
		"LaptopStandard",
		"DesktopLarge",
	}[dc]
}

// MinWidth returns the minimum width for this device class
func (dc DeviceClass) MinWidth() int {
	return [...]int{0, 40, 55, 70, 85, 100, 120, 140, 160}[dc]
}

// MaxWidth returns the maximum width for this device class (0 = unlimited)
func (dc DeviceClass) MaxWidth() int {
	return [...]int{39, 54, 69, 84, 99, 119, 139, 159, 0}[dc]
}

// IsMobile returns true if this is a mobile device (phone)
func (dc DeviceClass) IsMobile() bool {
	return dc == PhoneTiny || dc == PhoneCompact || dc == PhoneStandard
}

// IsTablet returns true if this is a tablet device
func (dc DeviceClass) IsTablet() bool {
	return dc == TabletSmall || dc == TabletMedium || dc == TabletLarge
}

// IsLaptop returns true if this is a laptop device
func (dc DeviceClass) IsLaptop() bool {
	return dc == LaptopSmall || dc == LaptopStandard
}

// IsDesktop returns true if this is a desktop/laptop device
func (dc DeviceClass) IsDesktop() bool {
	return dc.IsLaptop() || dc == DesktopLarge
}

// IsPortrait returns true if this is likely portrait orientation
// Based on width - narrow screens are usually portrait
func (dc DeviceClass) IsPortrait() bool {
	return dc.IsMobile() || dc == TabletSmall || dc == TabletMedium
}

// IsLandscape returns true if this is likely landscape orientation
func (dc DeviceClass) IsLandscape() bool {
	return dc == TabletLarge || dc.IsDesktop()
}

// SupportsMultiPane returns true if this device class supports multi-pane layouts
func (dc DeviceClass) SupportsMultiPane() bool {
	return dc == TabletLarge || dc.IsDesktop()
}

// SupportsSplitPane returns true if this device class supports split-pane layouts
func (dc DeviceClass) SupportsSplitPane() bool {
	return dc.IsTablet() || dc.IsDesktop()
}

// SupportsFileTree returns true if file tree should be shown by default
func (dc DeviceClass) SupportsFileTree() bool {
	return dc >= TabletSmall // Tablets and above
}

// GetRecommendedPanes returns the recommended number of panes for this device class
func (dc DeviceClass) GetRecommendedPanes() int {
	switch dc {
	case PhoneTiny, PhoneCompact, PhoneStandard:
		return 1 // Chat only
	case TabletSmall, TabletMedium:
		return 2 // File tree + chat
	case TabletLarge, LaptopSmall:
		return 2 // File tree + chat (comfortable)
	case LaptopStandard:
		return 3 // File tree + chat + info
	case DesktopLarge:
		return 3 // File tree + chat + preview/info
	default:
		return 1
	}
}

// GetFileTreeWidth returns the recommended file tree width in columns
func (dc DeviceClass) GetFileTreeWidth() int {
	switch dc {
	case PhoneTiny, PhoneCompact, PhoneStandard:
		return 0 // No file tree by default on phones
	case TabletSmall:
		return 20 // Compact file tree for iPad Mini
	case TabletMedium:
		return 25 // Standard file tree for iPad
	case TabletLarge:
		return 28 // Comfortable file tree
	case LaptopSmall:
		return 30 // Full file tree
	case LaptopStandard, DesktopLarge:
		return 35 // Spacious file tree
	default:
		return 30
	}
}

// GetMinimumChatWidth returns the minimum chat width needed for comfortable coding
func (dc DeviceClass) GetMinimumChatWidth() int {
	switch dc {
	case PhoneTiny:
		return 35 // Absolute minimum
	case PhoneCompact:
		return 40 // Tight but usable
	case PhoneStandard:
		return 50 // Comfortable phone coding
	case TabletSmall, TabletMedium:
		return 60 // Tablet coding
	default:
		return 80 // Desktop coding
	}
}

// ShouldShowFileTreeOverlay returns true if file tree should be an overlay instead of split
func (dc DeviceClass) ShouldShowFileTreeOverlay() bool {
	return dc.IsMobile() // Phones use overlay/drawer
}

// GetRecommendedFontScale returns font size multiplier for readability
func (dc DeviceClass) GetRecommendedFontScale() float64 {
	switch dc {
	case PhoneTiny:
		return 0.85 // Slightly smaller to fit more
	case PhoneCompact:
		return 0.9 // Compact but readable
	case PhoneStandard:
		return 1.0 // Standard
	case TabletSmall, TabletMedium:
		return 1.0 // Standard for tablets
	default:
		return 1.0 // Standard for laptops/desktops
	}
}

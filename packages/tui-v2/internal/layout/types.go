package layout

// DeviceClass represents different device categories based on terminal dimensions
type DeviceClass int

const (
	PhonePortrait DeviceClass = iota
	PhoneLandscape
	TabletPortrait
	TabletLandscape
	DesktopSmall
	DesktopLarge
)

// String returns the string representation of DeviceClass
func (dc DeviceClass) String() string {
	return [...]string{
		"PhonePortrait",
		"PhoneLandscape",
		"TabletPortrait",
		"TabletLandscape",
		"DesktopSmall",
		"DesktopLarge",
	}[dc]
}

// MinWidth returns the minimum width for this device class
func (dc DeviceClass) MinWidth() int {
	return [...]int{0, 60, 80, 100, 120, 160}[dc]
}

// MaxWidth returns the maximum width for this device class (0 = unlimited)
func (dc DeviceClass) MaxWidth() int {
	return [...]int{59, 79, 99, 119, 159, 0}[dc]
}

// IsMobile returns true if this is a mobile device (phone)
func (dc DeviceClass) IsMobile() bool {
	return dc == PhonePortrait || dc == PhoneLandscape
}

// IsTablet returns true if this is a tablet device
func (dc DeviceClass) IsTablet() bool {
	return dc == TabletPortrait || dc == TabletLandscape
}

// IsDesktop returns true if this is a desktop device
func (dc DeviceClass) IsDesktop() bool {
	return dc == DesktopSmall || dc == DesktopLarge
}

// IsPortrait returns true if this is a portrait orientation
func (dc DeviceClass) IsPortrait() bool {
	return dc == PhonePortrait || dc == TabletPortrait
}

// IsLandscape returns true if this is a landscape orientation
func (dc DeviceClass) IsLandscape() bool {
	return dc == PhoneLandscape || dc == TabletLandscape || dc.IsDesktop()
}

// SupportsMultiPane returns true if this device class supports multi-pane layouts
func (dc DeviceClass) SupportsMultiPane() bool {
	return dc.IsDesktop()
}

// SupportsSplitPane returns true if this device class supports split-pane layouts
func (dc DeviceClass) SupportsSplitPane() bool {
	return dc.IsTablet() || dc.IsDesktop()
}

// GetRecommendedPanes returns the recommended number of panes for this device class
func (dc DeviceClass) GetRecommendedPanes() int {
	switch dc {
	case PhonePortrait, PhoneLandscape:
		return 1
	case TabletPortrait, TabletLandscape:
		return 2
	case DesktopSmall:
		return 3
	case DesktopLarge:
		return 4
	default:
		return 1
	}
}

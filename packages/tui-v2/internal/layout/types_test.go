package layout

import (
	"testing"
)

func TestDeviceClass_String(t *testing.T) {
	tests := []struct {
		class DeviceClass
		want  string
	}{
		{PhoneTiny, "PhoneTiny"},
		{PhoneCompact, "PhoneCompact"},
		{PhoneStandard, "PhoneStandard"},
		{TabletSmall, "TabletSmall"},
		{TabletMedium, "TabletMedium"},
		{TabletLarge, "TabletLarge"},
		{LaptopSmall, "LaptopSmall"},
		{LaptopStandard, "LaptopStandard"},
		{DesktopLarge, "DesktopLarge"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.class.String(); got != tt.want {
				t.Errorf("DeviceClass.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeviceClass_MinWidth(t *testing.T) {
	tests := []struct {
		class DeviceClass
		want  int
	}{
		{PhoneTiny, 0},
		{PhoneCompact, 40},
		{PhoneStandard, 55},
		{TabletSmall, 70},
		{TabletMedium, 85},
		{TabletLarge, 100},
		{LaptopSmall, 120},
		{LaptopStandard, 140},
		{DesktopLarge, 160},
	}

	for _, tt := range tests {
		t.Run(tt.class.String(), func(t *testing.T) {
			if got := tt.class.MinWidth(); got != tt.want {
				t.Errorf("DeviceClass.MinWidth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeviceClass_MaxWidth(t *testing.T) {
	tests := []struct {
		class DeviceClass
		want  int
	}{
		{PhoneTiny, 39},
		{PhoneCompact, 54},
		{PhoneStandard, 69},
		{TabletSmall, 84},
		{TabletMedium, 99},
		{TabletLarge, 119},
		{LaptopSmall, 139},
		{LaptopStandard, 159},
		{DesktopLarge, 0}, // Unlimited
	}

	for _, tt := range tests {
		t.Run(tt.class.String(), func(t *testing.T) {
			if got := tt.class.MaxWidth(); got != tt.want {
				t.Errorf("DeviceClass.MaxWidth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeviceClass_IsMobile(t *testing.T) {
	tests := []struct {
		class DeviceClass
		want  bool
	}{
		{PhoneTiny, true},
		{PhoneCompact, true},
		{PhoneStandard, true},
		{TabletSmall, false},
		{TabletMedium, false},
		{TabletLarge, false},
		{LaptopSmall, false},
		{LaptopStandard, false},
		{DesktopLarge, false},
	}

	for _, tt := range tests {
		t.Run(tt.class.String(), func(t *testing.T) {
			if got := tt.class.IsMobile(); got != tt.want {
				t.Errorf("DeviceClass.IsMobile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeviceClass_IsTablet(t *testing.T) {
	tests := []struct {
		class DeviceClass
		want  bool
	}{
		{PhoneTiny, false},
		{PhoneCompact, false},
		{PhoneStandard, false},
		{TabletSmall, true},
		{TabletMedium, true},
		{TabletLarge, true},
		{LaptopSmall, false},
		{LaptopStandard, false},
		{DesktopLarge, false},
	}

	for _, tt := range tests {
		t.Run(tt.class.String(), func(t *testing.T) {
			if got := tt.class.IsTablet(); got != tt.want {
				t.Errorf("DeviceClass.IsTablet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeviceClass_IsLaptop(t *testing.T) {
	tests := []struct {
		class DeviceClass
		want  bool
	}{
		{PhoneTiny, false},
		{PhoneCompact, false},
		{PhoneStandard, false},
		{TabletSmall, false},
		{TabletMedium, false},
		{TabletLarge, false},
		{LaptopSmall, true},
		{LaptopStandard, true},
		{DesktopLarge, false},
	}

	for _, tt := range tests {
		t.Run(tt.class.String(), func(t *testing.T) {
			if got := tt.class.IsLaptop(); got != tt.want {
				t.Errorf("DeviceClass.IsLaptop() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeviceClass_IsDesktop(t *testing.T) {
	tests := []struct {
		class DeviceClass
		want  bool
	}{
		{PhoneTiny, false},
		{PhoneCompact, false},
		{PhoneStandard, false},
		{TabletSmall, false},
		{TabletMedium, false},
		{TabletLarge, false},
		{LaptopSmall, true},
		{LaptopStandard, true},
		{DesktopLarge, true},
	}

	for _, tt := range tests {
		t.Run(tt.class.String(), func(t *testing.T) {
			if got := tt.class.IsDesktop(); got != tt.want {
				t.Errorf("DeviceClass.IsDesktop() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeviceClass_GetRecommendedPanes(t *testing.T) {
	tests := []struct {
		class DeviceClass
		want  int
	}{
		{PhoneTiny, 1},
		{PhoneCompact, 1},
		{PhoneStandard, 1},
		{TabletSmall, 2},
		{TabletMedium, 2},
		{TabletLarge, 2},
		{LaptopSmall, 2},
		{LaptopStandard, 3},
		{DesktopLarge, 3},
	}

	for _, tt := range tests {
		t.Run(tt.class.String(), func(t *testing.T) {
			if got := tt.class.GetRecommendedPanes(); got != tt.want {
				t.Errorf("DeviceClass.GetRecommendedPanes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeviceClass_SupportsMultiPane(t *testing.T) {
	tests := []struct {
		class DeviceClass
		want  bool
	}{
		{PhoneTiny, false},
		{PhoneCompact, false},
		{PhoneStandard, false},
		{TabletSmall, false},
		{TabletMedium, false},
		{TabletLarge, true},
		{LaptopSmall, true},
		{LaptopStandard, true},
		{DesktopLarge, true},
	}

	for _, tt := range tests {
		t.Run(tt.class.String(), func(t *testing.T) {
			if got := tt.class.SupportsMultiPane(); got != tt.want {
				t.Errorf("DeviceClass.SupportsMultiPane() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeviceClass_SupportsFileTree(t *testing.T) {
	tests := []struct {
		class DeviceClass
		want  bool
	}{
		{PhoneTiny, false},
		{PhoneCompact, false},
		{PhoneStandard, false},
		{TabletSmall, true},
		{TabletMedium, true},
		{TabletLarge, true},
		{LaptopSmall, true},
		{LaptopStandard, true},
		{DesktopLarge, true},
	}

	for _, tt := range tests {
		t.Run(tt.class.String(), func(t *testing.T) {
			if got := tt.class.SupportsFileTree(); got != tt.want {
				t.Errorf("DeviceClass.SupportsFileTree() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeviceClass_GetFileTreeWidth(t *testing.T) {
	tests := []struct {
		class DeviceClass
		want  int
	}{
		{PhoneTiny, 0},
		{PhoneCompact, 0},
		{PhoneStandard, 0},
		{TabletSmall, 20},
		{TabletMedium, 25},
		{TabletLarge, 28},
		{LaptopSmall, 30},
		{LaptopStandard, 35},
		{DesktopLarge, 35},
	}

	for _, tt := range tests {
		t.Run(tt.class.String(), func(t *testing.T) {
			if got := tt.class.GetFileTreeWidth(); got != tt.want {
				t.Errorf("DeviceClass.GetFileTreeWidth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeviceClass_ShouldShowFileTreeOverlay(t *testing.T) {
	tests := []struct {
		class DeviceClass
		want  bool
	}{
		{PhoneTiny, true},
		{PhoneCompact, true},
		{PhoneStandard, true},
		{TabletSmall, false},
		{TabletMedium, false},
		{TabletLarge, false},
		{LaptopSmall, false},
		{LaptopStandard, false},
		{DesktopLarge, false},
	}

	for _, tt := range tests {
		t.Run(tt.class.String(), func(t *testing.T) {
			if got := tt.class.ShouldShowFileTreeOverlay(); got != tt.want {
				t.Errorf("DeviceClass.ShouldShowFileTreeOverlay() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test breakpoint transitions for education devices
func TestDeviceClass_EducationDevices(t *testing.T) {
	tests := []struct {
		name     string
		width    int
		wantClass DeviceClass
		wantDesc string
	}{
		{"iPhone SE with large font", 38, PhoneTiny, "Emergency tiny mode"},
		{"iPhone SE portrait", 45, PhoneCompact, "iPhone SE, Mini"},
		{"iPhone 12 portrait", 58, PhoneStandard, "iPhone 12-14"},
		{"iPad Mini portrait", 75, TabletSmall, "iPad Mini"},
		{"iPad portrait", 90, TabletMedium, "iPad 10.2/10.9"},
		{"iPad landscape", 110, TabletLarge, "iPad landscape"},
		{"Chromebook", 130, LaptopSmall, "Chromebook"},
		{"Standard laptop", 150, LaptopStandard, "MacBook Air"},
		{"Large monitor", 180, DesktopLarge, "iMac/external"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lm := NewLayoutManager(tt.width, 24)
			if got := lm.GetDeviceClass(); got != tt.wantClass {
				t.Errorf("%s (width=%d): got %v, want %v (%s)",
					tt.name, tt.width, got, tt.wantClass, tt.wantDesc)
			}
		})
	}
}

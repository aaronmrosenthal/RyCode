package layout

import (
	"testing"
)

func TestDeviceClass_String(t *testing.T) {
	tests := []struct {
		class DeviceClass
		want  string
	}{
		{PhonePortrait, "PhonePortrait"},
		{PhoneLandscape, "PhoneLandscape"},
		{TabletPortrait, "TabletPortrait"},
		{TabletLandscape, "TabletLandscape"},
		{DesktopSmall, "DesktopSmall"},
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
		{PhonePortrait, 0},
		{PhoneLandscape, 60},
		{TabletPortrait, 80},
		{TabletLandscape, 100},
		{DesktopSmall, 120},
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

func TestDeviceClass_IsMobile(t *testing.T) {
	tests := []struct {
		class DeviceClass
		want  bool
	}{
		{PhonePortrait, true},
		{PhoneLandscape, true},
		{TabletPortrait, false},
		{TabletLandscape, false},
		{DesktopSmall, false},
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
		{PhonePortrait, false},
		{PhoneLandscape, false},
		{TabletPortrait, true},
		{TabletLandscape, true},
		{DesktopSmall, false},
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

func TestDeviceClass_IsDesktop(t *testing.T) {
	tests := []struct {
		class DeviceClass
		want  bool
	}{
		{PhonePortrait, false},
		{PhoneLandscape, false},
		{TabletPortrait, false},
		{TabletLandscape, false},
		{DesktopSmall, true},
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
		{PhonePortrait, 1},
		{PhoneLandscape, 1},
		{TabletPortrait, 2},
		{TabletLandscape, 2},
		{DesktopSmall, 3},
		{DesktopLarge, 4},
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
		{PhonePortrait, false},
		{PhoneLandscape, false},
		{TabletPortrait, false},
		{TabletLandscape, false},
		{DesktopSmall, true},
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

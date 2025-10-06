package layout

import (
	"testing"
	"time"
)

func TestNewLayoutManager(t *testing.T) {
	lm := NewLayoutManager(100, 50)

	if lm.GetWidth() != 100 {
		t.Errorf("Expected width 100, got %d", lm.GetWidth())
	}

	if lm.GetHeight() != 50 {
		t.Errorf("Expected height 50, got %d", lm.GetHeight())
	}

	if lm.GetDeviceClass() != TabletLarge {
		t.Errorf("Expected TabletLarge, got %s", lm.GetDeviceClass())
	}
}

func TestLayoutManager_DetectDevice(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
		want   DeviceClass
	}{
		{"Phone tiny", 35, 80, PhoneTiny},
		{"Phone compact", 45, 80, PhoneCompact},
		{"Phone standard", 60, 80, PhoneStandard},
		{"Tablet small", 75, 100, TabletSmall},
		{"Tablet medium", 90, 100, TabletMedium},
		{"Tablet large", 110, 60, TabletLarge},
		{"Laptop small", 130, 60, LaptopSmall},
		{"Laptop standard", 150, 60, LaptopStandard},
		{"Desktop large", 180, 60, DesktopLarge},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lm := NewLayoutManager(tt.width, tt.height)
			if got := lm.GetDeviceClass(); got != tt.want {
				t.Errorf("DetectDevice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLayoutManager_Update(t *testing.T) {
	lm := NewLayoutManager(35, 80)

	if lm.GetDeviceClass() != PhoneTiny {
		t.Errorf("Initial device class should be PhoneTiny")
	}

	// Update to desktop size
	lm.Update(180, 60)

	if lm.GetDeviceClass() != DesktopLarge {
		t.Errorf("After update, device class should be DesktopLarge, got %s", lm.GetDeviceClass())
	}

	w, h := lm.GetDimensions()
	if w != 180 || h != 60 {
		t.Errorf("Dimensions should be 180x60, got %dx%d", w, h)
	}
}

func TestLayoutManager_OnChange(t *testing.T) {
	lm := NewLayoutManager(35, 80)

	callbackCalled := false
	var receivedClass DeviceClass

	lm.OnChange(func(dc DeviceClass) {
		callbackCalled = true
		receivedClass = dc
	})

	// Update to different device class
	lm.Update(180, 60)

	if !callbackCalled {
		t.Error("OnChange callback should have been called")
	}

	if receivedClass != DesktopLarge {
		t.Errorf("Callback should receive DesktopLarge, got %s", receivedClass)
	}
}

func TestLayoutManager_OnChangeNotCalledSameClass(t *testing.T) {
	lm := NewLayoutManager(45, 80)

	callbackCalled := false

	lm.OnChange(func(dc DeviceClass) {
		callbackCalled = true
	})

	// Update to same device class (still phone compact)
	lm.Update(50, 90)

	if callbackCalled {
		t.Error("OnChange callback should NOT be called for same device class")
	}
}

func TestLayoutManager_ShouldUseStackLayout(t *testing.T) {
	tests := []struct {
		width int
		want  bool
	}{
		{35, true},   // Phone tiny
		{45, true},   // Phone compact
		{60, true},   // Phone standard
		{75, false},  // Tablet small
		{90, false},  // Tablet medium
		{110, false}, // Tablet large
		{130, false}, // Laptop small
		{150, false}, // Laptop standard
		{180, false}, // Desktop large
	}

	for _, tt := range tests {
		lm := NewLayoutManager(tt.width, 60)
		if got := lm.ShouldUseStackLayout(); got != tt.want {
			t.Errorf("Width %d: ShouldUseStackLayout() = %v, want %v (class=%s)",
				tt.width, got, tt.want, lm.GetDeviceClass())
		}
	}
}

func TestLayoutManager_ShouldUseSplitLayout(t *testing.T) {
	tests := []struct {
		width int
		want  bool
	}{
		{35, false},  // Phone tiny
		{45, false},  // Phone compact
		{60, false},  // Phone standard
		{75, true},   // Tablet small
		{90, true},   // Tablet medium
		{110, true},  // Tablet large
		{130, false}, // Laptop small (uses desktop layout)
		{150, false}, // Laptop standard (uses desktop layout)
		{180, false}, // Desktop large
	}

	for _, tt := range tests {
		lm := NewLayoutManager(tt.width, 60)
		if got := lm.ShouldUseSplitLayout(); got != tt.want {
			t.Errorf("Width %d: ShouldUseSplitLayout() = %v, want %v (class=%s)",
				tt.width, got, tt.want, lm.GetDeviceClass())
		}
	}
}

func TestLayoutManager_ShouldUseMultiPaneLayout(t *testing.T) {
	tests := []struct {
		width int
		want  bool
	}{
		{35, false},  // Phone tiny
		{45, false},  // Phone compact
		{60, false},  // Phone standard
		{75, false},  // Tablet small
		{90, false},  // Tablet medium
		{110, false}, // Tablet large
		{130, true},  // Laptop small
		{150, true},  // Laptop standard
		{180, true},  // Desktop large
	}

	for _, tt := range tests {
		lm := NewLayoutManager(tt.width, 60)
		if got := lm.ShouldUseMultiPaneLayout(); got != tt.want {
			t.Errorf("Width %d: ShouldUseMultiPaneLayout() = %v, want %v (class=%s)",
				tt.width, got, tt.want, lm.GetDeviceClass())
		}
	}
}

func TestLayoutManager_GetRecommendedSplitRatio(t *testing.T) {
	tests := []struct {
		class DeviceClass
		want  float64
	}{
		{PhoneTiny, 1.0},
		{PhoneCompact, 1.0},
		{PhoneStandard, 1.0},
		{TabletSmall, 0.75},
		{TabletMedium, 0.70},
		{TabletLarge, 0.72},
		{LaptopSmall, 0.75},
		{LaptopStandard, 0.75},
		{DesktopLarge, 0.75},
	}

	for _, tt := range tests {
		lm := NewLayoutManager(tt.class.MinWidth()+10, 60)
		if got := lm.GetRecommendedSplitRatio(); got != tt.want {
			t.Errorf("%s: GetRecommendedSplitRatio() = %v, want %v", tt.class, got, tt.want)
		}
	}
}

func TestLayoutManager_ShouldShowFileTree(t *testing.T) {
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
		lm := NewLayoutManager(tt.class.MinWidth()+10, 60)
		if got := lm.ShouldShowFileTree(); got != tt.want {
			t.Errorf("%s: ShouldShowFileTree() = %v, want %v", tt.class, got, tt.want)
		}
	}
}

func TestLayoutManager_ShouldUseFileTreeOverlay(t *testing.T) {
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
		lm := NewLayoutManager(tt.class.MinWidth()+10, 60)
		if got := lm.ShouldUseFileTreeOverlay(); got != tt.want {
			t.Errorf("%s: ShouldUseFileTreeOverlay() = %v, want %v", tt.class, got, tt.want)
		}
	}
}

func TestLayoutManager_GetFileTreeWidthForDevice(t *testing.T) {
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
		lm := NewLayoutManager(tt.class.MinWidth()+10, 60)
		if got := lm.GetFileTreeWidthForDevice(); got != tt.want {
			t.Errorf("%s: GetFileTreeWidthForDevice() = %v, want %v", tt.class, got, tt.want)
		}
	}
}

func TestLayoutManager_CanFitWidth(t *testing.T) {
	lm := NewLayoutManager(100, 50)

	if !lm.CanFitWidth(80) {
		t.Error("Should fit 80 width")
	}

	if lm.CanFitWidth(120) {
		t.Error("Should NOT fit 120 width")
	}
}

func TestLayoutManager_CanFitHeight(t *testing.T) {
	lm := NewLayoutManager(100, 50)

	if !lm.CanFitHeight(40) {
		t.Error("Should fit 40 height")
	}

	if lm.CanFitHeight(60) {
		t.Error("Should NOT fit 60 height")
	}
}

func TestLayoutManager_GetSafeWidth(t *testing.T) {
	lm := NewLayoutManager(100, 50)

	if got := lm.GetSafeWidth(5); got != 90 {
		t.Errorf("GetSafeWidth(5) = %d, want 90", got)
	}

	if got := lm.GetSafeWidth(0); got != 100 {
		t.Errorf("GetSafeWidth(0) = %d, want 100", got)
	}
}

func TestLayoutManager_GetSafeHeight(t *testing.T) {
	lm := NewLayoutManager(100, 50)

	if got := lm.GetSafeHeight(5); got != 40 {
		t.Errorf("GetSafeHeight(5) = %d, want 40", got)
	}

	if got := lm.GetSafeHeight(0); got != 50 {
		t.Errorf("GetSafeHeight(0) = %d, want 50", got)
	}
}

func TestLayoutManager_GetLastUpdate(t *testing.T) {
	before := time.Now()
	lm := NewLayoutManager(100, 50)
	after := time.Now()

	lastUpdate := lm.GetLastUpdate()

	if lastUpdate.Before(before) || lastUpdate.After(after) {
		t.Error("LastUpdate should be between before and after timestamps")
	}

	time.Sleep(10 * time.Millisecond)
	beforeUpdate := time.Now()
	lm.Update(120, 60)
	afterUpdate := time.Now()

	lastUpdate = lm.GetLastUpdate()

	if lastUpdate.Before(beforeUpdate) || lastUpdate.After(afterUpdate) {
		t.Error("LastUpdate should be updated after Update() call")
	}
}

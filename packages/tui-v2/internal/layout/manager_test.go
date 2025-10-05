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

	if lm.GetDeviceClass() != TabletLandscape {
		t.Errorf("Expected TabletLandscape, got %s", lm.GetDeviceClass())
	}
}

func TestLayoutManager_DetectDevice(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
		want   DeviceClass
	}{
		{"Phone portrait", 40, 80, PhonePortrait},
		{"Phone landscape", 70, 40, PhoneLandscape},
		{"Tablet portrait", 90, 100, TabletPortrait},
		{"Tablet landscape", 110, 60, TabletLandscape},
		{"Desktop small", 140, 60, DesktopSmall},
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
	lm := NewLayoutManager(40, 80)

	if lm.GetDeviceClass() != PhonePortrait {
		t.Errorf("Initial device class should be PhonePortrait")
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
	lm := NewLayoutManager(40, 80)

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
	lm := NewLayoutManager(40, 80)

	callbackCalled := false

	lm.OnChange(func(dc DeviceClass) {
		callbackCalled = true
	})

	// Update to same device class (still phone portrait)
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
		{40, true},   // Phone portrait
		{70, true},   // Phone landscape
		{90, false},  // Tablet portrait
		{110, false}, // Tablet landscape
		{140, false}, // Desktop small
		{180, false}, // Desktop large
	}

	for _, tt := range tests {
		lm := NewLayoutManager(tt.width, 60)
		if got := lm.ShouldUseStackLayout(); got != tt.want {
			t.Errorf("Width %d: ShouldUseStackLayout() = %v, want %v", tt.width, got, tt.want)
		}
	}
}

func TestLayoutManager_ShouldUseSplitLayout(t *testing.T) {
	tests := []struct {
		width int
		want  bool
	}{
		{40, false},  // Phone portrait
		{70, false},  // Phone landscape
		{90, true},   // Tablet portrait
		{110, true},  // Tablet landscape
		{140, false}, // Desktop small
		{180, false}, // Desktop large
	}

	for _, tt := range tests {
		lm := NewLayoutManager(tt.width, 60)
		if got := lm.ShouldUseSplitLayout(); got != tt.want {
			t.Errorf("Width %d: ShouldUseSplitLayout() = %v, want %v", tt.width, got, tt.want)
		}
	}
}

func TestLayoutManager_ShouldUseMultiPaneLayout(t *testing.T) {
	tests := []struct {
		width int
		want  bool
	}{
		{40, false},  // Phone portrait
		{70, false},  // Phone landscape
		{90, false},  // Tablet portrait
		{110, false}, // Tablet landscape
		{140, true},  // Desktop small
		{180, true},  // Desktop large
	}

	for _, tt := range tests {
		lm := NewLayoutManager(tt.width, 60)
		if got := lm.ShouldUseMultiPaneLayout(); got != tt.want {
			t.Errorf("Width %d: ShouldUseMultiPaneLayout() = %v, want %v", tt.width, got, tt.want)
		}
	}
}

func TestLayoutManager_GetRecommendedSplitRatio(t *testing.T) {
	tests := []struct {
		class DeviceClass
		want  float64
	}{
		{PhonePortrait, 0.5},
		{PhoneLandscape, 0.5},
		{TabletPortrait, 0.5},
		{TabletLandscape, 0.6},
		{DesktopSmall, 0.7},
		{DesktopLarge, 0.7},
	}

	for _, tt := range tests {
		lm := NewLayoutManager(tt.class.MinWidth()+10, 60)
		if got := lm.GetRecommendedSplitRatio(); got != tt.want {
			t.Errorf("%s: GetRecommendedSplitRatio() = %v, want %v", tt.class, got, tt.want)
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

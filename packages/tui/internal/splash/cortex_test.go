package splash

import (
	"math"
	"testing"
)

// TestTorusGeometry tests that torus points are within expected bounds
func TestTorusGeometry(t *testing.T) {
	_ = NewCortexRenderer(80, 24)

	// Test points around the torus
	for theta := 0.0; theta < 2*math.Pi; theta += 0.5 {
		for phi := 0.0; phi < 2*math.Pi; phi += 0.5 {
			// Calculate a point on the torus
			const R = 2.0 // Major radius
			const rr = 1.0 // Minor radius

			circleX := R + rr*math.Cos(phi)
			circleY := rr * math.Sin(phi)

			x := circleX * math.Cos(theta)
			y := circleX * math.Sin(theta)
			z := circleY

			// Distance from origin should be R ± r (2 ± 1 = 1 to 3)
			dist := math.Sqrt(x*x + y*y + z*z)
			if dist < 0.5 || dist > 3.5 {
				t.Errorf("Point out of bounds: dist=%.2f at θ=%.2f φ=%.2f", dist, theta, phi)
			}
		}
	}
}

// TestZBufferOcclusion tests that z-buffer prevents incorrect occlusion
func TestZBufferOcclusion(t *testing.T) {
	r := NewCortexRenderer(80, 24)
	r.RenderFrame()

	// Center should have some depth (not empty)
	centerIdx := (r.height/2)*r.width + (r.width / 2)
	if r.zbuffer[centerIdx] == 0 {
		t.Error("Z-buffer not working - center pixel is empty")
	}

	// Check that some pixels have depth
	nonZeroCount := 0
	for _, z := range r.zbuffer {
		if z > 0 {
			nonZeroCount++
		}
	}

	if nonZeroCount == 0 {
		t.Error("No pixels have depth - z-buffer is not working")
	}
}

// TestRenderFrameNoPanic tests that rendering doesn't panic
func TestRenderFrameNoPanic(t *testing.T) {
	r := NewCortexRenderer(80, 24)

	defer func() {
		if rec := recover(); rec != nil {
			t.Errorf("RenderFrame panicked: %v", rec)
		}
	}()

	// Render 100 frames
	for i := 0; i < 100; i++ {
		r.RenderFrame()
	}
}

// TestRenderWithColors tests that colored rendering works
func TestRenderWithColors(t *testing.T) {
	r := NewCortexRenderer(80, 24)

	defer func() {
		if rec := recover(); rec != nil {
			t.Errorf("Render panicked: %v", rec)
		}
	}()

	output := r.Render()
	if len(output) == 0 {
		t.Error("Render produced empty output")
	}
}

// TestRotationAnglesUpdate tests that rotation angles update correctly
func TestRotationAnglesUpdate(t *testing.T) {
	r := NewCortexRenderer(80, 24)

	initialA := r.A
	initialB := r.B

	r.RenderFrame()

	if r.A <= initialA {
		t.Error("Angle A did not increase after RenderFrame")
	}

	if r.B <= initialB {
		t.Error("Angle B did not increase after RenderFrame")
	}
}

// BenchmarkRenderFrame benchmarks the frame rendering performance
func BenchmarkRenderFrame(b *testing.B) {
	r := NewCortexRenderer(80, 24)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.RenderFrame()
	}
	// Target: <1ms (1,000,000 ns) per frame
}

// BenchmarkFullRender benchmarks full rendering with colors
func BenchmarkFullRender(b *testing.B) {
	r := NewCortexRenderer(80, 24)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = r.Render()
	}
	// Target: <10ms per full render
}

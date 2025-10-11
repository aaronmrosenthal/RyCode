package performance

import (
	"testing"
	"time"
)

func TestPerformanceMonitor_FrameMetrics(t *testing.T) {
	pm := NewPerformanceMonitor()

	// Simulate frames
	for i := 0; i < 10; i++ {
		pm.StartFrame()
		time.Sleep(time.Millisecond * 10) // Simulate 10ms render
		pm.EndFrame()
	}

	metrics := pm.GetFrameMetrics()

	if metrics.TotalFrames != 10 {
		t.Errorf("Expected 10 frames, got %d", metrics.TotalFrames)
	}

	if metrics.FPS == 0 {
		t.Error("FPS should not be 0")
	}

	if metrics.AverageFrameTime == 0 {
		t.Error("Average frame time should not be 0")
	}
}

func TestPerformanceMonitor_DroppedFrames(t *testing.T) {
	pm := NewPerformanceMonitor()

	// Simulate slow frame (should be dropped)
	pm.StartFrame()
	time.Sleep(time.Millisecond * 20) // Exceeds 16ms budget
	pm.EndFrame()

	metrics := pm.GetFrameMetrics()

	if metrics.DroppedFrames != 1 {
		t.Errorf("Expected 1 dropped frame, got %d", metrics.DroppedFrames)
	}
}

func TestPerformanceMonitor_ComponentMetrics(t *testing.T) {
	pm := NewPerformanceMonitor()

	// Simulate component render
	start := pm.StartComponentRender("TestComponent")
	time.Sleep(time.Millisecond * 5)
	pm.EndComponentRender("TestComponent", start)

	metrics := pm.GetComponentMetrics()

	if len(metrics) != 1 {
		t.Errorf("Expected 1 component, got %d", len(metrics))
	}

	compMetric := metrics["TestComponent"]
	if compMetric == nil {
		t.Fatal("TestComponent metrics not found")
	}

	if compMetric.RenderCount != 1 {
		t.Errorf("Expected 1 render, got %d", compMetric.RenderCount)
	}
}

func TestPerformanceMonitor_Health(t *testing.T) {
	pm := NewPerformanceMonitor()

	// Simulate good frames to establish baseline
	for i := 0; i < 60; i++ {
		pm.StartFrame()
		time.Sleep(time.Millisecond * 5) // Fast frame (well under 16ms budget)
		pm.EndFrame()
	}

	health := pm.GetHealth()
	if health < 90 {
		t.Errorf("Expected health > 90 for good frames, got %.2f", health)
	}

	// Simulate bad frames
	pm.Reset()
	for i := 0; i < 60; i++ {
		pm.StartFrame()
		time.Sleep(time.Millisecond * 20) // Slow frame (exceeds 16ms budget)
		pm.EndFrame()
	}

	health = pm.GetHealth()
	if health > 50 {
		t.Errorf("Expected health < 50 for bad frames, got %.2f", health)
	}
}

func TestPerformanceMonitor_EnableDisable(t *testing.T) {
	pm := NewPerformanceMonitor()

	if !pm.IsEnabled() {
		t.Error("Monitor should be enabled by default")
	}

	pm.Disable()
	if pm.IsEnabled() {
		t.Error("Monitor should be disabled")
	}

	pm.Enable()
	if !pm.IsEnabled() {
		t.Error("Monitor should be enabled")
	}
}

func TestPerformanceMonitor_Reset(t *testing.T) {
	pm := NewPerformanceMonitor()

	// Generate some metrics
	pm.StartFrame()
	pm.EndFrame()

	start := pm.StartComponentRender("Test")
	pm.EndComponentRender("Test", start)

	pm.addWarning("Test warning")

	// Reset
	pm.Reset()

	metrics := pm.GetFrameMetrics()
	if metrics.TotalFrames != 0 {
		t.Error("Frames should be reset")
	}

	components := pm.GetComponentMetrics()
	if len(components) != 0 {
		t.Error("Components should be reset")
	}

	warnings := pm.GetWarnings()
	if len(warnings) != 0 {
		t.Error("Warnings should be reset")
	}
}

func TestPerformanceMonitor_DroppedFrameRate(t *testing.T) {
	pm := NewPerformanceMonitor()

	// Simulate 10 frames, 3 dropped
	for i := 0; i < 10; i++ {
		pm.StartFrame()
		if i < 3 {
			time.Sleep(time.Millisecond * 20) // Slow
		} else {
			time.Sleep(time.Millisecond * 5) // Fast
		}
		pm.EndFrame()
	}

	dropRate := pm.GetDroppedFrameRate()
	if dropRate < 25 || dropRate > 35 {
		t.Errorf("Expected drop rate around 30%%, got %.2f%%", dropRate)
	}
}

func TestPerformanceMonitor_MemorySnapshot(t *testing.T) {
	pm := NewPerformanceMonitor()

	pm.RecordMemorySnapshot()

	metrics := pm.GetMemoryMetrics()

	if metrics.Alloc == 0 {
		t.Error("Memory allocation should not be 0")
	}

	if metrics.Sys == 0 {
		t.Error("System memory should not be 0")
	}
}

func TestPerformanceMonitor_Warnings(t *testing.T) {
	pm := NewPerformanceMonitor()

	pm.addWarning("Test warning 1")
	pm.addWarning("Test warning 2")

	warnings := pm.GetWarnings()
	if len(warnings) != 2 {
		t.Errorf("Expected 2 warnings, got %d", len(warnings))
	}

	pm.ClearWarnings()
	warnings = pm.GetWarnings()
	if len(warnings) != 0 {
		t.Errorf("Expected 0 warnings after clear, got %d", len(warnings))
	}
}

func TestPerformanceMonitor_GetSummary(t *testing.T) {
	pm := NewPerformanceMonitor()

	// Generate some metrics
	pm.StartFrame()
	time.Sleep(time.Millisecond * 10)
	pm.EndFrame()

	pm.RecordMemorySnapshot()

	summary := pm.GetSummary()
	if summary == "" {
		t.Error("Summary should not be empty")
	}

	// Should contain key metrics
	if len(summary) < 50 {
		t.Error("Summary should contain detailed metrics")
	}
}

// Benchmark tests
func BenchmarkPerformanceMonitor_FrameCycle(b *testing.B) {
	pm := NewPerformanceMonitor()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pm.StartFrame()
		pm.EndFrame()
	}
}

func BenchmarkPerformanceMonitor_ComponentRender(b *testing.B) {
	pm := NewPerformanceMonitor()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		start := pm.StartComponentRender("BenchComponent")
		pm.EndComponentRender("BenchComponent", start)
	}
}

func BenchmarkPerformanceMonitor_MemorySnapshot(b *testing.B) {
	pm := NewPerformanceMonitor()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pm.RecordMemorySnapshot()
	}
}

func BenchmarkPerformanceMonitor_GetMetrics(b *testing.B) {
	pm := NewPerformanceMonitor()

	// Populate with data
	for i := 0; i < 100; i++ {
		pm.StartFrame()
		pm.EndFrame()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pm.GetFrameMetrics()
		pm.GetMemoryMetrics()
		pm.GetComponentMetrics()
	}
}

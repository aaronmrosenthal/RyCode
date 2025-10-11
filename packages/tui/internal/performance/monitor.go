package performance

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// FrameMetrics tracks rendering performance
type FrameMetrics struct {
	FrameTime       time.Duration // Time to render single frame
	FPS             float64       // Frames per second
	DroppedFrames   int           // Frames that exceeded budget
	TotalFrames     int           // Total frames rendered
	AverageFrameTime time.Duration // Average across all frames
}

// MemoryMetrics tracks memory usage
type MemoryMetrics struct {
	Alloc         uint64 // Bytes allocated and in use
	TotalAlloc    uint64 // Bytes allocated (cumulative)
	Sys           uint64 // Bytes obtained from system
	NumGC         uint32 // Number of GC runs
	LastGCPause   time.Duration // Last GC pause duration
	HeapObjects   uint64 // Number of allocated heap objects
}

// RenderMetrics tracks specific render operations
type RenderMetrics struct {
	ComponentName   string
	RenderTime      time.Duration
	RenderCount     int
	AverageRenderTime time.Duration
}

// PerformanceMonitor tracks app performance metrics
type PerformanceMonitor struct {
	mu sync.RWMutex

	// Frame tracking
	frameStartTime    time.Time
	frameTimes        []time.Duration
	maxFrameHistory   int
	targetFrameTime   time.Duration // 16.67ms for 60fps
	droppedFrames     int
	totalFrames       int

	// Component tracking
	componentMetrics  map[string]*RenderMetrics

	// Memory tracking
	memorySnapshots   []MemoryMetrics
	maxMemoryHistory  int

	// Enabled flag
	enabled           bool

	// Performance warnings
	warnings          []string
	maxWarnings       int
}

// NewPerformanceMonitor creates a new performance monitor
func NewPerformanceMonitor() *PerformanceMonitor {
	return &PerformanceMonitor{
		frameTimes:       make([]time.Duration, 0, 100),
		maxFrameHistory:  100,
		targetFrameTime:  time.Millisecond * 16, // ~60fps
		componentMetrics: make(map[string]*RenderMetrics),
		memorySnapshots:  make([]MemoryMetrics, 0, 50),
		maxMemoryHistory: 50,
		warnings:         make([]string, 0),
		maxWarnings:      10,
		enabled:          true, // Enable in dev, disable in production
	}
}

// StartFrame marks the beginning of a frame render
func (pm *PerformanceMonitor) StartFrame() {
	if !pm.enabled {
		return
	}

	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.frameStartTime = time.Now()
}

// EndFrame marks the end of a frame render and records metrics
func (pm *PerformanceMonitor) EndFrame() {
	if !pm.enabled {
		return
	}

	pm.mu.Lock()
	defer pm.mu.Unlock()

	frameTime := time.Since(pm.frameStartTime)

	// Record frame time
	pm.frameTimes = append(pm.frameTimes, frameTime)
	if len(pm.frameTimes) > pm.maxFrameHistory {
		pm.frameTimes = pm.frameTimes[1:]
	}

	// Track dropped frames
	pm.totalFrames++
	if frameTime > pm.targetFrameTime {
		pm.droppedFrames++

		// Warn if frame time is significantly over budget
		if frameTime > pm.targetFrameTime*2 {
			pm.addWarning(fmt.Sprintf("Slow frame: %.2fms (target: %.2fms)",
				float64(frameTime.Microseconds())/1000.0,
				float64(pm.targetFrameTime.Microseconds())/1000.0))
		}
	}
}

// StartComponentRender marks the start of a component render
func (pm *PerformanceMonitor) StartComponentRender(componentName string) time.Time {
	if !pm.enabled {
		return time.Time{}
	}
	return time.Now()
}

// EndComponentRender records component render time
func (pm *PerformanceMonitor) EndComponentRender(componentName string, startTime time.Time) {
	if !pm.enabled {
		return
	}

	pm.mu.Lock()
	defer pm.mu.Unlock()

	renderTime := time.Since(startTime)

	metrics, exists := pm.componentMetrics[componentName]
	if !exists {
		metrics = &RenderMetrics{
			ComponentName: componentName,
		}
		pm.componentMetrics[componentName] = metrics
	}

	metrics.RenderTime = renderTime
	metrics.RenderCount++

	// Calculate rolling average
	if metrics.AverageRenderTime == 0 {
		metrics.AverageRenderTime = renderTime
	} else {
		// Exponential moving average
		alpha := 0.1
		metrics.AverageRenderTime = time.Duration(
			float64(metrics.AverageRenderTime)*(1-alpha) +
			float64(renderTime)*alpha,
		)
	}

	// Warn about slow components
	if renderTime > time.Millisecond*5 {
		pm.addWarning(fmt.Sprintf("Slow component '%s': %.2fms",
			componentName, float64(renderTime.Microseconds())/1000.0))
	}
}

// RecordMemorySnapshot captures current memory usage
func (pm *PerformanceMonitor) RecordMemorySnapshot() {
	if !pm.enabled {
		return
	}

	pm.mu.Lock()
	defer pm.mu.Unlock()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	snapshot := MemoryMetrics{
		Alloc:       m.Alloc,
		TotalAlloc:  m.TotalAlloc,
		Sys:         m.Sys,
		NumGC:       m.NumGC,
		HeapObjects: m.HeapObjects,
	}

	// Calculate last GC pause
	if m.NumGC > 0 {
		snapshot.LastGCPause = time.Duration(m.PauseNs[(m.NumGC+255)%256])
	}

	pm.memorySnapshots = append(pm.memorySnapshots, snapshot)
	if len(pm.memorySnapshots) > pm.maxMemoryHistory {
		pm.memorySnapshots = pm.memorySnapshots[1:]
	}

	// Warn about high memory usage (>100MB)
	if snapshot.Alloc > 100*1024*1024 {
		pm.addWarning(fmt.Sprintf("High memory usage: %.2f MB",
			float64(snapshot.Alloc)/(1024*1024)))
	}
}

// GetFrameMetrics returns current frame performance metrics
func (pm *PerformanceMonitor) GetFrameMetrics() FrameMetrics {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	if len(pm.frameTimes) == 0 {
		return FrameMetrics{}
	}

	// Calculate average frame time
	var total time.Duration
	for _, ft := range pm.frameTimes {
		total += ft
	}
	avgFrameTime := total / time.Duration(len(pm.frameTimes))

	// Calculate FPS
	fps := 0.0
	if avgFrameTime > 0 {
		fps = float64(time.Second) / float64(avgFrameTime)
	}

	// Get last frame time
	lastFrameTime := pm.frameTimes[len(pm.frameTimes)-1]

	return FrameMetrics{
		FrameTime:       lastFrameTime,
		FPS:             fps,
		DroppedFrames:   pm.droppedFrames,
		TotalFrames:     pm.totalFrames,
		AverageFrameTime: avgFrameTime,
	}
}

// GetMemoryMetrics returns current memory metrics
func (pm *PerformanceMonitor) GetMemoryMetrics() MemoryMetrics {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	if len(pm.memorySnapshots) == 0 {
		return MemoryMetrics{}
	}

	return pm.memorySnapshots[len(pm.memorySnapshots)-1]
}

// GetComponentMetrics returns metrics for all components
func (pm *PerformanceMonitor) GetComponentMetrics() map[string]*RenderMetrics {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	// Return copy to avoid race conditions
	metrics := make(map[string]*RenderMetrics)
	for name, m := range pm.componentMetrics {
		metricsCopy := *m
		metrics[name] = &metricsCopy
	}

	return metrics
}

// GetWarnings returns performance warnings
func (pm *PerformanceMonitor) GetWarnings() []string {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	warnings := make([]string, len(pm.warnings))
	copy(warnings, pm.warnings)
	return warnings
}

// ClearWarnings clears all performance warnings
func (pm *PerformanceMonitor) ClearWarnings() {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.warnings = make([]string, 0)
}

// addWarning adds a performance warning
func (pm *PerformanceMonitor) addWarning(warning string) {
	pm.warnings = append(pm.warnings, warning)
	if len(pm.warnings) > pm.maxWarnings {
		pm.warnings = pm.warnings[1:]
	}
}

// Reset clears all metrics
func (pm *PerformanceMonitor) Reset() {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.frameTimes = make([]time.Duration, 0, pm.maxFrameHistory)
	pm.droppedFrames = 0
	pm.totalFrames = 0
	pm.componentMetrics = make(map[string]*RenderMetrics)
	pm.memorySnapshots = make([]MemoryMetrics, 0, pm.maxMemoryHistory)
	pm.warnings = make([]string, 0)
}

// Enable enables performance monitoring
func (pm *PerformanceMonitor) Enable() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.enabled = true
}

// Disable disables performance monitoring
func (pm *PerformanceMonitor) Disable() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.enabled = false
}

// IsEnabled returns whether monitoring is enabled
func (pm *PerformanceMonitor) IsEnabled() bool {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.enabled
}

// GetDroppedFrameRate returns percentage of dropped frames
func (pm *PerformanceMonitor) GetDroppedFrameRate() float64 {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	if pm.totalFrames == 0 {
		return 0
	}

	return float64(pm.droppedFrames) / float64(pm.totalFrames) * 100.0
}

// GetHealth returns an overall health score (0-100)
func (pm *PerformanceMonitor) GetHealth() float64 {
	metrics := pm.GetFrameMetrics()
	memory := pm.GetMemoryMetrics()

	health := 100.0

	// Deduct for low FPS
	if metrics.FPS < 60 {
		health -= (60 - metrics.FPS)
	}

	// Deduct for dropped frames
	dropRate := pm.GetDroppedFrameRate()
	health -= dropRate

	// Deduct for high memory usage (>50MB)
	memoryMB := float64(memory.Alloc) / (1024 * 1024)
	if memoryMB > 50 {
		health -= (memoryMB - 50) / 5
	}

	// Deduct for warnings
	warnings := len(pm.GetWarnings())
	health -= float64(warnings) * 2

	if health < 0 {
		health = 0
	}
	if health > 100 {
		health = 100
	}

	return health
}

// GetSummary returns a human-readable performance summary
func (pm *PerformanceMonitor) GetSummary() string {
	metrics := pm.GetFrameMetrics()
	memory := pm.GetMemoryMetrics()
	health := pm.GetHealth()

	return fmt.Sprintf(
		"FPS: %.1f | Frame: %.2fms | Dropped: %d/%d (%.1f%%) | Mem: %.2fMB | Health: %.0f%%",
		metrics.FPS,
		float64(metrics.AverageFrameTime.Microseconds())/1000.0,
		metrics.DroppedFrames,
		metrics.TotalFrames,
		pm.GetDroppedFrameRate(),
		float64(memory.Alloc)/(1024*1024),
		health,
	)
}

// Global performance monitor instance
var globalMonitor = NewPerformanceMonitor()

// Global accessor functions
func StartFrame() {
	globalMonitor.StartFrame()
}

func EndFrame() {
	globalMonitor.EndFrame()
}

func StartComponentRender(name string) time.Time {
	return globalMonitor.StartComponentRender(name)
}

func EndComponentRender(name string, start time.Time) {
	globalMonitor.EndComponentRender(name, start)
}

func RecordMemorySnapshot() {
	globalMonitor.RecordMemorySnapshot()
}

func GetMonitor() *PerformanceMonitor {
	return globalMonitor
}

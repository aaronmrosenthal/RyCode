package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/aaronmrosenthal/rycode/internal/theme"
)

// PerformanceTest represents a single performance measurement
type PerformanceTest struct {
	Name        string
	Iterations  int
	MaxDuration time.Duration
	Fn          func()
}

// Run executes the performance test
func (pt *PerformanceTest) Run() (time.Duration, bool) {
	// Warm up
	for i := 0; i < 10; i++ {
		pt.Fn()
	}

	// Force garbage collection before test
	runtime.GC()

	// Measure
	start := time.Now()
	for i := 0; i < pt.Iterations; i++ {
		pt.Fn()
	}
	elapsed := time.Since(start)

	avgDuration := elapsed / time.Duration(pt.Iterations)
	passed := avgDuration <= pt.MaxDuration

	return avgDuration, passed
}

func main() {
	fmt.Println("=== Theme Performance Benchmark ===")
	fmt.Println("Target: <10ms per theme switch")
	fmt.Println("Goal: Imperceptible to users (<16.67ms = 60fps)")
	fmt.Println()

	providers := []string{"claude", "gemini", "codex", "qwen"}
	allPassed := true

	// Test 1: Theme Switching Speed
	fmt.Println("[Test 1] Theme Switching Performance")
	fmt.Println("Switching between all 4 providers in sequence...")
	fmt.Println()

	test := PerformanceTest{
		Name:        "Sequential Theme Switch",
		Iterations:  1000,
		MaxDuration: 10 * time.Millisecond,
		Fn: func() {
			for _, providerID := range providers {
				theme.SwitchToProvider(providerID)
			}
		},
	}

	avgDuration, passed := test.Run()
	status := "✓ PASS"
	if !passed {
		status = "✗ FAIL"
		allPassed = false
	}

	perSwitchDuration := avgDuration / time.Duration(len(providers))
	fmt.Printf("  %s Average per switch: %v (target: <10ms)\n", status, perSwitchDuration)
	fmt.Printf("       Total 4 switches: %v\n", avgDuration)
	fmt.Printf("       Iterations: %d\n", test.Iterations)
	fmt.Println()

	// Test 2: Theme Retrieval Speed
	fmt.Println("[Test 2] Theme Retrieval Performance")
	fmt.Println("Getting current theme repeatedly...")
	fmt.Println()

	theme.SwitchToProvider("claude")
	test2 := PerformanceTest{
		Name:        "CurrentTheme() Call",
		Iterations:  100000,
		MaxDuration: 100 * time.Nanosecond, // Should be extremely fast
		Fn: func() {
			_ = theme.CurrentTheme()
		},
	}

	avgDuration2, passed2 := test2.Run()
	status2 := "✓ PASS"
	if !passed2 {
		status2 = "✗ FAIL"
		allPassed = false
	}

	fmt.Printf("  %s Average: %v (target: <100ns)\n", status2, avgDuration2)
	fmt.Printf("       Iterations: %d\n", test2.Iterations)
	fmt.Println()

	// Test 3: Color Access Speed
	fmt.Println("[Test 3] Color Access Performance")
	fmt.Println("Accessing theme colors (Primary, Text, Border)...")
	fmt.Println()

	test3 := PerformanceTest{
		Name:        "Color Access",
		Iterations:  100000,
		MaxDuration: 200 * time.Nanosecond,
		Fn: func() {
			t := theme.CurrentTheme()
			_ = t.Primary()
			_ = t.Text()
			_ = t.Border()
		},
	}

	avgDuration3, passed3 := test3.Run()
	status3 := "✓ PASS"
	if !passed3 {
		status3 = "✗ FAIL"
		allPassed = false
	}

	fmt.Printf("  %s Average: %v (target: <200ns)\n", status3, avgDuration3)
	fmt.Printf("       Iterations: %d\n", test3.Iterations)
	fmt.Println()

	// Test 4: Memory Allocation During Switch
	fmt.Println("[Test 4] Memory Allocation")
	fmt.Println("Measuring memory usage during theme switches...")
	fmt.Println()

	runtime.GC()
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	allocsBefore := memStats.Alloc

	for i := 0; i < 1000; i++ {
		for _, providerID := range providers {
			theme.SwitchToProvider(providerID)
		}
	}

	runtime.ReadMemStats(&memStats)
	allocsAfter := memStats.Alloc

	allocsPerSwitch := (allocsAfter - allocsBefore) / (1000 * uint64(len(providers)))

	// Target: <1KB per switch
	memPassed := allocsPerSwitch < 1024
	status4 := "✓ PASS"
	if !memPassed {
		status4 = "✗ FAIL"
		allPassed = false
	}

	fmt.Printf("  %s Allocs per switch: %d bytes (target: <1KB)\n", status4, allocsPerSwitch)
	fmt.Printf("       Total switches: %d\n", 1000*len(providers))
	fmt.Println()

	// Test 5: Rapid Switching Stress Test
	fmt.Println("[Test 5] Rapid Switching Stress Test")
	fmt.Println("Simulating rapid Tab presses...")
	fmt.Println()

	test5 := PerformanceTest{
		Name:        "Rapid Switch",
		Iterations:  100,
		MaxDuration: 5 * time.Millisecond, // Faster than Test 1 (single switch)
		Fn: func() {
			// Simulate user rapidly pressing Tab
			theme.SwitchToProvider("claude")
			theme.SwitchToProvider("gemini")
			theme.SwitchToProvider("codex")
			theme.SwitchToProvider("qwen")
		},
	}

	avgDuration5, passed5 := test5.Run()
	status5 := "✓ PASS"
	if !passed5 {
		status5 = "✗ FAIL"
		allPassed = false
	}

	perSwitchDuration5 := avgDuration5 / 4
	fmt.Printf("  %s Average per switch: %v (target: <5ms)\n", status5, perSwitchDuration5)
	fmt.Printf("       Total 4 switches: %v\n", avgDuration5)
	fmt.Printf("       Iterations: %d\n", test5.Iterations)
	fmt.Println()

	// Summary
	fmt.Println("=== Performance Summary ===")
	fmt.Println()

	if allPassed {
		fmt.Println("✅ All performance tests passed!")
		fmt.Println()
		fmt.Println("Performance Analysis:")
		fmt.Printf("  • Theme switching is %s (target: <10ms)\n", formatSpeed(perSwitchDuration))
		fmt.Printf("  • Theme retrieval is %s (target: <100ns)\n", formatSpeed(avgDuration2))
		fmt.Printf("  • Color access is %s (target: <200ns)\n", formatSpeed(avgDuration3))
		fmt.Printf("  • Memory usage is %s per switch\n", formatBytes(allocsPerSwitch))
		fmt.Println()
		fmt.Println("User Experience:")
		if perSwitchDuration < 16670*time.Microsecond {
			fmt.Println("  ⚡ Switching is faster than 60fps (imperceptible)")
		} else if perSwitchDuration < 33333*time.Microsecond {
			fmt.Println("  🚀 Switching is faster than 30fps (very smooth)")
		} else {
			fmt.Println("  ✓ Switching is acceptably fast")
		}
		fmt.Println()
		os.Exit(0)
	} else {
		fmt.Println("❌ Some performance tests failed!")
		fmt.Println()
		fmt.Println("Optimization needed:")
		if !passed {
			fmt.Println("  • Theme switching is too slow")
		}
		if !passed2 {
			fmt.Println("  • Theme retrieval needs optimization")
		}
		if !passed3 {
			fmt.Println("  • Color access needs caching")
		}
		if !memPassed {
			fmt.Println("  • Memory allocations too high")
		}
		fmt.Println()
		os.Exit(1)
	}
}

func formatSpeed(d time.Duration) string {
	if d < time.Microsecond {
		return fmt.Sprintf("%.0fns", float64(d.Nanoseconds()))
	} else if d < time.Millisecond {
		return fmt.Sprintf("%.2fμs", float64(d.Nanoseconds())/1000.0)
	} else {
		return fmt.Sprintf("%.2fms", float64(d.Nanoseconds())/1000000.0)
	}
}

func formatBytes(b uint64) string {
	if b < 1024 {
		return fmt.Sprintf("%d bytes", b)
	} else if b < 1024*1024 {
		return fmt.Sprintf("%.2f KB", float64(b)/1024.0)
	} else {
		return fmt.Sprintf("%.2f MB", float64(b)/(1024.0*1024.0))
	}
}

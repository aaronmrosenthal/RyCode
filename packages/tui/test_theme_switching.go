package main

import (
	"fmt"
	"os"

	"github.com/aaronmrosenthal/rycode/internal/theme"
)

func main() {
	fmt.Println("=== Theme Switching Test ===\n")

	// Test 1: Get initial theme (should be Claude by default)
	fmt.Println("[Test 1] Initial theme:")
	currentTheme := theme.CurrentTheme()
	if currentTheme == nil {
		fmt.Println("  ✗ ERROR: No theme loaded")
		os.Exit(1)
	}
	fmt.Printf("  ✓ Theme loaded: %s\n", currentTheme.Name())
	fmt.Printf("    Primary color: %v\n", currentTheme.Primary())
	fmt.Println()

	// Test 2: Switch to Gemini
	fmt.Println("[Test 2] Switch to Gemini:")
	changed := theme.SwitchToProvider("gemini")
	if !changed {
		fmt.Println("  ✗ ERROR: Theme did not change")
		os.Exit(1)
	}
	currentTheme = theme.CurrentTheme()
	fmt.Printf("  ✓ Switched to: %s\n", currentTheme.Name())
	fmt.Printf("    Primary color: %v (should be blue #4285F4)\n", currentTheme.Primary())
	fmt.Println()

	// Test 3: Switch to Codex
	fmt.Println("[Test 3] Switch to Codex:")
	changed = theme.SwitchToProvider("codex")
	if !changed {
		fmt.Println("  ✗ ERROR: Theme did not change")
		os.Exit(1)
	}
	currentTheme = theme.CurrentTheme()
	fmt.Printf("  ✓ Switched to: %s\n", currentTheme.Name())
	fmt.Printf("    Primary color: %v (should be teal #10A37F)\n", currentTheme.Primary())
	fmt.Println()

	// Test 4: Switch to Qwen
	fmt.Println("[Test 4] Switch to Qwen:")
	changed = theme.SwitchToProvider("qwen")
	if !changed {
		fmt.Println("  ✗ ERROR: Theme did not change")
		os.Exit(1)
	}
	currentTheme = theme.CurrentTheme()
	fmt.Printf("  ✓ Switched to: %s\n", currentTheme.Name())
	fmt.Printf("    Primary color: %v (should be orange #FF6A00)\n", currentTheme.Primary())
	fmt.Println()

	// Test 5: Switch back to Claude
	fmt.Println("[Test 5] Switch back to Claude:")
	changed = theme.SwitchToProvider("claude")
	if !changed {
		fmt.Println("  ✗ ERROR: Theme did not change")
		os.Exit(1)
	}
	currentTheme = theme.CurrentTheme()
	fmt.Printf("  ✓ Switched to: %s\n", currentTheme.Name())
	fmt.Printf("    Primary color: %v (should be copper #D4754C)\n", currentTheme.Primary())
	fmt.Println()

	// Test 6: Try invalid provider (should not panic)
	fmt.Println("[Test 6] Try invalid provider:")
	changed = theme.SwitchToProvider("invalid")
	if changed {
		fmt.Println("  ✗ ERROR: Theme changed for invalid provider")
		os.Exit(1)
	}
	fmt.Println("  ✓ Invalid provider correctly ignored")
	currentTheme = theme.CurrentTheme()
	fmt.Printf("    Still on: %s\n", currentTheme.Name())
	fmt.Println()

	// Test 7: Check that theme persists across calls
	fmt.Println("[Test 7] Verify theme persistence:")
	theme.SwitchToProvider("gemini")
	theme1 := theme.CurrentTheme()
	theme2 := theme.CurrentTheme()
	if theme1.Name() != theme2.Name() {
		fmt.Println("  ✗ ERROR: Theme not consistent across calls")
		os.Exit(1)
	}
	fmt.Printf("  ✓ Theme persists: %s\n", theme1.Name())
	fmt.Println()

	fmt.Println("✅ All tests passed! Theme switching works correctly.")
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Launch RyCode TUI")
	fmt.Println("  2. Press Tab to cycle through providers")
	fmt.Println("  3. Watch borders, badges, and UI elements change color")
	fmt.Println("  4. Verify Claude=copper, Gemini=blue/pink, Codex=teal, Qwen=orange")
}

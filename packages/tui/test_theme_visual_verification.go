package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aaronmrosenthal/rycode/internal/theme"
	"github.com/charmbracelet/lipgloss/v2/compat"
)

// ColorTest represents a single color verification test
type ColorTest struct {
	Provider string
	Element  string
	Expected string // Hex color like "#D4754C"
}

// All color tests across all providers
var colorTests = []ColorTest{
	// Claude Theme (from NewClaudeTheme)
	{"claude", "Primary", "#D4754C"},
	{"claude", "Text", "#E8D5C4"},
	{"claude", "TextMuted", "#9C8373"},
	{"claude", "Background", "#1A1816"},
	{"claude", "BackgroundPanel", "#2C2622"},
	{"claude", "Border", "#D4754C"},
	{"claude", "Accent", "#F08C5C"},
	{"claude", "Success", "#6FA86F"},
	{"claude", "Error", "#D47C7C"},
	{"claude", "Warning", "#E8A968"},
	{"claude", "Info", "#D4754C"},
	{"claude", "MarkdownHeading", "#F08C5C"},
	{"claude", "MarkdownLink", "#D4754C"},
	{"claude", "MarkdownCode", "#E8A968"},

	// Gemini Theme (from NewGeminiTheme)
	{"gemini", "Primary", "#4285F4"},
	{"gemini", "Text", "#E8EAED"},
	{"gemini", "TextMuted", "#9AA0A6"},
	{"gemini", "Background", "#0D0D0D"},
	{"gemini", "BackgroundPanel", "#1A1A1A"},
	{"gemini", "Border", "#4285F4"},
	{"gemini", "Accent", "#EA4335"},
	{"gemini", "Success", "#34A853"},
	{"gemini", "Error", "#EA4335"},
	{"gemini", "Warning", "#FBBC04"},
	{"gemini", "Info", "#4285F4"},
	{"gemini", "MarkdownHeading", "#4285F4"},
	{"gemini", "MarkdownLink", "#4285F4"},
	{"gemini", "MarkdownCode", "#FBBC04"},

	// Codex Theme (from NewCodexTheme)
	{"codex", "Primary", "#10A37F"},
	{"codex", "Text", "#ECECEC"},
	{"codex", "TextMuted", "#8E8E8E"},
	{"codex", "Background", "#0E0E0E"},
	{"codex", "BackgroundPanel", "#1C1C1C"},
	{"codex", "Border", "#10A37F"},
	{"codex", "Accent", "#1FC2AA"},
	{"codex", "Success", "#10A37F"},
	{"codex", "Error", "#EF4444"},
	{"codex", "Warning", "#F59E0B"},
	{"codex", "Info", "#3B82F6"},
	{"codex", "MarkdownHeading", "#1FC2AA"},
	{"codex", "MarkdownLink", "#10A37F"},
	{"codex", "MarkdownCode", "#F59E0B"},

	// Qwen Theme (from NewQwenTheme)
	{"qwen", "Primary", "#FF6A00"},
	{"qwen", "Text", "#F0E8DC"},
	{"qwen", "TextMuted", "#A0947C"},
	{"qwen", "Background", "#161410"},
	{"qwen", "BackgroundPanel", "#221E18"},
	{"qwen", "Border", "#FF6A00"},
	{"qwen", "Accent", "#FF8533"},
	{"qwen", "Success", "#52C41A"},
	{"qwen", "Error", "#FF4D4F"},
	{"qwen", "Warning", "#FAAD14"},
	{"qwen", "Info", "#1890FF"},
	{"qwen", "MarkdownHeading", "#FF6A00"},
	{"qwen", "MarkdownLink", "#1890FF"},
	{"qwen", "MarkdownCode", "#FAAD14"},
}

func main() {
	fmt.Println("=== Theme Visual Verification ===")
	fmt.Println("Verifying all theme colors match specifications...")
	fmt.Println()

	allPassed := true
	providerResults := make(map[string][]string)

	for _, test := range colorTests {
		// Switch to provider theme
		theme.SwitchToProvider(test.Provider)
		th := theme.CurrentTheme()

		// Get the actual color from the theme
		var actualColor compat.AdaptiveColor
		switch test.Element {
		case "Primary":
			actualColor = th.Primary()
		case "Text":
			actualColor = th.Text()
		case "TextMuted":
			actualColor = th.TextMuted()
		case "Background":
			actualColor = th.Background()
		case "BackgroundPanel":
			actualColor = th.BackgroundPanel()
		case "Border":
			actualColor = th.Border()
		case "Accent":
			actualColor = th.Accent()
		case "Success":
			actualColor = th.Success()
		case "Error":
			actualColor = th.Error()
		case "Warning":
			actualColor = th.Warning()
		case "Info":
			actualColor = th.Info()
		case "MarkdownHeading":
			actualColor = th.MarkdownHeading()
		case "MarkdownLink":
			actualColor = th.MarkdownLink()
		case "MarkdownCode":
			actualColor = th.MarkdownCode()
		default:
			fmt.Printf("  ✗ Unknown element: %s\n", test.Element)
			allPassed = false
			continue
		}

		// Convert to hex string
		actualHex := colorToHex(actualColor)

		// Compare
		passed := actualHex == test.Expected
		status := "✓"
		if !passed {
			status = "✗"
			allPassed = false
		}

		result := fmt.Sprintf("  %s %-20s expected %s, got %s", status, test.Element, test.Expected, actualHex)
		providerResults[test.Provider] = append(providerResults[test.Provider], result)

		if !passed {
			fmt.Println(result)
		}
	}

	// Print results grouped by provider
	providers := []string{"claude", "gemini", "codex", "qwen"}
	for _, provider := range providers {
		results := providerResults[provider]
		passed := 0
		failed := 0

		for _, result := range results {
			if strings.Contains(result, "✓") {
				passed++
			} else {
				failed++
			}
		}

		fmt.Printf("[%s Theme]\n", provider)
		if failed > 0 {
			// Print failed tests
			for _, result := range results {
				if strings.Contains(result, "✗") {
					fmt.Println(result)
				}
			}
		}
		fmt.Printf("  Summary: %d passed, %d failed\n\n", passed, failed)
	}

	// Final summary
	if allPassed {
		totalTests := len(colorTests)
		fmt.Println("=== Visual Verification Summary ===")
		fmt.Println()
		fmt.Printf("✅ All %d color tests passed!\n", totalTests)
		fmt.Println()
		fmt.Println("Theme Color Accuracy:")
		fmt.Println("  • All primary colors match specifications")
		fmt.Println("  • All text colors match specifications")
		fmt.Println("  • All UI element colors match specifications")
		fmt.Println("  • All markdown colors match specifications")
		fmt.Println()
		fmt.Println("Benefits:")
		fmt.Println("  • Visual consistency guaranteed")
		fmt.Println("  • Brand colors accurately replicated")
		fmt.Println("  • No color drift over time")
		fmt.Println("  • CI-ready for regression detection")
		fmt.Println()
		os.Exit(0)
	} else {
		fmt.Println("❌ Some color tests failed!")
		fmt.Println()
		fmt.Println("Action Required:")
		fmt.Println("  • Review failed colors above")
		fmt.Println("  • Update theme definitions in provider_themes.go")
		fmt.Println("  • Re-run this test to verify fixes")
		fmt.Println()
		os.Exit(1)
	}
}

// colorToHex converts an AdaptiveColor to a hex string
func colorToHex(ac compat.AdaptiveColor) string {
	// Use dark variant since RyCode is a dark TUI
	c := ac.Dark
	r, g, b, _ := c.RGBA()

	// Convert from 16-bit (0-65535) to 8-bit (0-255)
	r8 := uint8(r >> 8)
	g8 := uint8(g >> 8)
	b8 := uint8(b >> 8)

	return fmt.Sprintf("#%02X%02X%02X", r8, g8, b8)
}

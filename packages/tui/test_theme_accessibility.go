package main

import (
	"fmt"
	"image/color"
	"math"
	"os"

	"github.com/aaronmrosenthal/rycode/internal/theme"
	"github.com/charmbracelet/lipgloss/v2/compat"
)

// ContrastRatio calculates the WCAG contrast ratio between two colors
func ContrastRatio(c1, c2 color.Color) float64 {
	l1 := relativeLuminance(c1)
	l2 := relativeLuminance(c2)

	if l1 > l2 {
		return (l1 + 0.05) / (l2 + 0.05)
	}
	return (l2 + 0.05) / (l1 + 0.05)
}

// relativeLuminance calculates the relative luminance of a color
func relativeLuminance(c color.Color) float64 {
	r, g, b, _ := c.RGBA()

	// Convert to 0-1 range
	rNorm := float64(r) / 65535.0
	gNorm := float64(g) / 65535.0
	bNorm := float64(b) / 65535.0

	// Apply gamma correction
	rLinear := toLinear(rNorm)
	gLinear := toLinear(gNorm)
	bLinear := toLinear(bNorm)

	// Calculate luminance
	return 0.2126*rLinear + 0.7152*gLinear + 0.0722*bLinear
}

// toLinear applies gamma correction
func toLinear(v float64) float64 {
	if v <= 0.03928 {
		return v / 12.92
	}
	return math.Pow((v+0.055)/1.055, 2.4)
}

// adaptiveColorToColor extracts the dark variant from AdaptiveColor
func adaptiveColorToColor(ac compat.AdaptiveColor) color.Color {
	return ac.Dark
}

// ColorTest represents a single contrast test
type ColorTest struct {
	Name       string
	Foreground color.Color
	Background color.Color
	MinRatio   float64 // 4.5 for AA normal text, 3.0 for AA large text, 7.0 for AAA
}

func main() {
	fmt.Println("=== Theme Accessibility Audit ===")
	fmt.Println("WCAG 2.1 Contrast Requirements:")
	fmt.Println("  AA Normal Text: 4.5:1")
	fmt.Println("  AA Large Text:  3.0:1")
	fmt.Println("  AAA Normal Text: 7.0:1")
	fmt.Println()

	providers := []string{"claude", "gemini", "codex", "qwen"}
	allPassed := true

	for _, providerID := range providers {
		fmt.Printf("=== %s Theme ===\n", providerID)

		// Switch to provider theme
		theme.SwitchToProvider(providerID)
		t := theme.CurrentTheme()

		if t == nil {
			fmt.Printf("  ✗ ERROR: Could not load theme\n\n")
			allPassed = false
			continue
		}

		tests := []ColorTest{
			// Primary text on background
			{
				Name:       "Text on Background",
				Foreground: adaptiveColorToColor(t.Text()),
				Background: adaptiveColorToColor(t.Background()),
				MinRatio:   4.5,
			},
			// Muted text on background
			{
				Name:       "Muted Text on Background",
				Foreground: adaptiveColorToColor(t.TextMuted()),
				Background: adaptiveColorToColor(t.Background()),
				MinRatio:   4.5,
			},
			// Primary text on panel
			{
				Name:       "Text on Panel",
				Foreground: adaptiveColorToColor(t.Text()),
				Background: adaptiveColorToColor(t.BackgroundPanel()),
				MinRatio:   4.5,
			},
			// Border on background (large text/UI components)
			{
				Name:       "Border on Background",
				Foreground: adaptiveColorToColor(t.Border()),
				Background: adaptiveColorToColor(t.Background()),
				MinRatio:   3.0, // Large text standard
			},
			// Primary color on background
			{
				Name:       "Primary on Background",
				Foreground: adaptiveColorToColor(t.Primary()),
				Background: adaptiveColorToColor(t.Background()),
				MinRatio:   3.0, // For UI elements
			},
			// Success color on background
			{
				Name:       "Success on Background",
				Foreground: adaptiveColorToColor(t.Success()),
				Background: adaptiveColorToColor(t.Background()),
				MinRatio:   3.0,
			},
			// Error color on background
			{
				Name:       "Error on Background",
				Foreground: adaptiveColorToColor(t.Error()),
				Background: adaptiveColorToColor(t.Background()),
				MinRatio:   3.0,
			},
			// Warning color on background
			{
				Name:       "Warning on Background",
				Foreground: adaptiveColorToColor(t.Warning()),
				Background: adaptiveColorToColor(t.Background()),
				MinRatio:   3.0,
			},
			// Info color on background
			{
				Name:       "Info on Background",
				Foreground: adaptiveColorToColor(t.Info()),
				Background: adaptiveColorToColor(t.Background()),
				MinRatio:   3.0,
			},
			// Markdown heading on background
			{
				Name:       "Markdown Heading on Background",
				Foreground: adaptiveColorToColor(t.MarkdownHeading()),
				Background: adaptiveColorToColor(t.Background()),
				MinRatio:   4.5,
			},
			// Markdown link on background
			{
				Name:       "Markdown Link on Background",
				Foreground: adaptiveColorToColor(t.MarkdownLink()),
				Background: adaptiveColorToColor(t.Background()),
				MinRatio:   4.5,
			},
			// Markdown code on background
			{
				Name:       "Markdown Code on Background",
				Foreground: adaptiveColorToColor(t.MarkdownCode()),
				Background: adaptiveColorToColor(t.Background()),
				MinRatio:   4.5,
			},
		}

		passed := 0
		failed := 0

		for _, test := range tests {
			ratio := ContrastRatio(test.Foreground, test.Background)
			status := "✓"
			statusText := "PASS"

			if ratio < test.MinRatio {
				status = "✗"
				statusText = "FAIL"
				failed++
				allPassed = false
			} else {
				passed++
			}

			// Determine WCAG level
			wcagLevel := ""
			if ratio >= 7.0 {
				wcagLevel = "AAA"
			} else if ratio >= 4.5 {
				wcagLevel = "AA "
			} else if ratio >= 3.0 {
				wcagLevel = "AA Large"
			}

			fmt.Printf("  %s %-35s %.2f:1 (req %.1f:1) [%s] %s\n",
				status, test.Name, ratio, test.MinRatio, wcagLevel, statusText)
		}

		fmt.Printf("\nSummary: %d passed, %d failed\n\n", passed, failed)
	}

	if allPassed {
		fmt.Println("✅ All themes pass WCAG AA accessibility standards!")
		os.Exit(0)
	} else {
		fmt.Println("❌ Some themes have accessibility issues. See failures above.")
		os.Exit(1)
	}
}

package responsive

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss/v2"
)

// RGB represents an RGB color
type RGB struct {
	R uint8
	G uint8
	B uint8
}

// ContrastRatio calculates WCAG contrast ratio between two colors
func ContrastRatio(fg, bg RGB) float64 {
	l1 := relativeLuminance(fg)
	l2 := relativeLuminance(bg)

	lighter := math.Max(l1, l2)
	darker := math.Min(l1, l2)

	return (lighter + 0.05) / (darker + 0.05)
}

// relativeLuminance calculates relative luminance per WCAG spec
func relativeLuminance(c RGB) float64 {
	r := linearize(float64(c.R) / 255.0)
	g := linearize(float64(c.G) / 255.0)
	b := linearize(float64(c.B) / 255.0)

	return 0.2126*r + 0.7152*g + 0.0722*b
}

// linearize converts sRGB to linear RGB
func linearize(val float64) float64 {
	if val <= 0.03928 {
		return val / 12.92
	}
	return math.Pow((val+0.055)/1.055, 2.4)
}

// MeetsWCAG_AA checks if contrast ratio meets WCAG AA standard
func MeetsWCAG_AA(ratio float64, isLargeText bool) bool {
	if isLargeText {
		return ratio >= 3.0
	}
	return ratio >= 4.5
}

// MeetsWCAG_AAA checks if contrast ratio meets WCAG AAA standard
func MeetsWCAG_AAA(ratio float64, isLargeText bool) bool {
	if isLargeText {
		return ratio >= 4.5
	}
	return ratio >= 7.0
}

// ColorFromLipgloss converts lipgloss.Color to RGB
func ColorFromLipgloss(c lipgloss.Color) (RGB, error) {
	s := string(c)

	// Handle different color formats
	if strings.HasPrefix(s, "#") {
		return parseHex(s)
	}

	// ANSI color numbers (0-255)
	if num, err := strconv.Atoi(s); err == nil {
		return ansiToRGB(num), nil
	}

	// Named colors
	if rgb, ok := namedColors[strings.ToLower(s)]; ok {
		return rgb, nil
	}

	return RGB{}, fmt.Errorf("unsupported color format: %s", s)
}

// parseHex parses hex color (#RRGGBB or #RGB)
func parseHex(hex string) (RGB, error) {
	hex = strings.TrimPrefix(hex, "#")

	// Expand shorthand (#RGB -> #RRGGBB)
	if len(hex) == 3 {
		hex = string([]byte{
			hex[0], hex[0],
			hex[1], hex[1],
			hex[2], hex[2],
		})
	}

	if len(hex) != 6 {
		return RGB{}, fmt.Errorf("invalid hex color: %s", hex)
	}

	r, _ := strconv.ParseUint(hex[0:2], 16, 8)
	g, _ := strconv.ParseUint(hex[2:4], 16, 8)
	b, _ := strconv.ParseUint(hex[4:6], 16, 8)

	return RGB{uint8(r), uint8(g), uint8(b)}, nil
}

// ansiToRGB converts ANSI color number to RGB
func ansiToRGB(num int) RGB {
	if num < 0 || num > 255 {
		return RGB{128, 128, 128} // Gray fallback
	}

	// Standard 16 colors (0-15)
	if num < 16 {
		return ansi16Colors[num]
	}

	// 216 color cube (16-231)
	if num < 232 {
		num -= 16
		r := (num / 36) * 51
		g := ((num / 6) % 6) * 51
		b := (num % 6) * 51
		return RGB{uint8(r), uint8(g), uint8(b)}
	}

	// Grayscale (232-255)
	gray := uint8((num-232)*10 + 8)
	return RGB{gray, gray, gray}
}

// ANSI 16 color palette
var ansi16Colors = [16]RGB{
	{0, 0, 0},       // 0: Black
	{128, 0, 0},     // 1: Red
	{0, 128, 0},     // 2: Green
	{128, 128, 0},   // 3: Yellow
	{0, 0, 128},     // 4: Blue
	{128, 0, 128},   // 5: Magenta
	{0, 128, 128},   // 6: Cyan
	{192, 192, 192}, // 7: White
	{128, 128, 128}, // 8: Bright Black (Gray)
	{255, 0, 0},     // 9: Bright Red
	{0, 255, 0},     // 10: Bright Green
	{255, 255, 0},   // 11: Bright Yellow
	{0, 0, 255},     // 12: Bright Blue
	{255, 0, 255},   // 13: Bright Magenta
	{0, 255, 255},   // 14: Bright Cyan
	{255, 255, 255}, // 15: Bright White
}

// Named CSS colors
var namedColors = map[string]RGB{
	"black":   {0, 0, 0},
	"white":   {255, 255, 255},
	"red":     {255, 0, 0},
	"green":   {0, 255, 0},
	"blue":    {0, 0, 255},
	"yellow":  {255, 255, 0},
	"cyan":    {0, 255, 255},
	"magenta": {255, 0, 255},
	"gray":    {128, 128, 128},
	"grey":    {128, 128, 128},
	"orange":  {255, 165, 0},
	"purple":  {128, 0, 128},
	"pink":    {255, 192, 203},
	"brown":   {165, 42, 42},
}

// ContrastChecker validates contrast ratios
type ContrastChecker struct {
	issues []ContrastIssue
}

// ContrastIssue represents a contrast problem
type ContrastIssue struct {
	ForegroundColor string
	BackgroundColor string
	Ratio           float64
	MeetsAA         bool
	MeetsAAA        bool
	IsLargeText     bool
	Component       string
}

// NewContrastChecker creates a contrast checker
func NewContrastChecker() *ContrastChecker {
	return &ContrastChecker{
		issues: []ContrastIssue{},
	}
}

// Check checks contrast between two colors
func (cc *ContrastChecker) Check(component string, fg, bg lipgloss.Color, isLargeText bool) error {
	fgRGB, err := ColorFromLipgloss(fg)
	if err != nil {
		return err
	}

	bgRGB, err := ColorFromLipgloss(bg)
	if err != nil {
		return err
	}

	ratio := ContrastRatio(fgRGB, bgRGB)
	meetsAA := MeetsWCAG_AA(ratio, isLargeText)
	meetsAAA := MeetsWCAG_AAA(ratio, isLargeText)

	issue := ContrastIssue{
		ForegroundColor: string(fg),
		BackgroundColor: string(bg),
		Ratio:           ratio,
		MeetsAA:         meetsAA,
		MeetsAAA:        meetsAAA,
		IsLargeText:     isLargeText,
		Component:       component,
	}

	if !meetsAA {
		cc.issues = append(cc.issues, issue)
	}

	return nil
}

// GetIssues returns all found contrast issues
func (cc *ContrastChecker) GetIssues() []ContrastIssue {
	return cc.issues
}

// HasIssues returns whether any issues were found
func (cc *ContrastChecker) HasIssues() bool {
	return len(cc.issues) > 0
}

// Report generates a contrast report
func (cc *ContrastChecker) Report() string {
	if !cc.HasIssues() {
		return "✅ All contrast ratios meet WCAG AA standards!"
	}

	lines := []string{
		fmt.Sprintf("⚠️  Found %d contrast issues:", len(cc.issues)),
		"",
	}

	for i, issue := range cc.issues {
		level := "NORMAL"
		if issue.IsLargeText {
			level = "LARGE"
		}

		standard := "AA (4.5:1)"
		if issue.IsLargeText {
			standard = "AA (3.0:1)"
		}

		lines = append(lines, fmt.Sprintf(
			"%d. %s [%s text]\n"+
				"   FG: %s, BG: %s\n"+
				"   Ratio: %.2f:1 (needs %s)\n"+
				"   Meets AA: %v, AAA: %v",
			i+1,
			issue.Component,
			level,
			issue.ForegroundColor,
			issue.BackgroundColor,
			issue.Ratio,
			standard,
			issue.MeetsAA,
			issue.MeetsAAA,
		))
		lines = append(lines, "")
	}

	return strings.Join(lines, "\n")
}

// SuggestFix suggests a fix for low contrast
func SuggestFix(fg, bg RGB, targetRatio float64) RGB {
	// Simple approach: lighten or darken foreground
	currentRatio := ContrastRatio(fg, bg)

	if currentRatio >= targetRatio {
		return fg // Already meets target
	}

	// Determine if we should lighten or darken
	bgLuminance := relativeLuminance(bg)

	var suggested RGB
	if bgLuminance > 0.5 {
		// Light background, darken foreground
		suggested = darken(fg, 0.2)
	} else {
		// Dark background, lighten foreground
		suggested = lighten(fg, 0.2)
	}

	// Check if it helps
	newRatio := ContrastRatio(suggested, bg)
	if newRatio > currentRatio {
		return suggested
	}

	return fg
}

// lighten lightens a color by factor (0-1)
func lighten(c RGB, factor float64) RGB {
	return RGB{
		lightenComponent(c.R, factor),
		lightenComponent(c.G, factor),
		lightenComponent(c.B, factor),
	}
}

// darken darkens a color by factor (0-1)
func darken(c RGB, factor float64) RGB {
	return RGB{
		darkenComponent(c.R, factor),
		darkenComponent(c.G, factor),
		darkenComponent(c.B, factor),
	}
}

func lightenComponent(val uint8, factor float64) uint8 {
	f := float64(val)
	f = f + (255-f)*factor
	if f > 255 {
		f = 255
	}
	return uint8(f)
}

func darkenComponent(val uint8, factor float64) uint8 {
	f := float64(val)
	f = f * (1 - factor)
	if f < 0 {
		f = 0
	}
	return uint8(f)
}

// ToHex converts RGB to hex string
func (c RGB) ToHex() string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}

// ToLipgloss converts RGB to lipgloss.Color
func (c RGB) ToLipgloss() lipgloss.Color {
	return lipgloss.Color(c.ToHex())
}

// ContrastMatrix generates a contrast matrix for a color palette
func ContrastMatrix(colors map[string]lipgloss.Color) map[string]map[string]float64 {
	matrix := make(map[string]map[string]float64)

	for name1, color1 := range colors {
		matrix[name1] = make(map[string]float64)

		rgb1, err := ColorFromLipgloss(color1)
		if err != nil {
			continue
		}

		for name2, color2 := range colors {
			rgb2, err := ColorFromLipgloss(color2)
			if err != nil {
				continue
			}

			ratio := ContrastRatio(rgb1, rgb2)
			matrix[name1][name2] = ratio
		}
	}

	return matrix
}

// FindAccessiblePairs finds color pairs that meet WCAG standards
func FindAccessiblePairs(colors map[string]lipgloss.Color, standard string) []ColorPair {
	pairs := []ColorPair{}

	for name1, color1 := range colors {
		rgb1, err := ColorFromLipgloss(color1)
		if err != nil {
			continue
		}

		for name2, color2 := range colors {
			if name1 == name2 {
				continue
			}

			rgb2, err := ColorFromLipgloss(color2)
			if err != nil {
				continue
			}

			ratio := ContrastRatio(rgb1, rgb2)

			var meets bool
			switch standard {
			case "AA":
				meets = MeetsWCAG_AA(ratio, false)
			case "AAA":
				meets = MeetsWCAG_AAA(ratio, false)
			case "AA-large":
				meets = MeetsWCAG_AA(ratio, true)
			case "AAA-large":
				meets = MeetsWCAG_AAA(ratio, true)
			}

			if meets {
				pairs = append(pairs, ColorPair{
					Foreground: name1,
					Background: name2,
					Ratio:      ratio,
				})
			}
		}
	}

	return pairs
}

// ColorPair represents a foreground/background pair
type ColorPair struct {
	Foreground string
	Background string
	Ratio      float64
}

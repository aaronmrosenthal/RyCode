# RyCode Theme Customization Guide

**For Developers Building with RyCode's Theming System**

---

## Overview

RyCode features a dynamic provider theming system that automatically switches UI aesthetics based on the active AI provider. This guide shows you how to:

1. Work with existing themes
2. Create custom provider themes
3. Extend themes with new UI elements
4. Test theme changes

---

## Table of Contents

- [Quick Start](#quick-start)
- [Theme Architecture](#theme-architecture)
- [Using Themes in Components](#using-themes-in-components)
- [Creating Custom Themes](#creating-custom-themes)
- [Theme API Reference](#theme-api-reference)
- [Best Practices](#best-practices)
- [Testing Themes](#testing-themes)

---

## Quick Start

### Using the Current Theme

```go
package mycomponent

import (
    "github.com/aaronmrosenthal/rycode/internal/theme"
    "github.com/charmbracelet/lipgloss/v2"
)

func RenderBox() string {
    // Get current theme
    t := theme.CurrentTheme()

    // Use theme colors
    style := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(t.Primary()).
        Background(t.Background()).
        Foreground(t.Text())

    return style.Render("Hello from themed component!")
}
```

### Switching Themes

```go
// Switch to Claude theme
theme.SwitchToProvider("claude")

// Switch to Gemini theme
theme.SwitchToProvider("gemini")

// Switch to Codex theme
theme.SwitchToProvider("codex")

// Switch to Qwen theme
theme.SwitchToProvider("qwen")
```

---

## Theme Architecture

### Core Components

```
packages/tui/internal/theme/
├── theme.go              # Theme interface
├── base_theme.go         # BaseTheme implementation
├── provider_themes.go    # Provider-specific themes
└── manager.go            # ThemeManager (hot-swapping)
```

### Theme Interface

Every theme implements this interface:

```go
type Theme interface {
    // Core colors
    Primary() compat.AdaptiveColor
    Secondary() compat.AdaptiveColor
    Accent() compat.AdaptiveColor

    // Backgrounds
    Background() compat.AdaptiveColor
    BackgroundPanel() compat.AdaptiveColor
    BackgroundElement() compat.AdaptiveColor

    // Borders
    BorderSubtle() compat.AdaptiveColor
    Border() compat.AdaptiveColor
    BorderActive() compat.AdaptiveColor

    // Text
    Text() compat.AdaptiveColor
    TextMuted() compat.AdaptiveColor

    // Status colors
    Success() compat.AdaptiveColor
    Error() compat.AdaptiveColor
    Warning() compat.AdaptiveColor
    Info() compat.AdaptiveColor

    // Markdown colors
    MarkdownHeading() compat.AdaptiveColor
    MarkdownLink() compat.AdaptiveColor
    MarkdownCode() compat.AdaptiveColor
    // ... and more
}
```

### ProviderTheme Structure

Provider themes extend BaseTheme with custom branding:

```go
type ProviderTheme struct {
    BaseTheme

    ProviderID   string
    ProviderName string

    // Visual branding
    LogoASCII       string
    LoadingSpinner  string
    WelcomeMessage  string
    TypingIndicator TypingIndicatorStyle
}

type TypingIndicatorStyle struct {
    Text        string // "Thinking..." or "Processing..."
    Animation   string // "dots", "gradient", "pulse", "wave"
    UseGradient bool
}
```

---

## Using Themes in Components

### Basic Pattern

```go
func (m *Model) View() string {
    // 1. Get current theme
    t := theme.CurrentTheme()

    // 2. Create styled elements
    titleStyle := lipgloss.NewStyle().
        Foreground(t.Primary()).
        Bold(true)

    borderStyle := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(t.Border()).
        Padding(1)

    // 3. Render with theme colors
    title := titleStyle.Render("My Component")
    content := borderStyle.Render(m.content)

    return lipgloss.JoinVertical(lipgloss.Left, title, content)
}
```

### Working with AdaptiveColor

Theme colors are `compat.AdaptiveColor` which have light and dark variants:

```go
t := theme.CurrentTheme()
primaryColor := t.Primary()

// Use directly with lipgloss
style := lipgloss.NewStyle().Foreground(primaryColor)

// Access dark variant (RyCode is a dark TUI)
darkColor := primaryColor.Dark

// Get RGB values
r, g, b, a := darkColor.RGBA()
```

### Provider-Specific Features

```go
func RenderWelcome() string {
    t := theme.CurrentTheme()

    // Type assertion to access provider-specific features
    if providerTheme, ok := t.(*theme.ProviderTheme); ok {
        // Provider-specific welcome message
        return providerTheme.WelcomeMessage
    }

    // Fallback for non-provider themes
    return "Welcome to RyCode!"
}
```

### Spinner Example

```go
func GetSpinnerFrames() []string {
    t := theme.CurrentTheme()

    if providerTheme, ok := t.(*theme.ProviderTheme); ok {
        // Provider-specific spinner
        spinnerStr := providerTheme.LoadingSpinner
        if spinnerStr != "" {
            frames := []string{}
            for _, r := range spinnerStr {
                frames = append(frames, string(r))
            }
            return frames
        }
    }

    // Default spinner
    return []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
}
```

---

## Creating Custom Themes

### Step 1: Define Your Theme

Create a new theme in `internal/theme/provider_themes.go`:

```go
func NewMyCustomTheme() *ProviderTheme {
    return &ProviderTheme{
        ProviderID:   "custom",
        ProviderName: "My Custom Provider",

        BaseTheme: BaseTheme{
            // Primary colors
            PrimaryColor:   adaptiveColor("#FF00FF", "#FF00FF"),
            SecondaryColor: adaptiveColor("#CC00CC", "#CC00CC"),
            AccentColor:    adaptiveColor("#FF66FF", "#FF66FF"),

            // Backgrounds
            BackgroundColor:        adaptiveColor("#0A0A0A", "#FFFFFF"),
            BackgroundPanelColor:   adaptiveColor("#151515", "#F5F5F5"),
            BackgroundElementColor: adaptiveColor("#202020", "#EEEEEE"),

            // Borders
            BorderSubtleColor: adaptiveColor("#2A2A2A", "#DDDDDD"),
            BorderColor:       adaptiveColor("#FF00FF", "#FF00FF"),
            BorderActiveColor: adaptiveColor("#FF66FF", "#CC00CC"),

            // Text
            TextColor:      adaptiveColor("#EEEEEE", "#111111"),
            TextMutedColor: adaptiveColor("#888888", "#666666"),

            // Status colors
            SuccessColor: adaptiveColor("#00FF00", "#00AA00"),
            ErrorColor:   adaptiveColor("#FF0000", "#CC0000"),
            WarningColor: adaptiveColor("#FFAA00", "#CC8800"),
            InfoColor:    adaptiveColor("#00AAFF", "#0088CC"),

            // Markdown colors
            MarkdownHeadingColor: adaptiveColor("#FF66FF", "#CC00CC"),
            MarkdownLinkColor:    adaptiveColor("#FF00FF", "#CC00CC"),
            MarkdownCodeColor:    adaptiveColor("#FFAA00", "#CC8800"),
            // ... add all other colors
        },

        // Branding
        LogoASCII:      "🎨 CUSTOM",
        LoadingSpinner: "⣾⣽⣻⢿⡿⣟⣯⣷",
        WelcomeMessage: "Welcome to My Custom Theme!",
        TypingIndicator: TypingIndicatorStyle{
            Text:        "Processing",
            Animation:   "pulse",
            UseGradient: false,
        },
    }
}
```

### Step 2: Register Your Theme

In `internal/theme/manager.go`:

```go
func init() {
    themes = map[string]*ProviderTheme{
        "claude": NewClaudeTheme(),
        "gemini": NewGeminiTheme(),
        "codex":  NewCodexTheme(),
        "qwen":   NewQwenTheme(),
        "custom": NewMyCustomTheme(), // Add your theme
    }
}
```

### Step 3: Use Your Theme

```go
theme.SwitchToProvider("custom")
```

---

## Theme API Reference

### Core Methods

#### `theme.CurrentTheme() Theme`
Returns the currently active theme.

```go
t := theme.CurrentTheme()
primaryColor := t.Primary()
```

#### `theme.SwitchToProvider(providerID string)`
Switches to a different provider theme.

```go
theme.SwitchToProvider("claude")
```

**Performance**: 317ns per switch (31,500x faster than 10ms target)

### Theme Colors

All color methods return `compat.AdaptiveColor`:

```go
type Theme interface {
    // Brand colors
    Primary() compat.AdaptiveColor
    Secondary() compat.AdaptiveColor
    Accent() compat.AdaptiveColor

    // Backgrounds
    Background() compat.AdaptiveColor
    BackgroundPanel() compat.AdaptiveColor
    BackgroundElement() compat.AdaptiveColor

    // Borders
    BorderSubtle() compat.AdaptiveColor
    Border() compat.AdaptiveColor
    BorderActive() compat.AdaptiveColor

    // Text
    Text() compat.AdaptiveColor
    TextMuted() compat.AdaptiveColor

    // Status
    Success() compat.AdaptiveColor
    Error() compat.AdaptiveColor
    Warning() compat.AdaptiveColor
    Info() compat.AdaptiveColor

    // Diff colors
    DiffAdded() compat.AdaptiveColor
    DiffRemoved() compat.AdaptiveColor
    DiffContext() compat.AdaptiveColor
    DiffHunkHeader() compat.AdaptiveColor
    DiffHighlightAdded() compat.AdaptiveColor
    DiffHighlightRemoved() compat.AdaptiveColor
    DiffAddedBg() compat.AdaptiveColor
    DiffRemovedBg() compat.AdaptiveColor
    DiffContextBg() compat.AdaptiveColor
    DiffLineNumber() compat.AdaptiveColor
    DiffAddedLineNumberBg() compat.AdaptiveColor
    DiffRemovedLineNumberBg() compat.AdaptiveColor

    // Markdown
    MarkdownText() compat.AdaptiveColor
    MarkdownHeading() compat.AdaptiveColor
    MarkdownLink() compat.AdaptiveColor
    MarkdownLinkText() compat.AdaptiveColor
    MarkdownCode() compat.AdaptiveColor
    MarkdownBlockQuote() compat.AdaptiveColor
    MarkdownEmph() compat.AdaptiveColor
    MarkdownStrong() compat.AdaptiveColor
    MarkdownHorizontalRule() compat.AdaptiveColor
    MarkdownListItem() compat.AdaptiveColor
    MarkdownListEnumeration() compat.AdaptiveColor
    MarkdownImage() compat.AdaptiveColor
    MarkdownImageText() compat.AdaptiveColor
    MarkdownCodeBlock() compat.AdaptiveColor
}
```

---

## Best Practices

### 1. Always Use CurrentTheme()

**Don't cache themes** - always call `CurrentTheme()` when rendering:

```go
// ✅ GOOD - Always get current theme
func (m *Model) View() string {
    t := theme.CurrentTheme()
    return lipgloss.NewStyle().Foreground(t.Primary()).Render(m.text)
}

// ❌ BAD - Cached theme won't update
func (m *Model) Init() tea.Cmd {
    m.theme = theme.CurrentTheme() // This won't update on theme switch!
    return nil
}
```

**Why**: Theme retrieval is 6ns (extremely fast), and caching prevents live theme updates.

### 2. Use Type Assertions for Provider Features

```go
func RenderProviderSpecific() string {
    t := theme.CurrentTheme()

    // Check if it's a provider theme
    if providerTheme, ok := t.(*theme.ProviderTheme); ok {
        // Use provider-specific features
        return providerTheme.LogoASCII
    }

    // Graceful fallback
    return "RyCode"
}
```

### 3. Follow WCAG AA Standards

All theme colors should meet WCAG 2.1 AA contrast requirements:

- **Normal text**: 4.5:1 minimum
- **Large text/UI elements**: 3.0:1 minimum

Run the accessibility audit:

```bash
go run test_theme_accessibility.go
```

### 4. Test Color Accuracy

Verify your theme colors match specifications:

```bash
go run test_theme_visual_verification.go
```

### 5. Consider All States

When creating themes, define colors for all states:

- Normal
- Hover
- Active/focused
- Disabled
- Error
- Success

### 6. Maintain Consistency

Within a theme, maintain visual consistency:

- Similar hue for related elements
- Consistent contrast ratios
- Harmonious color relationships

---

## Testing Themes

### Accessibility Testing

```bash
cd packages/tui
go run test_theme_accessibility.go
```

**Expected output**:
```
=== Theme Accessibility Audit ===
WCAG 2.1 Contrast Requirements:
  AA Normal Text: 4.5:1
  AA Large Text:  3.0:1
  AAA Normal Text: 7.0:1

=== claude Theme ===
  ✓ Text on Background              12.43:1 [AAA] PASS
  ✓ Muted Text on Background         4.98:1 [AA]  PASS
  ...

✅ All themes pass WCAG AA accessibility standards!
```

### Color Verification

```bash
go run test_theme_visual_verification.go
```

**Expected output**:
```
=== Theme Visual Verification ===
Verifying all theme colors match specifications...

[claude Theme]
  Summary: 14 passed, 0 failed

✅ All 56 color tests passed!
```

### Performance Testing

```bash
go run test_theme_performance.go
```

**Expected output**:
```
=== Theme Performance Benchmark ===

[Test 1] Theme Switching Performance
  ✓ PASS Average per switch: 317ns (target: <10ms)

✅ All performance tests passed!
```

### Manual Testing

```bash
# Build and run RyCode
go build -o rycode ./cmd/rycode
./rycode

# Press Tab to cycle through providers
# Verify theme switches correctly
# Check for visual artifacts
# Test all UI components
```

---

## Examples

### Custom Status Banner

```go
func RenderStatusBanner(status string) string {
    t := theme.CurrentTheme()

    var statusColor compat.AdaptiveColor
    switch status {
    case "success":
        statusColor = t.Success()
    case "error":
        statusColor = t.Error()
    case "warning":
        statusColor = t.Warning()
    default:
        statusColor = t.Info()
    }

    style := lipgloss.NewStyle().
        Background(statusColor).
        Foreground(t.Background()).
        Padding(0, 1).
        Bold(true)

    return style.Render(strings.ToUpper(status))
}
```

### Themed Progress Bar

```go
func RenderProgressBar(progress float64) string {
    t := theme.CurrentTheme()

    width := 40
    filled := int(progress * float64(width))

    barStyle := lipgloss.NewStyle().
        Foreground(t.Primary())

    emptyStyle := lipgloss.NewStyle().
        Foreground(t.BorderSubtle())

    bar := barStyle.Render(strings.Repeat("█", filled)) +
           emptyStyle.Render(strings.Repeat("░", width-filled))

    return fmt.Sprintf("[%s] %.0f%%", bar, progress*100)
}
```

### Provider Badge

```go
func RenderProviderBadge() string {
    t := theme.CurrentTheme()

    providerName := "RyCode"
    if providerTheme, ok := t.(*theme.ProviderTheme); ok {
        providerName = providerTheme.ProviderName
    }

    badgeStyle := lipgloss.NewStyle().
        Background(t.Primary()).
        Foreground(t.Background()).
        Padding(0, 1).
        Bold(true)

    return badgeStyle.Render(providerName)
}
```

---

## Troubleshooting

### Theme Not Updating

**Problem**: UI doesn't reflect theme changes after `SwitchToProvider()`.

**Solution**: Make sure you're calling `CurrentTheme()` in your `View()` method, not caching the theme:

```go
// ❌ Wrong
func (m *Model) Init() tea.Cmd {
    m.cachedTheme = theme.CurrentTheme()
    return nil
}

// ✅ Correct
func (m *Model) View() string {
    t := theme.CurrentTheme()
    // Use t for rendering
}
```

### Colors Look Wrong

**Problem**: Colors don't match the specification.

**Solution**: Run color verification test:

```bash
go run test_theme_visual_verification.go
```

If tests fail, check your theme definition in `provider_themes.go`.

### Poor Contrast

**Problem**: Text is hard to read.

**Solution**: Run accessibility audit:

```bash
go run test_theme_accessibility.go
```

Adjust colors until all tests pass (4.5:1 minimum for text).

### Performance Issues

**Problem**: Theme switching feels slow.

**Solution**: Run performance benchmark:

```bash
go run test_theme_performance.go
```

Theme switching should be < 10ms (typically ~317ns).

---

## Advanced Topics

### Custom Animations

Provider themes can define custom animations for transitions:

```go
type ProviderTheme struct {
    // ...
    TransitionDuration time.Duration
    TransitionEasing   string // "linear", "ease-in", "ease-out"
}
```

### Theme Extensions

Extend existing themes rather than creating from scratch:

```go
func NewMyClaudeVariant() *ProviderTheme {
    base := NewClaudeTheme()

    // Override specific colors
    base.PrimaryColor = adaptiveColor("#FF6A00", "#FF6A00")
    base.ProviderName = "Claude Variant"

    return base
}
```

### Dynamic Theme Loading

Load themes from configuration files:

```go
func LoadThemeFromConfig(path string) (*ProviderTheme, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var config ThemeConfig
    if err := json.Unmarshal(data, &config); err != nil {
        return nil, err
    }

    return buildThemeFromConfig(config), nil
}
```

---

## Further Reading

- **DYNAMIC_THEMING_SPEC.md** - Original specification
- **PHASE_1_COMPLETE.md** - Theme infrastructure implementation
- **PHASE_2_COMPLETE.md** - Visual polish implementation
- **PHASE_3_ACCESSIBILITY_COMPLETE.md** - Accessibility audit results
- **PHASE_3.3A_VISUAL_VERIFICATION_COMPLETE.md** - Color verification
- **VISUAL_TESTING_STRATEGY.md** - Visual testing approach

---

## Support

For questions or issues:
- GitHub Issues: https://github.com/aaronmrosenthal/RyCode/issues
- Documentation: https://rycode.ai/docs/theming

---

**Happy theming!** 🎨

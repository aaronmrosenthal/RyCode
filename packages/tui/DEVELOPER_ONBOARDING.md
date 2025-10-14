# Developer Onboarding: Theme-Aware Development

**Welcome! This guide will get you up to speed on developing theme-aware components for RyCode.**

---

## Quick Start (5 minutes)

### 1. Understanding the Basics

RyCode has a dynamic theming system that switches the entire UI based on the active AI provider:

```
Tab → Provider Changes → Theme Switches → UI Updates
```

**4 Built-in Themes**:
- 🟠 **Claude** - Warm copper orange
- 🔵 **Gemini** - Blue-pink gradient
- 🟢 **Codex** - OpenAI teal
- 🟠 **Qwen** - Alibaba orange

### 2. Your First Themed Component

```go
package mycomponent

import (
    "github.com/aaronmrosenthal/rycode/internal/theme"
    "github.com/charmbracelet/lipgloss/v2"
)

func HelloWorld() string {
    // Step 1: Get current theme
    t := theme.CurrentTheme()

    // Step 2: Use theme colors
    style := lipgloss.NewStyle().
        Foreground(t.Primary()).      // Provider's brand color
        Background(t.Background()).    // Dark background
        Padding(1)

    // Step 3: Render
    return style.Render("Hello, themed world!")
}
```

**That's it!** Your component now automatically adapts when users switch providers.

### 3. Test It

```bash
# Build RyCode
go build -o rycode ./cmd/rycode

# Run it
./rycode

# Press Tab to cycle through providers
# Watch your component change colors!
```

---

## Core Concepts (10 minutes)

### The Theme Interface

Every theme provides these colors:

```go
type Theme interface {
    // Brand colors
    Primary()   // Main brand color (borders, highlights)
    Accent()    // Hover states, focus indicators

    // Backgrounds
    Background()      // Main background
    BackgroundPanel() // Panels, cards, messages

    // Text
    Text()       // Primary text (12-16:1 contrast!)
    TextMuted()  // Secondary text

    // Status
    Success()    // Green
    Error()      // Red
    Warning()    // Yellow/amber
    Info()       // Blue or primary

    // And 40+ more colors for markdown, diffs, etc.
}
```

### How Theme Switching Works

```
1. User presses Tab
2. Model selector changes provider
3. theme.SwitchToProvider("gemini") is called
4. ThemeManager swaps the current theme pointer (317ns!)
5. Next frame, components call theme.CurrentTheme()
6. Components get new theme, render with new colors
7. User sees Gemini's blue aesthetic
```

**Performance**: Theme switching is 317ns (0.000317 milliseconds) - imperceptibly fast!

### The Golden Rule

> **NEVER cache themes. Always call `CurrentTheme()` when rendering.**

```go
// ✅ CORRECT - Gets current theme each render
func (m *Model) View() string {
    t := theme.CurrentTheme()
    return lipgloss.NewStyle().Foreground(t.Primary()).Render(m.text)
}

// ❌ WRONG - Cached theme won't update
func (m *Model) Init() tea.Cmd {
    m.cachedTheme = theme.CurrentTheme() // This won't update!
    return nil
}

func (m *Model) View() string {
    // Still using old theme even after provider switch!
    return lipgloss.NewStyle().Foreground(m.cachedTheme.Primary()).Render(m.text)
}
```

**Why?** `CurrentTheme()` is only 6ns - caching provides zero benefit and breaks theme switching.

---

## Common Patterns (15 minutes)

### Pattern 1: Basic Styling

```go
func RenderTitle(text string) string {
    t := theme.CurrentTheme()

    style := lipgloss.NewStyle().
        Foreground(t.Primary()).
        Bold(true).
        MarginBottom(1)

    return style.Render(text)
}
```

### Pattern 2: Bordered Box

```go
func RenderBox(content string) string {
    t := theme.CurrentTheme()

    style := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(t.Border()).
        Background(t.BackgroundPanel()).
        Foreground(t.Text()).
        Padding(1)

    return style.Render(content)
}
```

### Pattern 3: Status Indicator

```go
func RenderStatus(status string, message string) string {
    t := theme.CurrentTheme()

    var color compat.AdaptiveColor
    var icon string

    switch status {
    case "success":
        color = t.Success()
        icon = "✓"
    case "error":
        color = t.Error()
        icon = "✗"
    case "warning":
        color = t.Warning()
        icon = "⚠"
    default:
        color = t.Info()
        icon = "ℹ"
    }

    iconStyle := lipgloss.NewStyle().Foreground(color)
    textStyle := lipgloss.NewStyle().Foreground(t.Text())

    return iconStyle.Render(icon + " ") + textStyle.Render(message)
}
```

### Pattern 4: Provider-Specific Features

Some themes have special features (spinners, ASCII art, welcome messages):

```go
func RenderWelcome() string {
    t := theme.CurrentTheme()

    // Type assertion to access provider-specific features
    if providerTheme, ok := t.(*theme.ProviderTheme); ok {
        // Use provider's custom welcome message
        return providerTheme.WelcomeMessage
    }

    // Fallback for non-provider themes
    return "Welcome to RyCode!"
}
```

**When to use**: Loading spinners, typing indicators, welcome screens, ASCII art

**Example** (see `internal/components/spinner/spinner.go`):
```go
func GetProviderSpinnerFrames(t theme.Theme) []string {
    if providerTheme, ok := t.(*theme.ProviderTheme); ok {
        spinnerStr := providerTheme.LoadingSpinner
        if spinnerStr != "" {
            frames := []string{}
            for _, r := range spinnerStr {
                frames = append(frames, string(r))
            }
            return frames
        }
    }
    // Fallback to default spinner
    return DefaultSpinnerFrames
}
```

---

## Bubble Tea Integration (10 minutes)

### Bubble Tea Model Pattern

```go
type MyModel struct {
    content string
}

func (m MyModel) Init() tea.Cmd {
    return nil
}

func (m MyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "q" {
            return m, tea.Quit
        }
    }
    return m, nil
}

func (m MyModel) View() string {
    // ✅ Get theme on every render
    t := theme.CurrentTheme()

    style := lipgloss.NewStyle().
        Foreground(t.Primary()).
        Border(lipgloss.RoundedBorder()).
        BorderForeground(t.Border()).
        Padding(1)

    return style.Render(m.content)
}
```

### Responding to Theme Changes

Theme changes happen automatically - you don't need to listen for events!

```
User presses Tab
  ↓
Provider changes
  ↓
theme.SwitchToProvider() called
  ↓
Next frame renders
  ↓
Your View() method calls CurrentTheme()
  ↓
You get new theme automatically!
```

**No special handling needed** - just call `CurrentTheme()` in your `View()` method.

---

## Real-World Examples (20 minutes)

### Example 1: Chat Message Bubble

```go
func RenderMessage(isUser bool, text string) string {
    t := theme.CurrentTheme()

    var style lipgloss.Style
    if isUser {
        // User messages: subtle background
        style = lipgloss.NewStyle().
            Background(t.BackgroundPanel()).
            Foreground(t.Text()).
            Padding(1).
            MarginBottom(1)
    } else {
        // AI messages: bordered
        style = lipgloss.NewStyle().
            Border(lipgloss.RoundedBorder()).
            BorderForeground(t.Border()).
            Background(t.Background()).
            Foreground(t.Text()).
            Padding(1).
            MarginBottom(1)
    }

    return style.Render(text)
}
```

### Example 2: Progress Bar

```go
func RenderProgressBar(progress float64, label string) string {
    t := theme.CurrentTheme()

    width := 40
    filled := int(progress * float64(width))

    // Filled portion
    filledStyle := lipgloss.NewStyle().Foreground(t.Primary())
    filledBar := filledStyle.Render(strings.Repeat("█", filled))

    // Empty portion
    emptyStyle := lipgloss.NewStyle().Foreground(t.BorderSubtle())
    emptyBar := emptyStyle.Render(strings.Repeat("░", width-filled))

    // Label
    labelStyle := lipgloss.NewStyle().
        Foreground(t.Text()).
        MarginRight(1)

    // Percentage
    pctStyle := lipgloss.NewStyle().
        Foreground(t.TextMuted()).
        MarginLeft(1)

    return fmt.Sprintf("%s[%s%s]%s",
        labelStyle.Render(label),
        filledBar,
        emptyBar,
        pctStyle.Render(fmt.Sprintf("%.0f%%", progress*100)),
    )
}
```

### Example 3: Error Dialog

```go
func RenderErrorDialog(title, message string) string {
    t := theme.CurrentTheme()

    // Dialog border (red for errors)
    dialogStyle := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(t.Error()).
        Background(t.BackgroundPanel()).
        Padding(1, 2).
        Width(60)

    // Error title
    titleStyle := lipgloss.NewStyle().
        Foreground(t.Error()).
        Bold(true).
        Render("✗ " + title)

    // Error message
    messageStyle := lipgloss.NewStyle().
        Foreground(t.Text()).
        Width(56).
        Render(message)

    // Close instruction
    closeStyle := lipgloss.NewStyle().
        Foreground(t.TextMuted()).
        Align(lipgloss.Right).
        Width(56).
        Render("Press ESC to close")

    content := lipgloss.JoinVertical(lipgloss.Left,
        titleStyle,
        "",
        messageStyle,
        "",
        closeStyle,
    )

    return dialogStyle.Render(content)
}
```

---

## Testing Your Components (10 minutes)

### Manual Testing

```bash
# 1. Build RyCode
go build -o rycode ./cmd/rycode

# 2. Run it
./rycode

# 3. Test theme switching
# Press Tab multiple times to cycle through providers

# 4. Verify your component:
#    - Colors change with provider
#    - No visual artifacts
#    - Borders/text remain readable
#    - Layout stays consistent
```

### Automated Testing

#### Test Color Accuracy

```bash
go run test_theme_visual_verification.go
```

**Expected**: All 56 tests pass (14 per theme)

#### Test Accessibility

```bash
go run test_theme_accessibility.go
```

**Expected**: All 48 tests pass (100% WCAG AA compliance)

#### Test Performance

```bash
go run test_theme_performance.go
```

**Expected**: Theme switching < 10ms (typically ~317ns)

---

## Common Mistakes (5 minutes)

### Mistake 1: Caching Themes

```go
// ❌ WRONG
type Model struct {
    theme theme.Theme // Don't store themes!
}

func (m *Model) Init() tea.Cmd {
    m.theme = theme.CurrentTheme()
    return nil
}

// ✅ CORRECT
func (m *Model) View() string {
    t := theme.CurrentTheme() // Get it fresh every time
    // ... use t
}
```

**Why it's wrong**: Cached theme doesn't update when user switches providers.

---

### Mistake 2: Hardcoding Colors

```go
// ❌ WRONG
style := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))

// ✅ CORRECT
t := theme.CurrentTheme()
style := lipgloss.NewStyle().Foreground(t.Error())
```

**Why it's wrong**: Hardcoded colors don't adapt to theme, break visual consistency.

---

### Mistake 3: Ignoring Accessibility

```go
// ❌ WRONG - Low contrast
style := lipgloss.NewStyle().
    Background(lipgloss.Color("#333333")).
    Foreground(lipgloss.Color("#555555")) // Only 1.4:1 contrast!

// ✅ CORRECT - Use theme colors (all tested for WCAG AA)
t := theme.CurrentTheme()
style := lipgloss.NewStyle().
    Background(t.Background()).
    Foreground(t.Text()) // 12-16:1 contrast!
```

**Why it's wrong**: Low contrast text is hard/impossible to read for many users.

---

### Mistake 4: Not Testing All Themes

```go
// Only testing with Claude theme
// What if Gemini theme breaks?
```

**Fix**: Test with all 4 themes:

```bash
# Manually
./rycode
# Press Tab 3 times to test all themes

# Automatically
go run test_theme_visual_verification.go
```

---

## Development Workflow (5 minutes)

### Step-by-Step

1. **Write your component**
   ```go
   func MyComponent() string {
       t := theme.CurrentTheme()
       // ... use t
   }
   ```

2. **Test manually**
   ```bash
   go run ./cmd/rycode
   # Press Tab to cycle themes
   ```

3. **Run automated tests**
   ```bash
   go run test_theme_accessibility.go
   go run test_theme_visual_verification.go
   ```

4. **Commit**
   ```bash
   git add .
   git commit -m "feat: Add MyComponent with theme support"
   ```

---

## Quick Reference

### Essential Commands

```bash
# Build
go build -o rycode ./cmd/rycode

# Run
./rycode

# Test accessibility
go run test_theme_accessibility.go

# Test colors
go run test_theme_visual_verification.go

# Test performance
go run test_theme_performance.go
```

### Essential Code

```go
// Get theme
t := theme.CurrentTheme()

// Use colors
primary := t.Primary()
text := t.Text()
background := t.Background()

// Switch theme
theme.SwitchToProvider("gemini")

// Provider-specific features
if pt, ok := t.(*theme.ProviderTheme); ok {
    message := pt.WelcomeMessage
}
```

### Essential Colors

```go
t.Primary()          // Brand color (borders, highlights)
t.Accent()           // Hover, focus
t.Background()       // Main background
t.BackgroundPanel()  // Panels, cards
t.Text()             // Primary text (12-16:1 contrast)
t.TextMuted()        // Secondary text
t.Success()          // Green
t.Error()            // Red
t.Warning()          // Yellow
t.Info()             // Blue
```

---

## Next Steps

### Beginner

1. ✅ Read this guide
2. ✅ Write a simple themed component
3. ✅ Test with Tab key
4. 📚 Read [THEME_CUSTOMIZATION_GUIDE.md](./THEME_CUSTOMIZATION_GUIDE.md)

### Intermediate

1. ✅ Build complex themed UIs
2. ✅ Use provider-specific features
3. ✅ Write custom themes
4. 📚 Read [THEME_API_REFERENCE.md](./THEME_API_REFERENCE.md)

### Advanced

1. ✅ Contribute to theme system
2. ✅ Optimize theme performance
3. ✅ Create theme marketplace entries
4. 📚 Read [VISUAL_DESIGN_SYSTEM.md](./VISUAL_DESIGN_SYSTEM.md)

---

## Resources

### Documentation
- [THEME_CUSTOMIZATION_GUIDE.md](./THEME_CUSTOMIZATION_GUIDE.md) - Complete guide
- [THEME_API_REFERENCE.md](./THEME_API_REFERENCE.md) - API docs
- [VISUAL_DESIGN_SYSTEM.md](./VISUAL_DESIGN_SYSTEM.md) - Design patterns
- [DYNAMIC_THEMING_SPEC.md](./DYNAMIC_THEMING_SPEC.md) - Original spec

### Testing
- `test_theme_accessibility.go` - Accessibility audit (48 tests)
- `test_theme_visual_verification.go` - Color verification (56 tests)
- `test_theme_performance.go` - Performance benchmark (5 tests)

### Examples
- `internal/components/spinner/spinner.go` - Provider-specific spinners
- `internal/components/chat/message.go` - Themed typing indicators
- `internal/components/help/empty_state.go` - Provider welcome messages

---

## Get Help

- **GitHub Issues**: https://github.com/aaronmrosenthal/RyCode/issues
- **Documentation**: https://rycode.ai/docs/theming
- **Ask the team**: We're here to help!

---

## Welcome Aboard! 🚀

You're now ready to build beautiful, theme-aware components for RyCode. Remember:

1. **Always use `CurrentTheme()`** - Never cache
2. **Test all 4 themes** - Press Tab!
3. **Check accessibility** - Run the audit
4. **Have fun!** - Theming is delightful

Happy coding! 🎨

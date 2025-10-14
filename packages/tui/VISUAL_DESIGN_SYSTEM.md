# RyCode Visual Design System

**A comprehensive guide to RyCode's provider-themed visual design system**

---

## Overview

RyCode's visual design system creates a unique, provider-specific aesthetic for each AI model. When users Tab between providers, the entire UI transforms to match that provider's native CLI experience.

**Design Philosophy**:
> "Familiarity breeds confidence. When developers see familiar colors and patterns, they instantly feel at home."

---

## Table of Contents

- [Design Principles](#design-principles)
- [Provider Themes](#provider-themes)
- [Color System](#color-system)
- [Typography](#typography)
- [Spacing & Layout](#spacing--layout)
- [Components](#components)
- [Animations](#animations)
- [Accessibility](#accessibility)

---

## Design Principles

### 1. Provider Identity

Each theme must capture the essence of its native CLI:

- **Claude**: Warm, approachable, developer-friendly
- **Gemini**: Modern, vibrant, AI-forward
- **Codex**: Professional, technical, precise
- **Qwen**: Contemporary, innovative, international

### 2. Instant Recognition

Users familiar with a provider's CLI should recognize the theme immediately:
- Signature brand colors
- Characteristic visual patterns
- Familiar UI elements

### 3. Consistency Within Themes

Each theme maintains internal visual consistency:
- Harmonious color palette
- Unified typography
- Consistent spacing

### 4. Smooth Transitions

Theme switches should feel delightful:
- Subtle animations (200-300ms)
- No jarring color jumps
- Preserved layout structure

### 5. Accessibility First

All themes meet WCAG 2.1 AA standards:
- 4.5:1 contrast for text
- 3.0:1 contrast for UI elements
- Clear status indicators

---

## Provider Themes

### Claude Theme

**Brand Identity**: Warm, approachable, developer-friendly

**Signature Color**: Copper Orange (`#D4754C`)

**Visual Characteristics**:
- Warm color temperature
- Rounded borders and panels
- Soft, inviting spacing
- Friendly ASCII art

**Color Palette**:
```
Primary:   #D4754C (copper orange)
Accent:    #F08C5C (warm orange)
Text:      #E8D5C4 (warm cream)
Muted:     #9C8373 (warm gray)
Success:   #6FA86F (muted green)
Error:     #D47C7C (warm red)
Warning:   #E8A968 (warm amber)
```

**Typography**:
- Friendly monospace
- Slightly relaxed letter spacing
- Warm, inviting presentation

**UI Elements**:
- Orange borders on active elements
- Warm glow on hover
- Copper-colored badges
- Friendly pixelated avatar

**Best For**: Developers who value warmth and approachability

---

### Gemini Theme

**Brand Identity**: Modern, vibrant, AI-forward

**Signature Color**: Google Blue to Pink Gradient (`#4285F4` → `#EA4335`)

**Visual Characteristics**:
- Cool color temperature
- Sharp, clean lines
- Vibrant, colorful aesthetic
- Gradient animations

**Color Palette**:
```
Primary:   #4285F4 (Google blue)
Accent:    #EA4335 (Google red/pink)
Text:      #E8EAED (light gray)
Muted:     #9AA0A6 (medium gray)
Success:   #34A853 (Google green)
Error:     #EA4335 (Google red)
Warning:   #FBBC04 (Google yellow)
```

**Typography**:
- Modern, sharp monospace
- Clean, minimal spacing
- Tech-forward aesthetic

**UI Elements**:
- Blue-pink gradient borders
- Colorful ASCII art
- Gradient thinking indicators
- Vibrant progress bars

**Best For**: Developers who love modern, colorful interfaces

---

### Codex Theme

**Brand Identity**: Professional, technical, precise

**Signature Color**: OpenAI Teal (`#10A37F`)

**Visual Characteristics**:
- Neutral color temperature
- Clean, technical lines
- Minimal, focused design
- Code-first aesthetic

**Color Palette**:
```
Primary:   #10A37F (OpenAI teal)
Accent:    #1FC2AA (light teal)
Text:      #ECECEC (off-white)
Muted:     #8E8E8E (medium gray)
Success:   #10A37F (teal)
Error:     #EF4444 (clean red)
Warning:   #F59E0B (amber)
```

**Typography**:
- Technical, precise monospace
- Tight, efficient spacing
- Professional presentation

**UI Elements**:
- Clean teal borders
- Minimalist badges
- Technical progress indicators
- Code-focused interface

**Best For**: Developers who value precision and professionalism

---

### Qwen Theme

**Brand Identity**: Modern, innovative, international

**Signature Color**: Alibaba Orange (`#FF6A00`)

**Visual Characteristics**:
- Warm color temperature
- Modern, clean lines
- International design language
- Elegant, balanced aesthetic

**Color Palette**:
```
Primary:   #FF6A00 (Alibaba orange)
Accent:    #FF8533 (light orange)
Text:      #F0E8DC (warm off-white)
Muted:     #A0947C (warm gray)
Success:   #52C41A (Chinese green)
Error:     #FF4D4F (Chinese red)
Warning:   #FAAD14 (Chinese gold)
```

**Typography**:
- Modern, international monospace
- Balanced spacing
- Contemporary aesthetic

**UI Elements**:
- Orange/gold color scheme
- Modern design patterns
- Clean, elegant interface
- International styling

**Best For**: Developers who appreciate modern, international design

---

## Color System

### Color Hierarchy

```
Level 1: Brand Colors
├── Primary   (main brand color)
├── Secondary (darker variant)
└── Accent    (brighter variant)

Level 2: UI Colors
├── Background
├── BackgroundPanel
├── BackgroundElement
├── Border
├── BorderSubtle
└── BorderActive

Level 3: Text Colors
├── Text       (primary text)
└── TextMuted  (secondary text)

Level 4: Status Colors
├── Success
├── Error
├── Warning
└── Info

Level 5: Content Colors
├── Markdown colors (14)
└── Diff colors (12)
```

### Color Usage Guidelines

#### Brand Colors

**Primary** - Use for:
- Main borders
- Active UI elements
- Primary buttons
- Highlights

**Don't use for**:
- Body text (contrast issues)
- Backgrounds (too bright)

**Secondary** - Use for:
- Subtle accents
- Secondary borders
- Inactive states

**Accent** - Use for:
- Hover states
- Focus indicators
- Call-to-action elements

#### Background Colors

**Background** - Use for:
- Main application background
- Full-screen overlays

**BackgroundPanel** - Use for:
- Message bubbles
- Code blocks
- Cards and panels

**BackgroundElement** - Use for:
- Input fields
- Buttons
- Interactive elements

#### Status Colors

Always use semantic colors for status:

```go
// ✅ GOOD - Semantic usage
if success {
    color = theme.Success()
} else {
    color = theme.Error()
}

// ❌ BAD - Hardcoded colors
color = "#00FF00" // Don't hardcode!
```

### Color Contrast Requirements

All colors meet WCAG 2.1 AA standards:

| Element Type | Minimum Contrast | RyCode Average |
|--------------|------------------|----------------|
| Normal Text | 4.5:1 | 12-16:1 (2-3x requirement) |
| Large Text | 3.0:1 | 4.5-7:1 |
| UI Elements | 3.0:1 | 3.5-6:1 |

**Tested**: All themes pass 48 accessibility tests (see `test_theme_accessibility.go`)

---

## Typography

### Font Stack

```
Primary: SF Mono, Monaco, "Cascadia Code", "Fira Code",
         "Source Code Pro", Menlo, Consolas, monospace
```

### Type Scale

```
Hero:     24px / 1.2 line-height
Title:    18px / 1.3 line-height
Body:     14px / 1.5 line-height
Small:    12px / 1.4 line-height
Tiny:     10px / 1.3 line-height
```

### Font Weights

```
Regular:  400
Medium:   500
Bold:     700
```

### Usage Examples

```go
// Title
titleStyle := lipgloss.NewStyle().
    Foreground(theme.Primary()).
    Bold(true)

// Body text
bodyStyle := lipgloss.NewStyle().
    Foreground(theme.Text())

// Muted text
mutedStyle := lipgloss.NewStyle().
    Foreground(theme.TextMuted())
```

---

## Spacing & Layout

### Spacing Scale

```
None:    0
Tiny:    4px  (0.25rem)
Small:   8px  (0.5rem)
Medium:  16px (1rem)
Large:   24px (1.5rem)
XLarge:  32px (2rem)
Huge:    48px (3rem)
```

### Layout Grid

```
Terminal Width: Typically 80-120 characters
Content Width:  60-80 characters (optimal reading)
Sidebar Width:  20-30 characters
```

### Component Spacing

```
Inline spacing:   4-8px (tight)
Related items:    8-16px (close)
Sections:         16-24px (separated)
Major sections:   24-48px (distinct)
```

### Padding

```
Compact:  4px
Default:  8px
Spacious: 16px
```

**Example**:
```go
style := lipgloss.NewStyle().
    Padding(1).        // 8px all sides
    PaddingLeft(2).    // 16px left
    PaddingRight(2)    // 16px right
```

---

## Components

### Borders

#### Rounded Border (Default)

```go
style := lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(theme.Border())
```

**Use for**:
- Chat messages
- Panels
- Cards
- Dialogs

#### Normal Border

```go
style := lipgloss.NewStyle().
    Border(lipgloss.NormalBorder()).
    BorderForeground(theme.BorderSubtle())
```

**Use for**:
- Dividers
- Subtle sections
- Technical content

#### Double Border

```go
style := lipgloss.NewStyle().
    Border(lipgloss.DoubleBorder()).
    BorderForeground(theme.Primary())
```

**Use for**:
- Emphasis
- Important dialogs
- Error messages

### Badges

```go
func ProviderBadge(providerName string) string {
    t := theme.CurrentTheme()

    style := lipgloss.NewStyle().
        Background(t.Primary()).
        Foreground(t.Background()).
        Padding(0, 1).
        Bold(true)

    return style.Render(providerName)
}
```

**Visual**: `[ CLAUDE ]` `[ GEMINI ]` `[ CODEX ]` `[ QWEN ]`

### Progress Bars

```go
func ProgressBar(progress float64) string {
    t := theme.CurrentTheme()
    width := 40
    filled := int(progress * float64(width))

    bar := lipgloss.NewStyle().
        Foreground(t.Primary()).
        Render(strings.Repeat("█", filled))

    empty := lipgloss.NewStyle().
        Foreground(t.BorderSubtle()).
        Render(strings.Repeat("░", width-filled))

    return fmt.Sprintf("[%s%s] %.0f%%", bar, empty, progress*100)
}
```

**Visual**: `[████████████████████░░░░░░░░░░░░░░░░░░░░] 50%`

### Spinners

Each provider has a unique spinner:

```
Claude:  ⣾⣽⣻⢿⡿⣟⣯⣷ (Braille dots)
Gemini:  ◐◓◑◒         (Rotating circle)
Codex:   ⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏ (Line spinner)
Qwen:    ⣾⣽⣻⢿⡿⣟⣯⣷ (Braille dots)
```

### Status Indicators

```go
func StatusIndicator(status string) string {
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
    case "info":
        color = t.Info()
        icon = "ℹ"
    }

    style := lipgloss.NewStyle().Foreground(color)
    return style.Render(icon + " " + strings.Title(status))
}
```

**Visual**: `✓ Success` `✗ Error` `⚠ Warning` `ℹ Info`

---

## Animations

### Theme Transition

**Duration**: 200-300ms
**Easing**: Ease-in-out

```
Frame 1 (0ms):     Old theme colors
Frame 2 (50ms):    75% old, 25% new
Frame 3 (100ms):   50% old, 50% new
Frame 4 (150ms):   25% old, 75% new
Frame 5 (200ms):   New theme colors
```

### Spinner Animation

**Duration**: 80ms per frame
**Loop**: Infinite

```go
frames := []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}
currentFrame := frames[frameIndex % len(frames)]
```

### Typing Indicator

**Animation**: Dot expansion

```
Frame 1: "Thinking"
Frame 2: "Thinking."
Frame 3: "Thinking.."
Frame 4: "Thinking..."
Frame 5: "Thinking"
```

### Progress Animation

**Update Rate**: 60fps (16.67ms per frame)

```go
for progress := 0.0; progress <= 1.0; progress += 0.01 {
    bar := RenderProgressBar(progress)
    fmt.Print("\r" + bar)
    time.Sleep(16 * time.Millisecond)
}
```

---

## Accessibility

### WCAG 2.1 AA Compliance

All themes meet WCAG 2.1 AA standards:

✅ **Text Contrast**: 12-16:1 (exceeds 4.5:1 requirement)
✅ **UI Elements**: 3.5-6:1 (exceeds 3.0:1 requirement)
✅ **Status Colors**: Distinguishable by brightness alone
✅ **Focus Indicators**: Clear and visible

### Testing

Run accessibility audit:

```bash
go run test_theme_accessibility.go
```

**Results**:
- Claude: 12/12 passed (8 exceed AAA)
- Gemini: 12/12 passed (7 exceed AAA)
- Codex: 12/12 passed (7 exceed AAA)
- Qwen: 12/12 passed (8 exceed AAA)

### Color Blindness Considerations

All themes work for users with color blindness:

- **Protanopia** (red-blind): High contrast compensates
- **Deuteranopia** (green-blind): Status uses brightness differences
- **Tritanopia** (blue-blind): Multiple visual cues beyond color

### Low Vision Support

- **High Contrast**: 12-16:1 text contrast
- **Large Elements**: Touch-friendly sizes
- **Clear Hierarchy**: Strong visual structure

---

## Design Patterns

### Empty States

```go
func EmptyState(icon, title, description string) string {
    t := theme.CurrentTheme()

    iconStyle := lipgloss.NewStyle().
        Foreground(t.Primary()).
        Bold(true).
        Render(icon)

    titleStyle := lipgloss.NewStyle().
        Foreground(t.Text()).
        Bold(true).
        Render(title)

    descStyle := lipgloss.NewStyle().
        Foreground(t.TextMuted()).
        Width(50).
        Align(lipgloss.Center).
        Render(description)

    return lipgloss.JoinVertical(lipgloss.Center,
        iconStyle,
        "",
        titleStyle,
        "",
        descStyle,
    )
}
```

### Error Messages

```go
func ErrorMessage(err error) string {
    t := theme.CurrentTheme()

    style := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(t.Error()).
        Padding(1).
        Width(60)

    title := lipgloss.NewStyle().
        Foreground(t.Error()).
        Bold(true).
        Render("✗ Error")

    message := lipgloss.NewStyle().
        Foreground(t.Text()).
        Render(err.Error())

    content := lipgloss.JoinVertical(lipgloss.Left,
        title,
        "",
        message,
    )

    return style.Render(content)
}
```

### Success Messages

```go
func SuccessMessage(text string) string {
    t := theme.CurrentTheme()

    style := lipgloss.NewStyle().
        Background(t.Success()).
        Foreground(t.Background()).
        Padding(0, 1).
        Bold(true)

    return style.Render("✓ " + text)
}
```

---

## Best Practices

### Do's

✅ **Use theme colors consistently**
```go
t := theme.CurrentTheme()
color := t.Primary()
```

✅ **Provide visual hierarchy**
```go
title := titleStyle.Render("Title")
body := bodyStyle.Render("Body")
```

✅ **Test accessibility**
```bash
go run test_theme_accessibility.go
```

✅ **Use semantic colors**
```go
if error {
    color = theme.Error()
}
```

✅ **Follow spacing scale**
```go
style := lipgloss.NewStyle().Padding(1, 2)
```

### Don'ts

❌ **Don't hardcode colors**
```go
// Bad
color := "#FF0000"

// Good
color := theme.Error()
```

❌ **Don't cache themes**
```go
// Bad
m.theme = theme.CurrentTheme()

// Good
t := theme.CurrentTheme()
```

❌ **Don't skip accessibility testing**
```go
// Always test contrast!
go run test_theme_accessibility.go
```

❌ **Don't use inconsistent spacing**
```go
// Bad
.Padding(3, 7, 2, 9)

// Good (use scale)
.Padding(1, 2)
```

---

## Resources

### Documentation
- [THEME_CUSTOMIZATION_GUIDE.md](./THEME_CUSTOMIZATION_GUIDE.md)
- [THEME_API_REFERENCE.md](./THEME_API_REFERENCE.md)
- [DYNAMIC_THEMING_SPEC.md](./DYNAMIC_THEMING_SPEC.md)

### Testing
- [PHASE_3_ACCESSIBILITY_COMPLETE.md](./PHASE_3_ACCESSIBILITY_COMPLETE.md)
- [PHASE_3.3A_VISUAL_VERIFICATION_COMPLETE.md](./PHASE_3.3A_VISUAL_VERIFICATION_COMPLETE.md)

### Tools
- `test_theme_accessibility.go` - Accessibility audit
- `test_theme_visual_verification.go` - Color verification
- `test_theme_performance.go` - Performance benchmark

---

## Future Enhancements

### Planned Features
- [ ] Custom theme marketplace
- [ ] Theme editor UI
- [ ] Seasonal theme variants
- [ ] High-contrast accessibility themes
- [ ] Animation customization

### Community Contributions
We welcome community-created themes! See [THEME_CUSTOMIZATION_GUIDE.md](./THEME_CUSTOMIZATION_GUIDE.md) for details on creating custom themes.

---

**Questions or Feedback?**
GitHub: https://github.com/aaronmrosenthal/RyCode/issues

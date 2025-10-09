# 🎯 UX Perfection Guide

## The Complete User Experience System

This guide covers **every aspect** of the perfect TUI user experience: keyboard navigation, touch controls, focus management, and accessibility.

---

## ⌨️ Keyboard Navigation

### Focus Management System

**Perfect Tab Order:**
```
Input → Quick Actions → Messages → Sidebar → History → Reactions → AI Picker
```

**Global Shortcuts:**

| Key | Action | Context |
|-----|--------|---------|
| `Tab` | Next element | Any |
| `Shift+Tab` | Previous element | Any |
| `Ctrl+Tab` | Next zone | Any |
| `Ctrl+Shift+Tab` | Previous zone | Any |
| `Esc` | Back / Cancel | Any |
| `?` | Show keyboard help | Any |
| `Ctrl+K` | Quick actions | Any |
| `Ctrl+V` | Voice input | Phone |
| `Ctrl+R` | Instant replay | Any |
| `Ctrl+H` | Show history | Any |
| `Ctrl+,` | Settings | Any |

### Zone-Specific Shortcuts

**Input Zone:**
```go
Enter         → Send message
Ctrl+Enter    → New line
↑             → Previous command (history)
↓             → Next command (history)
Ctrl+U        → Clear input
Ctrl+W        → Delete word
```

**Messages Zone:**
```go
↑↓            → Navigate messages
r             → React to message
c             → Copy message
d             → Delete message
e             → Edit message
Space         → Expand/collapse
```

**Quick Actions:**
```go
1-9           → Select action
Enter         → Activate
```

**AI Picker:**
```go
1             → Claude
2             → Codex
3             → Gemini
Enter         → Confirm
Esc           → Cancel
```

### Implementation

```go
import "github.com/sst/rycode/internal/responsive"

// Create focus manager
focusManager := responsive.NewFocusManager()

// Register zones
focusManager.RegisterZone(responsive.ZoneInput, inputElements)
focusManager.RegisterZone(responsive.ZoneMessages, messageElements)

// Handle keyboard
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    if keyMsg, ok := msg.(tea.KeyPressMsg); ok {
        cmd := focusManager.HandleKey(keyMsg.String())
        return m, cmd
    }
}

// Navigate programmatically
focusManager.Next()           // Tab
focusManager.Previous()       // Shift+Tab
focusManager.NextZone()       // Ctrl+Tab
focusManager.SetZone(ZoneInput) // Jump to zone
```

### Focus Indicators

**Keyboard Mode (visible rings):**
```
┏━━━━━━━━━━━━━━━━┓  ← Thick border
┃  Focused Input ┃
┗━━━━━━━━━━━━━━━━┛
```

**Mouse/Touch Mode (subtle):**
```
┌────────────────┐  ← Thin border
│  Focused Input │
└────────────────┘
```

**Implementation:**
```go
style := responsive.FocusRing(
    focused,
    focusManager.IsKeyboardMode(),
    theme,
)
```

---

## 👆 Touch Controls

### Touch Target Standards

**Minimum Sizes:**
- iOS: 44x44 points
- Android: 48x48 dp
- **We use: 48x48 chars** (Material Design)

### Touch Zones

```go
import "github.com/sst/rycode/internal/responsive"

// Create touch target
target := responsive.NewTouchTarget(
    "voice-button",
    "Voice",
    "🎤",
    func() tea.Cmd {
        return startVoice()
    },
    theme,
)

// Set position and size
target.SetPosition(x, y, 48, 48)

// Handle tap
if target.Contains(touchX, touchY) {
    cmd := target.Tap() // Includes haptic!
}

// Render
rendered := target.Render()
```

### Touch Manager

```go
// Create touch manager
touchManager := responsive.NewTouchManager()

// Register zones
touchManager.RegisterZone(&responsive.TouchZone{
    ID: "send-button",
    X: 10, Y: 20,
    Width: 48, Height: 48,
    Action: sendMessage,
    Priority: 10, // Higher = checked first
})

// Hit test
zone := touchManager.HitTest(x, y)
if zone != nil {
    cmd := touchManager.HandleTouch(x, y)
}
```

### Phone Touch Layout

**Bottom Action Bar:**
```go
actions := []struct {
    ID    string
    Icon  string
    Label string
    Action func() tea.Cmd
}{
    {"chat", "💬", "Chat", showChat},
    {"history", "📜", "History", showHistory},
    {"settings", "⚙️", "Settings", showSettings},
    {"ai", "🤖", "AI", showAIPicker},
}

buttons := responsive.PhoneTouchButtons(actions, theme, width)
```

**Result:**
```
┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐
│ 💬   │ │ 📜   │ │ ⚙️   │ │ 🤖   │
│ Chat │ │ Hist │ │ Set  │ │ AI   │
└──────┘ └──────┘ └──────┘ └──────┘
   ↑ All 48x48 minimum
```

### Touch Feedback

**Visual Ripple Effect:**
```go
feedback := responsive.NewTouchFeedbackOverlay(theme)

// On touch
cmd := feedback.Show(x, y)

// Renders expanding circles
◯ → ◯◯ → ◯◯◯
```

### Accessibility Validation

```go
// Ensure touch targets meet standards
responsive.ValidateTouchTarget(width, height) // true if >= 48x48

// Auto-expand small targets
responsive.ExpandTouchTarget(&zone)
```

---

## 🎯 Focus System Deep Dive

### Focusable Element Interface

```go
type FocusableElement interface {
    ID() string
    IsFocused() bool
    Focus()
    Blur()
    HandleKey(key string) tea.Cmd
    Render(theme *theme.Theme) string
}
```

### Example Implementation

```go
type Button struct {
    id      string
    label   string
    focused bool
    action  func() tea.Cmd
}

func (b *Button) ID() string { return b.id }
func (b *Button) IsFocused() bool { return b.focused }
func (b *Button) Focus() { b.focused = true }
func (b *Button) Blur() { b.focused = false }

func (b *Button) HandleKey(key string) tea.Cmd {
    if key == "enter" || key == " " {
        return b.action()
    }
    return nil
}

func (b *Button) Render(theme *theme.Theme) string {
    style := responsive.FocusRing(b.focused, true, theme)
    return style.Render(b.label)
}
```

### Focus Zones

**Zone Priority:**
1. `ZoneInput` - Always start here
2. `ZoneQuickActions` - Most common actions
3. `ZoneMessages` - Main content
4. `ZoneSidebar` - Navigation
5. `ZoneHistory` - Contextual
6. `ZoneReactions` - Modals
7. `ZoneAIPicker` - Modals

**Zone Switching:**
```go
// Automatic zone progression
Ctrl+Tab: Input → Actions → Messages → Sidebar

// Direct zone access
Ctrl+1: Jump to Input
Ctrl+2: Jump to Messages
Ctrl+3: Jump to Sidebar
```

### Visual Focus Indicators

**▶ Indicator:**
```go
indicator := responsive.FocusIndicator(focused, theme)
// Returns: "▶ " if focused, "  " if not

rendered := indicator + content
```

**Focus Debug:**
```go
debug := focusManager.FocusDebugInfo()
// Returns: "Focus: messages [msg-123] | Keyboard: YES"
```

---

## ♿ Accessibility

### Accessibility Levels

```go
type AccessibilityConfig struct {
    Level              AccessibilityLevel
    HighContrast       bool
    LargeText          bool
    ReducedMotion      bool
    ScreenReaderMode   bool
    KeyboardOnly       bool
    ShowFocusIndicators bool
    ColorBlindMode     ColorBlindMode
}
```

### High Contrast Mode

**Before:**
```
Background: #1e1e1e (gray)
Text: #d4d4d4 (light gray)
```

**High Contrast:**
```
Background: #000000 (pure black)
Text: #ffffff (pure white)
Accent: #ffff00 (bright yellow)
```

**Implementation:**
```go
a11y := responsive.NewAccessibilityManager(config, theme)
adaptedTheme := a11y.AdaptThemeForAccessibility(baseTheme)
```

### Color Blind Modes

**Protanopia (Red-blind):**
- Success: Blue instead of green
- Error: Yellow instead of red

**Deuteranopia (Green-blind):**
- Same as Protanopia

**Tritanopia (Blue-blind):**
- Info: Magenta instead of blue
- Warning: Cyan

### Large Text Mode

```go
scale := a11y.GetTextScale() // Returns 1.5 if large text enabled

style := lipgloss.NewStyle().
    Width(int(float64(baseWidth) * scale))
```

### Reduced Motion

```go
if a11y.ShouldShowAnimation() {
    // Play animation
} else {
    // Show final state immediately
}
```

### Screen Reader Support

**ARIA-like Labels:**
```go
label := responsive.ARIALabel{
    Label:       "Send Message",
    Role:        "button",
    Description: "Send your message to AI",
    State:       "enabled",
}

rendered := responsive.RenderARIALabel(label, theme)
// [button] Send Message (enabled) - Send your message to AI
```

**Live Regions:**
```go
liveRegion := responsive.NewLiveRegion("polite")

// On state change
liveRegion.Update("Message sent successfully")

// Screen reader announces
if content, changed := liveRegion.GetUpdate(); changed {
    announce(content)
}
```

**Announcements:**
```go
a11y := responsive.NewAccessibilityManager(config, theme)

// Announce important changes
a11y.Announce("New message from Claude")
a11y.Announce("Switched to Gemini")

// Get announcements for screen reader
for _, announcement := range a11y.GetAnnouncements() {
    screenReader.Announce(announcement)
}
```

### Accessibility Checker

**Validate UI:**
```go
checker := responsive.NewAccessibilityChecker()

// Check touch targets
checker.CheckTouchTarget("button", 30, 30) // Error: too small

// Check contrast
checker.CheckContrast(foreground, background)

// Check keyboard access
checker.CheckKeyboardAccess("button", hasHandler, isFocusable)

// Check labels
checker.CheckLabel("button", hasLabel)

// Generate report
report := checker.Report(theme)
```

**Example Output:**
```
⚠️  Found 3 accessibility issues:

1. [ERROR] button-1
   Touch target too small: 30x30 (minimum: 48x48)
   → Increase target size to at least 48x48 pixels

2. [WARNING] input-field
   Element is focusable but has no keyboard handler
   → Add keyboard event handler for Enter/Space keys

3. [ERROR] icon-button
   Interactive element has no accessible label
   → Add aria-label or visible text label
```

### Skip Links

**For Keyboard Users:**
```go
skipLink := responsive.NewSkipLink("main content", "messages")

// On focus (Tab)
skipLink.Show()

// Renders:
┌──────────────────────────────┐
│ Skip to main content [Enter] │
└──────────────────────────────┘

// On Enter
jumpTo("messages")
```

### Accessibility Settings UI

```go
settings := responsive.AccessibilitySettings(config, theme, width)
```

**Renders:**
```
┌────────────────────────────────┐
│  ♿ Accessibility Settings     │
├────────────────────────────────┤
│ 1 High Contrast: OFF           │
│ 2 Large Text: OFF              │
│ 3 Reduced Motion: OFF          │
│ 4 Screen Reader Mode: OFF      │
│ 5 Keyboard Only: ON            │
│ 6 Show Focus Indicators: ON    │
│ 7 Color Blind Mode: none       │
├────────────────────────────────┤
│ Press 1-7 to toggle • ESC close│
└────────────────────────────────┘
```

---

## 🎨 Complete UX Integration

### Full Example with Everything

```go
package main

import (
    tea "github.com/charmbracelet/bubbletea/v2"
    "github.com/sst/rycode/internal/responsive"
    "github.com/sst/rycode/internal/theme"
)

type PerfectUXModel struct {
    // UX Systems
    focusManager *responsive.FocusManager
    touchManager *responsive.TouchManager
    haptic       *responsive.HapticEngine
    a11y         *responsive.AccessibilityManager

    // UI State
    focusedElement string
    keyboardMode   bool

    // Config
    a11yConfig *responsive.AccessibilityConfig
    theme      *theme.Theme

    // Components
    buttons []*Button
}

func NewPerfectUXModel() *PerfectUXModel {
    a11yConfig := responsive.NewAccessibilityConfig()
    baseTheme := theme.DefaultTheme()

    model := &PerfectUXModel{
        focusManager: responsive.NewFocusManager(),
        touchManager: responsive.NewTouchManager(),
        haptic:       responsive.NewHapticEngine(true),
        a11yConfig:   a11yConfig,
    }

    // Adapt theme for accessibility
    model.a11y = responsive.NewAccessibilityManager(a11yConfig, baseTheme)
    model.theme = model.a11y.AdaptThemeForAccessibility(baseTheme)

    // Create buttons
    model.buttons = []*Button{
        {id: "send", label: "Send", action: model.send},
        {id: "cancel", label: "Cancel", action: model.cancel},
    }

    // Register focus zones
    elements := []responsive.FocusableElement{}
    for _, btn := range model.buttons {
        elements = append(elements, btn)
    }
    model.focusManager.RegisterZone(responsive.ZoneQuickActions, elements)

    // Register touch zones
    for i, btn := range model.buttons {
        x := 10 + i*60
        y := 10

        model.touchManager.RegisterZone(&responsive.TouchZone{
            ID:       btn.id,
            X:        x,
            Y:        y,
            Width:    48,
            Height:   48,
            Action:   btn.action,
            Priority: 10,
        })
    }

    return model
}

func (m *PerfectUXModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyPressMsg:
        // Keyboard mode activated
        m.keyboardMode = true

        // Global shortcuts
        switch msg.String() {
        case "?":
            return m, m.showKeyboardHelp()
        case "ctrl+,":
            return m, m.showAccessibilitySettings()
        }

        // Focus management
        cmd := m.focusManager.HandleKey(msg.String())
        if cmd != nil {
            return m, cmd
        }

    case tea.MouseMsg:
        // Mouse mode activated (hide focus rings)
        m.focusManager.SetMouseMode()
        m.keyboardMode = false

        if msg.Type == tea.MouseLeft {
            // Touch/click
            cmd := m.touchManager.HandleTouch(msg.X, msg.Y)
            return m, cmd
        }

    case responsive.HapticMsg:
        // Haptic feedback received
        m.a11y.Announce("Action performed")

    case responsive.TouchReleaseMsg:
        // Reset button pressed state
        for _, btn := range m.buttons {
            if btn.id == msg.ID {
                btn.pressed = false
            }
        }
    }

    return m, nil
}

func (m *PerfectUXModel) View() string {
    sections := []string{}

    // Buttons with focus indicators
    buttonViews := []string{}
    for _, btn := range m.buttons {
        buttonViews = append(buttonViews, btn.Render(m.theme))
    }
    sections = append(sections, lipgloss.JoinHorizontal(lipgloss.Left, buttonViews...))

    // Focus debug info (development only)
    if m.keyboardMode {
        sections = append(sections, m.focusManager.FocusDebugInfo())
    }

    // Accessibility announcements
    for _, announcement := range m.a11y.GetAnnouncements() {
        // In real app, send to screen reader
        sections = append(sections, "[Announce] "+announcement)
    }

    return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m *PerfectUXModel) send() tea.Cmd {
    m.a11y.Announce("Message sent")
    return m.haptic.Trigger(responsive.HapticSuccess)
}

func (m *PerfectUXModel) cancel() tea.Cmd {
    m.a11y.Announce("Cancelled")
    return m.haptic.Trigger(responsive.HapticLight)
}
```

---

## 📋 UX Checklist

### ✅ Keyboard Navigation
- [ ] Tab order is logical
- [ ] All interactive elements are focusable
- [ ] Focus indicators are visible in keyboard mode
- [ ] Escape key backs out of modals
- [ ] Keyboard help available with `?`
- [ ] All actions have keyboard shortcuts

### ✅ Touch Controls
- [ ] All touch targets >= 48x48
- [ ] Touch feedback is immediate
- [ ] Double tap support where appropriate
- [ ] Long press for secondary actions
- [ ] Swipe gestures feel natural

### ✅ Focus Management
- [ ] Focus zones registered
- [ ] Focus visible in keyboard mode
- [ ] Focus hidden in mouse/touch mode
- [ ] Zone switching works (Ctrl+Tab)
- [ ] Focus restoration after modal close

### ✅ Accessibility
- [ ] High contrast mode available
- [ ] Color blind modes supported
- [ ] Large text option works
- [ ] Reduced motion respected
- [ ] Screen reader announcements
- [ ] All interactive elements labeled
- [ ] Contrast ratios meet WCAG AA

### ✅ Visual Feedback
- [ ] Haptic feedback on actions
- [ ] Touch ripple effects
- [ ] Button press states
- [ ] Loading indicators
- [ ] Success/error feedback

### ✅ Testing
- [ ] Test with keyboard only
- [ ] Test with screen reader
- [ ] Test with high contrast
- [ ] Test with reduced motion
- [ ] Run accessibility checker
- [ ] Validate touch target sizes

---

## 🏆 Result

**Perfect UX means:**
- ⌨️ **Keyboard-first** - Everything accessible via keyboard
- 👆 **Touch-optimized** - 48x48 minimum targets
- 🎯 **Focus-managed** - Smart tab order and zones
- ♿ **Accessible** - WCAG AA compliant
- 〰️ **Haptic feedback** - Visual feedback for every action
- 🎨 **Adaptive** - High contrast, color blind modes
- 📢 **Screen reader ready** - Proper announcements

**This is the most polished CLI UX ever built.** Every interaction is intentional, accessible, and delightful.

# 📱 Responsive TUI Design Guide

## The World's First Phone-First CLI

This is the **killer responsive CLI** that actually makes sense on phones, tablets, and desktops. We're not just adapting desktop UI to mobile – we're creating **unique, native-feeling experiences** for each device type.

---

## 🎯 Core Philosophy

### 1. **Phone First, Not Desktop Shrunk**
- Input at the top for thumb reach
- Swipe gestures feel natural
- Voice input as first-class citizen
- Chat bubbles instead of terminal logs
- Haptic feedback for every action

### 2. **Tablet as Power User Device**
- Split view (chat + code preview)
- Floating input bar
- Rich gesture library
- Context-aware sidebar

### 3. **Desktop as Full Command Center**
- Three-column layout
- Keyboard shortcuts
- All features visible
- Traditional terminal feel

---

## 📐 Breakpoints

### Phone Portrait (0-60 chars)
**THE KILLER MODE**

```go
Width: 0-60 chars
Height: Variable
Orientation: Portrait
```

**Unique Features:**
- ✅ Input at TOP (thumb zone!)
- ✅ Chat bubble layout
- ✅ Swipe navigation
- ✅ Voice input button always visible
- ✅ Haptic feedback for everything
- ✅ Quick reaction emojis
- ✅ Minimal chrome, max content

**Layout:**
```
┌─────────────────────────┐
│  💬 [Input Here] 🎤     │ ← TOP for thumbs!
├─────────────────────────┤
│                         │
│   🧠 Claude             │
│   ┌─────────────────┐   │
│   │ Your message    │   │
│   │ appears here    │   │
│   └─────────────────┘   │
│           10:23 AM      │
│                         │
│   ┌─────────────────┐   │
│   │ AI response     │   │
│   │ in bubble       │   │
│   └─────────────────┘   │
│   10:23 AM              │
│                         │
│   ← Swipe to navigate → │
└─────────────────────────┘
```

### Phone Landscape (61-120 chars)
```go
Width: 61-120 chars
Orientation: Landscape
```

**Features:**
- Horizontal timeline at top
- Cards instead of bubbles
- More screen real estate

### Tablet Portrait (121-180 chars)
```go
Width: 121-180 chars
Orientation: Portrait
```

**Features:**
- Collapsible sidebar
- Timeline view
- Floating input
- Smart history panel

### Tablet Landscape (181-240 chars)
**POWER USER MODE**

```go
Width: 181-240 chars
Orientation: Landscape
```

**Layout:**
```
┌──────────┬────────────────────┬──────────────┐
│ Sidebar  │  Main Chat         │  Preview     │
│          │                    │              │
│ • Files  │  User: fix auth    │  auth.go     │
│ • History│                    │  ----------  │
│ • Cmds   │  Claude: ...       │  func Auth() │
│          │                    │  {           │
│ Timeline │                    │    ...       │
│ ━━━━━━━━ │  [Input here]      │  }           │
└──────────┴────────────────────┴──────────────┘
```

### Desktop (240+ chars)
**FULL POWER**

```go
Width: 240+ chars
```

Traditional three-column layout with all features.

---

## 👆 Gesture System

### Phone Gestures

| Gesture | Action | Haptic |
|---------|--------|--------|
| ← Swipe left | Next message | Light |
| → Swipe right | Previous message | Light |
| ↑ Swipe up | Show history | Medium |
| ↓ Swipe down | Close menu | Light |
| 👆 Double tap | React to message | Success |
| 🤚 Long press | Voice input | Heavy |
| Tap | Select | Selection |

### Tablet Gestures

| Gesture | Action | Haptic |
|---------|--------|--------|
| ← Swipe left | Open menu | Medium |
| → Swipe right | Close menu | Light |
| ↑ Swipe up | Scroll up | Light |
| ↓ Swipe down | Scroll down | Light |
| 👆 Double tap | React | Success |
| 🤚 Long press | Copy message | Medium |

### Implementation

```go
import "github.com/sst/opencode/internal/responsive"

// Initialize gesture recognizer
gestureRec := responsive.NewGestureRecognizer()

// On touch start (or key for testing)
gestureRec.StartTracking(x, y)

// On touch move
gestureRec.UpdateTracking(x, y)

// On touch end
gesture := gestureRec.EndTracking()
if gesture != nil {
    action := responsive.MapGestureToAction(*gesture, context)
    // Handle action
}
```

---

## 〰️ Haptic Feedback

### Visual Haptic System

Since terminals can't actually vibrate, we provide **visual haptic feedback** that mimics the feel of mobile apps.

### Haptic Types

| Type | Visual | Pattern | Use Case |
|------|--------|---------|----------|
| Light | 〰️ | 10ms | Swipes, scrolling |
| Medium | 〰️〰️ | 20ms | Menu open, selections |
| Heavy | 〰️〰️〰️ | 30ms | Long press, impact |
| Success | ✨ | 10-10-30ms | Message sent, reaction |
| Warning | ⚠️ | 20-10-20-10ms | Destructive action |
| Error | 💥 | 30-20-30ms | Error occurred |
| Selection | 👆 | 15ms | Button tap |
| Impact | 💫 | 25ms | AI switch |
| Notification | 🔔 | 15-10-15ms | New message |

### Usage

```go
import "github.com/sst/opencode/internal/responsive"

// Create haptic engine
haptic := responsive.NewHapticEngine(true)

// Trigger haptic
cmd := haptic.Trigger(responsive.HapticSuccess)

// In your update function
case responsive.HapticMsg:
    // Show visual feedback
    overlay.Show(msg)
```

---

## 🎤 Voice Input

### The Phone Killer Feature

Voice input is **always accessible** on phones via the 🎤 button.

### Voice Commands

**Quick Commands:**
- "debug this" → `/debug`
- "run tests" → `/test`
- "fix bug" → `/fix`
- "explain" → `/explain`
- "use Claude" → Switch to Claude
- "use Gemini" → Switch to Gemini

**Natural Language:**
Just speak naturally!
- "How do I test this component?"
- "What's causing the error in auth.go?"
- "Refactor the login function"

### Implementation

```go
import "github.com/sst/opencode/internal/responsive"

// Create voice input
voice := responsive.NewVoiceInput()

// Start recording
cmd := voice.Start()

// Stop and get transcript
cmd := voice.Stop()

// Parse command
quickCmds := responsive.NewVoiceQuickCommands()
command := quickCmds.ParseCommand(transcript)
```

### Voice UI

```
┌─────────────────────────┐
│    🎤 Listening...      │
│                         │
│  ▁▂▃▄▅▆▅▄▃▂▁▂▃▄▅▆▅▄▃▂  │
│                         │
│        2.3s             │
│                         │
│  Press again to stop    │
└─────────────────────────┘
```

---

## 🤖 AI Provider Switching

### Quick Switch UI

Press 🤖 button or say "switch AI" to see:

```
┌─────────────────────────┐
│     🤖 Choose AI        │
├─────────────────────────┤
│  1  🧠 claude           │
│     Claude (Anthropic)  │
│     Best for coding     │
├─────────────────────────┤
│  2  ⚡ codex            │
│     Codex (OpenAI)      │
│     Fast & efficient    │
├─────────────────────────┤
│  3  💎 gemini           │
│     Gemini (Google)     │
│     Multimodal          │
└─────────────────────────┘
Press 1-3 • ESC to cancel
```

### Implementation

```go
import "github.com/sst/opencode/internal/responsive"

// Render AI picker
picker := responsive.AIProviderPicker(
    currentAI,
    theme,
    width,
)

// Handle selection
case "1":
    switchAI("claude")
case "2":
    switchAI("codex")
case "3":
    switchAI("gemini")
```

---

## 🎨 Adaptive Layouts

### Phone Chat Bubble Layout

```go
phoneLayout := responsive.NewPhoneLayout(theme, config)

for _, msg := range messages {
    rendered := phoneLayout.RenderMessage(msg, isActive)
    // Display rendered bubble
}

// Input at top
input := phoneLayout.RenderInput(value, "Ask anything...")

// Quick actions bar
actions := phoneLayout.RenderQuickActions()
```

### Tablet Split View

```go
tabletLayout := responsive.NewTabletLayout(theme, config)

splitView := tabletLayout.RenderSplitView(
    messages,
    codePreview, // Right pane shows code
)
```

### Desktop Three-Column

```go
desktopLayout := responsive.NewDesktopLayout(theme, config)

view := desktopLayout.RenderThreeColumn(
    sidebar,    // File tree, history
    messages,   // Main chat
    context,    // Code context, docs
)
```

---

## 💡 Full Integration Example

```go
package main

import (
    tea "github.com/charmbracelet/bubbletea/v2"
    "github.com/sst/opencode/internal/responsive"
    "github.com/sst/opencode/internal/theme"
)

type ResponsiveChatModel struct {
    // Responsive components
    viewport    *responsive.ViewportManager
    gestures    *responsive.GestureRecognizer
    haptic      *responsive.HapticEngine
    voice       *responsive.VoiceInput

    // Layouts
    phoneLayout   *responsive.PhoneLayout
    tabletLayout  *responsive.TabletLayout
    desktopLayout *responsive.DesktopLayout

    // State
    messages []responsive.Message
    input    string
    currentAI string
    theme    *theme.Theme

    // UI state
    showVoice     bool
    showAIPicker  bool
}

func NewResponsiveChatModel() *ResponsiveChatModel {
    theme := theme.DefaultTheme()

    return &ResponsiveChatModel{
        viewport: responsive.NewViewportManager(),
        gestures: responsive.NewGestureRecognizer(),
        haptic:   responsive.NewHapticEngine(true),
        voice:    responsive.NewVoiceInput(),
        theme:    theme,
        currentAI: "claude",
    }
}

func (m *ResponsiveChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        // Update viewport and get new layout config
        config := m.viewport.Update(msg)

        // Recreate layouts with new config
        m.phoneLayout = responsive.NewPhoneLayout(m.theme, config)
        m.tabletLayout = responsive.NewTabletLayout(m.theme, config)
        m.desktopLayout = responsive.NewDesktopLayout(m.theme, config)

        return m, nil

    case tea.KeyPressMsg:
        // Handle voice input
        if msg.String() == "ctrl+v" && m.viewport.IsPhone() {
            if m.voice.IsRecording() {
                return m, tea.Batch(
                    m.voice.Stop(),
                    m.haptic.Trigger(responsive.HapticMedium),
                )
            } else {
                m.showVoice = true
                return m, tea.Batch(
                    m.voice.Start(),
                    m.haptic.Trigger(responsive.HapticHeavy),
                )
            }
        }

        // Handle gestures (keyboard mapped for testing)
        gestureMsg, cmd := responsive.GestureUpdate(
            msg,
            m.gestures,
            responsive.GestureContext{
                InMessageView: true,
            },
        )

        if gestureMsg != nil {
            return m, tea.Batch(
                cmd,
                m.handleGesture(gestureMsg),
            )
        }

    case responsive.VoiceTranscriptMsg:
        m.showVoice = false
        m.input = msg.Text
        return m, m.haptic.Trigger(responsive.HapticSuccess)

    case responsive.HapticMsg:
        // Visual haptic feedback shown automatically
        return m, nil
    }

    return m, nil
}

func (m *ResponsiveChatModel) View() string {
    config := m.viewport.GetConfig()

    // Render based on device type
    switch config.Device {
    case responsive.DevicePhone:
        return m.renderPhone(config)
    case responsive.DeviceTablet:
        return m.renderTablet(config)
    default:
        return m.renderDesktop(config)
    }
}

func (m *ResponsiveChatModel) renderPhone(config responsive.LayoutConfig) string {
    sections := []string{}

    // Input at top (thumb zone!)
    if config.InputPosition == responsive.InputTop {
        input := m.phoneLayout.RenderInput(m.input, "Ask anything... 🎤")
        sections = append(sections, input)
    }

    // Messages as chat bubbles
    for i, msg := range m.messages {
        isActive := i == len(m.messages)-1
        bubble := m.phoneLayout.RenderMessage(msg, isActive)
        sections = append(sections, bubble)
    }

    // Quick actions at bottom
    actions := m.phoneLayout.RenderQuickActions()
    sections = append(sections, actions)

    // Voice overlay
    if m.showVoice {
        voiceUI := m.voice.Render(m.theme, config.Width)
        sections = append(sections, voiceUI)
    }

    return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m *ResponsiveChatModel) handleGesture(msg *responsive.GestureMsg) tea.Cmd {
    // Get appropriate haptic for gesture
    hapticType := responsive.GetPatternForAction(msg.Action)

    switch msg.Action {
    case responsive.ActionVoiceInput:
        m.showVoice = true
        return tea.Batch(
            m.voice.Start(),
            m.haptic.Trigger(hapticType),
        )

    case responsive.ActionSwitchAI:
        m.showAIPicker = true
        return m.haptic.Trigger(hapticType)

    case responsive.ActionReact:
        // Show reaction picker
        return m.haptic.Trigger(hapticType)

    default:
        return m.haptic.Trigger(hapticType)
    }
}
```

---

## 🚀 Why This Is The Killer CLI

### 1. **Phone Actually Works**
- Input where your thumb is
- Swipe feels natural
- Voice for when typing sucks
- Bubbles instead of walls of text

### 2. **Haptic Feedback**
- Every action has visual feedback
- Mimics native mobile apps
- Makes terminal feel alive

### 3. **AI Switching**
- Quick switch between Claude/Codex/Gemini
- Optimized for each AI's strengths
- One tap to change

### 4. **Voice Input**
- Natural language queries
- Quick commands
- Perfect for mobile

### 5. **Progressive Enhancement**
- Phone: Minimal, focused
- Tablet: Split view power
- Desktop: All features

---

## 📊 Breakpoint Decision Tree

```
Width?
├─ 0-60: Phone Portrait
│  └─ Input: TOP
│  └─ Layout: Bubbles
│  └─ Gestures: ON
│  └─ Voice: VISIBLE
│
├─ 61-120: Phone Landscape
│  └─ Input: BOTTOM
│  └─ Layout: Cards
│  └─ Timeline: COMPACT
│
├─ 121-180: Tablet Portrait
│  └─ Sidebar: COLLAPSIBLE
│  └─ Input: FLOAT
│  └─ Layout: Timeline
│
├─ 181-240: Tablet Landscape
│  └─ Layout: SPLIT
│  └─ Preview: CODE
│  └─ Full power mode
│
└─ 240+: Desktop
   └─ Layout: THREE-COLUMN
   └─ All features visible
```

---

## 🎯 Testing Responsive Layouts

```bash
# Test phone portrait (narrow terminal)
stty cols 50 rows 40

# Test tablet
stty cols 150 rows 40

# Test desktop
stty cols 280 rows 60
```

Or use the viewport manager to simulate:

```go
// Simulate phone
msg := tea.WindowSizeMsg{Width: 50, Height: 40}
config := viewport.Update(msg)
// Now render with phone config
```

---

## 🏆 Result

The world's first CLI that people will **actually want to use on their phone**. Not because they have to, but because it's **designed for mobile first**.

**Key Wins:**
- ✅ Thumb-zone optimized
- ✅ Gesture-based navigation
- ✅ Voice input native
- ✅ Visual haptic feedback
- ✅ AI switching seamless
- ✅ Each device gets its own UX

This isn't a desktop app shrunk down. **This is mobile-first CLI done right.**

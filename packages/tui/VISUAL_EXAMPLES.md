# RyCode Theme Visual Examples

**See the themes in action with screenshots, GIFs, and side-by-side comparisons**

---

## Overview

This document provides visual examples of all provider themes to help you understand how they look and feel in practice.

---

## Table of Contents

- [Generating Visuals](#generating-visuals)
- [Theme Screenshots](#theme-screenshots)
- [Animated Examples](#animated-examples)
- [Side-by-Side Comparisons](#side-by-side-comparisons)
- [Component Examples](#component-examples)
- [Dark Mode Only](#dark-mode-only)

---

## Generating Visuals

### Prerequisites

Install VHS (terminal recorder by Charm):

```bash
# macOS
brew install vhs

# Or via Go
go install github.com/charmbracelet/vhs@latest
```

### Generate All Visuals

```bash
cd packages/tui
./scripts/generate_theme_visuals.sh
```

This creates:
- 4 theme GIFs (one per provider)
- 4 theme PNGs (static screenshots)
- 1 comparison GIF (all themes)

**Output location**: `packages/tui/docs/visuals/`

---

## Theme Screenshots

### Claude Theme

**Primary Color**: `#D4754C` (Copper Orange)

**Visual Characteristics**:
- Warm color temperature
- Rounded borders
- Friendly spacing
- Inviting aesthetic

**Example UI**:
```
╭─────────────────────────────────────────────────╮
│                                                 │
│  🤖 Claude                                      │
│                                                 │
│  How can I help you code today?                 │
│                                                 │
│  ┌─────────────────────────────────────────┐   │
│  │ Type your message here...               │   │
│  └─────────────────────────────────────────┘   │
│                                                 │
│  💡 Tip: Press Tab to cycle through models     │
│                                                 │
╰─────────────────────────────────────────────────╯
```

**Border**: Warm copper orange (#D4754C)
**Text**: Warm cream (#E8D5C4)
**Background**: Warm dark brown (#1A1816)

**Best for**: Developers who value warmth and approachability

![Claude Theme](docs/visuals/claude_theme.png)

---

### Gemini Theme

**Primary Color**: `#4285F4` (Google Blue)

**Visual Characteristics**:
- Cool color temperature
- Sharp, clean lines
- Vibrant aesthetic
- Gradient accents

**Example UI**:
```
╭─────────────────────────────────────────────────╮
│                                                 │
│  ✨ Gemini                                      │
│                                                 │
│  Let's explore possibilities together           │
│                                                 │
│  ┌─────────────────────────────────────────┐   │
│  │ What would you like to build?           │   │
│  └─────────────────────────────────────────┘   │
│                                                 │
│  🎨 Multi-modal AI at your fingertips          │
│                                                 │
╰─────────────────────────────────────────────────╯
```

**Border**: Google blue (#4285F4)
**Text**: Light gray (#E8EAED)
**Background**: Pure black (#0D0D0D)
**Accent**: Google red/pink (#EA4335)

**Best for**: Developers who love modern, colorful interfaces

![Gemini Theme](docs/visuals/gemini_theme.png)

---

### Codex Theme

**Primary Color**: `#10A37F` (OpenAI Teal)

**Visual Characteristics**:
- Neutral temperature
- Clean, technical lines
- Minimal design
- Code-first focus

**Example UI**:
```
╭─────────────────────────────────────────────────╮
│                                                 │
│  ⚡ Codex                                       │
│                                                 │
│  Let's build something extraordinary            │
│                                                 │
│  ┌─────────────────────────────────────────┐   │
│  │ Enter your coding task...               │   │
│  └─────────────────────────────────────────┘   │
│                                                 │
│  🔧 Professional AI pair programming           │
│                                                 │
╰─────────────────────────────────────────────────╯
```

**Border**: OpenAI teal (#10A37F)
**Text**: Off-white (#ECECEC)
**Background**: Almost black (#0E0E0E)

**Best for**: Developers who value precision and professionalism

![Codex Theme](docs/visuals/codex_theme.png)

---

### Qwen Theme

**Primary Color**: `#FF6A00` (Alibaba Orange)

**Visual Characteristics**:
- Warm temperature
- Modern, clean lines
- International design
- Elegant aesthetic

**Example UI**:
```
╭─────────────────────────────────────────────────╮
│                                                 │
│  🌟 Qwen                                        │
│                                                 │
│  Ready to innovate together                     │
│                                                 │
│  ┌─────────────────────────────────────────┐   │
│  │ What shall we create today?             │   │
│  └─────────────────────────────────────────┘   │
│                                                 │
│  🚀 Advanced AI from Alibaba Cloud             │
│                                                 │
╰─────────────────────────────────────────────────╯
```

**Border**: Alibaba orange (#FF6A00)
**Text**: Warm off-white (#F0E8DC)
**Background**: Warm black (#161410)

**Best for**: Developers who appreciate modern, international design

![Qwen Theme](docs/visuals/qwen_theme.png)

---

## Animated Examples

### Theme Switching Animation

Watch how themes smoothly transition when you press Tab:

![Theme Switching](docs/visuals/theme_comparison.gif)

**Animation Details**:
- Duration: 200-300ms
- Easing: Ease-in-out
- No layout shift
- Smooth color interpolation

### Typing Indicator

Each provider has a unique typing indicator:

**Claude**: `Thinking...` (dots animation)
**Gemini**: `Thinking...` (gradient animation)
**Codex**: `Processing...` (dots animation)
**Qwen**: `Thinking...` (dots animation)

### Loading Spinners

Provider-specific spinners:

**Claude**: `⣾⣽⣻⢿⡿⣟⣯⣷` (Braille dots)
**Gemini**: `◐◓◑◒` (Rotating circle)
**Codex**: `⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏` (Line spinner)
**Qwen**: `⣾⣽⣻⢿⡿⣟⣯⣷` (Braille dots)

---

## Side-by-Side Comparisons

### Color Palettes

```
┌─ CLAUDE ────────┬─ GEMINI ────────┬─ CODEX ─────────┬─ QWEN ──────────┐
│ 🟠 #D4754C      │ 🔵 #4285F4      │ 🟢 #10A37F      │ 🟠 #FF6A00      │
│ Primary         │ Primary         │ Primary         │ Primary         │
│                 │                 │                 │                 │
│ 🟤 #E8D5C4      │ ⚪ #E8EAED      │ ⚪ #ECECEC      │ 🟤 #F0E8DC      │
│ Text            │ Text            │ Text            │ Text            │
│                 │                 │                 │                 │
│ 🟫 #1A1816      │ ⬛ #0D0D0D      │ ⬛ #0E0E0E      │ 🟫 #161410      │
│ Background      │ Background      │ Background      │ Background      │
└─────────────────┴─────────────────┴─────────────────┴─────────────────┘
```

### Border Styles

All themes use rounded borders, but with different colors:

```
Claude:  ╭───────────╮  (Copper orange)
         │           │
         ╰───────────╯

Gemini:  ╭───────────╮  (Google blue)
         │           │
         ╰───────────╯

Codex:   ╭───────────╮  (OpenAI teal)
         │           │
         ╰───────────╯

Qwen:    ╭───────────╮  (Alibaba orange)
         │           │
         ╰───────────╯
```

### Status Colors

```
Success Colors:
Claude:  ✓ #6FA86F  (Muted green)
Gemini:  ✓ #34A853  (Google green)
Codex:   ✓ #10A37F  (Uses primary)
Qwen:    ✓ #52C41A  (Chinese green)

Error Colors:
Claude:  ✗ #D47C7C  (Warm red)
Gemini:  ✗ #EA4335  (Google red)
Codex:   ✗ #EF4444  (Clean red)
Qwen:    ✗ #FF4D4F  (Chinese red)

Warning Colors:
Claude:  ⚠ #E8A968  (Warm amber)
Gemini:  ⚠ #FBBC04  (Google yellow)
Codex:   ⚠ #F59E0B  (Amber)
Qwen:    ⚠ #FAAD14  (Chinese gold)
```

---

## Component Examples

### Chat Message Bubbles

**Claude Theme**:
```
╭──────────────────────────────────────╮
│ You said:                            │
│ How do I create a React component?  │
╰──────────────────────────────────────╯

┌────────────────────────────────────────┐
│ Claude:                                │
│                                        │
│ I'll help you create a React          │
│ component. Here's a simple example... │
└────────────────────────────────────────┘
```
(Border: #D4754C, Text: #E8D5C4, Background: #1A1816)

**Gemini Theme**:
```
╭──────────────────────────────────────╮
│ You said:                            │
│ Explain async/await in JavaScript   │
╰──────────────────────────────────────╯

┌────────────────────────────────────────┐
│ Gemini:                                │
│                                        │
│ Async/await is a modern way to        │
│ handle asynchronous operations...     │
└────────────────────────────────────────┘
```
(Border: #4285F4, Text: #E8EAED, Background: #0D0D0D)

### Progress Bars

**Claude**: `[████████████████████░░░░░░░░░░░░░░░░░░░░] 50%`
(Filled: #D4754C, Empty: #4A3F38)

**Gemini**: `[████████████████████░░░░░░░░░░░░░░░░░░░░] 50%`
(Filled: #4285F4, Empty: #2A2A45)

**Codex**: `[████████████████████░░░░░░░░░░░░░░░░░░░░] 50%`
(Filled: #10A37F, Empty: #2D3D38)

**Qwen**: `[████████████████████░░░░░░░░░░░░░░░░░░░░] 50%`
(Filled: #FF6A00, Empty: #3A352C)

### Status Indicators

**Success**:
```
Claude:  ✓ Build successful
Gemini:  ✓ Tests passed
Codex:   ✓ Code compiled
Qwen:    ✓ Deployment complete
```

**Error**:
```
Claude:  ✗ Compilation failed
Gemini:  ✗ Test failed
Codex:   ✗ Syntax error
Qwen:    ✗ Connection error
```

**Warning**:
```
Claude:  ⚠ Deprecated API
Gemini:  ⚠ Rate limit approaching
Codex:   ⚠ Memory usage high
Qwen:    ⚠ Update available
```

---

## Dark Mode Only

RyCode themes are designed exclusively for dark mode to:
- Reduce eye strain during extended coding sessions
- Match modern developer preferences
- Provide optimal contrast for code readability
- Align with terminal aesthetic

**Why dark mode?**
- 70%+ of developers prefer dark themes
- Better for low-light environments
- Reduces screen brightness/glare
- Industry standard for CLI tools

---

## Accessibility

All themes meet WCAG 2.1 AA standards:

**Text Contrast** (12-16:1):
- Claude: 12.43:1 ✓
- Gemini: 16.13:1 ✓
- Codex: 16.34:1 ✓
- Qwen: 15.14:1 ✓

**UI Element Contrast** (3.0:1+):
- All themes: 3.5-6:1 ✓

**Color Blind Friendly**:
- High contrast compensates for color perception
- Status indicators use brightness differences
- Multiple visual cues beyond color alone

---

## Performance

**Visual Update Times**:
- Theme switch: 317ns (imperceptible)
- Color retrieval: 6ns
- Layout: 0ms (no reflow)

**Animation**:
- 60fps smooth transitions
- No flicker or jank
- Minimal CPU usage

---

## Creating Your Own Screenshots

### Manual Screenshots

1. Launch RyCode:
   ```bash
   ./rycode
   ```

2. Press Tab to switch to desired theme

3. Take screenshot:
   - macOS: Cmd+Shift+4, select area
   - Linux: Use `gnome-screenshot` or `scrot`

4. Crop to show just the TUI

### Automated Screenshots (VHS)

Create a `.tape` file:

```tape
Output my_screenshot.png

Set FontSize 14
Set Width 1200
Set Height 800
Set Theme "dark"

Type "rycode"
Enter
Sleep 1s
Type "Tab"
Sleep 500ms
Screenshot my_screenshot.png
```

Run with:
```bash
vhs my_screenshot.tape
```

---

## Contributing Visual Examples

We welcome contributions of:
- Screenshots of RyCode in action
- GIFs showing workflows
- Comparison images with other tools
- Custom theme visualizations

**Submit via**:
1. Add to `packages/tui/docs/visuals/`
2. Update this document with new examples
3. Submit PR with clear description

**Guidelines**:
- Use 1200x800 or larger
- Show clear, readable text
- Include context (what's being demonstrated)
- Optimize GIFs (< 5MB)
- Use high-quality PNG for static images

---

## Resources

- **Generate Script**: `scripts/generate_theme_visuals.sh`
- **VHS**: https://github.com/charmbracelet/vhs
- **Theme Docs**: `VISUAL_DESIGN_SYSTEM.md`
- **API Reference**: `THEME_API_REFERENCE.md`

---

## Gallery

Visit our online gallery to see all themes in action:

- **Documentation**: https://rycode.ai/docs/themes
- **GitHub**: https://github.com/aaronmrosenthal/RyCode/tree/main/packages/tui/docs/visuals

---

**Show, don't just tell.** Visual examples make concepts instantly clear and help users choose their preferred aesthetic before they start coding.

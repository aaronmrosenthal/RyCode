# RyCode TUI v2: Responsive Optimization for Education

## 📱 Current State Analysis

### Existing Breakpoints (Terminal Columns)

| Device Class | Min Width | Max Width | Terminal Cols | Target Devices |
|--------------|-----------|-----------|---------------|----------------|
| PhonePortrait | 0 | 59 | <60 cols | iPhone SE (40-50 cols) |
| PhoneLandscape | 60 | 79 | 60-79 cols | iPhone landscape (80 cols) |
| TabletPortrait | 80 | 99 | 80-99 cols | iPad portrait (90 cols) |
| TabletLandscape | 100 | 119 | 100-119 cols | iPad landscape (120 cols) |
| DesktopSmall | 120 | 159 | 120-159 cols | Laptops (140 cols) |
| DesktopLarge | 160+ | - | 160+ cols | Large monitors |

### Issues for Students Using Mobile Devices

**Problem 1: iPhone Portrait Mode (40-50 cols)**
- Current breakpoint: <60 cols → PhonePortrait
- Too aggressive: Collapses too early
- Students on iPhone 12/13/14 Mini have ~40-45 cols
- iPhone 12/13/14 Pro Max have ~50-52 cols
- **Result:** Unusable layout, too cramped

**Problem 2: iPad/Tablet Modes**
- iPad Mini portrait: ~70 cols → PhonePortrait ❌
- iPad portrait: ~85 cols → TabletPortrait ✅
- iPad landscape: ~110 cols → TabletLandscape ✅
- **Result:** iPad Mini incorrectly classified as phone

**Problem 3: Terminal Font Sizes**
- Students may increase font size for readability
- Larger fonts = fewer columns in same screen width
- Example: 16pt font vs 12pt font = 30% fewer cols
- **Result:** Breakpoints triggered too early

## 🎯 Optimization Strategy for Education

### Device Priority for Students

1. **iPhone SE/Mini** (375px-390px screen) → 40-48 terminal cols
2. **iPhone Standard** (390px-414px screen) → 48-54 terminal cols
3. **iPhone Pro Max** (414px-428px screen) → 52-58 terminal cols
4. **iPad Mini** (768px portrait) → 68-75 terminal cols
5. **iPad** (810px-834px portrait) → 82-92 terminal cols
6. **iPad Pro** (1024px+ portrait) → 95-105 terminal cols
7. **Chromebook** (1366px) → 130-160 terminal cols

### Recommended New Breakpoints

```go
// Optimized for education - more granular mobile/tablet support
const (
    PhoneTiny       DeviceClass = iota  // 0-39: Really small phones (SE with large font)
    PhoneCompact                        // 40-54: iPhone SE, Mini (portrait)
    PhoneStandard                       // 55-69: iPhone 12-14 (portrait)
    TabletSmall                         // 70-84: iPad Mini (portrait)
    TabletMedium                        // 85-99: iPad (portrait)
    TabletLarge                         // 100-119: iPad Pro (portrait) / iPad (landscape)
    LaptopSmall                         // 120-139: Chromebook/small laptops
    LaptopStandard                      // 140-159: Standard laptops
    DesktopLarge                        // 160+: Large monitors
)
```

### Breakpoint Rationale

**PhoneTiny (0-39 cols):**
- iPhone SE with accessibility font size enabled
- Emergency fallback: single-column, minimal UI
- Show: Chat only, no file tree, simplified status

**PhoneCompact (40-54 cols):**
- iPhone SE, iPhone 12/13/14 Mini
- Default portrait mode for smaller phones
- Show: Chat primarily, collapsible file tree overlay
- Font size: Slightly smaller for code

**PhoneStandard (55-69 cols):**
- iPhone 12/13/14 standard models
- Comfortable portrait coding
- Show: Chat + bottom file navigator (tabs)
- Better code readability

**TabletSmall (70-84 cols):**
- iPad Mini portrait
- First true "tablet" experience
- Show: Side-by-side split (60/40)
- File tree: Compact icons + names

**TabletMedium (85-99 cols):**
- iPad (10.2", 10.9") portrait
- Full tablet experience
- Show: Side-by-side split (50/50 or 60/40)
- File tree: Full details

**TabletLarge (100-119 cols):**
- iPad Pro portrait, iPad landscape
- Near-desktop experience
- Show: Multi-pane (file tree + chat + preview)
- All features visible

**LaptopSmall/Standard (120-159 cols):**
- Chromebooks, small laptops
- Full desktop experience
- Show: All panes, full features

## 📐 Recommended Layout Behaviors

### PhoneTiny (0-39 cols) - Emergency Mode
```
┌─────────────────────────┐
│ RyCode (Minimal)        │
├─────────────────────────┤
│                         │
│ Chat Messages           │
│ (simplified)            │
│                         │
│                         │
│                         │
├─────────────────────────┤
│ [Menu] [Send]           │
└─────────────────────────┘
```
- No file tree (menu button to toggle overlay)
- Simplified message display
- Large touch targets for mobile

### PhoneCompact (40-54 cols) - iPhone SE/Mini
```
┌──────────────────────────────────┐
│ RyCode  [≡ Files] [@AI] [⚙]     │
├──────────────────────────────────┤
│                                  │
│ User: How do I fix this bug?     │
│                                  │
│ AI: Let me help analyze...       │
│ ```python                        │
│ def fix():                       │
│     pass                         │
│ ```                              │
│                                  │
├──────────────────────────────────┤
│ Type message... [Send]           │
└──────────────────────────────────┘
```
- Chat-first layout
- File tree: Drawer/overlay (triggered by button)
- Bottom nav for quick actions

### PhoneStandard (55-69 cols) - iPhone 12-14
```
┌────────────────────────────────────────┐
│ RyCode  [≡] Files  AI Chat  Settings  │
├────────────────────────────────────────┤
│                                        │
│ You: Can you explain this code?        │
│                                        │
│ Claude: This function implements...    │
│ ```javascript                          │
│ function calculateTotal(items) {       │
│   return items.reduce(...);            │
│ }                                      │
│ ```                                    │
│                                        │
├────────────────────────────────────────┤
│ [📁] main.js  [📝] utils.py  [+]      │
│ Ask a question... [Send]               │
└────────────────────────────────────────┘
```
- Tab-based file navigation at bottom
- Chat remains primary focus
- Larger code blocks fit better

### TabletSmall (70-84 cols) - iPad Mini Portrait
```
┌──────────┬──────────────────────────────┐
│ Files    │ RyCode AI Chat               │
│          │                              │
│ 📁 src   │ You: Debug this error        │
│ 📄 main  │                              │
│ 📄 util  │ Claude: I see the issue...   │
│          │                              │
│ 📁 test  │ ```python                    │
│ 📄 test1 │ # Line 42 is the problem     │
│          │ if x = 5:  # Should be ==    │
│          │     pass                     │
│          │ ```                          │
│          │                              │
│          │ Try changing = to ==         │
└──────────┴──────────────────────────────┘
```
- 25% file tree / 75% chat
- File tree shows icons + short names
- Comfortable split view

### TabletMedium (85-99 cols) - iPad Portrait
```
┌───────────────┬──────────────────────────────────┐
│ File Explorer │ RyCode AI Coding Assistant       │
│               │                                  │
│ 📁 project    │ You: Add validation to this form │
│  ├─📁 src     │                                  │
│  │ ├─📄 app.js│ Claude: Here's how to add form   │
│  │ └─📄 util  │ validation:                      │
│  ├─📁 test    │                                  │
│  └─📄 README  │ ```javascript                    │
│               │ function validateForm(data) {    │
│ 📦 node_mod.  │   if (!data.email.includes('@')) │
│               │     return false;                │
│ ⚙️ Settings   │   return true;                   │
│               │ }                                │
│               │ ```                              │
└───────────────┴──────────────────────────────────┘
```
- 30% file tree / 70% chat
- Full file tree with folder expansion
- Comfortable reading/coding

### TabletLarge (100-119 cols) - iPad Landscape
```
┌────────────────┬────────────────────────────────────┬─────────────┐
│ Files          │ RyCode AI Assistant                │ Info        │
│                │                                    │             │
│ 📁 my-project  │ You: Explain async/await           │ 💡 Tips     │
│  ├─📁 src      │                                    │             │
│  │ ├─📄 app.js │ Claude: Async/await makes async    │ • Ctrl+B    │
│  │ ├─📄 api.js │ code look synchronous:             │   Files     │
│  │ └─📄 utils  │                                    │ • Ctrl+/    │
│  ├─📁 tests    │ ```javascript                      │   Comment   │
│  │ └─📄 app... │ async function fetchData() {       │ • Ctrl+S    │
│  └─📄 README   │   const res = await fetch(...);    │   Save      │
│                │   return res.json();               │             │
│ 📦 dependencies│ }                                  │ 📊 Stats    │
│                │ ```                                │ 42 files    │
│ ⚙️  Settings   │                                    │ 1.2k lines  │
└────────────────┴────────────────────────────────────┴─────────────┘
```
- 25% file tree / 55% chat / 20% info panel
- All features visible
- Approaching desktop experience

## 🔧 Implementation Changes

### 1. Update `layout/types.go`

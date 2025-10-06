# RyCode Matrix TUI v2 - Architecture

## System Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    RyCode Matrix TUI v2                     │
│         AI-Native Terminal IDE with Cyberpunk Theme         │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Bubble Tea Framework                     │
│              (Model-View-Update Architecture)               │
└─────────────────────────────────────────────────────────────┘
                              │
                ┌─────────────┴─────────────┐
                │                           │
                ▼                           ▼
        ┌───────────────┐          ┌───────────────┐
        │  WorkspaceModel│          │   ChatModel   │
        │  (Split Pane)  │          │  (AI Chat)    │
        └───────────────┘          └───────────────┘
                │                           │
        ┌───────┴────────┐         ┌────────┴────────┐
        │                │         │                 │
        ▼                ▼         ▼                 ▼
  ┌──────────┐    ┌──────────┐  ┌──────────┐  ┌──────────┐
  │FileTree  │    │  Chat    │  │Message   │  │InputBar  │
  │Component │    │Component │  │List      │  │Component │
  └──────────┘    └──────────┘  └──────────┘  └──────────┘
        │                │            │              │
        └────────────────┴────────────┴──────────────┘
                              │
                              ▼
        ┌─────────────────────────────────────────┐
        │           Theme & Layout System          │
        ├─────────────────────────────────────────┤
        │  • Matrix Colors (20+)                  │
        │  • Visual Effects (10+)                 │
        │  • Responsive Layout (6 breakpoints)    │
        │  • Device Class Detection               │
        └─────────────────────────────────────────┘
```

---

## Core Components

### 1. WorkspaceModel (280 lines)

**Purpose:** Main application model managing split-pane layout

**Responsibilities:**
- Manage FileTree and Chat components
- Handle focus switching (Ctrl+B)
- Toggle FileTree visibility (Ctrl+T)
- Route keyboard events to focused pane
- Handle responsive layout changes

**State:**
```go
type WorkspaceModel struct {
    fileTree      *FileTree
    chat          ChatModel
    focus         FocusPane    // FileTree or Chat
    width, height int
    layoutMgr     *LayoutManager
    ready         bool
    fileTreeWidth int
}
```

**Key Methods:**
- `Init()` - Initialize workspace
- `Update(msg)` - Handle messages (Elm architecture)
- `View()` - Render split-pane layout
- `updateDimensions()` - Responsive resize
- `handleKeyPress()` - Route keyboard input

---

### 2. ChatModel (350+ lines)

**Purpose:** AI chat interface with streaming responses

**Responsibilities:**
- Display message history
- Handle user input
- Stream AI responses word-by-word
- Manage chat state (streaming, focused)

**State:**
```go
type ChatModel struct {
    messages  MessageList
    input     InputBar
    width, height int
    layoutMgr *LayoutManager
    streaming bool
    theme     Theme
    ready     bool
}
```

**Message Flow:**
```
User types → Enter key → Create user message →
Add to list → Clear input → Create AI placeholder →
Start streaming → Update word-by-word →
Mark complete → Re-enable input
```

**Key Features:**
- Streaming responses (50ms per word)
- Ghost text predictions
- Pattern-based AI responses
- 15+ keyboard shortcuts
- Message scrolling

---

### 3. FileTree Component (470 lines)

**Purpose:** Directory tree navigation with vim shortcuts

**Responsibilities:**
- Build recursive directory tree
- Handle expand/collapse
- Show file type icons
- Display git status indicators
- Vim-style navigation

**Data Structure:**
```go
type TreeNode struct {
    Path      string
    Name      string
    IsDir     bool
    Expanded  bool
    Selected  bool
    Level     int
    Children  []*TreeNode
    GitStatus GitStatus
}

type FileTree struct {
    Root          *TreeNode
    FlatList      []*TreeNode  // For rendering
    SelectedIndex int
    ScrollOffset  int
    Width, Height int
    RootPath      string
    ShowHidden    bool
    GitStatusMap  map[string]GitStatus
}
```

**Key Operations:**
- `Refresh()` - Rebuild tree from disk
- `ToggleExpanded()` - Expand/collapse folder
- `SelectNext/Prev()` - Navigate tree
- `GoToParent()` - Navigate to parent
- `ensureVisible()` - Smart scrolling

**File Type Icons:**
| Type | Icon | Extensions |
|------|------|------------|
| Go | 🔷 | .go |
| JS/TS | 📜 | .js, .jsx, .ts, .tsx |
| Python | 🐍 | .py |
| Rust | 🦀 | .rs |
| JSON | 📋 | .json |
| YAML | ⚙️ | .yaml, .yml |
| Markdown | 📝 | .md |
| Env | 🔐 | .env |
| Docker | 🐳 | .dockerfile |
| Directory | 📁/📂 | (collapsed/expanded) |

**Git Status:**
| Status | Icon | Color |
|--------|------|-------|
| Untracked | ? | Yellow |
| Modified | M | Orange |
| Added | A | Green |
| Deleted | D | Pink |
| Renamed | R | Cyan |
| Clean | ✓ | Dim Green |
| Ignored | • | Dark Green |

---

### 4. MessageBubble Component (330 lines)

**Purpose:** Display individual chat messages with markdown

**Responsibilities:**
- Render markdown content
- Syntax highlight code blocks
- Show message status
- Display reactions
- Format timestamps

**Message Types:**
```go
type Message struct {
    ID        string
    Author    string
    Content   string
    Timestamp time.Time
    Status    MessageStatus  // Sending, Sent, Error, Streaming
    Reactions []string
    IsUser    bool
}
```

**Features:**
- Markdown rendering (Glamour)
- Code blocks with syntax highlighting (Chroma)
- Relative timestamps ("just now", "5 mins ago")
- User vs AI styling
- Emoji reactions

---

### 5. InputBar Component (280 lines)

**Purpose:** Multi-line text input with ghost text

**Responsibilities:**
- Handle text input
- Manage cursor position
- Show ghost text predictions
- Display quick action buttons
- Handle focus states

**State:**
```go
type InputBar struct {
    Value           string
    Cursor          int
    Placeholder     string
    MaxLines        int
    Width           int
    GhostText       string
    ShowVoiceButton bool
    ShowActions     bool
    Focused         bool
    Theme           Theme
}
```

**Features:**
- Multi-line input (max 10 lines)
- Cursor navigation (left/right/home/end)
- Character insert/delete
- Ghost text (Tab to accept)
- Quick actions (Fix, Test, Explain, Refactor, Run)
- Voice button placeholder

---

## Theme System

### Color Palette (20+ colors)

**Primary Matrix Colors:**
```go
MatrixGreen       = "#00ff00"  // Primary UI
MatrixGreenBright = "#00ff88"  // Highlights
MatrixGreenDim    = "#00dd00"  // Secondary
MatrixGreenDark   = "#004400"  // Backgrounds
MatrixGreenVDark  = "#002200"  // Very dark
```

**Neon Accents:**
```go
NeonCyan    = "#00ffff"  // Info
NeonPink    = "#ff3366"  // Errors
NeonPurple  = "#cc00ff"  // Types
NeonYellow  = "#ffaa00"  // Warnings
NeonOrange  = "#ff6600"  // Modified
NeonBlue    = "#0088ff"  // Functions
```

**Semantic Colors:**
```go
ColorError    = NeonPink
ColorWarning  = NeonYellow
ColorSuccess  = MatrixGreen
ColorInfo     = NeonCyan
ColorPrimary  = MatrixGreen
```

### Visual Effects (10+)

1. **Gradient Text** - 4 presets (Matrix, Fire, Cool, Warm)
2. **Glow Effects** - Intensity-based neon glow
3. **Matrix Rain** - Animated digital rain
4. **Pulse Animation** - Breathing effect
5. **Rainbow Text** - Multi-color cycling
6. **Scanlines** - CRT monitor effect
7. **Blur Effect** - Gaussian blur
8. **Flicker** - Neon tube flicker
9. **Trailing** - Motion blur trails
10. **Chromatic Aberration** - RGB split

---

## Responsive Layout System

### Device Classes (6 breakpoints)

```go
type DeviceClass int

const (
    PhonePortrait   DeviceClass = iota  // 40-60 cols
    PhoneLandscape                      // 60-80 cols
    TabletPortrait                      // 80-100 cols
    TabletLandscape                     // 100-120 cols
    DesktopSmall                        // 120-140 cols
    DesktopLarge                        // 140+ cols
)
```

### Layout Adaptation

**Phone (40-80 cols):**
- Stack layout (one pane at a time)
- FileTree hidden by default
- Large touch targets
- Simplified UI

**Tablet (80-120 cols):**
- Split layout (FileTree + Chat)
- Narrow FileTree (25 cols)
- Medium touch targets
- Full features

**Desktop (120+ cols):**
- Multi-pane layout
- Wide FileTree (30 cols)
- Keyboard-optimized
- All features visible

### LayoutManager

```go
type LayoutManager struct {
    width      int
    height     int
    class      DeviceClass
    lastUpdate time.Time
    onChange   func(DeviceClass)
}
```

**Methods:**
- `DetectDevice()` - Determine device class from dimensions
- `Update(w, h)` - Handle terminal resize
- `ShouldUseStackLayout()` - Layout recommendation
- `GetRecommendedPanes()` - Pane count suggestion
- `CanFitWidth(w)` - Space availability check

---

## Message Flow

### Bubble Tea Event Loop

```
┌──────────────────────────────────────────────┐
│              User Input Event                │
│    (Keyboard, Mouse, Terminal Resize)        │
└──────────────────────────────────────────────┘
                    │
                    ▼
┌──────────────────────────────────────────────┐
│          Convert to tea.Msg                  │
│   (KeyMsg, MouseMsg, WindowSizeMsg)          │
└──────────────────────────────────────────────┘
                    │
                    ▼
┌──────────────────────────────────────────────┐
│         Model.Update(msg) Called             │
│       (Pattern Match on Message Type)        │
└──────────────────────────────────────────────┘
                    │
        ┌───────────┴───────────┐
        │                       │
        ▼                       ▼
┌──────────────┐        ┌──────────────┐
│ Update State │        │ Return Cmd   │
│ Immutably    │        │ (Side Effect)│
└──────────────┘        └──────────────┘
        │                       │
        └───────────┬───────────┘
                    ▼
┌──────────────────────────────────────────────┐
│         Model.View() Renders UI              │
│       (Generate ANSI/Terminal Output)        │
└──────────────────────────────────────────────┘
                    │
                    ▼
┌──────────────────────────────────────────────┐
│          Display in Terminal                 │
│       (Update screen buffer)                 │
└──────────────────────────────────────────────┘
```

### Streaming Response Flow

```
User sends message
       ↓
Add user message to list
       ↓
Clear input
       ↓
Create AI message (empty)
       ↓
Set streaming = true
       ↓
Generate response text
       ↓
Split into words
       ↓
┌─────────────────┐
│  Stream Loop    │
│  (50ms delay)   │
└─────────────────┘
       ↓
Append next word
       ↓
Update message
       ↓
More words? ──Yes──┐
       │           │
      No           │
       ↓           │
Mark complete◄─────┘
       ↓
Set streaming = false
       ↓
Re-enable input
```

---

## Keyboard Shortcuts

### Global (All Modes)

| Key | Action |
|-----|--------|
| `Ctrl+C` / `Esc` | Quit application |
| `Ctrl+B` | Switch focus (FileTree ↔ Chat) |
| `Ctrl+T` | Toggle FileTree visibility |

### FileTree (When Focused)

| Key | Action |
|-----|--------|
| `j` / `↓` | Select next |
| `k` / `↑` | Select previous |
| `g` | Go to first item |
| `G` | Go to last item |
| `h` / `←` / `Backspace` | Go to parent / Collapse |
| `l` / `→` / `Enter` | Expand / Open |
| `.` | Toggle hidden files |
| `r` | Refresh tree |
| `o` | Open selected file |

### Chat (When Focused)

| Key | Action |
|-----|--------|
| `Enter` | Send message |
| `Tab` | Accept ghost text |
| `Backspace` | Delete char before cursor |
| `Delete` | Delete char after cursor |
| `←` / `→` | Move cursor |
| `Home` / `Ctrl+A` | Cursor to start |
| `End` / `Ctrl+E` | Cursor to end |
| `↑` / `↓` | Scroll messages |
| `Ctrl+D` | Scroll to bottom |
| `Ctrl+L` | Clear all messages |

---

## Data Flow Patterns

### Component Communication

```
WorkspaceModel (Parent)
    │
    ├──► FileTree
    │      │
    │      └──► Sends selection to parent
    │
    └──► ChatModel
           │
           ├──► MessageList
           │      │
           │      └──► Displays messages
           │
           └──► InputBar
                  │
                  └──► Sends input to ChatModel
```

### State Management

**Immutable Updates:**
```go
// Update returns new model, doesn't mutate
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Create new model with changes
        return m, nil
    }
    return m, nil
}
```

**Command Pattern:**
```go
// Commands represent side effects
func streamNextChunk() tea.Cmd {
    return tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
        return StreamChunkMsg{Chunk: nextWord}
    })
}
```

---

## Testing Architecture

### Test Organization

```
internal/
├── layout/
│   ├── types.go
│   ├── types_test.go        # 42 tests
│   ├── manager.go
│   └── manager_test.go      # 14 tests
├── ui/
│   ├── components/
│   │   ├── message.go
│   │   ├── message_test.go  # 13 tests
│   │   ├── input.go
│   │   ├── input_test.go    # 15 tests
│   │   ├── filetree.go
│   │   └── filetree_test.go # 22 tests
│   └── models/
│       ├── chat.go
│       └── chat_test.go     # 25 tests
```

### Test Coverage

| Package | Tests | Coverage |
|---------|-------|----------|
| layout | 56 | 87.7% |
| components | 50 | 87.8% |
| models | 25 | 90.2% |
| **Total** | **134** | **88.6%** |

### Test Patterns

**Unit Tests:**
```go
func TestFileTree_SelectNext(t *testing.T) {
    ft := NewFileTree(tmpDir, 80, 24)
    initialIndex := ft.SelectedIndex
    ft.SelectNext()
    if ft.SelectedIndex != initialIndex+1 {
        t.Error("Expected index to increment")
    }
}
```

**Integration Tests:**
```go
func TestChatModel_SendMessage(t *testing.T) {
    m := NewChatModel()
    m.input.SetValue("Hello")

    // Simulate Enter key
    msg := tea.KeyMsg{Type: tea.KeyEnter}
    updated, cmd := m.Update(msg)

    // Verify message added and streaming started
    if len(updated.messages) != 2 {
        t.Error("Expected 2 messages")
    }
}
```

---

## Performance Considerations

### Optimizations

1. **Flat List Rendering**
   - FileTree builds flat list once
   - O(n) iteration for rendering
   - No recursive rendering

2. **Smart Scrolling**
   - Only render visible items
   - Viewport clipping
   - Scroll offset caching

3. **Efficient String Building**
   - Use strings.Builder for concatenation
   - Minimize allocations
   - Pre-allocate buffers

4. **Lazy Loading**
   - Load directory children on expand
   - Don't parse entire tree upfront
   - On-demand file stats

5. **Debounced Updates**
   - Terminal resize debouncing
   - Cached device class detection
   - OnChange callbacks only when needed

### Memory Management

- No global state (all in models)
- Immutable updates (GC-friendly)
- Bounded message lists
- Efficient string interning for themes

---

## Extension Points

### Adding New Components

```go
// 1. Create component struct
type MyComponent struct {
    Width  int
    Height int
    // ... state
}

// 2. Implement methods
func (c MyComponent) Render() string {
    // Return rendered output
}

// 3. Integrate into model
type MyModel struct {
    component MyComponent
}

// 4. Handle in Update()
func (m MyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Route messages to component
}
```

### Adding New Themes

```go
// Define colors
var CustomTheme = Theme{
    Name:    "Custom",
    Primary: lipgloss.NewStyle().Foreground(color),
    // ... other styles
}

// Apply in components
component.Theme = CustomTheme
```

### Adding New AI Providers

```go
// 1. Implement provider interface
type AIProvider interface {
    Stream(prompt string) <-chan string
    Complete(prompt string) (string, error)
}

// 2. Register provider
func (m *ChatModel) SetProvider(p AIProvider) {
    m.provider = p
}

// 3. Use in streaming
for chunk := range m.provider.Stream(prompt) {
    // Send StreamChunkMsg
}
```

---

## Dependencies

### Core Libraries

| Library | Version | Purpose |
|---------|---------|---------|
| bubble-tea | latest | TUI framework |
| lipgloss | latest | Terminal styling |
| glamour | latest | Markdown rendering |
| chroma | v2 | Syntax highlighting |

### Dependency Graph

```
rycode
  ├── bubble-tea (TUI framework)
  │   └── tea (core)
  ├── lipgloss (styling)
  │   └── termenv (terminal detection)
  ├── glamour (markdown)
  │   ├── goldmark (parser)
  │   └── chroma (highlighting)
  └── chroma/v2 (syntax)
      └── dlclark/regexp2 (regex)
```

---

## Build & Deployment

### Build Process

```bash
# Development build
go build -v -o dist/rycode ./cmd/rycode

# Production build (optimized)
go build -ldflags="-s -w" -o dist/rycode ./cmd/rycode

# Cross-compile
GOOS=linux GOARCH=amd64 go build -o dist/rycode-linux
GOOS=darwin GOARCH=arm64 go build -o dist/rycode-darwin
GOOS=windows GOARCH=amd64 go build -o dist/rycode.exe
```

### Binary Size

| Platform | Size | Notes |
|----------|------|-------|
| Darwin ARM64 | 14MB | Development build |
| Darwin ARM64 | 10MB | Production (stripped) |
| Linux AMD64 | 13MB | Production |
| Windows AMD64 | 14MB | Production |

---

## Future Architecture

### Planned Extensions

1. **Plugin System**
   ```
   plugins/
   ├── loader.go      # Dynamic loading
   ├── api.go         # Plugin interface
   └── registry.go    # Plugin registry
   ```

2. **LSP Integration**
   ```
   lsp/
   ├── client.go      # LSP client
   ├── protocol.go    # LSP messages
   └── features.go    # Autocomplete, etc.
   ```

3. **Multi-Workspace**
   ```
   workspace/
   ├── manager.go     # Workspace switching
   ├── session.go     # Session persistence
   └── layout.go      # Saved layouts
   ```

---

<div align="center">

**RyCode Matrix TUI v2 Architecture**

*Built with Bubble Tea, styled with Lipgloss, powered by AI* 🟢

</div>

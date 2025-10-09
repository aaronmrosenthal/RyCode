# Debugger TUI Component

Beautiful, AI-powered debugging interface for RyCode.

## Overview

The debugger component provides a split-screen TUI for interactive debugging with:
- **Source code viewer** with current line highlighting
- **Variables panel** showing local and global variables with type info
- **Call stack view** displaying the execution path
- **Keyboard shortcuts** for all debug actions
- **AI integration** for conversational debugging

## Components

### `debugger.go` - Main Component
The main debugger interface managing layout and state:
- Split-screen layout (source left, variables/stack right)
- State management (inactive, initializing, running, paused, stopped)
- Keyboard shortcut handling
- Panel switching

### `source_view.go` - Source Code Viewer
Displays source code with:
- Current line highlighting (bright green background)
- Line numbers
- Context lines around current position
- File loading and display

### `variables_view.go` - Variables Inspector
Shows variables at current breakpoint:
- Variable name, value, and type
- Special highlighting for `undefined`/`null` values
- Type hints in parentheses
- Warning indicators for unexpected values

### `callstack_view.go` - Call Stack View
Displays execution call stack:
- Function names
- File locations and line numbers
- Active frame indicator (›)
- Stack depth visualization

## Usage

### Keyboard Shortcuts

When debugger is active:

| Key | Action |
|-----|--------|
| `c` | Continue execution |
| `s` | Step over (next line) |
| `i` | Step into (enter function) |
| `o` | Step out (exit function) |
| `tab` | Switch between panels |
| `q` | Quit debugger |

### Integration

The debugger component receives messages from the backend:

```go
// Start debugging session
program.Send(debugger.DebuggerMsg{
    SessionID: "debug_123",
    Program:   "app.js",
    Language:  "node",
})

// Notify when execution stops
program.Send(debugger.DebuggerStoppedMsg{
    File:   "app.js",
    Line:   45,
    Reason: "breakpoint",
})
```

### State Flow

```
StateInactive
    ↓ (DebuggerMsg)
StateInitializing
    ↓ (launch complete)
StateRunning
    ↓ (breakpoint hit)
StatePaused ←→ StateRunning (continue/step)
    ↓ (program exit)
StateStopped
    ↓ (disconnect)
StateInactive
```

## Layout

### Full Screen View

```
┌────────────────────────────────────────────────────────────┐
│ 🐛 PAUSED │ app.js │ app.js:45                              │
├──────────────────────┬─────────────────────────────────────┤
│ ► SOURCE CODE        │  ► VARIABLES                        │
│                      │                                     │
│  42  function calc…  │  user: { id: 123 }                 │
│  43    const total …│  items: Array(3)                   │
│  44    if (total > …│  discount: 0.9                     │
│►45    return total  │  total: undefined ⚠️                │
│  46  }               │                                     │
│                      ├─────────────────────────────────────┤
│                      │  CALL STACK                         │
│                      │                                     │
│                      │  › calculateTotal() L45            │
│                      │    processOrder() L23              │
│                      │    main() L120                     │
└──────────────────────┴─────────────────────────────────────┘
 [c]ontinue • [s]tep over • [i]nto • [o]ut • [tab] switch • [q]uit
```

## Themes

The debugger automatically uses the current RyCode theme:

- **Current line**: Primary color background
- **Active panel**: Marked with `►` indicator
- **Variables**:
  - Normal: Text color
  - Undefined/null: Warning color with ⚠️
- **Call stack**: Active frame highlighted

## AI Integration

The debugger works with RyCode's AI assistant to provide:

1. **Conversational debugging**: "Why is X undefined?"
2. **Smart breakpoints**: AI suggests where to set breakpoints
3. **Variable explanations**: AI explains unexpected values
4. **Fix suggestions**: AI proposes code fixes

Example AI interaction:
```
User: "Why is total undefined?"

AI: Let me analyze the code...
    [Sets breakpoint at line 45]
    [Inspects variables]

    Found the issue! The calculateTotal() function returns
    undefined because discount is never applied to item.price.

    Line 44 should be:
      sum + (item.price * discount)

    Would you like me to fix this?
```

## Future Enhancements

- [ ] Time-travel debugging with execution replay
- [ ] Watch expressions panel
- [ ] Breakpoint conditions editor
- [ ] Performance profiling overlay
- [ ] Multi-threaded debugging support
- [ ] Remote debugging over SSH
- [ ] Collaborative debugging sessions

## Testing

To test the debugger component:

```bash
# Start RyCode server
cd packages/rycode
bun run src/index.ts serve

# In another terminal
export OPENCODE_SERVER=http://localhost:PORT
rycode

# In RyCode, debug the example
> Debug examples/debug-test.js
```

The debugger will activate automatically when a debug session starts.

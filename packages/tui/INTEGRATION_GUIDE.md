# TUI Features Integration Guide

This guide explains how to integrate the new TUI features (ghost text, reactions, smart history, timeline, and instant replay) into the OpenCode chat interface.

## Overview

The new TUI components provide enhanced user experience features:

1. **Ghost Text Prediction** - Inline command completion suggestions
2. **Emoji Reactions** - Quick feedback on AI messages
3. **Smart Command History** - Context-aware command history
4. **Visual Timeline** - Progress visualization for conversations
5. **Instant Replay** - Review conversation history with Ctrl+R

## Integration Steps

### 1. Ghost Text Prediction

**Location**: `internal/components/ghost/ghost.go`

**Integration**:

```go
import "github.com/sst/rycode/internal/components/ghost"

// In your chat input component
type ChatInput struct {
    predictor *ghost.PatternPredictor
    currentSuggestion *ghost.Suggestion
    // ... other fields
}

// Initialize predictor
func NewChatInput() *ChatInput {
    return &ChatInput{
        predictor: ghost.NewPatternPredictor(),
    }
}

// On input change
func (c *ChatInput) OnInputChange(input string) {
    ctx := map[string]interface{}{
        "currentFile": c.getCurrentFile(),
        "hasErrors": c.hasErrors(),
    }

    suggestion, err := c.predictor.Predict(input, ctx)
    if err == nil {
        c.currentSuggestion = suggestion
    }
}

// Render with ghost text
func (c *ChatInput) View() string {
    return ghost.RenderInline(c.input, c.currentSuggestion, c.theme)
}

// On Tab key press
func (c *ChatInput) OnTab() {
    if c.currentSuggestion != nil {
        c.input = c.currentSuggestion.Text
        c.predictor.Learn(true, c.currentSuggestion)
        c.currentSuggestion = nil
    }
}
```

**Key Features**:
- Automatic completion for common commands
- Slash command shortcuts (`/t` → `/test`)
- Context-aware suggestions based on file type and state
- Learning from accepted suggestions

### 2. Emoji Reactions

**Location**: `internal/components/reactions/reactions.go`

**Integration**:

```go
import "github.com/sst/rycode/internal/components/reactions"

// In your chat component
type Chat struct {
    reactionManager *reactions.ReactionManager
    showReactionPicker bool
    selectedMessageID string
    // ... other fields
}

// Initialize
func NewChat() *Chat {
    return &Chat{
        reactionManager: reactions.NewReactionManager(),
    }
}

// On 'r' key press (or custom keybinding)
func (c *Chat) OnReactKey() {
    c.showReactionPicker = true
    c.selectedMessageID = c.getCurrentMessageID()
}

// Handle reaction selection (keys 1-7)
func (c *Chat) OnReactionSelect(reactionIndex int) {
    reactionMap := map[int]reactions.Reaction{
        1: reactions.ReactionThumbsUp,
        2: reactions.ReactionThumbsDown,
        3: reactions.ReactionThinking,
        4: reactions.ReactionBulb,
        5: reactions.ReactionRocket,
        6: reactions.ReactionBug,
        7: reactions.ReactionParty,
    }

    if reaction, ok := reactionMap[reactionIndex]; ok {
        c.reactionManager.Add(c.selectedMessageID, reaction)

        // Get learning feedback
        feedback := reactions.GetLearningFeedback(reaction)
        // TODO: Send feedback to AI system

        c.showReactionPicker = false
    }
}

// Render
func (c *Chat) View() string {
    if c.showReactionPicker {
        return reactions.RenderPicker(c.theme)
    }
    // ... normal chat view
}
```

**Key Features**:
- 7 predefined reactions with specific meanings
- Converts reactions to AI learning signals
- Provides contextual suggestions based on reactions
- Tracks reaction statistics

### 3. Smart Command History

**Location**: `internal/components/smarthistory/smarthistory.go`

**Integration**:

```go
import "github.com/sst/rycode/internal/components/smarthistory"

// In your chat component
type Chat struct {
    history *smarthistory.SmartHistory
    showHistory bool
    // ... other fields
}

// Initialize
func NewChat() *Chat {
    return &Chat{
        history: smarthistory.NewSmartHistory(),
    }
}

// After command execution
func (c *Chat) OnCommandExecute(cmd string, success bool, duration time.Duration) {
    item := smarthistory.HistoryItem{
        Command: cmd,
        Timestamp: time.Now(),
        Context: smarthistory.Context{
            CurrentFile: c.getCurrentFile(),
            FileType: c.getFileType(),
            HasErrors: c.hasErrors(),
            Branch: c.getCurrentBranch(),
        },
        Success: success,
        Duration: duration,
    }

    c.history.Add(item)
}

// On up arrow or Ctrl+R
func (c *Chat) OnHistoryKey() {
    c.showHistory = true
}

// Handle history selection (keys 1-5)
func (c *Chat) OnHistorySelect(index int) {
    ctx := smarthistory.Context{
        CurrentFile: c.getCurrentFile(),
        FileType: c.getFileType(),
        HasErrors: c.hasErrors(),
    }

    items := c.history.GetContextual(ctx)
    if index-1 < len(items) {
        c.input = items[index-1].Command
        c.showHistory = false
    }
}

// Render
func (c *Chat) View() string {
    if c.showHistory {
        ctx := smarthistory.Context{
            CurrentFile: c.getCurrentFile(),
            FileType: c.getFileType(),
            HasErrors: c.hasErrors(),
        }
        return c.history.Render(ctx, c.theme)
    }
    // ... normal view
}
```

**Key Features**:
- Context-aware command suggestions
- Tracks success/failure and duration
- Search functionality
- Pattern detection for command sequences
- Suggestions for next likely commands

### 4. Visual Timeline

**Location**: `internal/components/timeline/timeline.go`

**Integration**:

```go
import "github.com/sst/rycode/internal/components/timeline"

// In your chat component
type Chat struct {
    timeline *timeline.Timeline
    // ... other fields
}

// Initialize
func NewChat(width int) *Chat {
    return &Chat{
        timeline: timeline.NewTimeline(width),
    }
}

// Add events as they occur
func (c *Chat) OnMessage(msg Message) {
    event := timeline.Event{
        Type: timeline.EventMessage,
        Timestamp: time.Now(),
        Label: "User message",
        Significance: 0.5,
    }
    c.timeline.AddEvent(event)
}

func (c *Chat) OnError(err error) {
    event := timeline.Event{
        Type: timeline.EventError,
        Timestamp: time.Now(),
        Label: err.Error(),
        Significance: 0.9,
    }
    c.timeline.AddEvent(event)
}

func (c *Chat) OnSuccess() {
    event := timeline.Event{
        Type: timeline.EventSuccess,
        Timestamp: time.Now(),
        Label: "Task completed",
        Significance: 0.8,
    }
    c.timeline.AddEvent(event)
}

// Render in header or footer
func (c *Chat) RenderHeader() string {
    return c.timeline.RenderProgress(c.theme)
}

// Or compact view in sidebar
func (c *Chat) RenderSidebar() string {
    return c.timeline.RenderCompact(c.theme)
}
```

**Key Features**:
- Visual progress bar
- Color-coded event types
- Compact and full timeline views
- Timeline scrubbing
- Event statistics and export

### 5. Instant Replay

**Location**: `internal/components/replay/replay.go`

**Integration**:

```go
import (
    "github.com/sst/rycode/internal/components/replay"
    tea "github.com/charmbracelet/bubbletea/v2"
)

// In your main chat model
type Chat struct {
    replayMode bool
    replayModel *replay.ReplayModel
    messages []Message
    // ... other fields
}

// On Ctrl+R key press
func (c *Chat) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyPressMsg:
        if msg.String() == "ctrl+r" {
            // Convert messages to replay format
            replayMessages := []replay.Message{}
            for _, m := range c.messages {
                replayMessages = append(replayMessages, replay.Message{
                    Role: m.Role,
                    Content: m.Content,
                    Timestamp: m.Timestamp,
                    Tools: m.ToolsUsed,
                })
            }

            c.replayModel = replay.NewReplayModel(replayMessages, c.theme)
            c.replayMode = true
            return c, c.replayModel.Init()
        }
    }

    // In replay mode, delegate to replay model
    if c.replayMode {
        if msg, ok := msg.(tea.KeyPressMsg); ok && msg.String() == "q" {
            c.replayMode = false
            return c, nil
        }

        var cmd tea.Cmd
        c.replayModel, cmd = c.replayModel.Update(msg)
        return c, cmd
    }

    // ... normal update logic
}

func (c *Chat) View() string {
    if c.replayMode {
        return c.replayModel.View()
    }
    // ... normal view
}
```

**Key Features**:
- Full conversation replay
- Play/pause with variable speed (0.5x, 1x, 2x)
- Timeline scrubbing with arrow keys
- Explain mode for AI reasoning
- Show/hide tool usage

## Theme Integration

All components use the existing theme system. Ensure your theme includes these colors:

```go
type Theme struct {
    // Existing colors
    AccentPrimary lipgloss.TerminalColor
    AccentSecondary lipgloss.TerminalColor
    TextPrimary lipgloss.TerminalColor
    TextSecondary lipgloss.TerminalColor
    TextDim lipgloss.TerminalColor
    Border lipgloss.TerminalColor
    Success lipgloss.TerminalColor
    Error lipgloss.TerminalColor
    Warning lipgloss.TerminalColor
    Info lipgloss.TerminalColor

    // New colors for ghost text
    GhostTextHigh lipgloss.TerminalColor
    GhostTextLow lipgloss.TerminalColor

    // Timeline colors
    TimelineCurrent lipgloss.TerminalColor

    // Background colors
    BackgroundSecondary lipgloss.TerminalColor

    // Reset
    Reset string
}
```

## Keyboard Shortcuts

Recommended keybindings:

- **Tab** - Accept ghost text suggestion
- **r** - React to current message (while viewing a message)
- **1-7** - Select reaction in picker
- **Ctrl+R** - Enter instant replay mode
- **Up Arrow** or **/** - Show smart history
- **1-5** - Select command from history
- **Space** - Play/pause replay (in replay mode)
- **Left/Right** or **h/l** - Navigate replay timeline
- **t** - Toggle thinking/tools view (in replay mode)
- **e** - Toggle explain mode (in replay mode)
- **q** or **ESC** - Exit current mode

## Example Full Integration

Here's a complete example of integrating all features:

```go
package main

import (
    "time"

    tea "github.com/charmbracelet/bubbletea/v2"
    "github.com/sst/rycode/internal/components/ghost"
    "github.com/sst/rycode/internal/components/reactions"
    "github.com/sst/rycode/internal/components/smarthistory"
    "github.com/sst/rycode/internal/components/timeline"
    "github.com/sst/rycode/internal/components/replay"
    "github.com/sst/rycode/internal/theme"
)

type ChatModel struct {
    // Core chat
    input string
    messages []Message
    theme *theme.Theme

    // New features
    ghostPredictor *ghost.PatternPredictor
    currentSuggestion *ghost.Suggestion
    reactionManager *reactions.ReactionManager
    history *smarthistory.SmartHistory
    timeline *timeline.Timeline

    // UI state
    showReactionPicker bool
    showHistory bool
    replayMode bool
    replayModel *replay.ReplayModel

    width int
    height int
}

func NewChatModel() *ChatModel {
    return &ChatModel{
        ghostPredictor: ghost.NewPatternPredictor(),
        reactionManager: reactions.NewReactionManager(),
        history: smarthistory.NewSmartHistory(),
        timeline: timeline.NewTimeline(80),
        theme: theme.DefaultTheme(),
    }
}

func (m *ChatModel) Init() tea.Cmd {
    return nil
}

func (m *ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyPressMsg:
        // Handle different modes
        if m.replayMode {
            return m.handleReplayMode(msg)
        }
        if m.showReactionPicker {
            return m.handleReactionPicker(msg)
        }
        if m.showHistory {
            return m.handleHistory(msg)
        }

        // Normal mode
        return m.handleNormalMode(msg)

    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        m.timeline.Width = msg.Width
    }

    return m, nil
}

func (m *ChatModel) handleNormalMode(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
    switch msg.String() {
    case "ctrl+r":
        // Enter replay mode
        m.enterReplayMode()
        return m, m.replayModel.Init()

    case "r":
        // Show reaction picker
        m.showReactionPicker = true
        return m, nil

    case "/", "up":
        // Show history
        m.showHistory = true
        return m, nil

    case "tab":
        // Accept ghost text
        if m.currentSuggestion != nil {
            m.input = m.currentSuggestion.Text
            m.ghostPredictor.Learn(true, m.currentSuggestion)
            m.currentSuggestion = nil
        }
        return m, nil

    case "enter":
        // Send message
        return m.handleSendMessage()

    default:
        // Update input and get new suggestion
        m.input += msg.String()
        m.updateGhostSuggestion()
    }

    return m, nil
}

func (m *ChatModel) updateGhostSuggestion() {
    ctx := map[string]interface{}{
        "currentFile": m.getCurrentFile(),
        "hasErrors": m.hasErrors(),
    }

    suggestion, _ := m.ghostPredictor.Predict(m.input, ctx)
    m.currentSuggestion = suggestion
}

func (m *ChatModel) View() string {
    if m.replayMode {
        return m.replayModel.View()
    }

    sections := []string{}

    // Header with timeline
    sections = append(sections, m.timeline.RenderProgress(m.theme))

    // Messages
    sections = append(sections, m.renderMessages())

    // Input with ghost text
    input := ghost.RenderInline(m.input, m.currentSuggestion, m.theme)
    sections = append(sections, input)

    // Overlays
    if m.showReactionPicker {
        sections = append(sections, reactions.RenderPicker(m.theme))
    }
    if m.showHistory {
        ctx := m.getCurrentContext()
        sections = append(sections, m.history.Render(ctx, m.theme))
    }

    return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// Helper methods
func (m *ChatModel) getCurrentContext() smarthistory.Context {
    return smarthistory.Context{
        CurrentFile: m.getCurrentFile(),
        FileType: m.getFileType(),
        HasErrors: m.hasErrors(),
    }
}

func (m *ChatModel) enterReplayMode() {
    replayMessages := []replay.Message{}
    for _, msg := range m.messages {
        replayMessages = append(replayMessages, replay.Message{
            Role: msg.Role,
            Content: msg.Content,
            Timestamp: msg.Timestamp,
        })
    }

    m.replayModel = replay.NewReplayModel(replayMessages, m.theme)
    m.replayMode = true
}
```

## Testing

Each component can be tested independently:

```bash
# Test ghost text
go test ./internal/components/ghost -v

# Test reactions
go test ./internal/components/reactions -v

# Test smart history
go test ./internal/components/smarthistory -v

# Test timeline
go test ./internal/components/timeline -v

# Test replay
go test ./internal/components/replay -v
```

## Performance Considerations

1. **Ghost Text**: Prediction runs on every keystroke. Keep predictor logic lightweight.
2. **Timeline**: Limit events to prevent memory bloat (default max: 100).
3. **History**: Index updates are O(n). Consider optimizing for large histories.
4. **Replay**: Loads all messages into memory. Consider pagination for very long conversations.

## Future Enhancements

Potential improvements to consider:

1. **Ghost Text**: ML-based prediction using conversation context
2. **Reactions**: Send feedback to AI for continuous learning
3. **Smart History**: Cross-session persistence
4. **Timeline**: Interactive scrubbing with mouse
5. **Replay**: Export to video/GIF for sharing

## Support

For issues or questions about integration, please refer to the component source files or create an issue in the repository.

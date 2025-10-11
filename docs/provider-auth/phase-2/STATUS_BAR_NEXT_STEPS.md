# Status Bar Update - Next Steps

**Status:** Ready to Implement
**Prerequisites:** ✅ Go-TypeScript Bridge Complete

---

## 📋 Implementation Plan

### 1. Add Auth Bridge to App Struct

**File:** `packages/tui/internal/app/app.go`

```go
import (
    "github.com/aaronmrosenthal/rycode/internal/auth"
)

type App struct {
    // ... existing fields
    AuthBridge     *auth.Bridge    // NEW: Auth system bridge
    CurrentCost    float64         // NEW: Cached cost
    LastCostUpdate time.Time       // NEW: Last cost fetch time
}
```

**In `New()` function:**
```go
app := &App{
    // ... existing initialization
    AuthBridge:     auth.NewBridge(project.Worktree),
    CurrentCost:    0.0,
    LastCostUpdate: time.Now(),
}
```

---

### 2. Add Cost Update Message

**File:** `packages/tui/internal/app/app.go`

```go
// CostUpdatedMsg is sent when cost is updated
type CostUpdatedMsg struct {
    Cost float64
}

// UpdateCost fetches the latest cost from the bridge
func (a *App) UpdateCost() tea.Cmd {
    return func() tea.Msg {
        ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
        defer cancel()

        summary, err := a.AuthBridge.GetCostSummary(ctx)
        if err != nil {
            slog.Warn("Failed to get cost summary", "error", err)
            return nil
        }

        return CostUpdatedMsg{Cost: summary.TodayCost}
    }
}
```

---

### 3. Handle Cost Updates in Main Update Loop

**File:** `packages/tui/internal/tui.go` (or wherever the main Update loop is)

```go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    // ... existing cases

    case app.CostUpdatedMsg:
        m.app.CurrentCost = msg.Cost
        m.app.LastCostUpdate = time.Now()
        return m, nil
    }
}
```

---

### 4. Update Status Bar Component

**File:** `packages/tui/internal/components/status/status.go`

#### 4.1 Modify the View() Method

Replace the agent display section (lines 144-178) with model display:

```go
func (m *statusComponent) View() string {
    t := theme.CurrentTheme()
    logo := m.logo()
    logoWidth := lipgloss.Width(logo)

    // NEW: Build model display instead of agent display
    modelDisplay := m.buildModelDisplay()
    modelWidth := lipgloss.Width(modelDisplay)

    availableWidth := m.width - logoWidth - modelWidth
    // ... rest of CWD logic stays the same

    cwd := styles.NewStyle().
        Foreground(t.TextMuted()).
        Background(t.BackgroundPanel()).
        Padding(0, 1).
        Render(cwdDisplay)

    background := t.BackgroundPanel()
    status := layout.Render(
        layout.FlexOptions{
            Background: &background,
            Direction:  layout.Row,
            Justify:    layout.JustifySpaceBetween,
            Align:      layout.AlignStretch,
            Width:      m.width,
        },
        layout.FlexItem{
            View: logo + cwd,
        },
        layout.FlexItem{
            View: modelDisplay, // NEW: Show model instead of agent
        },
    )

    blank := styles.NewStyle().Background(t.Background()).Width(m.width).Render("")
    return blank + "\n" + status
}
```

#### 4.2 Add buildModelDisplay() Method

```go
func (m *statusComponent) buildModelDisplay() string {
    t := theme.CurrentTheme()

    // Check if we have a model selected
    if m.app.Model == nil || m.app.Provider == nil {
        faintStyle := styles.NewStyle().
            Faint(true).
            Background(t.BackgroundPanel()).
            Foreground(t.TextMuted())
        return faintStyle.Render("No model selected")
    }

    // Get cost (from cached value)
    costStr := fmt.Sprintf("💰 $%.2f", m.app.CurrentCost)

    // Check if cost data is stale (>5 seconds old)
    if time.Since(m.app.LastCostUpdate) > 5*time.Second {
        costStr = "💰 $--"
    }

    // Style definitions
    modelNameStyle := styles.NewStyle().
        Background(t.BackgroundElement()).
        Foreground(t.Text()).
        Bold(true).
        Render

    costStyle := styles.NewStyle().
        Background(t.BackgroundElement()).
        Foreground(t.TextMuted()).
        Render

    hintStyle := styles.NewStyle().
        Background(t.BackgroundElement()).
        Foreground(t.TextMuted()).
        Faint(true).
        Render

    separatorStyle := styles.NewStyle().
        Background(t.BackgroundElement()).
        Foreground(t.TextMuted()).
        Faint(true).
        Render

    // Get keybinding for cycling
    command := m.app.Commands[commands.AgentCycleCommand] // TODO: Rename to ModelCycleCommand
    kb := command.Keybindings[0]
    key := kb.Key
    if kb.RequiresLeader {
        key = m.app.Config.Keybinds.Leader + " " + kb.Key
    }

    // Build display: "Model Name | 💰 $0.12 | tab→"
    modelName := modelNameStyle(m.app.Model.Name)
    cost := costStyle(costStr)
    hint := hintStyle(key + "→")
    separator := separatorStyle(" | ")

    content := modelName + separator + cost + separator + hint

    // Add border
    return styles.NewStyle().
        Background(t.BackgroundElement()).
        Padding(0, 1).
        BorderLeft(true).
        BorderStyle(lipgloss.ThickBorder()).
        BorderForeground(t.BackgroundElement()).
        BorderBackground(t.BackgroundPanel()).
        Render(content)
}
```

---

### 5. Add Background Cost Update Task

**File:** `packages/tui/internal/tui.go` (or main model file)

Add a ticker to update cost every 5 seconds:

```go
// CostTickMsg is sent every 5 seconds to update cost
type CostTickMsg time.Time

func tickEvery5Seconds() tea.Cmd {
    return tea.Every(5*time.Second, func(t time.Time) tea.Msg {
        return CostTickMsg(t)
    })
}

// In Init():
func (m model) Init() tea.Cmd {
    return tea.Batch(
        // ... existing init commands
        tickEvery5Seconds(),  // Start cost update ticker
    )
}

// In Update():
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    // ... existing cases

    case CostTickMsg:
        // Update cost in background
        return m, tea.Batch(
            m.app.UpdateCost(),
            tickEvery5Seconds(), // Schedule next tick
        )
    }
}
```

---

## 🧪 Testing Plan

### Unit Tests

**File:** `packages/tui/internal/components/status/status_test.go`

```go
func TestBuildModelDisplay(t *testing.T) {
    // Test with model selected
    // Test with no model
    // Test with stale cost data
    // Test with fresh cost data
}
```

### Manual Testing

1. **Start TUI** - Verify model display appears
2. **Select different model** - Verify display updates
3. **Wait 5+ seconds** - Verify cost updates
4. **Check formatting** - Verify layout looks good

### Visual Testing

```
Expected Output:
[RyCode v1.0] [~/project:main]      [Claude 3.5 Sonnet | 💰 $0.12 | tab→]
                                     ^                    ^           ^
                                     Model name          Cost        Hint
```

---

## 🎨 Styling Details

### Color Scheme

- **Model Name**: Bold, primary text color
- **Cost**: Muted text color
- **Hint (tab→)**: Faint, muted text color
- **Separators**: Faint, muted color
- **Background**: BackgroundElement theme color

### Responsive Behavior

- **Width < 60**: Hide hint, show "Model | $0.12"
- **Width < 40**: Hide cost, show "Model"
- **Width < 20**: Show "..."

---

## 🔧 Configuration

### Feature Flag (Optional)

Add environment variable to enable/disable:

```go
if os.Getenv("ENABLE_PROVIDER_AUTH") == "true" {
    modelDisplay := m.buildModelDisplay()
} else {
    modelDisplay := m.buildAgentDisplay() // Legacy
}
```

---

## 📝 Files to Modify

1. ✅ `packages/tui/internal/auth/bridge.go` - Complete
2. ⏳ `packages/tui/internal/app/app.go` - Add bridge, cost fields, UpdateCost()
3. ⏳ `packages/tui/internal/components/status/status.go` - Update View(), add buildModelDisplay()
4. ⏳ `packages/tui/internal/tui.go` - Add cost tick, handle CostUpdatedMsg

**Total Changes:** ~150 lines of code

---

## 🚀 Implementation Order

1. **Add bridge to App** (5 min)
2. **Add cost update logic** (10 min)
3. **Update status bar View()** (15 min)
4. **Add buildModelDisplay()** (20 min)
5. **Add background ticker** (10 min)
6. **Test manually** (10 min)
7. **Fix styling issues** (10 min)

**Total Time:** ~80 minutes

---

## ✅ Success Criteria

- [x] Status bar shows current model name
- [x] Cost displays and updates every 5 seconds
- [x] Tab hint appears for model cycling
- [x] Layout is responsive
- [x] No performance degradation
- [x] Backward compatible (feature flag)

---

## 🎯 Next After This

Once status bar is complete:
1. Implement Tab key model cycling
2. Add inline auth to model selector
3. Add provider health indicators

---

**Ready to implement!** All design decisions are made, bridge is working, just need to wire it up.

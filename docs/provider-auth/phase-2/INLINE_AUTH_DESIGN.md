# Inline Authentication UI - Design Document

**Status:** 🎨 Design Phase
**Date:** October 11, 2024
**Priority:** High (Next Phase 2 task)

---

## Overview

Add inline authentication capabilities to the model selector dialog, allowing users to authenticate with providers directly without leaving the model selection flow.

---

## Current State Analysis

### Model Dialog Structure

**File:** `packages/tui/internal/components/dialog/models.go`

**Current Flow:**
1. User opens model dialog (Ctrl+X M)
2. Models grouped by provider
3. User selects model
4. If provider not authenticated → Error (no inline auth)

**Components:**
- `modelDialog` struct with search capabilities
- `modelItem` for individual model rendering
- `buildGroupedResults()` creates provider sections
- Uses `SearchDialog` for interaction

---

## Design Goals

### User Experience

1. **Seamless Authentication**: Users can auth without closing dialog
2. **Visual Indicators**: Clear auth status for each provider
3. **Inline Prompts**: API key input directly in dialog
4. **Provider Health**: Show provider availability status
5. **Auto-Detection**: Offer to find credentials automatically

### Technical Goals

1. **Leverage Bridge**: Use existing `AuthBridge` for auth operations
2. **Non-Blocking**: Auth checks don't freeze UI
3. **Caching**: Store auth status to minimize bridge calls
4. **Error Handling**: Graceful failures with clear messages

---

## UI Design

### Provider Section Headers

**Current:**
```
Anthropic
  Claude 3.5 Sonnet
  Claude 3 Opus
```

**New:**
```
Anthropic [✓ Authenticated]
  Claude 3.5 Sonnet
  Claude 3 Opus

OpenAI [🔒 Not Authenticated - Press 'a' to auth]
  GPT-4 Turbo                    [locked]
  GPT-3.5 Turbo                  [locked]

Google [⚠ Degraded]
  Gemini Pro
```

### Auth Status Indicators

| Status | Icon | Color | Description |
|--------|------|-------|-------------|
| Authenticated | ✓ | Green | Provider ready to use |
| Not Authenticated | 🔒 | Yellow | Needs API key |
| Degraded | ⚠ | Orange | Circuit breaker half-open |
| Down | ✗ | Red | Circuit breaker open |
| Unknown | ? | Gray | Auth check failed |

### Locked Model Items

When provider not authenticated:
```
  GPT-4 Turbo                    [locked]
  ^                              ^
  Model name (grayed out)        Lock indicator
```

**Behavior:**
- Cannot select locked models
- Pressing Enter shows auth prompt
- Pressing 'a' on header starts auth flow

---

## Authentication Flow

### Flow 1: Direct Authentication

```
1. User opens model dialog
2. Sees "OpenAI [🔒 Not Authenticated]"
3. Presses 'a' on OpenAI header OR selects locked model
4. Dialog shows: "Enter OpenAI API Key:"
5. User enters key
6. Bridge validates key
7. Success → Models unlocked, toast shown
8. Failure → Error message, try again
```

### Flow 2: Auto-Detection

```
1. User opens model dialog
2. Sees "Anthropic [🔒 Not Authenticated]"
3. Presses 'd' for auto-detect
4. Bridge scans environment/files
5. Found → Auto-auth, toast shown
6. Not found → Prompt for manual entry
```

### Flow 3: Provider Selection

```
1. User navigates to locked model
2. Presses Enter
3. Prompt: "GPT-4 requires OpenAI authentication. Enter API key or press 'd' to auto-detect:"
4. User chooses auth method
5. After auth → Model automatically selected
```

---

## Implementation Plan

### Phase 1: Auth Status Display (Easy)

**Files to Modify:**
- `packages/tui/internal/components/dialog/models.go`

**Changes:**
1. Add `providerAuthStatus` map to `modelDialog`
2. Create `checkProviderAuth()` method
3. Update `buildGroupedResults()` to show auth status
4. Modify header rendering to include indicators

**Estimated Time:** 30 minutes

### Phase 2: Auth Status Checking (Medium)

**Files to Modify:**
- `packages/tui/internal/components/dialog/models.go`

**Changes:**
1. Call `app.AuthBridge.CheckAuthStatus()` for each provider
2. Cache results in `providerAuthStatus`
3. Refresh on dialog open (background)
4. Handle timeouts gracefully

**Estimated Time:** 20 minutes

### Phase 3: Locked Model Display (Easy)

**Files to Modify:**
- `packages/tui/internal/components/dialog/models.go`

**Changes:**
1. Update `modelItem.Render()` to check auth status
2. Gray out locked models
3. Add "[locked]" suffix
4. Make unselectable (return false from Selectable())

**Estimated Time:** 15 minutes

### Phase 4: Inline Auth Prompt (Hard)

**Files to Create:**
- `packages/tui/internal/components/dialog/auth_prompt.go`

**New Component:**
```go
type AuthPromptDialog struct {
    provider string
    input    textinput.Model
    error    string
}

func (a *AuthPromptDialog) View() string {
    // Render: "Enter {Provider} API Key:"
    // Show input box
    // Show error if any
    // Show hints: "Press 'd' for auto-detect"
}
```

**Estimated Time:** 40 minutes

### Phase 5: Auth Integration (Medium)

**Files to Modify:**
- `packages/tui/internal/components/dialog/models.go`

**Changes:**
1. Add auth prompt state to `modelDialog`
2. Handle 'a' key to start auth
3. Handle 'd' key for auto-detect
4. Call `app.AuthBridge.Authenticate()`
5. Show success/error messages
6. Refresh model list after auth

**Estimated Time:** 30 minutes

### Phase 6: Provider Health (Easy)

**Files to Modify:**
- `packages/tui/internal/components/dialog/models.go`

**Changes:**
1. Call `app.AuthBridge.GetProviderHealth()`
2. Show health status in headers
3. Update indicators (⚠, ✗)
4. Add tooltip/help text

**Estimated Time:** 20 minutes

---

## Code Examples

### Auth Status Checking

```go
type modelDialog struct {
    app               *app.App
    allModels         []ModelWithProvider
    providerAuthStatus map[string]AuthStatus // NEW
    // ... existing fields
}

type AuthStatus struct {
    IsAuthenticated bool
    Health          string // "healthy", "degraded", "down"
    ModelsCount     int
}

func (m *modelDialog) checkProviderAuth(providerID string) AuthStatus {
    if cached, ok := m.providerAuthStatus[providerID]; ok {
        return cached
    }

    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()

    status, err := m.app.AuthBridge.CheckAuthStatus(ctx, providerID)
    if err != nil {
        return AuthStatus{IsAuthenticated: false, Health: "unknown"}
    }

    result := AuthStatus{
        IsAuthenticated: status.IsAuthenticated,
        ModelsCount:     status.ModelsCount,
    }

    // Check health
    health, err := m.app.AuthBridge.GetProviderHealth(ctx, providerID)
    if err == nil {
        result.Health = health.Status
    }

    m.providerAuthStatus[providerID] = result
    return result
}
```

### Provider Header with Auth Status

```go
func (m *modelDialog) buildGroupedResults() []list.Item {
    var items []list.Item

    // ... recent section

    for _, providerName := range providerNames {
        models := providerGroups[providerName]
        providerID := models[0].Provider.ID

        // Check auth status
        authStatus := m.checkProviderAuth(providerID)

        // Create header with status
        header := m.buildProviderHeader(providerName, authStatus)
        items = append(items, header)

        // Add models
        for _, model := range models {
            item := modelItem{
                model:           model,
                isAuthenticated: authStatus.IsAuthenticated,
            }
            items = append(items, item)
        }
    }

    return items
}

func (m *modelDialog) buildProviderHeader(name string, status AuthStatus) list.Item {
    var indicator string
    if status.IsAuthenticated {
        if status.Health == "healthy" {
            indicator = " [✓ Authenticated]"
        } else if status.Health == "degraded" {
            indicator = " [⚠ Degraded]"
        } else if status.Health == "down" {
            indicator = " [✗ Down]"
        }
    } else {
        indicator = " [🔒 Not Authenticated - Press 'a' to auth]"
    }

    return list.HeaderItem(name + indicator)
}
```

### Locked Model Rendering

```go
type modelItem struct {
    model           ModelWithProvider
    isAuthenticated bool // NEW
}

func (m modelItem) Render(selected bool, width int, baseStyle styles.Style) string {
    t := theme.CurrentTheme()

    itemStyle := baseStyle.
        Background(t.BackgroundPanel()).
        Foreground(t.Text())

    // Gray out if locked
    if !m.isAuthenticated {
        itemStyle = itemStyle.Foreground(t.TextMuted()).Faint(true)
    } else if selected {
        itemStyle = itemStyle.Foreground(t.Primary())
    }

    providerStyle := baseStyle.
        Foreground(t.TextMuted()).
        Background(t.BackgroundPanel())

    modelPart := itemStyle.Render(m.model.Model.Name)
    providerPart := providerStyle.Render(fmt.Sprintf(" %s", m.model.Provider.Name))

    // Add lock indicator
    lockPart := ""
    if !m.isAuthenticated {
        lockPart = providerStyle.Render(" [locked]")
    }

    combinedText := modelPart + providerPart + lockPart
    return baseStyle.
        Background(t.BackgroundPanel()).
        PaddingLeft(1).
        Render(combinedText)
}

func (m modelItem) Selectable() bool {
    return m.isAuthenticated // Only selectable if authenticated
}
```

### Auth Prompt Dialog

```go
package dialog

import (
    "github.com/charmbracelet/bubbles/v2/textinput"
    tea "github.com/charmbracelet/bubbletea/v2"
)

type AuthPromptDialog struct {
    provider  string
    input     textinput.Model
    error     string
    autoDetectHint bool
}

func NewAuthPromptDialog(provider string) *AuthPromptDialog {
    ti := textinput.New()
    ti.Placeholder = "sk-..."
    ti.Focus()
    ti.CharLimit = 256
    ti.Width = 60
    ti.EchoMode = textinput.EchoPassword // Hide API key

    return &AuthPromptDialog{
        provider:       provider,
        input:          ti,
        autoDetectHint: true,
    }
}

func (a *AuthPromptDialog) Update(msg tea.Msg) (*AuthPromptDialog, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyPressMsg:
        switch msg.String() {
        case "enter":
            // Submit API key
            return a, util.CmdHandler(AuthSubmitMsg{
                Provider: a.provider,
                APIKey:   a.input.Value(),
            })
        case "esc":
            return a, util.CmdHandler(AuthCancelMsg{})
        case "ctrl+d":
            // Auto-detect
            return a, util.CmdHandler(AuthAutoDetectMsg{
                Provider: a.provider,
            })
        }
    }

    var cmd tea.Cmd
    a.input, cmd = a.input.Update(msg)
    return a, cmd
}

func (a *AuthPromptDialog) View() string {
    t := theme.CurrentTheme()

    titleStyle := styles.NewStyle().
        Foreground(t.Primary()).
        Bold(true)

    title := titleStyle.Render(fmt.Sprintf("Enter %s API Key:", a.provider))

    inputView := a.input.View()

    hintStyle := styles.NewStyle().
        Foreground(t.TextMuted()).
        Faint(true)

    hints := ""
    if a.autoDetectHint {
        hints = hintStyle.Render("Press Ctrl+D for auto-detect | Esc to cancel")
    }

    errorView := ""
    if a.error != "" {
        errorStyle := styles.NewStyle().
            Foreground(t.Error())
        errorView = "\n" + errorStyle.Render("Error: " + a.error)
    }

    return lipgloss.JoinVertical(
        lipgloss.Left,
        title,
        "",
        inputView,
        "",
        hints,
        errorView,
    )
}
```

---

## Message Types

```go
// Auth-related messages
type AuthSubmitMsg struct {
    Provider string
    APIKey   string
}

type AuthCancelMsg struct{}

type AuthAutoDetectMsg struct {
    Provider string
}

type AuthSuccessMsg struct {
    Provider string
    ModelsCount int
}

type AuthFailureMsg struct {
    Provider string
    Error    string
}

type AuthStatusUpdatedMsg struct {
    Provider string
    Status   AuthStatus
}
```

---

## Keyboard Shortcuts

| Key | Context | Action |
|-----|---------|--------|
| `a` | On provider header | Start authentication |
| `d` | On provider header | Auto-detect credentials |
| `Enter` | On locked model | Show auth prompt |
| `Ctrl+D` | In auth prompt | Auto-detect |
| `Enter` | In auth prompt | Submit API key |
| `Esc` | In auth prompt | Cancel |

---

## Error Handling

### Network Errors
```
"Failed to check authentication status. Please check your connection."
```

### Invalid API Key
```
"Invalid API key for OpenAI. Please check and try again."
```

### Auto-Detect Failure
```
"No credentials found. Please enter your API key manually."
```

### Provider Down
```
"OpenAI is currently unavailable (circuit breaker open). Try again later."
```

---

## Performance Considerations

### Auth Checks
- **Cache Duration:** 30 seconds (refresh on dialog open)
- **Timeout:** 1 second per provider
- **Parallel Checks:** Check all providers concurrently
- **Background Updates:** Don't block UI during checks

### UI Rendering
- **Lazy Loading:** Only check auth for visible providers
- **Debouncing:** Don't spam auth checks on rapid opens
- **Caching:** Store results in dialog struct

---

## Testing Strategy

### Unit Tests

1. **Auth Status Checking**
   ```go
   func TestCheckProviderAuth(t *testing.T) {
       // Test authenticated provider
       // Test unauthenticated provider
       // Test timeout handling
   }
   ```

2. **Locked Model Rendering**
   ```go
   func TestLockedModelRender(t *testing.T) {
       // Test locked model appears grayed
       // Test locked model shows indicator
       // Test locked model unselectable
   }
   ```

3. **Auth Prompt**
   ```go
   func TestAuthPrompt(t *testing.T) {
       // Test API key input
       // Test auto-detect trigger
       // Test validation
   }
   ```

### Integration Tests

1. **End-to-End Auth Flow**
   - Open dialog → See locked models
   - Press 'a' → Enter key
   - Verify → Models unlocked

2. **Auto-Detect Flow**
   - Press 'd' → Credentials found
   - Toast shown → Models unlocked

3. **Error Handling**
   - Invalid key → Error shown
   - Retry → Success

---

## Success Criteria

- [ ] Provider auth status visible in headers
- [ ] Locked models clearly indicated
- [ ] Auth prompt appears on 'a' press
- [ ] API key authentication works
- [ ] Auto-detect finds credentials
- [ ] Toast notifications on success/failure
- [ ] Models unlock after successful auth
- [ ] Provider health indicators show
- [ ] No UI freezing during auth checks
- [ ] Graceful error handling

---

## Timeline

| Phase | Task | Time | Status |
|-------|------|------|--------|
| 1 | Auth status display | 30min | 🔴 TODO |
| 2 | Auth status checking | 20min | 🔴 TODO |
| 3 | Locked model display | 15min | 🔴 TODO |
| 4 | Auth prompt component | 40min | 🔴 TODO |
| 5 | Auth integration | 30min | 🔴 TODO |
| 6 | Provider health | 20min | 🔴 TODO |
| 7 | Testing | 30min | 🔴 TODO |

**Total Estimated Time:** ~3 hours

---

## Next Actions

1. **Review Design**: Get feedback on UI/UX
2. **Start Implementation**: Begin with Phase 1 (auth status display)
3. **Incremental Testing**: Test each phase before proceeding
4. **Documentation**: Update user docs with new features

---

**Status:** ✅ Design Complete, Ready for Implementation
**Priority:** High - Critical for Phase 2 completion
**Dependencies:** AuthBridge (✅ Complete)

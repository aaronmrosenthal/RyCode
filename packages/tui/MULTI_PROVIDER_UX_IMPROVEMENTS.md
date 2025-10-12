# Multi-Provider UX Improvements

## User Story

As a user authenticated with multiple providers (Codex, Gemini, Qwen, Claude), I want RyCode to:
1. **Auto-detect** my existing auth on startup (already working!)
2. **Never bother me** with auth prompts if already authenticated
3. **Easily switch** between authenticated providers/models using Tab key

## Current State Analysis

### ✅ What's Already Working

1. **Auto-Detection on First Run** (`app.go:490-494`)
   - Checks if it's first run via `isFirstRun()`
   - Runs `autoDetectAllCredentials()` automatically
   - Shows success toast: "Found N provider(s). Ready to code!"

2. **Background Authentication** (`models.go:295-324`)
   - `tryAutoAuthThenPrompt()` attempts auto-auth before showing prompt
   - 3-second timeout, graceful fallback
   - Only prompts if auto-detect fails

3. **Provider Status Indicators** (`models.go:461-479`)
   - Headers show: "Provider ✓" (authenticated) or "Provider 🔒" (locked)
   - Health indicators: ✓ (healthy), ⚠ (degraded), ✗ (down)

4. **Recent Models Cycling** (`app.go:321-376`)
   - Tab already cycles through recent models
   - Shows toast: "Switched to Model (Provider)"

### ⚠️ Current Gaps

1. **Auto-detection only runs on "first run"**
   - If you've used RyCode before but add new providers, auto-detect doesn't run
   - You have to manually press 'd' in model dialog to trigger it

2. **Tab cycles recent models, not authenticated providers**
   - Tab only works if you've USED models before
   - Doesn't help you discover newly authenticated providers

3. **No visual "all authenticated" status**
   - User doesn't know which providers are ready without opening /model
   - Status bar could show: "4 providers ready"

4. **Manual step required after new auth**
   - After `rycode auth login`, must restart or manually refresh

## Proposed Improvements

### 1. Auto-Detect on EVERY Startup (Not Just First Run)

**Goal**: Always detect newly added providers automatically

**Change**: `app.go:490-494`
```go
// BEFORE (only first run):
if a.isFirstRun() {
    autoDetectCmd = a.autoDetectAllCredentials()
}

// AFTER (every startup, but silent if none found):
autoDetectCmd = a.autoDetectAllCredentialsQuiet()
```

**New Function**:
```go
// autoDetectAllCredentialsQuiet runs auto-detect silently (no toast if found=0)
func (a *App) autoDetectAllCredentialsQuiet() tea.Cmd {
    return func() tea.Msg {
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()

        result, err := a.AuthBridge.AutoDetect(ctx)
        if err != nil {
            slog.Debug("Auto-detect failed", "error", err)
            return nil
        }

        if result.Found > 0 {
            slog.Info("Auto-detected credentials", "count", result.Found)
            // Only show toast if NEW providers found (not already authenticated)
            return AuthStatusRefreshMsg{}
        }

        return nil
    }
}
```

**Benefits**:
- ✅ Detects new providers automatically
- ✅ Silent if nothing new
- ✅ No user interruption

### 2. Tab Cycles Through Authenticated Providers (Not Just Recent)

**Goal**: Use Tab to discover and switch between all authenticated providers

**Current**: Tab calls `CycleRecentModel()` which only works if models are in recent history

**Improvement**: Create new command `CycleAuthenticatedProviders`

**New Function** in `app.go`:
```go
// CycleAuthenticatedProviders cycles through all authenticated providers
func (a *App) CycleAuthenticatedProviders(forward bool) (*App, tea.Cmd) {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    // Get authentication status for all providers
    status, err := a.AuthBridge.GetAuthStatus(ctx)
    if err != nil {
        return a, toast.NewErrorToast("Failed to get provider status")
    }

    if len(status.Authenticated) == 0 {
        return a, toast.NewInfoToast("No authenticated providers. Press 'd' to auto-detect.")
    }

    if len(status.Authenticated) == 1 {
        return a, toast.NewInfoToast("Only one provider authenticated")
    }

    // Find current provider index
    currentIndex := -1
    for i, prov := range status.Authenticated {
        if a.Provider != nil && prov.ID == a.Provider.ID {
            currentIndex = i
            break
        }
    }

    // Calculate next index
    nextIndex := 0
    if currentIndex != -1 {
        if forward {
            nextIndex = (currentIndex + 1) % len(status.Authenticated)
        } else {
            nextIndex = (currentIndex - 1 + len(status.Authenticated)) % len(status.Authenticated)
        }
    }

    // Get next provider's default model
    nextProvider := status.Authenticated[nextIndex]

    // Find provider and default model in a.Providers
    provider, model := a.findProviderDefaultModel(nextProvider.ID)
    if provider == nil || model == nil {
        return a, toast.NewErrorToast("Provider or model not found")
    }

    a.Provider = provider
    a.Model = model
    a.State.AgentModel[a.Agent().Name] = AgentModel{
        ProviderID: provider.ID,
        ModelID:    model.ID,
    }
    a.State.UpdateModelUsage(provider.ID, model.ID)

    return a, tea.Sequence(
        a.SaveState(),
        toast.NewSuccessToast(
            fmt.Sprintf("→ %s: %s", provider.Name, model.Name),
        ),
    )
}
```

**Keybinding** (add to commands):
```go
commands.ModelCycleAuthenticatedCommand: {
    Name:        "cycle_authenticated_providers",
    Description: "Cycle authenticated providers",
    Keybindings: []Keybinding{
        {Key: "tab", RequiresLeader: false},
    },
}
```

**Benefits**:
- ✅ Tab works even if you haven't used models yet
- ✅ Discover all authenticated providers
- ✅ Fast switching: Tab → Codex, Tab → Gemini, Tab → Claude

### 3. Status Bar Shows Authenticated Provider Count

**Goal**: Immediately see how many providers are ready

**Change**: `status.go` (status component)

```go
// Add to status bar left side:
"4 providers ✓" // if all authenticated
"2/4 providers"  // if some authenticated
```

**Implementation**:
```go
func (s *StatusComponent) renderProviderStatus() string {
    ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
    defer cancel()

    status, err := s.app.AuthBridge.GetAuthStatus(ctx)
    if err != nil {
        return ""
    }

    total := len(s.app.Providers)
    authed := len(status.Authenticated)

    if authed == total {
        return fmt.Sprintf("%d providers ✓", authed)
    } else if authed > 0 {
        return fmt.Sprintf("%d/%d providers", authed, total)
    } else {
        return "No providers ✓"
    }
}
```

**Benefits**:
- ✅ Immediate visibility
- ✅ No need to open /model to check
- ✅ Motivates authentication if 0

### 4. Proactive Refresh After `rycode auth login`

**Goal**: Auto-refresh providers after CLI auth, no restart needed

**Currently**: User must restart RyCode after `rycode auth login`

**Improvement**: Add file watcher for `~/.local/share/rycode/auth.json`

**Implementation**:
```go
// In app initialization:
func (a *App) WatchAuthFile() tea.Cmd {
    return func() tea.Msg {
        watcher, err := fsnotify.NewWatcher()
        if err != nil {
            return nil
        }

        authPath := filepath.Join(os.Getenv("HOME"), ".local/share/rycode/auth.json")
        watcher.Add(authPath)

        go func() {
            for {
                select {
                case event := <-watcher.Events:
                    if event.Op&fsnotify.Write == fsnotify.Write {
                        // Auth file changed!
                        return AuthFileChangedMsg{}
                    }
                }
            }
        }()

        return nil
    }
}

// Handle message:
case AuthFileChangedMsg:
    return a, tea.Batch(
        a.autoDetectAllCredentialsQuiet(),
        toast.NewInfoToast("Auth updated. Refreshing providers..."),
    )
```

**Benefits**:
- ✅ No restart needed after auth
- ✅ Seamless workflow
- ✅ Instant feedback

### 5. Startup Toast: "All Providers Ready"

**Goal**: Immediately confirm all your providers are detected

**Change**: Improve the startup toast to be more informative

```go
// BEFORE:
"Found 3 provider(s). Ready to code!"

// AFTER:
"All providers ready: Codex, Gemini, Claude ✓"
// or if not all:
"Ready: Codex, Gemini (2/4 providers)"
```

**Benefits**:
- ✅ Clear visibility
- ✅ Confirms your setup
- ✅ Shows which providers are missing

## Implementation Priority

### Phase 1: Quick Wins (1-2 hours)
1. ✅ **Auto-detect on every startup** (not just first run)
2. ✅ **Improved startup toast** with provider names
3. ✅ **Status bar provider count**

### Phase 2: Tab Enhancement (2-3 hours)
4. ✅ **Tab cycles authenticated providers**
5. ✅ **Keybinding update**

### Phase 3: Advanced (optional)
6. ✅ **Auth file watcher** for instant refresh
7. ✅ **Model dialog shows "last used" per provider**

## User Workflow (After Implementation)

### Scenario: Fresh Start
```
1. Open RyCode → "All providers ready: Codex, Gemini, Qwen, Claude ✓"
2. Start typing → Uses last model (e.g., Claude Sonnet)
3. Press Tab → Switches to Gemini
4. Press Tab → Switches to Codex
5. Press Tab → Switches to Qwen
6. Press Tab → Back to Claude
```

### Scenario: New Provider Auth
```
1. Terminal: rycode auth login → Add DeepSeek
2. RyCode (still open) → "Auth updated. Refreshing providers..."
3. RyCode → "Ready: Codex, Gemini, Qwen, Claude, DeepSeek ✓"
4. Press Tab repeatedly → DeepSeek appears in rotation
```

### Scenario: Check Status
```
1. Look at status bar → "5 providers ✓"
2. No need to open /model unless choosing specific model
```

## Testing Checklist

- [ ] Fresh install: Auto-detect runs on first startup
- [ ] Existing install: Auto-detect runs silently on every startup
- [ ] New auth: Providers refresh without restart
- [ ] Tab key: Cycles through authenticated providers only
- [ ] Status bar: Shows correct provider count
- [ ] Toast messages: Clear and informative
- [ ] No providers: Graceful fallback messages

## Code Files to Modify

1. `packages/tui/internal/app/app.go` - Auto-detect logic, Tab cycling
2. `packages/tui/internal/components/status/status.go` - Provider count display
3. `packages/tui/internal/commands/commands.go` - New Tab command
4. `packages/tui/internal/tui/tui.go` - Handle Tab keybinding
5. `packages/tui/internal/auth/bridge.go` - GetAuthStatus improvements

## Expected UX Improvements

**Before**:
- Open RyCode → No indication of auth status
- Tab → Only works if used models before
- New auth → Must restart RyCode
- Check providers → Must open /model dialog

**After**:
- Open RyCode → "All providers ready: Codex, Gemini, Qwen, Claude ✓"
- Tab → Instantly switch: Claude → Gemini → Codex → Qwen
- New auth → Auto-refreshes, no restart
- Check providers → Status bar shows "4 providers ✓"

## Summary

This improvement plan focuses on **frictionless multi-provider workflow**:
1. ✅ Auto-detect runs on every startup (silent if nothing new)
2. ✅ Tab cycles through authenticated providers (not just recent)
3. ✅ Status bar shows provider count
4. ✅ Auth file watcher for instant refresh
5. ✅ Improved startup toast with provider names

**Result**: You never need to think about auth. Open RyCode, press Tab to switch providers, start coding.

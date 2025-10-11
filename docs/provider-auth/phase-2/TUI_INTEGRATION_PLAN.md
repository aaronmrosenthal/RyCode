# TUI Integration Plan - Phase 2

## 📋 Current Architecture Analysis

### Status Bar (packages/tui/internal/components/status/status.go)
**Current Implementation:**
```go
// Lines 144-178 in status.go
agentColor := util.GetAgentColor(m.app.AgentIndex)
agent := agentNameStyle(strings.ToUpper(m.app.Agent().Name)) + agentDescStyle(" AGENT")
agent = agentStyle.Padding(0, 1).BorderLeft(true).Render(agent)
agent = faintStyle.Render(key+" ") + agent
```

**What it shows now:**
```
[RyCode v1.0] [~/project:main]                    [tab BUILD AGENT]
```

**What we need:**
```
[RyCode v1.0] [~/project:main]      [Claude 3.5 Sonnet | 💰 $0.12 | tab→]
```

### Model Selector (packages/tui/internal/components/dialog/models.go)
**Current Implementation:**
- Groups models by provider (lines 350-390)
- Shows "Recent" section at top
- Uses fuzzy search filtering
- Displays model items with provider badge

**What we need to add:**
- Authentication status per provider
- Inline "Sign In" action items for unauthenticated providers
- Provider health indicators (circuit breaker status)
- Call TypeScript auth functions from Go

### App Structure (packages/tui/internal/app/app.go)
**Current State:**
```go
type App struct {
    Agents     []opencode.Agent    // Array of agents (build, plan, doc)
    AgentIndex int                 // Current agent index
    Provider   *opencode.Provider  // Current provider
    Model      *opencode.Model     // Current model
}
```

**What needs to change:**
- Remove dependency on `AgentIndex` for status display
- Make `Provider` and `Model` the primary display source
- Add cost tracking state
- Add provider health tracking

---

## 🎯 Integration Tasks

### Task 1: Update Status Bar Display
**File:** `packages/tui/internal/components/status/status.go`

**Changes Required:**
1. **Replace agent display with model display** (lines 144-178)
   ```go
   // OLD:
   agent := agentNameStyle(strings.ToUpper(m.app.Agent().Name)) + agentDescStyle(" AGENT")

   // NEW:
   modelDisplay := m.buildModelDisplay()
   // Shows: "Claude 3.5 Sonnet | 💰 $0.12 | tab→"
   ```

2. **Add cost tracking to App struct**
   ```go
   type App struct {
       // ... existing fields
       CurrentCost float64  // Today's cost
       CostTracker *CostTracker  // Bridge to TypeScript cost tracker
   }
   ```

3. **Create buildModelDisplay() method**
   ```go
   func (m *statusComponent) buildModelDisplay() string {
       if m.app.Model == nil || m.app.Provider == nil {
           return faintStyle.Render("No model selected")
       }

       modelName := m.app.Model.Name
       cost := fmt.Sprintf("💰 $%.2f", m.app.CurrentCost)
       hint := "[tab→]"

       return modelNameStyle(modelName) +
              faintStyle(" | ") +
              costStyle(cost) +
              faintStyle(" | ") +
              hintStyle(hint)
   }
   ```

**Impact:**
- User sees current model instead of agent
- Real-time cost visibility
- Tab key hint for model switching

---

### Task 2: Add Inline Authentication to Model Selector
**File:** `packages/tui/internal/components/dialog/models.go`

**Changes Required:**
1. **Add authentication status to ModelWithProvider**
   ```go
   type ModelWithProvider struct {
       Model           opencode.Model
       Provider        opencode.Provider
       IsAuthenticated bool      // NEW
       HealthStatus    string    // NEW: "healthy", "degraded", "down"
   }
   ```

2. **Modify buildGroupedResults() to show auth status** (lines 350-390)
   ```go
   func (m *modelDialog) buildGroupedResults() []list.Item {
       var items []list.Item

       // Group by provider
       for providerID, provider := range providerGroups {
           // Check auth status
           authStatus := m.checkAuthStatus(providerID)

           if !authStatus.IsAuthenticated {
               // Add "Sign In" action item
               items = append(items, authActionItem{
                   provider: provider,
                   action: "sign_in"
               })
           } else {
               // Show models
               for _, model := range provider.Models {
                   items = append(items, modelItem{model: model})
               }
           }
       }
   }
   ```

3. **Create auth action item type**
   ```go
   type authActionItem struct {
       provider opencode.Provider
       action   string  // "sign_in", "reconnect", "refresh"
   }

   func (a authActionItem) Render(w io.Writer, m list.Model, index int, item list.Item) {
       // Render as: "→ Sign in to Anthropic" with auth icon
       icon := "🔐"
       text := fmt.Sprintf("%s Sign in to %s", icon, a.provider.Name)
       // Style and write
   }
   ```

4. **Handle selection of auth action item**
   ```go
   func (m *modelDialog) handleAuthAction(provider opencode.Provider) tea.Cmd {
       // Call TypeScript auth system
       return func() tea.Msg {
           result := callAuthManager(provider.ID)
           if result.Success {
               return AuthSuccessMsg{Provider: provider}
           }
           return AuthFailureMsg{Error: result.Error}
       }
   }
   ```

**Impact:**
- Users can authenticate directly from model selector
- No need to leave TUI to configure auth
- Clear visual feedback on auth status

---

### Task 3: Implement Tab Key Model Cycling
**File:** `packages/tui/internal/app/app.go`

**Changes Required:**
1. **Modify cycleMode() to cycle models instead of agents** (lines 256-299)
   ```go
   // REPLACE: cycleMode() function
   func (a *App) cycleAuthenticatedModels(forward bool) (*App, tea.Cmd) {
       // Get list of authenticated providers
       authenticatedModels := a.getAuthenticatedModels()

       if len(authenticatedModels) == 0 {
           return a, toast.NewErrorToast("No authenticated providers")
       }

       // Find current model index
       currentIndex := -1
       for i, model := range authenticatedModels {
           if a.Model != nil && model.ID == a.Model.ID {
               currentIndex = i
               break
           }
       }

       // Cycle to next/previous
       var nextIndex int
       if forward {
           nextIndex = (currentIndex + 1) % len(authenticatedModels)
       } else {
           nextIndex = (currentIndex - 1 + len(authenticatedModels)) % len(authenticatedModels)
       }

       // Update model
       nextModel := authenticatedModels[nextIndex]
       a.Provider = nextModel.Provider
       a.Model = nextModel.Model

       return a, tea.Sequence(
           a.SaveState(),
           toast.NewSuccessToast(fmt.Sprintf("Switched to %s", nextModel.Model.Name)),
           a.UpdateCostDisplay(),  // NEW: Update cost in status bar
       )
   }
   ```

2. **Add helper to get authenticated models**
   ```go
   func (a *App) getAuthenticatedModels() []AuthenticatedModel {
       var models []AuthenticatedModel

       for _, provider := range a.Providers {
           // Call TypeScript auth system to check authentication
           authStatus := checkProviderAuth(provider.ID)

           if authStatus.IsAuthenticated {
               for _, model := range provider.Models {
                   models = append(models, AuthenticatedModel{
                       Provider: &provider,
                       Model:    &model,
                   })
               }
           }
       }

       return models
   }
   ```

3. **Update command binding** (keep existing Tab key binding)
   ```go
   // In commands package - map Tab to cycleAuthenticatedModels instead of cycleMode
   ```

**Impact:**
- Tab key cycles through all authenticated models
- No more agent concept in UX
- Seamless model switching

---

### Task 4: Go-TypeScript Bridge for Auth
**New File:** `packages/tui/internal/auth/bridge.go`

**Implementation:**
```go
package auth

import (
    "context"
    "encoding/json"
    "os/exec"
)

// AuthStatus represents the authentication status from TypeScript
type AuthStatus struct {
    IsAuthenticated bool   `json:"isAuthenticated"`
    Provider        string `json:"provider"`
    ModelsCount     int    `json:"modelsCount"`
    Error           string `json:"error,omitempty"`
}

// ProviderHealth represents circuit breaker health
type ProviderHealth struct {
    Provider      string `json:"provider"`
    Status        string `json:"status"` // "healthy", "degraded", "down"
    FailureCount  int    `json:"failureCount"`
    NextAttemptAt string `json:"nextAttemptAt,omitempty"`
}

// CostSummary represents cost tracking data
type CostSummary struct {
    TodayCost       float64 `json:"todayCost"`
    MonthCost       float64 `json:"monthCost"`
    Projection      float64 `json:"projection"`
    SavingsTip      string  `json:"savingsTip,omitempty"`
}

// CheckAuthStatus calls TypeScript authManager.isAuthenticated()
func CheckAuthStatus(ctx context.Context, provider string) (*AuthStatus, error) {
    cmd := exec.CommandContext(ctx, "bun", "run",
        "./packages/rycode/src/auth/cli.ts",
        "check", provider)

    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }

    var status AuthStatus
    if err := json.Unmarshal(output, &status); err != nil {
        return nil, err
    }

    return &status, nil
}

// Authenticate calls TypeScript authManager.authenticate()
func Authenticate(ctx context.Context, provider, apiKey string) error {
    cmd := exec.CommandContext(ctx, "bun", "run",
        "./packages/rycode/src/auth/cli.ts",
        "auth", provider, apiKey)

    _, err := cmd.Output()
    return err
}

// GetCostSummary calls TypeScript authManager.getCostSummary()
func GetCostSummary(ctx context.Context) (*CostSummary, error) {
    cmd := exec.CommandContext(ctx, "bun", "run",
        "./packages/rycode/src/auth/cli.ts",
        "cost")

    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }

    var summary CostSummary
    if err := json.Unmarshal(output, &summary); err != nil {
        return nil, err
    }

    return &summary, nil
}

// GetProviderHealth calls TypeScript circuitBreakerRegistry.getStats()
func GetProviderHealth(ctx context.Context, provider string) (*ProviderHealth, error) {
    cmd := exec.CommandContext(ctx, "bun", "run",
        "./packages/rycode/src/auth/cli.ts",
        "health", provider)

    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }

    var health ProviderHealth
    if err := json.Unmarshal(output, &health); err != nil {
        return nil, err
    }

    return &health, nil
}
```

**TypeScript CLI File:** `packages/rycode/src/auth/cli.ts`
```typescript
#!/usr/bin/env bun

import { authManager } from './auth-manager'

const command = process.argv[2]
const args = process.argv.slice(3)

async function main() {
  switch (command) {
    case 'check': {
      const provider = args[0]
      const isAuthenticated = authManager.isAuthenticated(provider)
      const models = isAuthenticated ? authManager.getAuthenticatedModels(provider) : []
      console.log(JSON.stringify({
        isAuthenticated,
        provider,
        modelsCount: models.length
      }))
      break
    }

    case 'auth': {
      const provider = args[0]
      const apiKey = args[1]
      await authManager.authenticate({ provider, apiKey })
      console.log(JSON.stringify({ success: true }))
      break
    }

    case 'cost': {
      const summary = authManager.getCostSummary()
      console.log(JSON.stringify(summary))
      break
    }

    case 'health': {
      const provider = args[0]
      const health = authManager.getProviderHealth(provider)
      console.log(JSON.stringify(health))
      break
    }

    default:
      console.error('Unknown command:', command)
      process.exit(1)
  }
}

main().catch(err => {
  console.error(err)
  process.exit(1)
})
```

**Impact:**
- Go can call TypeScript auth functions
- JSON-based communication
- No need to reimplement auth logic in Go

---

### Task 5: Provider Health Indicators
**File:** `packages/tui/internal/components/dialog/models.go`

**Changes Required:**
1. **Add health badge to provider headers**
   ```go
   func (m *modelDialog) renderProviderHeader(provider opencode.Provider, health *auth.ProviderHealth) string {
       name := provider.Name

       var healthBadge string
       switch health.Status {
       case "healthy":
           healthBadge = "✓"  // Green checkmark
       case "degraded":
           healthBadge = "⚠"  // Yellow warning
       case "down":
           healthBadge = "✗"  // Red X
       }

       return fmt.Sprintf("%s %s", healthBadge, name)
   }
   ```

2. **Show circuit breaker info on hover/selection**
   ```go
   func (m *modelDialog) renderProviderDetails(health *auth.ProviderHealth) string {
       if health.Status == "down" {
           return fmt.Sprintf(
               "Provider temporarily unavailable. Retry in %s",
               health.NextAttemptAt,
           )
       }
       return ""
   }
   ```

**Impact:**
- Users see provider health at a glance
- Know when providers are temporarily unavailable
- Understand circuit breaker behavior

---

## 🔄 Migration Strategy

### Backward Compatibility
**Keep agents in codebase for now:**
```go
// In app.go - dual mode support
if os.Getenv("ENABLE_PROVIDER_AUTH") == "true" {
    // Use new provider-based system
    status = m.buildModelDisplay()
} else {
    // Use legacy agent system
    status = m.buildAgentDisplay()
}
```

### Gradual Rollout
1. **Week 1:** Feature flag at 0% (infrastructure only)
2. **Week 2:** Enable for 10% of users
3. **Week 3:** Expand to 50% of users
4. **Week 4:** Full rollout at 100%

---

## 📊 Success Metrics

### Must Track:
- [ ] Auth success rate per provider
- [ ] Tab key usage frequency
- [ ] Cost display accuracy
- [ ] Model switch latency
- [ ] Circuit breaker events
- [ ] User satisfaction (surveys)

### Target Metrics:
- Auth success rate: >95%
- Model switch latency: <500ms
- Cost calculation accuracy: 100%
- Circuit breaker auto-recovery: <2 minutes

---

## 🎯 File Changes Summary

### Files to Modify:
1. `packages/tui/internal/components/status/status.go` (status bar)
2. `packages/tui/internal/components/dialog/models.go` (model selector)
3. `packages/tui/internal/app/app.go` (Tab key cycling)
4. `packages/tui/internal/commands/commands.go` (command mappings)

### Files to Create:
1. `packages/tui/internal/auth/bridge.go` (Go-TypeScript bridge)
2. `packages/rycode/src/auth/cli.ts` (TypeScript CLI interface)

### Total Lines of Code: ~800 lines Go + ~100 lines TypeScript

---

## ⚠️ Known Challenges

### Challenge 1: Cost Tracking Latency
**Problem:** Fetching cost from TypeScript on every render could be slow

**Solution:**
- Cache cost data in App struct
- Update every 5 seconds via background goroutine
- Show stale indicator if data is >10 seconds old

### Challenge 2: Authentication State Sync
**Problem:** Go and TypeScript need to share auth state

**Solution:**
- TypeScript writes auth state to JSON file
- Go reads JSON file on demand (with caching)
- File watching for real-time updates

### Challenge 3: Error Handling
**Problem:** TypeScript errors need to surface in Go TUI

**Solution:**
- Structured error JSON format
- Map TypeScript errors to Go toast messages
- Show actionable error messages from TypeScript

---

## 🚀 Next Steps

1. ✅ **Phase 1 Complete:** Infrastructure (5,045 LOC)
2. 🔄 **Phase 2 In Progress:** TUI Integration (this document)
3. ⏳ **Phase 3 Pending:** Migration wizard
4. ⏳ **Phase 4 Pending:** Testing (90% coverage)
5. ⏳ **Phase 5 Pending:** Launch (gradual rollout)

---

**Ready to implement!** 🎉

This design maintains the existing TUI structure while seamlessly integrating the new provider authentication system.

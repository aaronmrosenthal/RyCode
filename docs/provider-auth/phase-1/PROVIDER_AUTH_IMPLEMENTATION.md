# Provider Authentication Implementation Guide

## Quick Start Implementation Path

### Step 1: Create Provider Auth Commands

Create new CLI commands for provider authentication that can be called from the model selector:

```typescript
// packages/rycode/src/cli/cmd/auth-provider.ts

export const AuthProviderCommand = cmd({
  command: "auth <provider>",
  describe: "authenticate with an AI provider",
  builder: (yargs) =>
    yargs
      .positional("provider", {
        describe: "provider to authenticate with",
        type: "string",
        choices: ["anthropic", "google", "openai", "qwen", "grok"],
      })
      .option("method", {
        describe: "authentication method",
        type: "string",
        choices: ["api-key", "oauth", "browser", "cli"],
      }),
  handler: async (opts) => {
    const provider = opts.provider as string
    const method = opts.method || getDefaultMethod(provider)

    switch (provider) {
      case "anthropic":
        await authenticateAnthropic(method)
        break
      case "google":
        await authenticateGoogle(method)
        break
      case "openai":
        await authenticateOpenAI(method)
        break
      case "qwen":
        await authenticateQwen(method)
        break
      case "grok":
        await authenticateGrok(method)
        break
    }
  },
})
```

### Step 2: Modify Model Selector Dialog

Update the models dialog to include authentication UI:

```go
// packages/tui/internal/components/dialog/models.go

type ModelDialog struct {
    app          *app.App
    providers    []ProviderWithAuth
    searchDialog *SearchDialog
}

type ProviderWithAuth struct {
    Provider     opencode.Provider
    IsAuthenticated bool
    AuthMethod   string
    UserInfo     string // email or username if available
}

func (m *modelDialog) buildProviderSections() []list.Item {
    var items []list.Item

    for _, provider := range m.providers {
        if provider.IsAuthenticated {
            // Show authenticated provider with models
            items = append(items, list.HeaderItem(
                fmt.Sprintf("✓ %s (%s)", provider.Provider.Name, provider.UserInfo),
            ))
            for _, model := range provider.Provider.Models {
                items = append(items, modelItem{
                    model: model,
                    provider: provider.Provider,
                })
            }
            items = append(items, actionItem{
                label: "Manage Account",
                action: m.manageProviderAuth(provider),
            })
        } else {
            // Show unauthenticated provider
            items = append(items, list.HeaderItem(
                fmt.Sprintf("🔐 %s", provider.Provider.Name),
            ))
            items = append(items, actionItem{
                label: "Sign In",
                action: m.authenticateProvider(provider),
            })
        }
    }

    return items
}
```

### Step 3: Create Auth Dialogs

Create inline authentication dialogs for each auth method:

```go
// packages/tui/internal/components/dialog/auth.go

type AuthDialog interface {
    layout.Modal
}

type apiKeyAuthDialog struct {
    provider   string
    textInput  textinput.Model
    modal      *modal.Modal
}

func (a *apiKeyAuthDialog) View() string {
    t := theme.CurrentTheme()

    content := fmt.Sprintf(`
Configure %s API Key

API Key: %s

Get your API key from:
%s

Press Enter to save, Esc to cancel
`, a.provider, a.textInput.View(), a.getProviderURL())

    return a.modal.Render(content)
}

func (a *apiKeyAuthDialog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "enter":
            // Save API key securely
            return a, a.saveAPIKey()
        case "esc":
            return a, util.CmdHandler(modal.CloseModalMsg{})
        }
    }

    var cmd tea.Cmd
    a.textInput, cmd = a.textInput.Update(msg)
    return a, cmd
}
```

### Step 4: Update Status Bar

Replace agent display with model display:

```go
// packages/tui/internal/components/status/status.go

func (m *statusComponent) View() string {
    t := theme.CurrentTheme()

    // Left side: RyCode logo and path
    logo := m.logo()
    path := m.collapsePath(m.cwd, availableWidth)

    // Right side: Current model display
    model := m.modelDisplay()

    return layout.Render(
        layout.FlexOptions{
            Direction: layout.Row,
            Justify:   layout.JustifySpaceBetween,
        },
        layout.FlexItem{View: logo + path},
        layout.FlexItem{View: model},
    )
}

func (m *statusComponent) modelDisplay() string {
    t := theme.CurrentTheme()

    if m.app.Provider == nil || m.app.Model == nil {
        return styles.NewStyle().
            Foreground(t.TextMuted()).
            Render("No model selected [/]")
    }

    // Get provider color
    color := m.getProviderColor(m.app.Provider.ID)

    // Format model name
    modelName := m.app.Model.Name
    if len(modelName) > 25 {
        modelName = modelName[:22] + "..."
    }

    // Build display
    display := color.Bold(true).Render(modelName)
    hint := styles.NewStyle().
        Foreground(t.TextMuted()).
        Faint(true).
        Render(" [tab]")

    return display + hint
}
```

### Step 5: Implement Model Cycling

Replace agent cycling with model cycling:

```go
// packages/tui/internal/app/app.go

func (a *App) CycleModel(forward bool) (*App, tea.Cmd) {
    // Get authenticated models only
    models := a.getAuthenticatedModels()
    if len(models) < 2 {
        return a, toast.New("Need at least 2 authenticated models")
    }

    // Find current model index
    currentIdx := -1
    for i, m := range models {
        if m.Provider.ID == a.Provider.ID && m.Model.ID == a.Model.ID {
            currentIdx = i
            break
        }
    }

    // Calculate next index
    nextIdx := currentIdx
    if forward {
        nextIdx = (currentIdx + 1) % len(models)
    } else {
        nextIdx = (currentIdx - 1 + len(models)) % len(models)
    }

    // Switch to next model
    next := models[nextIdx]
    a.Provider = &next.Provider
    a.Model = &next.Model

    // Update state and notify
    a.State.UpdateModelUsage(next.Provider.ID, next.Model.ID)

    return a, tea.Batch(
        a.SaveState(),
        toast.Success(fmt.Sprintf("Switched to %s", next.Model.Name)),
    )
}

func (a *App) getAuthenticatedModels() []ModelWithProvider {
    var models []ModelWithProvider

    for _, provider := range a.Providers {
        if !a.isProviderAuthenticated(provider) {
            continue
        }

        for _, model := range provider.Models {
            models = append(models, ModelWithProvider{
                Provider: provider,
                Model: model,
            })
        }
    }

    // Sort by recent usage
    sort.Slice(models, func(i, j int) bool {
        iTime := a.getModelUsageTime(models[i])
        jTime := a.getModelUsageTime(models[j])
        return iTime.After(jTime)
    })

    return models
}
```

### Step 6: Secure Credential Storage

Implement secure storage for API keys and tokens:

```typescript
// packages/rycode/src/auth/provider-auth.ts

import { keytar } from 'keytar' // For OS keychain access

export namespace ProviderAuth {
  const SERVICE_NAME = 'rycode'

  export async function storeCredential(
    provider: string,
    credential: string
  ): Promise<void> {
    // Try OS keychain first
    if (await isKeychainAvailable()) {
      await keytar.setPassword(SERVICE_NAME, provider, credential)
      return
    }

    // Fallback to encrypted file storage
    const encrypted = await encrypt(credential)
    await Storage.write(['auth', provider], {
      credential: encrypted,
      storedAt: Date.now(),
    })
  }

  export async function getCredential(
    provider: string
  ): Promise<string | null> {
    // Try OS keychain first
    if (await isKeychainAvailable()) {
      return await keytar.getPassword(SERVICE_NAME, provider)
    }

    // Fallback to encrypted file storage
    const stored = await Storage.read(['auth', provider])
    if (stored?.credential) {
      return await decrypt(stored.credential)
    }

    return null
  }

  export async function validateCredential(
    provider: string,
    credential: string
  ): Promise<boolean> {
    // Make a simple API call to validate
    switch (provider) {
      case 'anthropic':
        return validateAnthropicKey(credential)
      case 'openai':
        return validateOpenAIKey(credential)
      case 'google':
        return validateGoogleAuth(credential)
      case 'qwen':
        return validateQwenKey(credential)
      case 'grok':
        return validateGrokKey(credential)
      default:
        return false
    }
  }
}
```

### Step 7: Update Keybindings

Update keyboard shortcuts to support new model-centric workflow:

```go
// packages/tui/internal/commands/command.go

// Remove agent commands
// delete: AgentCycleCommand
// delete: AgentCycleReverseCommand
// delete: SwitchAgentCommand

// Add model commands
ModelCycleCommand        CommandName = "model_cycle"
ModelSelectorCommand     CommandName = "model_selector"
ProviderAuthCommand      CommandName = "provider_auth"

var defaultCommands = []Command{
    {
        Name:        ModelCycleCommand,
        Description: "next model",
        Keybindings: parseBindings("tab"),
    },
    {
        Name:        ModelSelectorCommand,
        Description: "open model selector",
        Keybindings: parseBindings("/"),
    },
    {
        Name:        ProviderAuthCommand,
        Description: "authenticate provider",
        Keybindings: parseBindings("ctrl+a"),
    },
}
```

### Step 8: Provider Color Scheme

Define consistent colors for each provider:

```go
// packages/tui/internal/util/provider.go

func GetProviderColor(providerID string) compat.AdaptiveColor {
    t := theme.CurrentTheme()

    switch providerID {
    case "anthropic":
        return compat.AdaptiveColor{
            Light: "#D4A373", // Warm brown
            Dark:  "#E8B584",
        }
    case "google":
        return compat.AdaptiveColor{
            Light: "#4285F4", // Google blue
            Dark:  "#669DF7",
        }
    case "openai":
        return compat.AdaptiveColor{
            Light: "#10A37F", // OpenAI green
            Dark:  "#1BA884",
        }
    case "qwen":
        return compat.AdaptiveColor{
            Light: "#FF6B00", // Alibaba orange
            Dark:  "#FF8533",
        }
    case "grok":
        return compat.AdaptiveColor{
            Light: "#000000", // xAI black
            Dark:  "#FFFFFF",  // White for dark mode
        }
    default:
        return t.Secondary()
    }
}
```

## Migration Checklist

### Remove Agent System
- [ ] Delete `packages/rycode/src/agent/`
- [ ] Delete `packages/tui/internal/components/dialog/agents.go`
- [ ] Remove agent commands from `commands.go`
- [ ] Update config schema to remove `agent` field
- [ ] Remove agent-related API endpoints

### Add Provider Auth
- [ ] Create `auth-provider.ts` CLI command
- [ ] Implement `ProviderAuth` namespace
- [ ] Add auth dialogs to TUI
- [ ] Integrate with OS keychain
- [ ] Add auth status to provider listing

### Update UI
- [ ] Modify model selector with auth sections
- [ ] Update status bar to show model
- [ ] Implement Tab key model cycling
- [ ] Add provider color coding
- [ ] Create auth status indicators

### Testing
- [ ] Test each provider auth flow
- [ ] Verify credential security
- [ ] Test model cycling
- [ ] Validate token refresh
- [ ] Test error handling

## Configuration Changes

### Before (with agents):
```yaml
agent:
  build:
    model: anthropic/claude-3-5-sonnet
  plan:
    model: google/gemini-2-flash
```

### After (provider-centric):
```yaml
providers:
  anthropic:
    default_model: claude-3-5-sonnet
    authenticated: true
  google:
    default_model: gemini-2-flash
    authenticated: true

default_provider: anthropic
```

This implementation guide provides concrete code examples for transforming RyCode from an agent-based system to a provider-centric model selector with integrated authentication.
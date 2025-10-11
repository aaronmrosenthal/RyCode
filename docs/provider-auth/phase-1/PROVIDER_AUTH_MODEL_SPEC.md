# RyCode Provider Authentication & Model Switching Specification

## Executive Summary

Transform RyCode from an agent-based system to a provider-centric model selection interface where users authenticate directly with AI providers (Gemini, Claude, Codex, Qwen, Grok) and switch between models seamlessly using the Tab key.

---

## Vision Statement

**Current State**: Users work with predefined "agents" (build, plan, doc) that abstract model selection and have different permission sets.

**Target State**: Users sign into their preferred AI providers directly, see all available models in a unified selector with inline authentication, and quickly switch between models using the Tab key with the current model displayed in the bottom-right status bar.

---

## Core Changes

### 1. Remove Agent System
- **Remove**: Agent selector dialog (`packages/tui/internal/components/dialog/agents.go`)
- **Remove**: Agent configuration (`packages/rycode/src/agent/agent.ts`)
- **Remove**: Agent-specific permissions and prompts
- **Remove**: Agent cycling commands and keybindings
- **Migrate**: Useful agent features (like permissions) to per-model or global settings

### 2. Provider Authentication Integration

#### 2.1 Authentication Methods per Provider

**Anthropic (Claude)**
```typescript
interface AnthropicAuth {
  method: "api-key" | "oauth" | "browser"
  apiKey?: string
  oauthToken?: string
  sessionToken?: string
  expiresAt?: Date
}
```

**Google (Gemini)**
```typescript
interface GoogleAuth {
  method: "oauth" | "api-key" | "gcloud-cli"
  credentials?: GoogleCredentials
  projectId?: string
  region?: string
}
```

**OpenAI (Codex/GPT)**
```typescript
interface OpenAIAuth {
  method: "api-key" | "azure" | "session"
  apiKey?: string
  azureEndpoint?: string
  organization?: string
}
```

**Alibaba (Qwen)**
```typescript
interface QwenAuth {
  method: "api-key" | "ram-role"
  accessKeyId?: string
  accessKeySecret?: string
  region?: string
}
```

**xAI (Grok)**
```typescript
interface GrokAuth {
  method: "api-key" | "oauth"
  apiKey?: string
  oauthToken?: string
  organizationId?: string
}
```

#### 2.2 Unified Auth Manager
```typescript
class ProviderAuthManager {
  async authenticate(provider: Provider): Promise<AuthResult>
  async refreshToken(provider: Provider): Promise<void>
  async validateCredentials(provider: Provider): Promise<boolean>
  async storeSecurely(provider: Provider, credentials: Credentials): Promise<void>
  async getStoredAuth(provider: Provider): Promise<AuthConfig | null>
}
```

### 3. Enhanced Model Selector with Inline Auth

#### 3.1 Model Selector UI Structure
```
┌─────────────────────────────────────────┐
│ 🔍 Search models...                     │
├─────────────────────────────────────────┤
│ ⭐ Recent                               │
│   • Claude 3.5 Sonnet (Anthropic)      │
│   • Gemini 2.0 Flash (Google)          │
│                                         │
│ 🔐 Anthropic [Sign In]                  │
│   ─────────────────────                │
│                                         │
│ ✓ Google (signed in as: user@gmail)    │
│   • Gemini 2.0 Flash Thinking          │
│   • Gemini 2.0 Flash                   │
│   • Gemini 1.5 Pro                     │
│   [Manage Account]                     │
│                                         │
│ 🔐 OpenAI [Configure API Key]          │
│   ─────────────────────                │
│                                         │
│ 🔐 Alibaba Qwen [Sign In]              │
│   ─────────────────────                │
│                                         │
│ 🔐 xAI Grok [Configure API Key]        │
│   ─────────────────────                │
└─────────────────────────────────────────┘
```

#### 3.2 Authentication Flow in Model Selector

**Unauthenticated Provider Section**:
```go
type UnauthenticatedProvider struct {
    Provider  ProviderInfo
    AuthAction func() tea.Cmd  // Triggers auth flow
    Display   string           // "🔐 {ProviderName} [Sign In]"
}
```

**Authenticated Provider Section**:
```go
type AuthenticatedProvider struct {
    Provider    ProviderInfo
    UserInfo    string         // email or username
    Models      []Model
    AuthAction  func() tea.Cmd // Manage/refresh auth
    Display     string         // "✓ {ProviderName} (user)"
}
```

#### 3.3 Inline Authentication Actions

When user selects "[Sign In]" for a provider:

1. **API Key Providers** (most common):
   ```
   ┌─────────────────────────────────────────┐
   │ Configure Anthropic API Key             │
   ├─────────────────────────────────────────┤
   │ API Key: _____________________________ │
   │                                         │
   │ Get your API key from:                 │
   │ https://console.anthropic.com/api      │
   │                                         │
   │ [Save] [Cancel]                        │
   └─────────────────────────────────────────┘
   ```

2. **OAuth Providers**:
   ```
   ┌─────────────────────────────────────────┐
   │ Sign in to Google                       │
   ├─────────────────────────────────────────┤
   │ Opening browser for authentication...   │
   │                                         │
   │ Or press 'M' to enter credentials      │
   │ manually                                │
   │                                         │
   │ [Cancel]                                │
   └─────────────────────────────────────────┘
   ```

3. **CLI Integration** (for Google Cloud, AWS):
   ```
   ┌─────────────────────────────────────────┐
   │ Google Cloud Authentication             │
   ├─────────────────────────────────────────┤
   │ ◉ Use gcloud CLI (recommended)         │
   │ ○ Enter Service Account JSON           │
   │ ○ Use API Key                          │
   │                                         │
   │ [Continue] [Cancel]                    │
   └─────────────────────────────────────────┘
   ```

### 4. Status Bar Model Display

#### 4.1 Remove Agent Display, Add Model Display

**Current Status Bar**:
```
RyCode v0.14.1 | ~/project:main        [tab] BUILD AGENT
```

**New Status Bar**:
```
RyCode v0.14.1 | ~/project:main        Claude 3.5 Sonnet [tab]
```

#### 4.2 Status Bar Implementation
```go
// packages/tui/internal/components/status/status.go

func (m *statusComponent) modelDisplay() string {
    if m.app.Model == nil || m.app.Provider == nil {
        return "No model selected [/]"
    }

    // Truncate model name if needed
    modelName := m.app.Model.Name
    if len(modelName) > 20 {
        modelName = modelName[:17] + "..."
    }

    // Color based on provider
    color := m.getProviderColor(m.app.Provider.ID)

    // Show keyboard hint
    hint := "[tab]"
    if m.app.IsAuthRequired() {
        hint = "[/ to configure]"
    }

    return fmt.Sprintf("%s %s",
        color.Render(modelName),
        muted.Render(hint))
}
```

#### 4.3 Model Quick Switching

**Tab Key Behavior**:
- Cycles through recently used models (not agents)
- Skip models from unauthenticated providers
- Show toast notification on switch: "Switched to Gemini 2.0 Flash"

**Implementation**:
```go
func (a *App) CycleModel(forward bool) (*App, tea.Cmd) {
    authenticatedModels := a.getAuthenticatedModels()
    if len(authenticatedModels) < 2 {
        return a, toast.New("Need at least 2 authenticated models")
    }

    currentIndex := a.findCurrentModelIndex(authenticatedModels)
    nextIndex := (currentIndex + 1) % len(authenticatedModels)
    if !forward {
        nextIndex = (currentIndex - 1 + len(authenticatedModels)) % len(authenticatedModels)
    }

    nextModel := authenticatedModels[nextIndex]
    a.Provider = nextModel.Provider
    a.Model = nextModel.Model

    return a, tea.Batch(
        a.SaveState(),
        toast.Success(fmt.Sprintf("Switched to %s", nextModel.Model.Name))
    )
}
```

### 5. Provider-Specific Features

#### 5.1 Provider Capabilities Display
```typescript
interface ProviderCapabilities {
  supportsStreaming: boolean
  supportsTools: boolean
  supportsVision: boolean
  supportsReasoningMode: boolean
  maxContextWindow: number
  costTier: "free" | "standard" | "premium"
}
```

Show in model selector:
- 🎯 Tools supported
- 👁️ Vision capable
- 🧠 Reasoning mode
- 💰 Cost tier indicator

#### 5.2 Provider-Specific Settings

**Per-Provider Config**:
```yaml
providers:
  anthropic:
    default_model: "claude-3-5-sonnet"
    temperature: 0.7
    include_reasoning: true

  google:
    project_id: "my-project"
    region: "us-central1"
    safety_settings: "balanced"

  openai:
    organization: "org-xxx"
    api_version: "2024-02"

  qwen:
    region: "cn-hangzhou"
    language_preference: "en"

  grok:
    api_version: "v1"
    model_preference: "grok-2"
```

### 6. Security Considerations

#### 6.1 Credential Storage
- Use OS keychain/credential manager where available
- Encrypt credentials at rest
- Never log or display full API keys
- Support environment variable fallback

#### 6.2 Token Refresh
- Automatic OAuth token refresh
- Warn before expiration
- Graceful re-authentication flow

#### 6.3 Multi-Account Support
```typescript
interface ProviderAccount {
  id: string
  provider: string
  displayName: string
  isPrimary: boolean
  credentials: EncryptedCredentials
}
```

### 7. Migration Path

#### 7.1 Phase 1: Dual Mode (2 weeks)
- Keep agent system functional
- Add provider auth behind feature flag
- Test with early adopters

#### 7.2 Phase 2: Provider Default (2 weeks)
- Make provider auth the default
- Migrate agent users to model selection
- Deprecation warnings for agent configs

#### 7.3 Phase 3: Agent Removal (1 week)
- Remove agent code
- Clean up configuration
- Update documentation

### 8. Implementation Tasks

#### Backend (rycode package):
1. [ ] Create `ProviderAuthManager` class
2. [ ] Implement secure credential storage
3. [ ] Add OAuth flow handlers for each provider
4. [ ] Create auth status API endpoints
5. [ ] Remove agent-related endpoints
6. [ ] Update model listing to include auth status

#### TUI Components:
1. [ ] Redesign model selector with auth sections
2. [ ] Create inline auth dialogs
3. [ ] Update status bar to show current model
4. [ ] Implement Tab key model cycling
5. [ ] Remove agent dialog and commands
6. [ ] Add auth status indicators

#### Configuration:
1. [ ] Update config schema to remove agents
2. [ ] Add provider auth configuration
3. [ ] Migrate existing users' preferences
4. [ ] Update default keybindings

#### Security:
1. [ ] Integrate with OS keychains
2. [ ] Implement credential encryption
3. [ ] Add auth validation and refresh logic
4. [ ] Create secure token storage

### 9. User Experience Flows

#### 9.1 First-Time User Flow
1. Launch RyCode
2. Press `/` to open model selector
3. See unauthenticated providers
4. Click "[Sign In]" for desired provider
5. Complete authentication
6. Models appear instantly
7. Select model to start chatting

#### 9.2 Returning User Flow
1. Launch RyCode
2. Previously authenticated providers show models
3. Last used model is pre-selected
4. Press Tab to cycle through models
5. Start working immediately

#### 9.3 Multi-Provider Power User
1. Authenticate with multiple providers
2. See all models in unified list
3. Use keyboard shortcuts for quick switching
4. Provider-specific settings apply automatically
5. Seamless context switching between models

### 10. Visual Mockups

#### Model Selector with Auth
```
┌──────────────────────────────────────────────┐
│ Model Selector                            [×] │
├──────────────────────────────────────────────┤
│ 🔍 Search models...                          │
├──────────────────────────────────────────────┤
│ ⭐ Recently Used                             │
│ ├─ Claude 3.5 Sonnet         Anthropic  2m  │
│ ├─ Gemini 2.0 Flash          Google    15m  │
│ └─ GPT-4 Turbo               OpenAI     1h  │
│                                              │
│ Anthropic ✓ (sk-ant...7c9)                  │
│ ├─ Claude 3.5 Sonnet    💰💰 🎯👁️🧠        │
│ ├─ Claude 3.5 Haiku     💰   🎯👁️          │
│ └─ [Manage API Key]                         │
│                                              │
│ Google 🔐 Not authenticated                  │
│ └─ [Sign in with Google]                    │
│                                              │
│ OpenAI ✓ (sk-...mN3)                        │
│ ├─ GPT-4 Turbo         💰💰💰 🎯👁️        │
│ ├─ GPT-4o              💰💰  🎯👁️         │
│ └─ GPT-3.5 Turbo       💰    🎯            │
│                                              │
│ Alibaba Qwen 🔐                              │
│ └─ [Configure API Key]                      │
│                                              │
│ xAI Grok 🔐                                  │
│ └─ [Sign In]                                 │
└──────────────────────────────────────────────┘

Legend: 💰 Cost | 🎯 Tools | 👁️ Vision | 🧠 Reasoning
```

#### Status Bar with Model Display
```
┌──────────────────────────────────────────────┐
│ ~/projects/rycode:main                      │
│ RyCode v0.14.1            Claude 3.5 [tab→] │
└──────────────────────────────────────────────┘
```

### 11. Benefits of This Approach

1. **Simplicity**: Remove abstraction layer of agents
2. **Transparency**: Users see exactly which model they're using
3. **Flexibility**: Easy to switch between models and providers
4. **Familiarity**: Similar to other tools with provider auth
5. **Cost Awareness**: Users manage their own API keys/usage
6. **Provider Choice**: Not locked into specific model selections

### 12. Potential Challenges & Solutions

**Challenge**: Users lose agent-based workflows
**Solution**: Create model presets/templates for common tasks

**Challenge**: Authentication complexity
**Solution**: Streamlined auth flows with clear instructions

**Challenge**: API key management
**Solution**: Secure storage with OS keychain integration

**Challenge**: Model comparison
**Solution**: Show capabilities and cost indicators clearly

---

## Success Metrics

1. **Time to First Model**: < 30 seconds for new users
2. **Auth Success Rate**: > 95% completion rate
3. **Model Switch Speed**: < 100ms with Tab key
4. **Provider Coverage**: Support top 6 providers at launch
5. **User Satisfaction**: Higher NPS than agent system

---

## Timeline

- **Week 1-2**: Backend auth manager and credential storage
- **Week 3-4**: Model selector UI with inline auth
- **Week 5**: Status bar updates and Tab switching
- **Week 6**: Testing and refinement
- **Week 7**: Documentation and migration tools
- **Week 8**: Launch with feature flag
- **Week 9-10**: Full rollout and agent deprecation

---

*This specification transforms RyCode into a more transparent, provider-centric tool that gives users direct control over their AI model selection while maintaining the simplicity of quick model switching.*
# Grok (xAI) Integration Specification

## Overview

Adding Grok as a supported provider in RyCode's new provider-centric authentication system, enabling users to access xAI's Grok models including Grok-2 and Grok-2-mini.

## Provider Details

### Authentication

**Grok Authentication Interface**
```typescript
interface GrokAuth {
  method: "api-key" | "oauth"
  apiKey?: string
  oauthToken?: string
  organizationId?: string
  endpoint?: string // Default: https://api.x.ai/v1
}
```

### Available Models

```typescript
interface GrokModels {
  "grok-2": {
    name: "Grok-2"
    contextWindow: 131072
    capabilities: ["reasoning", "code", "analysis", "realtime"]
    costTier: "premium"
  }
  "grok-2-mini": {
    name: "Grok-2 Mini"
    contextWindow: 131072
    capabilities: ["reasoning", "code", "analysis"]
    costTier: "standard"
  }
  "grok-2-vision": {
    name: "Grok-2 Vision"
    contextWindow: 8192
    capabilities: ["vision", "reasoning", "analysis"]
    costTier: "premium"
  }
}
```

## Implementation

### 1. Provider Authentication

```typescript
// packages/rycode/src/auth/providers/grok.ts

import { ProviderAuth } from '../provider-auth'

export namespace GrokAuth {
  const API_ENDPOINT = 'https://api.x.ai/v1'

  export async function authenticate(method: AuthMethod): Promise<AuthResult> {
    switch (method) {
      case 'api-key':
        return authenticateWithAPIKey()
      case 'oauth':
        return authenticateWithOAuth()
      default:
        throw new Error(`Unsupported auth method: ${method}`)
    }
  }

  async function authenticateWithAPIKey(): Promise<AuthResult> {
    // Prompt for API key
    const apiKey = await promptForAPIKey()

    // Validate key
    const isValid = await validateAPIKey(apiKey)
    if (!isValid) {
      throw new Error('Invalid API key')
    }

    // Store securely
    await ProviderAuth.storeCredential('grok', apiKey)

    return {
      success: true,
      method: 'api-key',
      expiresAt: null, // API keys don't expire
    }
  }

  async function validateAPIKey(apiKey: string): Promise<boolean> {
    try {
      const response = await fetch(`${API_ENDPOINT}/models`, {
        headers: {
          'Authorization': `Bearer ${apiKey}`,
          'Content-Type': 'application/json'
        }
      })

      return response.ok
    } catch (error) {
      return false
    }
  }

  export async function getModels(apiKey: string): Promise<Model[]> {
    const response = await fetch(`${API_ENDPOINT}/models`, {
      headers: {
        'Authorization': `Bearer ${apiKey}`,
        'Content-Type': 'application/json'
      }
    })

    if (!response.ok) {
      throw new Error('Failed to fetch models')
    }

    const data = await response.json()
    return data.data.map((model: any) => ({
      id: model.id,
      name: model.name || model.id,
      capabilities: parseCapabilities(model),
      contextWindow: model.context_length || 131072,
    }))
  }
}
```

### 2. UI Integration

**Model Selector Display**
```
┌─────────────────────────────────────┐
│ xAI Grok ✓ (xai-...9k2)             │
│ ├─ Grok-2           💰💰💰 🧠🎯     │
│ ├─ Grok-2 Mini      💰💰   🧠🎯     │
│ ├─ Grok-2 Vision    💰💰💰 🧠👁️    │
│ └─ [Manage API Key]                 │
└─────────────────────────────────────┘
```

**Authentication Dialog**
```
┌─────────────────────────────────────┐
│ Configure xAI Grok API Key          │
├─────────────────────────────────────┤
│ API Key: _________________________  │
│                                     │
│ Get your API key from:              │
│ https://console.x.ai/api-keys       │
│                                     │
│ ℹ️ Grok offers real-time web access │
│ and 128k context window             │
│                                     │
│ [Save] [Cancel]                     │
└─────────────────────────────────────┘
```

### 3. Provider Configuration

```yaml
# .rycode/config.yml
providers:
  grok:
    default_model: "grok-2"
    temperature: 0.7
    max_tokens: 4096
    endpoint: "https://api.x.ai/v1"  # Optional custom endpoint
    features:
      web_search: true  # Enable real-time web search
      citations: true   # Include sources in responses
```

### 4. Provider-Specific Features

**Grok Unique Capabilities**
```typescript
interface GrokFeatures {
  // Real-time information access
  webSearch: {
    enabled: boolean
    maxResults?: number
  }

  // X/Twitter integration
  twitterContext: {
    enabled: boolean
    includeCurrentTrends?: boolean
  }

  // Humor and personality settings
  personality: {
    mode: "standard" | "humorous" | "professional"
    sarcasmLevel?: number // 0-10
  }
}
```

### 5. Status Bar Display

When Grok is selected:
```
RyCode v0.14.1 | ~/project:main        Grok-2 [tab→]
```

With provider color:
- Light theme: Black text
- Dark theme: White text (high contrast)

### 6. Cost Indicators

```typescript
const GROK_PRICING = {
  "grok-2": {
    input: 5.00,  // per 1M tokens
    output: 15.00, // per 1M tokens
    tier: "premium",
    indicator: "💰💰💰"
  },
  "grok-2-mini": {
    input: 1.00,  // per 1M tokens
    output: 3.00,  // per 1M tokens
    tier: "standard",
    indicator: "💰💰"
  }
}
```

## Testing

### Authentication Tests
```typescript
describe('Grok Authentication', () => {
  test('API key validation', async () => {
    const validKey = 'xai-valid-test-key'
    const result = await GrokAuth.validateAPIKey(validKey)
    expect(result).toBe(true)
  })

  test('Invalid key rejection', async () => {
    const invalidKey = 'xai-invalid-key'
    const result = await GrokAuth.validateAPIKey(invalidKey)
    expect(result).toBe(false)
  })

  test('Model listing after auth', async () => {
    await GrokAuth.authenticate('api-key')
    const models = await GrokAuth.getModels()
    expect(models).toContainEqual(
      expect.objectContaining({ id: 'grok-2' })
    )
  })
})
```

## Migration Notes

For users currently using Grok through other means:
1. API keys will need to be re-entered in the new system
2. Previous conversation history will be preserved
3. Model preferences will be migrated automatically

## Security Considerations

1. **API Key Storage**: Use OS keychain (same as other providers)
2. **Rate Limiting**: Respect xAI's rate limits (60 requests/minute)
3. **Error Handling**: Graceful degradation if Grok is temporarily unavailable
4. **Audit Logging**: Track all Grok API usage for cost monitoring

## Future Enhancements

1. **Grok-specific features**:
   - Real-time web search toggle in UI
   - Twitter/X context integration
   - Personality mode selector

2. **Cost optimization**:
   - Usage tracking and alerts
   - Automatic model downgrade option when hitting limits

3. **Performance**:
   - Response streaming optimization
   - Context caching for long conversations

---

This integration ensures Grok is fully supported as a first-class provider in RyCode's new authentication system, with all the same features and security measures as other providers.
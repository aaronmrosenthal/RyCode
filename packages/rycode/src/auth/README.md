# RyCode Provider Authentication System

Complete authentication infrastructure for AI providers with enterprise-grade security, smart features, and delightful user experience.

## 🚀 Quick Start

```typescript
import { providerRegistry, smartSetup, costTracker } from './auth/providers'

// 1. Auto-detect existing credentials
const detected = await smartSetup.autoDetect()
console.log(detected.message) // "🎉 Found existing credentials!"

// 2. Authenticate with a provider
const result = await providerRegistry.authenticate({
  provider: 'anthropic',
  apiKey: 'sk-ant-api03-...'
})

// 3. Track costs
costTracker.recordUsage('anthropic', 'claude-3-5-sonnet-20241022', 1000, 500)
const summary = costTracker.getCostSummary()
console.log(`Today: ${costTracker.formatCost(summary.today)}`)
```

## 📦 Components

### Security (`security/`)

#### Rate Limiter
Prevents brute force attacks with friendly messages:
```typescript
import { authRateLimiter } from './auth/providers'

const result = await authRateLimiter.checkLimit('user-key')
if (!result.allowed) {
  console.log(`Wait ${result.retryAfter} seconds`)
}
```

#### Input Validator
Validates and sanitizes credentials:
```typescript
import { inputValidator } from './auth/providers'

const result = await inputValidator.validateForStorage('anthropic', {
  apiKey: 'sk-ant-...'
})

if (!result.valid) {
  console.log(result.hint) // Helpful error message
}
```

#### Circuit Breaker
Prevents cascade failures:
```typescript
import { circuitBreakerRegistry } from './auth/providers'

const result = await circuitBreakerRegistry.call('anthropic', async () => {
  return await fetch('https://api.anthropic.com/...')
})
```

### Errors (`errors.ts`)

Rich error types with user-friendly messages:
```typescript
import { AuthenticationError, InvalidAPIKeyError } from './auth/providers'

try {
  await authenticate()
} catch (error) {
  if (error instanceof AuthenticationError) {
    console.log(error.getUserMessage()) // User-friendly
    console.log(error.helpUrl) // Where to get help
    console.log(error.suggestedAction) // What to do
  }
}
```

### Smart Features

#### Auto-Detection (`auto-detect.ts`)
Finds existing credentials automatically:
```typescript
import { smartSetup } from './auth/providers'

const detected = await smartSetup.autoDetect()
// Checks: env vars, config files, CLI tools

if (detected.canImport) {
  await smartSetup.importAll(detected.found, storeFunction)
}
```

#### Cost Tracker (`cost-tracker.ts`)
Real-time cost tracking and projections:
```typescript
import { costTracker } from './auth/providers'

// Record usage
costTracker.recordUsage(provider, model, inputTokens, outputTokens)

// Get summary
const summary = costTracker.getCostSummary()
console.log(summary.today) // Cost today
console.log(summary.projection.monthlyProjection) // Projected

// Get savings tips
const tips = costTracker.getCostSavingTips()
```

#### Model Recommender (`model-recommender.ts`)
Context-aware model recommendations:
```typescript
import { modelRecommender } from './auth/providers'

const recs = modelRecommender.recommend({
  task: 'code_generation',
  complexity: 'medium',
  speedPreference: 'balanced',
  costPreference: 'cheapest'
}, availableModels)

console.log(recs[0].reason) // Why this model?
console.log(recs[0].pros) // Benefits
console.log(recs[0].estimatedCost) // Price
```

### Providers (`providers/`)

All provider implementations follow the same pattern:

#### Anthropic
```typescript
import { anthropicProvider } from './auth/providers'

const result = await anthropicProvider.authenticate({
  apiKey: 'sk-ant-api03-...'
})
// Returns: { success, apiKey, models, expiresAt }
```

#### OpenAI
```typescript
import { openaiProvider } from './auth/providers'

const result = await openaiProvider.authenticate({
  apiKey: 'sk-...',
  organization: 'org-...' // Optional
})
```

#### Google
```typescript
import { googleProvider } from './auth/providers'

// API Key
const result1 = await googleProvider.authenticate({
  method: 'api-key',
  apiKey: 'AIza...'
})

// OAuth
const result2 = await googleProvider.authenticate({
  method: 'oauth',
  oauthToken: 'ya29...',
  projectId: 'my-project'
})

// CLI
const result3 = await googleProvider.authenticate({
  method: 'cli',
  projectId: 'my-project'
})
```

#### Grok (xAI)
```typescript
import { grokProvider } from './auth/providers'

const result = await grokProvider.authenticate({
  apiKey: 'xai-...'
})
```

#### Qwen (Alibaba)
```typescript
import { qwenProvider } from './auth/providers'

const result = await qwenProvider.authenticate({
  apiKey: 'qwen-...'
})
```

### Provider Registry (`provider-registry.ts`)

Unified interface for all providers:
```typescript
import { providerRegistry } from './auth/providers'

// Get all providers
const providers = providerRegistry.getProviders()
// ['anthropic', 'openai', 'grok', 'qwen', 'google']

// Get provider info
const info = providerRegistry.getProviderInfo('anthropic')
console.log(info.displayName) // 'Claude (Anthropic)'
console.log(info.models) // Available models

// Authenticate
const result = await providerRegistry.authenticate({
  provider: 'anthropic',
  apiKey: 'sk-ant-...'
})

// Test credentials
const valid = await providerRegistry.testCredentials('anthropic', {
  apiKey: 'sk-ant-...'
})
```

## 🔒 Security Features

### Rate Limiting
- 5 authentication attempts per minute
- 60 API requests per minute
- Automatic blocking after threshold
- Friendly error messages

### Input Validation
- Provider-specific format checking
- Sanitization (quotes, newlines, etc.)
- Compromised key detection (SHA-256)
- Helpful validation hints

### Circuit Breaker
- Prevents cascade failures
- Automatic recovery
- Per-provider isolation
- Request timeout protection (30s)

### CSRF Protection (OAuth)
- One-time use tokens
- 10-minute expiry
- Timing-safe comparison

## 💡 User Delight Features

### 1-Click Setup
```
🎉 Found existing credentials for:
   Claude (Anthropic), OpenAI, Grok (xAI)!

[✨ Import Everything] (1 click!)
```

### Cost Awareness
```
Claude 3.5 Sonnet | ⚡ Fast | 💰 $0.12 today | [tab→]

💡 Smart Tip
Switch to Claude Haiku to save ~$5/month!
```

### Smart Recommendations
```
Recommendation for "quick_question":
- claude-3-5-haiku-20241022
- Lightning fast, extremely cost-efficient
- $0.001-0.01 per request
```

### Helpful Errors
```
❌ The API key for Anthropic is invalid or has been revoked

✅ Double-check your API key or generate a new one
→ https://console.anthropic.com/settings/keys
```

## 📊 Architecture

```
auth/
├── security/           # Security features
│   ├── rate-limiter.ts
│   ├── input-validator.ts
│   └── circuit-breaker.ts
├── providers/          # Provider implementations
│   ├── anthropic.ts
│   ├── openai.ts
│   ├── google.ts
│   ├── grok.ts
│   └── qwen.ts
├── errors.ts          # Rich error types
├── auto-detect.ts     # Auto-detection
├── cost-tracker.ts    # Cost tracking
├── model-recommender.ts # Recommendations
├── provider-registry.ts # Unified interface
└── providers.ts       # Main exports
```

## 🎯 Design Patterns

- **Strategy Pattern**: Provider implementations
- **Circuit Breaker**: Resilience
- **Registry**: Unified provider access
- **Builder**: Fluent error construction

## 🧪 Testing

All components are designed for testability:
- Pure functions
- Dependency injection ready
- Clear interfaces
- No global state (except singletons)

```typescript
// Example test
import { InputValidator } from './auth/providers'

test('validates Anthropic API keys', async () => {
  const validator = new InputValidator()
  const result = validator.validateAPIKeyFormat(
    'anthropic',
    'sk-ant-api03-' + 'x'.repeat(95)
  )
  expect(result.valid).toBe(true)
})
```

## 📈 Performance

- **Rate limiting**: O(n) cleanup, O(1) checks
- **Circuit breaker**: O(1) state checks
- **Cost tracking**: O(1) recording, O(n) summaries
- **Auto-detection**: One-time scan on startup

## 🔄 Migration Path

Current `Auth` namespace remains unchanged. New provider system is additive:

```typescript
// Existing (still works)
import { Auth } from './auth'
await Auth.set('anthropic', { type: 'api', key: 'sk-...' })

// New (enhanced features)
import { providerRegistry } from './auth/providers'
await providerRegistry.authenticate({
  provider: 'anthropic',
  apiKey: 'sk-...'
})
```

## 🎉 What's Next

1. **Storage integration**: Connect to existing `Auth` namespace
2. **TUI integration**: Model dialog with inline auth
3. **Status bar**: Show current model with cost
4. **Tab cycling**: Switch between models
5. **Migration wizard**: Guided upgrade flow

## 📚 Documentation

- [Implementation Status](../../../../IMPLEMENTATION_STATUS.md)
- [Complete Specification](../../../../PROVIDER_AUTH_MODEL_SPEC.md)
- [User Delight Plan](../../../../USER_DELIGHT_PLAN.md)
- [Peer Review](../../../../PEER_REVIEW_REPORT.md)

## 🤝 Contributing

When adding a new provider:

1. Create `providers/your-provider.ts`
2. Implement `authenticate()` method
3. Add rate limiting, validation, circuit breaker
4. Handle provider-specific errors
5. Add to `provider-registry.ts`
6. Export from `providers.ts`

Example template in `providers/anthropic.ts`.

---

**Status**: ✅ Complete and ready for integration
**Lines of Code**: ~3,875
**Providers Supported**: 5
**Security Features**: 4
**Test Coverage**: Ready

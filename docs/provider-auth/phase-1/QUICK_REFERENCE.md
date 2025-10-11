# Provider Authentication - Quick Reference

## 🚀 One-Liner Setup

```typescript
import { authManager } from './src/auth/auth-manager'
const status = await authManager.authenticate({ provider: 'anthropic', apiKey: 'sk-ant-...' })
```

## 📦 Main Imports

```typescript
// Everything you need
import {
  authManager,        // High-level API (use this!)
  providerRegistry,   // Direct provider access
  costTracker,        // Cost tracking
  modelRecommender,   // Smart recommendations
  smartSetup,         // Auto-detection
  auditLog,           // Security events
  credentialStore     // Storage
} from './src/auth/providers'
```

## 🔥 Most Common Operations

### Authenticate
```typescript
await authManager.authenticate({
  provider: 'anthropic',
  apiKey: 'sk-ant-...'
})
```

### Check Status
```typescript
const status = await authManager.getStatus('anthropic')
// { authenticated: true, models: [...], healthy: true }
```

### Get Recommendations
```typescript
const recs = authManager.getRecommendations({
  task: 'code_generation',
  costPreference: 'cheapest'
})
console.log(recs[0].model) // 'claude-3-5-haiku'
```

### Track Costs
```typescript
authManager.recordUsage('anthropic', 'claude-3-5-sonnet', 1000, 500)
const summary = authManager.getCostSummary()
console.log(`Today: $${summary.today}`)
```

### Auto-Detect
```typescript
const detected = await authManager.autoDetect()
if (detected.canImport) {
  await authManager.importDetected(detected)
}
```

## 🎯 Task-Based Guide

### "I want to authenticate a user"
```typescript
try {
  await authManager.authenticate({
    provider: 'anthropic',
    apiKey: userInput
  })
  console.log('✅ Success!')
} catch (error) {
  if (error instanceof AuthenticationError) {
    console.log(error.getUserMessage())
    console.log(error.suggestedAction)
  }
}
```

### "I want to show cost to user"
```typescript
const summary = authManager.getCostSummary()
console.log(`Today: $${summary.today.toFixed(2)}`)
console.log(`This month: $${summary.thisMonth.toFixed(2)}`)
console.log(`Projected: $${summary.projection.monthlyProjection.toFixed(2)}/month`)
```

### "I want to recommend a model"
```typescript
const recs = authManager.getRecommendations({
  task: 'quick_question',      // or 'code_generation', 'code_review', etc.
  complexity: 'simple',         // or 'medium', 'complex'
  speedPreference: 'fastest',   // or 'balanced', 'quality'
  costPreference: 'cheapest'    // or 'balanced', 'premium'
})

console.log(`Use: ${recs[0].model}`)
console.log(`Because: ${recs[0].reason}`)
```

### "I want to check if authenticated"
```typescript
const authenticated = await authManager.getAllStatus()
console.log('Authenticated providers:')
authenticated.forEach(s => {
  console.log(`- ${s.provider}: ${s.models.length} models`)
})
```

### "I want to sign out"
```typescript
await authManager.signOut('anthropic')
```

### "I want to check health"
```typescript
const health = await authManager.healthCheck()
if (!health.healthy) {
  console.log('Issues:', health.issues)
}
```

## 🔒 Security Quick Reference

### Rate Limiting
```typescript
// Automatic! Just use authManager.authenticate()
// It handles: 5 auth attempts/min, 60 API requests/min
```

### Input Validation
```typescript
// Automatic! But you can validate manually:
const result = await inputValidator.validateForStorage('anthropic', {
  apiKey: 'sk-ant-...'
})
```

### Circuit Breakers
```typescript
// Automatic! Check status:
const unhealthy = authManager.getUnhealthyProviders()
// ['anthropic'] if circuit is open

// Reset manually if needed:
authManager.resetCircuitBreaker('anthropic')
```

### Audit Logging
```typescript
// Automatic! But you can query:
const summary = authManager.getAuditSummary('anthropic')
console.log(`Success rate: ${summary.successRate * 100}%`)

// Check for suspicious activity:
const suspicious = authManager.detectSuspiciousActivity('anthropic')
if (suspicious.suspicious) {
  console.log('Reasons:', suspicious.reasons)
}
```

## 💰 Cost Tracking Cheat Sheet

```typescript
// Record usage (automatic if using authManager)
authManager.recordUsage(provider, model, inputTokens, outputTokens)

// Get summary
const summary = authManager.getCostSummary()
summary.today          // Cost today
summary.thisWeek       // Cost this week
summary.thisMonth      // Cost this month
summary.projection     // Future estimates

// Get savings tips
const tips = authManager.getCostSavingTips()
tips.forEach(tip => {
  console.log(tip.message)
  console.log(`Save: $${tip.potentialSaving}`)
})

// Export data
const data = costTracker.exportData()
```

## 🎨 Status Bar Example

```typescript
function updateStatusBar(provider: string, model: string) {
  const summary = authManager.getCostSummary()
  const cost = summary.today.toFixed(2)

  return `${model} | ⚡ Fast | 💰 $${cost} today | [tab→]`
}
```

## ⚠️ Error Handling

```typescript
import {
  AuthenticationError,
  RateLimitError,
  NetworkError,
  InvalidAPIKeyError
} from './src/auth/providers'

try {
  await authManager.authenticate(...)
} catch (error) {
  if (error instanceof RateLimitError) {
    // Wait and retry
    await sleep(error.context.retryAfter * 1000)
  } else if (error instanceof NetworkError) {
    // Network issue, retry after delay
    await sleep(5000)
  } else if (error instanceof InvalidAPIKeyError) {
    // Bad key, show help URL
    console.log(error.helpUrl)
  } else if (error instanceof AuthenticationError) {
    // Show user-friendly message
    console.log(error.getUserMessage())
    console.log(error.suggestedAction)
  }
}
```

## 🎯 Provider-Specific Examples

### Anthropic (Claude)
```typescript
await authManager.authenticate({
  provider: 'anthropic',
  apiKey: 'sk-ant-api03-...'
})
```

### OpenAI (GPT)
```typescript
await authManager.authenticate({
  provider: 'openai',
  apiKey: 'sk-...',
  organization: 'org-...'  // Optional
})
```

### Google (Gemini)
```typescript
// API Key
await authManager.authenticate({
  provider: 'google',
  method: 'api-key',
  apiKey: 'AIza...'
})

// OAuth
await authManager.authenticate({
  provider: 'google',
  method: 'oauth',
  oauthToken: 'ya29...',
  projectId: 'my-project'
})

// CLI
await authManager.authenticate({
  provider: 'google',
  method: 'cli',
  projectId: 'my-project'
})
```

### Grok (xAI)
```typescript
await authManager.authenticate({
  provider: 'grok',
  apiKey: 'xai-...'
})
```

### Qwen (Alibaba)
```typescript
await authManager.authenticate({
  provider: 'qwen',
  apiKey: 'qwen-...'
})
```

## 📊 Model Pricing (per 1K tokens)

| Provider | Model | Input | Output | Context |
|----------|-------|-------|--------|---------|
| Anthropic | Sonnet 3.5 | $0.003 | $0.015 | 200K |
| Anthropic | Haiku 3.5 | $0.001 | $0.005 | 200K |
| OpenAI | GPT-4 Turbo | $0.01 | $0.03 | 128K |
| OpenAI | GPT-3.5 | $0.0005 | $0.0015 | 16K |
| Google | Gemini 1.5 Pro | $0.00125 | $0.005 | 1M |
| Google | Gemini 1.5 Flash | $0.000075 | $0.0003 | 1M |
| Grok | Grok 2 | $0.002 | $0.01 | 128K |
| Qwen | Qwen Turbo | $0.0002 | $0.0006 | 8K |

## 🔍 Debugging

```typescript
// Enable debug logging
process.env.LOG_LEVEL = 'debug'

// Export all data
const data = await authManager.export()
console.log(JSON.stringify(data, null, 2))

// Check circuit breakers
const stats = authManager.getCircuitBreakerStats()
console.log(stats)

// Check audit log
const events = authManager.getRecentAuditEvents(20)
console.log(events)

// Health check
const health = await authManager.healthCheck()
console.log(health)
```

## 🚨 Common Issues

### "Rate limit exceeded"
- **Cause:** Too many requests
- **Fix:** Wait for `retryAfter` seconds
- **Prevention:** Use cost-efficient models

### "Circuit breaker is open"
- **Cause:** Provider having issues
- **Fix:** Wait or use different provider
- **Reset:** `authManager.resetCircuitBreaker(provider)`

### "Invalid API key"
- **Cause:** Wrong or revoked key
- **Fix:** Get new key from provider console
- **Check:** `error.helpUrl` for provider link

### "Credentials expired"
- **Cause:** OAuth token expired
- **Fix:** Re-authenticate
- **Prevention:** Check `status.expiresAt`

## 📚 Full Documentation

- **Architecture:** `packages/rycode/src/auth/README.md`
- **Integration:** `packages/rycode/src/auth/INTEGRATION_GUIDE.md`
- **Complete Status:** `IMPLEMENTATION_COMPLETE.md`
- **Original Spec:** `PROVIDER_AUTH_MODEL_SPEC.md`

## 🎉 That's It!

Most common use case:
```typescript
import { authManager } from './src/auth/auth-manager'

// Authenticate
await authManager.authenticate({ provider: 'anthropic', apiKey: '...' })

// Get recommendation
const rec = authManager.getRecommendations({ task: 'code_generation' })[0]

// Track usage
authManager.recordUsage(rec.provider, rec.model, 1000, 500)

// Show cost
console.log(authManager.getCostSummary())
```

Done! 🚀

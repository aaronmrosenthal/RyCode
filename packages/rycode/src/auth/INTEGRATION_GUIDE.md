# Integration Guide: Provider Authentication System

Complete guide for integrating the new provider authentication system into RyCode.

## 🚀 Quick Start

```typescript
import { authManager } from './auth/auth-manager'

// 1. Auto-detect and import existing credentials
const detected = await authManager.autoDetect()
if (detected.canImport) {
  const result = await authManager.importDetected(detected)
  console.log(`Imported ${result.success} credentials`)
}

// 2. Or authenticate manually
await authManager.authenticate({
  provider: 'anthropic',
  apiKey: process.env.ANTHROPIC_API_KEY!
})

// 3. Use the provider
const status = await authManager.getStatus('anthropic')
console.log(`Authenticated: ${status?.authenticated}`)
console.log(`Models: ${status?.models.join(', ')}`)
```

## 📦 Complete Integration Example

### 1. Application Initialization

```typescript
// src/app/init.ts
import { authManager } from './auth/auth-manager'

export async function initializeAuth() {
  console.log('🔐 Initializing authentication...')

  // Auto-detect existing credentials
  const detected = await authManager.autoDetect()

  if (detected.canImport && detected.found.length > 0) {
    console.log(detected.message)
    console.log('Found credentials for:', detected.found.map(d => d.provider).join(', '))

    // Ask user if they want to import
    const shouldImport = await askUser('Import detected credentials?')

    if (shouldImport) {
      const result = await authManager.importDetected(detected)
      console.log(`✅ Imported ${result.success} credentials`)

      if (result.failed > 0) {
        console.error(`❌ Failed to import ${result.failed} credentials`)
        result.errors.forEach(err => console.error('  -', err))
      }
    }
  } else {
    console.log('No existing credentials found')
    console.log('Use the model selector to sign into providers')
  }

  // Show authenticated providers
  const authenticated = await authManager.getAllStatus()
  if (authenticated.length > 0) {
    console.log('✅ Authenticated providers:')
    for (const status of authenticated) {
      console.log(`  - ${status.provider}: ${status.models.length} models`)
    }
  }
}
```

### 2. Provider Selection UI

```typescript
// src/ui/provider-selector.ts
import { authManager } from './auth/auth-manager'

export async function showProviderSelector() {
  // Get all available providers
  const providers = authManager.getAvailableProviders()

  // Get authentication status
  const authenticated = await authManager.getAllStatus()
  const authMap = new Map(authenticated.map(s => [s.provider, s]))

  // Display providers
  for (const provider of providers) {
    const status = authMap.get(provider.id)

    if (status?.authenticated) {
      console.log(`✅ ${provider.displayName}`)
      console.log(`   ${status.models.length} models available`)

      if (status.expiresAt) {
        const remaining = Math.floor((status.expiresAt.getTime() - Date.now()) / 1000 / 60)
        console.log(`   Expires in ${remaining} minutes`)
      }

      if (!status.healthy) {
        console.log(`   ⚠️  Provider temporarily unavailable`)
      }
    } else {
      console.log(`🔓 ${provider.displayName} - Sign In`)
    }
  }
}
```

### 3. Authentication Flow

```typescript
// src/ui/authenticate.ts
import { authManager } from './auth/auth-manager'
import { AuthenticationError } from './auth/errors'

export async function authenticateProvider(
  provider: string,
  credentials: { apiKey?: string; oauthToken?: string }
) {
  try {
    console.log(`🔐 Authenticating with ${provider}...`)

    const status = await authManager.authenticate(
      {
        provider,
        ...credentials
      },
      {
        saveCredentials: true,
        testConnection: true
      }
    )

    console.log(`✅ Successfully authenticated!`)
    console.log(`   Models: ${status.models.join(', ')}`)

    return status
  } catch (error) {
    if (error instanceof AuthenticationError) {
      // Show user-friendly error
      console.error(`❌ ${error.getUserMessage()}`)

      if (error.suggestedAction) {
        console.log(`💡 ${error.suggestedAction}`)
      }

      if (error.helpUrl) {
        console.log(`🔗 ${error.helpUrl}`)
      }

      // Log to audit
      // (already done automatically by authManager)
    } else {
      console.error('Unexpected error:', error)
    }

    throw error
  }
}
```

### 4. Model Selection with Recommendations

```typescript
// src/ui/model-selector.ts
import { authManager } from './auth/auth-manager'

export async function selectModel(context: {
  task: string
  complexity?: string
}) {
  // Get authenticated providers
  const authenticated = await authManager.getAllStatus()
  const availableProviders = authenticated
    .filter(s => s.authenticated && s.healthy)
    .map(s => s.provider)

  if (availableProviders.length === 0) {
    console.log('❌ No authenticated providers')
    console.log('💡 Sign in to a provider first')
    return null
  }

  // Get recommendations
  const recommendations = authManager.getRecommendations(
    {
      task: context.task as any,
      complexity: context.complexity as any,
      speedPreference: 'balanced',
      costPreference: 'balanced'
    },
    availableProviders
  )

  if (recommendations.length === 0) {
    console.log('No model recommendations available')
    return null
  }

  // Show top recommendation
  const top = recommendations[0]
  console.log(`🎯 Recommended: ${top.model}`)
  console.log(`   ${top.reason}`)
  console.log(`   Pros: ${top.pros.join(', ')}`)
  console.log(`   Cost: ${top.estimatedCost}`)
  console.log(`   Speed: ${top.speed} | Quality: ${'⭐'.repeat(top.quality)}`)

  return top
}
```

### 5. Usage Tracking

```typescript
// src/provider/request.ts
import { authManager } from './auth/auth-manager'

export async function makeRequest(
  provider: string,
  model: string,
  prompt: string
) {
  // Make API request
  const response = await callProviderAPI(provider, model, prompt)

  // Track usage for cost monitoring
  authManager.recordUsage(
    provider,
    model,
    response.usage.inputTokens,
    response.usage.outputTokens
  )

  return response
}
```

### 6. Cost Dashboard

```typescript
// src/ui/cost-dashboard.ts
import { authManager } from './auth/auth-manager'

export function showCostDashboard() {
  const summary = authManager.getCostSummary()

  console.log('💰 Cost Summary')
  console.log('─'.repeat(40))
  console.log(`Today:        $${summary.today.toFixed(2)}`)
  console.log(`This Week:    $${summary.thisWeek.toFixed(2)}`)
  console.log(`This Month:   $${summary.thisMonth.toFixed(2)}`)
  console.log(``)
  console.log('📊 Projection')
  console.log(`Daily Avg:    $${summary.projection.dailyAverage.toFixed(2)}`)
  console.log(`Monthly Est:  $${summary.projection.monthlyProjection.toFixed(2)}`)
  console.log(`Yearly Est:   $${summary.projection.yearlyProjection.toFixed(2)}`)

  // Show cost-saving tips
  const tips = authManager.getCostSavingTips()
  if (tips.length > 0) {
    console.log('')
    console.log('💡 Cost-Saving Tips')
    tips.forEach(tip => {
      console.log(`   ${tip.message}`)
      if (tip.potentialSaving > 0) {
        console.log(`   Save ~$${tip.potentialSaving.toFixed(2)}/month`)
      }
    })
  }
}
```

### 7. Status Bar Integration

```typescript
// src/ui/status-bar.ts
import { authManager } from './auth/auth-manager'

export function updateStatusBar(currentProvider: string, currentModel: string) {
  const summary = authManager.getCostSummary()
  const costToday = summary.today

  // Format: "Claude 3.5 Sonnet | ⚡ Fast | 💰 $0.12 today | [tab→]"
  const status = [
    currentModel,
    '⚡ Fast',
    `💰 $${costToday.toFixed(2)} today`,
    '[tab→]'
  ].join(' | ')

  console.log(status)
}
```

### 8. Health Monitoring

```typescript
// src/monitoring/health.ts
import { authManager } from './auth/auth-manager'

export async function checkAuthHealth() {
  const health = await authManager.healthCheck()

  if (!health.healthy) {
    console.warn('⚠️  Authentication issues detected:')
    health.issues.forEach(issue => console.warn(`  - ${issue}`))
  }

  // Check circuit breakers
  const unhealthy = authManager.getUnhealthyProviders()
  if (unhealthy.length > 0) {
    console.warn(`Circuit breakers open for: ${unhealthy.join(', ')}`)
  }

  return health
}

// Run health check every 5 minutes
setInterval(checkAuthHealth, 5 * 60 * 1000)
```

### 9. Audit and Security

```typescript
// src/admin/audit.ts
import { authManager } from './auth/auth-manager'

export function showAuditReport(provider?: string) {
  const summary = authManager.getAuditSummary(provider)

  console.log('📋 Audit Summary')
  console.log('─'.repeat(40))
  console.log(`Total Events: ${summary.totalEvents}`)
  console.log(`Success Rate: ${(summary.successRate * 100).toFixed(1)}%`)
  console.log(`Failure Rate: ${(summary.failureRate * 100).toFixed(1)}%`)

  console.log('\nBy Provider:')
  for (const [provider, count] of Object.entries(summary.byProvider)) {
    console.log(`  ${provider}: ${count} events`)
  }

  if (summary.recentFailures.length > 0) {
    console.log('\n⚠️  Recent Failures:')
    summary.recentFailures.forEach(event => {
      console.log(`  ${event.provider}: ${event.reason} (${event.timestamp})`)
    })
  }

  if (summary.riskEvents.length > 0) {
    console.log('\n🚨 High-Risk Events:')
    summary.riskEvents.forEach(event => {
      console.log(`  ${event.provider}: ${event.eventType} (risk: ${event.riskScore}/10)`)
    })
  }
}

// Check for suspicious activity
export function checkSecurity() {
  const providers = ['anthropic', 'openai', 'grok', 'qwen', 'google']

  for (const provider of providers) {
    const suspicious = authManager.detectSuspiciousActivity(provider)

    if (suspicious.suspicious) {
      console.error(`🚨 Suspicious activity detected for ${provider}:`)
      suspicious.reasons.forEach(reason => console.error(`  - ${reason}`))
    }
  }
}
```

### 10. Error Handling Best Practices

```typescript
// src/util/error-handler.ts
import { AuthenticationError, RateLimitError, NetworkError } from './auth/errors'

export async function handleAuthError(error: unknown) {
  if (error instanceof RateLimitError) {
    // Rate limited - wait and retry
    console.log('⏳ Rate limited, waiting...')
    const retryAfter = error.context.retryAfter || 60
    await sleep(retryAfter * 1000)
    return { retry: true }
  }

  if (error instanceof NetworkError) {
    // Network issue - might be temporary
    console.log('🌐 Network issue, will retry...')
    return { retry: true, delay: 5000 }
  }

  if (error instanceof AuthenticationError) {
    // Show user-friendly error
    console.error('❌', error.getUserMessage())

    if (error.suggestedAction) {
      console.log('💡', error.suggestedAction)
    }

    if (error.helpUrl) {
      console.log('🔗', error.helpUrl)
    }

    // Don't retry non-retryable errors
    return { retry: error.isRetryable() }
  }

  // Unknown error
  console.error('Unexpected error:', error)
  return { retry: false }
}
```

## 🔄 Migration from Agent System

### Before (Agent System)
```typescript
// Old way
const agent = selectAgent('build') // build, plan, doc
await agent.run(task)
```

### After (Provider System)
```typescript
// New way
const recommendation = await authManager.getRecommendations({
  task: 'code_generation',
  complexity: 'medium'
})

const model = recommendation[0]
await makeRequest(model.provider, model.model, task)
```

## 📊 Monitoring Integration

```typescript
// src/monitoring/metrics.ts
import { authManager } from './auth/auth-manager'

export function collectMetrics() {
  return {
    auth: {
      authenticated: authManager.getAllStatus(),
      costs: authManager.getCostSummary(),
      audit: authManager.getAuditSummary(),
      health: authManager.healthCheck(),
      circuitBreakers: authManager.getCircuitBreakerStats()
    }
  }
}

// Export for monitoring dashboard
export async function exportMonitoringData() {
  const data = await authManager.export()

  // Save to file or send to monitoring service
  await Bun.write('monitoring-data.json', JSON.stringify(data, null, 2))
}
```

## 🧪 Testing

```typescript
// test/auth.test.ts
import { authManager } from '../src/auth/auth-manager'

describe('Authentication', () => {
  test('authenticates with valid API key', async () => {
    const status = await authManager.authenticate({
      provider: 'anthropic',
      apiKey: process.env.ANTHROPIC_API_KEY!
    })

    expect(status.authenticated).toBe(true)
    expect(status.models.length).toBeGreaterThan(0)
  })

  test('handles invalid API key gracefully', async () => {
    await expect(
      authManager.authenticate({
        provider: 'anthropic',
        apiKey: 'invalid-key'
      })
    ).rejects.toThrow(AuthenticationError)
  })

  test('tracks costs', () => {
    authManager.recordUsage('anthropic', 'claude-3-5-sonnet-20241022', 1000, 500)

    const summary = authManager.getCostSummary()
    expect(summary.today).toBeGreaterThan(0)
  })
})
```

## 🎯 Best Practices

1. **Always use AuthManager** - Don't call providers directly
2. **Enable auto-detection** - Import existing credentials automatically
3. **Track costs** - Monitor spending in real-time
4. **Use recommendations** - Let the system suggest optimal models
5. **Handle errors gracefully** - Show user-friendly messages
6. **Monitor health** - Check circuit breakers and audit logs
7. **Respect rate limits** - The system handles this automatically
8. **Test credentials** - Use `testConnection: true` when authenticating
9. **Check expiry** - OAuth tokens expire, check status regularly
10. **Log security events** - Review audit logs for suspicious activity

## 🚀 Next Steps

1. Integrate with existing TUI components
2. Add model selector with inline authentication
3. Update status bar to show model and cost
4. Implement Tab key model cycling
5. Create migration wizard for existing users
6. Add comprehensive tests
7. Update documentation

---

**Status**: ✅ Ready for integration
**Documentation**: Complete
**Examples**: Comprehensive
**Best Practices**: Documented

#!/usr/bin/env bun
/**
 * CLI Provider Authentication Test
 *
 * This script tests that RyCode can detect and use CLI-authenticated providers
 * (Codex, Gemini, Claude) without making actual API calls. It verifies:
 * - Authentication status via CLI auth
 * - Provider configuration
 * - Model availability
 * - No actual API calls are made (dry-run mode)
 *
 * Usage:
 *   bun run packages/rycode/test/provider-cli-test.ts
 */

import { Auth } from '../src/auth'
import { ModelsDev } from '../src/provider/models'

interface CLIProviderTest {
  provider: string
  displayName: string
  expectedModels?: string[]
}

const CLI_PROVIDERS: CLIProviderTest[] = [
  {
    provider: 'openai',
    displayName: 'OpenAI Codex',
    expectedModels: ['gpt-4', 'gpt-3.5-turbo', 'gpt-4-turbo'],
  },
  {
    provider: 'google',
    displayName: 'Google Gemini',
    expectedModels: ['gemini-1.5-pro', 'gemini-1.5-flash'],
  },
  {
    provider: 'anthropic',
    displayName: 'Anthropic Claude',
    expectedModels: ['claude-3-5-sonnet-20241022', 'claude-3-5-haiku-20241022'],
  },
]

interface TestResult {
  provider: string
  status: 'authenticated' | 'not_authenticated' | 'error'
  authType?: string
  modelCount?: number
  models?: string[]
  message?: string
}

async function testCLIAuthentication(): Promise<TestResult[]> {
  const results: TestResult[] = []

  console.log('🔐 Testing CLI-Authenticated Providers...\n')

  try {
    const auth = await Auth.all()
    const authEntries = Object.entries(auth)

    console.log(`Found ${authEntries.length} authenticated provider(s) in CLI auth:\n`)

    for (const test of CLI_PROVIDERS) {
      const authConfig = auth[test.provider]

      if (authConfig) {
        console.log(`  ✓ ${test.displayName} (${test.provider})`)
        console.log(`    Auth Type: ${authConfig.type}`)

        if (authConfig.type === 'api' && authConfig.key) {
          const maskedKey = authConfig.key.slice(0, 8) + '...' + authConfig.key.slice(-4)
          console.log(`    API Key: ${maskedKey}`)
        } else if (authConfig.type === 'oauth') {
          console.log(`    OAuth: Token present`)
        }

        results.push({
          provider: test.provider,
          status: 'authenticated',
          authType: authConfig.type,
        })
      } else {
        console.log(`  ✗ ${test.displayName} (${test.provider})`)
        console.log(`    Status: Not authenticated`)
        results.push({
          provider: test.provider,
          status: 'not_authenticated',
          message: 'Run: rycode auth login',
        })
      }
      console.log('')
    }
  } catch (error: any) {
    console.error('❌ Error reading auth config:', error.message)
    for (const test of CLI_PROVIDERS) {
      results.push({
        provider: test.provider,
        status: 'error',
        message: error.message,
      })
    }
  }

  return results
}

async function testModelAvailability(results: TestResult[]): Promise<void> {
  console.log('\n📦 Testing Model Availability...\n')

  try {
    await ModelsDev.refresh()
    const providers = await ModelsDev.get()

    for (const result of results) {
      if (result.status !== 'authenticated') continue

      const test = CLI_PROVIDERS.find(t => t.provider === result.provider)
      if (!test) continue

      const providerData = providers[result.provider]
      if (!providerData) {
        console.log(`  ✗ ${test.displayName}: Provider not found in models`)
        continue
      }

      const models = providerData.models ? Object.keys(providerData.models) : []
      result.modelCount = models.length
      result.models = models

      console.log(`  ✓ ${test.displayName}:`)
      console.log(`    Total Models: ${models.length}`)

      // Check for expected models
      if (test.expectedModels && test.expectedModels.length > 0) {
        const foundExpected = test.expectedModels.filter(expected =>
          models.some(m => m.includes(expected) || expected.includes(m))
        )
        console.log(`    Expected Models Found: ${foundExpected.length}/${test.expectedModels.length}`)

        if (foundExpected.length > 0) {
          console.log(`    Available: ${foundExpected.slice(0, 3).join(', ')}${foundExpected.length > 3 ? '...' : ''}`)
        }
      } else {
        // Show first 3 models
        if (models.length > 0) {
          const sampleModels = models.slice(0, 3)
          console.log(`    Sample Models: ${sampleModels.join(', ')}${models.length > 3 ? '...' : ''}`)
        }
      }
      console.log('')
    }
  } catch (error: any) {
    console.error('❌ Error fetching models:', error.message)
  }
}

async function testProviderConfiguration(results: TestResult[]): Promise<void> {
  console.log('\n⚙️  Testing Provider Configuration...\n')

  for (const result of results) {
    if (result.status !== 'authenticated') continue

    const test = CLI_PROVIDERS.find(t => t.provider === result.provider)
    if (!test) continue

    console.log(`  ${test.displayName}:`)

    // Check if provider is ready for use
    const isReady = result.authType && result.modelCount && result.modelCount > 0

    if (isReady) {
      console.log(`    ✅ Ready to use`)
      console.log(`    Auth: ${result.authType}`)
      console.log(`    Models: ${result.modelCount} available`)
    } else {
      console.log(`    ⚠️  Configuration incomplete`)
      if (!result.authType) console.log(`    Missing: Authentication`)
      if (!result.modelCount || result.modelCount === 0) console.log(`    Missing: Model list`)
    }
    console.log('')
  }
}

async function generateTestReport(results: TestResult[]): Promise<void> {
  console.log('\n📊 Test Summary:\n')

  const authenticated = results.filter(r => r.status === 'authenticated')
  const notAuthenticated = results.filter(r => r.status === 'not_authenticated')
  const errors = results.filter(r => r.status === 'error')

  console.log(`  ✅ Authenticated: ${authenticated.length}/${CLI_PROVIDERS.length}`)
  console.log(`  ❌ Not Authenticated: ${notAuthenticated.length}/${CLI_PROVIDERS.length}`)
  if (errors.length > 0) {
    console.log(`  ⚠️  Errors: ${errors.length}`)
  }

  // Show which providers are ready
  const ready = authenticated.filter(r => r.modelCount && r.modelCount > 0)
  console.log(`\n  🚀 Ready to use: ${ready.length}/${CLI_PROVIDERS.length}`)

  if (ready.length > 0) {
    console.log(`\n  Available providers:`)
    for (const result of ready) {
      const test = CLI_PROVIDERS.find(t => t.provider === result.provider)
      if (test) {
        console.log(`    • ${test.displayName} (${result.modelCount} models)`)
      }
    }
  }

  if (notAuthenticated.length > 0) {
    console.log(`\n  💡 To authenticate missing providers:`)
    console.log(`     rycode auth login`)
    console.log(`     Or set environment variables (OPENAI_API_KEY, etc.)`)
  }

  console.log('')

  // Detailed status for each provider
  if (authenticated.length > 0) {
    console.log('  📝 Provider Details:\n')
    for (const result of results) {
      const test = CLI_PROVIDERS.find(t => t.provider === result.provider)
      if (!test) continue

      const statusEmoji = result.status === 'authenticated' ? '✅' : '❌'
      const modelInfo = result.modelCount ? ` (${result.modelCount} models)` : ''

      console.log(`    ${statusEmoji} ${test.displayName}${modelInfo}`)

      if (result.status === 'authenticated' && result.authType) {
        console.log(`       Auth: ${result.authType}`)
      }
      if (result.status === 'not_authenticated' && result.message) {
        console.log(`       ${result.message}`)
      }
    }
    console.log('')
  }
}

async function main() {
  console.log(`
╔════════════════════════════════════════════════╗
║   RyCode CLI Provider Authentication Test      ║
║   Testing Local/CLI Authenticated Providers    ║
╚════════════════════════════════════════════════╝
`)

  try {
    // Step 1: Check CLI authentication
    const results = await testCLIAuthentication()

    // Step 2: Test model availability
    await testModelAvailability(results)

    // Step 3: Test provider configuration
    await testProviderConfiguration(results)

    // Step 4: Generate report
    await generateTestReport(results)

    // Exit with appropriate code
    const hasAuthenticatedProvider = results.some(r => r.status === 'authenticated')
    const allReady = results.every(
      r => r.status === 'authenticated' && r.modelCount && r.modelCount > 0
    )

    if (allReady) {
      console.log('✅ All CLI providers are authenticated and ready!\n')
      process.exit(0)
    } else if (hasAuthenticatedProvider) {
      console.log('⚠️  Some providers are authenticated but not all are ready.\n')
      process.exit(0)
    } else {
      console.log('❌ No CLI providers are authenticated.\n')
      process.exit(1)
    }
  } catch (error: any) {
    console.error('\n❌ Test failed:', error.message)
    console.error(error.stack)
    process.exit(1)
  }
}

// Run if called directly
if (import.meta.main) {
  main()
}

export { testCLIAuthentication, testModelAvailability, CLI_PROVIDERS }

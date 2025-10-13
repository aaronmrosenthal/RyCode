#!/usr/bin/env bun
/**
 * Server Provider Test
 *
 * This script tests provider authentication by querying the running RyCode server.
 * It verifies that providers are configured and available through the API.
 *
 * Prerequisites:
 * - RyCode server must be running (bun run packages/rycode/src/index.ts serve --port 4096)
 *
 * Usage:
 *   bun run packages/rycode/test/provider-server-test.ts
 */

interface Provider {
  id: string
  name: string
  models: Record<string, Model>
}

interface Model {
  id: string
  name: string
  context?: number
  maxOutput?: number
}

interface ProvidersResponse {
  providers: Provider[]
  default: Record<string, number>
}

const SERVER_URL = process.env['RYCODE_SERVER'] || 'http://127.0.0.1:4096'

async function testServerConnection(): Promise<boolean> {
  console.log('🔌 Testing server connection...\n')

  try {
    // Try the providers endpoint directly (no /health endpoint)
    const response = await fetch(`${SERVER_URL}/app/providers`, {
      method: 'GET',
      signal: AbortSignal.timeout(5000),
    })

    if (response.ok) {
      console.log(`  ✅ Server is running at ${SERVER_URL}`)
      return true
    } else {
      console.log(`  ❌ Server returned ${response.status}`)
      return false
    }
  } catch (error: any) {
    console.log(`  ❌ Cannot connect to server`)
    console.log(`  Error: ${error.message}`)
    console.log(`\n  💡 Start the server with:`)
    console.log(`     bun run packages/rycode/src/index.ts serve --port 4096`)
    return false
  }
}

async function fetchProviders(): Promise<Provider[] | null> {
  console.log('\n📦 Fetching providers from server...\n')

  try {
    const response = await fetch(`${SERVER_URL}/app/providers`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      signal: AbortSignal.timeout(10000),
    })

    if (!response.ok) {
      console.log(`  ❌ Failed to fetch providers: ${response.status}`)
      return null
    }

    const data = (await response.json()) as ProvidersResponse
    return data.providers
  } catch (error: any) {
    console.log(`  ❌ Error fetching providers: ${error.message}`)
    return null
  }
}

function analyzeProviders(providers: Provider[]): void {
  console.log(`Found ${providers.length} provider(s):\n`)

  // Filter to our target providers
  const targetProviders = ['openai', 'anthropic', 'google']
  const authenticatedTargets = providers.filter(p => targetProviders.includes(p.id))

  if (authenticatedTargets.length === 0) {
    console.log('  ⚠️  None of the target providers (OpenAI, Anthropic, Google) are authenticated\n')
  }

  // Show all providers with model counts
  for (const provider of providers) {
    const modelCount = Object.keys(provider.models || {}).length
    const isTarget = targetProviders.includes(provider.id)
    const marker = isTarget ? '🎯' : '  '

    console.log(`  ${marker} ${provider.name} (${provider.id})`)
    console.log(`     Models: ${modelCount}`)

    // Show sample models for target providers
    if (isTarget && modelCount > 0) {
      const modelNames = Object.values(provider.models).slice(0, 3).map(m => m.name)
      console.log(`     Sample: ${modelNames.join(', ')}${modelCount > 3 ? '...' : ''}`)
    }

    console.log('')
  }
}

function generateReport(providers: Provider[]): void {
  console.log('📊 Provider Status Report:\n')

  const targetProviders = [
    { id: 'openai', name: 'OpenAI Codex' },
    { id: 'google', name: 'Google Gemini' },
    { id: 'anthropic', name: 'Anthropic Claude' },
  ]

  let authenticatedCount = 0
  let totalModels = 0

  for (const target of targetProviders) {
    const provider = providers.find(p => p.id === target.id)

    if (provider) {
      const modelCount = Object.keys(provider.models || {}).length
      authenticatedCount++
      totalModels += modelCount

      console.log(`  ✅ ${target.name}`)
      console.log(`     Status: Authenticated`)
      console.log(`     Models: ${modelCount} available`)

      // Show model examples
      if (modelCount > 0) {
        const sampleModels = Object.entries(provider.models)
          .slice(0, 2)
          .map(([id, model]) => model.name || id)
        console.log(`     Examples: ${sampleModels.join(', ')}`)
      }
    } else {
      console.log(`  ❌ ${target.name}`)
      console.log(`     Status: Not authenticated`)
    }
    console.log('')
  }

  console.log(`\n  Summary:`)
  console.log(`  • Authenticated: ${authenticatedCount}/3 target providers`)
  console.log(`  • Total models: ${totalModels}`)
  console.log(`  • All providers: ${providers.length}`)

  if (authenticatedCount === 0) {
    console.log(`\n  💡 To authenticate providers:`)
    console.log(`     1. Run: rycode auth login`)
    console.log(`     2. Or set: OPENAI_API_KEY=sk-...`)
    console.log(`     3. Restart the server`)
  } else if (authenticatedCount < 3) {
    console.log(`\n  💡 Some providers missing. To add more:`)
    console.log(`     rycode auth login`)
  } else {
    console.log(`\n  🚀 All target providers are ready!`)
  }
}

async function testModelSelection(providers: Provider[]): Promise<void> {
  const targetProviders = providers.filter(p => ['openai', 'anthropic', 'google'].includes(p.id))

  if (targetProviders.length === 0) {
    console.log('\n  ⚠️  No target providers available for model selection test\n')
    return
  }

  console.log('\n🎯 Testing Model Selection:\n')

  for (const provider of targetProviders) {
    const models = Object.values(provider.models || {})
    if (models.length === 0) continue

    console.log(`  ${provider.name}:`)

    // Find recommended models
    const recommended = models.filter(m =>
      m.name.toLowerCase().includes('gpt-4') ||
      m.name.toLowerCase().includes('claude-3.5') ||
      m.name.toLowerCase().includes('gemini-1.5')
    )

    if (recommended.length > 0) {
      console.log(`    Recommended models:`)
      for (const model of recommended.slice(0, 2)) {
        const contextInfo = model.context ? ` (${model.context} tokens)` : ''
        console.log(`      • ${model.name}${contextInfo}`)
      }
    } else {
      console.log(`    Available: ${models[0].name}`)
    }
    console.log('')
  }
}

async function main() {
  console.log(`
╔════════════════════════════════════════════════╗
║   RyCode Server Provider Test                  ║
║   Testing Provider Authentication via API      ║
╚════════════════════════════════════════════════╝
`)

  // Step 1: Check server connection
  const serverOnline = await testServerConnection()
  if (!serverOnline) {
    console.log('\n❌ Cannot proceed without server connection\n')
    process.exit(1)
  }

  // Step 2: Fetch providers
  const providers = await fetchProviders()
  if (!providers || providers.length === 0) {
    console.log('\n❌ No providers found\n')
    process.exit(1)
  }

  // Step 3: Analyze providers
  analyzeProviders(providers)

  // Step 4: Test model selection
  await testModelSelection(providers)

  // Step 5: Generate report
  generateReport(providers)

  // Determine exit code
  const targetProviders = ['openai', 'anthropic', 'google']
  const authenticatedTargets = providers.filter(p => targetProviders.includes(p.id))

  if (authenticatedTargets.length === 3) {
    console.log('\n✅ All target providers authenticated and ready!\n')
    process.exit(0)
  } else if (authenticatedTargets.length > 0) {
    console.log('\n⚠️  Some providers authenticated, others missing\n')
    process.exit(0)
  } else {
    console.log('\n❌ No target providers authenticated\n')
    process.exit(1)
  }
}

// Run if called directly
if (import.meta.main) {
  main().catch(error => {
    console.error('\n❌ Test failed:', error.message)
    process.exit(1)
  })
}

export { testServerConnection, fetchProviders, SERVER_URL }

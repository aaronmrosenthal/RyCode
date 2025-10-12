#!/usr/bin/env bun
/**
 * Provider Authentication Test Script
 *
 * This script tests that RyCode can detect and use API keys from environment variables
 * for all supported providers, ensuring developers can use their own accounts.
 */

import { SmartProviderSetup } from '../src/auth/auto-detect'
import { Auth } from '../src/auth'
import { ModelsDev } from '../src/provider/models'

interface TestResult {
  provider: string
  status: 'detected' | 'configured' | 'not_found' | 'error'
  source?: string
  message?: string
}

async function testProviderDetection(): Promise<TestResult[]> {
  const results: TestResult[] = []
  const smartSetup = new SmartProviderSetup()

  console.log('🔍 Testing Provider Auto-Detection...\n')

  // Test auto-detection
  const detected = await smartSetup.autoDetect()

  console.log(`Found ${detected.found.length} credential(s) in environment:\n`)

  for (const cred of detected.found) {
    console.log(`  ✓ ${cred.provider} (from ${cred.source}: ${cred.metadata?.['envVar'] || 'unknown'})`)
    results.push({
      provider: cred.provider,
      status: 'detected',
      source: `${cred.source}:${cred.metadata?.['envVar'] || 'unknown'}`,
    })
  }

  // Test configured credentials
  console.log('\n📋 Checking Configured Credentials...\n')

  const auth = await Auth.all()
  const authEntries = Object.entries(auth)

  if (authEntries.length > 0) {
    for (const [providerID, config] of authEntries) {
      console.log(`  ✓ ${providerID} (configured: ${config.type})`)
      results.push({
        provider: providerID,
        status: 'configured',
        source: `auth.json:${config.type}`,
      })
    }
  } else {
    console.log('  No credentials configured in auth.json')
  }

  // Check expected providers
  console.log('\n🎯 Expected Providers Status:\n')

  const expectedProviders = ['openai', 'anthropic', 'google', 'qwen']

  for (const provider of expectedProviders) {
    const found = results.find(r => r.provider === provider)
    if (found) {
      console.log(`  ✓ ${provider}: ${found.status} (${found.source})`)
    } else {
      console.log(`  ✗ ${provider}: not found`)
      results.push({
        provider,
        status: 'not_found',
        message: `Set ${provider.toUpperCase()}_API_KEY environment variable`,
      })
    }
  }

  return results
}

async function testModelAvailability() {
  console.log('\n\n📦 Testing Model Availability...\n')

  try {
    await ModelsDev.refresh()
    const providers = await ModelsDev.get()

    const providerList = Object.entries(providers)
      .filter(([id]) => !id.includes('opencode'))
      .slice(0, 10)

    console.log(`Found ${providerList.length} external providers:\n`)

    for (const [id, provider] of providerList) {
      const modelCount = provider.models ? Object.keys(provider.models).length : 0
      console.log(`  ${provider.name || id}: ${modelCount} model(s)`)
    }
  } catch (error) {
    console.error('Error fetching models:', error)
  }
}

async function generateTestReport(results: TestResult[]) {
  console.log('\n\n📊 Test Summary:\n')

  const detected = results.filter(r => r.status === 'detected').length
  const configured = results.filter(r => r.status === 'configured').length
  const notFound = results.filter(r => r.status === 'not_found').length

  console.log(`  ✓ Detected from environment: ${detected}`)
  console.log(`  ✓ Configured in auth.json: ${configured}`)
  console.log(`  ✗ Not found: ${notFound}`)

  if (notFound > 0) {
    console.log('\n💡 To add missing providers:')
    console.log('   1. Set environment variables (e.g., ANTHROPIC_API_KEY=sk-...)')
    console.log('   2. Or run: bun run packages/rycode/src/index.ts auth login')
    console.log('   3. Or add to .env file in project root')
  }

  console.log('\n✅ Test Complete!\n')
}

async function main() {
  console.log(`
╔════════════════════════════════════════╗
║   RyCode Provider Authentication Test  ║
║   Testing Developer Account Support    ║
╚════════════════════════════════════════╝
`)

  try {
    const results = await testProviderDetection()
    await testModelAvailability()
    await generateTestReport(results)

    // Exit with appropriate code
    const hasAnyProvider = results.some(r => r.status === 'detected' || r.status === 'configured')
    process.exit(hasAnyProvider ? 0 : 1)
  } catch (error) {
    console.error('\n❌ Test failed:', error)
    process.exit(1)
  }
}

// Run if called directly
if (import.meta.main) {
  main()
}

export { testProviderDetection, testModelAvailability }

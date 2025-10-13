#!/usr/bin/env bun
/**
 * Test CLI Bridge
 *
 * Tests that the CLI provider bridge can detect and communicate with
 * installed CLI tools (claude, qwen, codex, gemini).
 */

import { cliProviderBridge } from './src/auth/providers/cli-bridge'

async function main() {
  console.log('🔍 Testing CLI Provider Bridge\n')

  // Test 1: Detect available CLI providers
  console.log('1️⃣ Detecting available CLI providers...')
  const providers = await cliProviderBridge.detectAvailableProviders()

  if (providers.length === 0) {
    console.log('   ❌ No CLI providers found')
    process.exit(1)
  }

  console.log(`   ✅ Found ${providers.length} CLI provider(s):\n`)
  for (const provider of providers) {
    console.log(`   • ${provider.name}`)
    console.log(`     Path: ${provider.cliCommand}`)
    console.log(`     Available: ${provider.available ? '✅' : '❌'}`)
    if (provider.version) {
      console.log(`     Version: ${provider.version}`)
    }
    console.log()
  }

  // Test 2: Get providers with models
  console.log('2️⃣ Getting providers with available models...')
  const providersWithModels = await cliProviderBridge.getAvailableProvidersWithModels()

  console.log(`   ✅ Found ${providersWithModels.length} provider(s) with models:\n`)
  for (const provider of providersWithModels) {
    console.log(`   • ${provider.provider} (via ${provider.source})`)
    console.log(`     Models: ${provider.models.length}`)
    console.log(`     Sample: ${provider.models.slice(0, 3).join(', ')}`)
    console.log()
  }

  // Test 3: Test provider communication (simple test)
  console.log('3️⃣ Testing provider communication...')
  for (const provider of providers) {
    try {
      console.log(`   Testing ${provider.name}...`)
      const isWorking = await cliProviderBridge.testProvider(provider.name)
      if (isWorking) {
        console.log(`   ✅ ${provider.name} is working`)
      } else {
        console.log(`   ⚠️  ${provider.name} test returned false`)
      }
    } catch (error: any) {
      console.log(`   ❌ ${provider.name} error: ${error.message}`)
    }
    console.log()
  }

  console.log('✅ CLI Bridge tests complete!\n')
}

main().catch(error => {
  console.error('❌ Test failed:', error.message)
  process.exit(1)
})

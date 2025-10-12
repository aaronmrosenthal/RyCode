#!/usr/bin/env bun
/**
 * End-to-End Provider Test
 *
 * This script makes actual API calls to verify that providers work correctly
 * with user-provided API keys. It tests basic completion requests to ensure
 * authentication and API integration is working.
 *
 * Usage:
 *   # Test all providers with environment variables
 *   OPENAI_API_KEY=sk-... ANTHROPIC_API_KEY=sk-ant-... bun run packages/rycode/test/provider-e2e-test.ts
 *
 *   # Test specific provider
 *   OPENAI_API_KEY=sk-... bun run packages/rycode/test/provider-e2e-test.ts openai
 */

import { createOpenAI } from '@ai-sdk/openai'
import { createAnthropic } from '@ai-sdk/anthropic'
import { createGoogleGenerativeAI } from '@ai-sdk/google'
import { generateText } from 'ai'

interface ProviderTest {
  name: string
  provider: string
  envVar: string
  modelId: string
  test: () => Promise<{ success: boolean; message: string; latency?: number }>
}

const TEST_PROMPT = 'Write a one-line comment explaining what a const is in JavaScript.'

const providers: ProviderTest[] = [
  {
    name: 'OpenAI (GPT-3.5)',
    provider: 'openai',
    envVar: 'OPENAI_API_KEY',
    modelId: 'gpt-3.5-turbo',
    test: async () => {
      const apiKey = process.env['OPENAI_API_KEY']
      if (!apiKey) {
        return { success: false, message: 'OPENAI_API_KEY not set' }
      }

      try {
        const startTime = Date.now()
        const openai = createOpenAI({ apiKey })
        const result = await generateText({
          model: openai('gpt-3.5-turbo'),
          prompt: TEST_PROMPT,
          maxOutputTokens: 50,
        })
        const latency = Date.now() - startTime

        if (result.text && result.text.length > 10) {
          return {
            success: true,
            message: `Received ${result.text.length} chars`,
            latency,
          }
        }
        return { success: false, message: 'Empty response' }
      } catch (error: any) {
        return {
          success: false,
          message: `API Error: ${error.message || 'Unknown error'}`,
        }
      }
    },
  },
  {
    name: 'Anthropic (Claude)',
    provider: 'anthropic',
    envVar: 'ANTHROPIC_API_KEY',
    modelId: 'claude-3-haiku-20240307',
    test: async () => {
      const apiKey = process.env['ANTHROPIC_API_KEY'] || process.env['CLAUDE_API_KEY']
      if (!apiKey) {
        return { success: false, message: 'ANTHROPIC_API_KEY not set' }
      }

      try {
        const startTime = Date.now()
        const anthropic = createAnthropic({ apiKey })
        const result = await generateText({
          model: anthropic('claude-3-haiku-20240307'),
          prompt: TEST_PROMPT,
          maxOutputTokens: 50,
        })
        const latency = Date.now() - startTime

        if (result.text && result.text.length > 10) {
          return {
            success: true,
            message: `Received ${result.text.length} chars`,
            latency,
          }
        }
        return { success: false, message: 'Empty response' }
      } catch (error: any) {
        return {
          success: false,
          message: `API Error: ${error.message || 'Unknown error'}`,
        }
      }
    },
  },
  {
    name: 'Google (Gemini)',
    provider: 'google',
    envVar: 'GOOGLE_API_KEY',
    modelId: 'gemini-1.5-flash',
    test: async () => {
      const apiKey = process.env['GOOGLE_API_KEY']
      if (!apiKey) {
        return { success: false, message: 'GOOGLE_API_KEY not set' }
      }

      try {
        const startTime = Date.now()
        const google = createGoogleGenerativeAI({ apiKey })
        const result = await generateText({
          model: google('gemini-1.5-flash'),
          prompt: TEST_PROMPT,
          maxOutputTokens: 50,
        })
        const latency = Date.now() - startTime

        if (result.text && result.text.length > 10) {
          return {
            success: true,
            message: `Received ${result.text.length} chars`,
            latency,
          }
        }
        return { success: false, message: 'Empty response' }
      } catch (error: any) {
        return {
          success: false,
          message: `API Error: ${error.message || 'Unknown error'}`,
        }
      }
    },
  },
]

async function runTest(test: ProviderTest): Promise<void> {
  const apiKey = process.env[test.envVar]
  const hasKey = !!apiKey

  process.stdout.write(`  Testing ${test.name}... `)

  if (!hasKey) {
    console.log(`⏭️  SKIPPED (${test.envVar} not set)`)
    return
  }

  try {
    const result = await test.test()

    if (result.success) {
      const latencyStr = result.latency ? ` (${result.latency}ms)` : ''
      console.log(`✅ PASS - ${result.message}${latencyStr}`)
    } else {
      console.log(`❌ FAIL - ${result.message}`)
    }
  } catch (error: any) {
    console.log(`❌ ERROR - ${error.message || 'Unknown error'}`)
  }
}

async function runAllTests(specificProvider?: string) {
  console.log(`
╔════════════════════════════════════════════╗
║   RyCode Provider End-to-End Test          ║
║   Testing actual API calls with user keys  ║
╚════════════════════════════════════════════╝
`)

  let testsToRun = providers
  if (specificProvider) {
    testsToRun = providers.filter(p => p.provider === specificProvider)
    if (testsToRun.length === 0) {
      console.error(`❌ Unknown provider: ${specificProvider}`)
      console.error(`Available providers: ${providers.map(p => p.provider).join(', ')}`)
      process.exit(1)
    }
  }

  console.log('🧪 Running provider tests...\n')

  for (const test of testsToRun) {
    await runTest(test)
  }

  console.log('\n📊 Test Summary:')
  const configured = providers.filter(p => process.env[p.envVar]).length
  const total = providers.length

  console.log(`  ✓ Configured: ${configured}/${total}`)
  console.log(`  ⏭️  Skipped: ${total - configured}/${total}`)

  if (configured === 0) {
    console.log('\n💡 No API keys configured. Set environment variables:')
    for (const test of providers) {
      console.log(`   export ${test.envVar}="your-key-here"`)
    }
    process.exit(1)
  }

  console.log('\n✅ Test Complete!\n')
}

// Parse CLI arguments
const specificProvider = process.argv[2]

// Run tests
if (import.meta.main) {
  runAllTests(specificProvider).catch(error => {
    console.error('\n❌ Test suite failed:', error)
    process.exit(1)
  })
}

export { runAllTests, providers }

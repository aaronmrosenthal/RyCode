/**
 * 1-Click Smart Setup - Auto-detect existing credentials
 *
 * Automatically detects API keys from environment variables, config files,
 * and CLI tools to provide a seamless onboarding experience.
 */

import { existsSync, readFileSync } from 'fs'
import { homedir } from 'os'
import { join } from 'path'
import { exec } from 'child_process'
import { promisify } from 'util'

const execAsync = promisify(exec)

export interface DetectedCredential {
  provider: string
  source: 'env' | 'config' | 'cli' | 'keychain'
  credential: string
  metadata?: Record<string, any>
}

export interface AutoDetectResult {
  found: DetectedCredential[]
  message: string
  canImport: boolean
}

export class SmartProviderSetup {
  /**
   * Auto-detect all available credentials
   */
  async autoDetect(): Promise<AutoDetectResult> {
    const found: DetectedCredential[] = []

    // Check environment variables
    const envCreds = await this.checkEnvironmentVariables()
    found.push(...envCreds)

    // Check existing config files
    const configCreds = await this.checkConfigFiles()
    found.push(...configCreds)

    // Check CLI tools
    const cliCreds = await this.checkCLITools()
    found.push(...cliCreds)

    // Generate appropriate message
    const message = this.generateMessage(found)

    return {
      found,
      message,
      canImport: found.length > 0
    }
  }

  /**
   * Check environment variables for API keys
   */
  private async checkEnvironmentVariables(): Promise<DetectedCredential[]> {
    const found: DetectedCredential[] = []

    const envVars = [
      { key: 'ANTHROPIC_API_KEY', provider: 'anthropic' },
      { key: 'CLAUDE_API_KEY', provider: 'anthropic' },
      { key: 'OPENAI_API_KEY', provider: 'openai' },
      { key: 'XAI_API_KEY', provider: 'grok' },
      { key: 'GROK_API_KEY', provider: 'grok' },
      { key: 'DASHSCOPE_API_KEY', provider: 'qwen' },
      { key: 'QWEN_API_KEY', provider: 'qwen' },
      { key: 'GOOGLE_API_KEY', provider: 'google' },
      { key: 'GOOGLE_APPLICATION_CREDENTIALS', provider: 'google' }
    ]

    for (const { key, provider } of envVars) {
      const value = process.env[key]
      if (value && value.length > 10) {
        // Check if it's a file path (for Google)
        if (key === 'GOOGLE_APPLICATION_CREDENTIALS' && existsSync(value)) {
          found.push({
            provider,
            source: 'env',
            credential: value,
            metadata: { type: 'service_account_file', envVar: key }
          })
        } else if (key !== 'GOOGLE_APPLICATION_CREDENTIALS') {
          found.push({
            provider,
            source: 'env',
            credential: value,
            metadata: { envVar: key }
          })
        }
      }
    }

    return found
  }

  /**
   * Check common config file locations
   */
  private async checkConfigFiles(): Promise<DetectedCredential[]> {
    const found: DetectedCredential[] = []
    const home = homedir()

    const configPaths = [
      // Anthropic
      { path: join(home, '.anthropic', 'config.json'), provider: 'anthropic', key: 'api_key' },
      { path: join(home, '.config', 'anthropic', 'config.json'), provider: 'anthropic', key: 'api_key' },

      // OpenAI
      { path: join(home, '.openai', 'config.json'), provider: 'openai', key: 'api_key' },
      { path: join(home, '.config', 'openai', 'config.json'), provider: 'openai', key: 'api_key' },

      // Grok / xAI
      { path: join(home, '.xai', 'config.json'), provider: 'grok', key: 'api_key' },
      { path: join(home, '.config', 'xai', 'config.json'), provider: 'grok', key: 'api_key' },

      // Qwen
      { path: join(home, '.dashscope', 'config.json'), provider: 'qwen', key: 'api_key' },

      // Google
      { path: join(home, '.config', 'gcloud', 'application_default_credentials.json'), provider: 'google', key: 'token' }
    ]

    for (const { path, provider, key } of configPaths) {
      try {
        if (existsSync(path)) {
          const content = readFileSync(path, 'utf-8')
          const config = JSON.parse(content)

          if (config[key]) {
            found.push({
              provider,
              source: 'config',
              credential: config[key],
              metadata: { configPath: path }
            })
          }
        }
      } catch (error) {
        // Ignore parsing errors
      }
    }

    return found
  }

  /**
   * Check CLI tools for authentication
   */
  private async checkCLITools(): Promise<DetectedCredential[]> {
    const found: DetectedCredential[] = []

    // Check for installed CLI tools using 'which'
    const cliTools = [
      { name: 'claude', provider: 'anthropic', versionCmd: 'claude --version 2>/dev/null' },
      { name: 'qwen', provider: 'qwen', versionCmd: 'qwen --version 2>/dev/null' },
      { name: 'codex', provider: 'openai', versionCmd: 'codex --version 2>/dev/null' },
      { name: 'gemini', provider: 'google', versionCmd: 'gemini --version 2>/dev/null' },
    ]

    for (const { name, provider, versionCmd } of cliTools) {
      try {
        // Check if CLI tool exists
        const { stdout: whichOut } = await execAsync(`which ${name} 2>/dev/null`)
        if (whichOut.trim()) {
          // CLI tool exists, try to get version to confirm it works
          try {
            const { stdout: versionOut } = await execAsync(versionCmd)
            found.push({
              provider,
              source: 'cli',
              credential: 'cli-authenticated', // Marker that CLI is available
              metadata: {
                tool: name,
                path: whichOut.trim(),
                version: versionOut.trim()
              }
            })
          } catch {
            // Version check failed, but CLI exists - still add it
            found.push({
              provider,
              source: 'cli',
              credential: 'cli-authenticated',
              metadata: {
                tool: name,
                path: whichOut.trim()
              }
            })
          }
        }
      } catch {
        // CLI tool not found
      }
    }

    // Also check Google Cloud CLI (gcloud) as alternative to gemini CLI
    try {
      const { stdout } = await execAsync('gcloud auth print-access-token 2>/dev/null')
      if (stdout.trim()) {
        // Only add gcloud if we haven't already found gemini CLI
        const hasGeminiCLI = found.some(c => c.metadata?.['tool'] === 'gemini')
        if (!hasGeminiCLI) {
          found.push({
            provider: 'google',
            source: 'cli',
            credential: stdout.trim(),
            metadata: { tool: 'gcloud' }
          })
        }
      }
    } catch {
      // gcloud not available or not authenticated
    }

    return found
  }

  /**
   * Generate user-friendly message based on findings
   */
  private generateMessage(found: DetectedCredential[]): string {
    if (found.length === 0) {
      return "👋 Let's get you started! We recommend Anthropic for the best experience."
    }

    const providers = [...new Set(found.map(c => c.provider))]

    if (providers.length === 1) {
      return `🎉 Found existing ${this.getProviderName(providers[0])} credentials! Import them?`
    }

    const providerList = providers
      .map(p => this.getProviderName(p))
      .join(', ')

    return `🎉 Found existing credentials for: ${providerList}! Import them all?`
  }

  /**
   * Get friendly provider name
   */
  private getProviderName(provider: string): string {
    const names: Record<string, string> = {
      anthropic: 'Claude (Anthropic)',
      openai: 'OpenAI',
      grok: 'Grok (xAI)',
      qwen: 'Qwen (Alibaba)',
      google: 'Google AI'
    }
    return names[provider] || provider
  }

  /**
   * Import detected credentials
   */
  async importAll(
    detected: DetectedCredential[],
    storeFunction: (provider: string, credential: string) => Promise<void>
  ): Promise<{ success: number; failed: number; errors: string[] }> {
    const results = { success: 0, failed: 0, errors: [] as string[] }

    for (const cred of detected) {
      try {
        await storeFunction(cred.provider, cred.credential)
        results.success++
      } catch (error) {
        results.failed++
        results.errors.push(
          `Failed to import ${cred.provider}: ${error instanceof Error ? error.message : 'Unknown error'}`
        )
      }
    }

    return results
  }

  /**
   * Get quick start recommendation
   */
  getQuickStartRecommendation(): {
    provider: string
    reason: string
    helpUrl: string
  } {
    return {
      provider: 'anthropic',
      reason: 'Claude offers the best balance of capabilities, speed, and cost for coding tasks',
      helpUrl: 'https://console.anthropic.com/settings/keys'
    }
  }

  /**
   * Generate onboarding UI data
   */
  async generateOnboardingUI(): Promise<{
    hasExisting: boolean
    providers: string[]
    message: string
    actions: Array<{ label: string; action: string; primary: boolean }>
  }> {
    const detected = await this.autoDetect()

    if (detected.canImport) {
      return {
        hasExisting: true,
        providers: [...new Set(detected.found.map(c => c.provider))],
        message: detected.message,
        actions: [
          { label: '✨ Import Everything', action: 'import_all', primary: true },
          { label: '🚀 Quick Setup', action: 'quick_setup', primary: false },
          { label: 'Skip', action: 'skip', primary: false }
        ]
      }
    }

    return {
      hasExisting: false,
      providers: [],
      message: "👋 Let's get you started!",
      actions: [
        { label: '🚀 Quick Setup (Recommended)', action: 'quick_setup', primary: true },
        { label: 'Manual Setup', action: 'manual', primary: false },
        { label: 'Skip for now', action: 'skip', primary: false }
      ]
    }
  }
}

/**
 * Global instance
 */
export const smartSetup = new SmartProviderSetup()

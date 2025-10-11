/**
 * Provider Registry
 *
 * Central registry for all AI providers, implementing the Strategy pattern.
 * Provides unified interface for authentication across all providers.
 */

import { anthropicProvider } from './providers/anthropic'
import { openaiProvider } from './providers/openai'
import { grokProvider } from './providers/grok'
import { qwenProvider } from './providers/qwen'
import { googleProvider } from './providers/google'
import { AuthenticationError } from './errors'

export interface ProviderAuthConfig {
  provider: string
  [key: string]: any
}

export interface ProviderAuthResult {
  success: boolean
  provider: string
  method: string
  models: string[]
  expiresAt?: Date
  [key: string]: any
}

export interface ProviderInfo {
  id: string
  name: string
  displayName: string
  supportsAPIKey: boolean
  supportsOAuth: boolean
  supportsCLI: boolean
  requiresProject?: boolean
  credentialExpiry: boolean
  helpUrl: string
  models: Array<{
    id: string
    name: string
    contextWindow: number
    supportsVision?: boolean
    inputPrice: number
    outputPrice: number
  }>
}

/**
 * Provider Registry - manages all AI providers
 */
export class ProviderRegistry {
  private providers = new Map<string, any>([
    ['anthropic', anthropicProvider],
    ['openai', openaiProvider],
    ['grok', grokProvider],
    ['qwen', qwenProvider],
    ['google', googleProvider]
  ])

  /**
   * Get all registered providers
   */
  getProviders(): string[] {
    return Array.from(this.providers.keys())
  }

  /**
   * Get provider info
   */
  getProviderInfo(providerId: string): ProviderInfo | null {
    const provider = this.providers.get(providerId)
    if (!provider) return null

    const capabilities = provider.getCapabilities()

    const displayNames: Record<string, string> = {
      anthropic: 'Claude (Anthropic)',
      openai: 'OpenAI',
      grok: 'Grok (xAI)',
      qwen: 'Qwen (Alibaba)',
      google: 'Google AI'
    }

    return {
      id: providerId,
      name: providerId,
      displayName: displayNames[providerId] || providerId,
      supportsAPIKey: capabilities.supportsAPIKey,
      supportsOAuth: capabilities.supportsOAuth || false,
      supportsCLI: capabilities.supportsCLI || false,
      requiresProject: capabilities.requiresProject,
      credentialExpiry: capabilities.credentialExpiry,
      helpUrl: capabilities.helpUrl,
      models: capabilities.models
    }
  }

  /**
   * Get all provider info
   */
  getAllProviderInfo(): ProviderInfo[] {
    return this.getProviders()
      .map(id => this.getProviderInfo(id))
      .filter((info): info is ProviderInfo => info !== null)
  }

  /**
   * Authenticate with a provider
   */
  async authenticate(config: ProviderAuthConfig): Promise<ProviderAuthResult> {
    const { provider: providerId, ...providerConfig } = config

    const provider = this.providers.get(providerId)

    if (!provider) {
      throw new AuthenticationError(
        'validation_failed',
        false,
        { provider: providerId, timestamp: new Date() },
        `Unknown provider: ${providerId}`,
        undefined,
        `Available providers: ${this.getProviders().join(', ')}`
      )
    }

    // Call provider's authenticate method
    const result = await provider.authenticate(providerConfig)

    return {
      ...result,
      provider: providerId
    }
  }

  /**
   * Test credentials for a provider
   */
  async testCredentials(providerId: string, credentials: any): Promise<boolean> {
    const provider = this.providers.get(providerId)

    if (!provider) {
      return false
    }

    try {
      if (provider.testCredentials) {
        return await provider.testCredentials(credentials)
      }

      // Fallback: try to authenticate
      await provider.authenticate(credentials)
      return true
    } catch {
      return false
    }
  }

  /**
   * Get available models for a provider
   */
  async getAvailableModels(providerId: string, credentials: any): Promise<string[]> {
    const provider = this.providers.get(providerId)

    if (!provider) {
      throw new AuthenticationError(
        'validation_failed',
        false,
        { provider: providerId, timestamp: new Date() },
        `Unknown provider: ${providerId}`
      )
    }

    if (provider.getAvailableModels) {
      return await provider.getAvailableModels(credentials)
    }

    // Fallback: return known models
    const info = this.getProviderInfo(providerId)
    return info?.models.map(m => m.id) || []
  }

  /**
   * Get recommended provider for new users
   */
  getRecommendedProvider(): ProviderInfo {
    return this.getProviderInfo('anthropic')!
  }

  /**
   * Check if provider is registered
   */
  hasProvider(providerId: string): boolean {
    return this.providers.has(providerId)
  }

  /**
   * Get provider by ID
   */
  getProvider(providerId: string): any {
    return this.providers.get(providerId)
  }
}

/**
 * Global provider registry instance
 */
export const providerRegistry = new ProviderRegistry()

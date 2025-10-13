/**
 * Unified Auth Manager
 *
 * High-level API that combines all authentication components:
 * - Provider authentication
 * - Credential storage
 * - Audit logging
 * - Cost tracking
 * - Smart recommendations
 * - Auto-detection
 *
 * This is the main entry point for authentication operations.
 */

import { providerRegistry } from './provider-registry'
import { credentialStore } from './storage/credential-store'
import { auditLog } from './storage/audit-log'
import { costTracker } from './cost-tracker'
import { modelRecommender } from './model-recommender'
import { smartSetup } from './auto-detect'
import { circuitBreakerRegistry } from './security/circuit-breaker'
import { cliProviderBridge } from './providers/cli-bridge'
import type { ProviderAuthConfig, ProviderInfo } from './provider-registry'
import type { TaskContext, ModelRecommendation } from './model-recommender'
import type { CostSummary, CostSavingTip } from './cost-tracker'
import type { AutoDetectResult } from './auto-detect'
import { AuthenticationError } from './errors'

export interface AuthManagerConfig {
  autoLog?: boolean // Automatically log to audit (default: true)
  autoTrackCost?: boolean // Automatically track costs (default: true)
}

export interface AuthenticateOptions {
  saveCredentials?: boolean // Save to storage (default: true)
  testConnection?: boolean // Test after auth (default: false)
}

export interface AuthStatus {
  authenticated: boolean
  provider: string
  method: string
  expiresAt?: Date
  models: string[]
  healthy: boolean
}

export class AuthManager {
  private config: AuthManagerConfig

  constructor(config: AuthManagerConfig = {}) {
    this.config = {
      autoLog: config.autoLog ?? true,
      autoTrackCost: config.autoTrackCost ?? true
    }
  }

  /**
   * Authenticate with a provider
   */
  async authenticate(
    config: ProviderAuthConfig,
    options: AuthenticateOptions = {}
  ): Promise<AuthStatus> {
    const { provider } = config
    const {
      saveCredentials = true,
      testConnection = false
    } = options

    // Log attempt
    if (this.config.autoLog) {
      await auditLog.recordAuthAttempt(
        provider,
        config['method'] || 'api-key',
        { saveCredentials, testConnection }
      )
    }

    try {
      // Authenticate with provider
      const result = await providerRegistry.authenticate(config)

      // Log success
      if (this.config.autoLog) {
        await auditLog.recordAuthSuccess(
          provider,
          result.method,
          { modelsCount: result.models.length }
        )
      }

      // Save credentials
      if (saveCredentials) {
        await credentialStore.store(provider, result)

        if (this.config.autoLog) {
          await auditLog.recordCredentialStored(provider, result.method)
        }
      }

      // Test connection if requested
      if (testConnection) {
        await this.testCredentials(provider)
      }

      // Check circuit breaker health
      const healthy = circuitBreakerRegistry.isProviderHealthy(provider)

      return {
        authenticated: true,
        provider,
        method: result.method,
        expiresAt: result.expiresAt,
        models: result.models,
        healthy
      }
    } catch (error) {
      // Log failure
      if (this.config.autoLog && error instanceof AuthenticationError) {
        await auditLog.recordAuthFailure(
          provider,
          config['method'] || 'api-key',
          error.reason,
          { error: error.message }
        )
      }

      throw error
    }
  }

  /**
   * Get authentication status for a provider
   */
  async getStatus(provider: string): Promise<AuthStatus | null> {
    try {
      // Check if credentials exist
      const hasCredentials = await credentialStore.has(provider)

      if (!hasCredentials) {
        return null
      }

      // Get stored credential
      const credential = await credentialStore.retrieve(provider)

      if (!credential) {
        return null
      }

      // Check if expired
      const expired = await credentialStore.isExpired(provider)

      if (expired) {
        return {
          authenticated: false,
          provider,
          method: credential.method,
          expiresAt: credential.expiresAt,
          models: [],
          healthy: false
        }
      }

      // Get available models
      const providerInfo = providerRegistry.getProviderInfo(provider)
      const models = providerInfo?.models.map(m => m.id) || []

      // Check circuit breaker health
      const healthy = circuitBreakerRegistry.isProviderHealthy(provider)

      return {
        authenticated: true,
        provider,
        method: credential.method,
        expiresAt: credential.expiresAt,
        models,
        healthy
      }
    } catch {
      return null
    }
  }

  /**
   * Get all authenticated providers
   */
  async getAllStatus(): Promise<AuthStatus[]> {
    const providers = await credentialStore.list()
    const statuses: AuthStatus[] = []

    for (const provider of providers) {
      const status = await this.getStatus(provider)
      if (status) {
        statuses.push(status)
      }
    }

    return statuses
  }

  /**
   * Test stored credentials
   */
  async testCredentials(provider: string): Promise<boolean> {
    try {
      const authInfo = await credentialStore.getAuthInfo(provider)

      if (!authInfo) {
        return false
      }

      // Build config for testing
      const config: any = { provider }

      if (authInfo.type === 'api') {
        config.apiKey = authInfo.key
      } else if (authInfo.type === 'oauth') {
        config.method = 'oauth'
        config.oauthToken = authInfo.access
        config.refreshToken = authInfo.refresh
      }

      // Test with provider
      return await providerRegistry.testCredentials(provider, config)
    } catch {
      return false
    }
  }

  /**
   * Remove provider credentials
   */
  async signOut(provider: string): Promise<boolean> {
    const removed = await credentialStore.remove(provider)

    if (removed && this.config.autoLog) {
      await auditLog.recordCredentialRemoved(provider)
    }

    return removed
  }

  /**
   * Auto-detect existing credentials
   */
  async autoDetect(): Promise<AutoDetectResult> {
    return await smartSetup.autoDetect()
  }

  /**
   * Detect available CLI providers
   */
  async detectCLIProviders() {
    return await cliProviderBridge.detectAvailableProviders()
  }

  /**
   * Get available providers with models (including CLI)
   */
  async getAvailableProvidersWithModels() {
    return await cliProviderBridge.getAvailableProvidersWithModels()
  }

  /**
   * Test if a CLI provider is working
   */
  async testCLIProvider(provider: string): Promise<boolean> {
    return await cliProviderBridge.testProvider(provider)
  }

  /**
   * Import detected credentials
   */
  async importDetected(detected: AutoDetectResult): Promise<{
    success: number
    failed: number
    errors: string[]
  }> {
    return await smartSetup.importAll(
      detected.found,
      async (provider, credential) => {
        await this.authenticate({
          provider,
          apiKey: credential
        })
      }
    )
  }

  /**
   * Get model recommendations
   */
  getRecommendations(
    context: TaskContext,
    providers?: string[]
  ): ModelRecommendation[] {
    // Get available models from authenticated providers
    const availableModels: Array<{ provider: string; model: string }> = []

    for (const provider of providers || this.getAuthenticatedProviders()) {
      const info = providerRegistry.getProviderInfo(provider)
      if (info) {
        for (const model of info.models) {
          availableModels.push({ provider, model: model.id })
        }
      }
    }

    return modelRecommender.recommend(context, availableModels)
  }

  /**
   * Record usage and track cost
   */
  recordUsage(
    provider: string,
    model: string,
    inputTokens: number,
    outputTokens: number
  ): void {
    if (this.config.autoTrackCost) {
      costTracker.recordUsage(provider, model, inputTokens, outputTokens)
    }
  }

  /**
   * Get cost summary
   */
  getCostSummary(): CostSummary {
    return costTracker.getCostSummary()
  }

  /**
   * Get cost-saving tips
   */
  getCostSavingTips(): CostSavingTip[] {
    return costTracker.getCostSavingTips()
  }

  /**
   * Get audit summary
   */
  getAuditSummary(provider?: string) {
    return auditLog.getSummary(provider)
  }

  /**
   * Get recent audit events
   */
  getRecentAuditEvents(limit: number = 10) {
    return auditLog.getRecent(limit)
  }

  /**
   * Detect suspicious activity
   */
  detectSuspiciousActivity(provider: string) {
    return auditLog.detectSuspiciousActivity(provider)
  }

  /**
   * Get all available providers
   */
  getAvailableProviders(): ProviderInfo[] {
    return providerRegistry.getAllProviderInfo()
  }

  /**
   * Get authenticated providers (synchronous helper)
   */
  private getAuthenticatedProviders(): string[] {
    // This is a sync helper, but credentialStore.list() is async
    // In practice, this would need to be called after initialization
    // For now, return all providers
    return ['anthropic', 'openai', 'grok', 'qwen', 'google']
  }

  /**
   * Get circuit breaker stats
   */
  getCircuitBreakerStats() {
    return circuitBreakerRegistry.getAllStats()
  }

  /**
   * Get unhealthy providers
   */
  getUnhealthyProviders(): string[] {
    return circuitBreakerRegistry.getUnhealthyProviders()
  }

  /**
   * Reset circuit breaker for a provider
   */
  resetCircuitBreaker(provider: string): void {
    circuitBreakerRegistry.reset(provider)
  }

  /**
   * Export all data (for backup/debugging)
   */
  async export() {
    return {
      credentials: await credentialStore.export(),
      audit: auditLog.export(),
      costs: costTracker.exportData(),
      circuitBreakers: Object.fromEntries(this.getCircuitBreakerStats())
    }
  }

  /**
   * Health check
   */
  async healthCheck(): Promise<{
    healthy: boolean
    providers: Record<string, boolean>
    issues: string[]
  }> {
    const issues: string[] = []
    const providers: Record<string, boolean> = {}

    // Check authenticated providers
    const authenticated = await this.getAllStatus()

    for (const status of authenticated) {
      providers[status.provider] = status.healthy && status.authenticated

      if (!status.healthy) {
        issues.push(`${status.provider}: Circuit breaker is open`)
      }

      if (!status.authenticated) {
        issues.push(`${status.provider}: Not authenticated`)
      }

      if (status.expiresAt && status.expiresAt < new Date()) {
        issues.push(`${status.provider}: Credentials expired`)
      }
    }

    // Check for suspicious activity
    for (const status of authenticated) {
      const suspicious = this.detectSuspiciousActivity(status.provider)
      if (suspicious.suspicious) {
        issues.push(`${status.provider}: ${suspicious.reasons.join(', ')}`)
      }
    }

    return {
      healthy: issues.length === 0,
      providers,
      issues
    }
  }
}

/**
 * Global auth manager instance
 */
export const authManager = new AuthManager()

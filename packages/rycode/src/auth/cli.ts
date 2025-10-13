#!/usr/bin/env bun
/**
 * CLI interface for the authentication system
 *
 * This script provides a command-line interface to the auth manager,
 * allowing Go code to interact with the TypeScript authentication system.
 *
 * Usage:
 *   bun run cli.ts check <provider>
 *   bun run cli.ts auth <provider> <apiKey>
 *   bun run cli.ts cost
 *   bun run cli.ts health <provider>
 *   bun run cli.ts list
 *   bun run cli.ts auto-detect
 */

import { authManager } from './auth-manager'

interface CLIResponse {
  success: boolean
  data?: any
  error?: string
}

async function main() {
  const command = process.argv[2]
  const args = process.argv.slice(3)

  try {
    let response: CLIResponse

    switch (command) {
      case 'check': {
        // Check authentication status for a provider
        const provider = args[0]
        if (!provider) {
          throw new Error('Provider name required')
        }

        // First check stored credentials
        const status = await authManager.getStatus(provider)
        let isAuthenticated = status !== null && status.authenticated
        let models = status?.models.length || 0

        // If not authenticated via stored credentials, check CLI providers
        if (!isAuthenticated) {
          const cliProviders = await authManager.getAvailableProvidersWithModels()
          const cliProvider = cliProviders.find(p => p.provider === provider)
          if (cliProvider) {
            isAuthenticated = true
            models = cliProvider.models.length
          }
        }

        response = {
          success: true,
          data: {
            isAuthenticated,
            provider,
            modelsCount: models
          }
        }
        break
      }

      case 'auth': {
        // Authenticate with a provider
        const provider = args[0]
        const apiKey = args[1]

        if (!provider || !apiKey) {
          throw new Error('Provider and API key required')
        }

        const result = await authManager.authenticate({
          provider: provider as any,
          apiKey
        })

        response = {
          success: true,
          data: {
            provider: result.provider,
            modelsCount: result.models.length,
            message: `Successfully authenticated with ${provider}`
          }
        }
        break
      }

      case 'cost': {
        // Get cost summary
        const summary = authManager.getCostSummary()
        const savingsTips = authManager.getCostSavingTips()

        response = {
          success: true,
          data: {
            todayCost: summary.today,
            monthCost: summary.thisMonth,
            projection: summary.projection,
            savingsTip: savingsTips.length > 0 ? savingsTips[0].message : undefined
          }
        }
        break
      }

      case 'health': {
        // Get provider health status
        const provider = args[0]
        if (!provider) {
          throw new Error('Provider name required')
        }

        const stats = authManager.getCircuitBreakerStats()
        const providerStats = stats.get(provider)

        if (!providerStats) {
          response = {
            success: true,
            data: {
              provider,
              status: 'unknown',
              failureCount: 0
            }
          }
        } else {
          const status = providerStats.state === 'closed' ? 'healthy' :
                        providerStats.state === 'half-open' ? 'degraded' : 'down'

          response = {
            success: true,
            data: {
              provider,
              status,
              failureCount: providerStats.failures,
              nextAttemptAt: providerStats.nextAttempt?.toISOString()
            }
          }
        }
        break
      }

      case 'list': {
        // List all authenticated providers
        const authenticated = await authManager.getAllStatus()

        response = {
          success: true,
          data: {
            providers: authenticated.map(status => ({
              id: status.provider,
              name: status.provider,
              modelsCount: status.models.length
            }))
          }
        }
        break
      }

      case 'auto-detect': {
        // Auto-detect credentials
        const result = await authManager.autoDetect()

        response = {
          success: true,
          data: {
            message: result.message,
            found: Object.keys(result.found).length,
            credentials: Object.entries(result.found).map(([provider, keys]) => ({
              provider,
              count: Array.isArray(keys) ? keys.length : 1
            }))
          }
        }
        break
      }

      case 'recommendations': {
        // Get model recommendations
        const task = args[0] || 'general'
        const recommendations = authManager.getRecommendations({
          task: task as any
        })

        response = {
          success: true,
          data: {
            recommendations: recommendations.slice(0, 3).map(r => ({
              provider: r.provider,
              model: r.model,
              confidence: r.confidence,
              reason: r.reason
            }))
          }
        }
        break
      }

      case 'cli-providers': {
        // Get CLI providers with models
        const cliProviders = await authManager.getAvailableProvidersWithModels()

        response = {
          success: true,
          data: {
            providers: cliProviders.map(p => ({
              provider: p.provider,
              models: p.models,
              source: p.source
            }))
          }
        }
        break
      }

      default:
        throw new Error(`Unknown command: ${command}\n\nAvailable commands:\n  - check <provider>\n  - auth <provider> <apiKey>\n  - cost\n  - health <provider>\n  - list\n  - auto-detect\n  - recommendations [task]\n  - cli-providers`)
    }

    // Output JSON response
    console.log(JSON.stringify(response.data))
    process.exit(0)

  } catch (error) {
    // Output error as JSON
    const errorResponse: CLIResponse = {
      success: false,
      error: error instanceof Error ? error.message : String(error)
    }

    console.error(JSON.stringify(errorResponse))
    process.exit(1)
  }
}

// Run CLI
main().catch(err => {
  console.error(JSON.stringify({
    success: false,
    error: err.message
  }))
  process.exit(1)
})

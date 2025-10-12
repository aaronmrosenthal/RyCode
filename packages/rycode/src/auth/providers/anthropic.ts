/**
 * Anthropic Provider Authentication
 *
 * Handles authentication for Claude models using API keys.
 * Includes rate limiting, input validation, and circuit breaker protection.
 */

import { authRateLimiter } from '../security/rate-limiter'
import { inputValidator } from '../security/input-validator'
import { circuitBreakerRegistry } from '../security/circuit-breaker'
import {
  InvalidAPIKeyError,
  NetworkError,
  ErrorHandler,
  AuthenticationError
} from '../errors'

export interface AnthropicAuthConfig {
  apiKey: string
}

export interface AnthropicAuthResult {
  success: boolean
  apiKey: string
  models: string[]
  expiresAt?: Date
}

export interface AnthropicModel {
  id: string
  display_name: string
  created_at: string
}

const PROVIDER = 'anthropic'
const API_BASE_URL = 'https://api.anthropic.com/v1'
const API_VERSION = '2023-06-01'

export class AnthropicProvider {
  /**
   * Authenticate with Anthropic API
   */
  async authenticate(config: AnthropicAuthConfig): Promise<AnthropicAuthResult> {
    const { apiKey } = config

    try {
      // Step 1: Rate limiting check
      const rateLimitResult = await authRateLimiter.checkLimit(`${PROVIDER}:${apiKey.substring(0, 10)}`)
      if (!rateLimitResult.allowed) {
        throw new AuthenticationError(
          'rate_limited',
          true,
          { provider: PROVIDER, timestamp: new Date() },
          rateLimitResult.retryAfter
            ? `Taking a quick breather! Try again in ${rateLimitResult.retryAfter} seconds. ☕`
            : 'Rate limit exceeded',
          undefined,
          'Pro tip: Batch your requests for better flow'
        )
      }

      // Step 2: Input validation and sanitization
      const sanitized = inputValidator.sanitizeAPIKey(apiKey)
      const validationResult = await inputValidator.validateForStorage(PROVIDER, {
        apiKey: sanitized
      })

      if (!validationResult.valid) {
        throw new AuthenticationError(
          'validation_failed',
          false,
          { provider: PROVIDER, timestamp: new Date() },
          validationResult.error || 'Validation failed',
          validationResult.helpUrl,
          validationResult.hint
        )
      }

      // Step 3: Verify API key with circuit breaker protection
      const models = await circuitBreakerRegistry.call(PROVIDER, async () => {
        return await this.verifyAPIKey(sanitized)
      })

      // Step 4: Record success
      authRateLimiter.recordSuccess(`${PROVIDER}:${apiKey.substring(0, 10)}`)

      return {
        success: true,
        apiKey: sanitized,
        models,
        expiresAt: undefined // API keys don't expire
      }
    } catch (error) {
      // Handle and transform error
      const authError = ErrorHandler.handle(PROVIDER, error)
      ErrorHandler.log(authError)
      throw authError
    }
  }

  /**
   * Verify API key by attempting to list models
   */
  private async verifyAPIKey(apiKey: string): Promise<string[]> {
    try {
      const response = await fetch(`${API_BASE_URL}/models`, {
        method: 'GET',
        headers: {
          'x-api-key': apiKey,
          'anthropic-version': API_VERSION,
          'content-type': 'application/json'
        },
        signal: AbortSignal.timeout(30000) // 30 second timeout
      })

      if (!response.ok) {
        throw this.handleHTTPError(response.status, await response.text())
      }

      const data = await response.json()

      // Extract model IDs
      const models = data.data?.map((model: AnthropicModel) => model.id) || []

      // Fallback to known models if API doesn't return them
      if (models.length === 0) {
        return [
          'claude-sonnet-4-5-20250929',
          'claude-opus-4-1-20250805',
          'claude-sonnet-4-20250514',
          'claude-3-7-sonnet-20250219',
          'claude-3-5-sonnet-20241022',
          'claude-3-5-haiku-20241022'
        ]
      }

      return models
    } catch (error) {
      if (error instanceof AuthenticationError) {
        throw error
      }

      // Network errors
      if (error instanceof Error) {
        if (error.name === 'AbortError' || error.message.includes('timeout')) {
          throw new AuthenticationError(
            'timeout',
            true,
            { provider: PROVIDER, timestamp: new Date(), originalError: error },
            'Request to Anthropic timed out',
            undefined,
            'Try again in a moment'
          )
        }

        throw new NetworkError(PROVIDER, { originalError: error })
      }

      throw error
    }
  }

  /**
   * Handle HTTP errors from Anthropic API
   */
  private handleHTTPError(status: number, body: string): AuthenticationError {
    let errorData: any = {}
    try {
      errorData = JSON.parse(body)
    } catch {
      // Ignore JSON parse errors
    }

    const context = {
      provider: PROVIDER,
      statusCode: status,
      timestamp: new Date()
    }

    switch (status) {
      case 401:
        return new InvalidAPIKeyError(PROVIDER, context)

      case 403:
        return new AuthenticationError(
          'forbidden',
          false,
          context,
          'This API key doesn\'t have permission to access models',
          'https://console.anthropic.com/settings/keys',
          'Check your API key permissions in the Anthropic Console'
        )

      case 429:
        const retryAfter = errorData.retry_after || 60
        return new AuthenticationError(
          'rate_limited',
          true,
          context,
          `Anthropic rate limit reached. Try again in ${retryAfter} seconds.`,
          undefined,
          'Consider upgrading your Anthropic plan for higher limits'
        )

      case 500:
      case 502:
      case 503:
      case 504:
        return new AuthenticationError(
          'server_error',
          true,
          context,
          'Anthropic is experiencing issues. Please try again later.',
          'https://status.anthropic.com',
          'Check Anthropic status page'
        )

      default:
        return new AuthenticationError(
          'unknown',
          true,
          context,
          `Anthropic returned error ${status}: ${errorData.error?.message || 'Unknown error'}`,
          'https://console.anthropic.com/settings/keys'
        )
    }
  }

  /**
   * Test if credentials are still valid
   */
  async testCredentials(apiKey: string): Promise<boolean> {
    try {
      await this.verifyAPIKey(apiKey)
      return true
    } catch {
      return false
    }
  }

  /**
   * Get available models for authenticated user
   */
  async getAvailableModels(apiKey: string): Promise<string[]> {
    return await circuitBreakerRegistry.call(PROVIDER, async () => {
      return await this.verifyAPIKey(apiKey)
    })
  }

  /**
   * Get provider capabilities
   */
  getCapabilities() {
    return {
      supportsAPIKey: true,
      supportsOAuth: false,
      supportsCLI: false,
      requiresProject: false,
      credentialExpiry: false,
      helpUrl: 'https://console.anthropic.com/settings/keys',
      models: [
        {
          id: 'claude-sonnet-4-5-20250929',
          name: 'Claude Sonnet 4.5',
          contextWindow: 200000,
          supportsVision: true,
          inputPrice: 0.003,
          outputPrice: 0.015
        },
        {
          id: 'claude-opus-4-1-20250805',
          name: 'Claude Opus 4.1',
          contextWindow: 200000,
          supportsVision: true,
          inputPrice: 0.015,
          outputPrice: 0.075
        },
        {
          id: 'claude-sonnet-4-20250514',
          name: 'Claude Sonnet 4',
          contextWindow: 200000,
          supportsVision: true,
          inputPrice: 0.003,
          outputPrice: 0.015
        },
        {
          id: 'claude-3-7-sonnet-20250219',
          name: 'Claude 3.7 Sonnet',
          contextWindow: 200000,
          supportsVision: true,
          inputPrice: 0.003,
          outputPrice: 0.015
        },
        {
          id: 'claude-3-5-sonnet-20241022',
          name: 'Claude 3.5 Sonnet',
          contextWindow: 200000,
          supportsVision: true,
          inputPrice: 0.003,
          outputPrice: 0.015
        },
        {
          id: 'claude-3-5-haiku-20241022',
          name: 'Claude 3.5 Haiku',
          contextWindow: 200000,
          supportsVision: false,
          inputPrice: 0.001,
          outputPrice: 0.005
        }
      ]
    }
  }
}

/**
 * Global Anthropic provider instance
 */
export const anthropicProvider = new AnthropicProvider()

/**
 * Grok (xAI) Provider Authentication
 *
 * Handles authentication for Grok models using API keys.
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

export interface GrokAuthConfig {
  apiKey: string
}

export interface GrokAuthResult {
  success: boolean
  apiKey: string
  models: string[]
  expiresAt?: Date
}

const PROVIDER = 'grok'
const API_BASE_URL = 'https://api.x.ai/v1'

export class GrokProvider {
  /**
   * Authenticate with Grok (xAI) API
   */
  async authenticate(config: GrokAuthConfig): Promise<GrokAuthResult> {
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
          'Authorization': `Bearer ${apiKey}`,
          'Content-Type': 'application/json'
        },
        signal: AbortSignal.timeout(30000) // 30 second timeout
      })

      if (!response.ok) {
        throw this.handleHTTPError(response.status, await response.text())
      }

      const data = await response.json()

      // Extract Grok model IDs
      const models = data.data?.map((model: any) => model.id) || []

      // Fallback to known models if API doesn't return them
      if (models.length === 0) {
        return [
          'grok-4-fast',
          'grok-4',
          'grok-3',
          'grok-2-1212',
          'grok-2-vision-1212'
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
            'Request to Grok (xAI) timed out',
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
   * Handle HTTP errors from Grok API
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

    const errorMessage = errorData.error?.message || errorData.message || 'Unknown error'

    switch (status) {
      case 401:
        return new InvalidAPIKeyError(PROVIDER, context)

      case 403:
        return new AuthenticationError(
          'forbidden',
          false,
          context,
          'This API key doesn\'t have permission to access Grok models',
          'https://console.x.ai/api-keys',
          'Check your API key permissions in the xAI Console'
        )

      case 429:
        const retryAfter = errorData.retry_after || 60
        return new AuthenticationError(
          'rate_limited',
          true,
          context,
          `Grok rate limit reached. Try again in ${retryAfter} seconds.`,
          undefined,
          'Consider upgrading your xAI plan for higher limits'
        )

      case 500:
      case 502:
      case 503:
      case 504:
        return new AuthenticationError(
          'server_error',
          true,
          context,
          'Grok (xAI) is experiencing issues. Please try again later.',
          'https://status.x.ai',
          'Check xAI status page'
        )

      default:
        return new AuthenticationError(
          'unknown',
          true,
          context,
          `Grok returned error ${status}: ${errorMessage}`,
          'https://console.x.ai/api-keys'
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
      helpUrl: 'https://console.x.ai/api-keys',
      specialFeatures: ['real-time-web-search', 'humor', 'x-twitter-context'],
      models: [
        {
          id: 'grok-4-fast',
          name: 'Grok 4 Fast',
          contextWindow: 2000000,
          supportsVision: false,
          supportsRealTime: true,
          inputPrice: 0.0004,
          outputPrice: 0.002
        },
        {
          id: 'grok-4',
          name: 'Grok 4',
          contextWindow: 128000,
          supportsVision: false,
          supportsRealTime: true,
          inputPrice: 0.002,
          outputPrice: 0.01
        },
        {
          id: 'grok-3',
          name: 'Grok 3',
          contextWindow: 128000,
          supportsVision: false,
          supportsRealTime: true,
          inputPrice: 0.002,
          outputPrice: 0.01
        },
        {
          id: 'grok-2-1212',
          name: 'Grok 2',
          contextWindow: 128000,
          supportsVision: false,
          supportsRealTime: true,
          inputPrice: 0.002,
          outputPrice: 0.01
        },
        {
          id: 'grok-2-vision-1212',
          name: 'Grok 2 Vision',
          contextWindow: 128000,
          supportsVision: true,
          supportsRealTime: true,
          inputPrice: 0.002,
          outputPrice: 0.01
        }
      ]
    }
  }
}

/**
 * Global Grok provider instance
 */
export const grokProvider = new GrokProvider()

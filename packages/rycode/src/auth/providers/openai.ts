/**
 * OpenAI Provider Authentication
 *
 * Handles authentication for GPT models using API keys.
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

export interface OpenAIAuthConfig {
  apiKey: string
  organization?: string
}

export interface OpenAIAuthResult {
  success: boolean
  apiKey: string
  organization?: string
  models: string[]
  expiresAt?: Date
}

export interface OpenAIModel {
  id: string
  object: string
  created: number
  owned_by: string
}

const PROVIDER = 'openai'
const API_BASE_URL = 'https://api.openai.com/v1'

export class OpenAIProvider {
  /**
   * Authenticate with OpenAI API
   */
  async authenticate(config: OpenAIAuthConfig): Promise<OpenAIAuthResult> {
    const { apiKey, organization } = config

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
        return await this.verifyAPIKey(sanitized, organization)
      })

      // Step 4: Record success
      authRateLimiter.recordSuccess(`${PROVIDER}:${apiKey.substring(0, 10)}`)

      return {
        success: true,
        apiKey: sanitized,
        organization,
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
  private async verifyAPIKey(apiKey: string, organization?: string): Promise<string[]> {
    try {
      const headers: Record<string, string> = {
        'Authorization': `Bearer ${apiKey}`,
        'Content-Type': 'application/json'
      }

      if (organization) {
        headers['OpenAI-Organization'] = organization
      }

      const response = await fetch(`${API_BASE_URL}/models`, {
        method: 'GET',
        headers,
        signal: AbortSignal.timeout(30000) // 30 second timeout
      })

      if (!response.ok) {
        throw this.handleHTTPError(response.status, await response.text())
      }

      const data = await response.json()

      // Extract GPT model IDs (filter for chat models)
      const allModels = data.data?.map((model: OpenAIModel) => model.id) || []
      const gptModels = allModels.filter((id: string) =>
        id.startsWith('gpt-') && !id.includes('instruct') && !id.includes('vision')
      )

      // Fallback to known models if API doesn't return them
      if (gptModels.length === 0) {
        return [
          'gpt-5',
          'gpt-4.1',
          'gpt-4.1-mini',
          'gpt-4o',
          'gpt-4-turbo',
          'gpt-3.5-turbo'
        ]
      }

      return gptModels
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
            'Request to OpenAI timed out',
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
   * Handle HTTP errors from OpenAI API
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

    const errorMessage = errorData.error?.message || 'Unknown error'
    const errorType = errorData.error?.type

    switch (status) {
      case 401:
        if (errorMessage.includes('Incorrect API key')) {
          return new InvalidAPIKeyError(PROVIDER, context)
        }
        return new AuthenticationError(
          'invalid_key',
          false,
          context,
          errorMessage,
          'https://platform.openai.com/api-keys',
          'Check your API key is correct and active'
        )

      case 403:
        return new AuthenticationError(
          'forbidden',
          false,
          context,
          errorMessage.includes('country')
            ? 'OpenAI is not available in your country'
            : 'This API key doesn\'t have permission to access models',
          'https://platform.openai.com/account/api-keys',
          errorMessage.includes('country')
            ? 'Consider using a VPN or alternative provider'
            : 'Check your API key permissions'
        )

      case 429:
        const isRateLimit = errorType === 'rate_limit_exceeded'
        const isQuotaExceeded = errorMessage.includes('quota') || errorMessage.includes('billing')

        if (isQuotaExceeded) {
          return new AuthenticationError(
            'forbidden',
            false,
            context,
            'Your OpenAI account has exceeded its quota or billing limit',
            'https://platform.openai.com/account/billing',
            'Add payment method or upgrade your plan'
          )
        }

        return new AuthenticationError(
          'rate_limited',
          true,
          context,
          isRateLimit
            ? 'OpenAI rate limit reached. Try again in a moment.'
            : errorMessage,
          undefined,
          'Consider upgrading your OpenAI plan for higher limits'
        )

      case 500:
      case 502:
      case 503:
      case 504:
        return new AuthenticationError(
          'server_error',
          true,
          context,
          'OpenAI is experiencing issues. Please try again later.',
          'https://status.openai.com',
          'Check OpenAI status page'
        )

      default:
        return new AuthenticationError(
          'unknown',
          true,
          context,
          `OpenAI returned error ${status}: ${errorMessage}`,
          'https://platform.openai.com/api-keys'
        )
    }
  }

  /**
   * Test if credentials are still valid
   */
  async testCredentials(apiKey: string, organization?: string): Promise<boolean> {
    try {
      await this.verifyAPIKey(apiKey, organization)
      return true
    } catch {
      return false
    }
  }

  /**
   * Get available models for authenticated user
   */
  async getAvailableModels(apiKey: string, organization?: string): Promise<string[]> {
    return await circuitBreakerRegistry.call(PROVIDER, async () => {
      return await this.verifyAPIKey(apiKey, organization)
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
      requiresOrganization: true,
      credentialExpiry: false,
      helpUrl: 'https://platform.openai.com/api-keys',
      models: [
        {
          id: 'gpt-5',
          name: 'GPT-5',
          contextWindow: 1000000,
          supportsVision: true,
          inputPrice: 0.01,
          outputPrice: 0.03
        },
        {
          id: 'gpt-4.1',
          name: 'GPT-4.1',
          contextWindow: 1000000,
          supportsVision: true,
          inputPrice: 0.008,
          outputPrice: 0.024
        },
        {
          id: 'gpt-4.1-mini',
          name: 'GPT-4.1 Mini',
          contextWindow: 1000000,
          supportsVision: true,
          inputPrice: 0.0015,
          outputPrice: 0.006
        },
        {
          id: 'gpt-4o',
          name: 'GPT-4o',
          contextWindow: 128000,
          supportsVision: true,
          inputPrice: 0.005,
          outputPrice: 0.015
        },
        {
          id: 'gpt-4-turbo',
          name: 'GPT-4 Turbo',
          contextWindow: 128000,
          supportsVision: true,
          inputPrice: 0.01,
          outputPrice: 0.03
        },
        {
          id: 'gpt-3.5-turbo',
          name: 'GPT-3.5 Turbo',
          contextWindow: 16385,
          supportsVision: false,
          inputPrice: 0.0005,
          outputPrice: 0.0015
        }
      ]
    }
  }
}

/**
 * Global OpenAI provider instance
 */
export const openaiProvider = new OpenAIProvider()

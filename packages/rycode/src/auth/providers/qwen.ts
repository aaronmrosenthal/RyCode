/**
 * Qwen (Alibaba Cloud) Provider Authentication
 *
 * Handles authentication for Qwen models using API keys from DashScope.
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

export interface QwenAuthConfig {
  apiKey: string
}

export interface QwenAuthResult {
  success: boolean
  apiKey: string
  models: string[]
  expiresAt?: Date
}

const PROVIDER = 'qwen'
const API_BASE_URL = 'https://dashscope.aliyuncs.com/api/v1'

export class QwenProvider {
  /**
   * Authenticate with Qwen (DashScope) API
   */
  async authenticate(config: QwenAuthConfig): Promise<QwenAuthResult> {
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
   * Verify API key by making a simple API call
   */
  private async verifyAPIKey(apiKey: string): Promise<string[]> {
    try {
      // DashScope uses a different endpoint structure
      // We'll use the text generation endpoint with a minimal request to verify
      const response = await fetch(`${API_BASE_URL}/services/aigc/text-generation/generation`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${apiKey}`,
          'Content-Type': 'application/json',
          'X-DashScope-SSE': 'disable'
        },
        body: JSON.stringify({
          model: 'qwen-turbo',
          input: {
            messages: [
              { role: 'user', content: 'test' }
            ]
          },
          parameters: {
            max_tokens: 1
          }
        }),
        signal: AbortSignal.timeout(30000) // 30 second timeout
      })

      if (!response.ok) {
        throw this.handleHTTPError(response.status, await response.text())
      }

      // If we get here, the API key is valid
      // Return known Qwen models
      return this.getKnownModels()
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
            'Request to Qwen (DashScope) timed out',
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
   * Get known Qwen models
   */
  private getKnownModels(): string[] {
    return [
      'qwen3-coder-480b',
      'qwen3-coder-30b',
      'qwen2.5-coder',
      'qwen-max',
      'qwen-plus',
      'qwen-turbo'
    ]
  }

  /**
   * Handle HTTP errors from Qwen API
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

    const errorMessage = errorData.message || errorData.error?.message || 'Unknown error'
    const errorCode = errorData.code || errorData.error?.code

    switch (status) {
      case 401:
        return new InvalidAPIKeyError(PROVIDER, context)

      case 403:
        return new AuthenticationError(
          'forbidden',
          false,
          context,
          errorCode === 'InsufficientBalance'
            ? 'Your Qwen account has insufficient balance'
            : 'This API key doesn\'t have permission to access Qwen models',
          'https://dashscope.console.aliyun.com/',
          errorCode === 'InsufficientBalance'
            ? 'Add funds to your DashScope account'
            : 'Check your API key permissions in DashScope Console'
        )

      case 429:
        const retryAfter = errorData.retry_after || 60
        return new AuthenticationError(
          'rate_limited',
          true,
          context,
          `Qwen rate limit reached. Try again in ${retryAfter} seconds.`,
          undefined,
          'Consider upgrading your DashScope plan for higher limits'
        )

      case 500:
      case 502:
      case 503:
      case 504:
        return new AuthenticationError(
          'server_error',
          true,
          context,
          'Qwen (DashScope) is experiencing issues. Please try again later.',
          undefined,
          'Check DashScope status or try again later'
        )

      default:
        return new AuthenticationError(
          'unknown',
          true,
          context,
          `Qwen returned error ${status}: ${errorMessage}`,
          'https://dashscope.console.aliyun.com/'
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
      helpUrl: 'https://dashscope.console.aliyun.com/',
      models: [
        {
          id: 'qwen3-coder-480b',
          name: 'Qwen3-Coder 480B',
          contextWindow: 256000,
          supportsVision: false,
          inputPrice: 0.003,
          outputPrice: 0.009
        },
        {
          id: 'qwen3-coder-30b',
          name: 'Qwen3-Coder 30B',
          contextWindow: 256000,
          supportsVision: false,
          inputPrice: 0.001,
          outputPrice: 0.003
        },
        {
          id: 'qwen2.5-coder',
          name: 'Qwen2.5-Coder',
          contextWindow: 128000,
          supportsVision: false,
          inputPrice: 0.0008,
          outputPrice: 0.0024
        },
        {
          id: 'qwen-max',
          name: 'Qwen Max',
          contextWindow: 8192,
          supportsVision: false,
          inputPrice: 0.002,
          outputPrice: 0.006
        },
        {
          id: 'qwen-plus',
          name: 'Qwen Plus',
          contextWindow: 32768,
          supportsVision: false,
          inputPrice: 0.0004,
          outputPrice: 0.0012
        },
        {
          id: 'qwen-turbo',
          name: 'Qwen Turbo',
          contextWindow: 8192,
          supportsVision: false,
          inputPrice: 0.0002,
          outputPrice: 0.0006
        }
      ]
    }
  }
}

/**
 * Global Qwen provider instance
 */
export const qwenProvider = new QwenProvider()

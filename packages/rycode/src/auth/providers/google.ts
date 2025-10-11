/**
 * Google AI Provider Authentication
 *
 * Handles authentication for Gemini models using OAuth2 or API keys.
 * Includes rate limiting, input validation, circuit breaker protection,
 * and CSRF protection for OAuth flows.
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
import { randomBytes } from 'crypto'

export interface GoogleAuthConfig {
  method: 'oauth' | 'api-key' | 'cli'
  apiKey?: string
  oauthToken?: string
  refreshToken?: string
  projectId?: string
}

export interface GoogleAuthResult {
  success: boolean
  method: 'oauth' | 'api-key' | 'cli'
  apiKey?: string
  oauthToken?: string
  refreshToken?: string
  projectId?: string
  models: string[]
  expiresAt?: Date
}

export interface CSRFToken {
  token: string
  expiresAt: Date
}

const PROVIDER = 'google'
const API_BASE_URL = 'https://generativelanguage.googleapis.com/v1beta'
const OAUTH_TOKEN_URL = 'https://oauth2.googleapis.com/token'
const OAUTH_REVOKE_URL = 'https://oauth2.googleapis.com/revoke'

export class GoogleProvider {
  private csrfTokens = new Map<string, Date>()

  /**
   * Authenticate with Google AI API
   */
  async authenticate(config: GoogleAuthConfig): Promise<GoogleAuthResult> {
    const { method } = config

    try {
      switch (method) {
        case 'api-key':
          return await this.authenticateWithAPIKey(config)
        case 'oauth':
          return await this.authenticateWithOAuth(config)
        case 'cli':
          return await this.authenticateWithCLI(config)
        default:
          throw new AuthenticationError(
            'validation_failed',
            false,
            { provider: PROVIDER, timestamp: new Date() },
            `Unsupported authentication method: ${method}`
          )
      }
    } catch (error) {
      const authError = ErrorHandler.handle(PROVIDER, error)
      ErrorHandler.log(authError)
      throw authError
    }
  }

  /**
   * Authenticate with API key
   */
  private async authenticateWithAPIKey(config: GoogleAuthConfig): Promise<GoogleAuthResult> {
    const { apiKey } = config

    if (!apiKey) {
      throw new AuthenticationError(
        'validation_failed',
        false,
        { provider: PROVIDER, timestamp: new Date() },
        'API key is required'
      )
    }

    // Rate limiting
    const rateLimitResult = await authRateLimiter.checkLimit(`${PROVIDER}:${apiKey.substring(0, 10)}`)
    if (!rateLimitResult.allowed) {
      throw new AuthenticationError(
        'rate_limited',
        true,
        { provider: PROVIDER, timestamp: new Date() },
        rateLimitResult.retryAfter
          ? `Taking a quick breather! Try again in ${rateLimitResult.retryAfter} seconds. ☕`
          : 'Rate limit exceeded'
      )
    }

    // Validation
    const sanitized = inputValidator.sanitizeAPIKey(apiKey)
    const validationResult = await inputValidator.validateForStorage(PROVIDER, { apiKey: sanitized })

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

    // Verify with circuit breaker
    const models = await circuitBreakerRegistry.call(PROVIDER, async () => {
      return await this.verifyAPIKey(sanitized)
    })

    authRateLimiter.recordSuccess(`${PROVIDER}:${apiKey.substring(0, 10)}`)

    return {
      success: true,
      method: 'api-key',
      apiKey: sanitized,
      models,
      expiresAt: undefined
    }
  }

  /**
   * Authenticate with OAuth
   */
  private async authenticateWithOAuth(config: GoogleAuthConfig): Promise<GoogleAuthResult> {
    const { oauthToken, refreshToken, projectId } = config

    if (!oauthToken && !refreshToken) {
      throw new AuthenticationError(
        'validation_failed',
        false,
        { provider: PROVIDER, timestamp: new Date() },
        'OAuth token or refresh token is required',
        'https://console.cloud.google.com/apis/credentials'
      )
    }

    // Validate project ID if provided
    if (projectId) {
      const projectResult = inputValidator.validateProjectId(projectId)
      if (!projectResult.valid) {
        throw new AuthenticationError(
          'validation_failed',
          false,
          { provider: PROVIDER, timestamp: new Date() },
          projectResult.error || 'Invalid project ID',
          undefined,
          projectResult.hint
        )
      }
    }

    // If we have a refresh token, get a new access token
    let accessToken = oauthToken
    let expiresAt: Date | undefined

    if (refreshToken && !oauthToken) {
      const tokenResult = await this.refreshAccessToken(refreshToken)
      accessToken = tokenResult.accessToken
      expiresAt = tokenResult.expiresAt
    }

    if (!accessToken) {
      throw new AuthenticationError(
        'validation_failed',
        false,
        { provider: PROVIDER, timestamp: new Date() },
        'Unable to obtain access token'
      )
    }

    // Verify with circuit breaker
    const models = await circuitBreakerRegistry.call(PROVIDER, async () => {
      return await this.verifyOAuthToken(accessToken!, projectId)
    })

    return {
      success: true,
      method: 'oauth',
      oauthToken: accessToken,
      refreshToken,
      projectId,
      models,
      expiresAt
    }
  }

  /**
   * Authenticate with gcloud CLI
   */
  private async authenticateWithCLI(config: GoogleAuthConfig): Promise<GoogleAuthResult> {
    const { projectId } = config

    try {
      // Try to get token from gcloud CLI
      const { exec } = await import('child_process')
      const { promisify } = await import('util')
      const execAsync = promisify(exec)

      const { stdout } = await execAsync('gcloud auth print-access-token')
      const accessToken = stdout.trim()

      if (!accessToken) {
        throw new AuthenticationError(
          'validation_failed',
          false,
          { provider: PROVIDER, timestamp: new Date() },
          'No access token from gcloud CLI',
          undefined,
          'Run: gcloud auth login'
        )
      }

      // Verify token
      const models = await circuitBreakerRegistry.call(PROVIDER, async () => {
        return await this.verifyOAuthToken(accessToken, projectId)
      })

      // Tokens from gcloud expire in 1 hour
      const expiresAt = new Date(Date.now() + 60 * 60 * 1000)

      return {
        success: true,
        method: 'cli',
        oauthToken: accessToken,
        projectId,
        models,
        expiresAt
      }
    } catch (error) {
      if (error instanceof AuthenticationError) {
        throw error
      }

      throw new AuthenticationError(
        'validation_failed',
        false,
        { provider: PROVIDER, timestamp: new Date(), originalError: error as Error },
        'Failed to authenticate with gcloud CLI',
        undefined,
        'Make sure gcloud is installed and authenticated: gcloud auth login'
      )
    }
  }

  /**
   * Verify API key
   */
  private async verifyAPIKey(apiKey: string): Promise<string[]> {
    try {
      const response = await fetch(`${API_BASE_URL}/models?key=${apiKey}`, {
        method: 'GET',
        signal: AbortSignal.timeout(30000)
      })

      if (!response.ok) {
        throw this.handleHTTPError(response.status, await response.text())
      }

      const data = await response.json()
      const models = data.models?.map((model: any) => model.name.replace('models/', '')) || []

      // Filter for Gemini models
      const geminiModels = models.filter((id: string) => id.startsWith('gemini-'))

      if (geminiModels.length === 0) {
        return this.getKnownModels()
      }

      return geminiModels
    } catch (error) {
      if (error instanceof AuthenticationError) {
        throw error
      }

      if (error instanceof Error && (error.name === 'AbortError' || error.message.includes('timeout'))) {
        throw new AuthenticationError(
          'timeout',
          true,
          { provider: PROVIDER, timestamp: new Date(), originalError: error },
          'Request to Google AI timed out'
        )
      }

      throw new NetworkError(PROVIDER, { originalError: error as Error })
    }
  }

  /**
   * Verify OAuth token
   */
  private async verifyOAuthToken(token: string, projectId?: string): Promise<string[]> {
    try {
      const url = projectId
        ? `${API_BASE_URL}/models?access_token=${token}&project=${projectId}`
        : `${API_BASE_URL}/models?access_token=${token}`

      const response = await fetch(url, {
        method: 'GET',
        signal: AbortSignal.timeout(30000)
      })

      if (!response.ok) {
        throw this.handleHTTPError(response.status, await response.text())
      }

      const data = await response.json()
      const models = data.models?.map((model: any) => model.name.replace('models/', '')) || []

      const geminiModels = models.filter((id: string) => id.startsWith('gemini-'))

      if (geminiModels.length === 0) {
        return this.getKnownModels()
      }

      return geminiModels
    } catch (error) {
      if (error instanceof AuthenticationError) {
        throw error
      }

      throw new NetworkError(PROVIDER, { originalError: error as Error })
    }
  }

  /**
   * Refresh OAuth access token
   */
  private async refreshAccessToken(refreshToken: string): Promise<{
    accessToken: string
    expiresAt: Date
  }> {
    try {
      const response = await fetch(OAUTH_TOKEN_URL, {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        body: new URLSearchParams({
          grant_type: 'refresh_token',
          refresh_token: refreshToken,
          client_id: process.env.GOOGLE_CLIENT_ID || '',
          client_secret: process.env.GOOGLE_CLIENT_SECRET || ''
        })
      })

      if (!response.ok) {
        throw new AuthenticationError(
          'expired',
          false,
          { provider: PROVIDER, statusCode: response.status, timestamp: new Date() },
          'Failed to refresh Google OAuth token',
          'https://console.cloud.google.com/apis/credentials',
          'Re-authenticate to continue'
        )
      }

      const data = await response.json()
      const expiresIn = data.expires_in || 3600
      const expiresAt = new Date(Date.now() + expiresIn * 1000)

      return {
        accessToken: data.access_token,
        expiresAt
      }
    } catch (error) {
      if (error instanceof AuthenticationError) {
        throw error
      }

      throw new NetworkError(PROVIDER, { originalError: error as Error })
    }
  }

  /**
   * Generate CSRF token for OAuth flow
   */
  generateCSRFToken(): CSRFToken {
    const token = randomBytes(32).toString('hex')
    const expiresAt = new Date(Date.now() + 10 * 60 * 1000) // 10 minutes

    this.csrfTokens.set(token, expiresAt)

    // Cleanup expired tokens
    this.cleanupCSRFTokens()

    return { token, expiresAt }
  }

  /**
   * Validate CSRF token
   */
  validateCSRFToken(token: string): boolean {
    const expiresAt = this.csrfTokens.get(token)

    if (!expiresAt) {
      return false
    }

    if (expiresAt < new Date()) {
      this.csrfTokens.delete(token)
      return false
    }

    // Token is valid, remove it (one-time use)
    this.csrfTokens.delete(token)
    return true
  }

  /**
   * Cleanup expired CSRF tokens
   */
  private cleanupCSRFTokens(): void {
    const now = new Date()
    for (const [token, expiresAt] of this.csrfTokens.entries()) {
      if (expiresAt < now) {
        this.csrfTokens.delete(token)
      }
    }
  }

  /**
   * Get known Gemini models
   */
  private getKnownModels(): string[] {
    return [
      'gemini-1.5-pro',
      'gemini-1.5-flash',
      'gemini-1.0-pro'
    ]
  }

  /**
   * Handle HTTP errors
   */
  private handleHTTPError(status: number, body: string): AuthenticationError {
    let errorData: any = {}
    try {
      errorData = JSON.parse(body)
    } catch {
      // Ignore
    }

    const context = {
      provider: PROVIDER,
      statusCode: status,
      timestamp: new Date()
    }

    const errorMessage = errorData.error?.message || 'Unknown error'

    switch (status) {
      case 401:
        return new InvalidAPIKeyError(PROVIDER, context)

      case 403:
        return new AuthenticationError(
          'forbidden',
          false,
          context,
          errorMessage.includes('API key')
            ? 'This API key doesn\'t have permission'
            : errorMessage,
          'https://console.cloud.google.com/apis/credentials',
          'Check your API key permissions or enable the Generative Language API'
        )

      case 429:
        return new AuthenticationError(
          'rate_limited',
          true,
          context,
          'Google AI rate limit reached',
          undefined,
          'Wait a moment or upgrade your quota'
        )

      case 500:
      case 502:
      case 503:
      case 504:
        return new AuthenticationError(
          'server_error',
          true,
          context,
          'Google AI is experiencing issues',
          'https://status.cloud.google.com'
        )

      default:
        return new AuthenticationError(
          'unknown',
          true,
          context,
          `Google returned error ${status}: ${errorMessage}`,
          'https://console.cloud.google.com/apis/credentials'
        )
    }
  }

  /**
   * Test credentials
   */
  async testCredentials(config: GoogleAuthConfig): Promise<boolean> {
    try {
      await this.authenticate(config)
      return true
    } catch {
      return false
    }
  }

  /**
   * Get provider capabilities
   */
  getCapabilities() {
    return {
      supportsAPIKey: true,
      supportsOAuth: true,
      supportsCLI: true,
      requiresProject: true,
      credentialExpiry: true, // OAuth tokens expire
      helpUrl: 'https://console.cloud.google.com/apis/credentials',
      models: [
        {
          id: 'gemini-1.5-pro',
          name: 'Gemini 1.5 Pro',
          contextWindow: 1000000,
          supportsVision: true,
          inputPrice: 0.00125,
          outputPrice: 0.005
        },
        {
          id: 'gemini-1.5-flash',
          name: 'Gemini 1.5 Flash',
          contextWindow: 1000000,
          supportsVision: true,
          inputPrice: 0.000075,
          outputPrice: 0.0003
        },
        {
          id: 'gemini-1.0-pro',
          name: 'Gemini 1.0 Pro',
          contextWindow: 30720,
          supportsVision: false,
          inputPrice: 0.0005,
          outputPrice: 0.0015
        }
      ]
    }
  }
}

/**
 * Global Google provider instance
 */
export const googleProvider = new GoogleProvider()

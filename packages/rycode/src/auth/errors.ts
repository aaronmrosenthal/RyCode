/**
 * Rich Error Types for Provider Authentication
 *
 * Provides detailed error information to enable better user experience
 * and easier debugging.
 */

export type ErrorReason =
  | 'invalid_key'
  | 'expired'
  | 'rate_limited'
  | 'network'
  | 'unauthorized'
  | 'forbidden'
  | 'not_found'
  | 'server_error'
  | 'timeout'
  | 'invalid_format'
  | 'compromised'
  | 'validation_failed'
  | 'storage_failed'
  | 'unknown'

export interface ErrorContext {
  provider?: string
  method?: string
  statusCode?: number
  originalError?: Error
  timestamp: Date
  requestId?: string
}

/**
 * Base authentication error class
 */
export class AuthenticationError extends Error {
  constructor(
    public reason: ErrorReason,
    public retryable: boolean,
    public context: ErrorContext,
    public userMessage: string,
    public helpUrl?: string,
    public suggestedAction?: string
  ) {
    super(`Authentication failed: ${reason}`)
    this.name = 'AuthenticationError'
    Object.setPrototypeOf(this, AuthenticationError.prototype)
  }

  /**
   * Convert to JSON for logging
   */
  toJSON() {
    return {
      name: this.name,
      reason: this.reason,
      retryable: this.retryable,
      userMessage: this.userMessage,
      helpUrl: this.helpUrl,
      suggestedAction: this.suggestedAction,
      context: {
        ...this.context,
        originalError: this.context.originalError?.message
      }
    }
  }

  /**
   * Get user-friendly error message
   */
  getUserMessage(): string {
    return this.userMessage
  }

  /**
   * Check if error is retryable
   */
  isRetryable(): boolean {
    return this.retryable
  }
}

/**
 * Invalid API key error
 */
export class InvalidAPIKeyError extends AuthenticationError {
  constructor(provider: string, context?: Partial<ErrorContext>) {
    super(
      'invalid_key',
      false,
      {
        provider,
        timestamp: new Date(),
        ...context
      },
      `The API key for ${provider} is invalid or has been revoked`,
      getProviderHelpUrl(provider),
      'Double-check your API key or generate a new one'
    )
    this.name = 'InvalidAPIKeyError'
  }
}

/**
 * Expired credentials error
 */
export class ExpiredCredentialsError extends AuthenticationError {
  constructor(provider: string, context?: Partial<ErrorContext>) {
    super(
      'expired',
      false,
      {
        provider,
        timestamp: new Date(),
        ...context
      },
      `Your credentials for ${provider} have expired`,
      getProviderHelpUrl(provider),
      'Re-authenticate to continue using this provider'
    )
    this.name = 'ExpiredCredentialsError'
  }
}

/**
 * Rate limit error
 */
export class RateLimitError extends AuthenticationError {
  constructor(
    provider: string,
    retryAfter: number,
    context?: Partial<ErrorContext>
  ) {
    super(
      'rate_limited',
      true,
      {
        provider,
        timestamp: new Date(),
        ...context
      },
      `Taking a quick breather! Try again in ${retryAfter} seconds. ☕`,
      undefined,
      'Pro tip: Batch your requests for better flow'
    )
    this.name = 'RateLimitError'
  }
}

/**
 * Network error
 */
export class NetworkError extends AuthenticationError {
  constructor(provider: string, context?: Partial<ErrorContext>) {
    super(
      'network',
      true,
      {
        provider,
        timestamp: new Date(),
        ...context
      },
      `Unable to connect to ${provider}. Check your internet connection.`,
      undefined,
      'Try again in a moment'
    )
    this.name = 'NetworkError'
  }
}

/**
 * Validation error
 */
export class ValidationError extends AuthenticationError {
  constructor(
    provider: string,
    message: string,
    hint?: string,
    context?: Partial<ErrorContext>
  ) {
    super(
      'validation_failed',
      false,
      {
        provider,
        timestamp: new Date(),
        ...context
      },
      message,
      getProviderHelpUrl(provider),
      hint
    )
    this.name = 'ValidationError'
  }
}

/**
 * Storage error
 */
export class StorageError extends AuthenticationError {
  constructor(provider: string, operation: string, context?: Partial<ErrorContext>) {
    super(
      'storage_failed',
      true,
      {
        provider,
        timestamp: new Date(),
        ...context
      },
      `Failed to ${operation} credentials for ${provider}`,
      undefined,
      'Check your system keychain permissions'
    )
    this.name = 'StorageError'
  }
}

/**
 * Compromised key error
 */
export class CompromisedKeyError extends AuthenticationError {
  constructor(provider: string, context?: Partial<ErrorContext>) {
    super(
      'compromised',
      false,
      {
        provider,
        timestamp: new Date(),
        ...context
      },
      `This API key has been compromised and cannot be used`,
      getProviderHelpUrl(provider),
      'Please generate a new API key immediately'
    )
    this.name = 'CompromisedKeyError'
  }
}

/**
 * Helper function to get provider help URLs
 */
function getProviderHelpUrl(provider: string): string {
  const urls: Record<string, string> = {
    anthropic: 'https://console.anthropic.com/settings/keys',
    openai: 'https://platform.openai.com/api-keys',
    grok: 'https://console.x.ai/api-keys',
    qwen: 'https://dashscope.console.aliyun.com/',
    google: 'https://console.cloud.google.com/apis/credentials'
  }
  return urls[provider] || ''
}

/**
 * Parse HTTP error response into AuthenticationError
 */
export function parseHTTPError(
  provider: string,
  statusCode: number,
  responseBody?: any,
  originalError?: Error
): AuthenticationError {
  const context: ErrorContext = {
    provider,
    statusCode,
    timestamp: new Date(),
    originalError
  }

  switch (statusCode) {
    case 401:
      return new InvalidAPIKeyError(provider, context)

    case 403:
      return new AuthenticationError(
        'forbidden',
        false,
        context,
        `You don't have permission to access this ${provider} resource`,
        getProviderHelpUrl(provider),
        'Check your API key permissions'
      )

    case 404:
      return new AuthenticationError(
        'not_found',
        false,
        context,
        `The ${provider} resource was not found`,
        getProviderHelpUrl(provider)
      )

    case 429:
      const retryAfter = responseBody?.retry_after || 60
      return new RateLimitError(provider, retryAfter, context)

    case 500:
    case 502:
    case 503:
    case 504:
      return new AuthenticationError(
        'server_error',
        true,
        context,
        `${provider} is experiencing issues. Please try again later.`,
        undefined,
        'Check the provider status page'
      )

    default:
      return new AuthenticationError(
        'unknown',
        true,
        context,
        `Something went wrong with ${provider} (${statusCode})`,
        getProviderHelpUrl(provider)
      )
  }
}

/**
 * Parse network error into AuthenticationError
 */
export function parseNetworkError(
  provider: string,
  error: Error
): AuthenticationError {
  if (error.message.includes('timeout')) {
    return new AuthenticationError(
      'timeout',
      true,
      { provider, timestamp: new Date(), originalError: error },
      `Request to ${provider} timed out`,
      undefined,
      'Try again in a moment'
    )
  }

  return new NetworkError(provider, { originalError: error })
}

/**
 * Error handler for authentication operations
 */
export class ErrorHandler {
  /**
   * Handle error and return user-friendly AuthenticationError
   */
  static handle(provider: string, error: unknown): AuthenticationError {
    // Already an AuthenticationError
    if (error instanceof AuthenticationError) {
      return error
    }

    // HTTP errors
    if (error && typeof error === 'object' && 'statusCode' in error) {
      const httpError = error as { statusCode: number; body?: any }
      return parseHTTPError(provider, httpError.statusCode, httpError.body)
    }

    // Network errors
    if (error instanceof Error) {
      if (
        error.message.includes('ECONNREFUSED') ||
        error.message.includes('ENOTFOUND') ||
        error.message.includes('ETIMEDOUT')
      ) {
        return parseNetworkError(provider, error)
      }
    }

    // Unknown errors
    return new AuthenticationError(
      'unknown',
      true,
      {
        provider,
        timestamp: new Date(),
        originalError: error instanceof Error ? error : undefined
      },
      `An unexpected error occurred with ${provider}`,
      getProviderHelpUrl(provider),
      'Please try again'
    )
  }

  /**
   * Log error with appropriate level
   */
  static log(error: AuthenticationError, logger: Console = console): void {
    const logData = {
      error: error.toJSON(),
      stack: error.stack
    }

    if (error.retryable) {
      logger.warn('Retryable authentication error:', logData)
    } else {
      logger.error('Non-retryable authentication error:', logData)
    }
  }
}

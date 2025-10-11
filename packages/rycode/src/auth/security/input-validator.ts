/**
 * Input Validation for Provider Authentication
 *
 * Validates and sanitizes API keys and OAuth tokens before storage or use.
 * Checks for known compromised keys and validates format.
 */

export interface ValidationResult {
  valid: boolean
  error?: string
  hint?: string
  helpUrl?: string
}

/**
 * Known API key format patterns for each provider
 */
const API_KEY_PATTERNS: Record<string, RegExp> = {
  anthropic: /^sk-ant-api03-[A-Za-z0-9_-]{95}$/,
  openai: /^sk-[A-Za-z0-9]{48}$/,
  grok: /^xai-[A-Za-z0-9]{32,}$/,
  qwen: /^qwen-[A-Za-z0-9-]{32,}$/,
  // Google uses OAuth primarily, but also supports API keys
  google: /^AIza[A-Za-z0-9_-]{35}$/
}

/**
 * Help URLs for getting API keys
 */
const HELP_URLS: Record<string, string> = {
  anthropic: 'https://console.anthropic.com/settings/keys',
  openai: 'https://platform.openai.com/api-keys',
  grok: 'https://console.x.ai/api-keys',
  qwen: 'https://dashscope.console.aliyun.com/',
  google: 'https://console.cloud.google.com/apis/credentials'
}

/**
 * Known compromised API key hashes (SHA-256)
 * In production, this would be fetched from a service
 */
const COMPROMISED_KEY_HASHES = new Set<string>([
  // This would be populated with actual compromised key hashes
  // For now, it's an empty set for demonstration
])

export class InputValidator {
  /**
   * Validate API key format
   */
  validateAPIKeyFormat(provider: string, apiKey: string): ValidationResult {
    // Basic checks
    if (!apiKey || typeof apiKey !== 'string') {
      return {
        valid: false,
        error: 'API key is required',
        hint: 'Please provide a valid API key'
      }
    }

    // Trim whitespace
    apiKey = apiKey.trim()

    // Check length
    if (apiKey.length < 20) {
      return {
        valid: false,
        error: 'API key too short',
        hint: 'API keys are typically longer than 20 characters. Did you copy the full key?'
      }
    }

    if (apiKey.length > 500) {
      return {
        valid: false,
        error: 'API key too long',
        hint: 'This doesn\'t look like a valid API key'
      }
    }

    // Check for suspicious characters
    if (!/^[A-Za-z0-9_-]+$/.test(apiKey)) {
      return {
        valid: false,
        error: 'Invalid characters in API key',
        hint: 'API keys should only contain letters, numbers, hyphens, and underscores'
      }
    }

    // Check provider-specific format
    const pattern = API_KEY_PATTERNS[provider]
    if (pattern && !pattern.test(apiKey)) {
      let hint = `This doesn't match the expected format for ${provider} API keys`

      // Provider-specific hints
      if (provider === 'anthropic' && !apiKey.startsWith('sk-ant-')) {
        hint = 'Anthropic API keys start with "sk-ant-api03-"'
      } else if (provider === 'openai' && !apiKey.startsWith('sk-')) {
        hint = 'OpenAI API keys start with "sk-"'
      } else if (provider === 'grok' && !apiKey.startsWith('xai-')) {
        hint = 'Grok API keys start with "xai-"'
      }

      return {
        valid: false,
        error: 'Invalid API key format',
        hint,
        helpUrl: HELP_URLS[provider]
      }
    }

    return { valid: true }
  }

  /**
   * Sanitize API key (remove common issues)
   */
  sanitizeAPIKey(apiKey: string): string {
    return apiKey
      .trim()
      .replace(/[\r\n\t]/g, '') // Remove newlines and tabs
      .replace(/^["']|["']$/g, '') // Remove quotes
  }

  /**
   * Check if an API key has been compromised
   * In production, this would call an external service
   */
  async isCompromisedKey(apiKey: string): Promise<boolean> {
    // Hash the key for comparison
    const hash = await this.hashKey(apiKey)
    return COMPROMISED_KEY_HASHES.has(hash)
  }

  /**
   * Validate OAuth token format
   */
  validateOAuthToken(provider: string, token: string): ValidationResult {
    if (!token || typeof token !== 'string') {
      return {
        valid: false,
        error: 'OAuth token is required'
      }
    }

    // Basic JWT format check (if it looks like a JWT)
    if (token.includes('.')) {
      const parts = token.split('.')
      if (parts.length !== 3) {
        return {
          valid: false,
          error: 'Invalid JWT token format',
          hint: 'JWT tokens should have exactly 3 parts separated by dots'
        }
      }

      // Check if parts are base64
      for (const part of parts) {
        if (!/^[A-Za-z0-9_-]+$/.test(part)) {
          return {
            valid: false,
            error: 'Invalid token encoding'
          }
        }
      }
    }

    return { valid: true }
  }

  /**
   * Validate project ID (for Google)
   */
  validateProjectId(projectId: string): ValidationResult {
    if (!projectId || typeof projectId !== 'string') {
      return {
        valid: false,
        error: 'Project ID is required'
      }
    }

    // Google Cloud project IDs have specific format
    if (!/^[a-z][a-z0-9-]{4,28}[a-z0-9]$/.test(projectId)) {
      return {
        valid: false,
        error: 'Invalid project ID format',
        hint: 'Project IDs must be 6-30 characters, start with a letter, and contain only lowercase letters, numbers, and hyphens'
      }
    }

    return { valid: true }
  }

  /**
   * Mask API key for logging (show only first/last 4 chars)
   */
  maskAPIKey(apiKey: string): string {
    if (!apiKey || apiKey.length < 12) {
      return '****'
    }

    const start = apiKey.substring(0, 4)
    const end = apiKey.substring(apiKey.length - 4)
    return `${start}...${end}`
  }

  /**
   * Hash an API key for comparison (SHA-256)
   */
  private async hashKey(key: string): Promise<string> {
    const encoder = new TextEncoder()
    const data = encoder.encode(key)
    const hashBuffer = await crypto.subtle.digest('SHA-256', data)
    const hashArray = Array.from(new Uint8Array(hashBuffer))
    return hashArray.map(b => b.toString(16).padStart(2, '0')).join('')
  }

  /**
   * Comprehensive validation for storing credentials
   */
  async validateForStorage(
    provider: string,
    credentials: { apiKey?: string; oauthToken?: string; projectId?: string }
  ): Promise<ValidationResult> {
    // Validate API key if provided
    if (credentials.apiKey) {
      const sanitized = this.sanitizeAPIKey(credentials.apiKey)
      const formatResult = this.validateAPIKeyFormat(provider, sanitized)

      if (!formatResult.valid) {
        return formatResult
      }

      // Check if compromised
      const compromised = await this.isCompromisedKey(sanitized)
      if (compromised) {
        return {
          valid: false,
          error: 'This API key has been compromised',
          hint: 'Please generate a new API key from your provider console',
          helpUrl: HELP_URLS[provider]
        }
      }
    }

    // Validate OAuth token if provided
    if (credentials.oauthToken) {
      const tokenResult = this.validateOAuthToken(provider, credentials.oauthToken)
      if (!tokenResult.valid) {
        return tokenResult
      }
    }

    // Validate project ID if provided (Google)
    if (credentials.projectId) {
      const projectResult = this.validateProjectId(credentials.projectId)
      if (!projectResult.valid) {
        return projectResult
      }
    }

    return { valid: true }
  }
}

export const inputValidator = new InputValidator()

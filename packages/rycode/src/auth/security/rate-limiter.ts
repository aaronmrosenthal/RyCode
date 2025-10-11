/**
 * Rate Limiter for Provider Authentication
 *
 * Implements token bucket algorithm to prevent brute force attacks
 * and excessive API calls to authentication endpoints.
 */

export interface RateLimitConfig {
  maxAttempts: number
  windowMs: number
  blockDurationMs?: number
}

export interface RateLimitResult {
  allowed: boolean
  remaining: number
  resetAt: Date
  retryAfter?: number
}

export class RateLimiter {
  private attempts = new Map<string, number[]>()
  private blockedUntil = new Map<string, Date>()

  constructor(private config: RateLimitConfig = {
    maxAttempts: 5,
    windowMs: 60000, // 1 minute
    blockDurationMs: 300000 // 5 minutes
  }) {}

  /**
   * Check if a request should be allowed
   */
  async checkLimit(key: string): Promise<RateLimitResult> {
    const now = Date.now()

    // Check if currently blocked
    const blockedUntil = this.blockedUntil.get(key)
    if (blockedUntil && blockedUntil.getTime() > now) {
      return {
        allowed: false,
        remaining: 0,
        resetAt: blockedUntil,
        retryAfter: Math.ceil((blockedUntil.getTime() - now) / 1000)
      }
    }

    // Clean up expired block
    if (blockedUntil) {
      this.blockedUntil.delete(key)
    }

    // Get recent attempts within the time window
    const userAttempts = this.attempts.get(key) || []
    const recentAttempts = userAttempts.filter(t => now - t < this.config.windowMs)

    // Check if limit exceeded
    if (recentAttempts.length >= this.config.maxAttempts) {
      const blockUntil = new Date(now + (this.config.blockDurationMs || this.config.windowMs))
      this.blockedUntil.set(key, blockUntil)

      return {
        allowed: false,
        remaining: 0,
        resetAt: blockUntil,
        retryAfter: Math.ceil((this.config.blockDurationMs || this.config.windowMs) / 1000)
      }
    }

    // Allow request and record attempt
    recentAttempts.push(now)
    this.attempts.set(key, recentAttempts)

    const resetAt = new Date(now + this.config.windowMs)

    return {
      allowed: true,
      remaining: this.config.maxAttempts - recentAttempts.length,
      resetAt
    }
  }

  /**
   * Record a successful operation (optional - for tracking)
   */
  recordSuccess(key: string): void {
    // Remove from blocked list on success
    this.blockedUntil.delete(key)

    // Optionally clear attempts on success (more forgiving)
    // this.attempts.delete(key)
  }

  /**
   * Manually reset limits for a key
   */
  reset(key: string): void {
    this.attempts.delete(key)
    this.blockedUntil.delete(key)
  }

  /**
   * Clean up old entries to prevent memory leaks
   */
  cleanup(): void {
    const now = Date.now()

    // Clean up attempts
    for (const [key, attempts] of this.attempts.entries()) {
      const recent = attempts.filter(t => now - t < this.config.windowMs)
      if (recent.length === 0) {
        this.attempts.delete(key)
      } else {
        this.attempts.set(key, recent)
      }
    }

    // Clean up expired blocks
    for (const [key, until] of this.blockedUntil.entries()) {
      if (until.getTime() <= now) {
        this.blockedUntil.delete(key)
      }
    }
  }

  /**
   * Get current status for a key
   */
  getStatus(key: string): {
    attempts: number
    blocked: boolean
    blockedUntil?: Date
  } {
    const now = Date.now()
    const userAttempts = this.attempts.get(key) || []
    const recentAttempts = userAttempts.filter(t => now - t < this.config.windowMs)
    const blockedUntil = this.blockedUntil.get(key)

    return {
      attempts: recentAttempts.length,
      blocked: blockedUntil ? blockedUntil.getTime() > now : false,
      blockedUntil: blockedUntil && blockedUntil.getTime() > now ? blockedUntil : undefined
    }
  }
}

/**
 * Global rate limiter instances for different operations
 */
export const authRateLimiter = new RateLimiter({
  maxAttempts: 5,
  windowMs: 60000, // 1 minute
  blockDurationMs: 300000 // 5 minutes
})

export const apiRateLimiter = new RateLimiter({
  maxAttempts: 60,
  windowMs: 60000, // 60 per minute
  blockDurationMs: 60000 // 1 minute
})

// Cleanup every 5 minutes
setInterval(() => {
  authRateLimiter.cleanup()
  apiRateLimiter.cleanup()
}, 300000)

/**
 * Circuit Breaker Pattern for Provider Resilience
 *
 * Prevents cascading failures by stopping requests to failing providers
 * and automatically recovering when the provider becomes healthy again.
 */

export type CircuitState = 'closed' | 'open' | 'half-open'

export interface CircuitBreakerConfig {
  failureThreshold: number
  successThreshold: number
  timeout: number
  resetTimeout: number
}

export interface CircuitBreakerStats {
  state: CircuitState
  failures: number
  successes: number
  lastFailure?: Date
  nextAttempt?: Date
}

export class CircuitBreakerError extends Error {
  constructor(
    public provider: string,
    public state: CircuitState,
    public nextAttempt?: Date
  ) {
    super(`Circuit breaker is ${state} for provider ${provider}`)
    this.name = 'CircuitBreakerError'
  }
}

export class CircuitBreaker {
  private state: CircuitState = 'closed'
  private failures = 0
  private successes = 0
  private lastFailure?: Date
  private nextAttempt?: Date

  constructor(
    private provider: string,
    private config: CircuitBreakerConfig = {
      failureThreshold: 5, // Open circuit after 5 failures
      successThreshold: 2, // Close circuit after 2 successes
      timeout: 60000, // 1 minute timeout for requests
      resetTimeout: 30000 // Try again after 30 seconds
    }
  ) {}

  /**
   * Execute a function with circuit breaker protection
   */
  async call<T>(fn: () => Promise<T>): Promise<T> {
    // Check circuit state
    if (this.state === 'open') {
      // Check if we should try half-open
      if (this.nextAttempt && Date.now() >= this.nextAttempt.getTime()) {
        this.state = 'half-open'
        this.successes = 0
      } else {
        throw new CircuitBreakerError(this.provider, this.state, this.nextAttempt)
      }
    }

    try {
      // Execute with timeout
      const result = await this.executeWithTimeout(fn, this.config.timeout)

      // Record success
      this.onSuccess()

      return result
    } catch (error) {
      // Record failure
      this.onFailure()

      throw error
    }
  }

  /**
   * Execute function with timeout
   */
  private async executeWithTimeout<T>(
    fn: () => Promise<T>,
    timeoutMs: number
  ): Promise<T> {
    return Promise.race([
      fn(),
      new Promise<T>((_, reject) =>
        setTimeout(() => reject(new Error('Request timeout')), timeoutMs)
      )
    ])
  }

  /**
   * Handle successful operation
   */
  private onSuccess(): void {
    this.failures = 0

    if (this.state === 'half-open') {
      this.successes++

      // Close circuit after threshold successes
      if (this.successes >= this.config.successThreshold) {
        this.state = 'closed'
        this.successes = 0
        this.lastFailure = undefined
        this.nextAttempt = undefined
      }
    }
  }

  /**
   * Handle failed operation
   */
  private onFailure(): void {
    this.failures++
    this.lastFailure = new Date()

    if (this.state === 'half-open') {
      // Go back to open on any failure in half-open state
      this.state = 'open'
      this.nextAttempt = new Date(Date.now() + this.config.resetTimeout)
    } else if (this.failures >= this.config.failureThreshold) {
      // Open circuit after threshold failures
      this.state = 'open'
      this.nextAttempt = new Date(Date.now() + this.config.resetTimeout)
    }
  }

  /**
   * Manually reset the circuit breaker
   */
  reset(): void {
    this.state = 'closed'
    this.failures = 0
    this.successes = 0
    this.lastFailure = undefined
    this.nextAttempt = undefined
  }

  /**
   * Get current circuit breaker stats
   */
  getStats(): CircuitBreakerStats {
    return {
      state: this.state,
      failures: this.failures,
      successes: this.successes,
      lastFailure: this.lastFailure,
      nextAttempt: this.nextAttempt
    }
  }

  /**
   * Check if circuit breaker is healthy
   */
  isHealthy(): boolean {
    return this.state === 'closed'
  }
}

/**
 * Circuit Breaker Registry for managing multiple providers
 */
export class CircuitBreakerRegistry {
  private breakers = new Map<string, CircuitBreaker>()

  constructor(private config?: CircuitBreakerConfig) {}

  /**
   * Get or create a circuit breaker for a provider
   */
  getBreaker(provider: string): CircuitBreaker {
    if (!this.breakers.has(provider)) {
      this.breakers.set(provider, new CircuitBreaker(provider, this.config))
    }
    return this.breakers.get(provider)!
  }

  /**
   * Execute a function with circuit breaker protection for a provider
   */
  async call<T>(provider: string, fn: () => Promise<T>): Promise<T> {
    const breaker = this.getBreaker(provider)
    return breaker.call(fn)
  }

  /**
   * Reset circuit breaker for a provider
   */
  reset(provider: string): void {
    const breaker = this.breakers.get(provider)
    if (breaker) {
      breaker.reset()
    }
  }

  /**
   * Get stats for all providers
   */
  getAllStats(): Map<string, CircuitBreakerStats> {
    const stats = new Map<string, CircuitBreakerStats>()
    for (const [provider, breaker] of this.breakers.entries()) {
      stats.set(provider, breaker.getStats())
    }
    return stats
  }

  /**
   * Get list of unhealthy providers
   */
  getUnhealthyProviders(): string[] {
    const unhealthy: string[] = []
    for (const [provider, breaker] of this.breakers.entries()) {
      if (!breaker.isHealthy()) {
        unhealthy.push(provider)
      }
    }
    return unhealthy
  }

  /**
   * Check if a provider is healthy
   */
  isProviderHealthy(provider: string): boolean {
    const breaker = this.breakers.get(provider)
    return breaker ? breaker.isHealthy() : true // Assume healthy if not tracked
  }
}

/**
 * Global circuit breaker registry
 */
export const circuitBreakerRegistry = new CircuitBreakerRegistry({
  failureThreshold: 5,
  successThreshold: 2,
  timeout: 30000, // 30 seconds
  resetTimeout: 60000 // 1 minute
})

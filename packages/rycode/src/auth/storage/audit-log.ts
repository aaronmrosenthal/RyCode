/**
 * Audit Log - Security Event Tracking
 *
 * Tracks authentication events for security monitoring, debugging,
 * and compliance. Provides insights into authentication patterns
 * and potential security issues.
 */

import { Log } from '../../util/log'
import path from 'path'
import { Global } from '../../global'

const log = Log.create({ service: 'auth-audit' })

export type AuditEventType =
  | 'auth_attempt'
  | 'auth_success'
  | 'auth_failure'
  | 'credential_stored'
  | 'credential_retrieved'
  | 'credential_removed'
  | 'credential_expired'
  | 'rate_limit_exceeded'
  | 'circuit_breaker_opened'
  | 'circuit_breaker_closed'
  | 'validation_failed'
  | 'token_refreshed'
  | 'migration_completed'

export interface AuditEvent {
  timestamp: Date
  eventType: AuditEventType
  provider: string
  method?: string
  success?: boolean
  reason?: string
  metadata?: Record<string, any>
  ipAddress?: string
  userAgent?: string
  riskScore?: number
}

export interface AuditQuery {
  provider?: string
  eventType?: AuditEventType
  startDate?: Date
  endDate?: Date
  success?: boolean
  limit?: number
}

export interface AuditSummary {
  totalEvents: number
  successRate: number
  failureRate: number
  byProvider: Record<string, number>
  byEventType: Record<string, number>
  recentFailures: AuditEvent[]
  riskEvents: AuditEvent[]
}

export class AuditLog {
  private events: AuditEvent[] = []
  private readonly maxEvents = 1000 // Keep last 1000 events in memory
  private readonly filepath = path.join(Global.Path.data, 'auth-audit.log')

  /**
   * Record an audit event
   */
  async record(event: Omit<AuditEvent, 'timestamp'>): Promise<void> {
    const auditEvent: AuditEvent = {
      ...event,
      timestamp: new Date()
    }

    // Add to in-memory store
    this.events.push(auditEvent)

    // Keep only last N events
    if (this.events.length > this.maxEvents) {
      this.events = this.events.slice(-this.maxEvents)
    }

    // Log based on event type
    this.logEvent(auditEvent)

    // Persist to disk (async, don't wait)
    this.persistEvent(auditEvent).catch(error => {
      log.error('Failed to persist audit event', { error })
    })
  }

  /**
   * Record authentication attempt
   */
  async recordAuthAttempt(
    provider: string,
    method: string,
    metadata?: Record<string, any>
  ): Promise<void> {
    await this.record({
      eventType: 'auth_attempt',
      provider,
      method,
      metadata
    })
  }

  /**
   * Record authentication success
   */
  async recordAuthSuccess(
    provider: string,
    method: string,
    metadata?: Record<string, any>
  ): Promise<void> {
    await this.record({
      eventType: 'auth_success',
      provider,
      method,
      success: true,
      metadata
    })
  }

  /**
   * Record authentication failure
   */
  async recordAuthFailure(
    provider: string,
    method: string,
    reason: string,
    metadata?: Record<string, any>
  ): Promise<void> {
    const riskScore = this.calculateRiskScore(provider, reason)

    await this.record({
      eventType: 'auth_failure',
      provider,
      method,
      success: false,
      reason,
      riskScore,
      metadata
    })
  }

  /**
   * Record credential storage
   */
  async recordCredentialStored(
    provider: string,
    method: string
  ): Promise<void> {
    await this.record({
      eventType: 'credential_stored',
      provider,
      method,
      success: true
    })
  }

  /**
   * Record credential retrieval
   */
  async recordCredentialRetrieved(provider: string): Promise<void> {
    await this.record({
      eventType: 'credential_retrieved',
      provider,
      success: true
    })
  }

  /**
   * Record credential removal
   */
  async recordCredentialRemoved(provider: string): Promise<void> {
    await this.record({
      eventType: 'credential_removed',
      provider,
      success: true
    })
  }

  /**
   * Record rate limit exceeded
   */
  async recordRateLimitExceeded(
    provider: string,
    metadata?: Record<string, any>
  ): Promise<void> {
    await this.record({
      eventType: 'rate_limit_exceeded',
      provider,
      success: false,
      reason: 'Rate limit exceeded',
      riskScore: 3,
      metadata
    })
  }

  /**
   * Record circuit breaker event
   */
  async recordCircuitBreakerEvent(
    provider: string,
    opened: boolean
  ): Promise<void> {
    await this.record({
      eventType: opened ? 'circuit_breaker_opened' : 'circuit_breaker_closed',
      provider,
      success: !opened,
      metadata: { state: opened ? 'open' : 'closed' }
    })
  }

  /**
   * Query audit events
   */
  query(filters: AuditQuery = {}): AuditEvent[] {
    let results = [...this.events]

    // Filter by provider
    if (filters.provider) {
      results = results.filter(e => e.provider === filters.provider)
    }

    // Filter by event type
    if (filters.eventType) {
      results = results.filter(e => e.eventType === filters.eventType)
    }

    // Filter by date range
    if (filters.startDate) {
      results = results.filter(e => e.timestamp >= filters.startDate!)
    }

    if (filters.endDate) {
      results = results.filter(e => e.timestamp <= filters.endDate!)
    }

    // Filter by success
    if (filters.success !== undefined) {
      results = results.filter(e => e.success === filters.success)
    }

    // Sort by timestamp (newest first)
    results.sort((a, b) => b.timestamp.getTime() - a.timestamp.getTime())

    // Apply limit
    if (filters.limit) {
      results = results.slice(0, filters.limit)
    }

    return results
  }

  /**
   * Get audit summary
   */
  getSummary(provider?: string): AuditSummary {
    const events = provider
      ? this.events.filter(e => e.provider === provider)
      : this.events

    const totalEvents = events.length
    const successCount = events.filter(e => e.success === true).length
    const failureCount = events.filter(e => e.success === false).length

    // By provider
    const byProvider: Record<string, number> = {}
    for (const event of events) {
      byProvider[event.provider] = (byProvider[event.provider] || 0) + 1
    }

    // By event type
    const byEventType: Record<string, number> = {}
    for (const event of events) {
      byEventType[event.eventType] = (byEventType[event.eventType] || 0) + 1
    }

    // Recent failures
    const recentFailures = events
      .filter(e => e.success === false)
      .sort((a, b) => b.timestamp.getTime() - a.timestamp.getTime())
      .slice(0, 10)

    // Risk events (high risk score)
    const riskEvents = events
      .filter(e => e.riskScore && e.riskScore >= 5)
      .sort((a, b) => (b.riskScore || 0) - (a.riskScore || 0))
      .slice(0, 10)

    return {
      totalEvents,
      successRate: totalEvents > 0 ? successCount / totalEvents : 0,
      failureRate: totalEvents > 0 ? failureCount / totalEvents : 0,
      byProvider,
      byEventType,
      recentFailures,
      riskEvents
    }
  }

  /**
   * Get recent events
   */
  getRecent(limit: number = 10): AuditEvent[] {
    return this.query({ limit })
  }

  /**
   * Get failures for a provider
   */
  getFailures(provider: string, limit: number = 10): AuditEvent[] {
    return this.query({ provider, success: false, limit })
  }

  /**
   * Check for suspicious activity
   */
  detectSuspiciousActivity(provider: string): {
    suspicious: boolean
    reasons: string[]
  } {
    const last5Min = new Date(Date.now() - 5 * 60 * 1000)
    const recentEvents = this.query({
      provider,
      startDate: last5Min
    })

    const reasons: string[] = []

    // Too many failures
    const failures = recentEvents.filter(e => e.success === false)
    if (failures.length >= 5) {
      reasons.push(`${failures.length} failures in last 5 minutes`)
    }

    // High risk events
    const highRisk = recentEvents.filter(e => e.riskScore && e.riskScore >= 5)
    if (highRisk.length > 0) {
      reasons.push(`${highRisk.length} high-risk events detected`)
    }

    // Rate limit exceeded multiple times
    const rateLimits = recentEvents.filter(e => e.eventType === 'rate_limit_exceeded')
    if (rateLimits.length >= 3) {
      reasons.push('Multiple rate limit violations')
    }

    return {
      suspicious: reasons.length > 0,
      reasons
    }
  }

  /**
   * Calculate risk score (1-10)
   */
  private calculateRiskScore(provider: string, reason: string): number {
    let score = 1

    // Check recent failures
    const last10Min = new Date(Date.now() - 10 * 60 * 1000)
    const recentFailures = this.query({
      provider,
      success: false,
      startDate: last10Min
    }).length

    score += Math.min(recentFailures, 5)

    // Increase score for specific reasons
    if (reason.includes('invalid') || reason.includes('compromised')) {
      score += 3
    }

    if (reason.includes('rate limit')) {
      score += 2
    }

    return Math.min(score, 10)
  }

  /**
   * Log event to console
   */
  private logEvent(event: AuditEvent): void {
    const logData = {
      provider: event.provider,
      eventType: event.eventType,
      success: event.success,
      reason: event.reason,
      riskScore: event.riskScore
    }

    if (event.riskScore && event.riskScore >= 7) {
      log.error('High-risk authentication event', logData)
    } else if (event.success === false) {
      log.warn('Authentication event failed', logData)
    } else {
      log.debug('Authentication event', logData)
    }
  }

  /**
   * Persist event to disk
   */
  private async persistEvent(event: AuditEvent): Promise<void> {
    try {
      const line = JSON.stringify(event) + '\n'
      await Bun.write(this.filepath, line, { append: true })
    } catch (error) {
      // Don't throw - logging failure shouldn't break auth
      log.error('Failed to persist audit event to disk', { error })
    }
  }

  /**
   * Clear all events (for testing)
   */
  clear(): void {
    this.events = []
  }

  /**
   * Export events to JSON
   */
  export(): AuditEvent[] {
    return [...this.events]
  }
}

/**
 * Global audit log instance
 */
export const auditLog = new AuditLog()

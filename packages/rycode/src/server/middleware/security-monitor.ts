import { Log } from "../../util/log"
import type { Context } from "hono"

export namespace SecurityMonitor {
  const log = Log.create({ service: "security-monitor" })

  export type SecurityEventType =
    | "auth_failure"
    | "rate_limit"
    | "path_traversal"
    | "request_too_large"
    | "suspicious_activity"
    | "weak_key_attempt"

  interface SecurityEvent {
    type: SecurityEventType
    ip: string
    path: string
    userAgent?: string
    timestamp: number
    details: Record<string, any>
  }

  const recentEvents = new Map<string, SecurityEvent[]>()
  const ALERT_THRESHOLD = 10
  const ALERT_WINDOW_MS = 60_000
  const CLEANUP_INTERVAL_MS = 300_000

  setInterval(() => {
    const now = Date.now()
    for (const [key, events] of recentEvents.entries()) {
      const recent = events.filter((e) => now - e.timestamp < ALERT_WINDOW_MS * 5)
      if (recent.length === 0) {
        recentEvents.delete(key)
      } else {
        recentEvents.set(key, recent)
      }
    }
  }, CLEANUP_INTERVAL_MS).unref()

  export function track(c: Context, type: SecurityEventType, details: Record<string, any> = {}) {
    const ip = getIP(c)
    const userAgent = c.req.header("User-Agent")

    const event: SecurityEvent = {
      type,
      ip,
      path: c.req.path,
      userAgent,
      timestamp: Date.now(),
      details,
    }

    log.warn("security event", {
      type: event.type,
      ip: event.ip,
      path: event.path,
      details: event.details,
    })

    const key = ip + ":" + type
    const events = recentEvents.get(key) || []
    events.push(event)
    recentEvents.set(key, events)

    const recentInWindow = events.filter((e) => Date.now() - e.timestamp < ALERT_WINDOW_MS)

    if (recentInWindow.length >= ALERT_THRESHOLD) {
      alert(ip, type, recentInWindow.length, details)
    }
  }

  function alert(ip: string, type: SecurityEventType, count: number, details: Record<string, any>) {
    log.error("🚨 SECURITY ALERT - Potential attack detected", {
      severity: "HIGH",
      ip,
      type,
      count,
      window: ALERT_WINDOW_MS / 1000 + "s",
      message: count + " " + type + " events from " + ip + " in last minute",
      action: "Consider blocking this IP",
      details,
    })
  }

  function getIP(c: Context): string {
    return (
      c.req.header("x-forwarded-for")?.split(",")[0].trim() ||
      c.req.header("x-real-ip") ||
      c.req.header("cf-connecting-ip") ||
      "unknown"
    )
  }

  export function getStats(): { totalEvents: number; recentIPs: number; alerts: Map<string, number> } {
    let totalEvents = 0
    const alerts = new Map<string, number>()

    for (const [key, events] of recentEvents) {
      totalEvents += events.length
      const recentInWindow = events.filter((e) => Date.now() - e.timestamp < ALERT_WINDOW_MS)
      if (recentInWindow.length >= ALERT_THRESHOLD) {
        alerts.set(key, recentInWindow.length)
      }
    }

    return {
      totalEvents,
      recentIPs: recentEvents.size,
      alerts,
    }
  }
}

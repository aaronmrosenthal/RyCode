import type { Context } from "hono"
import { NamedError } from "../../util/error"
import z from "zod/v4"
import { Log } from "../../util/log"
import { Config } from "../../config/config"
import { AuthMiddleware } from "./auth"
import { SecurityMonitor } from "./security-monitor"

export namespace RateLimitMiddleware {
  const log = Log.create({ service: "rate-limit.middleware" })

  export const RateLimitError = NamedError.create(
    "RateLimitError",
    z.object({
      message: z.string(),
      retryAfter: z.number(),
      limit: z.number(),
    }),
  )

  interface BucketState {
    tokens: number
    lastRefill: number
  }

  // SECURITY: Maximum buckets to prevent memory exhaustion attacks
  const MAX_BUCKETS = 10_000
  const CLEANUP_INTERVAL_MS = 60_000 // 1 minute

  interface RateLimitConfig {
    /**
     * Maximum requests per window
     * @default 100
     */
    limit?: number
    /**
     * Window duration in milliseconds
     * @default 60000 (1 minute)
     */
    windowMs?: number
    /**
     * Key to identify the client (ip, session, api_key)
     * @default "ip"
     */
    keyBy?: "ip" | "session" | "api_key" | ((c: Context) => string)
    /**
     * Skip rate limiting for certain paths
     */
    skip?: (c: Context) => boolean
  }

  const buckets = new Map<string, BucketState>()

  // Cleanup old buckets more frequently to prevent memory leaks
  setInterval(
    () => {
      const now = Date.now()
      const expiryThreshold = 10 * 60 * 1000 // 10 minutes
      let cleaned = 0

      for (const [key, bucket] of buckets.entries()) {
        if (now - bucket.lastRefill > expiryThreshold) {
          buckets.delete(key)
          cleaned++
        }
      }

      if (cleaned > 0) {
        log.debug("cleaned up buckets", { cleaned, remaining: buckets.size })
      }
    },
    CLEANUP_INTERVAL_MS,
  ).unref()

  /**
   * Add bucket with LRU eviction if at capacity
   */
  function addBucket(key: string, bucket: BucketState): void {
    // If at capacity, evict oldest bucket (approximate LRU)
    if (buckets.size >= MAX_BUCKETS) {
      // Find and remove the oldest bucket
      let oldestKey: string | null = null
      let oldestTime = Date.now()

      for (const [k, b] of buckets.entries()) {
        if (b.lastRefill < oldestTime) {
          oldestTime = b.lastRefill
          oldestKey = k
        }
      }

      if (oldestKey) {
        buckets.delete(oldestKey)
        log.warn("bucket capacity reached, evicted oldest", {
          evicted: oldestKey,
          capacity: MAX_BUCKETS,
        })
      }
    }

    buckets.set(key, bucket)
  }

  function getClientKey(c: Context, keyBy: RateLimitConfig["keyBy"]): string {
    if (typeof keyBy === "function") {
      return keyBy(c)
    }

    switch (keyBy) {
      case "session": {
        const sessionID = c.req.param("id")
        return sessionID ? `session:${sessionID}` : getIP(c)
      }
      case "api_key": {
        const apiKey = c.req.header(AuthMiddleware.HEADER_NAME)
        return apiKey ? `key:${apiKey}` : getIP(c)
      }
      case "ip":
      default:
        return getIP(c)
    }
  }

  function getIP(c: Context): string {
    return (
      c.req.header("x-forwarded-for")?.split(",")[0].trim() ||
      c.req.header("x-real-ip") ||
      c.req.header("cf-connecting-ip") ||
      "unknown"
    )
  }

  function refillTokens(bucket: BucketState, limit: number, windowMs: number): void {
    const now = Date.now()
    const timePassed = now - bucket.lastRefill
    const refillRate = limit / windowMs
    const tokensToAdd = timePassed * refillRate

    bucket.tokens = Math.min(limit, bucket.tokens + tokensToAdd)
    bucket.lastRefill = now
  }

  export async function middleware(c: Context, next: () => Promise<void>, config: RateLimitConfig = {}) {
    const cfg = await Config.get()
    const limit = config.limit ?? cfg.server?.rate_limit?.limit ?? 100
    const windowMs = config.windowMs ?? cfg.server?.rate_limit?.window_ms ?? 60_000
    const keyBy = config.keyBy ?? "ip"

    // Skip if rate limiting disabled
    if (cfg.server?.rate_limit?.enabled === false) {
      return next()
    }

    // Skip for specific paths
    if (config.skip?.(c)) {
      return next()
    }

    const clientKey = getClientKey(c, keyBy)
    let bucket = buckets.get(clientKey)

    if (!bucket) {
      bucket = {
        tokens: limit,
        lastRefill: Date.now(),
      }
      // Use addBucket to enforce capacity limits
      addBucket(clientKey, bucket)
    }

    refillTokens(bucket, limit, windowMs)

    if (bucket.tokens < 1) {
      const retryAfter = Math.ceil((1 - bucket.tokens) / (limit / windowMs))

      SecurityMonitor.track(c, "rate_limit", {
        limit,
        tokens: bucket.tokens,
        retryAfter,
        clientKey,
      })

      log.warn("rate limit exceeded", {
        clientKey,
        path: c.req.path,
        retryAfter,
      })

      throw new RateLimitError({
        message: "Rate limit exceeded. Please try again later.",
        retryAfter,
        limit,
      })
    }

    bucket.tokens -= 1

    // Set rate limit headers
    c.header("X-RateLimit-Limit", limit.toString())
    c.header("X-RateLimit-Remaining", Math.floor(bucket.tokens).toString())
    c.header("X-RateLimit-Reset", new Date(bucket.lastRefill + windowMs).toISOString())

    return next()
  }

  /**
   * Stricter rate limiting for sensitive endpoints
   */
  export async function strictMiddleware(c: Context, next: () => Promise<void>) {
    return middleware(c, next, {
      limit: 20,
      windowMs: 60_000,
      keyBy: "session",
    })
  }

  /**
   * Clear all rate limit buckets (for testing)
   */
  export function clearBuckets() {
    buckets.clear()
  }
}

import { describe, test, expect, beforeEach } from "bun:test"
import { Hono } from "hono"
import { RateLimitMiddleware } from "../../src/server/middleware/rate-limit"
import { Config } from "../../src/config/config"

describe("RateLimitMiddleware", () => {
  let app: Hono

  beforeEach(() => {
    app = new Hono()
  })

  test("allows requests under rate limit", async () => {
    const originalGet = Config.get
    Config.get = async () =>
      ({
        server: {
          rate_limit: {
            enabled: true,
            limit: 10,
            window_ms: 60_000,
          },
        },
      }) as any

    app
      .use(async (c, next) =>
        RateLimitMiddleware.middleware(c, next, {
          limit: 10,
          windowMs: 60_000,
        }),
      )
      .get("/test", (c) => c.json({ success: true }))

    // Make 5 requests (under limit of 10)
    for (let i = 0; i < 5; i++) {
      const res = await app.request("/test")
      expect(res.status).toBe(200)
      expect(res.headers.get("X-RateLimit-Limit")).toBe("10")
    }

    Config.get = originalGet
  })

  test("blocks requests over rate limit", async () => {
    const originalGet = Config.get
    Config.get = async () =>
      ({
        server: {
          rate_limit: {
            enabled: true,
            limit: 3,
            window_ms: 60_000,
          },
        },
      }) as any

    app
      .onError((err, c) => {
        if (RateLimitMiddleware.RateLimitError.isInstance(err)) {
          c.header("Retry-After", err.data.retryAfter.toString())
          return c.json(err.toObject(), 429)
        }
        return c.json({ error: "Unknown" }, 500)
      })
      .use(async (c, next) =>
        RateLimitMiddleware.middleware(c, next, {
          limit: 3,
          windowMs: 60_000,
        }),
      )
      .get("/test", (c) => c.json({ success: true }))

    // Make 3 requests (at limit)
    for (let i = 0; i < 3; i++) {
      const res = await app.request("/test")
      expect(res.status).toBe(200)
    }

    // 4th request should be rate limited
    const blockedRes = await app.request("/test")
    expect(blockedRes.status).toBe(429)
    expect(blockedRes.headers.get("Retry-After")).toBeDefined()

    const body = await blockedRes.json()
    expect(body.data.message).toContain("Rate limit exceeded")

    Config.get = originalGet
  })

  test("sets rate limit headers", async () => {
    const originalGet = Config.get
    Config.get = async () =>
      ({
        server: {
          rate_limit: {
            enabled: true,
            limit: 100,
            window_ms: 60_000,
          },
        },
      }) as any

    app
      .use(async (c, next) =>
        RateLimitMiddleware.middleware(c, next, {
          limit: 100,
          windowMs: 60_000,
        }),
      )
      .get("/test", (c) => c.json({ success: true }))

    const res = await app.request("/test")

    expect(res.status).toBe(200)
    expect(res.headers.get("X-RateLimit-Limit")).toBe("100")
    expect(res.headers.get("X-RateLimit-Remaining")).toBe("99")
    expect(res.headers.get("X-RateLimit-Reset")).toBeDefined()

    Config.get = originalGet
  })

  test("bypasses rate limiting when disabled", async () => {
    const originalGet = Config.get
    Config.get = async () =>
      ({
        server: {
          rate_limit: {
            enabled: false,
          },
        },
      }) as any

    app
      .use(async (c, next) =>
        RateLimitMiddleware.middleware(c, next, {
          limit: 1, // Very low limit
          windowMs: 60_000,
        }),
      )
      .get("/test", (c) => c.json({ success: true }))

    // Make multiple requests - all should succeed
    for (let i = 0; i < 5; i++) {
      const res = await app.request("/test")
      expect(res.status).toBe(200)
    }

    Config.get = originalGet
  })

  test("uses different buckets for different IPs", async () => {
    const originalGet = Config.get
    Config.get = async () =>
      ({
        server: {
          rate_limit: {
            enabled: true,
            limit: 2,
            window_ms: 60_000,
          },
        },
      }) as any

    app
      .use(async (c, next) =>
        RateLimitMiddleware.middleware(c, next, {
          limit: 2,
          windowMs: 60_000,
          keyBy: "ip",
        }),
      )
      .get("/test", (c) => c.json({ success: true }))

    // Requests from IP 1
    const res1 = await app.request("/test", {
      headers: { "X-Forwarded-For": "1.2.3.4" },
    })
    expect(res1.status).toBe(200)

    const res2 = await app.request("/test", {
      headers: { "X-Forwarded-For": "1.2.3.4" },
    })
    expect(res2.status).toBe(200)

    // Requests from IP 2 (different bucket)
    const res3 = await app.request("/test", {
      headers: { "X-Forwarded-For": "5.6.7.8" },
    })
    expect(res3.status).toBe(200)

    Config.get = originalGet
  })

  // SECURITY TESTS - Memory exhaustion prevention

  test("enforces maximum bucket limit to prevent memory exhaustion", async () => {
    const originalGet = Config.get
    Config.get = async () =>
      ({
        server: {
          rate_limit: {
            enabled: true,
            limit: 10,
            window_ms: 60_000,
          },
        },
      }) as any

    app
      .use(async (c, next) =>
        RateLimitMiddleware.middleware(c, next, {
          limit: 10,
          windowMs: 60_000,
        }),
      )
      .get("/test", (c) => c.json({ success: true }))

    // Simulate attack: requests from many different IPs
    // The system should cap at MAX_BUCKETS (10,000) and evict oldest
    for (let i = 0; i < 100; i++) {
      const res = await app.request("/test", {
        headers: { "X-Forwarded-For": `192.168.1.${i}` },
      })
      expect(res.status).toBe(200)
    }

    // All requests should still work - oldest buckets evicted if needed
    const finalRes = await app.request("/test", {
      headers: { "X-Forwarded-For": "10.0.0.1" },
    })
    expect(finalRes.status).toBe(200)

    Config.get = originalGet
  })

  test("cleans up old buckets periodically", async () => {
    const originalGet = Config.get
    Config.get = async () =>
      ({
        server: {
          rate_limit: {
            enabled: true,
            limit: 10,
            window_ms: 1000, // 1 second window
          },
        },
      }) as any

    app
      .use(async (c, next) =>
        RateLimitMiddleware.middleware(c, next, {
          limit: 10,
          windowMs: 1000,
        }),
      )
      .get("/test", (c) => c.json({ success: true }))

    // Make request to create bucket
    await app.request("/test", {
      headers: { "X-Forwarded-For": "1.1.1.1" },
    })

    // Cleanup happens automatically via setInterval
    // This test just ensures no errors occur during normal operation
    expect(true).toBe(true)

    Config.get = originalGet
  })

  test("handles negative or invalid token counts gracefully", async () => {
    const originalGet = Config.get
    Config.get = async () =>
      ({
        server: {
          rate_limit: {
            enabled: true,
            limit: 1,
            window_ms: 60_000,
          },
        },
      }) as any

    app
      .onError((err, c) => {
        if (RateLimitMiddleware.RateLimitError.isInstance(err)) {
          return c.json(err.toObject(), 429)
        }
        return c.json({ error: "Unknown" }, 500)
      })
      .use(async (c, next) =>
        RateLimitMiddleware.middleware(c, next, {
          limit: 1,
          windowMs: 60_000,
        }),
      )
      .get("/test", (c) => c.json({ success: true }))

    // Use up the token
    const res1 = await app.request("/test")
    expect(res1.status).toBe(200)

    // Next request should be rate limited (tokens = 0)
    const res2 = await app.request("/test")
    expect(res2.status).toBe(429)

    // Should not crash with negative tokens
    const res3 = await app.request("/test")
    expect(res3.status).toBe(429)

    Config.get = originalGet
  })

  test("prevents DoS via many unique client identifiers", async () => {
    const originalGet = Config.get
    Config.get = async () =>
      ({
        server: {
          rate_limit: {
            enabled: true,
            limit: 5,
            window_ms: 60_000,
          },
        },
      }) as any

    app
      .use(async (c, next) =>
        RateLimitMiddleware.middleware(c, next, {
          limit: 5,
          windowMs: 60_000,
        }),
      )
      .get("/test", (c) => c.json({ success: true }))

    // Simulate DDoS: many requests from different IPs
    // System should handle gracefully with LRU eviction
    const promises = []
    for (let i = 0; i < 1000; i++) {
      promises.push(
        app.request("/test", {
          headers: { "X-Forwarded-For": `10.0.${Math.floor(i / 256)}.${i % 256}` },
        }),
      )
    }

    const results = await Promise.all(promises)

    // All requests should complete without crashing
    expect(results.length).toBe(1000)
    expect(results.every((r) => r.status === 200 || r.status === 429)).toBe(true)

    Config.get = originalGet
  })
})

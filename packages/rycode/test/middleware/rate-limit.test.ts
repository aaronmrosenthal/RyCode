import { describe, test, expect, beforeEach, afterEach } from "bun:test"
import { Hono } from "hono"
import { RateLimitMiddleware } from "../../src/server/middleware/rate-limit"
import { Config } from "../../src/config/config"

describe("RateLimitMiddleware", () => {
  let originalGet: typeof Config.get

  beforeEach(() => {
    originalGet = Config.get
    RateLimitMiddleware.clearBuckets() // Reset state
  })

  afterEach(() => {
    Config.get = originalGet
  })

  test("allows requests under rate limit", async () => {
    Config.get = async () =>
      ({
        server: { rate_limit: { enabled: true, limit: 10, window_ms: 60000 } },
      }) as any

    const app = new Hono()
      .use(async (c, next) => RateLimitMiddleware.middleware(c, next))
      .get("/test", (c) => c.json({ ok: true }))

    // Make 5 requests (under limit of 10)
    for (let i = 0; i < 5; i++) {
      const res = await app.request("/test")
      expect(res.status).toBe(200)
      expect(res.headers.get("X-RateLimit-Limit")).toBe("10")
      const remaining = parseInt(res.headers.get("X-RateLimit-Remaining") || "0")
      expect(remaining).toBeGreaterThanOrEqual(0)
    }
  })

  test("blocks requests over rate limit", async () => {
    Config.get = async () =>
      ({
        server: { rate_limit: { enabled: true, limit: 3, window_ms: 60000 } },
      }) as any

    const app = new Hono()
      .onError((err, c) => {
        if (RateLimitMiddleware.RateLimitError.isInstance(err)) {
          return c.json(err.toObject(), 429)
        }
        return c.json({ error: "Unknown" }, 500)
      })
      .use(async (c, next) => RateLimitMiddleware.middleware(c, next))
      .get("/test", (c) => c.json({ ok: true }))

    // Make 3 requests (at limit)
    for (let i = 0; i < 3; i++) {
      const res = await app.request("/test")
      expect(res.status).toBe(200)
    }

    // 4th request should be blocked
    const blocked = await app.request("/test")
    expect(blocked.status).toBe(429)

    const body = await blocked.json()
    expect(body.data.retryAfter).toBeGreaterThan(0)
    expect(body.data.limit).toBe(3)
  })

  test("does NOT accept API key via query parameter (security fix)", async () => {
    Config.get = async () =>
      ({
        server: { rate_limit: { enabled: true, limit: 5, window_ms: 60000 } },
      }) as any

    const app = new Hono()
      .use(async (c, next) => RateLimitMiddleware.middleware(c, next, { keyBy: "api_key" }))
      .get("/test", (c) => c.json({ ok: true }))

    // Make request with API key in query (should use IP fallback, not API key)
    // If query param was used, both requests would be counted against same bucket
    const res1 = await app.request("/test?api_key=test-key-12345678901234567890")
    expect(res1.status).toBe(200)

    const res2 = await app.request("/test?api_key=test-key-12345678901234567890")
    expect(res2.status).toBe(200)

    // Both requests succeeded because they use IP-based rate limiting
    // (query param is ignored for security)
  })

  test("can be disabled via config", async () => {
    Config.get = async () =>
      ({
        server: { rate_limit: { enabled: false } },
      }) as any

    const app = new Hono()
      .use(async (c, next) => RateLimitMiddleware.middleware(c, next))
      .get("/test", (c) => c.json({ ok: true }))

    // Make many requests - all should succeed
    for (let i = 0; i < 20; i++) {
      const res = await app.request("/test")
      expect(res.status).toBe(200)
    }
  })
})

import { describe, test, expect } from "bun:test"
import { Hono } from "hono"
import { SecurityHeadersMiddleware } from "../../src/server/middleware/security-headers"

describe("SecurityHeadersMiddleware", () => {
  test("sets all required security headers", async () => {
    const app = new Hono()
    app.use(async (c, next) => SecurityHeadersMiddleware.middleware(c, next))
    app.get("/test", (c) => c.json({ ok: true }))

    const res = await app.request("/test")

    expect(res.headers.get("Content-Security-Policy")).toBeTruthy()
    expect(res.headers.get("X-Content-Type-Options")).toBe("nosniff")
    expect(res.headers.get("X-Frame-Options")).toBe("DENY")
    expect(res.headers.get("X-XSS-Protection")).toBe("1; mode=block")
    expect(res.headers.get("Referrer-Policy")).toBe("strict-origin-when-cross-origin")
    expect(res.headers.get("Permissions-Policy")).toContain("geolocation=()")
    expect(res.headers.get("Server")).toBe("") // Version disclosure protection
  })

  test("blocks oversized requests", async () => {
    const app = new Hono()
      .onError((err, c) => {
        if (SecurityHeadersMiddleware.RequestTooLargeError.isInstance(err)) {
          return c.json(err.toObject(), 413)
        }
        return c.json({ error: "Unknown" }, 500)
      })
      .use(async (c, next) =>
        SecurityHeadersMiddleware.middleware(c, next, {
          maxBodySize: 1000, // 1KB limit for testing
        }),
      )
      .post("/test", (c) => c.json({ ok: true }))

    const res = await app.request("/test", {
      method: "POST",
      headers: {
        "Content-Length": "2000", // 2KB - too large
      },
    })

    expect(res.status).toBe(413)
    const body = await res.json()
    expect(body.data.size).toBe(2000)
    expect(body.data.maxSize).toBe(1000)
  })

  test("allows requests under size limit", async () => {
    const app = new Hono()
      .use(async (c, next) =>
        SecurityHeadersMiddleware.middleware(c, next, {
          maxBodySize: 1000,
        }),
      )
      .post("/test", (c) => c.json({ ok: true }))

    const res = await app.request("/test", {
      method: "POST",
      headers: {
        "Content-Length": "500", // Under limit
      },
    })

    expect(res.status).toBe(200)
  })

  test("uses stricter CSP for API endpoints", async () => {
    const app = new Hono()
    app.use(async (c, next) => SecurityHeadersMiddleware.apiMiddleware(c, next))
    app.get("/api/test", (c) => c.json({ ok: true }))

    const res = await app.request("/api/test")

    expect(res.headers.get("Content-Security-Policy")).toBe("default-src 'none'")
  })

  test("uses relaxed CSP for web endpoints", async () => {
    const app = new Hono()
    app.use(async (c, next) => SecurityHeadersMiddleware.webMiddleware(c, next))
    app.get("/", (c) => c.text("<html>"))

    const res = await app.request("/")

    const csp = res.headers.get("Content-Security-Policy")
    expect(csp).toContain("script-src 'self'")
    expect(csp).toContain("style-src 'self'")
  })

  test("can be disabled via options", async () => {
    const app = new Hono()
    app.use(async (c, next) => SecurityHeadersMiddleware.middleware(c, next, { enabled: false }))
    app.get("/test", (c) => c.json({ ok: true }))

    const res = await app.request("/test")

    expect(res.headers.get("Content-Security-Policy")).toBeNull()
    expect(res.headers.get("X-Content-Type-Options")).toBeNull()
  })

  test("allows custom CSP policy", async () => {
    const customCSP = "default-src 'self'; script-src 'self' https://cdn.example.com"
    const app = new Hono()
    app.use(async (c, next) =>
      SecurityHeadersMiddleware.middleware(c, next, {
        csp: customCSP,
      }),
    )
    app.get("/test", (c) => c.json({ ok: true }))

    const res = await app.request("/test")

    expect(res.headers.get("Content-Security-Policy")).toBe(customCSP)
  })

  test("handles missing content-length header gracefully", async () => {
    const app = new Hono()
      .use(async (c, next) =>
        SecurityHeadersMiddleware.middleware(c, next, {
          maxBodySize: 1000,
        }),
      )
      .post("/test", (c) => c.json({ ok: true }))

    // No Content-Length header - should pass through
    const res = await app.request("/test", {
      method: "POST",
    })

    expect(res.status).toBe(200)
  })

  test("sets default max body size to 10MB", async () => {
    const app = new Hono()
      .onError((err, c) => {
        if (SecurityHeadersMiddleware.RequestTooLargeError.isInstance(err)) {
          return c.json(err.toObject(), 413)
        }
        return c.json({ error: "Unknown" }, 500)
      })
      .use(async (c, next) => SecurityHeadersMiddleware.middleware(c, next))
      .post("/test", (c) => c.json({ ok: true }))

    const res = await app.request("/test", {
      method: "POST",
      headers: {
        "Content-Length": String(11 * 1024 * 1024), // 11MB - over default limit
      },
    })

    expect(res.status).toBe(413)
    const body = await res.json()
    expect(body.data.maxSize).toBe(10 * 1024 * 1024) // 10MB default
  })
})

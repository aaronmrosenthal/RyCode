import { describe, test, expect, beforeEach } from "bun:test"
import { Hono } from "hono"
import { AuthMiddleware } from "../../src/server/middleware/auth"
import { Config } from "../../src/config/config"

describe("AuthMiddleware", () => {
  let app: Hono

  beforeEach(() => {
    app = new Hono()
  })

  test("allows requests when auth is disabled", async () => {
    // Mock config without auth
    const originalGet = Config.get
    Config.get = async () =>
      ({
        server: {
          require_auth: false,
        },
      }) as any

    app.use(async (c, next) => AuthMiddleware.middleware(c, next))
    app.get("/test", (c) => c.json({ success: true }))

    const res = await app.request("/test")

    expect(res.status).toBe(200)
    expect(await res.json()).toEqual({ success: true })

    Config.get = originalGet
  })

  test("blocks requests without API key when auth is enabled", async () => {
    const originalGet = Config.get
    Config.get = async () =>
      ({
        server: {
          require_auth: true,
          api_keys: ["test-key-123"],
        },
      }) as any

    app
      .onError((err, c) => {
        if (AuthMiddleware.UnauthorizedError.isInstance(err)) {
          return c.json(err.toObject(), 401)
        }
        return c.json({ error: "Unknown" }, 500)
      })
      .use(async (c, next) => AuthMiddleware.middleware(c, next))
      .get("/test", (c) => c.json({ success: true }))

    const res = await app.request("/test")

    expect(res.status).toBe(401)
    const body = await res.json()
    expect(body.data.message).toContain("Missing API key")

    Config.get = originalGet
  })

  test("allows requests with valid API key in header", async () => {
    const originalGet = Config.get
    Config.get = async () =>
      ({
        server: {
          require_auth: true,
          api_keys: ["test-key-123"],
        },
      }) as any

    app
      .onError((err, c) => {
        if (AuthMiddleware.UnauthorizedError.isInstance(err)) {
          return c.json(err.toObject(), 401)
        }
        return c.json({ error: "Unknown" }, 500)
      })
      .use(async (c, next) => AuthMiddleware.middleware(c, next))
      .get("/test", (c) => c.json({ success: true }))

    const res = await app.request("/test", {
      headers: {
        "X-OpenCode-API-Key": "test-key-123",
      },
    })

    expect(res.status).toBe(200)
    expect(await res.json()).toEqual({ success: true })

    Config.get = originalGet
  })

  test("allows requests with valid API key in query parameter", async () => {
    const originalGet = Config.get
    Config.get = async () =>
      ({
        server: {
          require_auth: true,
          api_keys: ["test-key-456"],
        },
      }) as any

    app
      .onError((err, c) => {
        if (AuthMiddleware.UnauthorizedError.isInstance(err)) {
          return c.json(err.toObject(), 401)
        }
        return c.json({ error: "Unknown" }, 500)
      })
      .use(async (c, next) => AuthMiddleware.middleware(c, next))
      .get("/test", (c) => c.json({ success: true }))

    const res = await app.request("/test?api_key=test-key-456")

    expect(res.status).toBe(200)
    expect(await res.json()).toEqual({ success: true })

    Config.get = originalGet
  })

  test("blocks requests with invalid API key", async () => {
    const originalGet = Config.get
    Config.get = async () =>
      ({
        server: {
          require_auth: true,
          api_keys: ["valid-key"],
        },
      }) as any

    app
      .onError((err, c) => {
        if (AuthMiddleware.UnauthorizedError.isInstance(err)) {
          return c.json(err.toObject(), 401)
        }
        return c.json({ error: "Unknown" }, 500)
      })
      .use(async (c, next) => AuthMiddleware.middleware(c, next))
      .get("/test", (c) => c.json({ success: true }))

    const res = await app.request("/test", {
      headers: {
        "X-OpenCode-API-Key": "invalid-key",
      },
    })

    expect(res.status).toBe(401)
    const body = await res.json()
    expect(body.data.message).toContain("Invalid API key")

    Config.get = originalGet
  })

  test("bypasses auth for public endpoints", async () => {
    const originalGet = Config.get
    Config.get = async () =>
      ({
        server: {
          require_auth: true,
          api_keys: ["test-key"],
        },
      }) as any

    app
      .use(async (c, next) =>
        AuthMiddleware.middleware(c, next, {
          publicEndpoints: ["/doc", "/config/providers"],
        }),
      )
      .get("/doc", (c) => c.json({ success: true }))
      .get("/secure", (c) => c.json({ success: true }))

    // Public endpoint should work without auth
    const publicRes = await app.request("/doc")
    expect(publicRes.status).toBe(200)

    Config.get = originalGet
  })
})

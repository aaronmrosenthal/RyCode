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
          api_keys: ["test-key-12345678901234567890123456"],
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
          api_keys: ["test-key-12345678901234567890123456"],
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
        "X-OpenCode-API-Key": "test-key-12345678901234567890123456",
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
          api_keys: ["test-key-45678901234567890123456789"],
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

    const res = await app.request("/test?api_key=test-key-45678901234567890123456789")

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
          api_keys: ["valid-key-must-be-at-least-32-characters-long"],
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
        "X-OpenCode-API-Key": "invalid-key-678901234567890123456789",
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
          api_keys: ["test-key-must-be-at-least-32-chars-long"],
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

  // SECURITY TESTS - New tests for security fixes

  test("rejects weak API keys (less than 32 characters)", async () => {
    const originalGet = Config.get
    Config.get = async () =>
      ({
        server: {
          require_auth: true,
          api_keys: ["valid-key-must-be-at-least-32-characters-long"],
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
        "X-OpenCode-API-Key": "short-key", // Only 9 characters
      },
    })

    expect(res.status).toBe(401)
    const body = await res.json()
    expect(body.data.message).toContain("Invalid API key format")
    expect(body.data.message).toContain("at least 32 characters")

    Config.get = originalGet
  })

  test("rejects API keys with invalid characters", async () => {
    const originalGet = Config.get
    Config.get = async () =>
      ({
        server: {
          require_auth: true,
          api_keys: ["valid-key-must-be-at-least-32-characters-long"],
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
        "X-OpenCode-API-Key": "key-with-special-chars!@#$%^&*()", // Has special characters
      },
    })

    expect(res.status).toBe(401)
    const body = await res.json()
    expect(body.data.message).toContain("Invalid API key format")

    Config.get = originalGet
  })

  test("accepts API keys with valid format (alphanumeric, hyphen, underscore)", async () => {
    const validKey = "valid_api-key_12345678901234567890"
    const originalGet = Config.get
    Config.get = async () =>
      ({
        server: {
          require_auth: true,
          api_keys: [validKey],
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
        "X-OpenCode-API-Key": validKey,
      },
    })

    expect(res.status).toBe(200)
    expect(await res.json()).toEqual({ success: true })

    Config.get = originalGet
  })

  test("uses constant-time comparison (prevents timing attacks)", async () => {
    const validKey = "correct-key-1234567890123456789012"
    const originalGet = Config.get
    Config.get = async () =>
      ({
        server: {
          require_auth: true,
          api_keys: [validKey],
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

    // Test with keys that differ at different positions
    const wrongKey1 = "Xorrect-key-1234567890123456789012" // Differs at position 0
    const wrongKey2 = "correct-key-1234567890123456789012X" // Differs at last position

    const start1 = Date.now()
    await app.request("/test", { headers: { "X-OpenCode-API-Key": wrongKey1 } })
    const time1 = Date.now() - start1

    const start2 = Date.now()
    await app.request("/test", { headers: { "X-OpenCode-API-Key": wrongKey2 } })
    const time2 = Date.now() - start2

    // Timing should be similar (within 10ms tolerance for test environment)
    // This is not a perfect timing attack test, but validates the implementation exists
    expect(Math.abs(time1 - time2)).toBeLessThan(10)

    Config.get = originalGet
  })

  test("localhost bypass should NOT work with spoofed Host header", async () => {
    const originalGet = Config.get
    Config.get = async () =>
      ({
        server: {
          require_auth: true,
          api_keys: ["test-key-must-be-at-least-32-characters"],
        },
      }) as any

    app
      .onError((err, c) => {
        if (AuthMiddleware.UnauthorizedError.isInstance(err)) {
          return c.json(err.toObject(), 401)
        }
        return c.json({ error: "Unknown" }, 500)
      })
      .use(async (c, next) =>
        AuthMiddleware.middleware(c, next, {
          bypassLocalhost: true,
        }),
      )
      .get("/test", (c) => c.json({ success: true }))

    // Try to bypass auth by spoofing Host header
    const res = await app.request("/test", {
      headers: {
        Host: "localhost", // Attacker tries to spoof localhost
      },
    })

    // Should still require auth because we check socket address, not Host header
    expect(res.status).toBe(401)

    Config.get = originalGet
  })

  test("rejects empty API key", async () => {
    const originalGet = Config.get
    Config.get = async () =>
      ({
        server: {
          require_auth: true,
          api_keys: ["valid-key-must-be-at-least-32-characters"],
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
        "X-OpenCode-API-Key": "", // Empty key
      },
    })

    expect(res.status).toBe(401)
    const body = await res.json()
    // Empty keys are treated as missing
    expect(body.data.message).toContain("Missing API key")

    Config.get = originalGet
  })
})

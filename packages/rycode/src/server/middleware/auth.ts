import type { Context } from "hono"
import { Config } from "../../config/config"
import { NamedError } from "../../util/error"
import z from "zod/v4"
import { Log } from "../../util/log"
import crypto from "crypto"
import { APIKey } from "./api-key"
import { SecurityMonitor } from "./security-monitor"

export namespace AuthMiddleware {
  const log = Log.create({ service: "auth.middleware" })

  export const UnauthorizedError = NamedError.create(
    "UnauthorizedError",
    z.object({
      message: z.string(),
    }),
  )

  export const HEADER_NAME = "X-OpenCode-API-Key"

  /**
   * Validate API key format to prevent weak keys
   */
  function validateApiKeyFormat(key: string): boolean {
    return (
      typeof key === "string" &&
      key.length >= 32 &&
      /^[A-Za-z0-9_-]+$/.test(key)
    )
  }

  /**
   * Constant-time comparison to prevent timing attacks
   * Supports both hashed (new) and plaintext (legacy) API keys
   */
  async function validateApiKey(provided: string, validKeys: string[]): Promise<boolean> {
    for (const storedKey of validKeys) {
      try {
        if (APIKey.isHashed(storedKey)) {
          // New hashed format - use scrypt verification
          if (await APIKey.verify(provided, storedKey)) {
            return true
          }
        } else {
          // Legacy plaintext - log warning and verify with constant-time comparison
          log.warn("⚠️  Plaintext API key detected - please migrate to hashed format", {
            keyPrefix: storedKey.substring(0, 8) + "...",
          })

          // Still validate for backwards compatibility
          const maxLen = Math.max(provided.length, storedKey.length)
          const providedBuf = Buffer.alloc(maxLen)
          const keyBuf = Buffer.alloc(maxLen)
          Buffer.from(provided).copy(providedBuf)
          Buffer.from(storedKey).copy(keyBuf)

          if (crypto.timingSafeEqual(providedBuf, keyBuf)) {
            // Auto-generate migration suggestion
            const hashed = await APIKey.hash(storedKey)
            log.warn("📝 Replace plaintext key in config with hashed version:", {
              hashed,
            })
            return true
          }
        }
      } catch (err) {
        // Continue checking other keys
        continue
      }
    }

    return false
  }

  interface Options {
    /**
     * Skip authentication for localhost requests
     * @default true in development, false in production
     */
    bypassLocalhost?: boolean
    /**
     * Endpoints that don't require authentication
     */
    publicEndpoints?: string[]
  }

  export async function middleware(c: Context, next: () => Promise<void>, options: Options = {}) {
    const config = await Config.get()
    const bypassLocalhost = options.bypassLocalhost ?? !process.env.NODE_ENV?.includes("production")
    const publicEndpoints = options.publicEndpoints ?? ["/doc", "/config/providers", "/path", "/event"]

    // Skip auth for public endpoints
    if (publicEndpoints.some((endpoint) => c.req.path === endpoint || c.req.path.startsWith(endpoint))) {
      return next()
    }

    // Bypass auth if not configured
    if (!config.server?.require_auth) {
      log.debug("auth disabled", { path: c.req.path })
      return next()
    }

    // Bypass localhost in development
    if (bypassLocalhost) {
      // SECURITY: Use actual remote address, not Host header which can be spoofed
      const remoteAddress = c.env?.incoming?.socket?.remoteAddress
      const isLocalhost =
        remoteAddress === "127.0.0.1" ||
        remoteAddress === "::1" ||
        remoteAddress === "::ffff:127.0.0.1"

      if (isLocalhost) {
        log.debug("bypassing auth for localhost", { path: c.req.path, remoteAddress })
        return next()
      }
    }

    // Check API key (header only for security)
    const apiKey = c.req.header(HEADER_NAME)

    if (!apiKey) {
      SecurityMonitor.track(c, "auth_failure", { reason: "missing_key" })
      log.warn("missing api key", { path: c.req.path })
      throw new UnauthorizedError({
        message: `Missing API key. Provide ${HEADER_NAME} header (query parameters not supported for security reasons).`,
      })
    }

    // Validate API key format
    if (!validateApiKeyFormat(apiKey)) {
      SecurityMonitor.track(c, "weak_key_attempt", { keyLength: apiKey.length })
      log.warn("invalid api key format", { path: c.req.path, keyLength: apiKey.length })
      throw new UnauthorizedError({
        message: "Invalid API key format. Keys must be at least 32 characters and contain only alphanumeric, hyphen, or underscore characters.",
      })
    }

    // Validate API key using constant-time comparison
    const validKeys = config.server?.api_keys ?? []
    if (!(await validateApiKey(apiKey, validKeys))) {
      SecurityMonitor.track(c, "auth_failure", { reason: "invalid_key" })
      log.warn("invalid api key", { path: c.req.path })
      throw new UnauthorizedError({
        message: "Invalid API key",
      })
    }

    log.debug("authenticated", { path: c.req.path })
    return next()
  }
}

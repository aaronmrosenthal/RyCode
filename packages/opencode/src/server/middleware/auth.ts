import { Context } from "hono"
import { Config } from "../../config/config"
import { NamedError } from "../../util/error"
import z from "zod/v4"
import { Log } from "../../util/log"
import crypto from "crypto"

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
   */
  function validateApiKey(provided: string, validKeys: string[]): boolean {
    let isValid = false

    for (const key of validKeys) {
      try {
        // Pad to same length to prevent length-based timing attacks
        const maxLen = Math.max(provided.length, key.length)
        const providedBuf = Buffer.alloc(maxLen)
        const keyBuf = Buffer.alloc(maxLen)

        Buffer.from(provided).copy(providedBuf)
        Buffer.from(key).copy(keyBuf)

        // Constant-time comparison
        const match = crypto.timingSafeEqual(providedBuf, keyBuf)
        isValid = isValid || match
      } catch {
        // Length mismatch or invalid buffer - continue checking
        continue
      }
    }

    return isValid
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

    // Check API key
    const apiKey = c.req.header(HEADER_NAME) || c.req.query("api_key")

    if (!apiKey) {
      log.warn("missing api key", { path: c.req.path })
      throw new UnauthorizedError({
        message: `Missing API key. Provide ${HEADER_NAME} header or api_key query parameter.`,
      })
    }

    // Validate API key format
    if (!validateApiKeyFormat(apiKey)) {
      log.warn("invalid api key format", { path: c.req.path, keyLength: apiKey.length })
      throw new UnauthorizedError({
        message: "Invalid API key format. Keys must be at least 32 characters and contain only alphanumeric, hyphen, or underscore characters.",
      })
    }

    // Validate API key using constant-time comparison
    const validKeys = config.server?.api_keys ?? []
    if (!validateApiKey(apiKey, validKeys)) {
      log.warn("invalid api key", { path: c.req.path })
      throw new UnauthorizedError({
        message: "Invalid API key",
      })
    }

    log.debug("authenticated", { path: c.req.path })
    return next()
  }
}

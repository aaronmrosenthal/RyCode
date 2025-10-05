import type { Context } from "hono"
import { Log } from "../../util/log"
import { NamedError } from "../../util/error"
import z from "zod/v4"

export namespace SecurityHeadersMiddleware {
  const log = Log.create({ service: "security-headers.middleware" })

  export const RequestTooLargeError = NamedError.create(
    "RequestTooLargeError",
    z.object({
      message: z.string(),
      size: z.number(),
      maxSize: z.number(),
    }),
  )

  interface Options {
    /**
     * Maximum request body size in bytes
     * @default 10485760 (10MB)
     */
    maxBodySize?: number
    /**
     * Content Security Policy directives
     * @default "default-src 'self'"
     */
    csp?: string
    /**
     * Enable security headers
     * @default true
     */
    enabled?: boolean
  }

  export async function middleware(c: Context, next: () => Promise<void>, options: Options = {}) {
    const maxBodySize = options.maxBodySize ?? 10 * 1024 * 1024 // 10MB
    const csp = options.csp ?? "default-src 'self'"
    const enabled = options.enabled ?? true

    if (!enabled) {
      return next()
    }

    // Check request body size
    const contentLength = c.req.header("content-length")
    if (contentLength) {
      const size = parseInt(contentLength, 10)
      if (size > maxBodySize) {
        log.warn("request too large", {
          path: c.req.path,
          size,
          maxSize: maxBodySize,
        })
        throw new RequestTooLargeError({
          message: `Request body too large. Maximum size: ${maxBodySize} bytes`,
          size,
          maxSize: maxBodySize,
        })
      }
    }

    // Set security headers
    c.header("Content-Security-Policy", csp)
    c.header("X-Content-Type-Options", "nosniff")
    c.header("X-Frame-Options", "DENY")
    c.header("X-XSS-Protection", "1; mode=block")
    c.header("Referrer-Policy", "strict-origin-when-cross-origin")
    c.header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

    // Remove server header to avoid version disclosure
    c.header("Server", "")

    return next()
  }

  /**
   * Stricter CSP for API endpoints
   */
  export async function apiMiddleware(c: Context, next: () => Promise<void>) {
    return middleware(c, next, {
      csp: "default-src 'none'", // API endpoints don't need to load resources
    })
  }

  /**
   * Relaxed CSP for web UI endpoints
   */
  export async function webMiddleware(c: Context, next: () => Promise<void>) {
    return middleware(c, next, {
      csp: "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:;",
    })
  }
}

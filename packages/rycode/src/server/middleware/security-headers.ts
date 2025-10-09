import type { Context } from "hono"
import { Log } from "../../util/log"
import { NamedError } from "../../util/error"
import z from "zod/v4"
import crypto from "crypto"
import { SecurityMonitor } from "./security-monitor"

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

  /**
   * Generate cryptographically secure CSP nonce
   */
  function generateNonce(): string {
    return crypto.randomBytes(16).toString("base64")
  }

  interface Options {
    /**
     * Maximum request body size in bytes
     * @default 10485760 (10MB)
     */
    maxBodySize?: number
    /**
     * Content Security Policy directives
     * Use {nonce} placeholder for nonce-based CSP
     * @default "default-src 'self'"
     */
    csp?: string
    /**
     * Enable security headers
     * @default true
     */
    enabled?: boolean
    /**
     * Use nonces for CSP (more secure than unsafe-inline)
     * @default true
     */
    useNonces?: boolean
  }

  export async function middleware(c: Context, next: () => Promise<void>, options: Options = {}) {
    const maxBodySize = options.maxBodySize ?? 10 * 1024 * 1024 // 10MB
    const useNonces = options.useNonces ?? true
    const enabled = options.enabled ?? true

    if (!enabled) {
      return next()
    }

    // Check request body size
    const contentLength = c.req.header("content-length")
    if (contentLength) {
      const size = parseInt(contentLength, 10)
      if (size > maxBodySize) {
        SecurityMonitor.track(c, "request_too_large", {
          size,
          maxSize: maxBodySize,
        })

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

    // Generate nonce for CSP if enabled
    let csp = options.csp ?? "default-src 'self'"
    if (useNonces) {
      const nonce = generateNonce()
      c.set("csp-nonce", nonce)
      // Replace {nonce} placeholder with actual nonce
      csp = csp.replace(/\{nonce\}/g, `'nonce-${nonce}'`)
    }

    // Set security headers
    c.header("Content-Security-Policy", csp)
    c.header("X-Content-Type-Options", "nosniff")
    c.header("X-Frame-Options", "DENY")
    c.header("X-XSS-Protection", "1; mode=block")
    c.header("Referrer-Policy", "strict-origin-when-cross-origin")
    c.header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

    // HSTS - Force HTTPS for 1 year including subdomains
    // Only set if connection is already HTTPS to avoid breaking HTTP development
    const proto = c.req.header("x-forwarded-proto") || (c.req.url.startsWith("https") ? "https" : "http")
    if (proto === "https") {
      c.header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
    }

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
   * Relaxed CSP for web UI endpoints with nonce-based inline script/style support
   */
  export async function webMiddleware(c: Context, next: () => Promise<void>) {
    return middleware(c, next, {
      csp: "default-src 'self'; script-src 'self' {nonce}; style-src 'self' {nonce}; img-src 'self' data: https:;",
      useNonces: true,
    })
  }
}

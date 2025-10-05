import { Context } from "hono"
import { Config } from "../../config/config"
import { NamedError } from "../../util/error"
import z from "zod/v4"
import { Log } from "../../util/log"

export namespace AuthMiddleware {
  const log = Log.create({ service: "auth.middleware" })

  export const UnauthorizedError = NamedError.create(
    "UnauthorizedError",
    z.object({
      message: z.string(),
    }),
  )

  export const HEADER_NAME = "X-OpenCode-API-Key"

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
      const hostname = c.req.header("host")?.split(":")[0]
      if (hostname === "localhost" || hostname === "127.0.0.1" || hostname === "::1") {
        log.debug("bypassing auth for localhost", { path: c.req.path })
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

    // Validate API key
    const validKeys = config.server?.api_keys ?? []
    if (!validKeys.includes(apiKey)) {
      log.warn("invalid api key", { path: c.req.path })
      throw new UnauthorizedError({
        message: "Invalid API key",
      })
    }

    log.debug("authenticated", { path: c.req.path })
    return next()
  }
}

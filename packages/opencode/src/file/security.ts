import path from "path"
import { NamedError } from "../util/error"
import z from "zod/v4"
import { Log } from "../util/log"
import { Instance } from "../project/instance"

export namespace FileSecurity {
  const log = Log.create({ service: "file.security" })

  export const PathTraversalError = NamedError.create(
    "PathTraversalError",
    z.object({
      requestedPath: z.string(),
      message: z.string(),
    }),
  )

  export const SensitiveFileError = NamedError.create(
    "SensitiveFileError",
    z.object({
      requestedPath: z.string(),
      message: z.string(),
    }),
  )

  /**
   * List of sensitive file patterns that should not be accessible
   */
  const SENSITIVE_PATTERNS = [
    // Credentials and secrets
    ".env",
    ".env.*",
    "*.pem",
    "*.key",
    "*.p12",
    "*.pfx",
    "*credentials*",
    "*secret*",
    "*password*",
    "id_rsa",
    "id_dsa",
    "id_ed25519",

    // System files (Unix/Linux)
    "/etc/passwd",
    "/etc/shadow",
    "/etc/hosts",
    "/etc/ssh/*",
    "/root/*",

    // System files (macOS)
    "/System/*",
    "/Library/Keychains/*",

    // System files (Windows)
    "C:\\Windows\\*",
    "C:\\Program Files\\*",

    // Cloud provider credentials
    ".aws/credentials",
    ".azure/credentials",
    ".gcp/credentials",
    "gcloud/credentials",

    // SSH
    ".ssh/*",

    // Git credentials
    ".git-credentials",
    ".netrc",

    // Database files
    "*.sqlite",
    "*.db",

    // Other sensitive configs
    "kubeconfig",
    ".kube/config",
  ]

  /**
   * Validates that a path is safe to access
   * @throws PathTraversalError if path escapes allowed directories
   * @throws SensitiveFileError if path matches sensitive file patterns
   */
  export function validatePath(requestedPath: string): string {
    const normalized = path.normalize(requestedPath)
    const resolved = path.resolve(Instance.directory, normalized)

    // Check for path traversal outside of worktree
    const worktree = path.resolve(Instance.worktree)
    const directory = path.resolve(Instance.directory)

    // Path must be within either the current directory or worktree
    const isInDirectory = resolved.startsWith(directory)
    const isInWorktree = resolved.startsWith(worktree)

    if (!isInDirectory && !isInWorktree) {
      log.warn("path traversal attempt", {
        requestedPath,
        resolved,
        directory,
        worktree,
      })
      throw new PathTraversalError({
        requestedPath,
        message: `Path '${requestedPath}' is outside allowed directories`,
      })
    }

    // Check for sensitive files
    if (isSensitiveFile(resolved)) {
      log.warn("sensitive file access attempt", {
        requestedPath,
        resolved,
      })
      throw new SensitiveFileError({
        requestedPath,
        message: `Access to sensitive file '${requestedPath}' is not allowed`,
      })
    }

    return resolved
  }

  /**
   * Checks if a path matches sensitive file patterns
   */
  function isSensitiveFile(filePath: string): boolean {
    const normalized = filePath.toLowerCase()

    for (const pattern of SENSITIVE_PATTERNS) {
      if (matchPattern(normalized, pattern.toLowerCase())) {
        return true
      }
    }

    return false
  }

  /**
   * Simple glob pattern matching
   */
  function matchPattern(filePath: string, pattern: string): boolean {
    // Exact match
    if (filePath === pattern) return true

    // Contains match for patterns with wildcards
    if (pattern.includes("*")) {
      const regex = new RegExp("^" + pattern.replace(/\*/g, ".*") + "$")
      return regex.test(filePath)
    }

    // Substring match for partial patterns
    return filePath.includes(pattern)
  }

  /**
   * Validates multiple paths in a batch
   */
  export function validatePaths(paths: string[]): string[] {
    return paths.map((p) => validatePath(p))
  }

  /**
   * Checks if a path is safe without throwing
   */
  export function isPathSafe(requestedPath: string): boolean {
    try {
      validatePath(requestedPath)
      return true
    } catch {
      return false
    }
  }

  /**
   * Get the list of sensitive patterns (for documentation/config)
   */
  export function getSensitivePatterns(): readonly string[] {
    return SENSITIVE_PATTERNS
  }
}

/**
 * Plugin Security System
 * Provides allowlist, capability-based permissions, and verification for plugins
 */

import z from "zod/v4"
import { NamedError } from "../util/error"
import { Log } from "../util/log"
import crypto from "crypto"

export namespace PluginSecurity {
  const log = Log.create({ service: "plugin.security" })

  /**
   * Plugin capabilities define what system resources a plugin can access
   */
  export const Capabilities = z.object({
    /** File system read access */
    fileSystemRead: z.boolean().default(true),
    /** File system write access */
    fileSystemWrite: z.boolean().default(false),
    /** Network/HTTP access */
    network: z.boolean().default(true),
    /** Shell command execution */
    shell: z.boolean().default(false),
    /** Environment variable access */
    env: z.boolean().default(false),
    /** Access to project metadata */
    projectMetadata: z.boolean().default(true),
    /** Access to AI client */
    aiClient: z.boolean().default(true),
  })
  export type Capabilities = z.infer<typeof Capabilities>

  /**
   * Trusted plugin configuration
   */
  export const TrustedPlugin = z.object({
    /** Package name */
    name: z.string(),
    /** Allowed versions (semver or "latest") */
    versions: z.array(z.string()).default(["latest"]),
    /** Plugin capabilities */
    capabilities: Capabilities,
    /** Optional: SHA-256 hash for integrity verification */
    hash: z.string().optional(),
    /** Optional: GPG signature for verification */
    signature: z.string().optional(),
    /** Is this a first-party official plugin? */
    official: z.boolean().default(false),
  })
  export type TrustedPlugin = z.infer<typeof TrustedPlugin>

  /**
   * Plugin security policy
   */
  export const Policy = z.object({
    /** Enforcement mode */
    mode: z.enum(["strict", "warn", "permissive"]).default("warn"),
    /** List of trusted plugins */
    trustedPlugins: z.array(TrustedPlugin).default([]),
    /** Default capabilities for untrusted plugins */
    defaultCapabilities: Capabilities,
    /** Require user approval for untrusted plugins */
    requireApproval: z.boolean().default(true),
    /** Enable integrity verification (hash checking) */
    verifyIntegrity: z.boolean().default(true),
  })
  export type Policy = z.infer<typeof Policy>

  /**
   * Official first-party plugins (always trusted)
   */
  export const OFFICIAL_PLUGINS: TrustedPlugin[] = [
    {
      name: "opencode-copilot-auth",
      versions: ["0.0.3", "latest"],
      official: true,
      capabilities: {
        fileSystemRead: true,
        fileSystemWrite: false,
        network: true,
        shell: false,
        env: true,
        projectMetadata: true,
        aiClient: true,
      },
    },
    {
      name: "opencode-anthropic-auth",
      versions: ["0.0.2", "latest"],
      official: true,
      capabilities: {
        fileSystemRead: true,
        fileSystemWrite: false,
        network: true,
        shell: false,
        env: true,
        projectMetadata: true,
        aiClient: true,
      },
    },
  ]

  /**
   * Default security policy (warn mode with official plugins trusted)
   */
  export const DEFAULT_POLICY: Policy = {
    mode: "warn",
    trustedPlugins: OFFICIAL_PLUGINS,
    defaultCapabilities: {
      fileSystemRead: true,
      fileSystemWrite: false,
      network: false,
      shell: false,
      env: false,
      projectMetadata: true,
      aiClient: false,
    },
    requireApproval: true,
    verifyIntegrity: false, // Disabled by default until we have a registry
  }

  /**
   * Errors
   */
  export const UntrustedPluginError = NamedError.create(
    "UntrustedPluginError",
    z.object({
      plugin: z.string(),
      version: z.string(),
      message: z.string(),
    }),
  )

  export const CapabilityDeniedError = NamedError.create(
    "CapabilityDeniedError",
    z.object({
      plugin: z.string(),
      capability: z.string(),
      message: z.string(),
    }),
  )

  export const IntegrityCheckFailedError = NamedError.create(
    "IntegrityCheckFailedError",
    z.object({
      plugin: z.string(),
      expected: z.string(),
      actual: z.string(),
    }),
  )

  /**
   * Check if a plugin is trusted according to policy
   */
  export function isTrusted(
    packageName: string,
    version: string,
    policy: Policy = DEFAULT_POLICY,
  ): { trusted: boolean; config?: TrustedPlugin } {
    const config = policy.trustedPlugins.find((p) => p.name === packageName)

    if (!config) {
      return { trusted: false }
    }

    // Check version compatibility
    const versionMatches =
      config.versions.includes("latest") ||
      config.versions.includes(version) ||
      config.versions.some((v) => matchesVersion(version, v))

    if (!versionMatches) {
      log.warn("plugin version not in allowlist", {
        plugin: packageName,
        version,
        allowedVersions: config.versions,
      })
      return { trusted: false }
    }

    return { trusted: true, config }
  }

  /**
   * Simple semver range matching (basic implementation)
   */
  function matchesVersion(version: string, pattern: string): boolean {
    if (pattern === "*" || pattern === "latest") return true
    if (pattern.startsWith("^")) {
      const base = pattern.slice(1)
      return version.startsWith(base.split(".")[0])
    }
    if (pattern.startsWith("~")) {
      const base = pattern.slice(1)
      const baseParts = base.split(".")
      return version.startsWith(`${baseParts[0]}.${baseParts[1]}`)
    }
    return version === pattern
  }

  /**
   * Get capabilities for a plugin
   */
  export function getCapabilities(
    packageName: string,
    version: string,
    policy: Policy = DEFAULT_POLICY,
  ): Capabilities {
    const { trusted, config } = isTrusted(packageName, version, policy)

    if (trusted && config) {
      log.info("using trusted plugin capabilities", {
        plugin: packageName,
        capabilities: config.capabilities,
      })
      return config.capabilities
    }

    log.warn("using default capabilities for untrusted plugin", {
      plugin: packageName,
      capabilities: policy.defaultCapabilities,
    })
    return policy.defaultCapabilities
  }

  /**
   * Verify plugin integrity using SHA-256 hash
   */
  export async function verifyIntegrity(
    pluginPath: string,
    expectedHash?: string,
  ): Promise<boolean> {
    if (!expectedHash) {
      log.debug("no hash provided, skipping integrity check")
      return true
    }

    try {
      const file = Bun.file(pluginPath)
      const content = await file.arrayBuffer()
      const hash = crypto.createHash("sha256").update(Buffer.from(content)).digest("hex")

      if (hash !== expectedHash) {
        log.error("integrity check failed", {
          path: pluginPath,
          expected: expectedHash,
          actual: hash,
        })
        return false
      }

      log.info("integrity check passed", { path: pluginPath })
      return true
    } catch (error) {
      log.error("integrity check error", { path: pluginPath, error })
      return false
    }
  }

  /**
   * Validate plugin can perform an action based on capabilities
   */
  export function checkCapability(
    pluginName: string,
    capability: keyof Capabilities,
    capabilities: Capabilities,
  ): void {
    if (!capabilities[capability]) {
      throw new CapabilityDeniedError({
        plugin: pluginName,
        capability,
        message: `Plugin "${pluginName}" does not have permission for: ${capability}`,
      })
    }
  }

  /**
   * Create a sandboxed plugin input with capability restrictions
   */
  export function createSandboxedInput(
    pluginName: string,
    baseInput: any,
    capabilities: Capabilities,
  ): any {
    const sandboxed: any = {
      project: capabilities.projectMetadata ? baseInput.project : undefined,
      worktree: capabilities.fileSystemRead ? baseInput.worktree : undefined,
      directory: capabilities.fileSystemRead ? baseInput.directory : undefined,
    }

    // Conditional client access
    if (capabilities.aiClient) {
      sandboxed.client = baseInput.client
    }

    // Restricted shell access
    if (capabilities.shell) {
      sandboxed.$ = baseInput.$
    } else {
      // Provide no-op shell
      sandboxed.$ = new Proxy(
        {},
        {
          get() {
            throw new CapabilityDeniedError({
              plugin: pluginName,
              capability: "shell",
              message: `Plugin "${pluginName}" does not have shell execution permission`,
            })
          },
        },
      )
    }

    // Restricted environment access
    if (!capabilities.env) {
      // Override process.env access
      Object.defineProperty(sandboxed, "env", {
        get() {
          throw new CapabilityDeniedError({
            plugin: pluginName,
            capability: "env",
            message: `Plugin "${pluginName}" does not have environment variable access`,
          })
        },
      })
    }

    return sandboxed
  }

  /**
   * Generate SHA-256 hash for a plugin file (for allowlist creation)
   */
  export async function generateHash(pluginPath: string): Promise<string> {
    const file = Bun.file(pluginPath)
    const content = await file.arrayBuffer()
    return crypto.createHash("sha256").update(Buffer.from(content)).digest("hex")
  }

  /**
   * Security audit log entry
   */
  export interface AuditEntry {
    timestamp: number
    plugin: string
    version: string
    action: "loaded" | "denied" | "capability_check"
    trusted: boolean
    capabilities?: Capabilities
    reason?: string
  }

  const auditLog: AuditEntry[] = []

  /**
   * Log security event
   */
  export function audit(entry: Omit<AuditEntry, "timestamp">): void {
    const fullEntry = {
      ...entry,
      timestamp: Date.now(),
    }
    auditLog.push(fullEntry)
    log.info("plugin security audit", fullEntry)
  }

  /**
   * Get security audit log
   */
  export function getAuditLog(): readonly AuditEntry[] {
    return auditLog
  }

  /**
   * Clear audit log (for testing)
   */
  export function clearAuditLog(): void {
    auditLog.length = 0
  }
}

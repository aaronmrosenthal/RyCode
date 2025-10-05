/**
 * Plugin Registry System
 *
 * Centralized registry of verified plugin hashes and metadata.
 * Provides a trusted source for plugin verification and discovery.
 */

import z from "zod/v4"
import { NamedError } from "../util/error"
import { Log } from "../util/log"
import path from "path"
import { existsSync } from "fs"
import { PluginSignature } from "./signature"

export namespace PluginRegistry {
  const log = Log.create({ service: "plugin.registry" })

  /**
   * Registry save error
   */
  export const RegistrySaveError = NamedError.create(
    "RegistrySaveError",
    z.object({
      filePath: z.string(),
    })
  )
  export type RegistrySaveError = InstanceType<typeof RegistrySaveError>

  /**
   * Plugin entry in the registry
   */
  export const RegistryEntry = z.object({
    /** Package name */
    name: z.string(),
    /** Version */
    version: z.string(),
    /** SHA-256 hash */
    hash: z.string().regex(/^[a-f0-9]{64}$/),
    /** Description */
    description: z.string().optional(),
    /** Author */
    author: z.string().optional(),
    /** Homepage URL */
    homepage: z.string().url().optional(),
    /** Repository URL */
    repository: z.string().url().optional(),
    /** Verified by (official/community/user) */
    verifiedBy: z.enum(["official", "community", "user"]).default("user"),
    /** Timestamp when added */
    timestamp: z.number(),
    /** Capabilities required */
    capabilities: z.record(z.string(), z.boolean()).optional(),
    /** Optional: Cryptographic signature */
    signature: PluginSignature.Signature.optional(),
  })
  export type RegistryEntry = z.infer<typeof RegistryEntry>

  /**
   * Registry data structure
   */
  export const Registry = z.object({
    /** Format version */
    version: z.string().default("1.0.0"),
    /** Registry entries */
    entries: z.array(RegistryEntry).default([]),
    /** Last updated timestamp */
    lastUpdated: z.number(),
  })
  export type Registry = z.infer<typeof Registry>

  /**
   * Registry configuration
   */
  export interface RegistryConfig {
    /** Local registry path */
    localPath?: string
    /** Remote registry URL */
    remoteUrl?: string
    /** Auto-update from remote */
    autoUpdate?: boolean
    /** Cache TTL in milliseconds */
    cacheTTL?: number
  }

  /**
   * In-memory registry cache
   */
  let registryCache: Registry | null = null
  let cacheTimestamp: number = 0

  /**
   * Default registry configuration
   */
  const DEFAULT_CONFIG: RegistryConfig = {
    localPath: path.join(process.env["HOME"] || "~", ".rycode", "plugin-registry.json"),
    remoteUrl: "https://registry.rycode.ai/plugins.json",
    autoUpdate: true,
    cacheTTL: 3600000, // 1 hour
  }

  /**
   * Load registry from local or remote source
   */
  export async function load(config: RegistryConfig = {}): Promise<Registry> {
    const cfg = { ...DEFAULT_CONFIG, ...config }

    // Check cache
    if (registryCache && cfg.cacheTTL && Date.now() - cacheTimestamp < cfg.cacheTTL) {
      log.debug("registry.load", { source: "cache" })
      return registryCache
    }

    let registry: Registry | null = null

    // Try local file first
    if (cfg.localPath && existsSync(cfg.localPath)) {
      try {
        const file = Bun.file(cfg.localPath)
      const content = await file.text()
        const data = JSON.parse(content)
        registry = Registry.parse(data)
        log.info("registry.load", { source: "local", path: cfg.localPath, entries: registry.entries.length })
      } catch (error) {
        log.warn("registry.load.local.failed", {
          error: error instanceof Error ? error.message : String(error),
        })
      }
    }

    // Try remote if configured and auto-update enabled
    if (cfg.autoUpdate && cfg.remoteUrl && (!registry || shouldUpdate(registry, cfg))) {
      try {
        const remote = await fetchRemote(cfg.remoteUrl)
        if (remote && (!registry || remote.lastUpdated > registry.lastUpdated)) {
          registry = remote
          // Save to local cache
          if (cfg.localPath) {
            await save(registry, cfg.localPath)
          }
          log.info("registry.load", { source: "remote", url: cfg.remoteUrl, entries: registry.entries.length })
        }
      } catch (error) {
        log.warn("registry.load.remote.failed", {
          error: error instanceof Error ? error.message : String(error),
        })
      }
    }

    // Create empty registry if none found
    if (!registry) {
      registry = {
        version: "1.0.0",
        entries: [],
        lastUpdated: Date.now(),
      }
      log.info("registry.load", { source: "empty" })
    }

    // Update cache
    registryCache = registry
    cacheTimestamp = Date.now()

    return registry
  }

  /**
   * Save registry to local file
   */
  export async function save(registry: Registry, filePath: string): Promise<void> {
    try {
      const dir = path.dirname(filePath)
      if (!existsSync(dir)) {
        const { mkdir } = await import("fs/promises")
        await mkdir(dir, { recursive: true })
      }

      await Bun.write(filePath, JSON.stringify(registry, null, 2))
      log.info("registry.save", { path: filePath, entries: registry.entries.length })

      // Update cache
      registryCache = registry
      cacheTimestamp = Date.now()
    } catch (error) {
      log.error("registry.save.failed", {
        error: error instanceof Error ? error.message : String(error),
      })
      throw new RegistrySaveError({ filePath })
    }
  }

  /**
   * Fetch registry from remote URL
   */
  async function fetchRemote(url: string): Promise<Registry | null> {
    try {
      const response = await fetch(url, {
        headers: {
          "User-Agent": "RyCode-Plugin-Registry/1.0",
          "Accept": "application/json",
        },
      })

      if (!response.ok) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`)
      }

      const data = await response.json()
      return Registry.parse(data)
    } catch (error) {
      log.error("registry.fetch.failed", {
        url,
        error: error instanceof Error ? error.message : String(error),
      })
      return null
    }
  }

  /**
   * Check if registry should be updated
   */
  function shouldUpdate(registry: Registry, config: RegistryConfig): boolean {
    if (!config.cacheTTL) return false
    return Date.now() - registry.lastUpdated > config.cacheTTL
  }

  /**
   * Find entry in registry
   */
  export async function find(
    name: string,
    version: string,
    config?: RegistryConfig
  ): Promise<RegistryEntry | null> {
    const registry = await load(config)
    return registry.entries.find(e => e.name === name && e.version === version) || null
  }

  /**
   * Find all entries for a plugin (all versions)
   */
  export async function findAll(
    name: string,
    config?: RegistryConfig
  ): Promise<RegistryEntry[]> {
    const registry = await load(config)
    return registry.entries.filter(e => e.name === name)
  }

  /**
   * Add entry to registry
   */
  export async function add(
    entry: Omit<RegistryEntry, "timestamp">,
    config?: RegistryConfig
  ): Promise<void> {
    const cfg = { ...DEFAULT_CONFIG, ...config }
    const registry = await load(cfg)

    // Check for duplicate
    const existing = registry.entries.findIndex(
      e => e.name === entry.name && e.version === entry.version
    )

    const newEntry: RegistryEntry = {
      ...entry,
      timestamp: Date.now(),
    }

    if (existing >= 0) {
      // Update existing entry
      registry.entries[existing] = newEntry
      log.info("registry.update", { name: entry.name, version: entry.version })
    } else {
      // Add new entry
      registry.entries.push(newEntry)
      log.info("registry.add", { name: entry.name, version: entry.version })
    }

    registry.lastUpdated = Date.now()

    // Save to local file
    if (cfg.localPath) {
      await save(registry, cfg.localPath)
    }
  }

  /**
   * Remove entry from registry
   */
  export async function remove(
    name: string,
    version: string,
    config?: RegistryConfig
  ): Promise<boolean> {
    const cfg = { ...DEFAULT_CONFIG, ...config }
    const registry = await load(cfg)

    const index = registry.entries.findIndex(
      e => e.name === name && e.version === version
    )

    if (index < 0) {
      return false
    }

    registry.entries.splice(index, 1)
    registry.lastUpdated = Date.now()

    log.info("registry.remove", { name, version })

    // Save to local file
    if (cfg.localPath) {
      await save(registry, cfg.localPath)
    }

    return true
  }

  /**
   * Verify hash against registry
   */
  export async function verify(
    name: string,
    version: string,
    hash: string,
    config?: RegistryConfig
  ): Promise<{ verified: boolean; entry: RegistryEntry | null }> {
    const entry = await find(name, version, config)

    if (!entry) {
      log.debug("registry.verify.not_found", { name, version })
      return { verified: false, entry: null }
    }

    const verified = entry.hash === hash

    log.info("registry.verify", {
      name,
      version,
      verified,
      verifiedBy: entry.verifiedBy,
    })

    return { verified, entry }
  }

  /**
   * Verify plugin with both hash and signature
   */
  export async function verifyComplete(
    name: string,
    version: string,
    filePath: string,
    config?: RegistryConfig
  ): Promise<{
    hashVerified: boolean
    signatureVerified: boolean
    entry: RegistryEntry | null
    signatureError?: string
  }> {
    const { PluginSecurity } = await import("./security")

    // Find entry
    const entry = await find(name, version, config)

    if (!entry) {
      return {
        hashVerified: false,
        signatureVerified: false,
        entry: null,
      }
    }

    // Verify hash
    const actualHash = await PluginSecurity.generateHash(filePath)
    const hashVerified = entry.hash === actualHash

    // Verify signature if present
    let signatureVerified = false
    let signatureError: string | undefined

    if (entry.signature) {
      const sigResult = await PluginSignature.verifyCryptoSignature(
        filePath,
        entry.signature,
        entry.signature.publicKey || ""
      )

      signatureVerified = sigResult.valid
      signatureError = sigResult.error

      log.info("complete verification", {
        name,
        version,
        hashVerified,
        signatureVerified,
      })
    } else {
      // No signature in entry
      signatureVerified = true // Don't fail if signature is optional
      log.debug("no signature in entry", { name, version })
    }

    return {
      hashVerified,
      signatureVerified,
      entry,
      signatureError,
    }
  }

  /**
   * Search registry by name pattern
   */
  export async function search(
    pattern: string,
    config?: RegistryConfig
  ): Promise<RegistryEntry[]> {
    const registry = await load(config)
    const regex = new RegExp(pattern, "i")

    return registry.entries.filter(e =>
      regex.test(e.name) ||
      (e.description && regex.test(e.description))
    )
  }

  /**
   * Clear cache
   */
  export function clearCache(): void {
    registryCache = null
    cacheTimestamp = 0
    log.debug("registry.cache.cleared")
  }

  /**
   * Get registry statistics
   */
  export async function stats(config?: RegistryConfig): Promise<{
    totalEntries: number
    officialCount: number
    communityCount: number
    userCount: number
    uniquePlugins: number
  }> {
    const registry = await load(config)

    const official = registry.entries.filter(e => e.verifiedBy === "official").length
    const community = registry.entries.filter(e => e.verifiedBy === "community").length
    const user = registry.entries.filter(e => e.verifiedBy === "user").length
    const uniqueNames = new Set(registry.entries.map(e => e.name))

    return {
      totalEntries: registry.entries.length,
      officialCount: official,
      communityCount: community,
      userCount: user,
      uniquePlugins: uniqueNames.size,
    }
  }
}

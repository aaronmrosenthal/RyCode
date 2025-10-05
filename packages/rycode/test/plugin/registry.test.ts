import { describe, expect, test, beforeEach } from "bun:test"
import { PluginRegistry } from "../../src/plugin/registry"
import { PluginSecurity } from "../../src/plugin/security"
import path from "path"
import { tmpdir } from "../fixture/fixture"

/**
 * Tests for Plugin Registry System
 */

describe("PluginRegistry", () => {
  beforeEach(() => {
    PluginRegistry.clearCache()
  })

  describe("Registry Management", () => {
    test("should create empty registry", async () => {
      await using tmp = await tmpdir()
      const registryPath = path.join(tmp.path, "registry.json")

      const registry = await PluginRegistry.load({
        localPath: registryPath,
        autoUpdate: false,
      })

      expect(registry.version).toBe("1.0.0")
      expect(registry.entries).toEqual([])
      expect(registry.lastUpdated).toBeGreaterThan(0)
    })

    test("should add entry to registry", async () => {
      await using tmp = await tmpdir()
      const registryPath = path.join(tmp.path, "registry.json")
      const pluginPath = path.join(tmp.path, "test-plugin.js")
      await Bun.write(pluginPath, "console.log('test')")

      const hash = await PluginSecurity.generateHash(pluginPath)

      await PluginRegistry.add(
        {
          name: "test-plugin",
          version: "1.0.0",
          hash,
          description: "Test plugin",
          verifiedBy: "user",
        },
        { localPath: registryPath, autoUpdate: false }
      )

      const registry = await PluginRegistry.load({
        localPath: registryPath,
        autoUpdate: false,
      })

      expect(registry.entries.length).toBe(1)
      expect(registry.entries[0].name).toBe("test-plugin")
      expect(registry.entries[0].version).toBe("1.0.0")
      expect(registry.entries[0].hash).toBe(hash)
    })

    test("should update existing entry", async () => {
      await using tmp = await tmpdir()
      const registryPath = path.join(tmp.path, "registry.json")

      // Add first entry
      await PluginRegistry.add(
        {
          name: "test-plugin",
          version: "1.0.0",
          hash: "a".repeat(64),
          description: "Old description",
          verifiedBy: "user",
        },
        { localPath: registryPath, autoUpdate: false }
      )

      // Update with new hash
      await PluginRegistry.add(
        {
          name: "test-plugin",
          version: "1.0.0",
          hash: "b".repeat(64),
          description: "New description",
          verifiedBy: "community",
        },
        { localPath: registryPath, autoUpdate: false }
      )

      const registry = await PluginRegistry.load({
        localPath: registryPath,
        autoUpdate: false,
      })

      expect(registry.entries.length).toBe(1)
      expect(registry.entries[0].hash).toBe("b".repeat(64))
      expect(registry.entries[0].description).toBe("New description")
      expect(registry.entries[0].verifiedBy).toBe("community")
    })

    test("should remove entry from registry", async () => {
      await using tmp = await tmpdir()
      const registryPath = path.join(tmp.path, "registry.json")

      // Add entry
      await PluginRegistry.add(
        {
          name: "test-plugin",
          version: "1.0.0",
          hash: "a".repeat(64),
          verifiedBy: "user",
        },
        { localPath: registryPath, autoUpdate: false }
      )

      // Remove it
      const removed = await PluginRegistry.remove(
        "test-plugin",
        "1.0.0",
        { localPath: registryPath, autoUpdate: false }
      )

      expect(removed).toBe(true)

      const registry = await PluginRegistry.load({
        localPath: registryPath,
        autoUpdate: false,
      })

      expect(registry.entries.length).toBe(0)
    })

    test("should return false when removing non-existent entry", async () => {
      await using tmp = await tmpdir()
      const registryPath = path.join(tmp.path, "registry.json")

      const removed = await PluginRegistry.remove(
        "nonexistent",
        "1.0.0",
        { localPath: registryPath, autoUpdate: false }
      )

      expect(removed).toBe(false)
    })
  })

  describe("Registry Search", () => {
    test("should find entry by name and version", async () => {
      await using tmp = await tmpdir()
      const registryPath = path.join(tmp.path, "registry.json")

      await PluginRegistry.add(
        {
          name: "test-plugin",
          version: "1.0.0",
          hash: "a".repeat(64),
          verifiedBy: "user",
        },
        { localPath: registryPath, autoUpdate: false }
      )

      const entry = await PluginRegistry.find(
        "test-plugin",
        "1.0.0",
        { localPath: registryPath, autoUpdate: false }
      )

      expect(entry).not.toBeNull()
      expect(entry?.name).toBe("test-plugin")
      expect(entry?.version).toBe("1.0.0")
    })

    test("should return null for non-existent entry", async () => {
      await using tmp = await tmpdir()
      const registryPath = path.join(tmp.path, "registry.json")

      const entry = await PluginRegistry.find(
        "nonexistent",
        "1.0.0",
        { localPath: registryPath, autoUpdate: false }
      )

      expect(entry).toBeNull()
    })

    test("should find all versions of a plugin", async () => {
      await using tmp = await tmpdir()
      const registryPath = path.join(tmp.path, "registry.json")

      await PluginRegistry.add(
        {
          name: "test-plugin",
          version: "1.0.0",
          hash: "a".repeat(64),
          verifiedBy: "user",
        },
        { localPath: registryPath, autoUpdate: false }
      )

      await PluginRegistry.add(
        {
          name: "test-plugin",
          version: "2.0.0",
          hash: "b".repeat(64),
          verifiedBy: "user",
        },
        { localPath: registryPath, autoUpdate: false }
      )

      const entries = await PluginRegistry.findAll(
        "test-plugin",
        { localPath: registryPath, autoUpdate: false }
      )

      expect(entries.length).toBe(2)
      expect(entries[0].version).toBe("1.0.0")
      expect(entries[1].version).toBe("2.0.0")
    })

    test("should search by pattern", async () => {
      await using tmp = await tmpdir()
      const registryPath = path.join(tmp.path, "registry.json")

      await PluginRegistry.add(
        {
          name: "opencode-auth",
          version: "1.0.0",
          hash: "a".repeat(64),
          description: "Authentication plugin",
          verifiedBy: "official",
        },
        { localPath: registryPath, autoUpdate: false }
      )

      await PluginRegistry.add(
        {
          name: "opencode-formatter",
          version: "1.0.0",
          hash: "b".repeat(64),
          description: "Code formatter",
          verifiedBy: "community",
        },
        { localPath: registryPath, autoUpdate: false }
      )

      const results = await PluginRegistry.search(
        "opencode",
        { localPath: registryPath, autoUpdate: false }
      )

      expect(results.length).toBe(2)

      const authResults = await PluginRegistry.search(
        "auth",
        { localPath: registryPath, autoUpdate: false }
      )

      expect(authResults.length).toBe(1)
      expect(authResults[0].name).toBe("opencode-auth")
    })
  })

  describe("Hash Verification", () => {
    test("should verify correct hash", async () => {
      await using tmp = await tmpdir()
      const registryPath = path.join(tmp.path, "registry.json")
      const pluginPath = path.join(tmp.path, "plugin.js")
      await Bun.write(pluginPath, "console.log('verify')")

      const hash = await PluginSecurity.generateHash(pluginPath)

      await PluginRegistry.add(
        {
          name: "test-plugin",
          version: "1.0.0",
          hash,
          verifiedBy: "user",
        },
        { localPath: registryPath, autoUpdate: false }
      )

      const { verified, entry } = await PluginRegistry.verify(
        "test-plugin",
        "1.0.0",
        hash,
        { localPath: registryPath, autoUpdate: false }
      )

      expect(verified).toBe(true)
      expect(entry).not.toBeNull()
      expect(entry?.hash).toBe(hash)
    })

    test("should fail verification with wrong hash", async () => {
      await using tmp = await tmpdir()
      const registryPath = path.join(tmp.path, "registry.json")

      await PluginRegistry.add(
        {
          name: "test-plugin",
          version: "1.0.0",
          hash: "a".repeat(64),
          verifiedBy: "user",
        },
        { localPath: registryPath, autoUpdate: false }
      )

      const { verified, entry } = await PluginRegistry.verify(
        "test-plugin",
        "1.0.0",
        "b".repeat(64),
        { localPath: registryPath, autoUpdate: false }
      )

      expect(verified).toBe(false)
      expect(entry).not.toBeNull()
      expect(entry?.hash).toBe("a".repeat(64))
    })

    test("should fail verification for non-existent plugin", async () => {
      await using tmp = await tmpdir()
      const registryPath = path.join(tmp.path, "registry.json")

      const { verified, entry } = await PluginRegistry.verify(
        "nonexistent",
        "1.0.0",
        "a".repeat(64),
        { localPath: registryPath, autoUpdate: false }
      )

      expect(verified).toBe(false)
      expect(entry).toBeNull()
    })
  })

  describe("Registry Statistics", () => {
    test("should calculate stats correctly", async () => {
      await using tmp = await tmpdir()
      const registryPath = path.join(tmp.path, "registry.json")

      // Add official plugin
      await PluginRegistry.add(
        {
          name: "official-plugin",
          version: "1.0.0",
          hash: "a".repeat(64),
          verifiedBy: "official",
        },
        { localPath: registryPath, autoUpdate: false }
      )

      // Add community plugin
      await PluginRegistry.add(
        {
          name: "community-plugin",
          version: "1.0.0",
          hash: "b".repeat(64),
          verifiedBy: "community",
        },
        { localPath: registryPath, autoUpdate: false }
      )

      // Add user plugin (2 versions)
      await PluginRegistry.add(
        {
          name: "user-plugin",
          version: "1.0.0",
          hash: "c".repeat(64),
          verifiedBy: "user",
        },
        { localPath: registryPath, autoUpdate: false }
      )

      await PluginRegistry.add(
        {
          name: "user-plugin",
          version: "2.0.0",
          hash: "d".repeat(64),
          verifiedBy: "user",
        },
        { localPath: registryPath, autoUpdate: false }
      )

      const stats = await PluginRegistry.stats({
        localPath: registryPath,
        autoUpdate: false,
      })

      expect(stats.totalEntries).toBe(4)
      expect(stats.officialCount).toBe(1)
      expect(stats.communityCount).toBe(1)
      expect(stats.userCount).toBe(2)
      expect(stats.uniquePlugins).toBe(3)
    })
  })

  describe("Cache Management", () => {
    test("should use cached registry", async () => {
      await using tmp = await tmpdir()
      const registryPath = path.join(tmp.path, "registry.json")

      // First load
      const registry1 = await PluginRegistry.load({
        localPath: registryPath,
        autoUpdate: false,
        cacheTTL: 3600000, // 1 hour
      })

      // Second load (should be from cache)
      const registry2 = await PluginRegistry.load({
        localPath: registryPath,
        autoUpdate: false,
        cacheTTL: 3600000,
      })

      // Should be the same object reference (from cache)
      expect(registry2).toBe(registry1)
    })

    test("should clear cache", async () => {
      await using tmp = await tmpdir()
      const registryPath = path.join(tmp.path, "registry.json")

      await PluginRegistry.load({
        localPath: registryPath,
        autoUpdate: false,
      })

      PluginRegistry.clearCache()

      // After clearing cache, load should create new instance
      const registry = await PluginRegistry.load({
        localPath: registryPath,
        autoUpdate: false,
      })

      expect(registry).toBeDefined()
    })
  })

  describe("Integration with PluginSecurity", () => {
    test("should verify plugin with registry", async () => {
      await using tmp = await tmpdir()
      const registryPath = path.join(tmp.path, "registry.json")
      const pluginPath = path.join(tmp.path, "plugin.js")
      await Bun.write(pluginPath, "console.log('test')")

      const hash = await PluginSecurity.generateHash(pluginPath)

      await PluginRegistry.add(
        {
          name: "test-plugin",
          version: "1.0.0",
          hash,
          verifiedBy: "official",
        },
        { localPath: registryPath, autoUpdate: false }
      )

      const result = await PluginSecurity.verifyWithRegistry(
        "test-plugin",
        "1.0.0",
        pluginPath,
        { localPath: registryPath, autoUpdate: false }
      )

      expect(result.verified).toBe(true)
      expect(result.hashMatch).toBe(true)
      expect(result.entry).not.toBeNull()
      expect(result.entry?.verifiedBy).toBe("official")
    })

    test("should fail verification for tampered plugin", async () => {
      await using tmp = await tmpdir()
      const registryPath = path.join(tmp.path, "registry.json")
      const pluginPath = path.join(tmp.path, "plugin.js")
      await Bun.write(pluginPath, "console.log('original')")

      const hash = await PluginSecurity.generateHash(pluginPath)

      await PluginRegistry.add(
        {
          name: "test-plugin",
          version: "1.0.0",
          hash,
          verifiedBy: "official",
        },
        { localPath: registryPath, autoUpdate: false }
      )

      // Tamper with the plugin
      await Bun.write(pluginPath, "console.log('tampered')")

      const result = await PluginSecurity.verifyWithRegistry(
        "test-plugin",
        "1.0.0",
        pluginPath,
        { localPath: registryPath, autoUpdate: false }
      )

      expect(result.verified).toBe(false)
      expect(result.hashMatch).toBe(false)
    })
  })
})

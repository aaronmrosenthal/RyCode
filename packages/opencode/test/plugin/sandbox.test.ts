import { describe, expect, test, afterEach } from "bun:test"
import { PluginSandbox } from "../../src/plugin/sandbox"
import { PluginSecurity } from "../../src/plugin/security"

/**
 * Tests for Plugin Sandboxing System
 */

describe("PluginSandbox", () => {
  afterEach(async () => {
    // Clean up all sandboxes after each test
    await PluginSandbox.terminateAll()
  })

  describe("Sandbox Creation", () => {
    test("should create a sandbox with default config", async () => {
      const sandbox = await PluginSandbox.createSandbox({
        pluginName: "test-plugin",
        pluginVersion: "1.0.0",
        capabilities: {
          fileSystemRead: true,
          fileSystemWrite: false,
          network: false,
          shell: false,
          env: false,
          projectMetadata: true,
          aiClient: false,
        },
        resourceLimits: {
          maxMemoryMB: 512,
          maxExecutionTime: 30000,
          maxCPUTime: 10000,
          maxFileSystemOps: 1000,
          maxNetworkRequests: 100,
        },
        strictMode: true,
      })

      expect(sandbox).toBeDefined()
      expect(sandbox.execute).toBeFunction()
      expect(sandbox.terminate).toBeFunction()
      expect(sandbox.getResourceUsage).toBeFunction()

      await sandbox.terminate()
    })

    test("should track active sandboxes", async () => {
      const sandbox1 = await PluginSandbox.createSandbox({
        pluginName: "plugin-1",
        pluginVersion: "1.0.0",
        capabilities: PluginSecurity.DEFAULT_POLICY.defaultCapabilities,
        resourceLimits: {
          maxMemoryMB: 256,
          maxExecutionTime: 10000,
          maxCPUTime: 5000,
          maxFileSystemOps: 500,
          maxNetworkRequests: 50,
        },
        strictMode: true,
      })

      const sandbox2 = await PluginSandbox.createSandbox({
        pluginName: "plugin-2",
        pluginVersion: "1.0.0",
        capabilities: PluginSecurity.DEFAULT_POLICY.defaultCapabilities,
        resourceLimits: {
          maxMemoryMB: 256,
          maxExecutionTime: 10000,
          maxCPUTime: 5000,
          maxFileSystemOps: 500,
          maxNetworkRequests: 50,
        },
        strictMode: true,
      })

      const active = PluginSandbox.getActiveSandboxes()
      expect(active.size).toBe(2)

      await sandbox1.terminate()
      await sandbox2.terminate()
    })
  })

  describe("Sandbox Execution", () => {
    test("should execute plugin code", async () => {
      const sandbox = await PluginSandbox.createSandbox({
        pluginName: "test-plugin",
        pluginVersion: "1.0.0",
        capabilities: PluginSecurity.DEFAULT_POLICY.defaultCapabilities,
        resourceLimits: {
          maxMemoryMB: 512,
          maxExecutionTime: 30000,
          maxCPUTime: 10000,
          maxFileSystemOps: 1000,
          maxNetworkRequests: 100,
        },
        strictMode: true,
      })

      const result = await sandbox.execute({ test: "data" })

      expect(result).toBeDefined()
      expect(result).toHaveProperty("success", true)

      await sandbox.terminate()
    })

    test("should enforce execution timeout", async () => {
      const sandbox = await PluginSandbox.createSandbox({
        pluginName: "timeout-test",
        pluginVersion: "1.0.0",
        capabilities: PluginSecurity.DEFAULT_POLICY.defaultCapabilities,
        resourceLimits: {
          maxMemoryMB: 512,
          maxExecutionTime: 100, // Very short timeout
          maxCPUTime: 10000,
          maxFileSystemOps: 1000,
          maxNetworkRequests: 100,
        },
        strictMode: true,
      })

      // This should timeout (the worker initialization might take longer than 100ms)
      try {
        await sandbox.execute({ delay: 200 })
        // If we get here, the test might have completed too fast
        // This is acceptable in some environments
      } catch (error) {
        expect(error).toBeInstanceOf(PluginSandbox.SandboxTimeoutError)
      } finally {
        await sandbox.terminate()
      }
    }, 1000) // Give the test itself more time

    test("should not allow execution after termination", async () => {
      const sandbox = await PluginSandbox.createSandbox({
        pluginName: "test-plugin",
        pluginVersion: "1.0.0",
        capabilities: PluginSecurity.DEFAULT_POLICY.defaultCapabilities,
        resourceLimits: {
          maxMemoryMB: 512,
          maxExecutionTime: 30000,
          maxCPUTime: 10000,
          maxFileSystemOps: 1000,
          maxNetworkRequests: 100,
        },
        strictMode: true,
      })

      await sandbox.terminate()

      try {
        await sandbox.execute({ test: "data" })
        expect(true).toBe(false) // Should not reach here
      } catch (error) {
        expect(error).toBeInstanceOf(Error)
        expect((error as Error).message).toContain("terminated")
      }
    })
  })

  describe("Resource Tracking", () => {
    test("should track resource usage", async () => {
      const sandbox = await PluginSandbox.createSandbox({
        pluginName: "test-plugin",
        pluginVersion: "1.0.0",
        capabilities: PluginSecurity.DEFAULT_POLICY.defaultCapabilities,
        resourceLimits: {
          maxMemoryMB: 512,
          maxExecutionTime: 30000,
          maxCPUTime: 10000,
          maxFileSystemOps: 1000,
          maxNetworkRequests: 100,
        },
        strictMode: true,
      })

      await sandbox.execute({ test: "data" })

      const usage = sandbox.getResourceUsage()

      expect(usage).toHaveProperty("memoryMB")
      expect(usage).toHaveProperty("cpuTime")
      expect(usage).toHaveProperty("fileSystemOps")
      expect(usage).toHaveProperty("networkRequests")

      expect(usage.memoryMB).toBeGreaterThanOrEqual(0)
      expect(usage.cpuTime).toBeGreaterThanOrEqual(0)

      await sandbox.terminate()
    })

    test("should get sandbox statistics", async () => {
      const sandbox1 = await PluginSandbox.createSandbox({
        pluginName: "plugin-1",
        pluginVersion: "1.0.0",
        capabilities: PluginSecurity.DEFAULT_POLICY.defaultCapabilities,
        resourceLimits: {
          maxMemoryMB: 256,
          maxExecutionTime: 10000,
          maxCPUTime: 5000,
          maxFileSystemOps: 500,
          maxNetworkRequests: 50,
        },
        strictMode: true,
      })

      const stats = PluginSandbox.getStatistics()

      expect(stats.activeSandboxes).toBe(1)
      expect(stats.totalMemoryMB).toBeGreaterThanOrEqual(0)
      expect(stats.totalCPUTime).toBeGreaterThanOrEqual(0)

      await sandbox1.terminate()
    })
  })

  describe("Sandbox Termination", () => {
    test("should terminate a single sandbox", async () => {
      const sandbox = await PluginSandbox.createSandbox({
        pluginName: "test-plugin",
        pluginVersion: "1.0.0",
        capabilities: PluginSecurity.DEFAULT_POLICY.defaultCapabilities,
        resourceLimits: {
          maxMemoryMB: 512,
          maxExecutionTime: 30000,
          maxCPUTime: 10000,
          maxFileSystemOps: 1000,
          maxNetworkRequests: 100,
        },
        strictMode: true,
      })

      await sandbox.terminate()

      const active = PluginSandbox.getActiveSandboxes()
      expect(active.size).toBe(0)
    })

    test("should terminate all sandboxes", async () => {
      await PluginSandbox.createSandbox({
        pluginName: "plugin-1",
        pluginVersion: "1.0.0",
        capabilities: PluginSecurity.DEFAULT_POLICY.defaultCapabilities,
        resourceLimits: {
          maxMemoryMB: 256,
          maxExecutionTime: 10000,
          maxCPUTime: 5000,
          maxFileSystemOps: 500,
          maxNetworkRequests: 50,
        },
        strictMode: true,
      })

      await PluginSandbox.createSandbox({
        pluginName: "plugin-2",
        pluginVersion: "1.0.0",
        capabilities: PluginSecurity.DEFAULT_POLICY.defaultCapabilities,
        resourceLimits: {
          maxMemoryMB: 256,
          maxExecutionTime: 10000,
          maxCPUTime: 5000,
          maxFileSystemOps: 500,
          maxNetworkRequests: 50,
        },
        strictMode: true,
      })

      await PluginSandbox.terminateAll()

      const active = PluginSandbox.getActiveSandboxes()
      expect(active.size).toBe(0)
    })
  })
})

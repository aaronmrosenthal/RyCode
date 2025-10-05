import { describe, test, expect, beforeEach } from "bun:test"
import { PluginSecurity } from "../../src/plugin/security"

describe("PluginSecurity", () => {
  beforeEach(() => {
    PluginSecurity.clearAuditLog()
  })

  describe("isTrusted", () => {
    test("should trust official plugins", () => {
      const { trusted } = PluginSecurity.isTrusted(
        "opencode-copilot-auth",
        "0.0.3",
        PluginSecurity.DEFAULT_POLICY
      )
      expect(trusted).toBe(true)
    })

    test("should not trust unlisted plugins", () => {
      const { trusted } = PluginSecurity.isTrusted(
        "unknown-plugin",
        "1.0.0",
        PluginSecurity.DEFAULT_POLICY
      )
      expect(trusted).toBe(false)
    })

    test("should match exact version", () => {
      const policy: PluginSecurity.Policy = {
        ...PluginSecurity.DEFAULT_POLICY,
        trustedPlugins: [
          {
            name: "test-plugin",
            versions: ["1.2.3"],
            official: false,
            capabilities: PluginSecurity.DEFAULT_POLICY.defaultCapabilities,
          },
        ],
      }

      const { trusted } = PluginSecurity.isTrusted("test-plugin", "1.2.3", policy)
      expect(trusted).toBe(true)

      const { trusted: notTrusted } = PluginSecurity.isTrusted("test-plugin", "1.2.4", policy)
      expect(notTrusted).toBe(false)
    })

    test("should match latest version", () => {
      const policy: PluginSecurity.Policy = {
        ...PluginSecurity.DEFAULT_POLICY,
        trustedPlugins: [
          {
            name: "test-plugin",
            versions: ["latest"],
            official: false,
            capabilities: PluginSecurity.DEFAULT_POLICY.defaultCapabilities,
          },
        ],
      }

      const { trusted } = PluginSecurity.isTrusted("test-plugin", "99.99.99", policy)
      expect(trusted).toBe(true)
    })

    test("should match caret range", () => {
      const policy: PluginSecurity.Policy = {
        ...PluginSecurity.DEFAULT_POLICY,
        trustedPlugins: [
          {
            name: "test-plugin",
            versions: ["^1.2.0"],
            official: false,
            capabilities: PluginSecurity.DEFAULT_POLICY.defaultCapabilities,
          },
        ],
      }

      const { trusted: v123 } = PluginSecurity.isTrusted("test-plugin", "1.2.3", policy)
      expect(v123).toBe(true)

      const { trusted: v199 } = PluginSecurity.isTrusted("test-plugin", "1.9.9", policy)
      expect(v199).toBe(true)

      const { trusted: v200 } = PluginSecurity.isTrusted("test-plugin", "2.0.0", policy)
      expect(v200).toBe(false)
    })

    test("should match tilde range", () => {
      const policy: PluginSecurity.Policy = {
        ...PluginSecurity.DEFAULT_POLICY,
        trustedPlugins: [
          {
            name: "test-plugin",
            versions: ["~1.2.0"],
            official: false,
            capabilities: PluginSecurity.DEFAULT_POLICY.defaultCapabilities,
          },
        ],
      }

      const { trusted: v123 } = PluginSecurity.isTrusted("test-plugin", "1.2.3", policy)
      expect(v123).toBe(true)

      const { trusted: v130 } = PluginSecurity.isTrusted("test-plugin", "1.3.0", policy)
      expect(v130).toBe(false)
    })
  })

  describe("getCapabilities", () => {
    test("should return trusted plugin capabilities", () => {
      const capabilities = PluginSecurity.getCapabilities(
        "opencode-copilot-auth",
        "0.0.3",
        PluginSecurity.DEFAULT_POLICY
      )

      expect(capabilities.fileSystemRead).toBe(true)
      expect(capabilities.network).toBe(true)
      expect(capabilities.env).toBe(true)
    })

    test("should return default capabilities for untrusted plugins", () => {
      const capabilities = PluginSecurity.getCapabilities(
        "unknown-plugin",
        "1.0.0",
        PluginSecurity.DEFAULT_POLICY
      )

      expect(capabilities.fileSystemRead).toBe(true)
      expect(capabilities.fileSystemWrite).toBe(false)
      expect(capabilities.network).toBe(false)
      expect(capabilities.shell).toBe(false)
      expect(capabilities.env).toBe(false)
    })

    test("should respect custom default capabilities", () => {
      const policy: PluginSecurity.Policy = {
        ...PluginSecurity.DEFAULT_POLICY,
        defaultCapabilities: {
          fileSystemRead: false,
          fileSystemWrite: false,
          network: true,
          shell: false,
          env: false,
          projectMetadata: true,
          aiClient: true,
        },
      }

      const capabilities = PluginSecurity.getCapabilities("unknown-plugin", "1.0.0", policy)

      expect(capabilities.network).toBe(true)
      expect(capabilities.fileSystemRead).toBe(false)
    })
  })

  describe("checkCapability", () => {
    test("should allow permitted capability", () => {
      const capabilities: PluginSecurity.Capabilities = {
        fileSystemRead: true,
        fileSystemWrite: false,
        network: false,
        shell: false,
        env: false,
        projectMetadata: true,
        aiClient: false,
      }

      expect(() => {
        PluginSecurity.checkCapability("test-plugin", "fileSystemRead", capabilities)
      }).not.toThrow()
    })

    test("should deny forbidden capability", () => {
      const capabilities: PluginSecurity.Capabilities = {
        fileSystemRead: true,
        fileSystemWrite: false,
        network: false,
        shell: false,
        env: false,
        projectMetadata: true,
        aiClient: false,
      }

      expect(() => {
        PluginSecurity.checkCapability("test-plugin", "shell", capabilities)
      }).toThrow(PluginSecurity.CapabilityDeniedError)
    })
  })

  describe("createSandboxedInput", () => {
    test("should provide allowed resources", () => {
      const baseInput = {
        project: "test-project",
        worktree: "/test/worktree",
        directory: "/test/dir",
        client: {},
        $: {},
      }

      const capabilities: PluginSecurity.Capabilities = {
        fileSystemRead: true,
        fileSystemWrite: false,
        network: false,
        shell: false,
        env: false,
        projectMetadata: true,
        aiClient: true,
      }

      const sandboxed = PluginSecurity.createSandboxedInput("test-plugin", baseInput, capabilities)

      expect(sandboxed.project).toBe("test-project")
      expect(sandboxed.worktree).toBe("/test/worktree")
      expect(sandboxed.directory).toBe("/test/dir")
      expect(sandboxed.client).toBeDefined()
    })

    test("should block restricted resources", () => {
      const baseInput = {
        project: "test-project",
        worktree: "/test/worktree",
        directory: "/test/dir",
        client: {},
        $: {},
      }

      const capabilities: PluginSecurity.Capabilities = {
        fileSystemRead: false,
        fileSystemWrite: false,
        network: false,
        shell: false,
        env: false,
        projectMetadata: false,
        aiClient: false,
      }

      const sandboxed = PluginSecurity.createSandboxedInput("test-plugin", baseInput, capabilities)

      expect(sandboxed.project).toBeUndefined()
      expect(sandboxed.worktree).toBeUndefined()
      expect(sandboxed.directory).toBeUndefined()
      expect(sandboxed.client).toBeUndefined()
    })

    test("should throw when accessing forbidden shell", () => {
      const baseInput = {
        $: Bun.$,
      }

      const capabilities: PluginSecurity.Capabilities = {
        fileSystemRead: true,
        fileSystemWrite: false,
        network: false,
        shell: false,
        env: false,
        projectMetadata: true,
        aiClient: false,
      }

      const sandboxed = PluginSecurity.createSandboxedInput("test-plugin", baseInput, capabilities)

      expect(() => {
        // Try to access shell
        sandboxed.$.template
      }).toThrow(PluginSecurity.CapabilityDeniedError)
    })
  })

  describe("audit", () => {
    test("should log security events", () => {
      PluginSecurity.audit({
        plugin: "test-plugin",
        version: "1.0.0",
        action: "loaded",
        trusted: true,
      })

      const log = PluginSecurity.getAuditLog()
      expect(log.length).toBe(1)
      expect(log[0].plugin).toBe("test-plugin")
      expect(log[0].action).toBe("loaded")
    })

    test("should include timestamp in audit log", () => {
      const before = Date.now()
      PluginSecurity.audit({
        plugin: "test-plugin",
        version: "1.0.0",
        action: "loaded",
        trusted: true,
      })
      const after = Date.now()

      const log = PluginSecurity.getAuditLog()
      expect(log[0].timestamp).toBeGreaterThanOrEqual(before)
      expect(log[0].timestamp).toBeLessThanOrEqual(after)
    })

    test("should accumulate multiple events", () => {
      PluginSecurity.audit({
        plugin: "plugin-1",
        version: "1.0.0",
        action: "loaded",
        trusted: true,
      })

      PluginSecurity.audit({
        plugin: "plugin-2",
        version: "2.0.0",
        action: "denied",
        trusted: false,
        reason: "user_denied",
      })

      const log = PluginSecurity.getAuditLog()
      expect(log.length).toBe(2)
    })

    test("should clear audit log", () => {
      PluginSecurity.audit({
        plugin: "test-plugin",
        version: "1.0.0",
        action: "loaded",
        trusted: true,
      })

      expect(PluginSecurity.getAuditLog().length).toBe(1)

      PluginSecurity.clearAuditLog()

      expect(PluginSecurity.getAuditLog().length).toBe(0)
    })
  })

  describe("verifyIntegrity", () => {
    test("should pass when no hash provided", async () => {
      const result = await PluginSecurity.verifyIntegrity("/fake/path")
      expect(result).toBe(true)
    })

    // Note: Actual hash verification tests would require test files
  })

  describe("DEFAULT_POLICY", () => {
    test("should have official plugins trusted", () => {
      const officialPlugins = PluginSecurity.DEFAULT_POLICY.trustedPlugins.filter(
        (p) => p.official
      )
      expect(officialPlugins.length).toBeGreaterThan(0)
    })

    test("should default to warn mode", () => {
      expect(PluginSecurity.DEFAULT_POLICY.mode).toBe("warn")
    })

    test("should require approval by default", () => {
      expect(PluginSecurity.DEFAULT_POLICY.requireApproval).toBe(true)
    })

    test("should have restrictive default capabilities", () => {
      const caps = PluginSecurity.DEFAULT_POLICY.defaultCapabilities
      expect(caps.fileSystemWrite).toBe(false)
      expect(caps.network).toBe(false)
      expect(caps.shell).toBe(false)
      expect(caps.env).toBe(false)
    })
  })
})

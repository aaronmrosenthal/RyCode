import { describe, expect, test, beforeEach } from "bun:test"
import { $ } from "bun"
import path from "path"
import { tmpdir } from "../fixture/fixture"
import { PluginSecurity } from "../../src/plugin/security"

/**
 * Integration tests for plugin CLI commands
 */

const cliPath = path.join(__dirname, "../../src/index.ts")

// Helper to run CLI and capture output
async function runCLI(command: string) {
  // Parse command properly handling quoted arguments
  const argArray: string[] = []
  let current = ""
  let inQuote = false

  for (let i = 0; i < command.length; i++) {
    const char = command[i]
    if (char === '"' || char === "'") {
      inQuote = !inQuote
    } else if (char === " " && !inQuote) {
      if (current) {
        argArray.push(current)
        current = ""
      }
    } else {
      current += char
    }
  }
  if (current) argArray.push(current)

  const proc = await $`bun ${cliPath} ${argArray}`.quiet().nothrow()
  const output = proc.stdout.toString() + proc.stderr.toString()
  return { output, exitCode: proc.exitCode, proc }
}

describe("plugin CLI commands", () => {
  beforeEach(() => {
    PluginSecurity.clearAuditLog()
  })

  describe("plugin:hash", () => {
    test("should generate hash for a file", async () => {
      await using tmp = await tmpdir()
      const testFile = path.join(tmp.path, "test-plugin.js")
      await Bun.write(testFile, "console.log('test')")

      const { output } = await runCLI(`plugin:hash ${testFile}`)

      expect(output).toContain("Plugin Hash Generated")
      expect(output).toContain("File:")
      expect(output).toContain("Hash:")
      expect(output).toMatch(/[0-9a-f]{64}/)
    })

    test("should output JSON format", async () => {
      await using tmp = await tmpdir()
      const testFile = path.join(tmp.path, "test-plugin.js")
      await Bun.write(testFile, "console.log('test')")

      const { output } = await runCLI(`plugin:hash ${testFile} --json`)

      const json = JSON.parse(output)
      expect(json).toHaveProperty("path")
      expect(json).toHaveProperty("hash")
      expect(json.hash).toMatch(/^[0-9a-f]{64}$/)
    })

    test("should fail for non-existent file", async () => {
      const { output, exitCode } = await runCLI(`plugin:hash /nonexistent/file.js`)

      expect(exitCode).toBe(1)
      expect(output).toContain("Plugin file not found")
    })

    test("should generate consistent hashes", async () => {
      await using tmp = await tmpdir()
      const testFile = path.join(tmp.path, "test-plugin.js")
      await Bun.write(testFile, "console.log('consistent')")

      const { output: output1 } = await runCLI(`plugin:hash ${testFile} --json`)
      const { output: output2 } = await runCLI(`plugin:hash ${testFile} --json`)

      const hash1 = JSON.parse(output1).hash
      const hash2 = JSON.parse(output2).hash

      expect(hash1).toBe(hash2)
    })
  })

  describe("plugin:check", () => {
    test("should check trusted official plugin", async () => {
      const { output } = await runCLI(`plugin:check opencode-copilot-auth 0.0.3`)

      expect(output).toContain("Plugin Trust Status")
      expect(output).toContain("TRUSTED")
      expect(output).toContain("opencode-copilot-auth")
      expect(output).toContain("0.0.3")
      expect(output).toContain("Official:")
      expect(output).toContain("Yes")
    })

    test("should check untrusted plugin", async () => {
      const { output } = await runCLI(`plugin:check unknown-plugin 1.0.0`)

      expect(output).toContain("UNTRUSTED")
      expect(output).toContain("unknown-plugin")
      expect(output).toContain("To trust this plugin")
    })

    test("should output JSON format", async () => {
      const { output } = await runCLI(`plugin:check opencode-copilot-auth 0.0.3 --json`)

      const json = JSON.parse(output)
      expect(json).toHaveProperty("plugin", "opencode-copilot-auth")
      expect(json).toHaveProperty("version", "0.0.3")
      expect(json).toHaveProperty("trusted", true)
      expect(json).toHaveProperty("official", true)
      expect(json.capabilities).toHaveProperty("fileSystemRead")
    })
  })

  describe("plugin:verify", () => {
    test("should verify plugin with correct hash", async () => {
      await using tmp = await tmpdir()
      const testFile = path.join(tmp.path, "test-plugin.js")
      await Bun.write(testFile, "console.log('verify me')")

      const hash = await PluginSecurity.generateHash(testFile)
      const { output } = await runCLI(`plugin:verify ${testFile} --hash ${hash}`)

      expect(output).toContain("Integrity check PASSED")
      expect(output).toContain("Plugin has not been tampered with")
    })

    test("should fail verification with wrong hash", async () => {
      await using tmp = await tmpdir()
      const testFile = path.join(tmp.path, "test-plugin.js")
      await Bun.write(testFile, "console.log('verify me')")

      const wrongHash = "0".repeat(64)
      const { output, exitCode } = await runCLI(`plugin:verify ${testFile} --hash ${wrongHash}`)

      expect(exitCode).toBe(1)
      expect(output).toContain("Integrity check FAILED")
      expect(output).toContain("Plugin may have been tampered with")
    })

    test("should output JSON format", async () => {
      await using tmp = await tmpdir()
      const testFile = path.join(tmp.path, "test-plugin.js")
      await Bun.write(testFile, "console.log('json')")

      const hash = await PluginSecurity.generateHash(testFile)
      const { output } = await runCLI(`plugin:verify ${testFile} --hash ${hash} --json`)

      const json = JSON.parse(output)
      expect(json).toHaveProperty("valid", true)
    })
  })

  describe("plugin:audit", () => {
    test("should show empty audit log", async () => {
      const { output } = await runCLI(`plugin:audit`)

      expect(output).toContain("No audit log entries found")
    })

    test("should output JSON format for empty log", async () => {
      const { output } = await runCLI(`plugin:audit --json`)

      const json = JSON.parse(output)
      expect(Array.isArray(json)).toBe(true)
    })
  })

  describe("edge cases", () => {
    test("should handle special characters in paths", async () => {
      await using tmp = await tmpdir()
      const testFile = path.join(tmp.path, "test plugin (1).js")
      await Bun.write(testFile, "console.log('special')")

      const { output } = await runCLI(`plugin:hash "${testFile}"`)

      expect(output).toContain("Plugin Hash Generated")
    })

    test("should handle binary files", async () => {
      await using tmp = await tmpdir()
      const testFile = path.join(tmp.path, "binary.bin")
      await Bun.write(testFile, new Uint8Array([0x00, 0x01, 0x02, 0xff]))

      const { output } = await runCLI(`plugin:hash ${testFile}`)

      expect(output).toContain("Plugin Hash Generated")
      expect(output).toMatch(/[0-9a-f]{64}/)
    })
  })

  describe("JSON output consistency", () => {
    test("all commands should support --json", async () => {
      await using tmp = await tmpdir()
      const testFile = path.join(tmp.path, "test.js")
      await Bun.write(testFile, "test")
      const hash = await PluginSecurity.generateHash(testFile)

      const { output: hashJson } = await runCLI(`plugin:hash ${testFile} --json`)
      const { output: checkJson } = await runCLI(`plugin:check test-plugin 1.0.0 --json`)
      const { output: verifyJson } = await runCLI(`plugin:verify ${testFile} --hash ${hash} --json`)
      const { output: auditJson } = await runCLI(`plugin:audit --json`)

      expect(() => JSON.parse(hashJson)).not.toThrow()
      expect(() => JSON.parse(checkJson)).not.toThrow()
      expect(() => JSON.parse(verifyJson)).not.toThrow()
      expect(() => JSON.parse(auditJson)).not.toThrow()
    })
  })
})

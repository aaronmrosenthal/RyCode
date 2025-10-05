import { cmd } from "./cmd"
import { PluginSecurity } from "../../plugin/security"
import { UI } from "../ui"
import path from "path"
import { existsSync } from "fs"

/**
 * Plugin Security CLI Commands
 *
 * Provides commands to manage plugin security:
 * - plugin:hash    Generate SHA-256 hash for a plugin
 * - plugin:audit   View security audit log
 * - plugin:check   Check if a plugin is trusted
 * - plugin:verify  Verify plugin integrity
 */

export const PluginHashCommand = cmd({
  command: "plugin:hash <plugin-path>",
  describe: "Generate SHA-256 hash for a plugin file",
  builder: (yargs) =>
    yargs
      .positional("plugin-path", {
        describe: "Path to the plugin file",
        type: "string",
        demandOption: true,
      })
      .option("json", {
        describe: "Output as JSON",
        type: "boolean",
        default: false,
      }),
  handler: async (args) => {
    const pluginPath = path.resolve(args.pluginPath as string)

    if (!existsSync(pluginPath)) {
      UI.error(`Plugin file not found: ${pluginPath}`)
      process.exit(1)
    }

    try {
      const hash = await PluginSecurity.generateHash(pluginPath)

      if (args.json) {
        console.log(JSON.stringify({ path: pluginPath, hash }, null, 2))
      } else {
        UI.println()
        UI.println(UI.Style.BOLD + "Plugin Hash Generated" + UI.Style.RESET)
        UI.println()
        UI.println(UI.Style.TEXT_INFO + "File:  " + UI.Style.RESET + pluginPath)
        UI.println(UI.Style.TEXT_INFO + "Hash:  " + UI.Style.RESET + UI.Style.TEXT_SUCCESS + hash + UI.Style.RESET)
        UI.println()
        UI.println(UI.Style.DIM + "Add to your .rycode.json:" + UI.Style.RESET)
        UI.println()
        UI.println(UI.Style.DIM + JSON.stringify({
          plugin_security: {
            trustedPlugins: [{
              name: path.basename(pluginPath, path.extname(pluginPath)),
              hash: hash,
            }]
          }
        }, null, 2) + UI.Style.RESET)
        UI.println()
      }
    } catch (error) {
      UI.error(`Failed to generate hash: ${error instanceof Error ? error.message : String(error)}`)
      process.exit(1)
    }
  },
})

export const PluginAuditCommand = cmd({
  command: "plugin:audit",
  describe: "View security audit log",
  builder: (yargs) =>
    yargs
      .option("json", {
        describe: "Output as JSON",
        type: "boolean",
        default: false,
      })
      .option("filter", {
        describe: "Filter by action (loaded, denied, capability_check)",
        type: "string",
        choices: ["loaded", "denied", "capability_check"],
      })
      .option("limit", {
        describe: "Limit number of entries",
        type: "number",
      }),
  handler: async (args) => {
    const auditLog = PluginSecurity.getAuditLog()

    let filtered = Array.from(auditLog)

    if (args.filter) {
      filtered = filtered.filter(e => e.action === args.filter)
    }

    if (args.limit) {
      filtered = filtered.slice(-args.limit)
    }

    if (args.json) {
      console.log(JSON.stringify(filtered, null, 2))
    } else {
      if (filtered.length === 0) {
        UI.println()
        UI.println(UI.Style.TEXT_INFO + "No audit log entries found." + UI.Style.RESET)
        UI.println()
        return
      }

      UI.println()
      UI.println(UI.Style.BOLD + "Plugin Security Audit Log" + UI.Style.RESET)
      UI.println(UI.Style.DIM + `${filtered.length} ${filtered.length === 1 ? 'entry' : 'entries'}` + UI.Style.RESET)
      UI.println()

      for (const entry of filtered) {
        const timestamp = new Date(entry.timestamp).toISOString()
        const actionColor = entry.action === "denied" ? UI.Style.TEXT_DANGER :
                           entry.action === "loaded" ? UI.Style.TEXT_SUCCESS :
                           UI.Style.TEXT_WARNING

        const trustIcon = entry.trusted ? "✓" : "✗"
        const trustColor = entry.trusted ? UI.Style.TEXT_SUCCESS : UI.Style.TEXT_DANGER

        UI.println(UI.Style.DIM + timestamp + UI.Style.RESET)
        UI.println(`  ${trustColor}${trustIcon}${UI.Style.RESET} ${UI.Style.BOLD}${entry.plugin}${UI.Style.RESET}@${entry.version}`)
        UI.println(`  ${actionColor}${entry.action}${UI.Style.RESET}`)

        if (entry.reason) {
          UI.println(`  ${UI.Style.DIM}Reason: ${entry.reason}${UI.Style.RESET}`)
        }

        if (entry.capabilities) {
          const caps = Object.entries(entry.capabilities)
            .filter(([_, v]) => v)
            .map(([k, _]) => k)
            .join(", ")
          if (caps) {
            UI.println(`  ${UI.Style.DIM}Capabilities: ${caps}${UI.Style.RESET}`)
          }
        }

        UI.println()
      }
    }
  },
})

export const PluginCheckCommand = cmd({
  command: "plugin:check <plugin-name> <plugin-version>",
  describe: "Check if a plugin is trusted",
  builder: (yargs) =>
    yargs
      .positional("plugin-name", {
        describe: "Plugin package name",
        type: "string",
        demandOption: true,
      })
      .positional("plugin-version", {
        describe: "Plugin version",
        type: "string",
        demandOption: true,
      })
      .option("json", {
        describe: "Output as JSON",
        type: "boolean",
        default: false,
      }),
  handler: async (args) => {
    const pluginName = args.pluginName as string
    const version = args.pluginVersion as string

    const { trusted, config } = PluginSecurity.isTrusted(
      pluginName,
      version,
      PluginSecurity.DEFAULT_POLICY
    )

    const capabilities = PluginSecurity.getCapabilities(
      pluginName,
      version,
      PluginSecurity.DEFAULT_POLICY
    )

    if (args.json) {
      console.log(JSON.stringify({
        plugin: pluginName,
        version,
        trusted,
        official: config?.official || false,
        capabilities,
      }, null, 2))
    } else {
      UI.println()
      UI.println(UI.Style.BOLD + "Plugin Trust Status" + UI.Style.RESET)
      UI.println()

      const statusIcon = trusted ? "✓" : "✗"
      const statusColor = trusted ? UI.Style.TEXT_SUCCESS : UI.Style.TEXT_DANGER
      const statusText = trusted ? "TRUSTED" : "UNTRUSTED"

      UI.println(`  ${statusColor}${statusIcon} ${statusText}${UI.Style.RESET}`)
      UI.println()
      UI.println(`  ${UI.Style.TEXT_INFO}Plugin:${UI.Style.RESET}  ${pluginName}`)
      UI.println(`  ${UI.Style.TEXT_INFO}Version:${UI.Style.RESET} ${version}`)

      if (config?.official) {
        UI.println(`  ${UI.Style.TEXT_SUCCESS}Official:${UI.Style.RESET} Yes`)
      }

      UI.println()
      UI.println(UI.Style.BOLD + "Capabilities:" + UI.Style.RESET)
      UI.println()

      for (const [cap, enabled] of Object.entries(capabilities)) {
        const icon = enabled ? "✓" : "✗"
        const color = enabled ? UI.Style.TEXT_SUCCESS : UI.Style.DIM
        UI.println(`  ${color}${icon} ${cap}${UI.Style.RESET}`)
      }

      UI.println()

      if (!trusted) {
        UI.println(UI.Style.TEXT_WARNING + "💡 To trust this plugin, add it to .rycode.json:" + UI.Style.RESET)
        UI.println()
        UI.println(UI.Style.DIM + JSON.stringify({
          plugin_security: {
            trustedPlugins: [{
              name: pluginName,
              versions: [version],
              capabilities: capabilities,
            }]
          }
        }, null, 2) + UI.Style.RESET)
        UI.println()
      }
    }
  },
})

export const PluginVerifyCommand = cmd({
  command: "plugin:verify <plugin-path>",
  describe: "Verify plugin integrity using SHA-256 hash",
  builder: (yargs) =>
    yargs
      .positional("plugin-path", {
        describe: "Path to the plugin file",
        type: "string",
        demandOption: true,
      })
      .option("hash", {
        describe: "Expected SHA-256 hash",
        type: "string",
        demandOption: true,
      })
      .option("json", {
        describe: "Output as JSON",
        type: "boolean",
        default: false,
      }),
  handler: async (args) => {
    const pluginPath = path.resolve(args.pluginPath as string)
    const expectedHash = args.hash as string

    if (!existsSync(pluginPath)) {
      UI.error(`Plugin file not found: ${pluginPath}`)
      process.exit(1)
    }

    try {
      const actualHash = await PluginSecurity.generateHash(pluginPath)
      const isValid = await PluginSecurity.verifyIntegrity(pluginPath, expectedHash)

      if (args.json) {
        console.log(JSON.stringify({
          path: pluginPath,
          expected: expectedHash,
          actual: actualHash,
          valid: isValid,
        }, null, 2))
      } else {
        UI.println()
        UI.println(UI.Style.BOLD + "Plugin Integrity Verification" + UI.Style.RESET)
        UI.println()

        UI.println(UI.Style.TEXT_INFO + "File:     " + UI.Style.RESET + pluginPath)
        UI.println(UI.Style.TEXT_INFO + "Expected: " + UI.Style.RESET + expectedHash)
        UI.println(UI.Style.TEXT_INFO + "Actual:   " + UI.Style.RESET + actualHash)
        UI.println()

        if (isValid) {
          UI.println(UI.Style.TEXT_SUCCESS + "✓ Integrity check PASSED" + UI.Style.RESET)
          UI.println(UI.Style.TEXT_SUCCESS + "  Plugin has not been tampered with." + UI.Style.RESET)
        } else {
          UI.println(UI.Style.TEXT_DANGER + "✗ Integrity check FAILED" + UI.Style.RESET)
          UI.println(UI.Style.TEXT_DANGER + "  Plugin may have been tampered with!" + UI.Style.RESET)
          UI.println()
          UI.println(UI.Style.TEXT_WARNING + "⚠ DO NOT use this plugin." + UI.Style.RESET)
          process.exitCode = 1
        }

        UI.println()
      }
    } catch (error) {
      UI.error(`Failed to verify plugin: ${error instanceof Error ? error.message : String(error)}`)
      process.exit(1)
    }
  },
})

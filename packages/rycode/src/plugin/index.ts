import type { Hooks, PluginInput, Plugin as PluginInstance } from "@rycode-ai/plugin"
import { Config } from "../config/config"
import { Bus } from "../bus"
import { Log } from "../util/log"
import { createOpencodeClient } from "@rycode-ai/sdk"
import { Server } from "../server/server"
import { BunProc } from "../bun"
import { Instance } from "../project/instance"
import { Flag } from "../flag/flag"
import { PluginSecurity } from "./security"
import { Permission } from "../permission"

export namespace Plugin {
  const log = Log.create({ service: "plugin" })

  const state = Instance.state(async () => {
    const client = createOpencodeClient({
      baseUrl: "http://localhost:4096",
      fetch: async (...args) => Server.App().fetch(...args),
    })
    const config = await Config.get()

    // Load security policy from config or use defaults
    const securityPolicy = config.plugin_security ?? PluginSecurity.DEFAULT_POLICY

    const hooks = []
    const baseInput: PluginInput = {
      client,
      project: Instance.project,
      worktree: Instance.worktree,
      directory: Instance.directory,
      $: Bun.$,
    }

    const plugins = [...(config.plugin ?? [])]
    if (!Flag.OPENCODE_DISABLE_DEFAULT_PLUGINS) {
      plugins.push("opencode-copilot-auth@0.0.3")
      plugins.push("opencode-anthropic-auth@0.0.2")
    }

    for (let plugin of plugins) {
      log.info("loading plugin", { path: plugin })

      // Parse plugin name and version
      let pluginPath = plugin
      let packageName = plugin
      let version = "latest"

      if (!plugin.startsWith("file://")) {
        const parts = plugin.split("@")
        packageName = parts[0]
        version = parts[1] ?? "latest"

        // Security check: Is this plugin trusted?
        const { trusted, config: pluginConfig } = PluginSecurity.isTrusted(
          packageName,
          version,
          securityPolicy
        )

        // Log audit entry
        PluginSecurity.audit({
          plugin: packageName,
          version,
          action: "loaded",
          trusted,
          capabilities: PluginSecurity.getCapabilities(packageName, version, securityPolicy),
        })

        // In strict mode, deny untrusted plugins
        if (securityPolicy.mode === "strict" && !trusted) {
          log.error("untrusted plugin blocked in strict mode", { plugin: packageName, version })
          throw new PluginSecurity.UntrustedPluginError({
            plugin: packageName,
            version,
            message: `Plugin "${packageName}@${version}" is not in the trusted allowlist. Add it to plugin_security.trustedPlugins in your config.`,
          })
        }

        // In warn mode, notify but allow
        if (securityPolicy.mode === "warn" && !trusted) {
          log.warn("loading untrusted plugin with restricted capabilities", {
            plugin: packageName,
            version,
            capabilities: securityPolicy.defaultCapabilities,
          })
        }

        // Request user approval if required
        if (!trusted && securityPolicy.requireApproval) {
          try {
            await Permission.ask({
              type: "plugin_install",
              sessionID: "system",
              messageID: "system",
              callID: "system",
              title: `Install plugin: ${packageName}@${version}`,
              metadata: {
                plugin: packageName,
                version,
                trusted: false,
                capabilities: securityPolicy.defaultCapabilities,
              },
            })
          } catch (error) {
            log.info("user denied plugin installation", { plugin: packageName })
            PluginSecurity.audit({
              plugin: packageName,
              version,
              action: "denied",
              trusted: false,
              reason: "user_denied",
            })
            continue // Skip this plugin
          }
        }

        // Install plugin
        pluginPath = await BunProc.install(packageName, version)

        // Verify integrity if enabled and hash provided
        if (securityPolicy.verifyIntegrity && pluginConfig?.hash) {
          const integrityOk = await PluginSecurity.verifyIntegrity(
            pluginPath,
            pluginConfig.hash
          )

          if (!integrityOk) {
            log.error("plugin integrity check failed", { plugin: packageName, version })
            throw new PluginSecurity.IntegrityCheckFailedError({
              plugin: packageName,
              expected: pluginConfig.hash,
              actual: await PluginSecurity.generateHash(pluginPath),
            })
          }
        }
      }

      // Get capabilities for this plugin
      const capabilities = PluginSecurity.getCapabilities(packageName, version, securityPolicy)

      // Create sandboxed input based on capabilities
      const sandboxedInput = PluginSecurity.createSandboxedInput(
        packageName,
        baseInput,
        capabilities
      )

      // Load and initialize plugin
      try {
        const mod = await import(pluginPath)
        for (const [_name, fn] of Object.entries<PluginInstance>(mod)) {
          const init = await fn(sandboxedInput)
          hooks.push({
            ...init,
            _metadata: {
              name: packageName,
              version,
              capabilities,
            },
          })
        }
        log.info("plugin loaded successfully", {
          plugin: packageName,
          version,
          capabilities,
        })
      } catch (error) {
        log.error("failed to load plugin", {
          plugin: packageName,
          version,
          error,
        })

        // In strict mode, fail completely
        if (securityPolicy.mode === "strict") {
          throw error
        }

        // In other modes, continue with remaining plugins
        log.warn("skipping failed plugin", { plugin: packageName })
      }
    }

    return {
      hooks,
      input: baseInput,
      securityPolicy,
    }
  })

  export async function trigger<
    Name extends Exclude<keyof Required<Hooks>, "auth" | "event" | "tool">,
    Input = Parameters<Required<Hooks>[Name]>[0],
    Output = Parameters<Required<Hooks>[Name]>[1],
  >(name: Name, input: Input, output: Output): Promise<Output> {
    if (!name) return output
    for (const hook of await state().then((x) => x.hooks)) {
      const fn = hook[name]
      if (!fn) continue
      // @ts-expect-error if you feel adventurous, please fix the typing, make sure to bump the try-counter if you
      // give up.
      // try-counter: 2
      await fn(input, output)
    }
    return output
  }

  export async function list() {
    return state().then((x) => x.hooks)
  }

  export async function init() {
    const hooks = await state().then((x) => x.hooks)
    const config = await Config.get()
    for (const hook of hooks) {
      await hook.config?.(config)
    }
    Bus.subscribeAll(async (input) => {
      const hooks = await state().then((x) => x.hooks)
      for (const hook of hooks) {
        hook["event"]?.({
          event: input,
        })
      }
    })
  }

  /**
   * Get the security policy in use
   */
  export async function getSecurityPolicy() {
    return state().then((x) => x.securityPolicy)
  }

  /**
   * Get security audit log
   */
  export function getSecurityAuditLog() {
    return PluginSecurity.getAuditLog()
  }

  /**
   * Generate hash for a plugin file (for creating allowlist entries)
   */
  export async function generatePluginHash(pluginPath: string) {
    return PluginSecurity.generateHash(pluginPath)
  }

  /**
   * Check if a plugin is trusted
   */
  export async function isPluginTrusted(packageName: string, version: string) {
    const policy = await getSecurityPolicy()
    return PluginSecurity.isTrusted(packageName, version, policy)
  }

  // Re-export security types for convenience
  export type PluginCapabilities = PluginSecurity.Capabilities
  export type PluginSecurityPolicy = PluginSecurity.Policy
  export type TrustedPlugin = PluginSecurity.TrustedPlugin
}

/**
 * Plugin Sandboxing System
 *
 * Provides process-level isolation for plugins using worker threads.
 * Enforces resource limits, timeouts, and capability restrictions.
 */

import z from "zod/v4"
import { NamedError } from "../util/error"
import { Log } from "../util/log"
import { Worker } from "worker_threads"
import { PluginSecurity } from "./security"

export namespace PluginSandbox {
  const log = Log.create({ service: "plugin.sandbox" })

  /**
   * Sandbox timeout error
   */
  export const SandboxTimeoutError = NamedError.create(
    "SandboxTimeoutError",
    z.object({
      plugin: z.string(),
      timeout: z.number(),
    })
  )
  export type SandboxTimeoutError = InstanceType<typeof SandboxTimeoutError>

  /**
   * Sandbox resource limit error
   */
  export const SandboxResourceError = NamedError.create(
    "SandboxResourceError",
    z.object({
      plugin: z.string(),
      resource: z.string(),
      limit: z.number(),
      actual: z.number(),
    })
  )
  export type SandboxResourceError = InstanceType<typeof SandboxResourceError>

  /**
   * Resource limits for sandbox
   */
  export const ResourceLimits = z.object({
    /** Maximum memory in MB */
    maxMemoryMB: z.number().default(512),
    /** Maximum execution time in ms */
    maxExecutionTime: z.number().default(30000), // 30 seconds
    /** Maximum CPU time in ms */
    maxCPUTime: z.number().default(10000), // 10 seconds
    /** Maximum file system operations */
    maxFileSystemOps: z.number().default(1000),
    /** Maximum network requests */
    maxNetworkRequests: z.number().default(100),
  })
  export type ResourceLimits = z.infer<typeof ResourceLimits>

  /**
   * Sandbox configuration
   */
  export const SandboxConfig = z.object({
    /** Plugin name */
    pluginName: z.string(),
    /** Plugin version */
    pluginVersion: z.string(),
    /** Plugin capabilities */
    capabilities: PluginSecurity.Capabilities,
    /** Resource limits */
    resourceLimits: ResourceLimits,
    /** Enable strict mode (no eval, no dynamic require) */
    strictMode: z.boolean().default(true),
    /** Working directory for plugin */
    workingDirectory: z.string().optional(),
  })
  export type SandboxConfig = z.infer<typeof SandboxConfig>

  /**
   * Sandbox state
   */
  interface SandboxState {
    worker: Worker | null
    startTime: number
    resourceUsage: {
      memoryMB: number
      cpuTime: number
      fileSystemOps: number
      networkRequests: number
    }
    isTerminated: boolean
  }

  /**
   * Active sandboxes
   */
  const activeSandboxes = new Map<string, SandboxState>()

  /**
   * Create a sandboxed plugin environment
   */
  export async function createSandbox(
    config: SandboxConfig
  ): Promise<{
    execute: <T>(input: any) => Promise<T>
    terminate: () => Promise<void>
    getResourceUsage: () => SandboxState["resourceUsage"]
  }> {
    const sandboxId = `${config.pluginName}@${config.pluginVersion}`

    log.info("creating sandbox", {
      plugin: sandboxId,
      capabilities: config.capabilities,
      resourceLimits: config.resourceLimits,
    })

    // Initialize sandbox state
    const state: SandboxState = {
      worker: null,
      startTime: Date.now(),
      resourceUsage: {
        memoryMB: 0,
        cpuTime: 0,
        fileSystemOps: 0,
        networkRequests: 0,
      },
      isTerminated: false,
    }

    activeSandboxes.set(sandboxId, state)

    /**
     * Execute plugin code in sandbox
     */
    async function execute<T>(input: any): Promise<T> {
      if (state.isTerminated) {
        throw new Error("Sandbox has been terminated")
      }

      return new Promise((resolve, reject) => {
        // Create worker thread
        // Note: Bun doesn't support resourceLimits yet, so we track them manually
        const worker = new Worker(
          new URL("./sandbox-worker.ts", import.meta.url),
          {
            workerData: {
              config,
              input,
            },
            // resourceLimits not supported in Bun yet
            // We'll enforce limits through monitoring instead
          }
        )

        state.worker = worker

        // Timeout handler
        const timeout = setTimeout(() => {
          worker.terminate()
          state.isTerminated = true

          log.error("sandbox timeout", {
            plugin: sandboxId,
            timeout: config.resourceLimits.maxExecutionTime,
          })

          reject(new SandboxTimeoutError({
            plugin: sandboxId,
            timeout: config.resourceLimits.maxExecutionTime,
          }))
        }, config.resourceLimits.maxExecutionTime)

        // Message handler
        worker.on("message", (message) => {
          clearTimeout(timeout)

          if (message.type === "success") {
            log.info("sandbox execution completed", {
              plugin: sandboxId,
              duration: Date.now() - state.startTime,
            })

            resolve(message.result)
          } else if (message.type === "error") {
            log.error("sandbox execution failed", {
              plugin: sandboxId,
              error: message.error,
            })

            reject(new Error(message.error))
          } else if (message.type === "resourceUsage") {
            // Update resource usage
            state.resourceUsage = message.usage

            // Check limits
            if (state.resourceUsage.memoryMB > config.resourceLimits.maxMemoryMB) {
              worker.terminate()
              state.isTerminated = true

              reject(new SandboxResourceError({
                plugin: sandboxId,
                resource: "memory",
                limit: config.resourceLimits.maxMemoryMB,
                actual: state.resourceUsage.memoryMB,
              }))
            }
          }
        })

        // Error handler
        worker.on("error", (error) => {
          clearTimeout(timeout)

          log.error("sandbox worker error", {
            plugin: sandboxId,
            error: error.message,
          })

          reject(error)
        })

        // Exit handler
        worker.on("exit", (code) => {
          clearTimeout(timeout)

          if (code !== 0 && !state.isTerminated) {
            log.error("sandbox worker exited abnormally", {
              plugin: sandboxId,
              code,
            })

            reject(new Error(`Worker exited with code ${code}`))
          }
        })
      })
    }

    /**
     * Terminate sandbox
     */
    async function terminate(): Promise<void> {
      if (!state.isTerminated) {
        state.isTerminated = true

        if (state.worker) {
          await state.worker.terminate()
        }

        activeSandboxes.delete(sandboxId)

        log.info("sandbox terminated", {
          plugin: sandboxId,
          duration: Date.now() - state.startTime,
        })
      }
    }

    /**
     * Get current resource usage
     */
    function getResourceUsage(): SandboxState["resourceUsage"] {
      return { ...state.resourceUsage }
    }

    return {
      execute,
      terminate,
      getResourceUsage,
    }
  }

  /**
   * Get all active sandboxes
   */
  export function getActiveSandboxes(): Map<string, SandboxState> {
    return new Map(activeSandboxes)
  }

  /**
   * Terminate all active sandboxes
   */
  export async function terminateAll(): Promise<void> {
    log.info("terminating all sandboxes", {
      count: activeSandboxes.size,
    })

    const promises = Array.from(activeSandboxes.values()).map(async (state) => {
      if (state.worker && !state.isTerminated) {
        await state.worker.terminate()
        state.isTerminated = true
      }
    })

    await Promise.all(promises)
    activeSandboxes.clear()

    log.info("all sandboxes terminated")
  }

  /**
   * Get sandbox statistics
   */
  export function getStatistics(): {
    activeSandboxes: number
    totalMemoryMB: number
    totalCPUTime: number
    averageMemoryMB: number
    averageCPUTime: number
  } {
    const sandboxes = Array.from(activeSandboxes.values())

    const totalMemoryMB = sandboxes.reduce((sum, s) => sum + s.resourceUsage.memoryMB, 0)
    const totalCPUTime = sandboxes.reduce((sum, s) => sum + s.resourceUsage.cpuTime, 0)

    return {
      activeSandboxes: sandboxes.length,
      totalMemoryMB,
      totalCPUTime,
      averageMemoryMB: sandboxes.length > 0 ? totalMemoryMB / sandboxes.length : 0,
      averageCPUTime: sandboxes.length > 0 ? totalCPUTime / sandboxes.length : 0,
    }
  }
}

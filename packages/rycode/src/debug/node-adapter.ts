import { spawn } from "child_process"
import { DebugAdapter } from "./adapter"
import type { DebugAdapterConfig } from "./types"
import { Log } from "../util/log"

const log = Log.create({ service: "node-debug-adapter" })

/**
 * Debug adapter for Node.js and Bun using the built-in inspector
 */
export class NodeDebugAdapter extends DebugAdapter {
  private inspectorPort: number
  private runtime: "node" | "bun"

  constructor(config: DebugAdapterConfig) {
    super(config)
    this.runtime = config.language === "bun" ? "bun" : "node"
    this.inspectorPort = 9229 + Math.floor(Math.random() * 1000)
  }

  /**
   * Start the Node.js debug adapter
   * Node.js has a built-in DAP implementation via --inspect
   */
  async start(): Promise<void> {
    log.info("starting node debug adapter", {
      program: this.session.program,
      port: this.inspectorPort,
      runtime: this.runtime,
    })

    // For Node.js, we spawn the process with --inspect flag
    // The inspector will expose a WebSocket endpoint that speaks DAP
    const args = [
      `--inspect-brk=${this.inspectorPort}`,
      this.session.program,
    ]

    this.process = spawn(this.runtime, args, {
      cwd: process.cwd(),
      stdio: ["pipe", "pipe", "pipe"],
    })

    this.session.pid = this.process.pid
    this.session.port = this.inspectorPort

    // Handle process output
    if (this.process.stdout) {
      this.process.stdout.on("data", (data) => {
        log.debug("stdout", { data: data.toString() })
        this.emit("stdout", data.toString())
      })
    }

    if (this.process.stderr) {
      this.process.stderr.on("data", (data) => {
        const output = data.toString()
        log.debug("stderr", { data: output })

        // Look for "Debugger listening on" message
        if (output.includes("Debugger listening")) {
          this.emit("ready", { port: this.inspectorPort })
        }

        this.emit("stderr", output)
      })
    }

    this.process.on("error", (error) => {
      log.error("process error", { error })
      this.emit("error", error)
    })

    this.process.on("exit", (code, signal) => {
      log.info("process exited", { code, signal })
      this.session.status = "stopped"
      this.emit("exit", { code, signal })
    })

    // Wait for debugger to be ready
    await this.waitForReady()

    // Now we need to connect to the WebSocket endpoint
    // For a complete implementation, we'd use a WebSocket client here
    // For MVP, we'll simulate the connection
    log.info("node debug adapter ready", { port: this.inspectorPort })
  }

  /**
   * Wait for the inspector to be ready
   */
  private waitForReady(): Promise<void> {
    return new Promise((resolve, reject) => {
      const timeout = setTimeout(() => {
        reject(new Error("Debugger failed to start within timeout"))
      }, 10000)

      const handler = () => {
        clearTimeout(timeout)
        this.off("ready", handler)
        resolve()
      }

      this.once("ready", handler)
    })
  }

  /**
   * Get the inspector URL
   */
  getInspectorUrl(): string {
    return `ws://127.0.0.1:${this.inspectorPort}`
  }

  /**
   * Get the Chrome DevTools URL
   */
  getDevToolsUrl(): string {
    return `devtools://devtools/bundled/inspector.html?ws=127.0.0.1:${this.inspectorPort}`
  }
}

/**
 * Create a Node.js debug adapter
 */
export async function createNodeAdapter(config: DebugAdapterConfig): Promise<NodeDebugAdapter> {
  const adapter = new NodeDebugAdapter(config)
  await adapter.start()
  await adapter.initialize()

  return adapter
}

import { spawn, type ChildProcess } from "child_process"
import { EventEmitter } from "events"
import type { DAP, DebugSession, DebugAdapterConfig } from "./types"
import { Log } from "../util/log"
import { Bus } from "../bus"

const log = Log.create({ service: "debug-adapter" })

/**
 * Base class for Debug Adapter Protocol clients
 */
export abstract class DebugAdapter extends EventEmitter {
  protected process?: ChildProcess
  protected messageSeq = 1
  protected session: DebugSession
  protected messageBuffer = ""

  constructor(config: DebugAdapterConfig) {
    super()
    this.session = {
      id: this.generateSessionId(),
      language: config.language,
      program: config.program,
      status: "initializing",
      breakpoints: new Map(),
      threads: new Map(),
      stackFrames: new Map(),
      variables: new Map(),
    }
  }

  /**
   * Start the debug adapter process
   */
  abstract start(): Promise<void>

  /**
   * Send a DAP request and wait for response
   */
  protected async sendRequest<T = any>(command: string, args?: any): Promise<T> {
    const seq = this.messageSeq++
    const request: DAP.Request = {
      seq,
      type: "request",
      command,
      arguments: args,
    }

    return new Promise((resolve, reject) => {
      const timeout = setTimeout(() => {
        reject(new Error(`Request ${command} timed out`))
      }, 5000)

      const handler = (response: DAP.Response) => {
        if (response.request_seq === seq) {
          clearTimeout(timeout)
          this.off("response", handler)

          if (response.success) {
            resolve(response.body as T)
          } else {
            reject(new Error(response.message || `Request ${command} failed`))
          }
        }
      }

      this.on("response", handler)
      this.sendMessage(request)
    })
  }

  /**
   * Send a DAP message to the adapter
   */
  protected sendMessage(message: DAP.ProtocolMessage): void {
    if (!this.process || !this.process.stdin) {
      throw new Error("Debug adapter not started")
    }

    const json = JSON.stringify(message)
    const header = `Content-Length: ${Buffer.byteLength(json, "utf8")}\r\n\r\n`
    const data = header + json

    log.debug("sending message", { command: (message as any).command, seq: message.seq })
    this.process.stdin.write(data, "utf8")
  }

  /**
   * Handle incoming data from adapter
   */
  protected handleData(data: Buffer): void {
    this.messageBuffer += data.toString("utf8")

    while (true) {
      // Look for Content-Length header
      const headerMatch = this.messageBuffer.match(/Content-Length: (\d+)\r\n\r\n/)
      if (!headerMatch) break

      const contentLength = parseInt(headerMatch[1])
      const headerLength = headerMatch[0].length
      const totalLength = headerLength + contentLength

      // Check if we have the complete message
      if (this.messageBuffer.length < totalLength) break

      // Extract and parse message
      const messageData = this.messageBuffer.substring(headerLength, totalLength)
      this.messageBuffer = this.messageBuffer.substring(totalLength)

      try {
        const message = JSON.parse(messageData) as DAP.ProtocolMessage
        this.handleMessage(message)
      } catch (error) {
        log.error("failed to parse message", { error, data: messageData })
      }
    }
  }

  /**
   * Handle a parsed DAP message
   */
  protected handleMessage(message: DAP.ProtocolMessage): void {
    log.debug("received message", { type: message.type, seq: message.seq })

    if (message.type === "response") {
      this.emit("response", message as DAP.Response)
    } else if (message.type === "event") {
      this.handleEvent(message as DAP.Event)
    }
  }

  /**
   * Handle DAP events
   */
  protected handleEvent(event: DAP.Event): void {
    log.info("received event", { event: event.event })

    switch (event.event) {
      case "initialized":
        this.session.status = "running"
        this.emit("initialized")
        break

      case "stopped":
        this.session.status = "paused"
        this.emit("stopped", event.body)

        // Emit Bus event for TUI
        if (event.body?.source?.path && event.body?.line !== undefined) {
          Bus.emit("debug.stopped", {
            sessionId: this.session.id,
            file: event.body.source.path,
            line: event.body.line,
            reason: event.body.reason || "breakpoint",
          })
        }
        break

      case "continued":
        this.session.status = "running"
        this.emit("continued", event.body)

        // Emit Bus event for TUI
        Bus.emit("debug.continued", {
          sessionId: this.session.id,
        })
        break

      case "thread":
        if (event.body?.thread) {
          this.session.threads.set(event.body.thread.id, event.body.thread)
        }
        this.emit("thread", event.body)
        break

      case "breakpoint":
        this.emit("breakpoint", event.body)
        break

      case "terminated":
      case "exited":
        this.session.status = "stopped"
        this.emit("terminated", event.body)

        // Emit Bus event for TUI
        Bus.emit("debug.terminated", {
          sessionId: this.session.id,
        })
        break

      default:
        this.emit("event", event)
    }
  }

  /**
   * Initialize the debug session
   */
  async initialize(): Promise<void> {
    await this.sendRequest("initialize", {
      clientID: "rycode",
      clientName: "RyCode",
      adapterID: this.session.language,
      locale: "en-US",
      linesStartAt1: true,
      columnsStartAt1: true,
      pathFormat: "path",
      supportsVariableType: true,
      supportsVariablePaging: false,
      supportsRunInTerminalRequest: false,
    })

    log.info("debug adapter initialized", { language: this.session.language })
  }

  /**
   * Launch or attach to the program
   */
  async launch(args: DAP.LaunchRequestArguments): Promise<void> {
    await this.sendRequest("launch", args)
    log.info("debug session launched", { program: args.program })
  }

  /**
   * Set breakpoints in a source file
   */
  async setBreakpoints(
    source: DAP.Source,
    breakpoints: DAP.SourceBreakpoint[],
  ): Promise<DAP.Breakpoint[]> {
    const response = await this.sendRequest<{ breakpoints: DAP.Breakpoint[] }>("setBreakpoints", {
      source,
      breakpoints,
      sourceModified: false,
    })

    if (source.path) {
      this.session.breakpoints.set(source.path, response.breakpoints)
    }

    return response.breakpoints
  }

  /**
   * Continue execution
   */
  async continue(threadId: number = 1): Promise<void> {
    await this.sendRequest("continue", { threadId })
  }

  /**
   * Step over
   */
  async stepOver(threadId: number = 1): Promise<void> {
    await this.sendRequest("next", { threadId })
  }

  /**
   * Step into
   */
  async stepInto(threadId: number = 1): Promise<void> {
    await this.sendRequest("stepIn", { threadId })
  }

  /**
   * Step out
   */
  async stepOut(threadId: number = 1): Promise<void> {
    await this.sendRequest("stepOut", { threadId })
  }

  /**
   * Get stack trace
   */
  async stackTrace(threadId: number): Promise<DAP.StackFrame[]> {
    const response = await this.sendRequest<{ stackFrames: DAP.StackFrame[] }>("stackTrace", {
      threadId,
    })

    this.session.stackFrames.set(threadId, response.stackFrames)
    return response.stackFrames
  }

  /**
   * Get scopes for a stack frame
   */
  async scopes(frameId: number): Promise<DAP.Scope[]> {
    const response = await this.sendRequest<{ scopes: DAP.Scope[] }>("scopes", { frameId })
    return response.scopes
  }

  /**
   * Get variables
   */
  async variables(variablesReference: number): Promise<DAP.Variable[]> {
    const response = await this.sendRequest<{ variables: DAP.Variable[] }>("variables", {
      variablesReference,
    })

    this.session.variables.set(variablesReference, response.variables)
    return response.variables
  }

  /**
   * Evaluate an expression
   */
  async evaluate(expression: string, frameId?: number): Promise<{ result: string; type?: string }> {
    const response = await this.sendRequest<{ result: string; type?: string }>("evaluate", {
      expression,
      frameId,
      context: "watch",
    })

    return response
  }

  /**
   * Disconnect from the debug session
   */
  async disconnect(): Promise<void> {
    try {
      await this.sendRequest("disconnect", {
        terminateDebuggee: true,
      })
    } catch (error) {
      log.error("disconnect failed", { error })
    }

    if (this.process) {
      this.process.kill()
      this.process = undefined
    }

    this.session.status = "stopped"
  }

  /**
   * Get current session info
   */
  getSession(): DebugSession {
    return this.session
  }

  /**
   * Generate a unique session ID
   */
  private generateSessionId(): string {
    return `debug_${Date.now()}_${Math.random().toString(36).substring(2, 9)}`
  }
}

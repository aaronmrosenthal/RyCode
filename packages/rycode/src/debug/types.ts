/**
 * Debug Adapter Protocol (DAP) types and interfaces
 * Based on: https://microsoft.github.io/debug-adapter-protocol/specification
 */

export namespace DAP {
  // Base protocol message
  export interface ProtocolMessage {
    seq: number
    type: "request" | "response" | "event"
  }

  // Request message
  export interface Request extends ProtocolMessage {
    type: "request"
    command: string
    arguments?: any
  }

  // Response message
  export interface Response extends ProtocolMessage {
    type: "response"
    request_seq: number
    success: boolean
    command: string
    message?: string
    body?: any
  }

  // Event message
  export interface Event extends ProtocolMessage {
    type: "event"
    event: string
    body?: any
  }

  // Breakpoint
  export interface Breakpoint {
    id?: number
    verified: boolean
    line: number
    column?: number
    message?: string
  }

  export interface SourceBreakpoint {
    line: number
    column?: number
    condition?: string
    hitCondition?: string
    logMessage?: string
  }

  // Source
  export interface Source {
    name?: string
    path?: string
    sourceReference?: number
  }

  // Stack frame
  export interface StackFrame {
    id: number
    name: string
    source?: Source
    line: number
    column: number
    endLine?: number
    endColumn?: number
  }

  // Scope
  export interface Scope {
    name: string
    variablesReference: number
    expensive: boolean
    source?: Source
    line?: number
    column?: number
    endLine?: number
    endColumn?: number
  }

  // Variable
  export interface Variable {
    name: string
    value: string
    type?: string
    variablesReference: number
    evaluateName?: string
  }

  // Thread
  export interface Thread {
    id: number
    name: string
  }

  // Launch/Attach request arguments
  export interface LaunchRequestArguments {
    program: string
    args?: string[]
    cwd?: string
    env?: { [key: string]: string }
    stopOnEntry?: boolean
    [key: string]: any
  }
}

// Adapter state
export interface DebugSession {
  id: string
  language: string
  program: string
  pid?: number
  port?: number
  status: "initializing" | "running" | "paused" | "stopped"
  breakpoints: Map<string, DAP.Breakpoint[]>
  threads: Map<number, DAP.Thread>
  stackFrames: Map<number, DAP.StackFrame[]>
  variables: Map<number, DAP.Variable[]>
}

// Adapter configuration
export interface DebugAdapterConfig {
  language: string
  program: string
  args?: string[]
  cwd?: string
  env?: Record<string, string>
  stopOnEntry?: boolean
}

import z from "zod/v4"
import { Tool } from "./tool"
import DESCRIPTION from "./debug.txt"
import { Log } from "../util/log"
import { Instance } from "../project/instance"
import { createNodeAdapter } from "../debug/node-adapter"
import type { DebugAdapter } from "../debug/adapter"

const log = Log.create({ service: "debug-tool" })

// Store active debug sessions
const activeSessions = new Map<string, DebugAdapter>()

export const DebugTool = Tool.define("debug", {
  description: DESCRIPTION,
  parameters: z.object({
    language: z
      .enum(["node", "python", "go", "rust", "bun"])
      .describe("The language/runtime to debug"),
    program: z.string().describe("The file or command to debug"),
    args: z.array(z.string()).optional().describe("Command-line arguments to pass to the program"),
    breakpoints: z
      .array(
        z.object({
          file: z.string().describe("File path for the breakpoint"),
          line: z.number().describe("Line number for the breakpoint"),
          condition: z.string().optional().describe("Optional condition for the breakpoint"),
        }),
      )
      .optional()
      .describe("Initial breakpoints to set"),
    watch: z
      .array(z.string())
      .optional()
      .describe("Expressions to watch during debugging"),
    cwd: z.string().optional().describe("Working directory for the debugger"),
  }),
  async execute(params) {
    log.info("starting debug session", {
      language: params.language,
      program: params.program,
    })

    // Validate the program file exists
    const programPath = params.program.startsWith("/")
      ? params.program
      : `${params.cwd || Instance.directory}/${params.program}`

    try {
      await Bun.file(programPath).exists()
    } catch (error) {
      throw new Error(`Program file not found: ${programPath}`)
    }

    // Create the debug adapter
    const adapter = await createDebugAdapter(params.language, {
      language: params.language,
      program: programPath,
      args: params.args || [],
      cwd: params.cwd || Instance.directory,
    })

    const session = adapter.getSession()

    // Store the session
    activeSessions.set(session.id, adapter)

    // Launch the program
    await adapter.launch({
      program: programPath,
      args: params.args || [],
      cwd: params.cwd || Instance.directory,
      stopOnEntry: true, // Stop at the beginning
    })

    // Set initial breakpoints if provided
    const breakpoints: string[] = []
    if (params.breakpoints && params.breakpoints.length > 0) {
      for (const bp of params.breakpoints) {
        const bps = await adapter.setBreakpoints(
          { path: bp.file },
          [
            {
              line: bp.line,
              condition: bp.condition,
            },
          ],
        )
        breakpoints.push(...bps.map((b) => `${bp.file}:${b.line}`))
      }
    }

    // Set up watch expressions tracking if provided
    const watches: string[] = []
    if (params.watch && params.watch.length > 0) {
      watches.push(...params.watch)
    }

    log.info("debug session initialized", {
      sessionId: session.id,
      breakpoints: breakpoints.length,
      watches: watches.length,
    })

    return {
      title: `Debugging ${params.program}`,
      metadata: {
        sessionId: session.id,
        language: params.language,
        program: params.program,
        port: session.port,
        pid: session.pid,
        breakpoints,
        watches,
        status: session.status,
      },
      output: `Debug session started for ${params.program}

Session ID: ${session.id}
Language: ${params.language}
${session.port ? `Debug port: localhost:${session.port}` : ""}
${session.pid ? `Process ID: ${session.pid}` : ""}
Breakpoints: ${breakpoints.length}
Watch expressions: ${watches.length}

The debugger is paused at the entry point. You can:
- Continue execution (the AI can call 'continue')
- Step through code (the AI can call 'step')
- Inspect variables (the AI can call 'inspect')
- Ask questions: "Why is X undefined?", "Show me when Y changed"

The AI assistant will help guide you through the debugging process.`,
    }
  },
})

// Helper function to create the appropriate debug adapter
async function createDebugAdapter(
  language: string,
  config: { language: string; program: string; args: string[]; cwd: string },
): Promise<DebugAdapter> {
  log.info("creating debug adapter", { language, program: config.program })

  switch (language) {
    case "node":
    case "bun":
      return await createNodeAdapter(config)

    case "python":
      // TODO: Implement Python debugpy adapter
      throw new Error("Python debugging not yet implemented. Coming soon!")

    case "go":
      // TODO: Implement Go delve adapter
      throw new Error("Go debugging not yet implemented. Coming soon!")

    case "rust":
      // TODO: Implement Rust CodeLLDB adapter
      throw new Error("Rust debugging not yet implemented. Coming soon!")

    default:
      throw new Error(`Unsupported language for debugging: ${language}`)
  }
}

// Export function to get active session by ID
export function getDebugSession(sessionId: string): DebugAdapter | undefined {
  return activeSessions.get(sessionId)
}

// Export function to close a debug session
export async function closeDebugSession(sessionId: string): Promise<void> {
  const adapter = activeSessions.get(sessionId)
  if (adapter) {
    await adapter.disconnect()
    activeSessions.delete(sessionId)
    log.info("debug session closed", { sessionId })
  }
}

import { Hono } from "hono"
import { describeRoute, validator, resolver } from "hono-openapi"
import z from "zod/v4"
import { getDebugSession, closeDebugSession } from "../tool/debug"
import { Log } from "../util/log"
import { Bus } from "../bus"

const log = Log.create({ service: "debug-route" })

const ERRORS = {
  400: {
    description: "Bad request",
    content: {
      "application/json": {
        schema: resolver(
          z.object({
            data: z.record(z.string(), z.any()),
          }),
        ),
      },
    },
  },
  404: {
    description: "Debug session not found",
    content: {
      "application/json": {
        schema: resolver(
          z.object({
            data: z.object({
              message: z.string(),
            }),
          }),
        ),
      },
    },
  },
} as const

// Define event types for debug events
export namespace DebugEvent {
  export const Stopped = Bus.event(
    "debug.stopped",
    z.object({
      sessionId: z.string(),
      file: z.string(),
      line: z.number(),
      reason: z.string(),
    }),
  )

  export const Continued = Bus.event(
    "debug.continued",
    z.object({
      sessionId: z.string(),
    }),
  )

  export const Terminated = Bus.event(
    "debug.terminated",
    z.object({
      sessionId: z.string(),
    }),
  )
}

export const DebugRoute = new Hono()
  // Continue execution
  .post(
    "/:sessionId/continue",
    describeRoute({
      description: "Continue execution",
      operationId: "debug.continue",
      responses: {
        200: {
          description: "Execution continued",
          content: {
            "application/json": {
              schema: resolver(z.boolean()),
            },
          },
        },
        ...ERRORS,
      },
    }),
    validator(
      "param",
      z.object({
        sessionId: z.string(),
      }),
    ),
    validator(
      "json",
      z
        .object({
          threadId: z.number().optional().default(1),
        })
        .optional(),
    ),
    async (c) => {
      const { sessionId } = c.req.valid("param")
      const body = c.req.valid("json") as { threadId?: number } | undefined

      const adapter = getDebugSession(sessionId)
      if (!adapter) {
        return c.json({ data: { message: "Debug session not found" } }, 404)
      }

      log.info("continuing execution", { sessionId })
      await adapter.continue(body?.threadId ?? 1)

      return c.json(true)
    },
  )

  // Step over
  .post(
    "/:sessionId/step-over",
    describeRoute({
      description: "Step over (execute next line)",
      operationId: "debug.stepOver",
      responses: {
        200: {
          description: "Stepped over",
          content: {
            "application/json": {
              schema: resolver(z.boolean()),
            },
          },
        },
        ...ERRORS,
      },
    }),
    validator(
      "param",
      z.object({
        sessionId: z.string(),
      }),
    ),
    validator(
      "json",
      z
        .object({
          threadId: z.number().optional().default(1),
        })
        .optional(),
    ),
    async (c) => {
      const { sessionId } = c.req.valid("param")
      const body = c.req.valid("json") as { threadId?: number } | undefined

      const adapter = getDebugSession(sessionId)
      if (!adapter) {
        return c.json({ data: { message: "Debug session not found" } }, 404)
      }

      log.info("stepping over", { sessionId })
      await adapter.stepOver(body?.threadId ?? 1)

      return c.json(true)
    },
  )

  // Step into
  .post(
    "/:sessionId/step-into",
    describeRoute({
      description: "Step into (enter function)",
      operationId: "debug.stepInto",
      responses: {
        200: {
          description: "Stepped into",
          content: {
            "application/json": {
              schema: resolver(z.boolean()),
            },
          },
        },
        ...ERRORS,
      },
    }),
    validator(
      "param",
      z.object({
        sessionId: z.string(),
      }),
    ),
    validator(
      "json",
      z
        .object({
          threadId: z.number().optional().default(1),
        })
        .optional(),
    ),
    async (c) => {
      const { sessionId } = c.req.valid("param")
      const body = c.req.valid("json") as { threadId?: number } | undefined

      const adapter = getDebugSession(sessionId)
      if (!adapter) {
        return c.json({ data: { message: "Debug session not found" } }, 404)
      }

      log.info("stepping into", { sessionId })
      await adapter.stepInto(body?.threadId ?? 1)

      return c.json(true)
    },
  )

  // Step out
  .post(
    "/:sessionId/step-out",
    describeRoute({
      description: "Step out (exit current function)",
      operationId: "debug.stepOut",
      responses: {
        200: {
          description: "Stepped out",
          content: {
            "application/json": {
              schema: resolver(z.boolean()),
            },
          },
        },
        ...ERRORS,
      },
    }),
    validator(
      "param",
      z.object({
        sessionId: z.string(),
      }),
    ),
    validator(
      "json",
      z
        .object({
          threadId: z.number().optional().default(1),
        })
        .optional(),
    ),
    async (c) => {
      const { sessionId } = c.req.valid("param")
      const body = c.req.valid("json") as { threadId?: number } | undefined

      const adapter = getDebugSession(sessionId)
      if (!adapter) {
        return c.json({ data: { message: "Debug session not found" } }, 404)
      }

      log.info("stepping out", { sessionId })
      await adapter.stepOut(body?.threadId ?? 1)

      return c.json(true)
    },
  )

  // Disconnect
  .post(
    "/:sessionId/disconnect",
    describeRoute({
      description: "Disconnect from debug session",
      operationId: "debug.disconnect",
      responses: {
        200: {
          description: "Disconnected",
          content: {
            "application/json": {
              schema: resolver(z.boolean()),
            },
          },
        },
        ...ERRORS,
      },
    }),
    validator(
      "param",
      z.object({
        sessionId: z.string(),
      }),
    ),
    async (c) => {
      const { sessionId } = c.req.valid("param")

      log.info("disconnecting", { sessionId })
      await closeDebugSession(sessionId)

      // Emit terminated event
      await Bus.publish(DebugEvent.Terminated, { sessionId })

      return c.json(true)
    },
  )

  // Get stack trace
  .post(
    "/:sessionId/stack-trace",
    describeRoute({
      description: "Get stack trace for current execution point",
      operationId: "debug.stackTrace",
      responses: {
        200: {
          description: "Stack trace",
          content: {
            "application/json": {
              schema: resolver(
                z.array(
                  z.object({
                    id: z.number(),
                    name: z.string(),
                    source: z
                      .object({
                        path: z.string().optional(),
                      })
                      .optional(),
                    line: z.number(),
                    column: z.number(),
                  }),
                ),
              ),
            },
          },
        },
        ...ERRORS,
      },
    }),
    validator(
      "param",
      z.object({
        sessionId: z.string(),
      }),
    ),
    validator(
      "json",
      z.object({
        threadId: z.number().optional().default(1),
      }),
    ),
    async (c) => {
      const { sessionId } = c.req.valid("param")
      const body = c.req.valid("json") as { threadId?: number }
      const threadId = body?.threadId ?? 1

      const adapter = getDebugSession(sessionId)
      if (!adapter) {
        return c.json({ data: { message: "Debug session not found" } }, 404)
      }

      log.info("getting stack trace", { sessionId, threadId })
      const stackFrames = await adapter.stackTrace(threadId)

      return c.json(stackFrames)
    },
  )

  // Get variables
  .post(
    "/:sessionId/variables",
    describeRoute({
      description: "Get variables for a stack frame",
      operationId: "debug.variables",
      responses: {
        200: {
          description: "Variables",
          content: {
            "application/json": {
              schema: resolver(
                z.array(
                  z.object({
                    name: z.string(),
                    value: z.string(),
                    type: z.string().optional(),
                    variablesReference: z.number().optional(),
                  }),
                ),
              ),
            },
          },
        },
        ...ERRORS,
      },
    }),
    validator(
      "param",
      z.object({
        sessionId: z.string(),
      }),
    ),
    validator(
      "json",
      z.object({
        variablesReference: z.number(),
      }),
    ),
    async (c) => {
      const { sessionId } = c.req.valid("param")
      const { variablesReference } = c.req.valid("json")

      const adapter = getDebugSession(sessionId)
      if (!adapter) {
        return c.json({ data: { message: "Debug session not found" } }, 404)
      }

      log.info("getting variables", { sessionId, variablesReference })
      const variables = await adapter.variables(variablesReference)

      return c.json(variables)
    },
  )

  // Get scopes
  .post(
    "/:sessionId/scopes",
    describeRoute({
      description: "Get scopes for a stack frame",
      operationId: "debug.scopes",
      responses: {
        200: {
          description: "Scopes",
          content: {
            "application/json": {
              schema: resolver(
                z.array(
                  z.object({
                    name: z.string(),
                    variablesReference: z.number(),
                    expensive: z.boolean().optional(),
                  }),
                ),
              ),
            },
          },
        },
        ...ERRORS,
      },
    }),
    validator(
      "param",
      z.object({
        sessionId: z.string(),
      }),
    ),
    validator(
      "json",
      z.object({
        frameId: z.number(),
      }),
    ),
    async (c) => {
      const { sessionId } = c.req.valid("param")
      const { frameId } = c.req.valid("json")

      const adapter = getDebugSession(sessionId)
      if (!adapter) {
        return c.json({ data: { message: "Debug session not found" } }, 404)
      }

      log.info("getting scopes", { sessionId, frameId })
      const scopes = await adapter.scopes(frameId)

      return c.json(scopes)
    },
  )

  // Get session status
  .get(
    "/:sessionId/status",
    describeRoute({
      description: "Get debug session status",
      operationId: "debug.status",
      responses: {
        200: {
          description: "Session status",
          content: {
            "application/json": {
              schema: resolver(
                z.object({
                  id: z.string(),
                  language: z.string(),
                  program: z.string(),
                  status: z.enum(["initializing", "running", "paused", "stopped"]),
                  port: z.number().optional(),
                  pid: z.number().optional(),
                }),
              ),
            },
          },
        },
        ...ERRORS,
      },
    }),
    validator(
      "param",
      z.object({
        sessionId: z.string(),
      }),
    ),
    async (c) => {
      const { sessionId } = c.req.valid("param")

      const adapter = getDebugSession(sessionId)
      if (!adapter) {
        return c.json({ data: { message: "Debug session not found" } }, 404)
      }

      const session = adapter.getSession()
      return c.json({
        id: session.id,
        language: session.language,
        program: session.program,
        status: session.status,
        port: session.port,
        pid: session.pid,
      })
    },
  )

import { Config } from "../config/config"
import z from "zod/v4"
import { Provider } from "../provider/provider"
import { generateObject, type ModelMessage } from "ai"
import PROMPT_GENERATE from "./generate.txt"
import { SystemPrompt } from "../session/system"
import { Instance } from "../project/instance"
import { mergeDeep } from "remeda"

export namespace Agent {
  export const Info = z
    .object({
      name: z.string(),
      description: z.string().optional(),
      mode: z.union([z.literal("subagent"), z.literal("primary"), z.literal("all")]),
      builtIn: z.boolean(),
      topP: z.number().optional(),
      temperature: z.number().optional(),
      permission: z.object({
        edit: Config.Permission,
        bash: z.record(z.string(), Config.Permission),
        webfetch: Config.Permission.optional(),
      }),
      model: z
        .object({
          modelID: z.string(),
          providerID: z.string(),
        })
        .optional(),
      prompt: z.string().optional(),
      tools: z.record(z.string(), z.boolean()),
      options: z.record(z.string(), z.any()),
    })
    .meta({
      ref: "Agent",
    })
  export type Info = z.infer<typeof Info>

  const state = Instance.state(async () => {
    const cfg = await Config.get()
    const defaultTools = cfg.tools ?? {}
    const defaultPermission: Info["permission"] = {
      edit: "allow",
      bash: {
        "*": "allow",
      },
      webfetch: "allow",
    }
    const agentPermission = mergeAgentPermissions(defaultPermission, cfg.permission ?? {})

    const planPermission = mergeAgentPermissions(
      {
        edit: "deny",
        bash: "ask",
        webfetch: "allow",
      },
      cfg.permission ?? {},
    )

    const result: Record<string, Info> = {
      general: {
        name: "general",
        description:
          "General-purpose agent for researching complex questions, searching for code, and executing multi-step tasks. When you are searching for a keyword or file and are not confident that you will find the right match in the first few tries use this agent to perform the search for you.",
        tools: {
          todoread: false,
          todowrite: false,
          ...defaultTools,
        },
        options: {},
        permission: agentPermission,
        mode: "subagent",
        builtIn: true,
      },
      build: {
        name: "build",
        tools: { ...defaultTools },
        options: {},
        permission: agentPermission,
        mode: "primary",
        builtIn: true,
      },
      plan: {
        name: "plan",
        options: {},
        permission: planPermission,
        tools: {
          ...defaultTools,
        },
        mode: "primary",
        builtIn: true,
      },
    }
    for (const [key, value] of Object.entries(cfg.agent ?? {})) {
      if (value.disable) {
        delete result[key]
        continue
      }
      let item = result[key]
      if (!item)
        item = result[key] = {
          name: key,
          mode: "all",
          permission: agentPermission,
          options: {},
          tools: {},
          builtIn: false,
        }
      const { name, model, prompt, tools, description, temperature, top_p, mode, permission, ...extra } = value
      item.options = {
        ...item.options,
        ...extra,
      }
      if (model) item.model = Provider.parseModel(model)
      if (prompt) item.prompt = prompt
      if (tools)
        item.tools = {
          ...item.tools,
          ...tools,
        }
      item.tools = {
        ...defaultTools,
        ...item.tools,
      }
      if (description) item.description = description
      if (temperature != undefined) item.temperature = temperature
      if (top_p != undefined) item.topP = top_p
      if (mode) item.mode = mode
      // just here for consistency & to prevent it from being added as an option
      if (name) item.name = name

      if (permission ?? cfg.permission) {
        item.permission = mergeAgentPermissions(cfg.permission ?? {}, permission ?? {})
      }
    }
    return result
  })

  /**
   * Retrieves a specific agent configuration by name.
   *
   * @param agent - Agent identifier (e.g., "general", "plan", "build")
   * @returns Agent configuration or undefined if not found
   */
  export async function get(agent: string): Promise<Info | undefined> {
    return state().then((x) => x[agent])
  }

  /**
   * Lists all available agent configurations.
   *
   * Includes both built-in agents (general, plan, build) and custom agents
   * defined in the config file.
   *
   * @returns Array of agent configurations
   */
  export async function list(): Promise<Info[]> {
    return state().then((x) => Object.values(x))
  }

  /**
   * Generates a new agent configuration from a natural language description.
   *
   * Uses AI to create an agent with appropriate tools, permissions, and prompts
   * based on the user's description. Ensures unique identifiers.
   *
   * @param input - Object with description of desired agent behavior
   * @returns Generated agent configuration (identifier, whenToUse, systemPrompt)
   *
   * @example
   * ```typescript
   * const config = await Agent.generate({
   *   description: "Create an agent that specializes in writing tests"
   * })
   * // Returns: { identifier: "test-writer", whenToUse: "...", systemPrompt: "..." }
   * ```
   */
  export async function generate(input: { description: string }): Promise<{ identifier: string; whenToUse: string; systemPrompt: string }> {
    const defaultModel = await Provider.defaultModel()
    const model = await Provider.getModel(defaultModel.providerID, defaultModel.modelID)
    const system = SystemPrompt.header(defaultModel.providerID)
    system.push(PROMPT_GENERATE)
    const existing = await list()
    const result = await generateObject({
      temperature: 0.3,
      prompt: [
        ...system.map(
          (item): ModelMessage => ({
            role: "system",
            content: item,
          }),
        ),
        {
          role: "user",
          content: `Create an agent configuration based on this request: \"${input.description}\".\n\nIMPORTANT: The following identifiers already exist and must NOT be used: ${existing.map((i) => i.name).join(", ")}\n  Return ONLY the JSON object, no other text, do not wrap in backticks`,
        },
      ],
      model: model.language,
      schema: z.object({
        identifier: z.string(),
        whenToUse: z.string(),
        systemPrompt: z.string(),
      }),
    })
    return result.object
  }
}

function mergeAgentPermissions(basePermission: any, overridePermission: any): Agent.Info["permission"] {
  const merged = mergeDeep(basePermission ?? {}, overridePermission ?? {}) as any
  let mergedBash
  if (merged.bash) {
    if (typeof merged.bash === "string") {
      mergedBash = {
        "*": merged.bash,
      }
    }
    // if granular permissions are provided, default to "ask"
    if (typeof merged.bash === "object") {
      mergedBash = mergeDeep(
        {
          "*": "ask",
        },
        merged.bash,
      )
    }
  }

  const result: Agent.Info["permission"] = {
    edit: merged.edit ?? "allow",
    webfetch: merged.webfetch ?? "allow",
    bash: mergedBash ?? { "*": "allow" },
  }

  return result
}

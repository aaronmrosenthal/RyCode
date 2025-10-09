import { Decimal } from "decimal.js"
import z from "zod/v4"
import { type LanguageModelUsage, type ProviderMetadata } from "ai"

import PROMPT_INITIALIZE from "../session/prompt/initialize.txt"

import { Bus } from "../bus"
import { Config } from "../config/config"
import { Flag } from "../flag/flag"
import { Identifier } from "../id/id"
import { Installation } from "../installation"
import type { ModelsDev } from "../provider/models"
import { Share } from "../share/share"
import { Storage } from "../storage/storage"
import { Log } from "../util/log"
import { MessageV2 } from "./message-v2"
import { Project } from "../project/project"
import { Instance } from "../project/instance"
import { SessionPrompt } from "./prompt"

export namespace Session {
  const log = Log.create({ service: "session" })

  const parentSessionTitlePrefix = "New session - "
  const childSessionTitlePrefix = "Child session - "

  function createDefaultTitle(isChild = false) {
    return (isChild ? childSessionTitlePrefix : parentSessionTitlePrefix) + new Date().toISOString()
  }

  export const Info = z
    .object({
      id: Identifier.schema("session"),
      projectID: z.string(),
      directory: z.string(),
      parentID: Identifier.schema("session").optional(),
      share: z
        .object({
          url: z.string(),
        })
        .optional(),
      title: z.string(),
      version: z.string(),
      time: z.object({
        created: z.number(),
        updated: z.number(),
        compacting: z.number().optional(),
      }),
      revert: z
        .object({
          messageID: z.string(),
          partID: z.string().optional(),
          snapshot: z.string().optional(),
          diff: z.string().optional(),
        })
        .optional(),
    })
    .meta({
      ref: "Session",
    })
  export type Info = z.output<typeof Info>

  export const ShareInfo = z
    .object({
      secret: z.string(),
      url: z.string(),
    })
    .meta({
      ref: "SessionShare",
    })
  export type ShareInfo = z.output<typeof ShareInfo>

  export const Event = {
    Updated: Bus.event(
      "session.updated",
      z.object({
        info: Info,
      }),
    ),
    Deleted: Bus.event(
      "session.deleted",
      z.object({
        info: Info,
      }),
    ),
    Error: Bus.event(
      "session.error",
      z.object({
        sessionID: z.string().optional(),
        error: MessageV2.Assistant.shape.error,
      }),
    ),
  }

  /**
   * Creates a new AI coding session in the current project.
   *
   * Sessions represent individual conversations with the AI. Child sessions
   * are used for subagents that work on specific subtasks.
   *
   * @param parentID - Optional parent session ID for creating child sessions (subagents)
   * @param title - Optional human-readable session title (auto-generated if not provided)
   * @returns Session.Info with unique identifier and metadata
   *
   * @example
   * ```typescript
   * // Create root session
   * const session = await Session.create()
   *
   * // Create child session for subagent
   * const child = await Session.create(session.id, "Refactoring Task")
   * ```
   */
  export async function create(parentID?: string, title?: string): Promise<Info> {
    return createNext({
      parentID,
      directory: Instance.directory,
      title,
    })
  }

  /**
   * Updates the session's last-modified timestamp.
   *
   * Used to track session activity and sort by recency.
   *
   * @param sessionID - Unique session identifier
   */
  export async function touch(sessionID: string): Promise<void> {
    await update(sessionID, (draft) => {
      draft.time.updated = Date.now()
    })
  }

  /**
   * Internal helper to create a session with custom directory.
   *
   * Used internally and by migration code. Most callers should use `create()` instead.
   *
   * @param input - Session creation parameters including optional custom directory
   * @returns Session.Info with unique identifier and metadata
   * @internal
   */
  export async function createNext(input: { id?: string; title?: string; parentID?: string; directory: string }): Promise<Info> {
    const result: Info = {
      id: Identifier.descending("session", input.id),
      version: Installation.VERSION,
      projectID: Instance.project.id,
      directory: input.directory,
      parentID: input.parentID,
      title: input.title ?? createDefaultTitle(!!input.parentID),
      time: {
        created: Date.now(),
        updated: Date.now(),
      },
    }
    log.info("created", result)
    await Storage.write(["session", Instance.project.id, result.id], result)
    const cfg = await Config.get()
    if (!result.parentID && (Flag.OPENCODE_AUTO_SHARE || cfg.share === "auto"))
      share(result.id)
        .then((share) => {
          update(result.id, (draft) => {
            draft.share = share
          })
        })
        .catch(() => {
          // Silently ignore sharing errors during session creation
        })
    Bus.publish(Event.Updated, {
      info: result,
    })
    return result
  }

  /**
   * Retrieves session information by ID.
   *
   * @param id - Unique session identifier
   * @returns Session.Info containing session metadata
   * @throws Error if session not found
   */
  export async function get(id: string): Promise<Info> {
    const read = await Storage.read<Info>(["session", Instance.project.id, id])
    return read as Info
  }

  /**
   * Retrieves sharing information for a session.
   *
   * @param id - Unique session identifier
   * @returns ShareInfo containing secret and public URL
   * @throws Error if session is not shared
   */
  export async function getShare(id: string): Promise<ShareInfo> {
    return Storage.read<ShareInfo>(["share", id])
  }

  /**
   * Creates or retrieves a shareable URL for a session.
   *
   * Synchronizes session data to remote share service if not already shared.
   * Automatically syncs all messages and parts to the share service.
   *
   * @param id - Unique session identifier
   * @returns ShareInfo with secret and public URL
   * @throws Error if sharing is disabled in configuration
   */
  export async function share(id: string): Promise<ShareInfo> {
    const cfg = await Config.get()
    if (cfg.share === "disabled") {
      throw new Error("Sharing is disabled in configuration")
    }

    const session = await get(id)
    if (session.share) {
      // Session already shared, retrieve the secret from storage
      const shareInfo = await getShare(id)
      return shareInfo
    }
    const share = await Share.create(id)
    await update(id, (draft) => {
      draft.share = {
        url: share.url,
      }
    })
    await Storage.write(["share", id], share)
    await Share.sync("session/info/" + id, session)
    for (const msg of await messages(id)) {
      await Share.sync("session/message/" + id + "/" + msg.info.id, msg.info)
      for (const part of msg.parts) {
        await Share.sync("session/part/" + id + "/" + msg.info.id + "/" + part.id, part)
      }
    }
    return share
  }

  /**
   * Removes sharing for a session, deleting the public URL.
   *
   * Uses atomic transactions to ensure data consistency during unshare.
   * Remote deletion is best-effort and won't fail the operation if it errors.
   *
   * @param id - Unique session identifier
   * @throws Error if unshare transaction fails
   */
  export async function unshare(id: string): Promise<void> {
    const share = await getShare(id).catch(() => undefined)
    if (!share) return

    // BUG FIX: Use transaction for atomic local updates
    const tx = Storage.transaction()

    try {
      // Remove share metadata
      await tx.remove(["share", id])

      // Update session to remove share info
      const session = await get(id)
      session.share = undefined
      await tx.write(["session", Instance.project.id, id], session)

      // Commit atomic local changes
      await tx.commit()

      // Remote deletion (best effort, don't fail if this errors)
      await Share.remove(id, share.secret).catch((error) => {
        log.warn("Failed to remove remote share", {
          sessionID: id,
          error: error.message,
        })
      })

      // Publish update event
      Bus.publish(Event.Updated, {
        info: session,
      })
    } catch (error: any) {
      await tx.rollback()
      log.error("Failed to unshare session", {
        sessionID: id,
        error: error.message,
      })
      throw new Error(`Failed to unshare session ${id}: ${error.message}`, { cause: error })
    }
  }

  /**
   * Updates session metadata using a mutation function.
   *
   * Automatically updates the last-modified timestamp and publishes update events.
   *
   * @param id - Unique session identifier
   * @param editor - Mutation function to modify session data
   * @returns Updated Session.Info
   *
   * @example
   * ```typescript
   * await Session.update(sessionId, (draft) => {
   *   draft.title = "New Title"
   *   draft.share = { url: "https://..." }
   * })
   * ```
   */
  export async function update(id: string, editor: (session: Info) => void): Promise<Info> {
    const project = Instance.project
    const result = await Storage.update<Info>(["session", project.id, id], (draft) => {
      editor(draft)
      draft.time.updated = Date.now()
    })
    Bus.publish(Event.Updated, {
      info: result,
    })
    return result
  }

  /**
   * Retrieves all messages in a session, sorted chronologically.
   *
   * Each message includes its parts (tool calls, text blocks, etc.).
   *
   * @param sessionID - Unique session identifier
   * @returns Array of messages with their parts, sorted by creation time
   */
  export async function messages(sessionID: string): Promise<MessageV2.WithParts[]> {
    const result = [] as MessageV2.WithParts[]
    for (const p of await Storage.list(["message", sessionID])) {
      const read = await Storage.read<MessageV2.Info>(p)
      result.push({
        info: read,
        parts: await getParts(read.id),
      })
    }
    result.sort((a, b) => (a.info.id > b.info.id ? 1 : -1))
    return result
  }

  /**
   * Retrieves a specific message with its parts.
   *
   * @param sessionID - Unique session identifier
   * @param messageID - Unique message identifier
   * @returns Message with all its parts
   * @throws Error if message not found
   */
  export async function getMessage(sessionID: string, messageID: string): Promise<MessageV2.WithParts> {
    return {
      info: await Storage.read<MessageV2.Info>(["message", sessionID, messageID]),
      parts: await getParts(messageID),
    }
  }

  /**
   * Retrieves all parts of a message, sorted by creation order.
   *
   * Parts include tool calls, text blocks, errors, and other message components.
   *
   * @param messageID - Unique message identifier
   * @returns Array of message parts, sorted chronologically
   */
  export async function getParts(messageID: string): Promise<MessageV2.Part[]> {
    const result = [] as MessageV2.Part[]
    for (const item of await Storage.list(["part", messageID])) {
      const read = await Storage.read<MessageV2.Part>(item)
      result.push(read)
    }
    result.sort((a, b) => (a.id > b.id ? 1 : -1))
    return result
  }

  export async function* list() {
    const project = Instance.project
    for (const item of await Storage.list(["session", project.id])) {
      yield Storage.read<Info>(item)
    }
  }

  /**
   * Retrieves all child sessions of a parent session.
   *
   * Child sessions are created by subagents working on specific subtasks.
   *
   * @param parentID - Unique parent session identifier
   * @returns Array of child session info objects
   */
  export async function children(parentID: string): Promise<Session.Info[]> {
    const project = Instance.project
    const result = [] as Session.Info[]
    for (const item of await Storage.list(["session", project.id])) {
      const session = await Storage.read<Info>(item)
      if (session.parentID !== parentID) continue
      result.push(session)
    }
    return result
  }

  /**
   * Deletes a session and all its descendants recursively.
   *
   * Uses atomic transactions to ensure data consistency. Automatically unshares
   * the session and removes all messages, parts, and child sessions.
   *
   * @param sessionID - Unique session identifier
   * @param emitEvent - Whether to publish deletion event (default: true)
   * @throws Error if session deletion transaction fails
   */
  export async function remove(sessionID: string, emitEvent = true): Promise<void> {
    const project = Instance.project
    let session: Info | undefined

    try {
      session = await get(sessionID)

      // Collect all descendant sessions recursively
      const allSessions = await collectAllDescendants(sessionID)

      // BUG FIX: Use transaction for atomic deletion
      const tx = Storage.transaction()

      try {
        // Delete all sessions and their data atomically
        for (const sid of allSessions) {
          // Delete all messages and parts for this session
          for (const msg of await Storage.list(["message", sid])) {
            for (const part of await Storage.list(["part", msg.at(-1)!])) {
              await tx.remove(part)
            }
            await tx.remove(msg)
          }
          // Delete the session itself
          await tx.remove(["session", project.id, sid])
        }

        await tx.commit()

        // BUG FIX: Unshare after successful deletion (best effort)
        await unshare(sessionID).catch((error) => {
          log.warn("Failed to unshare session", { sessionID, error: error.message })
        })

        if (emitEvent && session) {
          Bus.publish(Event.Deleted, {
            info: session,
          })
        }
      } catch (txError) {
        await tx.rollback()
        throw txError
      }
    } catch (e: any) {
      log.error("Failed to remove session", {
        sessionID,
        error: e.message,
        stack: e.stack,
      })
      throw new Error(`Failed to remove session ${sessionID}: ${e.message}`, { cause: e })
    }
  }

  // Helper function to collect all descendant sessions
  async function collectAllDescendants(sessionID: string): Promise<string[]> {
    const result = [sessionID]
    const childSessions = await children(sessionID)

    for (const child of childSessions) {
      const descendants = await collectAllDescendants(child.id)
      result.push(...descendants)
    }

    return result
  }

  /**
   * Updates a message and publishes update event.
   *
   * @param msg - Message info object with updated data
   * @returns Updated message info
   */
  export async function updateMessage(msg: MessageV2.Info): Promise<MessageV2.Info> {
    await Storage.write(["message", msg.sessionID, msg.id], msg)
    Bus.publish(MessageV2.Event.Updated, {
      info: msg,
    })
    return msg
  }

  /**
   * Deletes a message from a session.
   *
   * @param sessionID - Unique session identifier
   * @param messageID - Unique message identifier to delete
   * @returns The deleted message ID
   */
  export async function removeMessage(sessionID: string, messageID: string): Promise<string> {
    await Storage.remove(["message", sessionID, messageID])
    Bus.publish(MessageV2.Event.Removed, {
      sessionID,
      messageID,
    })
    return messageID
  }

  /**
   * Updates a message part and publishes update event.
   *
   * @param part - Message part object with updated data
   * @returns Updated message part
   */
  export async function updatePart(part: MessageV2.Part): Promise<MessageV2.Part> {
    await Storage.write(["part", part.messageID, part.id], part)
    Bus.publish(MessageV2.Event.PartUpdated, {
      part,
    })
    return part
  }

  /**
   * Calculates token usage and cost for a model invocation.
   *
   * Handles cache tokens and reasoning tokens where applicable.
   * Costs are calculated based on per-million-token pricing.
   *
   * @param model - Model configuration with pricing info
   * @param usage - Token usage from the LLM response
   * @param metadata - Optional provider-specific metadata (e.g., Anthropic cache tokens)
   * @returns Object containing cost in dollars and detailed token breakdown
   */
  export function getUsage(
    model: ModelsDev.Model,
    usage: LanguageModelUsage,
    metadata?: ProviderMetadata
  ): { cost: number; tokens: { input: number; output: number; reasoning: number; cache: { write: number; read: number } } } {
    const tokens = {
      input: usage.inputTokens ?? 0,
      output: usage.outputTokens ?? 0,
      reasoning: usage?.reasoningTokens ?? 0,
      cache: {
        write: (metadata?.["anthropic"]?.["cacheCreationInputTokens"] ??
          // @ts-expect-error
          metadata?.["bedrock"]?.["usage"]?.["cacheWriteInputTokens"] ??
          0) as number,
        read: usage.cachedInputTokens ?? 0,
      },
    }
    return {
      cost: new Decimal(0)
        .add(new Decimal(tokens.input).mul(model.cost?.input ?? 0).div(1_000_000))
        .add(new Decimal(tokens.output).mul(model.cost?.output ?? 0).div(1_000_000))
        .add(new Decimal(tokens.cache.read).mul(model.cost?.cache_read ?? 0).div(1_000_000))
        .add(new Decimal(tokens.cache.write).mul(model.cost?.cache_write ?? 0).div(1_000_000))
        .toNumber(),
      tokens,
    }
  }

  export class BusyError extends Error {
    constructor(public readonly sessionID: string) {
      super(`Session ${sessionID} is busy`)
    }
  }

  /**
   * Initializes a session with the system prompt.
   *
   * Sends the initial prompt that sets up the AI's context and capabilities.
   * Marks the project as initialized after completion.
   *
   * @param input - Initialization parameters including session, model, and message IDs
   */
  export async function initialize(input: {
    sessionID: string
    modelID: string
    providerID: string
    messageID: string
  }): Promise<void> {
    await SessionPrompt.prompt({
      sessionID: input.sessionID,
      messageID: input.messageID,
      model: {
        providerID: input.providerID,
        modelID: input.modelID,
      },
      parts: [
        {
          id: Identifier.ascending("part"),
          type: "text",
          text: PROMPT_INITIALIZE.replace("${path}", Instance.worktree),
        },
      ],
    })
    await Project.setInitialized(Instance.project.id)
  }
}

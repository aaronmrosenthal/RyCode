/**
 * Branded types for domain-specific identifiers.
 *
 * These types provide compile-time safety by preventing accidental mixing of
 * different ID types (e.g., passing a SessionID where a ProjectID is expected).
 *
 * @example
 * ```typescript
 * function getSession(id: SessionID): Promise<Session.Info> {
 *   // TypeScript ensures only SessionID can be passed here
 * }
 *
 * const sessionId = "abc123" as SessionID
 * const projectId = "xyz789" as ProjectID
 *
 * getSession(sessionId) // ✓ OK
 * getSession(projectId) // ✗ Compile error
 * ```
 */

/**
 * Unique identifier for a session in the AI coding environment.
 * Sessions represent individual conversations with the AI.
 */
export type SessionID = string & { readonly __brand: "SessionID" }

/**
 * Unique identifier for a project/repository.
 * Projects are typically derived from git repository root commits.
 */
export type ProjectID = string & { readonly __brand: "ProjectID" }

/**
 * Unique identifier for an LLM provider (e.g., "anthropic", "openai", "google-vertex").
 */
export type ProviderID = string & { readonly __brand: "ProviderID" }

/**
 * Unique identifier for a specific model within a provider (e.g., "claude-sonnet-4", "gpt-5").
 */
export type ModelID = string & { readonly __brand: "ModelID" }

/**
 * Unique identifier for a message within a session.
 * Messages are the individual turns in the AI conversation.
 */
export type MessageID = string & { readonly __brand: "MessageID" }

/**
 * Unique identifier for a message part (e.g., tool call, text block, etc.).
 * Parts are components within a message.
 */
export type PartID = string & { readonly __brand: "PartID" }

/**
 * Unique identifier for an agent configuration.
 * Agents are specialized AI personas with specific permissions and tools.
 */
export type AgentID = string & { readonly __brand: "AgentID" }

/**
 * File path relative to project root.
 * Distinguished from absolute paths for semantic clarity.
 */
export type RelativeFilePath = string & { readonly __brand: "RelativeFilePath" }

/**
 * Absolute file system path.
 * Distinguished from relative paths for semantic clarity.
 */
export type AbsoluteFilePath = string & { readonly __brand: "AbsoluteFilePath" }

/**
 * Helper to create branded types from plain strings.
 * Use with caution - only when you're certain the string is valid.
 */
export const Branded = {
  sessionID: (id: string): SessionID => id as SessionID,
  projectID: (id: string): ProjectID => id as ProjectID,
  providerID: (id: string): ProviderID => id as ProviderID,
  modelID: (id: string): ModelID => id as ModelID,
  messageID: (id: string): MessageID => id as MessageID,
  partID: (id: string): PartID => id as PartID,
  agentID: (id: string): AgentID => id as AgentID,
  relativeFilePath: (path: string): RelativeFilePath => path as RelativeFilePath,
  absoluteFilePath: (path: string): AbsoluteFilePath => path as AbsoluteFilePath,
} as const

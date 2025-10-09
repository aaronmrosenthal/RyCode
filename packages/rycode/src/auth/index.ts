import path from "path"
import { Global } from "../global"
import fs from "fs/promises"
import z from "zod/v4"
import { SecureStorage } from "../storage/secure-storage"
import { Integrity } from "../storage/integrity"
import { Log } from "../util/log"

export namespace Auth {
  const log = Log.create({ service: "auth" })
  export const Oauth = z
    .object({
      type: z.literal("oauth"),
      refresh: z.string(),
      access: z.string(),
      expires: z.number(),
    })
    .meta({ ref: "OAuth" })

  export const Api = z
    .object({
      type: z.literal("api"),
      key: z.string(),
    })
    .meta({ ref: "ApiAuth" })

  export const WellKnown = z
    .object({
      type: z.literal("wellknown"),
      key: z.string(),
      token: z.string(),
    })
    .meta({ ref: "WellKnownAuth" })

  export const Info = z.discriminatedUnion("type", [Oauth, Api, WellKnown]).meta({ ref: "Auth" })
  export type Info = z.infer<typeof Info>

  const filepath = path.join(Global.Path.data, "auth.json")

  /**
   * Reads and decrypts auth data from storage.
   *
   * Supports both encrypted (with RYCODE_ENCRYPTION_KEY) and plaintext
   * formats for backward compatibility.
   *
   * @returns Decrypted auth data object
   */
  async function readAuthData(): Promise<Record<string, Info>> {
    const file = Bun.file(filepath)

    try {
      const rawData = await file.text()

      // Handle empty file
      if (!rawData || rawData.trim() === "") {
        return {}
      }

      // Check if data has integrity wrapper
      let data: string = rawData
      if (Integrity.hasIntegrity(rawData)) {
        try {
          data = Integrity.unwrap(rawData)
        } catch (error: any) {
          log.warn("integrity check failed on auth data", { error: error.message })
          // Continue with wrapped data - will fail on decrypt/parse
        }
      }

      // Try to decrypt if encrypted
      if (SecureStorage.isEncrypted(data)) {
        data = await SecureStorage.decrypt(data)
      }

      // Parse JSON
      return JSON.parse(data)
    } catch (error: any) {
      log.debug("failed to read auth data, returning empty", { error: error.message })
      return {}
    }
  }

  /**
   * Encrypts and writes auth data to storage.
   *
   * Uses encryption (if RYCODE_ENCRYPTION_KEY set) and integrity verification.
   *
   * @param data - Auth data to write
   */
  async function writeAuthData(data: Record<string, Info>): Promise<void> {
    const file = Bun.file(filepath)

    // Serialize to JSON
    let content = JSON.stringify(data, null, 2)

    // Encrypt if encryption key available
    content = await SecureStorage.encrypt(content)

    // Add integrity wrapper
    content = Integrity.wrap(content)

    // Write atomically
    await Bun.write(file, content)

    // Set restrictive permissions (owner read/write only)
    await fs.chmod(file.name!, 0o600)

    log.debug("auth data written with encryption and integrity", {
      encrypted: SecureStorage.isEncrypted(content),
      integrity: Integrity.hasIntegrity(content),
    })
  }

  /**
   * Retrieves authentication info for a provider.
   *
   * @param providerID - Provider identifier (e.g., "anthropic", "openai")
   * @returns Auth info or undefined if not found
   */
  export async function get(providerID: string): Promise<Info | undefined> {
    const data = await readAuthData()
    return data[providerID]
  }

  /**
   * Retrieves all stored authentication credentials.
   *
   * @returns Record of provider ID to auth info
   */
  export async function all(): Promise<Record<string, Info>> {
    return readAuthData()
  }

  /**
   * Stores authentication info for a provider.
   *
   * Data is encrypted if RYCODE_ENCRYPTION_KEY environment variable is set.
   *
   * @param key - Provider identifier
   * @param info - Authentication info to store
   * @throws Error if key is empty or info is invalid
   */
  export async function set(key: string, info: Info): Promise<void> {
    if (!key || typeof key !== "string" || key.trim() === "") {
      throw new Error("Provider key must be a non-empty string")
    }

    // Validate info against schema
    const validatedInfo = Info.parse(info)

    const data = await readAuthData()
    data[key] = validatedInfo
    await writeAuthData(data)
  }

  /**
   * Removes authentication info for a provider.
   *
   * @param key - Provider identifier
   * @returns true if credential was removed, false if it didn't exist
   */
  export async function remove(key: string): Promise<boolean> {
    if (!key || typeof key !== "string") {
      throw new Error("Provider key must be a non-empty string")
    }

    const data = await readAuthData()
    const existed = key in data
    delete data[key]
    await writeAuthData(data)

    if (existed) {
      log.info("removed auth credential", { provider: key })
    }

    return existed
  }

  /**
   * Migrates plaintext auth data to encrypted format.
   *
   * Call this after setting RYCODE_ENCRYPTION_KEY to encrypt existing data.
   *
   * @returns Number of credentials migrated
   */
  export async function migrateToEncrypted(): Promise<number> {
    if (!process.env['RYCODE_ENCRYPTION_KEY']) {
      throw new Error("RYCODE_ENCRYPTION_KEY must be set to migrate to encrypted storage")
    }

    const data = await readAuthData()
    const count = Object.keys(data).length

    if (count === 0) {
      return 0
    }

    // Re-write with encryption
    await writeAuthData(data)

    log.info("migrated auth credentials to encrypted storage", { count })
    return count
  }
}

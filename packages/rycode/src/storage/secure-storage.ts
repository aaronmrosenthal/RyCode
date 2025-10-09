import crypto from "crypto"
import { Log } from "../util/log"

/**
 * Secure storage module providing encryption at rest for sensitive data.
 *
 * Uses AES-256-GCM with authenticated encryption to prevent tampering.
 * Keys are derived from environment variable using PBKDF2.
 */
export namespace SecureStorage {
  const log = Log.create({ service: "secure-storage" })

  const ALGORITHM = "aes-256-gcm"
  const KEY_LENGTH = 32 // 256 bits
  const IV_LENGTH = 16 // 128 bits (recommended for GCM mode)
  const SALT_LENGTH = 32 // 256 bits for key derivation
  const AUTH_TAG_LENGTH = 16 // 128 bits for GCM authentication
  const PBKDF2_ITERATIONS = 100_000 // OWASP recommendation (2023)

  /**
   * Derives encryption key from master password using PBKDF2.
   *
   * @param masterKey - Master encryption key from environment variable
   * @param salt - Salt for key derivation (unique per installation)
   * @returns Derived 256-bit encryption key
   */
  async function deriveKey(masterKey: string, salt: Buffer): Promise<Buffer> {
    return new Promise((resolve, reject) => {
      crypto.pbkdf2(masterKey, salt, PBKDF2_ITERATIONS, KEY_LENGTH, "sha256", (err, derivedKey) => {
        if (err) reject(err)
        else resolve(derivedKey)
      })
    })
  }

  /**
   * Encrypts data using AES-256-GCM authenticated encryption.
   *
   * Format: salt:iv:authTag:encryptedData (all hex-encoded)
   *
   * @param data - Plain text data to encrypt
   * @param masterKey - Master encryption key (from RYCODE_ENCRYPTION_KEY env var)
   * @returns Encrypted data with salt, IV, and auth tag
   *
   * @example
   * ```typescript
   * const encrypted = await SecureStorage.encrypt(
   *   JSON.stringify({ apiKey: "secret" }),
   *   process.env.RYCODE_ENCRYPTION_KEY!
   * )
   * ```
   */
  export async function encrypt(data: string, masterKey?: string): Promise<string> {
    if (!masterKey) {
      masterKey = process.env['RYCODE_ENCRYPTION_KEY']
    }

    // If no encryption key configured, return plaintext with warning
    if (!masterKey) {
      log.warn("RYCODE_ENCRYPTION_KEY not set - data stored unencrypted")
      return `plaintext:${data}`
    }

    try {
      // Generate random salt for key derivation
      const salt = crypto.randomBytes(SALT_LENGTH)
      const key = await deriveKey(masterKey, salt)

      // Generate random IV
      const iv = crypto.randomBytes(IV_LENGTH)

      // Create cipher and encrypt
      const cipher = crypto.createCipheriv(ALGORITHM, key, iv)
      const encrypted = Buffer.concat([cipher.update(data, "utf8"), cipher.final()])

      // Get authentication tag
      const authTag = cipher.getAuthTag()

      // Format: salt:iv:authTag:encryptedData
      return [salt, iv, authTag, encrypted].map((buf) => buf.toString("hex")).join(":")
    } catch (error: any) {
      log.error("encryption failed", { error: error.message })
      throw new Error("Failed to encrypt data", { cause: error })
    }
  }

  /**
   * Decrypts data encrypted with AES-256-GCM.
   *
   * @param encrypted - Encrypted data in format: salt:iv:authTag:encryptedData
   * @param masterKey - Master encryption key (must match encryption key)
   * @returns Decrypted plain text data
   * @throws Error if decryption fails or auth tag verification fails
   *
   * @example
   * ```typescript
   * const decrypted = await SecureStorage.decrypt(
   *   encryptedData,
   *   process.env.RYCODE_ENCRYPTION_KEY!
   * )
   * const data = JSON.parse(decrypted)
   * ```
   */
  export async function decrypt(encrypted: string, masterKey?: string): Promise<string> {
    // Handle plaintext fallback for backward compatibility
    const PLAINTEXT_PREFIX = "plaintext:"
    if (encrypted.startsWith(PLAINTEXT_PREFIX)) {
      log.warn("reading unencrypted data - consider re-encrypting")
      return encrypted.substring(PLAINTEXT_PREFIX.length)
    }

    if (!masterKey) {
      masterKey = process.env['RYCODE_ENCRYPTION_KEY']
    }

    if (!masterKey) {
      throw new Error("RYCODE_ENCRYPTION_KEY required for decryption")
    }

    try {
      // Parse encrypted data format: salt:iv:authTag:encryptedData
      const parts = encrypted.split(":")
      const EXPECTED_PARTS = 4
      if (parts.length !== EXPECTED_PARTS) {
        throw new Error(`Invalid encrypted data format: expected ${EXPECTED_PARTS} parts, got ${parts.length}`)
      }

      const [saltHex, ivHex, authTagHex, dataHex] = parts
      const [salt, iv, authTag, data] = [saltHex, ivHex, authTagHex, dataHex].map((hex) =>
        Buffer.from(hex, "hex"),
      )

      // Derive key using same salt
      const key = await deriveKey(masterKey, salt)

      // Create decipher
      const decipher = crypto.createDecipheriv(ALGORITHM, key, iv)
      decipher.setAuthTag(authTag)

      // Decrypt and verify
      const decrypted = Buffer.concat([decipher.update(data), decipher.final()])

      return decrypted.toString("utf8")
    } catch (error: any) {
      log.error("decryption failed", { error: error.message })
      throw new Error("Failed to decrypt data - wrong key or corrupted data", { cause: error })
    }
  }

  /**
   * Checks if data is encrypted or plaintext.
   *
   * @param data - Data to check
   * @returns true if data appears to be encrypted
   */
  export function isEncrypted(data: string): boolean {
    // Plaintext marker
    const PLAINTEXT_PREFIX = "plaintext:"
    if (data.startsWith(PLAINTEXT_PREFIX)) return false

    // Encrypted format: 4 hex strings separated by colons
    const parts = data.split(":")
    const EXPECTED_PARTS = 4
    if (parts.length !== EXPECTED_PARTS) return false

    // Validate minimum lengths for encrypted components
    const [salt, iv, authTag, encryptedData] = parts
    if (
      salt.length !== SALT_LENGTH * 2 || // hex encoding doubles length
      iv.length !== IV_LENGTH * 2 ||
      authTag.length !== AUTH_TAG_LENGTH * 2 ||
      encryptedData.length < 2 // At least 1 byte of data
    ) {
      return false
    }

    // Check if all parts are valid hex
    return parts.every((part) => /^[0-9a-fA-F]+$/.test(part))
  }

  /**
   * Re-encrypts plaintext data with encryption enabled.
   *
   * Useful for migration from unencrypted to encrypted storage.
   *
   * @param data - Data (encrypted or plaintext)
   * @param masterKey - Master encryption key
   * @returns Re-encrypted data
   */
  export async function reencrypt(data: string, masterKey?: string): Promise<string> {
    // If already encrypted with current format, return as-is
    if (isEncrypted(data) && !data.startsWith("plaintext:")) {
      return data
    }

    // Decrypt if needed, then re-encrypt
    const plaintext = await decrypt(data, masterKey)
    return encrypt(plaintext, masterKey)
  }

  /**
   * Generates a secure random encryption key suitable for RYCODE_ENCRYPTION_KEY.
   *
   * @returns Base64-encoded 256-bit random key
   *
   * @example
   * ```typescript
   * const key = SecureStorage.generateKey()
   * // Set in environment: export RYCODE_ENCRYPTION_KEY="<key>"
   * console.log("Set this key:", key)
   * ```
   */
  export function generateKey(): string {
    const key = crypto.randomBytes(KEY_LENGTH)
    return key.toString("base64")
  }

  /**
   * Validates encryption key format.
   *
   * @param key - Key to validate
   * @returns true if key is valid format
   */
  export function isValidKey(key: string): boolean {
    if (!key || typeof key !== "string") return false

    try {
      const decoded = Buffer.from(key, "base64")
      return decoded.length >= KEY_LENGTH
    } catch {
      return false
    }
  }

  /**
   * Securely wipes sensitive data from memory.
   * Overwrites buffer with zeros before releasing.
   *
   * @param buffer - Buffer to wipe
   */
  export function secureWipe(buffer: Buffer): void {
    if (buffer && Buffer.isBuffer(buffer)) {
      buffer.fill(0)
    }
  }
}

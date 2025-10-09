import crypto from "crypto"
import { Log } from "../util/log"

/**
 * File integrity verification using SHA-256 checksums.
 *
 * Provides cryptographic verification that storage files haven't been
 * tampered with between write and read operations.
 */
export namespace Integrity {
  const log = Log.create({ service: "integrity" })

  const ALGORITHM = "sha256"
  const ENCODING = "hex"

  /**
   * Computes SHA-256 checksum of data.
   *
   * @param data - Data to checksum
   * @returns Hex-encoded SHA-256 hash (64 characters)
   * @throws Error if data is not a string
   *
   * @example
   * ```typescript
   * const checksum = Integrity.computeChecksum(JSON.stringify(data))
   * // Returns: "a3c2f1b8..." (64 hex characters)
   * ```
   */
  export function computeChecksum(data: string): string {
    if (typeof data !== "string") {
      throw new Error("Data must be a string")
    }

    const hash = crypto.createHash(ALGORITHM)
    hash.update(data, "utf8")
    return hash.digest(ENCODING)
  }

  /**
   * Verifies data integrity against a known checksum.
   *
   * Uses constant-time comparison to prevent timing attacks.
   *
   * @param data - Data to verify
   * @param expectedChecksum - Expected SHA-256 checksum
   * @returns true if checksums match, false otherwise
   *
   * @example
   * ```typescript
   * if (!Integrity.verifyChecksum(data, storedChecksum)) {
   *   throw new IntegrityError("Data has been tampered with")
   * }
   * ```
   */
  export function verifyChecksum(data: string, expectedChecksum: string): boolean {
    const actualChecksum = computeChecksum(data)

    try {
      // Use constant-time comparison to prevent timing attacks
      const actualBuf = Buffer.from(actualChecksum, ENCODING)
      const expectedBuf = Buffer.from(expectedChecksum, ENCODING)

      // Ensure buffers are same length
      if (actualBuf.length !== expectedBuf.length) {
        return false
      }

      return crypto.timingSafeEqual(actualBuf, expectedBuf)
    } catch (error: any) {
      log.error("checksum verification failed", { error: error.message })
      return false
    }
  }

  /**
   * Wraps data with integrity checksum.
   *
   * Format: checksum:data (checksum is 64 hex chars)
   *
   * @param data - Data to wrap
   * @returns Data prefixed with checksum
   */
  export function wrap(data: string): string {
    const checksum = computeChecksum(data)
    return `${checksum}:${data}`
  }

  /**
   * Unwraps and verifies data with integrity checksum.
   *
   * @param wrapped - Data wrapped with checksum
   * @returns Unwrapped data if checksum valid
   * @throws IntegrityError if checksum verification fails
   */
  export function unwrap(wrapped: string): string {
    // Parse checksum:data format
    const CHECKSUM_LENGTH = 64
    const colonIndex = wrapped.indexOf(":")

    if (colonIndex === -1 || colonIndex !== CHECKSUM_LENGTH) {
      throw new IntegrityError({
        message: "Invalid integrity format - missing or malformed checksum",
      })
    }

    const expectedChecksum = wrapped.substring(0, CHECKSUM_LENGTH)
    const data = wrapped.substring(CHECKSUM_LENGTH + 1) // Skip colon

    // Verify checksum
    if (!verifyChecksum(data, expectedChecksum)) {
      log.warn("integrity check failed", {
        expectedChecksum,
        actualChecksum: computeChecksum(data),
      })
      throw new IntegrityError({
        message: "Data integrity check failed - data may have been tampered with",
      })
    }

    return data
  }

  /**
   * Checks if data has integrity wrapper.
   *
   * @param data - Data to check
   * @returns true if data appears to have integrity checksum
   */
  export function hasIntegrity(data: string): boolean {
    const CHECKSUM_LENGTH = 64
    const MIN_LENGTH = CHECKSUM_LENGTH + 1 // checksum + colon

    // Check for checksum:data format (64 hex chars + colon)
    if (data.length < MIN_LENGTH) return false
    if (data.charAt(CHECKSUM_LENGTH) !== ":") return false
    return /^[0-9a-fA-F]{64}:/.test(data)
  }

  /**
   * Integrity verification error.
   *
   * Thrown when data fails checksum verification, indicating tampering or corruption.
   *
   * @example
   * ```typescript
   * try {
   *   const data = Integrity.unwrap(wrapped)
   * } catch (error) {
   *   if (error instanceof Integrity.IntegrityError) {
   *     console.error("Data has been tampered with!")
   *   }
   * }
   * ```
   */
  export class IntegrityError extends Error {
    constructor(public readonly details: { message: string }) {
      super(details.message)
      this.name = "IntegrityError"

      // Maintain proper stack trace (V8 engines)
      if (Error.captureStackTrace) {
        Error.captureStackTrace(this, IntegrityError)
      }
    }
  }

  /**
   * Generates integrity metadata for a file.
   *
   * Includes checksum, file size, and timestamp for comprehensive verification.
   *
   * @param data - File data
   * @returns Integrity metadata
   */
  export function generateMetadata(data: string): IntegrityMetadata {
    return {
      checksum: computeChecksum(data),
      size: Buffer.byteLength(data, "utf8"),
      timestamp: Date.now(),
      algorithm: ALGORITHM,
    }
  }

  /**
   * Verifies file integrity using metadata.
   *
   * @param data - File data
   * @param metadata - Previously generated metadata
   * @returns true if all checks pass
   */
  export function verifyMetadata(data: string, metadata: IntegrityMetadata): boolean {
    // Check size
    const actualSize = Buffer.byteLength(data, "utf8")
    if (actualSize !== metadata.size) {
      log.warn("size mismatch", {
        expected: metadata.size,
        actual: actualSize,
      })
      return false
    }

    // Check algorithm
    if (metadata.algorithm !== ALGORITHM) {
      log.warn("algorithm mismatch", {
        expected: ALGORITHM,
        actual: metadata.algorithm,
      })
      return false
    }

    // Check checksum
    return verifyChecksum(data, metadata.checksum)
  }

  export interface IntegrityMetadata {
    /** SHA-256 checksum (hex) */
    checksum: string
    /** File size in bytes */
    size: number
    /** Timestamp when metadata was generated */
    timestamp: number
    /** Hash algorithm used */
    algorithm: string
  }
}

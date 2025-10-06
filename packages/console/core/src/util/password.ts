import bcrypt from "bcrypt"
import { z } from "zod"

/**
 * Password utilities for secure hashing and verification
 * Uses bcrypt with recommended salt rounds
 */
export namespace Password {
  /**
   * Number of salt rounds for bcrypt
   * 12 rounds provides good security/performance balance
   * Takes ~250ms on modern hardware
   */
  const SALT_ROUNDS = 12

  /**
   * Password requirements schema
   */
  export const schema = z
    .string()
    .min(8, "Password must be at least 8 characters")
    .max(128, "Password must not exceed 128 characters")
    .refine(
      (password) => /[A-Z]/.test(password),
      "Password must contain at least one uppercase letter",
    )
    .refine(
      (password) => /[a-z]/.test(password),
      "Password must contain at least one lowercase letter",
    )
    .refine((password) => /[0-9]/.test(password), "Password must contain at least one number")
    .refine(
      (password) => /[!@#$%^&*(),.?":{}|<>]/.test(password),
      "Password must contain at least one special character",
    )

  /**
   * Hash a password using bcrypt
   * @param password - Plain text password to hash
   * @returns Promise resolving to bcrypt hash string
   */
  export async function hash(password: string): Promise<string> {
    // Validate password meets requirements
    schema.parse(password)

    // Generate salt and hash
    const hash = await bcrypt.hash(password, SALT_ROUNDS)

    return hash
  }

  /**
   * Verify a password against a bcrypt hash
   * @param password - Plain text password to verify
   * @param hash - Bcrypt hash to compare against
   * @returns Promise resolving to true if password matches
   */
  export async function verify(password: string, hash: string): Promise<boolean> {
    try {
      return await bcrypt.compare(password, hash)
    } catch (error) {
      // Invalid hash format or other bcrypt error
      console.error("Password verification error:", error)
      return false
    }
  }

  /**
   * Check if a hash needs to be rehashed (e.g., salt rounds changed)
   * @param hash - Bcrypt hash to check
   * @returns Promise resolving to true if rehash needed
   */
  export async function needsRehash(hash: string): Promise<boolean> {
    try {
      // Extract rounds from hash (bcrypt format: $2b$rounds$...)
      const rounds = parseInt(hash.split("$")[2] || "0", 10)
      return rounds !== SALT_ROUNDS
    } catch (error) {
      // Invalid hash format
      return true
    }
  }

  /**
   * Validate password meets strength requirements without hashing
   * @param password - Password to validate
   * @returns Object with validation result and errors
   */
  export function validate(password: string): {
    valid: boolean
    errors: string[]
  } {
    const errors: string[] = []

    if (password.length < 8) {
      errors.push("Password must be at least 8 characters")
    }

    if (password.length > 128) {
      errors.push("Password must not exceed 128 characters")
    }

    if (!/[A-Z]/.test(password)) {
      errors.push("Password must contain at least one uppercase letter")
    }

    if (!/[a-z]/.test(password)) {
      errors.push("Password must contain at least one lowercase letter")
    }

    if (!/[0-9]/.test(password)) {
      errors.push("Password must contain at least one number")
    }

    if (!/[!@#$%^&*(),.?":{}|<>]/.test(password)) {
      errors.push("Password must contain at least one special character")
    }

    return {
      valid: errors.length === 0,
      errors,
    }
  }

  /**
   * Calculate password strength score (0-4)
   * 0: Very weak
   * 1: Weak
   * 2: Fair
   * 3: Strong
   * 4: Very strong
   */
  export function strength(password: string): {
    score: number
    feedback: string
  } {
    let score = 0
    const feedback: string[] = []

    // Length bonus
    if (password.length >= 8) score++
    if (password.length >= 12) score++
    if (password.length >= 16) score++

    // Character variety
    if (/[A-Z]/.test(password) && /[a-z]/.test(password)) {
      score++
      feedback.push("Good mix of cases")
    }

    if (/[0-9]/.test(password)) {
      score++
      feedback.push("Contains numbers")
    }

    if (/[!@#$%^&*(),.?":{}|<>]/.test(password)) {
      score++
      feedback.push("Contains special characters")
    }

    // Common patterns penalty
    if (/(.)\1{2,}/.test(password)) {
      score--
      feedback.push("Avoid repeating characters")
    }

    if (/^[0-9]+$/.test(password)) {
      score -= 2
      feedback.push("Don't use only numbers")
    }

    if (/(password|12345|qwerty)/i.test(password)) {
      score -= 2
      feedback.push("Avoid common words/patterns")
    }

    // Normalize score to 0-4
    score = Math.max(0, Math.min(4, score))

    const strengthLabels = ["Very weak", "Weak", "Fair", "Strong", "Very strong"]

    return {
      score,
      feedback: `${strengthLabels[score]}${feedback.length > 0 ? ": " + feedback.join(", ") : ""}`,
    }
  }
}

import { z } from "zod"
import { eq, and, sql } from "drizzle-orm"
import { fn } from "./util/fn"
import { Database } from "./drizzle"
import { AccountTable } from "./schema/account.sql"
import { Identifier } from "./identifier"
import { render } from "@jsx-email/render"
import { AWS } from "./aws"
import { Password } from "./util/password"
import crypto from "crypto"

export namespace PasswordReset {
  // Store password reset tokens in-memory for now
  // WARNING: In production, use Redis or database table for:
  // - Persistence across server restarts
  // - Distributed systems support
  // - Better memory management
  const resetTokens = new Map<
    string,
    {
      accountID: string
      email: string
      expiresAt: number
      attempts: number // Track verification attempts to prevent brute force
    }
  >()

  const MAX_RESET_ATTEMPTS = 5
  const RATE_LIMIT_WINDOW = 15 * 60 * 1000 // 15 minutes
  const rateLimitMap = new Map<string, { count: number; resetAt: number }>()

  /**
   * Request a password reset email
   */
  export const requestReset = fn(
    z.object({
      email: z.string().email(),
    }),
    async ({ email }) => {
      // Rate limiting: Prevent spam/abuse
      const normalizedEmail = email.toLowerCase().trim()
      const rateLimit = rateLimitMap.get(normalizedEmail)
      const now = Date.now()

      if (rateLimit) {
        if (now < rateLimit.resetAt) {
          if (rateLimit.count >= 3) {
            // Too many requests - but still return success to prevent enumeration
            console.warn(`Rate limit exceeded for password reset: ${normalizedEmail}`)
            return { success: true }
          }
          rateLimit.count++
        } else {
          // Reset window expired
          rateLimitMap.set(normalizedEmail, { count: 1, resetAt: now + RATE_LIMIT_WINDOW })
        }
      } else {
        rateLimitMap.set(normalizedEmail, { count: 1, resetAt: now + RATE_LIMIT_WINDOW })
      }

      // Find account by email (case-insensitive)
      const account = await Database.use((tx) =>
        tx
          .select()
          .from(AccountTable)
          .where(sql`LOWER(${AccountTable.email}) = ${normalizedEmail}`)
          .then((rows) => rows[0]),
      )

      // Always return success to prevent email enumeration attacks
      if (!account) {
        console.log(`Password reset requested for non-existent email: ${normalizedEmail}`)
        return { success: true }
      }

      // Generate secure reset token
      const token = crypto.randomBytes(32).toString("hex")
      const expiresAt = Date.now() + 60 * 60 * 1000 // 1 hour

      // Store token with attempt tracking
      resetTokens.set(token, {
        accountID: account.id,
        email: account.email,
        expiresAt,
        attempts: 0,
      })

      // Clean up expired tokens periodically
      cleanupExpiredTokens()

      // Send reset email with proper HTML escaping
      const resetUrl = `${process.env.AUTH_FRONTEND_URL}/auth/reset-password?token=${encodeURIComponent(token)}`

      try {
        await AWS.sendEmail({
          to: account.email, // Use account email, not user input
          subject: "Reset your OpenCode password",
          body: `
            <h1>Reset your password</h1>
            <p>You requested to reset your password. Click the link below to continue:</p>
            <p><a href="${resetUrl}">Reset Password</a></p>
            <p>This link will expire in 1 hour.</p>
            <p>If you didn't request this, you can safely ignore this email.</p>
            <p><small>For security, this link can only be used once.</small></p>
          `,
        })
      } catch (e) {
        console.error("Failed to send password reset email:", e)
        // Don't expose email sending errors to prevent enumeration
        return { success: true }
      }

      return { success: true }
    },
  )

  /**
   * Verify a reset token is valid
   */
  export const verifyToken = fn(z.string(), (token) => {
    const resetData = resetTokens.get(token)

    if (!resetData) {
      return { valid: false, reason: "Token not found" }
    }

    if (resetData.expiresAt < Date.now()) {
      resetTokens.delete(token)
      return { valid: false, reason: "Token expired" }
    }

    // Check for brute force attempts
    if (resetData.attempts >= MAX_RESET_ATTEMPTS) {
      resetTokens.delete(token)
      console.warn(`Max verification attempts exceeded for token`)
      return { valid: false, reason: "Token locked due to too many attempts" }
    }

    // Increment attempt counter
    resetData.attempts++

    return {
      valid: true,
      accountID: resetData.accountID,
      email: resetData.email,
    }
  })

  /**
   * Reset password using a valid token
   */
  export const resetPassword = fn(
    z.object({
      token: z.string(),
      newPassword: Password.schema, // Uses comprehensive password validation
    }),
    async ({ token, newPassword }) => {
      const verification = verifyToken(token)

      if (!verification.valid) {
        throw new Error(verification.reason || "Invalid token")
      }

      // Hash the password using bcrypt (12 rounds)
      const passwordHash = await Password.hash(newPassword)

      // Update account with hashed password
      await Database.use((tx) =>
        tx
          .update(AccountTable)
          .set({
            passwordHash,
            timeUpdated: sql`now()`,
          })
          .where(eq(AccountTable.id, verification.accountID!)),
      )

      // Invalidate the token immediately after successful use
      resetTokens.delete(token)

      // Clear rate limit for this email
      const resetData = resetTokens.get(token)
      if (resetData) {
        rateLimitMap.delete(resetData.email.toLowerCase().trim())
      }

      return { success: true }
    },
  )

  /**
   * Clean up expired tokens
   */
  function cleanupExpiredTokens() {
    const now = Date.now()
    for (const [token, data] of resetTokens.entries()) {
      if (data.expiresAt < now) {
        resetTokens.delete(token)
      }
    }
  }

  /**
   * For testing: clear all tokens
   */
  export const clearAllTokens = fn(z.void(), () => {
    resetTokens.clear()
  })
}

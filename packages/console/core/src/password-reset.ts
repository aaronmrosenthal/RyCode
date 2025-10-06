import { z } from "zod"
import { eq, and, sql } from "drizzle-orm"
import { fn } from "./util/fn"
import { Database } from "./drizzle"
import { AccountTable } from "./schema/account.sql"
import { Identifier } from "./identifier"
import { render } from "@jsx-email/render"
import { AWS } from "./aws"
import crypto from "crypto"

export namespace PasswordReset {
  // Store password reset tokens in-memory for now
  // In production, you'd want to use Redis or a database table
  const resetTokens = new Map<
    string,
    {
      accountID: string
      email: string
      expiresAt: number
    }
  >()

  /**
   * Request a password reset email
   */
  export const requestReset = fn(
    z.object({
      email: z.string().email(),
    }),
    async ({ email }) => {
      // Find account by email
      const account = await Database.use((tx) =>
        tx
          .select()
          .from(AccountTable)
          .where(eq(AccountTable.email, email))
          .then((rows) => rows[0]),
      )

      // Always return success to prevent email enumeration attacks
      if (!account) {
        console.log(`Password reset requested for non-existent email: ${email}`)
        return { success: true }
      }

      // Generate secure reset token
      const token = crypto.randomBytes(32).toString("hex")
      const expiresAt = Date.now() + 60 * 60 * 1000 // 1 hour

      // Store token
      resetTokens.set(token, {
        accountID: account.id,
        email: account.email,
        expiresAt,
      })

      // Clean up expired tokens periodically
      cleanupExpiredTokens()

      // Send reset email
      try {
        await AWS.sendEmail({
          to: email,
          subject: "Reset your OpenCode password",
          body: `
            <h1>Reset your password</h1>
            <p>You requested to reset your password. Click the link below to continue:</p>
            <p><a href="${process.env.AUTH_FRONTEND_URL}/auth/reset-password?token=${token}">Reset Password</a></p>
            <p>This link will expire in 1 hour.</p>
            <p>If you didn't request this, you can safely ignore this email.</p>
          `,
        })
      } catch (e) {
        console.error("Failed to send password reset email:", e)
        throw new Error("Failed to send password reset email")
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
      newPassword: z.string().min(8),
    }),
    async ({ token, newPassword }) => {
      const verification = verifyToken(token)

      if (!verification.valid) {
        throw new Error(verification.reason || "Invalid token")
      }

      // In a real implementation, you would hash the password here
      // For now, this is a placeholder since the auth system uses OAuth
      // and doesn't store passwords directly

      // Update account (placeholder - adjust based on your schema)
      await Database.use((tx) =>
        tx
          .update(AccountTable)
          .set({
            // Add password field to schema if needed
            timeUpdated: sql`now()`,
          })
          .where(eq(AccountTable.id, verification.accountID!)),
      )

      // Invalidate the token
      resetTokens.delete(token)

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

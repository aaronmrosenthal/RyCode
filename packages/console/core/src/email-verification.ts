import { z } from "zod"
import { eq, and, sql } from "drizzle-orm"
import { fn } from "./util/fn"
import { Database } from "./drizzle"
import { AccountTable } from "./schema/account.sql"
import { Identifier } from "./identifier"
import { render } from "@jsx-email/render"
import { AWS } from "./aws"
import crypto from "crypto"

export namespace EmailVerification {
  // Store verification tokens in-memory for now
  // In production, you'd want to use Redis or a database table
  const verificationTokens = new Map<
    string,
    {
      accountID: string
      email: string
      expiresAt: number
    }
  >()

  /**
   * Send verification email to a new account
   */
  export const sendVerification = fn(
    z.object({
      accountID: z.string(),
      email: z.string().email(),
    }),
    async ({ accountID, email }) => {
      // Generate secure verification token
      const token = crypto.randomBytes(32).toString("hex")
      const expiresAt = Date.now() + 24 * 60 * 60 * 1000 // 24 hours

      // Store token
      verificationTokens.set(token, {
        accountID,
        email,
        expiresAt,
      })

      // Clean up expired tokens periodically
      cleanupExpiredTokens()

      // Send verification email
      try {
        await AWS.sendEmail({
          to: email,
          subject: "Verify your OpenCode email address",
          body: `
            <h1>Welcome to OpenCode!</h1>
            <p>Please verify your email address by clicking the link below:</p>
            <p><a href="${process.env.AUTH_FRONTEND_URL}/auth/verify-email?token=${token}">Verify Email</a></p>
            <p>This link will expire in 24 hours.</p>
            <p>If you didn't create an account, you can safely ignore this email.</p>
          `,
        })
      } catch (e) {
        console.error("Failed to send verification email:", e)
        throw new Error("Failed to send verification email")
      }

      return { success: true, token }
    },
  )

  /**
   * Verify an email using the token
   */
  export const verifyEmail = fn(z.string(), async (token) => {
    const verificationData = verificationTokens.get(token)

    if (!verificationData) {
      throw new Error("Invalid verification token")
    }

    if (verificationData.expiresAt < Date.now()) {
      verificationTokens.delete(token)
      throw new Error("Verification token expired")
    }

    // Mark email as verified in the database
    await Database.use((tx) =>
      tx
        .update(AccountTable)
        .set({
          // Add emailVerified field to schema if needed
          // emailVerified: true,
          timeUpdated: sql`now()`,
        })
        .where(eq(AccountTable.id, verificationData.accountID)),
    )

    // Invalidate the token
    verificationTokens.delete(token)

    return {
      success: true,
      accountID: verificationData.accountID,
      email: verificationData.email,
    }
  })

  /**
   * Check if an email is verified
   */
  export const isVerified = fn(z.string(), async (accountID) => {
    const account = await Database.use((tx) =>
      tx
        .select()
        .from(AccountTable)
        .where(eq(AccountTable.id, accountID))
        .then((rows) => rows[0]),
    )

    if (!account) {
      return false
    }

    // Check emailVerified field (add to schema if needed)
    // return account.emailVerified === true
    return true // Placeholder
  })

  /**
   * Resend verification email
   */
  export const resendVerification = fn(
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

      if (!account) {
        // Don't reveal if email exists
        return { success: true }
      }

      // Check if already verified
      const verified = await isVerified(account.id)
      if (verified) {
        return { success: true, alreadyVerified: true }
      }

      // Send new verification email
      return await sendVerification({
        accountID: account.id,
        email: account.email,
      })
    },
  )

  /**
   * Clean up expired tokens
   */
  function cleanupExpiredTokens() {
    const now = Date.now()
    for (const [token, data] of verificationTokens.entries()) {
      if (data.expiresAt < now) {
        verificationTokens.delete(token)
      }
    }
  }

  /**
   * For testing: clear all tokens
   */
  export const clearAllTokens = fn(z.void(), () => {
    verificationTokens.clear()
  })

  /**
   * For testing: get token for account
   */
  export const getTokenForAccount = fn(z.string(), (accountID) => {
    for (const [token, data] of verificationTokens.entries()) {
      if (data.accountID === accountID) {
        return token
      }
    }
    return null
  })
}

import { z } from "zod"
import { eq, and, sql } from "drizzle-orm"
import { fn } from "./util/fn"
import { Database } from "./drizzle"
import { AccountTable } from "./schema/account.sql"
import { Password } from "./util/password"

/**
 * Password-based authentication
 * Complements OAuth authentication for accounts with passwords
 */
export namespace PasswordAuth {
  // Track failed login attempts for rate limiting
  const failedAttempts = new Map<
    string,
    {
      count: number
      lockedUntil: number | null
    }
  >()

  const MAX_FAILED_ATTEMPTS = 5
  const LOCKOUT_DURATION = 15 * 60 * 1000 // 15 minutes
  const ATTEMPT_WINDOW = 15 * 60 * 1000 // 15 minutes

  /**
   * Authenticate user with email and password
   */
  export const authenticate = fn(
    z.object({
      email: z.string().email(),
      password: z.string(),
    }),
    async ({ email, password }) => {
      const normalizedEmail = email.toLowerCase().trim()

      // Check if account is locked
      const attempts = failedAttempts.get(normalizedEmail)
      if (attempts?.lockedUntil && attempts.lockedUntil > Date.now()) {
        const minutesLeft = Math.ceil((attempts.lockedUntil - Date.now()) / 60000)
        throw new Error(
          `Account temporarily locked due to too many failed attempts. Try again in ${minutesLeft} minute${minutesLeft > 1 ? "s" : ""}.`,
        )
      }

      // Find account
      const account = await Database.use((tx) =>
        tx
          .select()
          .from(AccountTable)
          .where(sql`LOWER(${AccountTable.email}) = ${normalizedEmail}`)
          .then((rows) => rows[0]),
      )

      // Always check password even if account doesn't exist (prevent timing attacks)
      const dummyHash = "$2b$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5eDLYJWOvz9/2" // Dummy bcrypt hash

      const passwordHash = account?.passwordHash || dummyHash
      const isValid = await Password.verify(password, passwordHash)

      // Check if account exists and password is correct
      if (!account || !account.passwordHash || !isValid) {
        // Increment failed attempts
        const current = failedAttempts.get(normalizedEmail) || { count: 0, lockedUntil: null }
        current.count++

        if (current.count >= MAX_FAILED_ATTEMPTS) {
          current.lockedUntil = Date.now() + LOCKOUT_DURATION
          console.warn(`Account locked after ${MAX_FAILED_ATTEMPTS} failed attempts: ${normalizedEmail}`)
        }

        failedAttempts.set(normalizedEmail, current)

        // Generic error message (don't reveal if account exists)
        throw new Error("Invalid email or password")
      }

      // Check if email is verified
      if (!account.emailVerified) {
        throw new Error("Please verify your email address before logging in")
      }

      // Clear failed attempts on successful login
      failedAttempts.delete(normalizedEmail)

      // Check if password needs rehashing (e.g., salt rounds changed)
      if (await Password.needsRehash(account.passwordHash)) {
        // Rehash with current settings
        const newHash = await Password.hash(password)
        await Database.use((tx) =>
          tx
            .update(AccountTable)
            .set({
              passwordHash: newHash,
              timeUpdated: sql`now()`,
            })
            .where(eq(AccountTable.id, account.id)),
        )
      }

      return {
        success: true,
        accountID: account.id,
        email: account.email,
      }
    },
  )

  /**
   * Check if an email has a password set (vs OAuth-only)
   */
  export const hasPassword = fn(z.string().email(), async (email) => {
    const normalizedEmail = email.toLowerCase().trim()

    const account = await Database.use((tx) =>
      tx
        .select({ passwordHash: AccountTable.passwordHash })
        .from(AccountTable)
        .where(sql`LOWER(${AccountTable.email}) = ${normalizedEmail}`)
        .then((rows) => rows[0]),
    )

    return account?.passwordHash !== null && account?.passwordHash !== undefined
  })

  /**
   * Set initial password for an account (e.g., during signup)
   */
  export const setPassword = fn(
    z.object({
      accountID: z.string(),
      password: Password.schema,
    }),
    async ({ accountID, password }) => {
      // Check if password already exists
      const account = await Database.use((tx) =>
        tx
          .select({ passwordHash: AccountTable.passwordHash })
          .from(AccountTable)
          .where(eq(AccountTable.id, accountID))
          .then((rows) => rows[0]),
      )

      if (!account) {
        throw new Error("Account not found")
      }

      if (account.passwordHash) {
        throw new Error("Password already set. Use password reset to change it.")
      }

      // Hash and store password
      const passwordHash = await Password.hash(password)

      await Database.use((tx) =>
        tx
          .update(AccountTable)
          .set({
            passwordHash,
            timeUpdated: sql`now()`,
          })
          .where(eq(AccountTable.id, accountID)),
      )

      return { success: true }
    },
  )

  /**
   * Change password (requires current password verification)
   */
  export const changePassword = fn(
    z.object({
      accountID: z.string(),
      currentPassword: z.string(),
      newPassword: Password.schema,
    }),
    async ({ accountID, currentPassword, newPassword }) => {
      // Get account
      const account = await Database.use((tx) =>
        tx
          .select()
          .from(AccountTable)
          .where(eq(AccountTable.id, accountID))
          .then((rows) => rows[0]),
      )

      if (!account) {
        throw new Error("Account not found")
      }

      if (!account.passwordHash) {
        throw new Error("No password set for this account")
      }

      // Verify current password
      const isValid = await Password.verify(currentPassword, account.passwordHash)
      if (!isValid) {
        throw new Error("Current password is incorrect")
      }

      // Prevent reusing the same password
      const isSamePassword = await Password.verify(newPassword, account.passwordHash)
      if (isSamePassword) {
        throw new Error("New password must be different from current password")
      }

      // Hash and update password
      const passwordHash = await Password.hash(newPassword)

      await Database.use((tx) =>
        tx
          .update(AccountTable)
          .set({
            passwordHash,
            timeUpdated: sql`now()`,
          })
          .where(eq(AccountTable.id, accountID)),
      )

      return { success: true }
    },
  )

  /**
   * Clear failed login attempts (for testing or admin use)
   */
  export const clearFailedAttempts = fn(z.string().email(), (email) => {
    const normalizedEmail = email.toLowerCase().trim()
    failedAttempts.delete(normalizedEmail)
    return { success: true }
  })

  /**
   * Get failed attempt status (for security monitoring)
   */
  export const getFailedAttempts = fn(z.string().email(), (email) => {
    const normalizedEmail = email.toLowerCase().trim()
    const attempts = failedAttempts.get(normalizedEmail)

    if (!attempts) {
      return { count: 0, locked: false }
    }

    const locked = attempts.lockedUntil !== null && attempts.lockedUntil > Date.now()
    const minutesLeft = locked ? Math.ceil((attempts.lockedUntil! - Date.now()) / 60000) : 0

    return {
      count: attempts.count,
      locked,
      minutesLeft,
    }
  })
}

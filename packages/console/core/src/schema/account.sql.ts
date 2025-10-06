import { mysqlTable, uniqueIndex, varchar, boolean, index } from "drizzle-orm/mysql-core"
import { id, timestamps } from "../drizzle/types"

export const AccountTable = mysqlTable(
  "account",
  {
    id: id(),
    ...timestamps,
    email: varchar("email", { length: 255 }).notNull(),
    emailVerified: boolean("email_verified").notNull().default(false),
    passwordHash: varchar("password_hash", { length: 255 }),
  },
  (table) => [
    uniqueIndex("email").on(table.email),
    index("idx_account_email_verified").on(table.emailVerified),
    index("idx_account_password_exists").on(table.passwordHash),
  ],
)

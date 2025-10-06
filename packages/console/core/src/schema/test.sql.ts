import { mysqlTable, varchar, text, index } from "drizzle-orm/mysql-core"
import { timestamps, ulid } from "../drizzle/types"
import { workspaceIndexes } from "./workspace.sql"

export const TestTable = mysqlTable(
  "test",
  {
    id: ulid("id").notNull(),
    workspaceID: varchar("workspace_id", { length: 255 }).notNull(),
    name: varchar("name", { length: 255 }).notNull(),
    description: text("description"),
    status: varchar("status", { length: 50 }).notNull().default("active"),
    ...timestamps,
  },
  workspaceIndexes,
)

export const TestStatus = ["active", "inactive", "archived"] as const

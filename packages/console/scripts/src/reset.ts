import { Resource } from "@rycode-ai/console-resource"
import { Database } from "@rycode-ai/console-core/drizzle/index.js"
import { UserTable } from "@rycode-ai/console-core/schema/user.sql.js"
import { AccountTable } from "@rycode-ai/console-core/schema/account.sql.js"
import { WorkspaceTable } from "@rycode-ai/console-core/schema/workspace.sql.js"
import { BillingTable, PaymentTable, UsageTable } from "@rycode-ai/console-core/schema/billing.sql.js"
import { KeyTable } from "@rycode-ai/console-core/schema/key.sql.js"

if (Resource.App.stage !== "frank") throw new Error("This script is only for frank")

for (const table of [AccountTable, BillingTable, KeyTable, PaymentTable, UsageTable, UserTable, WorkspaceTable]) {
  await Database.use((tx) => tx.delete(table))
}

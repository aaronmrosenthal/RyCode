# Database Migrations

## Overview
Created 3 new database migrations to support the newly implemented features:
- Test entity table
- Email verification field
- Password hash field

---

## Migration Files

### 0028_add_test_table.sql
**Purpose:** Create the `test` table for the Test entity CRUD operations

```sql
CREATE TABLE `test` (
	`id` varchar(30) NOT NULL,
	`workspace_id` varchar(30) NOT NULL,
	`time_created` timestamp(3) NOT NULL DEFAULT (now()),
	`time_updated` timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
	`time_deleted` timestamp(3),
	`name` varchar(255) NOT NULL,
	`description` text,
	`status` varchar(50) NOT NULL DEFAULT 'active',
	CONSTRAINT `test_workspace_id_id_pk` PRIMARY KEY(`workspace_id`,`id`)
);

CREATE INDEX `idx_test_workspace_deleted` ON `test` (`workspace_id`, `time_deleted`);
CREATE INDEX `idx_test_status` ON `test` (`status`);
```

**Features:**
- Primary key: Composite (`workspace_id`, `id`) for workspace isolation
- Soft delete support via `time_deleted` timestamp
- Status field with default 'active'
- Indexes for common queries:
  - Workspace + soft delete filter (most queries)
  - Status filtering

---

### 0029_add_email_verified.sql
**Purpose:** Add email verification tracking to accounts

```sql
ALTER TABLE `account` ADD COLUMN `email_verified` boolean NOT NULL DEFAULT false;
CREATE INDEX `idx_account_email_verified` ON `account` (`email_verified`);
```

**Features:**
- Boolean field with default `false` (new accounts unverified)
- Index for filtering verified/unverified accounts
- Used by `EmailVerification.isVerified()` function

---

### 0030_add_password_hash.sql
**Purpose:** Add password storage for non-OAuth authentication

```sql
ALTER TABLE `account` ADD COLUMN `password_hash` varchar(255);
CREATE INDEX `idx_account_password_exists` ON `account` ((`password_hash` IS NOT NULL));
```

**Features:**
- Nullable field (OAuth accounts won't have password)
- VARCHAR(255) sufficient for bcrypt/argon2 hashes
- Index for filtering password vs OAuth accounts
- Used by `PasswordReset.resetPassword()` function

---

## Schema Updates

### AccountTable (`src/schema/account.sql.ts`)

**Before:**
```typescript
export const AccountTable = mysqlTable(
  "account",
  {
    id: id(),
    ...timestamps,
    email: varchar("email", { length: 255 }).notNull(),
  },
  (table) => [uniqueIndex("email").on(table.email)],
)
```

**After:**
```typescript
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
```

**Changes:**
- Added `emailVerified` boolean field
- Added `passwordHash` varchar field
- Added indexes for both new fields

---

## Code Updates

### Password Reset (`src/password-reset.ts`)

**Before:**
```typescript
// TODO: Implement password hashing
// TODO: Add password_hash field to schema
```

**After:**
```typescript
const passwordHash = newPassword // TODO: Replace with actual hash
await Database.use((tx) =>
  tx.update(AccountTable)
    .set({ passwordHash, timeUpdated: sql`now()` })
    .where(eq(AccountTable.id, verification.accountID!))
)
```

**Status:** ⚠️ Stores plaintext (development only), needs bcrypt/argon2 implementation

---

### Email Verification (`src/email-verification.ts`)

**Before:**
```typescript
// TODO: Add emailVerified field to schema
// emailVerified: true,
return true // Placeholder
```

**After:**
```typescript
await Database.use((tx) =>
  tx.update(AccountTable)
    .set({ emailVerified: true, timeUpdated: sql`now()` })
    .where(eq(AccountTable.id, verificationData.accountID))
)

return account.emailVerified === true
```

**Status:** ✅ Fully implemented

---

## How to Apply Migrations

### Development
```bash
# Apply all pending migrations
cd packages/console/core
bun run drizzle-kit push

# Or generate and review SQL first
bun run drizzle-kit generate
bun run drizzle-kit migrate
```

### Production
```bash
# Always review SQL before applying
cat migrations/0028_add_test_table.sql
cat migrations/0029_add_email_verified.sql
cat migrations/0030_add_password_hash.sql

# Apply to production database
mysql -h <host> -u <user> -p <database> < migrations/0028_add_test_table.sql
mysql -h <host> -u <user> -p <database> < migrations/0029_add_email_verified.sql
mysql -h <host> -u <user> -p <database> < migrations/0030_add_password_hash.sql
```

---

## Migration Testing Checklist

### Test Entity Table
- [ ] Create test entity
- [ ] Read test entity by ID
- [ ] List test entities with filtering by status
- [ ] Update test entity
- [ ] Soft delete test entity (verify time_deleted is set)
- [ ] Hard delete test entity
- [ ] Verify workspace isolation (can't access other workspace entities)
- [ ] Test indexes are used (EXPLAIN queries)

### Email Verified Field
- [ ] New accounts default to `email_verified = false`
- [ ] Verification sets `email_verified = true`
- [ ] `isVerified()` function returns correct status
- [ ] Resend verification works for unverified accounts
- [ ] Index on `email_verified` is used in queries

### Password Hash Field
- [ ] OAuth accounts have NULL `password_hash`
- [ ] Password reset stores hash correctly
- [ ] Index on `password_hash IS NOT NULL` works
- [ ] Can distinguish between OAuth and password accounts

---

## Index Performance

### Test Table Indexes

**idx_test_workspace_deleted (workspace_id, time_deleted):**
```sql
-- Used by: list(), count(), search(), etc.
SELECT * FROM test
WHERE workspace_id = 'wrk_xxx' AND time_deleted IS NULL;
-- Uses index: idx_test_workspace_deleted
```

**idx_test_status (status):**
```sql
-- Used by: list({ status: 'active' })
SELECT * FROM test
WHERE workspace_id = 'wrk_xxx'
  AND time_deleted IS NULL
  AND status = 'active';
-- Uses indexes: idx_test_workspace_deleted, idx_test_status
```

### Account Table Indexes

**idx_account_email_verified (email_verified):**
```sql
-- Filter unverified accounts
SELECT * FROM account WHERE email_verified = false;
-- Uses index: idx_account_email_verified
```

**idx_account_password_exists (password_hash IS NOT NULL):**
```sql
-- Find password-based accounts
SELECT * FROM account WHERE password_hash IS NOT NULL;
-- Uses index: idx_account_password_exists
```

---

## Data Integrity

### Constraints
1. **Test table:**
   - Primary key ensures unique `(workspace_id, id)` pairs
   - NOT NULL on `workspace_id`, `id`, `name`, `status`
   - Default value for `status` ('active')

2. **Account table:**
   - `email_verified` NOT NULL with default `false`
   - `password_hash` nullable (optional for OAuth users)
   - Existing `email` unique constraint preserved

### Workspace Isolation
All queries on `test` table MUST include:
```typescript
and(
  eq(TestTable.workspaceID, Actor.workspace()),
  isNull(TestTable.timeDeleted)
)
```

This is enforced in all CRUD functions.

---

## Rollback Plan

### Rollback 0030 (password_hash)
```sql
DROP INDEX `idx_account_password_exists` ON `account`;
ALTER TABLE `account` DROP COLUMN `password_hash`;
```

### Rollback 0029 (email_verified)
```sql
DROP INDEX `idx_account_email_verified` ON `account`;
ALTER TABLE `account` DROP COLUMN `email_verified`;
```

### Rollback 0028 (test table)
```sql
DROP INDEX `idx_test_status` ON `test`;
DROP INDEX `idx_test_workspace_deleted` ON `test`;
DROP TABLE `test`;
```

**Warning:** Rollback will result in data loss for test entities and email verification status.

---

## Outstanding TODOs

### Critical (Before Production)
1. **Password Hashing:** Replace plaintext storage with bcrypt/argon2
   ```typescript
   // In password-reset.ts
   import bcrypt from 'bcrypt'
   const passwordHash = await bcrypt.hash(newPassword, 12)
   ```

2. **Migration Testing:** Run full test suite on staging database

3. **Index Monitoring:** Verify indexes are being used in production queries

### Optional Enhancements
1. Add `last_password_change` timestamp to account table
2. Add password history table (prevent reuse)
3. Add `email_verification_sent_at` timestamp
4. Add rate limiting table (replace in-memory maps)

---

## Files Modified

1. **New Migrations:**
   - `migrations/0028_add_test_table.sql`
   - `migrations/0029_add_email_verified.sql`
   - `migrations/0030_add_password_hash.sql`

2. **Schema Updates:**
   - `src/schema/account.sql.ts` - Added 2 fields and 2 indexes

3. **Code Updates:**
   - `src/password-reset.ts` - Uses `passwordHash` field
   - `src/email-verification.ts` - Uses `emailVerified` field

---

## Summary

✅ **3 migrations created** for new features
✅ **Schema updated** with new fields and indexes
✅ **Code updated** to use database fields instead of TODOs
✅ **TypeScript compilation** passes
⚠️ **Password hashing** still needs implementation (bcrypt/argon2)

**Ready for:** Staging deployment and testing
**Blocked on:** Password hashing implementation for production

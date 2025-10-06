ALTER TABLE `account` ADD COLUMN `email_verified` boolean NOT NULL DEFAULT false;
--> statement-breakpoint
CREATE INDEX `idx_account_email_verified` ON `account` (`email_verified`);

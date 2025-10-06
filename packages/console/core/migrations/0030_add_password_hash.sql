ALTER TABLE `account` ADD COLUMN `password_hash` varchar(255);
--> statement-breakpoint
CREATE INDEX `idx_account_password_exists` ON `account` ((`password_hash` IS NOT NULL));

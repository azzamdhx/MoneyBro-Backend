-- Rollback optimization indexes

DROP INDEX IF EXISTS idx_transactions_user_reftype_date;
DROP INDEX IF EXISTS idx_accounts_user_type;
DROP INDEX IF EXISTS idx_transaction_entries_account;
DROP INDEX IF EXISTS idx_expenses_category_total;

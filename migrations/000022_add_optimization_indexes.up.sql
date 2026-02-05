-- Optimization indexes for improved query performance

-- Index for GetMonthlyObligationByReferenceType queries
CREATE INDEX IF NOT EXISTS idx_transactions_user_reftype_date 
ON transactions(user_id, reference_type, transaction_date);

-- Index for GetAccountsByType queries
CREATE INDEX IF NOT EXISTS idx_accounts_user_type 
ON accounts(user_id, account_type);

-- Index for balance calculations on transaction_entries
CREATE INDEX IF NOT EXISTS idx_transaction_entries_account 
ON transaction_entries(account_id);

-- Index for category expense aggregation queries
CREATE INDEX IF NOT EXISTS idx_expenses_category_total 
ON expenses(category_id, unit_price, quantity);

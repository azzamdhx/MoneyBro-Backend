CREATE TABLE IF NOT EXISTS transaction_entries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    transaction_id UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    debit BIGINT NOT NULL DEFAULT 0,
    credit BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT chk_debit_or_credit CHECK (
        (debit > 0 AND credit = 0) OR (debit = 0 AND credit > 0)
    )
);

CREATE INDEX idx_transaction_entries_transaction_id ON transaction_entries(transaction_id);
CREATE INDEX idx_transaction_entries_account_id ON transaction_entries(account_id);

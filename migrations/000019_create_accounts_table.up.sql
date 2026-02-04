CREATE TABLE IF NOT EXISTS accounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    account_type VARCHAR(20) NOT NULL CHECK (account_type IN ('ASSET', 'LIABILITY', 'INCOME', 'EXPENSE')),
    current_balance BIGINT NOT NULL DEFAULT 0,
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    reference_id UUID,
    reference_type VARCHAR(50),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_accounts_user_id ON accounts(user_id);
CREATE INDEX idx_accounts_user_type ON accounts(user_id, account_type);
CREATE INDEX idx_accounts_reference ON accounts(reference_id, reference_type);
CREATE UNIQUE INDEX idx_accounts_user_default ON accounts(user_id) WHERE is_default = TRUE;

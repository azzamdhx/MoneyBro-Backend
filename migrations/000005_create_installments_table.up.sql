CREATE TABLE installments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    actual_amount BIGINT NOT NULL,
    loan_amount BIGINT NOT NULL,
    monthly_payment BIGINT NOT NULL,
    tenor INT NOT NULL,
    start_date DATE NOT NULL,
    due_day INT NOT NULL CHECK (due_day >= 1 AND due_day <= 31),
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'COMPLETED')),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE debts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    person_name VARCHAR(255) NOT NULL,
    actual_amount BIGINT NOT NULL,
    loan_amount BIGINT,
    payment_type VARCHAR(20) NOT NULL CHECK (payment_type IN ('ONE_TIME', 'INSTALLMENT')),
    monthly_payment BIGINT,
    tenor INT,
    due_date DATE,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'COMPLETED')),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

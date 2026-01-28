CREATE TABLE installment_payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    installment_id UUID NOT NULL REFERENCES installments(id) ON DELETE CASCADE,
    payment_number INT NOT NULL,
    amount BIGINT NOT NULL,
    paid_at DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

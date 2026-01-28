CREATE TABLE incomes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES income_categories(id) ON DELETE RESTRICT,
    source_name VARCHAR(255) NOT NULL,
    amount BIGINT NOT NULL,
    income_type VARCHAR(50) NOT NULL CHECK (income_type IN (
        'SALARY', 
        'FREELANCE', 
        'INVESTMENT', 
        'GIFT', 
        'BONUS', 
        'REFUND', 
        'BUSINESS', 
        'OTHER'
    )),
    income_date DATE NOT NULL,
    is_recurring BOOLEAN NOT NULL DEFAULT FALSE,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

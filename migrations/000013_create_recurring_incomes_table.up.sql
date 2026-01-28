CREATE TABLE recurring_incomes (
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
    recurring_day INT NOT NULL CHECK (recurring_day >= 1 AND recurring_day <= 31),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

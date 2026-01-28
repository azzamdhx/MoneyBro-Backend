CREATE TABLE expense_templates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    item_name VARCHAR(255) NOT NULL,
    unit_price BIGINT NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    recurring_day INT CHECK (recurring_day >= 1 AND recurring_day <= 31),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create expense_template_groups table
CREATE TABLE IF NOT EXISTS expense_template_groups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    recurring_day INTEGER,
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create expense_template_items table
CREATE TABLE IF NOT EXISTS expense_template_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    group_id UUID NOT NULL REFERENCES expense_template_groups(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    item_name VARCHAR(255) NOT NULL,
    unit_price BIGINT NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_expense_template_groups_user_id ON expense_template_groups(user_id);
CREATE INDEX IF NOT EXISTS idx_expense_template_items_group_id ON expense_template_items(group_id);
CREATE INDEX IF NOT EXISTS idx_expense_template_items_category_id ON expense_template_items(category_id);

-- Drop old expense_templates table
DROP TABLE IF EXISTS expense_templates;

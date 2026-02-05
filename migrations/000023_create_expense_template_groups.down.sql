-- Drop indexes
DROP INDEX IF EXISTS idx_expense_template_items_category_id;
DROP INDEX IF EXISTS idx_expense_template_items_group_id;
DROP INDEX IF EXISTS idx_expense_template_groups_user_id;

-- Drop new tables
DROP TABLE IF EXISTS expense_template_items;
DROP TABLE IF EXISTS expense_template_groups;

-- Recreate old expense_templates table
CREATE TABLE IF NOT EXISTS expense_templates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    item_name VARCHAR(255) NOT NULL,
    unit_price BIGINT NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 1,
    recurring_day INTEGER,
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

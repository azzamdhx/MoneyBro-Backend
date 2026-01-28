CREATE INDEX idx_income_categories_user_id ON income_categories(user_id);

CREATE INDEX idx_incomes_user_id ON incomes(user_id);
CREATE INDEX idx_incomes_category_id ON incomes(category_id);
CREATE INDEX idx_incomes_user_date ON incomes(user_id, income_date);

CREATE INDEX idx_recurring_incomes_user_id ON recurring_incomes(user_id);
CREATE INDEX idx_recurring_incomes_user_active ON recurring_incomes(user_id, is_active);
CREATE INDEX idx_recurring_incomes_day ON recurring_incomes(recurring_day);

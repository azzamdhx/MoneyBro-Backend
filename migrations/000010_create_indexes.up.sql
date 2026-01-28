CREATE INDEX idx_categories_user_id ON categories(user_id);

CREATE INDEX idx_expenses_user_id ON expenses(user_id);
CREATE INDEX idx_expenses_category_id ON expenses(category_id);
CREATE INDEX idx_expenses_user_date ON expenses(user_id, expense_date);

CREATE INDEX idx_expense_templates_user_id ON expense_templates(user_id);

CREATE INDEX idx_installments_user_id ON installments(user_id);
CREATE INDEX idx_installments_user_status ON installments(user_id, status);
CREATE INDEX idx_installments_due_day ON installments(due_day);

CREATE INDEX idx_installment_payments_installment_id ON installment_payments(installment_id);

CREATE INDEX idx_debts_user_id ON debts(user_id);
CREATE INDEX idx_debts_user_status ON debts(user_id, status);
CREATE INDEX idx_debts_due_date ON debts(due_date);

CREATE INDEX idx_debt_payments_debt_id ON debt_payments(debt_id);

CREATE INDEX idx_notification_logs_user_id ON notification_logs(user_id);
CREATE INDEX idx_notification_logs_reference ON notification_logs(reference_id, type);

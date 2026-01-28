DROP INDEX IF EXISTS idx_categories_user_id;

DROP INDEX IF EXISTS idx_expenses_user_id;
DROP INDEX IF EXISTS idx_expenses_category_id;
DROP INDEX IF EXISTS idx_expenses_user_date;

DROP INDEX IF EXISTS idx_expense_templates_user_id;

DROP INDEX IF EXISTS idx_installments_user_id;
DROP INDEX IF EXISTS idx_installments_user_status;
DROP INDEX IF EXISTS idx_installments_due_day;

DROP INDEX IF EXISTS idx_installment_payments_installment_id;

DROP INDEX IF EXISTS idx_debts_user_id;
DROP INDEX IF EXISTS idx_debts_user_status;
DROP INDEX IF EXISTS idx_debts_due_date;

DROP INDEX IF EXISTS idx_debt_payments_debt_id;

DROP INDEX IF EXISTS idx_notification_logs_user_id;
DROP INDEX IF EXISTS idx_notification_logs_reference;

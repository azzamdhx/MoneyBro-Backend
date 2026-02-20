-- Rollback: Remove demo account and all related data
DO $$
DECLARE
  v_user_id UUID := 'a0000000-0000-0000-0000-000000000001';
BEGIN
  DELETE FROM savings_contributions WHERE savings_goal_id IN (SELECT id FROM savings_goals WHERE user_id = v_user_id);
  DELETE FROM savings_goals WHERE user_id = v_user_id;
  DELETE FROM debt_payments WHERE debt_id IN (SELECT id FROM debts WHERE user_id = v_user_id);
  DELETE FROM debts WHERE user_id = v_user_id;
  DELETE FROM installment_payments WHERE installment_id IN (SELECT id FROM installments WHERE user_id = v_user_id);
  DELETE FROM installments WHERE user_id = v_user_id;
  DELETE FROM expense_template_items WHERE group_id IN (SELECT id FROM expense_template_groups WHERE user_id = v_user_id);
  DELETE FROM expense_template_groups WHERE user_id = v_user_id;
  DELETE FROM expenses WHERE user_id = v_user_id;
  DELETE FROM incomes WHERE user_id = v_user_id;
  DELETE FROM recurring_incomes WHERE user_id = v_user_id;
  DELETE FROM transaction_entries WHERE transaction_id IN (SELECT id FROM transactions WHERE user_id = v_user_id);
  DELETE FROM transactions WHERE user_id = v_user_id;
  DELETE FROM notification_logs WHERE user_id = v_user_id;
  DELETE FROM two_fa_codes WHERE user_id = v_user_id;
  DELETE FROM password_reset_tokens WHERE user_id = v_user_id;
  DELETE FROM accounts WHERE user_id = v_user_id;
  DELETE FROM categories WHERE user_id = v_user_id;
  DELETE FROM income_categories WHERE user_id = v_user_id;
  DELETE FROM users WHERE id = v_user_id;
END;
$$;

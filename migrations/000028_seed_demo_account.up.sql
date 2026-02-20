-- ============================================
-- MoneyBro Demo Account Seed
-- Email: demo@moneybro.my.id
-- Password: Moneybro123$$
-- ============================================
-- Includes transactions + transaction_entries
-- for the double-entry ledger system.
-- ============================================

DO $$
DECLARE
  v_user_id UUID := 'a0000000-0000-0000-0000-000000000001';
  v_cash_account_id UUID := 'a0000000-0000-0000-0000-000000000010';
  -- Expense categories
  v_cat_makanan UUID := 'a0000000-0000-0000-0000-000000000101';
  v_cat_transport UUID := 'a0000000-0000-0000-0000-000000000102';
  v_cat_belanja UUID := 'a0000000-0000-0000-0000-000000000103';
  v_cat_hiburan UUID := 'a0000000-0000-0000-0000-000000000104';
  v_cat_kesehatan UUID := 'a0000000-0000-0000-0000-000000000105';
  v_cat_pendidikan UUID := 'a0000000-0000-0000-0000-000000000106';
  v_cat_tagihan UUID := 'a0000000-0000-0000-0000-000000000107';
  v_cat_lainnya UUID := 'a0000000-0000-0000-0000-000000000108';
  -- Income categories
  v_icat_gaji UUID := 'a0000000-0000-0000-0000-000000000201';
  v_icat_freelance UUID := 'a0000000-0000-0000-0000-000000000202';
  v_icat_investasi UUID := 'a0000000-0000-0000-0000-000000000203';
  v_icat_lainnya UUID := 'a0000000-0000-0000-0000-000000000204';
  -- Installments
  v_inst_hp UUID := 'a0000000-0000-0000-0000-000000000301';
  v_inst_laptop UUID := 'a0000000-0000-0000-0000-000000000302';
  -- Debts
  v_debt_teman UUID := 'a0000000-0000-0000-0000-000000000401';
  v_debt_keluarga UUID := 'a0000000-0000-0000-0000-000000000402';
  -- Savings
  v_save_liburan UUID := 'a0000000-0000-0000-0000-000000000501';
  v_save_darurat UUID := 'a0000000-0000-0000-0000-000000000502';
  -- Expense template
  v_template_bulanan UUID := 'a0000000-0000-0000-0000-000000000601';
  -- Accounts
  v_acc_inst_hp UUID := 'a0000000-0000-0000-0000-000000000701';
  v_acc_inst_laptop UUID := 'a0000000-0000-0000-0000-000000000702';
  v_acc_debt_teman UUID := 'a0000000-0000-0000-0000-000000000703';
  v_acc_debt_keluarga UUID := 'a0000000-0000-0000-0000-000000000704';
  v_acc_save_liburan UUID := 'a0000000-0000-0000-0000-000000000705';
  v_acc_save_darurat UUID := 'a0000000-0000-0000-0000-000000000706';
  v_acc_income UUID := 'a0000000-0000-0000-0000-000000000707';
  v_acc_expense UUID := 'a0000000-0000-0000-0000-000000000708';
  -- Installment payments (fixed UUIDs for ledger references)
  v_ip_hp_1 UUID := 'a0000000-0000-0000-0000-000000000811';
  v_ip_hp_2 UUID := 'a0000000-0000-0000-0000-000000000812';
  v_ip_hp_3 UUID := 'a0000000-0000-0000-0000-000000000813';
  v_ip_hp_4 UUID := 'a0000000-0000-0000-0000-000000000814';
  v_ip_lt_1 UUID := 'a0000000-0000-0000-0000-000000000821';
  v_ip_lt_2 UUID := 'a0000000-0000-0000-0000-000000000822';
  v_ip_lt_3 UUID := 'a0000000-0000-0000-0000-000000000823';
  v_ip_lt_4 UUID := 'a0000000-0000-0000-0000-000000000824';
  v_ip_lt_5 UUID := 'a0000000-0000-0000-0000-000000000825';
  v_ip_lt_6 UUID := 'a0000000-0000-0000-0000-000000000826';
  v_ip_lt_7 UUID := 'a0000000-0000-0000-0000-000000000827';
  v_ip_lt_8 UUID := 'a0000000-0000-0000-0000-000000000828';
  -- Debt payments (fixed UUIDs)
  v_dp_budi_1 UUID := 'a0000000-0000-0000-0000-000000000911';
  v_dp_budi_2 UUID := 'a0000000-0000-0000-0000-000000000912';
  v_dp_kakak_1 UUID := 'a0000000-0000-0000-0000-000000000921';
  -- Savings contributions (fixed UUIDs)
  v_sc_lib_1 UUID := 'a0000000-0000-0000-0000-000000000a11';
  v_sc_lib_2 UUID := 'a0000000-0000-0000-0000-000000000a12';
  v_sc_lib_3 UUID := 'a0000000-0000-0000-0000-000000000a13';
  v_sc_dar_1 UUID := 'a0000000-0000-0000-0000-000000000a21';
  v_sc_dar_2 UUID := 'a0000000-0000-0000-0000-000000000a22';
  v_sc_dar_3 UUID := 'a0000000-0000-0000-0000-000000000a23';
  v_sc_dar_4 UUID := 'a0000000-0000-0000-0000-000000000a24';

  -- !! FIXED dates anchored to 2026-02-15 — frontend overrides Date to match !!
  v_now TIMESTAMP := '2026-02-15 12:00:00'::timestamp;
  v_today DATE := '2026-02-15'::date;
  v_this_month DATE := '2026-02-01'::date;
  v_last_month DATE := '2026-01-01'::date;
  -- 6-month spread: 3 months before + anchor + 3 months after
  v_month_nov25 DATE := '2025-11-01'::date;
  v_month_dec25 DATE := '2025-12-01'::date;
  -- v_last_month = Jan 2026, v_this_month = Feb 2026 (already defined)
  v_month_mar26 DATE := '2026-03-01'::date;
  v_month_apr26 DATE := '2026-04-01'::date;
  v_month_may26 DATE := '2026-05-01'::date;

  -- Installment start dates: adjusted so latest payment falls in "current month" (Feb 2026)
  -- Backend formula: periodDate = start_date + (payment_number - 1) months
  -- iPhone: 4 payments, start 3 months ago (Nov 2025) → payment 4 period = Feb 2026
  v_hp_start DATE := '2025-11-01'::date;
  -- MacBook: 8 payments, start 7 months ago (Jul 2025) → payment 8 period = Feb 2026
  v_lt_start DATE := '2025-07-01'::date;

  v_tx_id UUID;
BEGIN

-- ============ USER ============
INSERT INTO users (id, email, password_hash, name, profile_image, two_fa_enabled, created_at)
VALUES (v_user_id, 'demo@moneybro.my.id', '$2a$12$hHQUo9a5zwzUDV8n7/5xLe/uvOh6zL6S0qSiiWaAApHbNd7VtDBiW', 'Demo User', 'BRO-3-B', false, v_now);

-- ============ ACCOUNTS ============
INSERT INTO accounts (id, user_id, name, account_type, current_balance, is_default, created_at)
VALUES (v_cash_account_id, v_user_id, 'Cash', 'ASSET', 0, true, v_now);

INSERT INTO accounts (id, user_id, name, account_type, current_balance, is_default, created_at)
VALUES (v_acc_income, v_user_id, 'Income', 'INCOME', 0, false, v_now);

INSERT INTO accounts (id, user_id, name, account_type, current_balance, is_default, created_at)
VALUES (v_acc_expense, v_user_id, 'Expense', 'EXPENSE', 0, false, v_now);

INSERT INTO accounts (id, user_id, name, account_type, current_balance, is_default, reference_id, reference_type, created_at)
VALUES
(v_acc_inst_hp, v_user_id, 'Cicilan iPhone 15', 'LIABILITY', 0, false, v_inst_hp, 'installment', v_now),
(v_acc_inst_laptop, v_user_id, 'Cicilan MacBook Air', 'LIABILITY', 0, false, v_inst_laptop, 'installment', v_now),
(v_acc_debt_teman, v_user_id, 'Hutang Budi', 'LIABILITY', 0, false, v_debt_teman, 'debt', v_now),
(v_acc_debt_keluarga, v_user_id, 'Hutang Kakak', 'LIABILITY', 0, false, v_debt_keluarga, 'debt', v_now),
(v_acc_save_liburan, v_user_id, 'Tabungan Liburan Bali', 'ASSET', 0, false, v_save_liburan, 'savings_goal', v_now),
(v_acc_save_darurat, v_user_id, 'Dana Darurat', 'ASSET', 0, false, v_save_darurat, 'savings_goal', v_now);

-- ============ EXPENSE CATEGORIES ============
INSERT INTO categories (id, user_id, name, created_at) VALUES
(v_cat_makanan, v_user_id, 'Makanan & Minuman', v_now),
(v_cat_transport, v_user_id, 'Transportasi', v_now),
(v_cat_belanja, v_user_id, 'Belanja', v_now),
(v_cat_hiburan, v_user_id, 'Hiburan', v_now),
(v_cat_kesehatan, v_user_id, 'Kesehatan', v_now),
(v_cat_pendidikan, v_user_id, 'Pendidikan', v_now),
(v_cat_tagihan, v_user_id, 'Tagihan & Utilitas', v_now),
(v_cat_lainnya, v_user_id, 'Lainnya', v_now);

-- ============ INCOME CATEGORIES ============
INSERT INTO income_categories (id, user_id, name, created_at) VALUES
(v_icat_gaji, v_user_id, 'Gaji', v_now),
(v_icat_freelance, v_user_id, 'Freelance', v_now),
(v_icat_investasi, v_user_id, 'Investasi', v_now),
(v_icat_lainnya, v_user_id, 'Lainnya', v_now);

-- ============ EXPENSES - THIS MONTH ============
INSERT INTO expenses (id, user_id, category_id, item_name, unit_price, quantity, expense_date, created_at) VALUES
(uuid_generate_v4(), v_user_id, v_cat_makanan, 'Makan siang warteg', 15000, 1, v_this_month + 1, v_now),
(uuid_generate_v4(), v_user_id, v_cat_makanan, 'Kopi Starbucks', 55000, 1, v_this_month + 2, v_now),
(uuid_generate_v4(), v_user_id, v_cat_makanan, 'Makan malam Padang', 35000, 1, v_this_month + 3, v_now),
(uuid_generate_v4(), v_user_id, v_cat_makanan, 'Groceries Indomaret', 120000, 1, v_this_month + 5, v_now),
(uuid_generate_v4(), v_user_id, v_cat_transport, 'Grab ke kantor', 25000, 1, v_this_month + 1, v_now),
(uuid_generate_v4(), v_user_id, v_cat_transport, 'Bensin motor', 50000, 1, v_this_month + 4, v_now),
(uuid_generate_v4(), v_user_id, v_cat_transport, 'Gojek pulang kerja', 18000, 1, v_this_month + 6, v_now),
(uuid_generate_v4(), v_user_id, v_cat_belanja, 'Kaos Uniqlo', 199000, 1, v_this_month + 3, v_now),
(uuid_generate_v4(), v_user_id, v_cat_hiburan, 'Nonton bioskop', 50000, 2, v_this_month + 7, v_now),
(uuid_generate_v4(), v_user_id, v_cat_hiburan, 'Langganan Spotify', 54990, 1, v_this_month + 1, v_now),
(uuid_generate_v4(), v_user_id, v_cat_tagihan, 'Listrik PLN', 350000, 1, v_this_month + 5, v_now),
(uuid_generate_v4(), v_user_id, v_cat_tagihan, 'Internet IndiHome', 399000, 1, v_this_month + 5, v_now),
(uuid_generate_v4(), v_user_id, v_cat_kesehatan, 'Vitamin C & D', 85000, 1, v_this_month + 2, v_now);

-- ============ EXPENSES - LAST MONTH ============
INSERT INTO expenses (id, user_id, category_id, item_name, unit_price, quantity, expense_date, created_at) VALUES
(uuid_generate_v4(), v_user_id, v_cat_makanan, 'Makan siang kantin', 20000, 1, v_last_month + 2, v_now),
(uuid_generate_v4(), v_user_id, v_cat_makanan, 'Makan malam sushi', 150000, 1, v_last_month + 8, v_now),
(uuid_generate_v4(), v_user_id, v_cat_makanan, 'Kopi Janji Jiwa', 24000, 1, v_last_month + 10, v_now),
(uuid_generate_v4(), v_user_id, v_cat_makanan, 'Groceries Alfamart', 95000, 1, v_last_month + 15, v_now),
(uuid_generate_v4(), v_user_id, v_cat_transport, 'Bensin motor', 50000, 1, v_last_month + 5, v_now),
(uuid_generate_v4(), v_user_id, v_cat_transport, 'Grab ke mall', 30000, 1, v_last_month + 12, v_now),
(uuid_generate_v4(), v_user_id, v_cat_belanja, 'Celana jeans', 299000, 1, v_last_month + 10, v_now),
(uuid_generate_v4(), v_user_id, v_cat_hiburan, 'Langganan Spotify', 54990, 1, v_last_month + 1, v_now),
(uuid_generate_v4(), v_user_id, v_cat_hiburan, 'Langganan Netflix', 54000, 1, v_last_month + 1, v_now),
(uuid_generate_v4(), v_user_id, v_cat_tagihan, 'Listrik PLN', 320000, 1, v_last_month + 5, v_now),
(uuid_generate_v4(), v_user_id, v_cat_tagihan, 'Internet IndiHome', 399000, 1, v_last_month + 5, v_now),
(uuid_generate_v4(), v_user_id, v_cat_pendidikan, 'Kursus online Udemy', 149000, 1, v_last_month + 7, v_now),
(uuid_generate_v4(), v_user_id, v_cat_kesehatan, 'Obat flu', 45000, 1, v_last_month + 18, v_now);

-- ============ EXPENSES - NOVEMBER 2025 ============
INSERT INTO expenses (id, user_id, category_id, item_name, unit_price, quantity, expense_date, created_at) VALUES
(uuid_generate_v4(), v_user_id, v_cat_makanan, 'Makan siang nasi goreng', 18000, 1, v_month_nov25 + 2, v_now),
(uuid_generate_v4(), v_user_id, v_cat_makanan, 'Kopi Kenangan', 28000, 1, v_month_nov25 + 3, v_now),
(uuid_generate_v4(), v_user_id, v_cat_makanan, 'Groceries Superindo', 180000, 1, v_month_nov25 + 7, v_now),
(uuid_generate_v4(), v_user_id, v_cat_makanan, 'Makan malam ayam geprek', 25000, 1, v_month_nov25 + 10, v_now),
(uuid_generate_v4(), v_user_id, v_cat_makanan, 'Snack Indomaret', 32000, 1, v_month_nov25 + 18, v_now),
(uuid_generate_v4(), v_user_id, v_cat_transport, 'Bensin motor', 55000, 1, v_month_nov25 + 5, v_now),
(uuid_generate_v4(), v_user_id, v_cat_transport, 'Grab ke kantor', 22000, 1, v_month_nov25 + 8, v_now),
(uuid_generate_v4(), v_user_id, v_cat_belanja, 'Sepatu olahraga', 350000, 1, v_month_nov25 + 12, v_now),
(uuid_generate_v4(), v_user_id, v_cat_hiburan, 'Langganan Spotify', 54990, 1, v_month_nov25 + 1, v_now),
(uuid_generate_v4(), v_user_id, v_cat_hiburan, 'Langganan Netflix', 54000, 1, v_month_nov25 + 1, v_now),
(uuid_generate_v4(), v_user_id, v_cat_tagihan, 'Listrik PLN', 310000, 1, v_month_nov25 + 5, v_now),
(uuid_generate_v4(), v_user_id, v_cat_tagihan, 'Internet IndiHome', 399000, 1, v_month_nov25 + 5, v_now),
(uuid_generate_v4(), v_user_id, v_cat_kesehatan, 'Multivitamin', 95000, 1, v_month_nov25 + 14, v_now),
(uuid_generate_v4(), v_user_id, v_cat_pendidikan, 'Buku programming', 185000, 1, v_month_nov25 + 20, v_now);

-- ============ EXPENSES - DECEMBER 2025 ============
INSERT INTO expenses (id, user_id, category_id, item_name, unit_price, quantity, expense_date, created_at) VALUES
(uuid_generate_v4(), v_user_id, v_cat_makanan, 'Makan siang bakso', 20000, 1, v_month_dec25 + 2, v_now),
(uuid_generate_v4(), v_user_id, v_cat_makanan, 'Kopi Starbucks Christmas', 65000, 1, v_month_dec25 + 5, v_now),
(uuid_generate_v4(), v_user_id, v_cat_makanan, 'Makan malam BBQ keluarga', 250000, 1, v_month_dec25 + 24, v_now),
(uuid_generate_v4(), v_user_id, v_cat_makanan, 'Groceries Alfamart', 150000, 1, v_month_dec25 + 10, v_now),
(uuid_generate_v4(), v_user_id, v_cat_makanan, 'Snack & kue natal', 85000, 1, v_month_dec25 + 22, v_now),
(uuid_generate_v4(), v_user_id, v_cat_transport, 'Bensin motor', 60000, 1, v_month_dec25 + 6, v_now),
(uuid_generate_v4(), v_user_id, v_cat_transport, 'Grab ke mall', 35000, 1, v_month_dec25 + 15, v_now),
(uuid_generate_v4(), v_user_id, v_cat_transport, 'Tiket kereta mudik', 350000, 1, v_month_dec25 + 20, v_now),
(uuid_generate_v4(), v_user_id, v_cat_belanja, 'Kado natal', 200000, 1, v_month_dec25 + 18, v_now),
(uuid_generate_v4(), v_user_id, v_cat_belanja, 'Baju baru akhir tahun', 250000, 1, v_month_dec25 + 14, v_now),
(uuid_generate_v4(), v_user_id, v_cat_hiburan, 'Langganan Spotify', 54990, 1, v_month_dec25 + 1, v_now),
(uuid_generate_v4(), v_user_id, v_cat_hiburan, 'Langganan Netflix', 54000, 1, v_month_dec25 + 1, v_now),
(uuid_generate_v4(), v_user_id, v_cat_hiburan, 'Nonton bioskop', 50000, 2, v_month_dec25 + 26, v_now),
(uuid_generate_v4(), v_user_id, v_cat_tagihan, 'Listrik PLN', 380000, 1, v_month_dec25 + 5, v_now),
(uuid_generate_v4(), v_user_id, v_cat_tagihan, 'Internet IndiHome', 399000, 1, v_month_dec25 + 5, v_now),
(uuid_generate_v4(), v_user_id, v_cat_lainnya, 'Angpao natal', 500000, 1, v_month_dec25 + 25, v_now);

-- ============ EXPENSES - MARCH 2026 ============
INSERT INTO expenses (id, user_id, category_id, item_name, unit_price, quantity, expense_date, created_at) VALUES
(uuid_generate_v4(), v_user_id, v_cat_makanan, 'Makan siang soto', 20000, 1, v_month_mar26 + 2, v_now),
(uuid_generate_v4(), v_user_id, v_cat_makanan, 'Kopi Fore Coffee', 32000, 1, v_month_mar26 + 4, v_now),
(uuid_generate_v4(), v_user_id, v_cat_makanan, 'Groceries Indomaret', 135000, 1, v_month_mar26 + 8, v_now),
(uuid_generate_v4(), v_user_id, v_cat_makanan, 'Makan malam pizza', 89000, 1, v_month_mar26 + 12, v_now),
(uuid_generate_v4(), v_user_id, v_cat_makanan, 'Snack kantor', 28000, 1, v_month_mar26 + 15, v_now),
(uuid_generate_v4(), v_user_id, v_cat_transport, 'Bensin motor', 50000, 1, v_month_mar26 + 3, v_now),
(uuid_generate_v4(), v_user_id, v_cat_transport, 'Gojek ke meeting', 28000, 1, v_month_mar26 + 9, v_now),
(uuid_generate_v4(), v_user_id, v_cat_transport, 'Grab pulang lembur', 20000, 1, v_month_mar26 + 14, v_now),
(uuid_generate_v4(), v_user_id, v_cat_belanja, 'Kemeja kerja baru', 179000, 1, v_month_mar26 + 10, v_now),
(uuid_generate_v4(), v_user_id, v_cat_hiburan, 'Langganan Spotify', 54990, 1, v_month_mar26 + 1, v_now),
(uuid_generate_v4(), v_user_id, v_cat_hiburan, 'Langganan Netflix', 54000, 1, v_month_mar26 + 1, v_now),
(uuid_generate_v4(), v_user_id, v_cat_hiburan, 'Top-up game mobile', 99000, 1, v_month_mar26 + 16, v_now),
(uuid_generate_v4(), v_user_id, v_cat_tagihan, 'Listrik PLN', 340000, 1, v_month_mar26 + 5, v_now),
(uuid_generate_v4(), v_user_id, v_cat_tagihan, 'Internet IndiHome', 399000, 1, v_month_mar26 + 5, v_now),
(uuid_generate_v4(), v_user_id, v_cat_kesehatan, 'Check-up dokter', 250000, 1, v_month_mar26 + 18, v_now),
(uuid_generate_v4(), v_user_id, v_cat_pendidikan, 'Langganan Coursera', 199000, 1, v_month_mar26 + 7, v_now);

-- ============ EXPENSES - APRIL 2026 ============
INSERT INTO expenses (id, user_id, category_id, item_name, unit_price, quantity, expense_date, created_at) VALUES
(uuid_generate_v4(), v_user_id, v_cat_makanan, 'Makan siang pecel lele', 17000, 1, v_month_apr26 + 1, v_now),
(uuid_generate_v4(), v_user_id, v_cat_makanan, 'Kopi Janji Jiwa', 25000, 1, v_month_apr26 + 3, v_now),
(uuid_generate_v4(), v_user_id, v_cat_makanan, 'Groceries Giant', 165000, 1, v_month_apr26 + 7, v_now),
(uuid_generate_v4(), v_user_id, v_cat_makanan, 'Makan malam ramen', 75000, 1, v_month_apr26 + 11, v_now),
(uuid_generate_v4(), v_user_id, v_cat_makanan, 'Snack kantor', 45000, 1, v_month_apr26 + 15, v_now),
(uuid_generate_v4(), v_user_id, v_cat_transport, 'Bensin motor', 55000, 1, v_month_apr26 + 4, v_now),
(uuid_generate_v4(), v_user_id, v_cat_transport, 'Grab bolak-balik', 40000, 1, v_month_apr26 + 8, v_now),
(uuid_generate_v4(), v_user_id, v_cat_belanja, 'Tas ransel baru', 299000, 1, v_month_apr26 + 13, v_now),
(uuid_generate_v4(), v_user_id, v_cat_hiburan, 'Langganan Spotify', 54990, 1, v_month_apr26 + 1, v_now),
(uuid_generate_v4(), v_user_id, v_cat_hiburan, 'Langganan Netflix', 54000, 1, v_month_apr26 + 1, v_now),
(uuid_generate_v4(), v_user_id, v_cat_hiburan, 'Langganan Disney+', 39000, 1, v_month_apr26 + 1, v_now),
(uuid_generate_v4(), v_user_id, v_cat_tagihan, 'Listrik PLN', 360000, 1, v_month_apr26 + 5, v_now),
(uuid_generate_v4(), v_user_id, v_cat_tagihan, 'Internet IndiHome', 399000, 1, v_month_apr26 + 5, v_now),
(uuid_generate_v4(), v_user_id, v_cat_kesehatan, 'Obat alergi', 65000, 1, v_month_apr26 + 10, v_now),
(uuid_generate_v4(), v_user_id, v_cat_lainnya, 'Laundry kiloan', 75000, 1, v_month_apr26 + 6, v_now),
(uuid_generate_v4(), v_user_id, v_cat_lainnya, 'Parkir bulanan', 150000, 1, v_month_apr26 + 2, v_now);

-- ============ EXPENSES - MAY 2026 ============
INSERT INTO expenses (id, user_id, category_id, item_name, unit_price, quantity, expense_date, created_at) VALUES
(uuid_generate_v4(), v_user_id, v_cat_makanan, 'Makan siang warteg', 15000, 1, v_month_may26 + 1, v_now),
(uuid_generate_v4(), v_user_id, v_cat_makanan, 'Kopi Kulo', 22000, 1, v_month_may26 + 4, v_now),
(uuid_generate_v4(), v_user_id, v_cat_makanan, 'Groceries Alfamart', 140000, 1, v_month_may26 + 8, v_now),
(uuid_generate_v4(), v_user_id, v_cat_makanan, 'Makan malam seafood', 120000, 1, v_month_may26 + 14, v_now),
(uuid_generate_v4(), v_user_id, v_cat_makanan, 'Snack Indomaret', 35000, 1, v_month_may26 + 18, v_now),
(uuid_generate_v4(), v_user_id, v_cat_transport, 'Bensin motor', 60000, 1, v_month_may26 + 5, v_now),
(uuid_generate_v4(), v_user_id, v_cat_transport, 'Grab ke kantor', 25000, 1, v_month_may26 + 9, v_now),
(uuid_generate_v4(), v_user_id, v_cat_transport, 'Parkir motor', 30000, 1, v_month_may26 + 12, v_now),
(uuid_generate_v4(), v_user_id, v_cat_belanja, 'Celana chino', 249000, 1, v_month_may26 + 10, v_now),
(uuid_generate_v4(), v_user_id, v_cat_belanja, 'Sandal kulit', 129000, 1, v_month_may26 + 16, v_now),
(uuid_generate_v4(), v_user_id, v_cat_hiburan, 'Langganan Spotify', 54990, 1, v_month_may26 + 1, v_now),
(uuid_generate_v4(), v_user_id, v_cat_hiburan, 'Langganan Netflix', 54000, 1, v_month_may26 + 1, v_now),
(uuid_generate_v4(), v_user_id, v_cat_hiburan, 'Tiket konser', 350000, 1, v_month_may26 + 20, v_now),
(uuid_generate_v4(), v_user_id, v_cat_tagihan, 'Listrik PLN', 330000, 1, v_month_may26 + 5, v_now),
(uuid_generate_v4(), v_user_id, v_cat_tagihan, 'Internet IndiHome', 399000, 1, v_month_may26 + 5, v_now),
(uuid_generate_v4(), v_user_id, v_cat_pendidikan, 'Kursus bahasa Inggris', 250000, 1, v_month_may26 + 7, v_now),
(uuid_generate_v4(), v_user_id, v_cat_lainnya, 'Potong rambut', 50000, 1, v_month_may26 + 3, v_now);

-- ============ INCOMES - THIS MONTH ============
INSERT INTO incomes (id, user_id, category_id, source_name, amount, income_type, income_date, is_recurring, created_at) VALUES
(uuid_generate_v4(), v_user_id, v_icat_gaji, 'Gaji PT Maju Jaya', 8500000, 'SALARY', v_this_month + 1, true, v_now),
(uuid_generate_v4(), v_user_id, v_icat_freelance, 'Project website klien', 2500000, 'FREELANCE', v_this_month + 5, false, v_now),
(uuid_generate_v4(), v_user_id, v_icat_investasi, 'Dividen saham BBCA', 150000, 'INVESTMENT', v_this_month + 10, false, v_now);

-- ============ INCOMES - LAST MONTH ============
INSERT INTO incomes (id, user_id, category_id, source_name, amount, income_type, income_date, is_recurring, created_at) VALUES
(uuid_generate_v4(), v_user_id, v_icat_gaji, 'Gaji PT Maju Jaya', 8500000, 'SALARY', v_last_month + 1, true, v_now),
(uuid_generate_v4(), v_user_id, v_icat_freelance, 'Design logo', 1000000, 'FREELANCE', v_last_month + 12, false, v_now),
(uuid_generate_v4(), v_user_id, v_icat_lainnya, 'Cashback promo', 50000, 'OTHER', v_last_month + 20, false, v_now);

-- ============ INCOMES - NOVEMBER 2025 ============
INSERT INTO incomes (id, user_id, category_id, source_name, amount, income_type, income_date, is_recurring, created_at) VALUES
(uuid_generate_v4(), v_user_id, v_icat_gaji, 'Gaji PT Maju Jaya', 8500000, 'SALARY', v_month_nov25 + 1, true, v_now),
(uuid_generate_v4(), v_user_id, v_icat_freelance, 'Project design UI', 1800000, 'FREELANCE', v_month_nov25 + 10, false, v_now);

-- ============ INCOMES - DECEMBER 2025 ============
INSERT INTO incomes (id, user_id, category_id, source_name, amount, income_type, income_date, is_recurring, created_at) VALUES
(uuid_generate_v4(), v_user_id, v_icat_gaji, 'Gaji PT Maju Jaya', 8500000, 'SALARY', v_month_dec25 + 1, true, v_now),
(uuid_generate_v4(), v_user_id, v_icat_gaji, 'Bonus THR', 4250000, 'BONUS', v_month_dec25 + 15, false, v_now),
(uuid_generate_v4(), v_user_id, v_icat_freelance, 'Freelance landing page', 800000, 'FREELANCE', v_month_dec25 + 8, false, v_now);

-- ============ INCOMES - MARCH 2026 ============
INSERT INTO incomes (id, user_id, category_id, source_name, amount, income_type, income_date, is_recurring, created_at) VALUES
(uuid_generate_v4(), v_user_id, v_icat_gaji, 'Gaji PT Maju Jaya', 8500000, 'SALARY', v_month_mar26 + 1, true, v_now),
(uuid_generate_v4(), v_user_id, v_icat_freelance, 'Project mobile app', 3000000, 'FREELANCE', v_month_mar26 + 7, false, v_now);

-- ============ INCOMES - APRIL 2026 ============
INSERT INTO incomes (id, user_id, category_id, source_name, amount, income_type, income_date, is_recurring, created_at) VALUES
(uuid_generate_v4(), v_user_id, v_icat_gaji, 'Gaji PT Maju Jaya', 8500000, 'SALARY', v_month_apr26 + 1, true, v_now),
(uuid_generate_v4(), v_user_id, v_icat_freelance, 'Freelance branding', 1200000, 'FREELANCE', v_month_apr26 + 12, false, v_now);

-- ============ INCOMES - MAY 2026 ============
INSERT INTO incomes (id, user_id, category_id, source_name, amount, income_type, income_date, is_recurring, created_at) VALUES
(uuid_generate_v4(), v_user_id, v_icat_gaji, 'Gaji PT Maju Jaya', 8500000, 'SALARY', v_month_may26 + 1, true, v_now),
(uuid_generate_v4(), v_user_id, v_icat_investasi, 'Dividen saham BBRI', 200000, 'INVESTMENT', v_month_may26 + 5, false, v_now),
(uuid_generate_v4(), v_user_id, v_icat_lainnya, 'Cashback e-wallet', 75000, 'OTHER', v_month_may26 + 18, false, v_now);

-- ============ RECURRING INCOMES ============
INSERT INTO recurring_incomes (id, user_id, category_id, source_name, amount, income_type, recurring_day, is_active, created_at) VALUES
(uuid_generate_v4(), v_user_id, v_icat_gaji, 'Gaji PT Maju Jaya', 8500000, 'SALARY', 1, true, v_now),
(uuid_generate_v4(), v_user_id, v_icat_freelance, 'Retainer klien tetap', 1500000, 'FREELANCE', 15, true, v_now);

-- ============ INSTALLMENTS ============
-- iPhone 15 - 12 bulan, start 3 months ago, 4 payments (payment 4 period = current month)
INSERT INTO installments (id, user_id, name, actual_amount, loan_amount, monthly_payment, tenor, start_date, due_day, status, notes, created_at)
VALUES (v_inst_hp, v_user_id, 'Cicilan iPhone 15', 14999000, 16800000, 1400000, 12, v_hp_start, 15, 'ACTIVE', 'Cicilan HP dari Tokopedia', v_now);

INSERT INTO installment_payments (id, installment_id, payment_number, amount, paid_at, created_at) VALUES
(v_ip_hp_1, v_inst_hp, 1, 1400000, (v_hp_start)::date, v_now),
(v_ip_hp_2, v_inst_hp, 2, 1400000, (v_hp_start + INTERVAL '1 month')::date, v_now),
(v_ip_hp_3, v_inst_hp, 3, 1400000, (v_hp_start + INTERVAL '2 months')::date, v_now),
(v_ip_hp_4, v_inst_hp, 4, 1400000, (v_hp_start + INTERVAL '3 months')::date, v_now);

-- MacBook Air M3 - 24 bulan, start 7 months ago, 8 payments (payment 8 period = current month)
INSERT INTO installments (id, user_id, name, actual_amount, loan_amount, monthly_payment, tenor, start_date, due_day, status, notes, created_at)
VALUES (v_inst_laptop, v_user_id, 'Cicilan MacBook Air M3', 18999000, 21600000, 900000, 24, v_lt_start, 20, 'ACTIVE', 'Cicilan laptop dari iBox', v_now);

INSERT INTO installment_payments (id, installment_id, payment_number, amount, paid_at, created_at) VALUES
(v_ip_lt_1, v_inst_laptop, 1, 900000, (v_lt_start)::date, v_now),
(v_ip_lt_2, v_inst_laptop, 2, 900000, (v_lt_start + INTERVAL '1 month')::date, v_now),
(v_ip_lt_3, v_inst_laptop, 3, 900000, (v_lt_start + INTERVAL '2 months')::date, v_now),
(v_ip_lt_4, v_inst_laptop, 4, 900000, (v_lt_start + INTERVAL '3 months')::date, v_now),
(v_ip_lt_5, v_inst_laptop, 5, 900000, (v_lt_start + INTERVAL '4 months')::date, v_now),
(v_ip_lt_6, v_inst_laptop, 6, 900000, (v_lt_start + INTERVAL '5 months')::date, v_now),
(v_ip_lt_7, v_inst_laptop, 7, 900000, (v_lt_start + INTERVAL '6 months')::date, v_now),
(v_ip_lt_8, v_inst_laptop, 8, 900000, (v_lt_start + INTERVAL '7 months')::date, v_now);

-- ============ INSTALLMENT PAYMENT LEDGER ENTRIES ============
-- Pattern: Debit LIABILITY (reduce debt), Credit CASH (pay out)
-- Transaction date = start_date + (payment_number - 1) months

-- iPhone payment 1 (period: hp_start + 0)
INSERT INTO transactions (id, user_id, transaction_date, description, reference_id, reference_type, created_at)
VALUES (uuid_generate_v4(), v_user_id, v_hp_start, 'Installment Payment: Cicilan iPhone 15', v_ip_hp_1, 'installment_payment', v_now)
RETURNING id INTO v_tx_id;
INSERT INTO transaction_entries (id, transaction_id, account_id, debit, credit, created_at) VALUES
(uuid_generate_v4(), v_tx_id, v_acc_inst_hp, 1400000, 0, v_now),
(uuid_generate_v4(), v_tx_id, v_cash_account_id, 0, 1400000, v_now);

-- iPhone payment 2 (period: hp_start + 1 month)
INSERT INTO transactions (id, user_id, transaction_date, description, reference_id, reference_type, created_at)
VALUES (uuid_generate_v4(), v_user_id, (v_hp_start + INTERVAL '1 month')::date, 'Installment Payment: Cicilan iPhone 15', v_ip_hp_2, 'installment_payment', v_now)
RETURNING id INTO v_tx_id;
INSERT INTO transaction_entries (id, transaction_id, account_id, debit, credit, created_at) VALUES
(uuid_generate_v4(), v_tx_id, v_acc_inst_hp, 1400000, 0, v_now),
(uuid_generate_v4(), v_tx_id, v_cash_account_id, 0, 1400000, v_now);

-- iPhone payment 3 (period: hp_start + 2 months)
INSERT INTO transactions (id, user_id, transaction_date, description, reference_id, reference_type, created_at)
VALUES (uuid_generate_v4(), v_user_id, (v_hp_start + INTERVAL '2 months')::date, 'Installment Payment: Cicilan iPhone 15', v_ip_hp_3, 'installment_payment', v_now)
RETURNING id INTO v_tx_id;
INSERT INTO transaction_entries (id, transaction_id, account_id, debit, credit, created_at) VALUES
(uuid_generate_v4(), v_tx_id, v_acc_inst_hp, 1400000, 0, v_now),
(uuid_generate_v4(), v_tx_id, v_cash_account_id, 0, 1400000, v_now);

-- iPhone payment 4 (period: hp_start + 3 months = CURRENT MONTH)
INSERT INTO transactions (id, user_id, transaction_date, description, reference_id, reference_type, created_at)
VALUES (uuid_generate_v4(), v_user_id, (v_hp_start + INTERVAL '3 months')::date, 'Installment Payment: Cicilan iPhone 15', v_ip_hp_4, 'installment_payment', v_now)
RETURNING id INTO v_tx_id;
INSERT INTO transaction_entries (id, transaction_id, account_id, debit, credit, created_at) VALUES
(uuid_generate_v4(), v_tx_id, v_acc_inst_hp, 1400000, 0, v_now),
(uuid_generate_v4(), v_tx_id, v_cash_account_id, 0, 1400000, v_now);

-- MacBook payment 1
INSERT INTO transactions (id, user_id, transaction_date, description, reference_id, reference_type, created_at)
VALUES (uuid_generate_v4(), v_user_id, v_lt_start, 'Installment Payment: Cicilan MacBook Air M3', v_ip_lt_1, 'installment_payment', v_now)
RETURNING id INTO v_tx_id;
INSERT INTO transaction_entries (id, transaction_id, account_id, debit, credit, created_at) VALUES
(uuid_generate_v4(), v_tx_id, v_acc_inst_laptop, 900000, 0, v_now),
(uuid_generate_v4(), v_tx_id, v_cash_account_id, 0, 900000, v_now);

-- MacBook payment 2
INSERT INTO transactions (id, user_id, transaction_date, description, reference_id, reference_type, created_at)
VALUES (uuid_generate_v4(), v_user_id, (v_lt_start + INTERVAL '1 month')::date, 'Installment Payment: Cicilan MacBook Air M3', v_ip_lt_2, 'installment_payment', v_now)
RETURNING id INTO v_tx_id;
INSERT INTO transaction_entries (id, transaction_id, account_id, debit, credit, created_at) VALUES
(uuid_generate_v4(), v_tx_id, v_acc_inst_laptop, 900000, 0, v_now),
(uuid_generate_v4(), v_tx_id, v_cash_account_id, 0, 900000, v_now);

-- MacBook payment 3
INSERT INTO transactions (id, user_id, transaction_date, description, reference_id, reference_type, created_at)
VALUES (uuid_generate_v4(), v_user_id, (v_lt_start + INTERVAL '2 months')::date, 'Installment Payment: Cicilan MacBook Air M3', v_ip_lt_3, 'installment_payment', v_now)
RETURNING id INTO v_tx_id;
INSERT INTO transaction_entries (id, transaction_id, account_id, debit, credit, created_at) VALUES
(uuid_generate_v4(), v_tx_id, v_acc_inst_laptop, 900000, 0, v_now),
(uuid_generate_v4(), v_tx_id, v_cash_account_id, 0, 900000, v_now);

-- MacBook payment 4
INSERT INTO transactions (id, user_id, transaction_date, description, reference_id, reference_type, created_at)
VALUES (uuid_generate_v4(), v_user_id, (v_lt_start + INTERVAL '3 months')::date, 'Installment Payment: Cicilan MacBook Air M3', v_ip_lt_4, 'installment_payment', v_now)
RETURNING id INTO v_tx_id;
INSERT INTO transaction_entries (id, transaction_id, account_id, debit, credit, created_at) VALUES
(uuid_generate_v4(), v_tx_id, v_acc_inst_laptop, 900000, 0, v_now),
(uuid_generate_v4(), v_tx_id, v_cash_account_id, 0, 900000, v_now);

-- MacBook payment 5
INSERT INTO transactions (id, user_id, transaction_date, description, reference_id, reference_type, created_at)
VALUES (uuid_generate_v4(), v_user_id, (v_lt_start + INTERVAL '4 months')::date, 'Installment Payment: Cicilan MacBook Air M3', v_ip_lt_5, 'installment_payment', v_now)
RETURNING id INTO v_tx_id;
INSERT INTO transaction_entries (id, transaction_id, account_id, debit, credit, created_at) VALUES
(uuid_generate_v4(), v_tx_id, v_acc_inst_laptop, 900000, 0, v_now),
(uuid_generate_v4(), v_tx_id, v_cash_account_id, 0, 900000, v_now);

-- MacBook payment 6
INSERT INTO transactions (id, user_id, transaction_date, description, reference_id, reference_type, created_at)
VALUES (uuid_generate_v4(), v_user_id, (v_lt_start + INTERVAL '5 months')::date, 'Installment Payment: Cicilan MacBook Air M3', v_ip_lt_6, 'installment_payment', v_now)
RETURNING id INTO v_tx_id;
INSERT INTO transaction_entries (id, transaction_id, account_id, debit, credit, created_at) VALUES
(uuid_generate_v4(), v_tx_id, v_acc_inst_laptop, 900000, 0, v_now),
(uuid_generate_v4(), v_tx_id, v_cash_account_id, 0, 900000, v_now);

-- MacBook payment 7
INSERT INTO transactions (id, user_id, transaction_date, description, reference_id, reference_type, created_at)
VALUES (uuid_generate_v4(), v_user_id, (v_lt_start + INTERVAL '6 months')::date, 'Installment Payment: Cicilan MacBook Air M3', v_ip_lt_7, 'installment_payment', v_now)
RETURNING id INTO v_tx_id;
INSERT INTO transaction_entries (id, transaction_id, account_id, debit, credit, created_at) VALUES
(uuid_generate_v4(), v_tx_id, v_acc_inst_laptop, 900000, 0, v_now),
(uuid_generate_v4(), v_tx_id, v_cash_account_id, 0, 900000, v_now);

-- MacBook payment 8 (period: lt_start + 7 months = CURRENT MONTH)
INSERT INTO transactions (id, user_id, transaction_date, description, reference_id, reference_type, created_at)
VALUES (uuid_generate_v4(), v_user_id, (v_lt_start + INTERVAL '7 months')::date, 'Installment Payment: Cicilan MacBook Air M3', v_ip_lt_8, 'installment_payment', v_now)
RETURNING id INTO v_tx_id;
INSERT INTO transaction_entries (id, transaction_id, account_id, debit, credit, created_at) VALUES
(uuid_generate_v4(), v_tx_id, v_acc_inst_laptop, 900000, 0, v_now),
(uuid_generate_v4(), v_tx_id, v_cash_account_id, 0, 900000, v_now);

-- ============ DEBTS ============
INSERT INTO debts (id, user_id, person_name, actual_amount, loan_amount, payment_type, monthly_payment, tenor, due_date, status, notes, created_at)
VALUES (v_debt_teman, v_user_id, 'Budi', 3000000, 3000000, 'INSTALLMENT', 500000, 6, (v_today + INTERVAL '4 months')::date, 'ACTIVE', 'Pinjam untuk bayar kos', v_now);

INSERT INTO debt_payments (id, debt_id, payment_number, amount, paid_at, created_at) VALUES
(v_dp_budi_1, v_debt_teman, 1, 500000, (v_last_month + 15)::date, v_now),
(v_dp_budi_2, v_debt_teman, 2, 500000, (v_this_month + 1)::date, v_now);

INSERT INTO debts (id, user_id, person_name, actual_amount, payment_type, status, notes, created_at)
VALUES (v_debt_keluarga, v_user_id, 'Kakak', 1500000, 'ONE_TIME', 'COMPLETED', 'Pinjam darurat bulan lalu', v_now);

INSERT INTO debt_payments (id, debt_id, payment_number, amount, paid_at, created_at) VALUES
(v_dp_kakak_1, v_debt_keluarga, 1, 1500000, (v_last_month + 20)::date, v_now);

-- ============ DEBT PAYMENT LEDGER ENTRIES ============
-- Pattern: Debit LIABILITY (reduce debt), Credit CASH (pay out)
-- Transaction date = paid_at

-- Budi payment 1 (last month)
INSERT INTO transactions (id, user_id, transaction_date, description, reference_id, reference_type, created_at)
VALUES (uuid_generate_v4(), v_user_id, (v_last_month + 15)::date, 'Debt Payment: Budi', v_dp_budi_1, 'debt_payment', v_now)
RETURNING id INTO v_tx_id;
INSERT INTO transaction_entries (id, transaction_id, account_id, debit, credit, created_at) VALUES
(uuid_generate_v4(), v_tx_id, v_acc_debt_teman, 500000, 0, v_now),
(uuid_generate_v4(), v_tx_id, v_cash_account_id, 0, 500000, v_now);

-- Budi payment 2 (this month)
INSERT INTO transactions (id, user_id, transaction_date, description, reference_id, reference_type, created_at)
VALUES (uuid_generate_v4(), v_user_id, (v_this_month + 1)::date, 'Debt Payment: Budi', v_dp_budi_2, 'debt_payment', v_now)
RETURNING id INTO v_tx_id;
INSERT INTO transaction_entries (id, transaction_id, account_id, debit, credit, created_at) VALUES
(uuid_generate_v4(), v_tx_id, v_acc_debt_teman, 500000, 0, v_now),
(uuid_generate_v4(), v_tx_id, v_cash_account_id, 0, 500000, v_now);

-- Kakak payment 1 (last month)
INSERT INTO transactions (id, user_id, transaction_date, description, reference_id, reference_type, created_at)
VALUES (uuid_generate_v4(), v_user_id, (v_last_month + 20)::date, 'Debt Payment: Kakak', v_dp_kakak_1, 'debt_payment', v_now)
RETURNING id INTO v_tx_id;
INSERT INTO transaction_entries (id, transaction_id, account_id, debit, credit, created_at) VALUES
(uuid_generate_v4(), v_tx_id, v_acc_debt_keluarga, 1500000, 0, v_now),
(uuid_generate_v4(), v_tx_id, v_cash_account_id, 0, 1500000, v_now);

-- ============ SAVINGS GOALS ============
INSERT INTO savings_goals (id, user_id, name, target_amount, current_amount, target_date, icon, status, notes, created_at)
VALUES (v_save_liburan, v_user_id, 'Liburan Bali', 5000000, 3200000, (v_today + INTERVAL '3 months')::date, 'plane', 'ACTIVE', 'Liburan akhir tahun ke Bali', v_now);

INSERT INTO savings_contributions (id, savings_goal_id, amount, contribution_date, created_at) VALUES
(v_sc_lib_1, v_save_liburan, 1000000, (v_today - INTERVAL '2 months')::date, v_now),
(v_sc_lib_2, v_save_liburan, 1200000, (v_today - INTERVAL '1 month')::date, v_now),
(v_sc_lib_3, v_save_liburan, 1000000, (v_this_month + 1)::date, v_now);

INSERT INTO savings_goals (id, user_id, name, target_amount, current_amount, target_date, icon, status, notes, created_at)
VALUES (v_save_darurat, v_user_id, 'Dana Darurat', 15000000, 12000000, (v_today + INTERVAL '6 months')::date, 'shield', 'ACTIVE', 'Target 3x pengeluaran bulanan', v_now);

INSERT INTO savings_contributions (id, savings_goal_id, amount, contribution_date, created_at) VALUES
(v_sc_dar_1, v_save_darurat, 3000000, (v_today - INTERVAL '3 months')::date, v_now),
(v_sc_dar_2, v_save_darurat, 3000000, (v_today - INTERVAL '2 months')::date, v_now),
(v_sc_dar_3, v_save_darurat, 3000000, (v_today - INTERVAL '1 month')::date, v_now),
(v_sc_dar_4, v_save_darurat, 3000000, (v_this_month + 1)::date, v_now);

-- ============ SAVINGS CONTRIBUTION LEDGER ENTRIES ============
-- Pattern: Debit ASSET/savings (increase savings), Credit CASH (pay out)
-- Transaction date = contribution_date

-- Liburan 1
INSERT INTO transactions (id, user_id, transaction_date, description, reference_id, reference_type, created_at)
VALUES (uuid_generate_v4(), v_user_id, (v_today - INTERVAL '2 months')::date, 'Savings Contribution: Liburan Bali', v_sc_lib_1, 'savings_contribution', v_now)
RETURNING id INTO v_tx_id;
INSERT INTO transaction_entries (id, transaction_id, account_id, debit, credit, created_at) VALUES
(uuid_generate_v4(), v_tx_id, v_acc_save_liburan, 1000000, 0, v_now),
(uuid_generate_v4(), v_tx_id, v_cash_account_id, 0, 1000000, v_now);

-- Liburan 2
INSERT INTO transactions (id, user_id, transaction_date, description, reference_id, reference_type, created_at)
VALUES (uuid_generate_v4(), v_user_id, (v_today - INTERVAL '1 month')::date, 'Savings Contribution: Liburan Bali', v_sc_lib_2, 'savings_contribution', v_now)
RETURNING id INTO v_tx_id;
INSERT INTO transaction_entries (id, transaction_id, account_id, debit, credit, created_at) VALUES
(uuid_generate_v4(), v_tx_id, v_acc_save_liburan, 1200000, 0, v_now),
(uuid_generate_v4(), v_tx_id, v_cash_account_id, 0, 1200000, v_now);

-- Liburan 3 (this month)
INSERT INTO transactions (id, user_id, transaction_date, description, reference_id, reference_type, created_at)
VALUES (uuid_generate_v4(), v_user_id, (v_this_month + 1)::date, 'Savings Contribution: Liburan Bali', v_sc_lib_3, 'savings_contribution', v_now)
RETURNING id INTO v_tx_id;
INSERT INTO transaction_entries (id, transaction_id, account_id, debit, credit, created_at) VALUES
(uuid_generate_v4(), v_tx_id, v_acc_save_liburan, 1000000, 0, v_now),
(uuid_generate_v4(), v_tx_id, v_cash_account_id, 0, 1000000, v_now);

-- Darurat 1
INSERT INTO transactions (id, user_id, transaction_date, description, reference_id, reference_type, created_at)
VALUES (uuid_generate_v4(), v_user_id, (v_today - INTERVAL '3 months')::date, 'Savings Contribution: Dana Darurat', v_sc_dar_1, 'savings_contribution', v_now)
RETURNING id INTO v_tx_id;
INSERT INTO transaction_entries (id, transaction_id, account_id, debit, credit, created_at) VALUES
(uuid_generate_v4(), v_tx_id, v_acc_save_darurat, 3000000, 0, v_now),
(uuid_generate_v4(), v_tx_id, v_cash_account_id, 0, 3000000, v_now);

-- Darurat 2
INSERT INTO transactions (id, user_id, transaction_date, description, reference_id, reference_type, created_at)
VALUES (uuid_generate_v4(), v_user_id, (v_today - INTERVAL '2 months')::date, 'Savings Contribution: Dana Darurat', v_sc_dar_2, 'savings_contribution', v_now)
RETURNING id INTO v_tx_id;
INSERT INTO transaction_entries (id, transaction_id, account_id, debit, credit, created_at) VALUES
(uuid_generate_v4(), v_tx_id, v_acc_save_darurat, 3000000, 0, v_now),
(uuid_generate_v4(), v_tx_id, v_cash_account_id, 0, 3000000, v_now);

-- Darurat 3
INSERT INTO transactions (id, user_id, transaction_date, description, reference_id, reference_type, created_at)
VALUES (uuid_generate_v4(), v_user_id, (v_today - INTERVAL '1 month')::date, 'Savings Contribution: Dana Darurat', v_sc_dar_3, 'savings_contribution', v_now)
RETURNING id INTO v_tx_id;
INSERT INTO transaction_entries (id, transaction_id, account_id, debit, credit, created_at) VALUES
(uuid_generate_v4(), v_tx_id, v_acc_save_darurat, 3000000, 0, v_now),
(uuid_generate_v4(), v_tx_id, v_cash_account_id, 0, 3000000, v_now);

-- Darurat 4 (this month)
INSERT INTO transactions (id, user_id, transaction_date, description, reference_id, reference_type, created_at)
VALUES (uuid_generate_v4(), v_user_id, (v_this_month + 1)::date, 'Savings Contribution: Dana Darurat', v_sc_dar_4, 'savings_contribution', v_now)
RETURNING id INTO v_tx_id;
INSERT INTO transaction_entries (id, transaction_id, account_id, debit, credit, created_at) VALUES
(uuid_generate_v4(), v_tx_id, v_acc_save_darurat, 3000000, 0, v_now),
(uuid_generate_v4(), v_tx_id, v_cash_account_id, 0, 3000000, v_now);

-- ============ EXPENSE TEMPLATE ============
INSERT INTO expense_template_groups (id, user_id, name, recurring_day, notes, created_at)
VALUES (v_template_bulanan, v_user_id, 'Pengeluaran Rutin Bulanan', 1, 'Template untuk tagihan bulanan', v_now);

INSERT INTO expense_template_items (id, group_id, category_id, item_name, unit_price, quantity, created_at) VALUES
(uuid_generate_v4(), v_template_bulanan, v_cat_tagihan, 'Listrik PLN', 350000, 1, v_now),
(uuid_generate_v4(), v_template_bulanan, v_cat_tagihan, 'Internet IndiHome', 399000, 1, v_now),
(uuid_generate_v4(), v_template_bulanan, v_cat_hiburan, 'Langganan Spotify', 54990, 1, v_now),
(uuid_generate_v4(), v_template_bulanan, v_cat_hiburan, 'Langganan Netflix', 54000, 1, v_now);

-- ============ UPDATE ACCOUNT BALANCES ============
-- Cash: set to a reasonable positive balance
UPDATE accounts SET current_balance = 4500000 WHERE id = v_cash_account_id;
-- Installment liability (remaining owed): iPhone (12-4)*1400000=11200000, MacBook (24-8)*900000=14400000
UPDATE accounts SET current_balance = 11200000 WHERE id = v_acc_inst_hp;
UPDATE accounts SET current_balance = 14400000 WHERE id = v_acc_inst_laptop;
-- Debt liability: Budi 3000000-(500000*2)=2000000, Kakak completed=0
UPDATE accounts SET current_balance = 2000000 WHERE id = v_acc_debt_teman;
UPDATE accounts SET current_balance = 0 WHERE id = v_acc_debt_keluarga;
-- Savings: Liburan 3200000, Darurat 12000000
UPDATE accounts SET current_balance = 3200000 WHERE id = v_acc_save_liburan;
UPDATE accounts SET current_balance = 12000000 WHERE id = v_acc_save_darurat;

END;
$$;

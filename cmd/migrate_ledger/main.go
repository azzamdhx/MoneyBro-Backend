package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"github.com/azzamdhx/moneybro/backend/internal/config"
	"github.com/azzamdhx/moneybro/backend/internal/database"
	"github.com/azzamdhx/moneybro/backend/internal/models"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := config.Load()
	db, err := database.NewPostgres(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Starting ledger migration...")

	// Step 1: Create default accounts for all existing users
	if err := migrateUserAccounts(db); err != nil {
		log.Fatalf("Failed to migrate user accounts: %v", err)
	}

	// Step 2: Create expense accounts for all existing categories
	if err := migrateCategoryAccounts(db); err != nil {
		log.Fatalf("Failed to migrate category accounts: %v", err)
	}

	// Step 3: Create income accounts for all existing income categories
	if err := migrateIncomeCategoryAccounts(db); err != nil {
		log.Fatalf("Failed to migrate income category accounts: %v", err)
	}

	// Step 4: Create liability accounts for all existing installments
	if err := migrateInstallmentAccounts(db); err != nil {
		log.Fatalf("Failed to migrate installment accounts: %v", err)
	}

	// Step 5: Create liability accounts for all existing debts
	if err := migrateDebtAccounts(db); err != nil {
		log.Fatalf("Failed to migrate debt accounts: %v", err)
	}

	// Step 6: Migrate existing expenses to ledger entries
	if err := migrateExpenses(db); err != nil {
		log.Fatalf("Failed to migrate expenses: %v", err)
	}

	// Step 7: Migrate existing incomes to ledger entries
	if err := migrateIncomes(db); err != nil {
		log.Fatalf("Failed to migrate incomes: %v", err)
	}

	// Step 8: Migrate existing installment payments to ledger entries
	if err := migrateInstallmentPayments(db); err != nil {
		log.Fatalf("Failed to migrate installment payments: %v", err)
	}

	// Step 9: Migrate existing debt payments to ledger entries
	if err := migrateDebtPayments(db); err != nil {
		log.Fatalf("Failed to migrate debt payments: %v", err)
	}

	// Step 10: Recalculate all account balances
	if err := recalculateBalances(db); err != nil {
		log.Fatalf("Failed to recalculate balances: %v", err)
	}

	log.Println("Ledger migration completed successfully!")
}

func migrateUserAccounts(db *gorm.DB) error {
	log.Println("Migrating user accounts...")
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		return err
	}

	for _, user := range users {
		var existing models.Account
		if err := db.Where("user_id = ? AND is_default = ?", user.ID, true).First(&existing).Error; err == nil {
			continue // Already has default account
		}

		account := models.Account{
			ID:          uuid.New(),
			UserID:      user.ID,
			Name:        "Cash",
			AccountType: models.AccountTypeAsset,
			IsDefault:   true,
		}
		if err := db.Create(&account).Error; err != nil {
			return fmt.Errorf("failed to create default account for user %s: %w", user.ID, err)
		}
		log.Printf("Created default account for user %s", user.ID)
	}
	return nil
}

func migrateCategoryAccounts(db *gorm.DB) error {
	log.Println("Migrating category accounts...")
	var categories []models.Category
	if err := db.Find(&categories).Error; err != nil {
		return err
	}

	for _, cat := range categories {
		var existing models.Account
		if err := db.Where("reference_id = ? AND reference_type = ?", cat.ID, "category").First(&existing).Error; err == nil {
			continue
		}

		account := models.Account{
			ID:            uuid.New(),
			UserID:        cat.UserID,
			Name:          cat.Name,
			AccountType:   models.AccountTypeExpense,
			IsDefault:     false,
			ReferenceID:   &cat.ID,
			ReferenceType: strPtr("category"),
		}
		if err := db.Create(&account).Error; err != nil {
			return fmt.Errorf("failed to create expense account for category %s: %w", cat.ID, err)
		}
		log.Printf("Created expense account for category %s", cat.Name)
	}
	return nil
}

func migrateIncomeCategoryAccounts(db *gorm.DB) error {
	log.Println("Migrating income category accounts...")
	var categories []models.IncomeCategory
	if err := db.Find(&categories).Error; err != nil {
		return err
	}

	for _, cat := range categories {
		var existing models.Account
		if err := db.Where("reference_id = ? AND reference_type = ?", cat.ID, "income_category").First(&existing).Error; err == nil {
			continue
		}

		account := models.Account{
			ID:            uuid.New(),
			UserID:        cat.UserID,
			Name:          cat.Name,
			AccountType:   models.AccountTypeIncome,
			IsDefault:     false,
			ReferenceID:   &cat.ID,
			ReferenceType: strPtr("income_category"),
		}
		if err := db.Create(&account).Error; err != nil {
			return fmt.Errorf("failed to create income account for category %s: %w", cat.ID, err)
		}
		log.Printf("Created income account for category %s", cat.Name)
	}
	return nil
}

func migrateInstallmentAccounts(db *gorm.DB) error {
	log.Println("Migrating installment accounts...")
	var installments []models.Installment
	if err := db.Find(&installments).Error; err != nil {
		return err
	}

	for _, inst := range installments {
		var existing models.Account
		if err := db.Where("reference_id = ? AND reference_type = ?", inst.ID, "installment").First(&existing).Error; err == nil {
			continue
		}

		account := models.Account{
			ID:            uuid.New(),
			UserID:        inst.UserID,
			Name:          inst.Name,
			AccountType:   models.AccountTypeLiability,
			IsDefault:     false,
			ReferenceID:   &inst.ID,
			ReferenceType: strPtr("installment"),
		}
		if err := db.Create(&account).Error; err != nil {
			return fmt.Errorf("failed to create liability account for installment %s: %w", inst.ID, err)
		}
		log.Printf("Created liability account for installment %s", inst.Name)
	}
	return nil
}

func migrateDebtAccounts(db *gorm.DB) error {
	log.Println("Migrating debt accounts...")
	var debts []models.Debt
	if err := db.Find(&debts).Error; err != nil {
		return err
	}

	for _, debt := range debts {
		var existing models.Account
		if err := db.Where("reference_id = ? AND reference_type = ?", debt.ID, "debt").First(&existing).Error; err == nil {
			continue
		}

		account := models.Account{
			ID:            uuid.New(),
			UserID:        debt.UserID,
			Name:          "Hutang: " + debt.PersonName,
			AccountType:   models.AccountTypeLiability,
			IsDefault:     false,
			ReferenceID:   &debt.ID,
			ReferenceType: strPtr("debt"),
		}
		if err := db.Create(&account).Error; err != nil {
			return fmt.Errorf("failed to create liability account for debt %s: %w", debt.ID, err)
		}
		log.Printf("Created liability account for debt %s", debt.PersonName)
	}
	return nil
}

func migrateExpenses(db *gorm.DB) error {
	log.Println("Migrating expenses to ledger...")
	var expenses []models.Expense
	if err := db.Find(&expenses).Error; err != nil {
		return err
	}

	for _, exp := range expenses {
		var existing models.Transaction
		if err := db.Where("reference_id = ? AND reference_type = ?", exp.ID, "expense").First(&existing).Error; err == nil {
			continue
		}

		var expenseAccount models.Account
		if err := db.Where("reference_id = ? AND reference_type = ?", exp.CategoryID, "category").First(&expenseAccount).Error; err != nil {
			log.Printf("Warning: expense account not found for category %s, skipping expense %s", exp.CategoryID, exp.ID)
			continue
		}

		var cashAccount models.Account
		if err := db.Where("user_id = ? AND is_default = ?", exp.UserID, true).First(&cashAccount).Error; err != nil {
			log.Printf("Warning: cash account not found for user %s, skipping expense %s", exp.UserID, exp.ID)
			continue
		}

		txDate := time.Now()
		if exp.ExpenseDate != nil {
			txDate = *exp.ExpenseDate
		}

		tx := models.Transaction{
			ID:              uuid.New(),
			UserID:          exp.UserID,
			TransactionDate: txDate,
			Description:     "Expense: " + exp.ItemName,
			ReferenceID:     &exp.ID,
			ReferenceType:   strPtr("expense"),
		}
		if err := db.Create(&tx).Error; err != nil {
			return fmt.Errorf("failed to create transaction for expense %s: %w", exp.ID, err)
		}

		amount := exp.Total()
		entries := []models.TransactionEntry{
			{ID: uuid.New(), TransactionID: tx.ID, AccountID: expenseAccount.ID, Debit: amount, Credit: 0},
			{ID: uuid.New(), TransactionID: tx.ID, AccountID: cashAccount.ID, Debit: 0, Credit: amount},
		}
		for _, entry := range entries {
			if err := db.Create(&entry).Error; err != nil {
				return fmt.Errorf("failed to create transaction entry for expense %s: %w", exp.ID, err)
			}
		}
		log.Printf("Migrated expense %s", exp.ItemName)
	}
	return nil
}

func migrateIncomes(db *gorm.DB) error {
	log.Println("Migrating incomes to ledger...")
	var incomes []models.Income
	if err := db.Find(&incomes).Error; err != nil {
		return err
	}

	for _, inc := range incomes {
		var existing models.Transaction
		if err := db.Where("reference_id = ? AND reference_type = ?", inc.ID, "income").First(&existing).Error; err == nil {
			continue
		}

		var incomeAccount models.Account
		if err := db.Where("reference_id = ? AND reference_type = ?", inc.CategoryID, "income_category").First(&incomeAccount).Error; err != nil {
			log.Printf("Warning: income account not found for category %s, skipping income %s", inc.CategoryID, inc.ID)
			continue
		}

		var cashAccount models.Account
		if err := db.Where("user_id = ? AND is_default = ?", inc.UserID, true).First(&cashAccount).Error; err != nil {
			log.Printf("Warning: cash account not found for user %s, skipping income %s", inc.UserID, inc.ID)
			continue
		}

		tx := models.Transaction{
			ID:              uuid.New(),
			UserID:          inc.UserID,
			TransactionDate: inc.IncomeDate,
			Description:     "Income: " + inc.SourceName,
			ReferenceID:     &inc.ID,
			ReferenceType:   strPtr("income"),
		}
		if err := db.Create(&tx).Error; err != nil {
			return fmt.Errorf("failed to create transaction for income %s: %w", inc.ID, err)
		}

		entries := []models.TransactionEntry{
			{ID: uuid.New(), TransactionID: tx.ID, AccountID: cashAccount.ID, Debit: inc.Amount, Credit: 0},
			{ID: uuid.New(), TransactionID: tx.ID, AccountID: incomeAccount.ID, Debit: 0, Credit: inc.Amount},
		}
		for _, entry := range entries {
			if err := db.Create(&entry).Error; err != nil {
				return fmt.Errorf("failed to create transaction entry for income %s: %w", inc.ID, err)
			}
		}
		log.Printf("Migrated income %s", inc.SourceName)
	}
	return nil
}

func migrateInstallmentPayments(db *gorm.DB) error {
	log.Println("Migrating installment payments to ledger...")
	var payments []models.InstallmentPayment
	if err := db.Preload("Installment").Find(&payments).Error; err != nil {
		return err
	}

	for _, payment := range payments {
		var existing models.Transaction
		if err := db.Where("reference_id = ? AND reference_type = ?", payment.ID, "installment_payment").First(&existing).Error; err == nil {
			continue
		}

		var liabilityAccount models.Account
		if err := db.Where("reference_id = ? AND reference_type = ?", payment.InstallmentID, "installment").First(&liabilityAccount).Error; err != nil {
			log.Printf("Warning: liability account not found for installment %s, skipping payment %s", payment.InstallmentID, payment.ID)
			continue
		}

		var cashAccount models.Account
		if err := db.Where("user_id = ? AND is_default = ?", payment.Installment.UserID, true).First(&cashAccount).Error; err != nil {
			log.Printf("Warning: cash account not found, skipping payment %s", payment.ID)
			continue
		}

		tx := models.Transaction{
			ID:              uuid.New(),
			UserID:          payment.Installment.UserID,
			TransactionDate: payment.PaidAt,
			Description:     fmt.Sprintf("Installment Payment: %s #%d", payment.Installment.Name, payment.PaymentNumber),
			ReferenceID:     &payment.ID,
			ReferenceType:   strPtr("installment_payment"),
		}
		if err := db.Create(&tx).Error; err != nil {
			return fmt.Errorf("failed to create transaction for installment payment %s: %w", payment.ID, err)
		}

		entries := []models.TransactionEntry{
			{ID: uuid.New(), TransactionID: tx.ID, AccountID: liabilityAccount.ID, Debit: payment.Amount, Credit: 0},
			{ID: uuid.New(), TransactionID: tx.ID, AccountID: cashAccount.ID, Debit: 0, Credit: payment.Amount},
		}
		for _, entry := range entries {
			if err := db.Create(&entry).Error; err != nil {
				return fmt.Errorf("failed to create transaction entry for installment payment %s: %w", payment.ID, err)
			}
		}
		log.Printf("Migrated installment payment #%d for %s", payment.PaymentNumber, payment.Installment.Name)
	}
	return nil
}

func migrateDebtPayments(db *gorm.DB) error {
	log.Println("Migrating debt payments to ledger...")
	var payments []models.DebtPayment
	if err := db.Preload("Debt").Find(&payments).Error; err != nil {
		return err
	}

	for _, payment := range payments {
		var existing models.Transaction
		if err := db.Where("reference_id = ? AND reference_type = ?", payment.ID, "debt_payment").First(&existing).Error; err == nil {
			continue
		}

		var liabilityAccount models.Account
		if err := db.Where("reference_id = ? AND reference_type = ?", payment.DebtID, "debt").First(&liabilityAccount).Error; err != nil {
			log.Printf("Warning: liability account not found for debt %s, skipping payment %s", payment.DebtID, payment.ID)
			continue
		}

		var cashAccount models.Account
		if err := db.Where("user_id = ? AND is_default = ?", payment.Debt.UserID, true).First(&cashAccount).Error; err != nil {
			log.Printf("Warning: cash account not found, skipping payment %s", payment.ID)
			continue
		}

		tx := models.Transaction{
			ID:              uuid.New(),
			UserID:          payment.Debt.UserID,
			TransactionDate: payment.PaidAt,
			Description:     fmt.Sprintf("Debt Payment: %s #%d", payment.Debt.PersonName, payment.PaymentNumber),
			ReferenceID:     &payment.ID,
			ReferenceType:   strPtr("debt_payment"),
		}
		if err := db.Create(&tx).Error; err != nil {
			return fmt.Errorf("failed to create transaction for debt payment %s: %w", payment.ID, err)
		}

		entries := []models.TransactionEntry{
			{ID: uuid.New(), TransactionID: tx.ID, AccountID: liabilityAccount.ID, Debit: payment.Amount, Credit: 0},
			{ID: uuid.New(), TransactionID: tx.ID, AccountID: cashAccount.ID, Debit: 0, Credit: payment.Amount},
		}
		for _, entry := range entries {
			if err := db.Create(&entry).Error; err != nil {
				return fmt.Errorf("failed to create transaction entry for debt payment %s: %w", payment.ID, err)
			}
		}
		log.Printf("Migrated debt payment #%d for %s", payment.PaymentNumber, payment.Debt.PersonName)
	}
	return nil
}

func recalculateBalances(db *gorm.DB) error {
	log.Println("Recalculating account balances...")
	var accounts []models.Account
	if err := db.Find(&accounts).Error; err != nil {
		return err
	}

	for _, account := range accounts {
		var debitSum, creditSum int64

		db.Model(&models.TransactionEntry{}).
			Where("account_id = ?", account.ID).
			Select("COALESCE(SUM(debit), 0)").
			Scan(&debitSum)

		db.Model(&models.TransactionEntry{}).
			Where("account_id = ?", account.ID).
			Select("COALESCE(SUM(credit), 0)").
			Scan(&creditSum)

		var balance int64
		switch account.AccountType {
		case models.AccountTypeAsset, models.AccountTypeExpense:
			balance = debitSum - creditSum
		case models.AccountTypeLiability, models.AccountTypeIncome:
			balance = creditSum - debitSum
		}

		if err := db.Model(&account).Update("current_balance", balance).Error; err != nil {
			return fmt.Errorf("failed to update balance for account %s: %w", account.ID, err)
		}
		log.Printf("Updated balance for account %s: %d", account.Name, balance)
	}
	return nil
}

func strPtr(s string) *string {
	return &s
}

package services

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/azzamdhx/moneybro/backend/internal/models"
	"github.com/azzamdhx/moneybro/backend/internal/repository"
)

type ExpenseService struct {
	expenseRepo   repository.ExpenseRepository
	categoryRepo  repository.CategoryRepository
	accountRepo   repository.AccountRepository
	ledgerService *LedgerService
}

func NewExpenseService(
	expenseRepo repository.ExpenseRepository,
	categoryRepo repository.CategoryRepository,
	accountRepo repository.AccountRepository,
	ledgerService *LedgerService,
) *ExpenseService {
	return &ExpenseService{
		expenseRepo:   expenseRepo,
		categoryRepo:  categoryRepo,
		accountRepo:   accountRepo,
		ledgerService: ledgerService,
	}
}

type CreateExpenseInput struct {
	CategoryID  uuid.UUID
	ItemName    string
	UnitPrice   int64
	Quantity    int
	Notes       *string
	ExpenseDate *time.Time
}

func (s *ExpenseService) Create(userID uuid.UUID, input CreateExpenseInput) (*models.Expense, error) {
	if input.ItemName == "" {
		return nil, errors.New("item name is required")
	}
	if input.UnitPrice <= 0 {
		return nil, errors.New("unit price must be positive")
	}
	if input.Quantity <= 0 {
		input.Quantity = 1
	}

	expense := &models.Expense{
		ID:          uuid.New(),
		UserID:      userID,
		CategoryID:  input.CategoryID,
		ItemName:    input.ItemName,
		UnitPrice:   input.UnitPrice,
		Quantity:    input.Quantity,
		Notes:       input.Notes,
		ExpenseDate: input.ExpenseDate,
	}

	if err := s.expenseRepo.Create(expense); err != nil {
		return nil, err
	}

	// Create ledger entry: DEBIT Expense Account, CREDIT Cash Account
	if err := s.createLedgerEntry(userID, expense); err != nil {
		return nil, err
	}

	return s.expenseRepo.GetByID(expense.ID)
}

func (s *ExpenseService) GetByID(id uuid.UUID) (*models.Expense, error) {
	return s.expenseRepo.GetByID(id)
}

func (s *ExpenseService) GetByUserID(userID uuid.UUID, filter *repository.ExpenseFilter) ([]models.Expense, error) {
	return s.expenseRepo.GetByUserID(userID, filter)
}

type UpdateExpenseInput struct {
	CategoryID  *uuid.UUID
	ItemName    *string
	UnitPrice   *int64
	Quantity    *int
	Notes       *string
	ExpenseDate *time.Time
}

func (s *ExpenseService) Update(id uuid.UUID, input UpdateExpenseInput) (*models.Expense, error) {
	expense, err := s.expenseRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if input.CategoryID != nil {
		expense.CategoryID = *input.CategoryID
	}
	if input.ItemName != nil {
		expense.ItemName = *input.ItemName
	}
	if input.UnitPrice != nil {
		expense.UnitPrice = *input.UnitPrice
	}
	if input.Quantity != nil {
		expense.Quantity = *input.Quantity
	}
	if input.Notes != nil {
		expense.Notes = input.Notes
	}
	if input.ExpenseDate != nil {
		expense.ExpenseDate = input.ExpenseDate
	}

	if err := s.expenseRepo.Update(expense); err != nil {
		return nil, err
	}

	// Update ledger entry
	if err := s.updateLedgerEntry(expense); err != nil {
		return nil, err
	}

	return s.expenseRepo.GetByID(expense.ID)
}

func (s *ExpenseService) Delete(id uuid.UUID) error {
	expense, err := s.expenseRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Delete ledger entry first
	if err := s.ledgerService.DeleteByReference(expense.ID, "expense"); err != nil {
		return err
	}

	return s.expenseRepo.Delete(id)
}

func (s *ExpenseService) createLedgerEntry(userID uuid.UUID, expense *models.Expense) error {
	// Get expense account (linked to category)
	expenseAccount, err := s.accountRepo.GetByReference(expense.CategoryID, "category")
	if err != nil {
		return err
	}

	// Get default cash account
	cashAccount, err := s.accountRepo.GetDefaultByUserID(userID)
	if err != nil {
		return err
	}

	amount := expense.Total()
	expenseDate := time.Now()
	if expense.ExpenseDate != nil {
		expenseDate = *expense.ExpenseDate
	}

	entries := []LedgerEntry{
		{AccountID: expenseAccount.ID, Debit: amount, Credit: 0},
		{AccountID: cashAccount.ID, Debit: 0, Credit: amount},
	}

	_, err = s.ledgerService.CreateJournalEntry(
		userID,
		expenseDate,
		"Expense: "+expense.ItemName,
		entries,
		&expense.ID,
		"expense",
	)
	return err
}

func (s *ExpenseService) updateLedgerEntry(expense *models.Expense) error {
	tx, err := s.ledgerService.GetTransactionByReference(expense.ID, "expense")
	if err != nil {
		return err
	}

	expenseAccount, err := s.accountRepo.GetByReference(expense.CategoryID, "category")
	if err != nil {
		return err
	}

	cashAccount, err := s.accountRepo.GetDefaultByUserID(expense.UserID)
	if err != nil {
		return err
	}

	amount := expense.Total()
	expenseDate := time.Now()
	if expense.ExpenseDate != nil {
		expenseDate = *expense.ExpenseDate
	}

	entries := []LedgerEntry{
		{AccountID: expenseAccount.ID, Debit: amount, Credit: 0},
		{AccountID: cashAccount.ID, Debit: 0, Credit: amount},
	}

	_, err = s.ledgerService.UpdateJournalEntry(tx.ID, expenseDate, "Expense: "+expense.ItemName, entries)
	return err
}

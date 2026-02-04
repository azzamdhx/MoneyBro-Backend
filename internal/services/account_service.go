package services

import (
	"errors"

	"github.com/google/uuid"

	"github.com/azzamdhx/moneybro/backend/internal/models"
	"github.com/azzamdhx/moneybro/backend/internal/repository"
)

type AccountService struct {
	accountRepo repository.AccountRepository
}

func NewAccountService(accountRepo repository.AccountRepository) *AccountService {
	return &AccountService{accountRepo: accountRepo}
}

func (s *AccountService) CreateAccount(userID uuid.UUID, name string, accountType models.AccountType, isDefault bool) (*models.Account, error) {
	account := &models.Account{
		UserID:      userID,
		Name:        name,
		AccountType: accountType,
		IsDefault:   isDefault,
	}
	if err := s.accountRepo.Create(account); err != nil {
		return nil, err
	}
	return account, nil
}

func (s *AccountService) CreateDefaultAccount(userID uuid.UUID) (*models.Account, error) {
	return s.CreateAccount(userID, "Cash", models.AccountTypeAsset, true)
}

func (s *AccountService) CreateLinkedAccount(userID uuid.UUID, name string, accountType models.AccountType, referenceID uuid.UUID, referenceType string) (*models.Account, error) {
	account := &models.Account{
		UserID:        userID,
		Name:          name,
		AccountType:   accountType,
		IsDefault:     false,
		ReferenceID:   &referenceID,
		ReferenceType: &referenceType,
	}
	if err := s.accountRepo.Create(account); err != nil {
		return nil, err
	}
	return account, nil
}

func (s *AccountService) GetAccounts(userID uuid.UUID) ([]models.Account, error) {
	return s.accountRepo.GetByUserID(userID)
}

func (s *AccountService) GetAccountsByType(userID uuid.UUID, accountType models.AccountType) ([]models.Account, error) {
	return s.accountRepo.GetByUserIDAndType(userID, accountType)
}

func (s *AccountService) GetAccount(id uuid.UUID) (*models.Account, error) {
	return s.accountRepo.GetByID(id)
}

func (s *AccountService) GetDefaultAccount(userID uuid.UUID) (*models.Account, error) {
	return s.accountRepo.GetDefaultByUserID(userID)
}

func (s *AccountService) GetAccountByReference(referenceID uuid.UUID, referenceType string) (*models.Account, error) {
	return s.accountRepo.GetByReference(referenceID, referenceType)
}

func (s *AccountService) UpdateAccount(id uuid.UUID, name string) (*models.Account, error) {
	account, err := s.accountRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	account.Name = name
	if err := s.accountRepo.Update(account); err != nil {
		return nil, err
	}
	return account, nil
}

func (s *AccountService) DeleteAccount(id uuid.UUID) error {
	account, err := s.accountRepo.GetByID(id)
	if err != nil {
		return err
	}
	if account.IsDefault {
		return errors.New("cannot delete default account")
	}
	return s.accountRepo.Delete(id)
}

func (s *AccountService) DeleteAccountByReference(referenceID uuid.UUID, referenceType string) error {
	return s.accountRepo.DeleteByReference(referenceID, referenceType)
}

func (s *AccountService) RecalculateBalance(accountID uuid.UUID, entryRepo repository.TransactionEntryRepository) error {
	debit, credit, err := entryRepo.SumByAccountID(accountID)
	if err != nil {
		return err
	}

	account, err := s.accountRepo.GetByID(accountID)
	if err != nil {
		return err
	}

	var balance int64
	switch account.AccountType {
	case models.AccountTypeAsset, models.AccountTypeExpense:
		balance = debit - credit
	case models.AccountTypeLiability, models.AccountTypeIncome:
		balance = credit - debit
	}

	return s.accountRepo.UpdateBalance(accountID, balance)
}

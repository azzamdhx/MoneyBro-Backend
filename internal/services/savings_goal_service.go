package services

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/azzamdhx/moneybro/backend/internal/models"
	"github.com/azzamdhx/moneybro/backend/internal/repository"
)

type SavingsGoalService struct {
	goalRepo         repository.SavingsGoalRepository
	contributionRepo repository.SavingsContributionRepository
	accountRepo      repository.AccountRepository
	accountService   *AccountService
	ledgerService    *LedgerService
}

func NewSavingsGoalService(
	goalRepo repository.SavingsGoalRepository,
	contributionRepo repository.SavingsContributionRepository,
	accountRepo repository.AccountRepository,
	accountService *AccountService,
	ledgerService *LedgerService,
) *SavingsGoalService {
	return &SavingsGoalService{
		goalRepo:         goalRepo,
		contributionRepo: contributionRepo,
		accountRepo:      accountRepo,
		accountService:   accountService,
		ledgerService:    ledgerService,
	}
}

type CreateSavingsGoalInput struct {
	Name         string
	TargetAmount int64
	TargetDate   time.Time
	Icon         *string
	Notes        *string
}

func (s *SavingsGoalService) Create(userID uuid.UUID, input CreateSavingsGoalInput) (*models.SavingsGoal, error) {
	if input.Name == "" {
		return nil, errors.New("name is required")
	}
	if input.TargetAmount <= 0 {
		return nil, errors.New("target amount must be positive")
	}
	if input.TargetDate.Before(time.Now().Truncate(24 * time.Hour)) {
		return nil, errors.New("target date must be in the future")
	}

	goal := &models.SavingsGoal{
		ID:           uuid.New(),
		UserID:       userID,
		Name:         input.Name,
		TargetAmount: input.TargetAmount,
		TargetDate:   input.TargetDate,
		Icon:         input.Icon,
		Status:       models.SavingsGoalStatusActive,
		Notes:        input.Notes,
	}

	if err := s.goalRepo.Create(goal); err != nil {
		return nil, err
	}

	// Create linked ASSET account for this savings goal
	accountName := "Tabungan: " + input.Name
	if _, err := s.accountService.CreateLinkedAccount(userID, accountName, models.AccountTypeAsset, goal.ID, "savings_goal"); err != nil {
		return nil, err
	}

	return s.goalRepo.GetByID(goal.ID)
}

func (s *SavingsGoalService) GetByID(id uuid.UUID) (*models.SavingsGoal, error) {
	return s.goalRepo.GetByID(id)
}

func (s *SavingsGoalService) GetByUserID(userID uuid.UUID, status *models.SavingsGoalStatus) ([]models.SavingsGoal, error) {
	return s.goalRepo.GetByUserID(userID, status)
}

func (s *SavingsGoalService) GetActiveByUserID(userID uuid.UUID) ([]models.SavingsGoal, error) {
	return s.goalRepo.GetActiveByUserID(userID)
}

type UpdateSavingsGoalInput struct {
	Name         *string
	TargetAmount *int64
	TargetDate   *time.Time
	Icon         *string
	Notes        *string
	Status       *models.SavingsGoalStatus
}

func (s *SavingsGoalService) Update(id uuid.UUID, input UpdateSavingsGoalInput) (*models.SavingsGoal, error) {
	goal, err := s.goalRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		goal.Name = *input.Name
	}
	if input.TargetAmount != nil {
		if *input.TargetAmount <= 0 {
			return nil, errors.New("target amount must be positive")
		}
		goal.TargetAmount = *input.TargetAmount
	}
	if input.TargetDate != nil {
		goal.TargetDate = *input.TargetDate
	}
	if input.Icon != nil {
		goal.Icon = input.Icon
	}
	if input.Notes != nil {
		goal.Notes = input.Notes
	}
	if input.Status != nil {
		goal.Status = *input.Status
	}

	if err := s.goalRepo.Update(goal); err != nil {
		return nil, err
	}

	return s.goalRepo.GetByID(goal.ID)
}

func (s *SavingsGoalService) Delete(id uuid.UUID) error {
	goal, err := s.goalRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Delete all contribution transactions first (before CASCADE deletes contributions)
	for _, contribution := range goal.Contributions {
		_ = s.ledgerService.DeleteByReference(contribution.ID, "savings_contribution")
	}

	// Delete linked account
	if err := s.accountService.DeleteAccountByReference(id, "savings_goal"); err != nil {
		return err
	}

	return s.goalRepo.Delete(id)
}

func (s *SavingsGoalService) AddContribution(goalID uuid.UUID, amount int64, contributionDate time.Time, notes *string) (*models.SavingsContribution, error) {
	goal, err := s.goalRepo.GetByID(goalID)
	if err != nil {
		return nil, err
	}

	if goal.Status != models.SavingsGoalStatusActive {
		return nil, errors.New("cannot add contribution to non-active savings goal")
	}

	if amount <= 0 {
		return nil, errors.New("contribution amount must be positive")
	}

	contribution := &models.SavingsContribution{
		ID:               uuid.New(),
		SavingsGoalID:    goalID,
		Amount:           amount,
		ContributionDate: contributionDate,
		Notes:            notes,
	}

	if err := s.contributionRepo.Create(contribution); err != nil {
		return nil, err
	}

	// Update current_amount on the goal
	goal.CurrentAmount += amount
	if err := s.goalRepo.Update(goal); err != nil {
		return nil, err
	}

	// Create ledger entry: DEBIT Savings Account (Asset), CREDIT Cash Account (Asset)
	if err := s.createContributionLedgerEntry(goal.UserID, goal, contribution); err != nil {
		return nil, err
	}

	// Auto-complete if target reached
	if goal.CurrentAmount >= goal.TargetAmount {
		goal.Status = models.SavingsGoalStatusCompleted
		if err := s.goalRepo.Update(goal); err != nil {
			return nil, err
		}
	}

	return contribution, nil
}

func (s *SavingsGoalService) WithdrawContribution(contributionID uuid.UUID) error {
	contribution, err := s.contributionRepo.GetByID(contributionID)
	if err != nil {
		return err
	}

	goal, err := s.goalRepo.GetByID(contribution.SavingsGoalID)
	if err != nil {
		return err
	}

	// Delete ledger entry
	if err := s.ledgerService.DeleteByReference(contributionID, "savings_contribution"); err != nil {
		return err
	}

	// Update current_amount on the goal
	goal.CurrentAmount -= contribution.Amount
	if goal.CurrentAmount < 0 {
		goal.CurrentAmount = 0
	}

	// If it was completed, reactivate
	if goal.Status == models.SavingsGoalStatusCompleted {
		goal.Status = models.SavingsGoalStatusActive
	}

	if err := s.goalRepo.Update(goal); err != nil {
		return err
	}

	return s.contributionRepo.Delete(contributionID)
}

func (s *SavingsGoalService) MarkComplete(id uuid.UUID) (*models.SavingsGoal, error) {
	goal, err := s.goalRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	goal.Status = models.SavingsGoalStatusCompleted
	if err := s.goalRepo.Update(goal); err != nil {
		return nil, err
	}

	return s.goalRepo.GetByID(id)
}

func (s *SavingsGoalService) createContributionLedgerEntry(userID uuid.UUID, goal *models.SavingsGoal, contribution *models.SavingsContribution) error {
	// Get savings asset account (linked to goal)
	savingsAccount, err := s.accountRepo.GetByReference(goal.ID, "savings_goal")
	if err != nil {
		return err
	}

	// Get default cash account
	cashAccount, err := s.accountRepo.GetDefaultByUserID(userID)
	if err != nil {
		return err
	}

	entries := []LedgerEntry{
		{AccountID: savingsAccount.ID, Debit: contribution.Amount, Credit: 0},
		{AccountID: cashAccount.ID, Debit: 0, Credit: contribution.Amount},
	}

	_, err = s.ledgerService.CreateJournalEntry(
		userID,
		contribution.ContributionDate,
		"Savings Contribution: "+goal.Name,
		entries,
		&contribution.ID,
		"savings_contribution",
	)
	return err
}

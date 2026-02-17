package services

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/azzamdhx/moneybro/backend/internal/models"
	"github.com/azzamdhx/moneybro/backend/internal/repository"
)

type NotificationService struct {
	repos        *repository.Repositories
	emailService *EmailService
}

func NewNotificationService(repos *repository.Repositories, emailService *EmailService) *NotificationService {
	return &NotificationService{
		repos:        repos,
		emailService: emailService,
	}
}

func (s *NotificationService) GetByUserID(userID uuid.UUID) ([]models.NotificationLog, error) {
	return s.repos.NotificationLog.GetByUserID(userID)
}

func (s *NotificationService) SendDueReminders(ctx context.Context) error {
	now := time.Now()

	// Get all users with notifications enabled
	users, err := s.repos.User.GetAllWithNotificationsEnabled()
	if err != nil {
		log.Printf("Error getting users with notifications: %v", err)
		return err
	}

	for _, user := range users {
		// Send installment reminders based on user preference
		if user.NotifyInstallment {
			s.sendInstallmentRemindersForUser(ctx, user, now)
		}

		// Send debt reminders based on user preference
		if user.NotifyDebt {
			s.sendDebtRemindersForUser(ctx, user, now)
		}

		// Send savings goal reminders based on user preference
		if user.NotifySavingsGoal {
			s.sendSavingsGoalRemindersForUser(ctx, user, now)
		}
	}

	return nil
}

func (s *NotificationService) sendInstallmentRemindersForUser(ctx context.Context, user models.User, now time.Time) {
	today := now.Day()

	// Check for due dates within user's preferred notification window
	for daysAhead := 0; daysAhead <= user.NotifyDaysBefore; daysAhead++ {
		targetDay := today + daysAhead
		if targetDay > 31 {
			targetDay -= 31
		}

		installments, err := s.repos.Installment.GetByDueDay(targetDay, models.InstallmentStatusActive)
		if err != nil {
			log.Printf("Error getting installments for day %d: %v", targetDay, err)
			continue
		}

		for _, inst := range installments {
			// Only process this user's installments
			if inst.UserID != user.ID {
				continue
			}

			// Check if the installment is still within its payment period
			if !isInstallmentDueInMonth(inst.StartDate, inst.Tenor, int(now.Month()), now.Year()) {
				continue
			}

			exists, err := s.repos.NotificationLog.ExistsForToday(inst.UserID, inst.ID, models.NotificationTypeInstallmentReminder)
			if err != nil {
				log.Printf("Error checking notification log: %v", err)
				continue
			}
			if exists {
				continue
			}

			err = s.emailService.SendInstallmentReminder(ctx, user.Email, inst.Name, daysAhead, inst.MonthlyPayment)
			if err != nil {
				log.Printf("Error sending installment reminder to %s: %v", user.Email, err)
				continue
			}

			subject := "Reminder: Cicilan " + inst.Name + " jatuh tempo"
			logEntry := &models.NotificationLog{
				UserID:       inst.UserID,
				Type:         models.NotificationTypeInstallmentReminder,
				ReferenceID:  inst.ID,
				SentAt:       now,
				EmailSubject: &subject,
			}
			if err := s.repos.NotificationLog.Create(logEntry); err != nil {
				log.Printf("Error creating notification log: %v", err)
			}
			log.Printf("Sent installment reminder to %s for %s (due in %d days)", user.Email, inst.Name, daysAhead)
		}
	}
}

func (s *NotificationService) sendDebtRemindersForUser(ctx context.Context, user models.User, now time.Time) {
	// Check for due dates within user's preferred notification window
	for daysAhead := 0; daysAhead <= user.NotifyDaysBefore; daysAhead++ {
		targetDate := now.AddDate(0, 0, daysAhead)
		startDate := targetDate.Format("2006-01-02")
		endDate := targetDate.Format("2006-01-02")

		debts, err := s.repos.Debt.GetByDueDateRange(startDate, endDate, models.DebtStatusActive)
		if err != nil {
			log.Printf("Error getting debts for date %s: %v", startDate, err)
			continue
		}

		for _, debt := range debts {
			// Only process this user's debts
			if debt.UserID != user.ID {
				continue
			}

			exists, err := s.repos.NotificationLog.ExistsForToday(debt.UserID, debt.ID, models.NotificationTypeDebtReminder)
			if err != nil {
				log.Printf("Error checking notification log: %v", err)
				continue
			}
			if exists {
				continue
			}

			err = s.emailService.SendDebtReminder(ctx, user.Email, debt.PersonName, daysAhead, debt.RemainingAmount())
			if err != nil {
				log.Printf("Error sending debt reminder to %s: %v", user.Email, err)
				continue
			}

			subject := "Reminder: Hutang ke " + debt.PersonName + " jatuh tempo"
			logEntry := &models.NotificationLog{
				UserID:       debt.UserID,
				Type:         models.NotificationTypeDebtReminder,
				ReferenceID:  debt.ID,
				SentAt:       now,
				EmailSubject: &subject,
			}
			if err := s.repos.NotificationLog.Create(logEntry); err != nil {
				log.Printf("Error creating notification log: %v", err)
			}
			log.Printf("Sent debt reminder to %s for %s (due in %d days)", user.Email, debt.PersonName, daysAhead)
		}
	}
}

func (s *NotificationService) sendSavingsGoalRemindersForUser(ctx context.Context, user models.User, now time.Time) {
	// Check for savings goals with target dates within user's preferred notification window
	for daysAhead := 0; daysAhead <= user.NotifyDaysBefore; daysAhead++ {
		targetDate := now.AddDate(0, 0, daysAhead)
		startDate := targetDate.Format("2006-01-02")
		endDate := targetDate.Format("2006-01-02")

		goals, err := s.repos.SavingsGoal.GetByTargetDateRange(startDate, endDate, models.SavingsGoalStatusActive)
		if err != nil {
			log.Printf("Error getting savings goals for date %s: %v", startDate, err)
			continue
		}

		for _, goal := range goals {
			// Only process this user's goals
			if goal.UserID != user.ID {
				continue
			}

			exists, err := s.repos.NotificationLog.ExistsForToday(goal.UserID, goal.ID, models.NotificationTypeSavingsGoalReminder)
			if err != nil {
				log.Printf("Error checking notification log: %v", err)
				continue
			}
			if exists {
				continue
			}

			err = s.emailService.SendSavingsGoalReminder(ctx, user.Email, goal.Name, daysAhead, goal.RemainingAmount(), goal.CurrentAmount, goal.Progress())
			if err != nil {
				log.Printf("Error sending savings goal reminder to %s: %v", user.Email, err)
				continue
			}

			subject := "Reminder: Target tabungan " + goal.Name + " deadline"
			logEntry := &models.NotificationLog{
				UserID:       goal.UserID,
				Type:         models.NotificationTypeSavingsGoalReminder,
				ReferenceID:  goal.ID,
				SentAt:       now,
				EmailSubject: &subject,
			}
			if err := s.repos.NotificationLog.Create(logEntry); err != nil {
				log.Printf("Error creating notification log: %v", err)
			}
			log.Printf("Sent savings goal reminder to %s for %s (deadline in %d days)", user.Email, goal.Name, daysAhead)
		}
	}
}

// isInstallmentDueInMonth checks if an installment has a payment due in the specified month
func isInstallmentDueInMonth(startDate time.Time, tenor int, month int, year int) bool {
	targetMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	// Check if the installment has started by the target month
	if startDate.After(targetMonth.AddDate(0, 1, -1)) {
		return false
	}

	// If startDate is in the future relative to target month, not due yet
	if startDate.Year() > year || (startDate.Year() == year && int(startDate.Month()) > month) {
		return false
	}

	// Check if the installment has ended before the target month
	// End month is startDate + (tenor - 1) months (since first payment is in start month)
	endDate := time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, time.UTC).AddDate(0, tenor, 0)
	if targetMonth.After(endDate) || targetMonth.Equal(endDate) {
		return false
	}

	return true
}

package services

import (
	"context"
	"log"
	"time"

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

func (s *NotificationService) SendDueReminders(ctx context.Context) error {
	now := time.Now()
	today := now.Day()

	for daysAhead := 1; daysAhead <= 3; daysAhead++ {
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
			if inst.User == nil {
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

			err = s.emailService.SendInstallmentReminder(ctx, inst.User.Email, inst.Name, daysAhead, inst.MonthlyPayment)
			if err != nil {
				log.Printf("Error sending installment reminder: %v", err)
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
		}
	}

	for daysAhead := 1; daysAhead <= 3; daysAhead++ {
		targetDate := now.AddDate(0, 0, daysAhead)
		startDate := targetDate.Format("2006-01-02")
		endDate := targetDate.Format("2006-01-02")

		debts, err := s.repos.Debt.GetByDueDateRange(startDate, endDate, models.DebtStatusActive)
		if err != nil {
			log.Printf("Error getting debts for date %s: %v", startDate, err)
			continue
		}

		for _, debt := range debts {
			if debt.User == nil {
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

			err = s.emailService.SendDebtReminder(ctx, debt.User.Email, debt.PersonName, daysAhead, debt.RemainingAmount())
			if err != nil {
				log.Printf("Error sending debt reminder: %v", err)
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
		}
	}

	return nil
}

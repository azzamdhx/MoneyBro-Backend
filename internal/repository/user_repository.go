package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/azzamdhx/moneybro/backend/internal/models"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetAllWithNotificationsEnabled() ([]models.User, error) {
	var users []models.User
	err := r.db.Where("notify_installment = ? OR notify_debt = ?", true, true).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.User{}, "id = ?", id).Error
}

func (r *userRepository) DeleteAllUserData(userID uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Delete in order to respect foreign key constraints

		// First delete transaction_entries (references transactions and accounts)
		if err := tx.Exec("DELETE FROM transaction_entries WHERE transaction_id IN (SELECT id FROM transactions WHERE user_id = ?)", userID).Error; err != nil {
			return err
		}

		// Delete transactions (audit trail)
		if err := tx.Exec("DELETE FROM transactions WHERE user_id = ?", userID).Error; err != nil {
			return err
		}

		// Delete payment records
		if err := tx.Exec("DELETE FROM installment_payments WHERE installment_id IN (SELECT id FROM installments WHERE user_id = ?)", userID).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM debt_payments WHERE debt_id IN (SELECT id FROM debts WHERE user_id = ?)", userID).Error; err != nil {
			return err
		}

		// Delete main records
		if err := tx.Exec("DELETE FROM expenses WHERE user_id = ?", userID).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM expense_templates WHERE user_id = ?", userID).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM installments WHERE user_id = ?", userID).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM debts WHERE user_id = ?", userID).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM incomes WHERE user_id = ?", userID).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM recurring_incomes WHERE user_id = ?", userID).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM notification_logs WHERE user_id = ?", userID).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM password_reset_tokens WHERE user_id = ?", userID).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM two_fa_codes WHERE user_id = ?", userID).Error; err != nil {
			return err
		}

		// Delete accounts (ledger accounts)
		if err := tx.Exec("DELETE FROM accounts WHERE user_id = ?", userID).Error; err != nil {
			return err
		}

		// Delete categories
		if err := tx.Exec("DELETE FROM categories WHERE user_id = ?", userID).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM income_categories WHERE user_id = ?", userID).Error; err != nil {
			return err
		}

		// Finally delete user
		if err := tx.Delete(&models.User{}, "id = ?", userID).Error; err != nil {
			return err
		}

		return nil
	})
}

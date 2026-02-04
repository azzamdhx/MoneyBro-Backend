package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/azzamdhx/moneybro/backend/internal/models"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Starting to fix transaction dates for installment payments...")

	if err := fixInstallmentPaymentDates(db); err != nil {
		log.Fatalf("Failed to fix installment payment dates: %v", err)
	}

	log.Println("Successfully fixed all transaction dates!")
}

func fixInstallmentPaymentDates(db *gorm.DB) error {
	var payments []models.InstallmentPayment
	if err := db.Preload("Installment").Find(&payments).Error; err != nil {
		return err
	}

	for _, payment := range payments {
		if payment.Installment == nil {
			log.Printf("Warning: installment not found for payment %s, skipping", payment.ID)
			continue
		}

		// Calculate correct period date based on start_date + (payment_number - 1) months
		periodDate := payment.Installment.StartDate.AddDate(0, payment.PaymentNumber-1, 0)

		// Find the transaction for this payment
		var tx models.Transaction
		if err := db.Where("reference_id = ? AND reference_type = ?", payment.ID, "installment_payment").First(&tx).Error; err != nil {
			log.Printf("Warning: transaction not found for payment %s, skipping", payment.ID)
			continue
		}

		// Update the transaction date
		if err := db.Model(&tx).Update("transaction_date", periodDate).Error; err != nil {
			return fmt.Errorf("failed to update transaction %s: %w", tx.ID, err)
		}

		log.Printf("Fixed payment #%d for %s: %s -> %s",
			payment.PaymentNumber,
			payment.Installment.Name,
			tx.TransactionDate.Format("2006-01-02"),
			periodDate.Format("2006-01-02"))
	}

	return nil
}

package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/resend/resend-go/v3"
)

type EmailService struct {
	client       *resend.Client
	fromEmail    string
	templatesDir string
}

func NewEmailService(apiKey, templatesDir string) *EmailService {
	return &EmailService{
		client:       resend.NewClient(apiKey),
		fromEmail:    "MoneyBro <noreply@moneybro.my.id>",
		templatesDir: templatesDir,
	}
}

type EmailParams struct {
	To      string
	Subject string
	HTML    string
}

func (s *EmailService) Send(ctx context.Context, params EmailParams) error {
	_, err := s.client.Emails.SendWithContext(ctx, &resend.SendEmailRequest{
		From:    s.fromEmail,
		To:      []string{params.To},
		Subject: params.Subject,
		Html:    params.HTML,
	})
	return err
}

func (s *EmailService) loadTemplate(name string) (string, error) {
	path := filepath.Join(s.templatesDir, name)
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to load template %s: %w", name, err)
	}
	return string(content), nil
}

func (s *EmailService) renderTemplate(template string, variables map[string]interface{}) string {
	result := template
	re := regexp.MustCompile(`\{\{\{(\w+)\}\}\}`)

	result = re.ReplaceAllStringFunc(result, func(match string) string {
		key := re.FindStringSubmatch(match)[1]
		if val, ok := variables[key]; ok {
			switch v := val.(type) {
			case string:
				return v
			case int:
				return strconv.Itoa(v)
			case int64:
				return formatCurrency(v)
			default:
				return fmt.Sprintf("%v", v)
			}
		}
		return match
	})

	return result
}

func formatCurrency(amount int64) string {
	str := strconv.FormatInt(amount, 10)
	n := len(str)
	if n <= 3 {
		return str
	}

	var result []byte
	for i, digit := range str {
		if i > 0 && (n-i)%3 == 0 {
			result = append(result, '.')
		}
		result = append(result, byte(digit))
	}
	return string(result)
}

func (s *EmailService) SendInstallmentReminder(ctx context.Context, to, name string, daysUntil int, amount int64) error {
	template, err := s.loadTemplate("installment_reminder.html")
	if err != nil {
		return err
	}

	html := s.renderTemplate(template, map[string]interface{}{
		"name":       name,
		"days_until": daysUntil,
		"amount":     amount,
	})

	return s.Send(ctx, EmailParams{
		To:      to,
		Subject: fmt.Sprintf("Reminder: Cicilan %s jatuh tempo dalam %d hari", name, daysUntil),
		HTML:    html,
	})
}

func (s *EmailService) SendDebtReminder(ctx context.Context, to, personName string, daysUntil int, amount int64) error {
	template, err := s.loadTemplate("debt_reminder.html")
	if err != nil {
		return err
	}

	html := s.renderTemplate(template, map[string]interface{}{
		"person_name": personName,
		"days_until":  daysUntil,
		"amount":      amount,
	})

	return s.Send(ctx, EmailParams{
		To:      to,
		Subject: fmt.Sprintf("Reminder: Hutang ke %s jatuh tempo dalam %d hari", personName, daysUntil),
		HTML:    html,
	})
}

func (s *EmailService) SendPasswordResetEmail(ctx context.Context, to, name, resetLink string) error {
	template, err := s.loadTemplate("password_reset.html")
	if err != nil {
		return err
	}

	html := s.renderTemplate(template, map[string]interface{}{
		"name":       name,
		"reset_link": resetLink,
	})

	return s.Send(ctx, EmailParams{
		To:      to,
		Subject: "Reset Password MoneyBro",
		HTML:    html,
	})
}

func (s *EmailService) Send2FACodeEmail(ctx context.Context, to, name, code string) error {
	template, err := s.loadTemplate("2fa_code.html")
	if err != nil {
		return err
	}

	html := s.renderTemplate(template, map[string]interface{}{
		"name": name,
		"code": code,
	})

	return s.Send(ctx, EmailParams{
		To:      to,
		Subject: "Kode Verifikasi MoneyBro",
		HTML:    html,
	})
}

package services

import (
	"context"
	"fmt"

	"github.com/resend/resend-go/v2"
)

type EmailService struct {
	client *resend.Client
}

func NewEmailService(apiKey string) *EmailService {
	return &EmailService{
		client: resend.NewClient(apiKey),
	}
}

type EmailParams struct {
	To      string
	Subject string
	HTML    string
}

func (s *EmailService) Send(ctx context.Context, params EmailParams) error {
	_, err := s.client.Emails.Send(&resend.SendEmailRequest{
		From:    "MoneyBro <noreply@moneybro.app>",
		To:      []string{params.To},
		Subject: params.Subject,
		Html:    params.HTML,
	})
	return err
}

func (s *EmailService) SendInstallmentReminder(ctx context.Context, to, name string, daysUntil int, amount int64) error {
	subject := fmt.Sprintf("Reminder: Cicilan %s jatuh tempo dalam %d hari", name, daysUntil)
	html := fmt.Sprintf(`
		<h2>Reminder Pembayaran Cicilan</h2>
		<p>Halo,</p>
		<p>Cicilan <strong>%s</strong> akan jatuh tempo dalam <strong>%d hari</strong>.</p>
		<p>Jumlah pembayaran: <strong>Rp %d</strong></p>
		<p>Pastikan saldo Anda cukup untuk pembayaran.</p>
		<br>
		<p>Salam,<br>MoneyBro</p>
	`, name, daysUntil, amount)

	return s.Send(ctx, EmailParams{
		To:      to,
		Subject: subject,
		HTML:    html,
	})
}

func (s *EmailService) SendDebtReminder(ctx context.Context, to, personName string, daysUntil int, amount int64) error {
	subject := fmt.Sprintf("Reminder: Hutang ke %s jatuh tempo dalam %d hari", personName, daysUntil)
	html := fmt.Sprintf(`
		<h2>Reminder Pembayaran Hutang</h2>
		<p>Halo,</p>
		<p>Hutang kepada <strong>%s</strong> akan jatuh tempo dalam <strong>%d hari</strong>.</p>
		<p>Sisa hutang: <strong>Rp %d</strong></p>
		<p>Jangan lupa untuk melakukan pembayaran.</p>
		<br>
		<p>Salam,<br>MoneyBro</p>
	`, personName, daysUntil, amount)

	return s.Send(ctx, EmailParams{
		To:      to,
		Subject: subject,
		HTML:    html,
	})
}

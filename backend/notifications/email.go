package notifications

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"
)

type EmailNotifier struct {
	smtpHost string
	smtpPort int
	username string
	password string
	from     string
	to       []string
}

func NewEmailNotifier(smtpHost string, smtpPort int, username, password, from string, to []string) *EmailNotifier {
	return &EmailNotifier{
		smtpHost: smtpHost,
		smtpPort: smtpPort,
		username: username,
		password: password,
		from:     from,
		to:       to,
	}
}

func (e *EmailNotifier) Send(_ context.Context, subject string, body string) error {
	addr := fmt.Sprintf("%s:%d", e.smtpHost, e.smtpPort)
	auth := smtp.PlainAuth("", e.username, e.password, e.smtpHost)

	headers := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n",
		e.from, strings.Join(e.to, ","), subject)

	msg := []byte(headers + body)
	return smtp.SendMail(addr, auth, e.from, e.to, msg)
}

func (e *EmailNotifier) HealthCheck(_ context.Context) error {
	addr := fmt.Sprintf("%s:%d", e.smtpHost, e.smtpPort)
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("smtp connect failed: %w", err)
	}
	defer client.Close()
	return client.Hello("localhost")
}

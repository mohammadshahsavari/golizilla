package service

import (
	"fmt"
	"golizilla/config"
)

// IEmailService defines methods for sending emails.
type IEmailService interface {
	SendEmail(to string, subject string, body string) error
}

// EmailService is a basic implementation of IEmailService.
type EmailService struct {
	cfg *config.Config
}

// NewEmailService creates a new EmailService instance.
func NewEmailService(cfg *config.Config) IEmailService {
	return &EmailService{cfg: cfg}
}

// SendEmail sends an email (placeholder implementation).
func (s *EmailService) SendEmail(to string, subject string, body string) error {
	// TODO: Implement actual email sending logic.
	fmt.Printf("Sending email to: %s\nSubject: %s\nBody:\n%s\n", to, subject, body)
	return nil
}

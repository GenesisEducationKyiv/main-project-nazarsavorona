package services

import (
	"context"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/models"

	"golang.org/x/sync/errgroup"
)

type (
	EmailSender interface {
		SendEmail(to, subject, body string) error
	}

	EmailService struct {
		emailSender EmailSender
	}
)

func NewEmailService(emailSender EmailSender) *EmailService {
	return &EmailService{
		emailSender: emailSender,
	}
}

func (s *EmailService) SendEmails(ctx context.Context, emails []string, message *models.Message) error {
	group, _ := errgroup.WithContext(ctx)

	for _, email := range emails {
		email := email
		group.Go(func() error {
			return s.emailSender.SendEmail(email, message.Subject, message.Body)
		})
	}

	err := group.Wait()
	if err != nil {
		return err
	}

	return nil
}

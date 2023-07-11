package services

import (
	"context"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/models"

	"golang.org/x/sync/errgroup"
)

type (
	EmailSender interface {
		SendEmail(string, string, string) error
	}

	EmailService struct {
		mailSender EmailSender
	}
)

func NewEmailService(mailSender EmailSender) *EmailService {
	return &EmailService{mailSender: mailSender}
}

func (s *EmailService) SendEmails(ctx context.Context, emails []string, message *models.Message) error {
	// since we want to try to send all emails, we don't want to stop on first error
	group, _ := errgroup.WithContext(ctx)

	for _, currentEmail := range emails {
		currentEmail := currentEmail
		group.Go(func() error {
			return s.mailSender.SendEmail(currentEmail, message.Subject, message.Body)
		})
	}

	err := group.Wait()
	if err != nil {
		return err
	}

	return nil
}

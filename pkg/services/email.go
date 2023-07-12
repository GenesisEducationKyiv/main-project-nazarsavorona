package services

import (
	"context"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/email"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/models"

	"golang.org/x/sync/errgroup"
)

type (
	EmailSender interface {
		SendEmail(to, subject, body string, strategy email.MessageConstructStrategy) error
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

	for _, e := range emails {
		e := e
		group.Go(func() error {
			return s.emailSender.SendEmail(e, message.Subject, message.Body, &email.HTMLMessageBuilder{})
		})
	}

	err := group.Wait()
	if err != nil {
		return err
	}

	return nil
}

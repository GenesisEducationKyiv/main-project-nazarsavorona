package services

import (
	"context"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/models"

	"golang.org/x/sync/errgroup"
)

type (
	EmailSender interface {
		SendEmail(to string, message []byte) error
	}

	MessageConstructStrategy interface {
		Construct(message *models.Message) []byte
	}

	EmailService struct {
		emailSender EmailSender
		strategy    MessageConstructStrategy
	}
)

func NewEmailService(emailSender EmailSender, strategy MessageConstructStrategy) *EmailService {
	return &EmailService{
		emailSender: emailSender,
		strategy:    strategy,
	}
}

func (s *EmailService) SendEmails(ctx context.Context, emails []string, message *models.Message) error {
	group, _ := errgroup.WithContext(ctx)

	for _, e := range emails {
		e := e
		group.Go(func() error {
			messageBytes := s.strategy.Construct(message)
			return s.emailSender.SendEmail(e, messageBytes)
		})
	}

	err := group.Wait()
	if err != nil {
		return err
	}

	return nil
}

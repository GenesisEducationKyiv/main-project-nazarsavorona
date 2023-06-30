package service

import (
	"context"
	"fmt"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/clients"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/email"
)

type Service struct {
	fromCurrency string
	toCurrency   string

	repository *email.Repository
	mailSender *email.Sender
	rateGetter *clients.BinanceClient
}

func NewService(from, to string,
	repository *email.Repository,
	mailSender *email.Sender,
	rateGetter *clients.BinanceClient) *Service {
	return &Service{
		fromCurrency: from,
		toCurrency:   to,
		repository:   repository,
		mailSender:   mailSender,
		rateGetter:   rateGetter,
	}
}

func (s *Service) Subscribe(email string) error {
	emails := s.repository.EmailList()
	for _, currentEmail := range emails {
		if currentEmail == email {
			return fmt.Errorf("email %s is already subscribed", email)
		}
	}

	err := s.repository.AddEmail(email)
	if err != nil {
		return err
	}

	err = s.mailSender.SendEmail(email, "Subscription", "You have successfully subscribed to the service")
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) SendEmails(ctx context.Context) error {
	emails := s.repository.EmailList()

	rate, err := s.Rate(ctx)
	if err != nil {
		return err
	}

	errs := make(chan error, len(emails))

	for _, currentEmail := range emails {
		go func(email string) {
			errs <- s.mailSender.SendEmail(email,
				fmt.Sprintf("%s rate", s.fromCurrency),
				fmt.Sprintf("1 %s = %.2f %s", s.fromCurrency, rate, s.toCurrency))
		}(currentEmail)
	}

	for i := 0; i < len(emails); i++ {
		err = <-errs
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) Rate(ctx context.Context) (float64, error) {
	return s.rateGetter.Rate(ctx)
}

func (s *Service) EmailList() []string {
	return s.repository.EmailList()
}

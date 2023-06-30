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

func NewService(smtpHost, smtpPort, accountEmail, accountPassword, from, to string, db email.Database) *Service {
	return &Service{
		fromCurrency: from,
		toCurrency:   to,
		repository:   email.NewRepository(db),
		mailSender:   email.NewSender(smtpHost, smtpPort, accountEmail, accountPassword),
		// TODO: consider where the api url should come from
		rateGetter: clients.NewBinanceClient(from, to, "https://api.binance.com/api/v3/"),
	}
}

func (s *Service) Subscribe(email string) error {
	emails := s.repository.GetEmailList()
	for _, currentEmail := range emails {
		if currentEmail == email {
			return fmt.Errorf("email %s is already subscribed", email)
		}
	}

	err := s.repository.AddNewEmail(email)
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
	emails := s.repository.GetEmailList()

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
	return s.repository.GetEmailList()
}

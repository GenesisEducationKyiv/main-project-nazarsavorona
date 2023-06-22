package service

import (
	"fmt"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/email"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/currency"
)

type Service struct {
	fromCurrency string
	toCurrency   string

	repository *email.Repository
	mailSender *email.Sender
	rateGetter *currency.BinanceGetter
}

func NewService(smtpHost, smtpPort, accountEmail, accountPassword, from, to string, db email.Database) *Service {
	return &Service{
		fromCurrency: from,
		toCurrency:   to,
		repository:   email.NewRepository(db),
		mailSender:   email.NewSender(smtpHost, smtpPort, accountEmail, accountPassword),
		rateGetter:   currency.NewGetter(from, to),
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

func (s *Service) SendEmails() error {
	emails := s.repository.GetEmailList()

	rate, err := s.GetRate()
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

func (s *Service) GetRate() (float64, error) {
	return s.rateGetter.GetRate()
}

func (s *Service) GetEmailList() []string {
	return s.repository.GetEmailList()
}

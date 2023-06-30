package service

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/sync/errgroup"

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

var ErrAlreadySubscribed = fmt.Errorf("email is already subscribed")

func (s *Service) Subscribe(candidateEmail string) error {
	err := s.repository.AddEmail(candidateEmail)
	if err != nil {
		if errors.Is(err, email.ErrAlreadyExists) {
			return ErrAlreadySubscribed
		}
		return err
	}

	err = s.mailSender.SendEmail(candidateEmail, "Subscription", "You have successfully subscribed to the service")
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

	// since we want to try to send all emails, we don't want to stop on first error
	group, _ := errgroup.WithContext(ctx)

	for _, currentEmail := range emails {
		currentEmail := currentEmail
		group.Go(func() error {
			return s.mailSender.SendEmail(currentEmail,
				fmt.Sprintf("%s rate", s.fromCurrency),
				fmt.Sprintf("1 %s = %.2f %s", s.fromCurrency, rate, s.toCurrency))
		})
	}

	err = group.Wait()
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Rate(ctx context.Context) (float64, error) {
	return s.rateGetter.Rate(ctx)
}

func (s *Service) EmailList() []string {
	return s.repository.EmailList()
}

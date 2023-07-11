package services

import (
	"errors"
	"fmt"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/email"
)

type (
	EmailRepository interface {
		AddEmail(string) error
		EmailList() []string
	}

	SubscribeService struct {
		repository EmailRepository
	}
)

func NewSubscribeService(repository EmailRepository) *SubscribeService {
	return &SubscribeService{repository: repository}
}

var ErrAlreadySubscribed = fmt.Errorf("email is already subscribed")

func (s *SubscribeService) Subscribe(candidateEmail string) error {
	err := s.repository.AddEmail(candidateEmail)
	if err != nil {
		if errors.Is(err, email.ErrAlreadyExists) {
			return ErrAlreadySubscribed
		}
		return err
	}

	return nil
}

func (s *SubscribeService) EmailList() []string {
	return s.repository.EmailList()
}

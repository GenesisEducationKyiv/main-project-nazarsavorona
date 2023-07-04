package services

import (
	"context"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/models"
)

type (
	RateGetter interface {
		Rate(context.Context) (*models.Rate, error)
	}

	RateService struct {
		rateGetter RateGetter
	}
)

func NewRateService(rateGetter RateGetter) *RateService {
	return &RateService{rateGetter: rateGetter}
}

func (s *RateService) Rate(ctx context.Context) (*models.Rate, error) {
	return s.rateGetter.Rate(ctx)
}

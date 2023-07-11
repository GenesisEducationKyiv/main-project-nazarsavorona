package services

import (
	"context"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/models"
)

type (
	RateGetter interface {
		Rate(ctx context.Context, from, to string) (*models.Rate, error)
	}

	RateService struct {
		from       string
		to         string
		rateGetter RateGetter
	}
)

func NewRateService(from, to string, rateGetter RateGetter) *RateService {
	return &RateService{
		from:       from,
		to:         to,
		rateGetter: rateGetter,
	}
}

func (s *RateService) Rate(ctx context.Context) (*models.Rate, error) {
	return s.rateGetter.Rate(ctx, s.from, s.to)
}

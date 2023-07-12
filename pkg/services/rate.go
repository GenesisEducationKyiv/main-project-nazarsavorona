package services

import (
	"context"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/models"
)

type (
	RateFetcher interface {
		Rate(ctx context.Context, from, to string) (*models.Rate, error)
	}

	RateService struct {
		from    string
		to      string
		fetcher RateFetcher
	}
)

func NewRateService(from, to string, fetcher RateFetcher) *RateService {
	return &RateService{
		from:    from,
		to:      to,
		fetcher: fetcher,
	}
}

func (s *RateService) Rate(ctx context.Context) (*models.Rate, error) {
	return s.fetcher.Rate(ctx, s.from, s.to)
}

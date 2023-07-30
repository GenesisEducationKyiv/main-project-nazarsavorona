package chain

import (
	"context"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/models"
	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/services"
)

type Chain interface {
	services.RateFetcher
	SetNext(Chain)
}

type BaseChain struct {
	rateGetter services.RateFetcher
	next       Chain
}

func NewBaseChain(rateGetter services.RateFetcher) *BaseChain {
	return &BaseChain{
		rateGetter: rateGetter,
	}
}

func (b *BaseChain) SetNext(next Chain) {
	b.next = next
}

func (b *BaseChain) Rate(ctx context.Context, from, to string) (*models.Rate, error) {
	rate, err := b.rateGetter.Rate(ctx, from, to)
	if err != nil {
		next := b.next
		if next == nil {
			return nil, err
		}

		return next.Rate(ctx, from, to)
	}

	return rate, nil
}

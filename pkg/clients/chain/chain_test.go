package chain_test

import (
	"context"
	"errors"
	"testing"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/clients/chain"
	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/models"
	"github.com/stretchr/testify/require"
)

type RateFetcherMock struct {
	RateFunc func(ctx context.Context, from, to string) (*models.Rate, error)
}

func (m *RateFetcherMock) Rate(ctx context.Context, from, to string) (*models.Rate, error) {
	return m.RateFunc(ctx, from, to)
}

func TestBaseChain(t *testing.T) {
	t.Parallel()

	rate := models.NewRate("", "", 1)

	success := func(ctx context.Context, from, to string) (*models.Rate, error) {
		return rate, nil
	}

	fail := func(ctx context.Context, from, to string) (*models.Rate, error) {
		return nil, errors.New("fail")
	}

	tests := []struct {
		name     string
		fetchers []*RateFetcherMock
		want     *models.Rate
		wantErr  require.ErrorAssertionFunc
	}{
		{
			name: "success",
			fetchers: []*RateFetcherMock{
				{
					RateFunc: success,
				},
				{
					RateFunc: success,
				},
			},
			want:    rate,
			wantErr: require.NoError,
		},
		{
			name: "first fails",
			fetchers: []*RateFetcherMock{
				{
					RateFunc: fail,
				},
				{
					RateFunc: success,
				},
			},
			want:    rate,
			wantErr: require.NoError,
		},
		{
			name: "each success fails",
			fetchers: []*RateFetcherMock{
				{
					RateFunc: fail,
				},
				{
					RateFunc: fail,
				},
			},
			want:    nil,
			wantErr: require.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c1 := chain.NewBaseChain(tt.fetchers[0])
			c1.SetNext(chain.NewBaseChain(tt.fetchers[1]))

			got, err := c1.Rate(context.Background(), "", "")
			tt.wantErr(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

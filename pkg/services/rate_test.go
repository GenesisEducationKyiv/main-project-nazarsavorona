package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/models"
	"github.com/stretchr/testify/require"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/services"
)

type testRateGetter struct {
	rate *models.Rate
	err  error
}

func newTestRateGetter(rate *models.Rate, err error) *testRateGetter {
	return &testRateGetter{
		rate: rate,
		err:  err,
	}
}

func (t *testRateGetter) Rate(_ context.Context, _, _ string) (*models.Rate, error) {
	return t.rate, t.err
}

func TestRateService(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		rateGetter services.RateFetcher
		want       *models.Rate
		expectErr  require.ErrorAssertionFunc
	}{
		{
			name: "success",
			rateGetter: newTestRateGetter(&models.Rate{
				From: "USD",
				To:   "UAH",
				Rate: 2,
			}, nil),
			want: &models.Rate{
				From: "USD",
				To:   "UAH",
				Rate: 2,
			},
			expectErr: require.NoError,
		},
		{
			name:       "error",
			rateGetter: newTestRateGetter(nil, errors.New("test error")),
			want:       nil,
			expectErr:  require.Error,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := services.NewRateService("", "", tt.rateGetter)
			got, err := s.Rate(context.Background())
			tt.expectErr(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

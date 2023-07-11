package models_test

import (
	"testing"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/models"
	"github.com/stretchr/testify/require"
)

func TestNewMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		r    *models.Rate
		want *models.Message
	}{
		{
			name: "check message",
			r: &models.Rate{
				From: "USD",
				To:   "UAH",
				Rate: 27.5,
			},
			want: &models.Message{
				Subject: "USD rate",
				Body:    "1 USD = 27.50 UAH",
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, models.NewMessageFromRate(tt.r))
		})
	}
}

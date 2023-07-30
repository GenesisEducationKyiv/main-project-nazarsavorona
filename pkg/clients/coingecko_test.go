package clients_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/clients"
	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/models"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
)

type mockCoingeckoClient struct{}

func (m *mockCoingeckoClient) Do(_ *http.Request) (*http.Response, error) {
	response := &http.Response{
		Body: io.NopCloser(bytes.NewReader(
			[]byte(fmt.Sprintf(` {
										  "bitcoin": { 
											"uah": %f
										  }
										}`, testValue)))),
		StatusCode: http.StatusOK,
	}

	return response, nil
}

func TestCoingeckoClient_Rate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		from    string
		to      string
		client  clients.HTTPClient
		wantErr require.ErrorAssertionFunc
	}{
		{
			name:    "success",
			from:    "BTC",
			to:      "UAH",
			client:  &mockCoingeckoClient{},
			wantErr: require.NoError,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			want := models.NewRate(tt.from, tt.to, testValue)

			client := clients.NewCoingeckoClient("", tt.client)
			got, err := client.Rate(context.Background(), tt.from, tt.to)

			tt.wantErr(t, err)
			require.Equal(t, want, got)
		})
	}
}

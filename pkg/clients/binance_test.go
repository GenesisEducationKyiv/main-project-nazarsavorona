package clients_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/clients"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/models"
)

type mockBinanceClient struct{}

const testValue = 31415.926

func (m *mockBinanceClient) Do(_ *http.Request) (*http.Response, error) {
	response := &http.Response{
		Body: io.NopCloser(bytes.NewReader(
			[]byte(fmt.Sprintf(`{"price": "%f"}`, testValue)))),
		StatusCode: http.StatusOK,
	}

	return response, nil
}

func TestBinanceClient_Rate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		client  clients.HTTPClient
		want    *models.Rate
		wantErr require.ErrorAssertionFunc
	}{
		{
			name:    "success",
			client:  &mockBinanceClient{},
			want:    models.NewRate("", "", testValue),
			wantErr: require.NoError,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := clients.NewBinanceClient("", tt.client)

			got, err := client.Rate(context.Background(), "", "")

			tt.wantErr(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

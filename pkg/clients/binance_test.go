package clients_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/clients"
)

type mockClient struct{}

const testValue = 31415.926

func (m *mockClient) Do(_ *http.Request) (*http.Response, error) {
	responseData := struct {
		Price string `json:"price"`
	}{
		Price: fmt.Sprintf("%f", testValue),
	}

	responseBody, err := json.Marshal(responseData)
	if err != nil {
		return nil, err
	}

	response := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBuffer(responseBody)),
	}

	return response, nil
}

func TestBinanceClient_Rate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		client  clients.HTTPClient
		want    float64
		wantErr require.ErrorAssertionFunc
	}{
		{
			name:    "success",
			client:  &mockClient{},
			want:    testValue,
			wantErr: require.NoError,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			binanceClient := clients.NewBinanceClient("", "", "", tt.client)

			got, err := binanceClient.Rate(context.Background())

			tt.wantErr(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

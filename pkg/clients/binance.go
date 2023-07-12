package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/models"
)

type (
	HTTPClient interface {
		Do(*http.Request) (*http.Response, error)
	}

	BinanceClient struct {
		apiURL string
		client HTTPClient
	}

	rateDTO struct {
		Price string `json:"price"`
	}
)

func NewBinanceClient(apiURL string, client HTTPClient) *BinanceClient {
	return &BinanceClient{
		apiURL: apiURL,
		client: client,
	}
}

func (b *BinanceClient) Rate(ctx context.Context, from, to string) (*models.Rate, error) {
	url := fmt.Sprintf("%sticker/price?symbol=%s%s", b.apiURL, from, to)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	response, err := b.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %s", response.Status)
	}

	var r rateDTO
	err = json.NewDecoder(response.Body).Decode(&r)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	price, err := strconv.ParseFloat(r.Price, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse rate price: %w", err)
	}

	return models.NewRate(from, to, price), nil
}

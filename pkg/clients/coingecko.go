package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/models"
)

type CoingeckoClient struct {
	apiURL      string
	client      HTTPClient
	currencyMap map[string]string
}

func NewCoingeckoClient(apiURL string, client HTTPClient) *CoingeckoClient {
	return &CoingeckoClient{
		apiURL: apiURL,
		client: client,
		currencyMap: map[string]string{
			"UAH": "uah",
			"BTC": "bitcoin",
		},
	}
}

func (c *CoingeckoClient) Rate(ctx context.Context, from, to string) (*models.Rate, error) {
	fromCurrency := c.currencyMap[from]
	toCurrency := c.currencyMap[to]

	url := fmt.Sprintf("%ssimple/price?ids=%s&vs_currencies=%s",
		c.apiURL, fromCurrency, toCurrency)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	response, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %s", response.Status)
	}

	var responseData map[string]map[string]float64
	if err = json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		return nil, err
	}

	price := responseData[fromCurrency][toCurrency]

	return models.NewRate(from, to, price), nil
}

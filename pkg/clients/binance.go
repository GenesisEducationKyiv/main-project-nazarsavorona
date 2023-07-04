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
		apiURL       string
		fromCurrency string
		toCurrency   string
		client       HTTPClient
	}

	rateDTO struct {
		Price string `json:"price"`
	}
)

func NewBinanceClient(from, to, apiURL string, client HTTPClient) *BinanceClient {
	return &BinanceClient{
		apiURL:       apiURL,
		fromCurrency: from,
		toCurrency:   to,
		client:       client,
	}
}

func (g *BinanceClient) Rate(ctx context.Context) (*models.Rate, error) {
	url := fmt.Sprintf("%sticker/price?symbol=%s%s", g.apiURL, g.fromCurrency, g.toCurrency)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}

	response, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var r rateDTO

	err = json.NewDecoder(response.Body).Decode(&r)

	if err != nil {
		return nil, err
	}

	price, err := strconv.ParseFloat(r.Price, 64)
	if err != nil {
		return nil, err
	}

	return &models.Rate{
		From: g.fromCurrency,
		To:   g.toCurrency,
		Rate: price,
	}, nil
}

package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

	rate struct {
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

func (g *BinanceClient) Rate(ctx context.Context) (float64, error) {
	url := fmt.Sprintf("%sticker/price?symbol=%s%s", g.apiURL, g.fromCurrency, g.toCurrency)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return 0, err
	}

	response, err := g.client.Do(req)
	if err != nil {
		return 0, err
	}

	defer response.Body.Close()

	var r rate

	err = json.NewDecoder(response.Body).Decode(&r)

	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(r.Price, 64)
}

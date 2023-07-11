package clients

import (
	"context"
	"encoding/json"
	"errors"
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

func (g *BinanceClient) Rate(ctx context.Context, from, to string) (*models.Rate, error) {
	url := fmt.Sprintf("%sticker/price?symbol=%s%s", g.apiURL, from, to)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	response, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("request failed with status: " + response.Status)
	}

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
		From: from,
		To:   to,
		Rate: price,
	}, nil
}

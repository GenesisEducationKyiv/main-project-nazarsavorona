package currency

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type BinanceGetter struct {
	fromCurrency string
	toCurrency   string
}

type rate struct {
	Price string `json:"price"`
}

func NewGetter(from, to string) *BinanceGetter {
	return &BinanceGetter{
		fromCurrency: from,
		toCurrency:   to,
	}
}

const binanceAPIURL = "https://api.binance.com/api/v3/"

func (g *BinanceGetter) GetRate() (float64, error) {
	url := fmt.Sprintf("%sticker/price?symbol=%s%s", binanceAPIURL, g.fromCurrency, g.toCurrency)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)

	if err != nil {
		return 0, err
	}

	client := &http.Client{}
	response, err := client.Do(req)
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

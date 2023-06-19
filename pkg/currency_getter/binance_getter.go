package currency_getter

import (
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

func (g *BinanceGetter) GetRate() (float64, error) {
	response, err := http.Get(fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%s%s", g.fromCurrency, g.toCurrency))

	if err != nil {
		return 0, err
	}

	var r rate

	err = json.NewDecoder(response.Body).Decode(&r)

	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(r.Price, 64)
}

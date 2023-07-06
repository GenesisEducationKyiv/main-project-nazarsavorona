package funtional_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/smtp"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/clients"
	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/database"
	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/email"
	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/server"
	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/server/handlers"
	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/services"
	"github.com/stretchr/testify/require"
)

func TestRate(t *testing.T) {
	t.Parallel()

	filename := "test-rate.db"
	file, err := os.Create(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer func(name string) {
		err = file.Close()
		if err != nil {
			t.Fatal(err)
		}
		err = os.Remove(name)
		if err != nil {
			t.Fatal(err)
		}
	}(filename)

	s := prepareServer(t, file)

	testServer := httptest.NewServer(s)

	defer testServer.Close()

	url := testServer.URL + "/api/rate"
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	_, err = strconv.ParseFloat(string(body[:len(body)-1]), 64)
	require.NoError(t, err)
}

func TestSubscribe(t *testing.T) {
	t.Parallel()

	filename := "test-subscribe.db"
	file, err := os.Create(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer func(name string) {
		err = file.Close()
		if err != nil {
			t.Fatal(err)
		}
		err = os.Remove(name)
		if err != nil {
			t.Fatal(err)
		}
	}(filename)

	s := prepareServer(t, file)

	testServer := httptest.NewServer(s)

	defer testServer.Close()

	serverURL := testServer.URL + "/api/subscribe"
	formData := url.Values{}
	formData.Set("email", "test@email.com")

	resp, err := http.Post(serverURL,
		"application/x-www-form-urlencoded",
		strings.NewReader(formData.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = http.Post(serverURL,
		"application/x-www-form-urlencoded",
		strings.NewReader(formData.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	require.Equal(t, http.StatusConflict, resp.StatusCode)
}

func prepareServer(t *testing.T, file *os.File) *server.Server {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	senderEmail := ""
	senderPassword := ""

	fromCurrency := "BTC"
	toCurrency := "USDT"

	binanceURL := "https://api.binance.com/api/v3/"

	db := database.NewFileDatabase(file)
	if db == nil {
		t.Fatal("Error creating database")
	}

	repository := email.NewRepository(db)
	mailSender := email.NewSender(smtpHost, smtpPort, senderEmail, senderPassword,
		func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
			return nil
		})
	rateGetter := clients.NewBinanceClient(fromCurrency, toCurrency, binanceURL, &http.Client{})

	subscribeService := services.NewSubscribeService(repository)
	rateService := services.NewRateService(rateGetter)
	emailService := services.NewEmailService(mailSender)

	api := handlers.NewAPIHandlers(emailService, rateService, subscribeService)

	s := server.NewServer(api, nil)
	return s
}

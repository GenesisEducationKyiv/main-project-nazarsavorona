package main

import (
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/server"
	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/server/handlers"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/email"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/clients"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/database"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/services"
)

func main() {
	envValues, err := getEnvironmentValues()
	if err != nil {
		log.Panicf("Error fetching environment values: %v\n", err)
	}

	port := envValues["PORT"]

	smtpHost := envValues["SMTP_HOST"]
	smtpPort := envValues["SMTP_PORT"]

	senderEmail := envValues["EMAIL"]
	senderPassword := envValues["EMAIL_PASSWORD"]

	fromCurrency := envValues["FROM_CURRENCY"]
	toCurrency := envValues["TO_CURRENCY"]

	dbFileFolder := envValues["DB_FILE_FOLDER"]
	dbFileName := envValues["DB_FILE_NAME"]

	binanceURL := envValues["BINANCE_API_URL"]

	file, err := PrepareFile(dbFileFolder, dbFileName)
	if err != nil {
		log.Panicln(err.Error())
	}

	db := database.NewFileDatabase(file)
	if db == nil {
		log.Panicln("Error creating database")
	}

	repository := email.NewRepository(db)
	mailSender := email.NewSender(smtpHost, smtpPort, senderEmail, senderPassword, smtp.SendMail)
	rateGetter := clients.NewBinanceClient(fromCurrency, toCurrency, binanceURL, &http.Client{})

	subscribeService := services.NewSubscribeService(repository)
	rateService := services.NewRateService(rateGetter)
	emailService := services.NewEmailService(mailSender)

	api := handlers.NewAPIHandlers(emailService, rateService, subscribeService)
	web := handlers.NewWebHandlers(emailService, rateService, subscribeService)

	s := server.NewServer(api, web)

	log.Printf("Service listens port: %s", port)

	err = s.Start(":" + port)
	if err != nil {
		log.Panicln(err.Error())
	}
}

func getEnvironmentValues() (map[string]string, error) {
	requiredEnvVars := []string{
		"PORT",
		"SMTP_HOST",
		"SMTP_PORT",
		"EMAIL",
		"EMAIL_PASSWORD",
		"FROM_CURRENCY",
		"TO_CURRENCY",
		"DB_FILE_FOLDER",
		"DB_FILE_NAME",
		"BINANCE_API_URL",
	}

	envValues := make(map[string]string)

	for _, envVar := range requiredEnvVars {
		value := os.Getenv(envVar)
		if value == "" {
			return nil, fmt.Errorf("environment variable %s is not set", envVar)
		}
		envValues[envVar] = value
	}

	return envValues, nil
}

func PrepareFile(dbFileFolder string, dbFileName string) (*os.File, error) {
	err := os.MkdirAll(dbFileFolder, 0666)
	if err != nil {
		log.Panicln(err.Error())
	}

	file, err := os.OpenFile(dbFileFolder+dbFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Panicln(err.Error())
	}
	return file, err
}

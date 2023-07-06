package main

import (
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
	port := os.Getenv("PORT")

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	senderEmail := os.Getenv("EMAIL")
	senderPassword := os.Getenv("EMAIL_PASSWORD")

	fromCurrency := os.Getenv("FROM_CURRENCY")
	toCurrency := os.Getenv("TO_CURRENCY")

	dbFileFolder := os.Getenv("DB_FILE_FOLDER")
	dbFileName := os.Getenv("DB_FILE_NAME")

	binanceURL := os.Getenv("BINANCE_API_URL")

	err := os.MkdirAll(dbFileFolder, 0666)
	if err != nil {
		log.Panicln(err.Error())
	}

	file, err := os.OpenFile(dbFileFolder+dbFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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

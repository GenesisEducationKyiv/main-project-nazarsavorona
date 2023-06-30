package main

import (
	"log"
	"os"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/email"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/clients"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/application"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/database"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/service"
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
	mailSender := email.NewSender(smtpHost, smtpPort, senderEmail, senderPassword)
	rateGetter := clients.NewBinanceClient(fromCurrency, toCurrency, binanceURL)

	s := service.NewService(fromCurrency, toCurrency, repository, mailSender, rateGetter)
	app := application.NewApplication(s)

	log.Printf("Service listens port: %s", port)

	err = app.Run(":" + port)
	if err != nil {
		log.Panicln(err.Error())
	}
}

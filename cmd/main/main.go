package main

import (
	"log"
	"os"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/application"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/database"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/service"
)

func main() {
	port := os.Getenv("PORT")

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	email := os.Getenv("EMAIL")
	password := os.Getenv("EMAIL_PASSWORD")

	fromCurrency := os.Getenv("FROM_CURRENCY")
	toCurrency := os.Getenv("TO_CURRENCY")

	dbFileFolder := os.Getenv("DB_FILE_FOLDER")
	dbFilePath := os.Getenv("DB_FILE_PATH")

	err := os.MkdirAll(dbFileFolder, 0666)
	if err != nil {
		log.Panicln(err.Error())
	}

	file, err := os.OpenFile(dbFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Panicln(err.Error())
	}

	db := database.NewFileDatabase(file)
	if db == nil {
		log.Panicln("Error creating database")
	}

	s := service.NewService(smtpHost, smtpPort, email, password, fromCurrency, toCurrency, db)
	app := application.NewApplication(s)

	log.Printf("Service listens port: %s", port)

	err = app.Run(":" + port)
	if err != nil {
		log.Panicln(err.Error())
	}
}

package main

import (
	"github.com/nazarsavorona/btc-rate-check-service/pkg/application"
	"github.com/nazarsavorona/btc-rate-check-service/pkg/file_database"
	"github.com/nazarsavorona/btc-rate-check-service/pkg/service"
	"log"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	email := os.Getenv("EMAIL")
	password := os.Getenv("EMAIL_PASSWORD")

	fromCurrency := os.Getenv("FROM_CURRENCY")
	toCurrency := os.Getenv("TO_CURRENCY")

	dbFilePath := os.Getenv("DB_FILE_PATH")

	log.Printf("Service listens port: %s", port)

	s := service.NewService(email, password, fromCurrency, toCurrency, file_database.NewFileDatabase(dbFilePath))
	app := application.NewApplication(s)

	err := app.Run(":" + port)
	if err != nil {
		log.Panicln(err.Error())
	}
}

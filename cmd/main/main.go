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
	email := os.Getenv("EMAIL")
	password := os.Getenv("EMAIL_PASSWORD")

	fromCurrency := os.Getenv("FROM_CURRENCY")
	toCurrency := os.Getenv("TO_CURRENCY")

	dbFileFolder := os.Getenv("DB_FILE_FOLDER")
	dbFilePath := os.Getenv("DB_FILE_PATH")

	db := database.NewFileDatabase(dbFileFolder, dbFilePath)
	if db == nil {
		log.Panicln("Error creating database")
	}

	s := service.NewService(email, password, fromCurrency, toCurrency, db)
	app := application.NewApplication(s)

	log.Printf("Service listens port: %s", port)

	err := app.Run(":" + port)
	if err != nil {
		log.Panicln(err.Error())
	}
}

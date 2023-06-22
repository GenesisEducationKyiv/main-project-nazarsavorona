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
	emailHostURI := os.Getenv("EMAIL_HOST_URI")
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

	s := service.NewService(emailHostURI, email, password, fromCurrency, toCurrency, db)
	app := application.NewApplication(s)

	log.Printf("Service listens port: %s", port)

	err := app.Run(":" + port)
	if err != nil {
		log.Panicln(err.Error())
	}
}

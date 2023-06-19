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

	dbFilePath := os.Getenv("DB_FILE_PATH")

	log.Printf("Service listens port: %s", port)

	// create resources dir if not exists
	_, err := os.Stat("resources")
	if os.IsNotExist(err) {
		err = os.Mkdir("resources", 0777)
		if err != nil {
			log.Panicln(err.Error())
		}
	}

	s := service.NewService(email, password, fromCurrency, toCurrency, database.NewFileDatabase(dbFilePath))
	app := application.NewApplication(s)

	err = app.Run(":" + port)
	if err != nil {
		log.Panicln(err.Error())
	}
}

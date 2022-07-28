package main

import (
	"BTCRateCheckService/api"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	log.Printf("BTC Rate check service listens port: %s", port)

	server := api.NewServer(os.Getenv("EMAIL"), os.Getenv("EMAIL_PASSWORD"))
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), server)

	if err != nil {
		panic(err.Error())
	}
}

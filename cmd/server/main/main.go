package main

import (
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/internal/config"
	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/logger"
	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/clients/chain"
	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/server"
	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/server/handlers"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/clients"
	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/database"
	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/email"
	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/services"
)

func main() {
	conf, err := config.NewConfig()
	if err != nil {
		log.Panicln(err)
	}

	repository, err := prepareRepository(conf.DBFileFolder, conf.DBFileName)
	if err != nil {
		log.Panicln(err)
	}

	mailSender := email.NewSender(conf.Email,
		fmt.Sprintf("%s:%s", conf.SMTPHost, conf.SMTPPort),
		smtp.PlainAuth("", conf.Email, conf.EmailPassword, conf.SMTPHost))

	url := rabbitmq.ConstructRabbitMQURL(conf.RabbitMQHost, conf.RabbitMQPort,
		conf.RabbitMQUsername, conf.RabbitMQPassword)

	rabbitConn, err := amqp.Dial(url)
	if err != nil {
		log.Panicln(err)
	}

	defer rabbitConn.Close()

	rabbitChannel, err := rabbitConn.Channel()
	if err != nil {
		log.Panicln(err)
	}

	defer rabbitChannel.Close()

	rabbitLogger := rabbitmq.NewLogger(rabbitChannel)
	binanceChain := prepareClients(rabbitLogger, conf.BinanceAPIURL, conf.CoingeckoAPIURL)

	subscribeService := services.NewSubscribeService(repository)
	rateService := services.NewRateService(conf.FromCurrency, conf.ToCurrency, binanceChain)
	emailService := services.NewEmailService(mailSender, &email.HTMLMessageBuilder{})

	api := handlers.NewAPIHandlers(emailService, rateService, subscribeService)
	web := handlers.NewWebHandlers(emailService, rateService, subscribeService)

	s := server.NewServer(api, web)

	rabbitLogger.Log(logger.Info, fmt.Sprintf("Service listens on port: %s", conf.Port))

	err = s.Start(":" + conf.Port)
	if err != nil {
		log.Panicln(err)
	}
}

func prepareRepository(dbFileFolder string, dbFileName string) (*email.Repository, error) {
	file, err := prepareFile(dbFileFolder, dbFileName)
	if err != nil {
		return nil, fmt.Errorf("error preparing file: %w", err)
	}

	db := database.NewFileDatabase(file)
	if db == nil {
		return nil, fmt.Errorf("error creating database: %w", err)
	}

	repository := email.NewRepository(db)
	return repository, nil
}

func prepareFile(dbFileFolder string, dbFileName string) (*os.File, error) {
	err := os.MkdirAll(dbFileFolder, os.ModeDir)
	if err != nil {
		return nil, fmt.Errorf("error creating directory: %w", err)
	}

	file, err := os.OpenFile(dbFileFolder+"/"+dbFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}

	return file, nil
}

func prepareClients(logger logger.Logger, binanceURL string,
	coingeckoURL string) chain.Chain {
	binanceRateGetter := clients.NewLoggingClient("binance",
		clients.NewBinanceClient(binanceURL, &http.Client{}), logger)
	coingeckoRateGetter := clients.NewLoggingClient("coingecko",
		clients.NewCoingeckoClient(coingeckoURL, &http.Client{}), logger)

	binanceChain := chain.NewBaseChain(binanceRateGetter)
	coingeckoChain := chain.NewBaseChain(coingeckoRateGetter)
	binanceChain.SetNext(coingeckoChain)

	return binanceChain
}

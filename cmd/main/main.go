package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/smtp"
	"os"
	"time"

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

const defaultTimeout = 15

func main() {
	mode := flag.String("mode", "server", "mode to run the app or fetch logs")
	filter := flag.String("filter", "#", "filter to apply to the logs")
	timeout := flag.Int64("timeout", defaultTimeout, "timeout in seconds for the logs fetching")
	flag.Parse()

	switch *mode {
	case "server":
		err := startRateServer()
		if err != nil {
			log.Fatalf("Error starting rate server: %v", err)
		}
	case "logs":
		err := fetchLogs(*filter, *timeout)
		if err != nil {
			log.Fatalf("Error fetching logs: %v", err)
		}
	default:
		log.Fatalf("Unknown mode: %s", *mode)
	}
}

func fetchLogs(filter string, timeout int64) error {
	envValues, err := getEnvironmentValues()
	if err != nil {
		return fmt.Errorf("error fetching environment values: %w", err)
	}

	url := constructRabbitMQURL(envValues)
	rabbitConn, err := amqp.Dial(url)
	if err != nil {
		return fmt.Errorf("error connecting to rabbitmq: %w", err)
	}

	defer rabbitConn.Close()

	rabbitChannel, err := rabbitConn.Channel()
	if err != nil {
		return fmt.Errorf("error opening rabbitmq channel: %w", err)
	}

	defer rabbitChannel.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	consumer := rabbitmq.NewConsumer(rabbitChannel)
	consumer.ConsumeErrorLevelMsgs(ctx, filter, func(message string) {
		log.Printf("Received message: %s", message)
	})

	return nil
}

func startRateServer() error {
	envValues, err := getEnvironmentValues()
	if err != nil {
		return fmt.Errorf("error fetching environment values: %w", err)
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
	coingeckoURL := envValues["COINGECKO_API_URL"]

	repository, err := prepareRepository(dbFileFolder, dbFileName)
	if err != nil {
		return err
	}

	mailSender := email.NewSender(senderEmail,
		fmt.Sprintf("%s:%s", smtpHost, smtpPort),
		smtp.PlainAuth("", senderEmail, senderPassword, smtpHost))

	url := constructRabbitMQURL(envValues)
	rabbitConn, err := amqp.Dial(url)
	if err != nil {
		return fmt.Errorf("error connecting to rabbitmq: %w", err)
	}

	defer rabbitConn.Close()

	rabbitChannel, err := rabbitConn.Channel()
	if err != nil {
		return fmt.Errorf("error opening channel: %w", err)
	}

	defer rabbitChannel.Close()

	binanceRateGetter, coingeckoRateGetter := prepareClients(rabbitChannel, binanceURL, coingeckoURL)

	binanceChain := chain.NewBaseChain(binanceRateGetter)
	coingeckoChain := chain.NewBaseChain(coingeckoRateGetter)
	binanceChain.SetNext(coingeckoChain)

	subscribeService := services.NewSubscribeService(repository)
	rateService := services.NewRateService(fromCurrency, toCurrency, binanceChain)
	emailService := services.NewEmailService(mailSender, &email.HTMLMessageBuilder{})

	api := handlers.NewAPIHandlers(emailService, rateService, subscribeService)
	web := handlers.NewWebHandlers(emailService, rateService, subscribeService)

	s := server.NewServer(api, web)

	log.Printf("Service listens on port: %s", port)

	err = s.Start(":" + port)
	if err != nil {
		return fmt.Errorf("error starting server: %w", err)
	}

	return nil
}

func constructRabbitMQURL(envValues map[string]string) string {
	rabbitmqHost := envValues["RABBITMQ_HOST"]
	rabbitmqPort := envValues["RABBITMQ_PORT"]
	rabbitmqUsername := envValues["RABBITMQ_USERNAME"]
	rabbitmqPassword := envValues["RABBITMQ_PASSWORD"]

	rabbitHostPort := net.JoinHostPort(rabbitmqHost, rabbitmqPort)

	url := fmt.Sprintf("amqp://%s:%s@%s/",
		rabbitmqUsername, rabbitmqPassword, rabbitHostPort)
	return url
}

func prepareClients(rabbitChannel *amqp.Channel, binanceURL string,
	coingeckoURL string) (*clients.LoggingClient, *clients.LoggingClient) {
	logger := rabbitmq.NewLogger(rabbitChannel)

	binanceRateGetter := clients.NewLoggingClient("binance",
		clients.NewBinanceClient(binanceURL, &http.Client{}), logger)
	coingeckoRateGetter := clients.NewLoggingClient("coingecko",
		clients.NewCoingeckoClient(coingeckoURL, &http.Client{}), logger)
	return binanceRateGetter, coingeckoRateGetter
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
		"COINGECKO_API_URL",
		"RABBITMQ_HOST",
		"RABBITMQ_PORT",
		"RABBITMQ_USERNAME",
		"RABBITMQ_PASSWORD",
	}

	envValues := make(map[string]string)

	for _, envVar := range requiredEnvVars {
		value, exists := os.LookupEnv(envVar)
		if !exists {
			return nil, fmt.Errorf("environment variable %s is not set", envVar)
		}
		envValues[envVar] = value
	}

	return envValues, nil
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

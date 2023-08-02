package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port             string
	SMTPHost         string
	SMTPPort         string
	Email            string
	EmailPassword    string
	FromCurrency     string
	ToCurrency       string
	DBFileFolder     string
	DBFileName       string
	BinanceAPIURL    string
	CoingeckoAPIURL  string
	RabbitMQHost     string
	RabbitMQPort     string
	RabbitMQUsername string
	RabbitMQPassword string
}

func NewConfig() (*Config, error) {
	values, err := getEnvironmentValues()
	if err != nil {
		return nil, fmt.Errorf("error fetching environment values: %w", err)
	}

	return &Config{
		Port:             values["PORT"],
		SMTPHost:         values["SMTP_HOST"],
		SMTPPort:         values["SMTP_PORT"],
		Email:            values["EMAIL"],
		EmailPassword:    values["EMAIL_PASSWORD"],
		FromCurrency:     values["FROM_CURRENCY"],
		ToCurrency:       values["TO_CURRENCY"],
		DBFileFolder:     values["DB_FILE_FOLDER"],
		DBFileName:       values["DB_FILE_NAME"],
		BinanceAPIURL:    values["BINANCE_API_URL"],
		CoingeckoAPIURL:  values["COINGECKO_API_URL"],
		RabbitMQHost:     values["RABBITMQ_HOST"],
		RabbitMQPort:     values["RABBITMQ_PORT"],
		RabbitMQUsername: values["RABBITMQ_USERNAME"],
		RabbitMQPassword: values["RABBITMQ_PASSWORD"],
	}, nil
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

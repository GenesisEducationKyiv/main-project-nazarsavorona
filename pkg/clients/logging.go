package clients

import (
	"context"
	"log"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/models"
	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/services"
)

type LoggingClient struct {
	name        string
	rateFetcher services.RateFetcher
}

func NewLoggingClient(name string, rateFetcher services.RateFetcher) *LoggingClient {
	return &LoggingClient{
		name:        name,
		rateFetcher: rateFetcher,
	}
}

func (l *LoggingClient) Rate(ctx context.Context, from, to string) (*models.Rate, error) {
	rate, err := l.rateFetcher.Rate(ctx, from, to)
	if err != nil {
		log.Printf("%s: error: %v", l.name, err)
		return nil, err
	}

	log.Printf("%s: rate: %+v", l.name, *rate)
	return rate, nil
}

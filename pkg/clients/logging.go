package clients

import (
	"context"
	"fmt"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/logger"
	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/services"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/models"
)

type LoggingClient struct {
	name        string
	rateFetcher services.RateFetcher
	log         logger.Logger
}

func NewLoggingClient(name string, rateFetcher services.RateFetcher, l logger.Logger) *LoggingClient {
	return &LoggingClient{
		name:        name,
		rateFetcher: rateFetcher,
		log:         l,
	}
}

func (l *LoggingClient) Rate(ctx context.Context, from, to string) (*models.Rate, error) {
	rate, err := l.rateFetcher.Rate(ctx, from, to)
	if err != nil {
		l.log.Log(logger.Error, fmt.Sprintf("%s client: %s", l.name, err))
		return nil, err
	}

	l.log.Log(logger.Info, fmt.Sprintf("%s client: rate: %+v", l.name, *rate))

	return rate, nil
}

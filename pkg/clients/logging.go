package clients

import (
	"context"
	"fmt"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/services"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/models"
)

type Logger interface {
	Log(level models.Level, message string)
}

type LoggingClient struct {
	name        string
	rateFetcher services.RateFetcher
	log         Logger
}

func NewLoggingClient(name string, rateFetcher services.RateFetcher, l Logger) *LoggingClient {
	return &LoggingClient{
		name:        name,
		rateFetcher: rateFetcher,
		log:         l,
	}
}

func (l *LoggingClient) Rate(ctx context.Context, from, to string) (*models.Rate, error) {
	rate, err := l.rateFetcher.Rate(ctx, from, to)
	if err != nil {
		l.log.Log(models.Error, fmt.Sprintf("%s client: %s", l.name, err))
		return nil, err
	}

	l.log.Log(models.Info, fmt.Sprintf("%s client: rate: %+v", l.name, *rate))

	return rate, nil
}

package application

import (
	"github.com/nazarsavorona/btc-rate-check-service/pkg/http_server"
	"github.com/nazarsavorona/btc-rate-check-service/pkg/service"
)

type Application struct {
	server *http_server.Server
}

func NewApplication(s *service.Service) *Application {
	return &Application{
		server: http_server.NewServer(s),
	}
}

func (a *Application) Run(address string) error {
	return a.server.Start(address)
}

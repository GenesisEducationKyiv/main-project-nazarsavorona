package application

import (
	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/service"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/http"
)

type Application struct {
	server *http.Server
}

func NewApplication(s *service.Service) *Application {
	return &Application{
		server: http.NewServer(s),
	}
}

func (a *Application) Run(address string) error {
	return a.server.Start(address)
}

package server

import (
	"github.com/labstack/echo/v4"
)

type (
	apiHandlers interface {
		Rate(c echo.Context) error
		Subscribe(c echo.Context) error
		SendEmails(c echo.Context) error
	}

	webHandlers interface {
		Index(c echo.Context) error
		Subscribe(c echo.Context) error
		SendEmails(c echo.Context) error
		Conflict(c echo.Context) error
	}

	Server struct {
		router *echo.Echo

		api apiHandlers
		web webHandlers
	}
)

func NewServer(api apiHandlers, web webHandlers) *Server {
	e := echo.New()
	e.HideBanner = true

	s := &Server{
		router: e,
		api:    api,
		web:    web,
	}

	s.routes()

	return s
}

func (s *Server) routes() {
	s.router.GET("/api/rate", s.api.Rate)
	s.router.POST("/api/subscribe", s.api.Subscribe)
	s.router.POST("/api/sendEmails", s.api.SendEmails)

	s.router.GET("/", s.web.Index)
	s.router.POST("/subscribe", s.web.Subscribe)
	s.router.POST("/sendEmails", s.web.SendEmails)
	s.router.GET("/conflict", s.web.Conflict)
}

func (s *Server) Start(address string) error {
	return s.router.Start(address)
}

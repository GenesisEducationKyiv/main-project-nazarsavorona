package http

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strings"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/models"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/services"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type emailService interface {
	SendEmails(context.Context, []string, *models.Message) error
}

type rateService interface {
	Rate(ctx context.Context) (*models.Rate, error)
}

type subscribeService interface {
	Subscribe(email string) error
	EmailList() []string
}

type Server struct {
	router *echo.Echo

	emailService     emailService
	rateService      rateService
	subscribeService subscribeService

	template *template.Template
}

func NewServer(emailService emailService,
	rateService rateService,
	subscribeService subscribeService) *Server {
	functionMap := template.FuncMap{"add": func(x, y int) int { return x + y }}

	e := echo.New()
	e.HideBanner = true

	server := &Server{
		router: e,

		emailService:     emailService,
		rateService:      rateService,
		subscribeService: subscribeService,

		template: template.Must(template.New("").Funcs(functionMap).ParseGlob("./templates/*.gohtml")),
	}

	e.Use(middleware.Recover(), middleware.Logger())

	server.routes()

	return server
}

func (s *Server) Start(address string) error {
	return s.router.Start(address)
}

func (s *Server) routes() {
	s.router.GET("/api/rate", s.rate)
	s.router.POST("/api/subscribe", s.apiSubscribe)
	s.router.POST("/api/sendEmails", s.sendEmails)

	s.router.GET("/", s.index)
	s.router.POST("/subscribe", s.webSubscribe)
	s.router.POST("/sendEmails", s.webSendEmails)
	s.router.GET("/conflict", s.conflict)
}

func (s *Server) index(c echo.Context) error {
	emails := s.subscribeService.EmailList()
	sort.Strings(emails)

	r, err := s.rateService.Rate(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	indexData := struct {
		Rate   string
		Emails []string
	}{fmt.Sprintf("%.2f", r.Rate), emails}

	err = s.template.ExecuteTemplate(c.Response().Writer, "index.gohtml", indexData)
	if err != nil {
		http.Redirect(c.Response().Writer, c.Request(), "/", http.StatusInternalServerError)
	}

	return nil
}

func (s *Server) conflict(c echo.Context) error {
	return s.template.ExecuteTemplate(c.Response().Writer, "conflict.gohtml", nil)
}

var greetingsMessage = &models.Message{
	Subject: "Subscription",
	Body:    "You have successfully subscribed to the service",
}

func (s *Server) webSubscribe(c echo.Context) error {
	email := extractEmail(c)

	err := s.subscribeService.Subscribe(email)
	if err != nil {
		if errors.Is(err, services.ErrAlreadySubscribed) {
			http.Redirect(c.Response().Writer, c.Request(), "/conflict", http.StatusSeeOther)
			return nil
		}
		http.Redirect(c.Response().Writer, c.Request(), "/", http.StatusBadRequest)
	}

	err = s.emailService.SendEmails(c.Request().Context(), []string{email}, greetingsMessage)
	if err != nil {
		http.Redirect(c.Response().Writer, c.Request(), "/", http.StatusBadRequest)
		return err
	}

	http.Redirect(c.Response().Writer, c.Request(), "/", http.StatusSeeOther)
	return nil
}

func (s *Server) webSendEmails(c echo.Context) error {
	r, err := s.rateService.Rate(c.Request().Context())
	if err != nil {
		http.Redirect(c.Response().Writer, c.Request(), "/", http.StatusBadRequest)
		return err
	}

	err = s.emailService.SendEmails(c.Request().Context(), s.subscribeService.EmailList(), models.NewMessage(r))
	if err != nil {
		http.Redirect(c.Response().Writer, c.Request(), "/", http.StatusBadRequest)
		return err
	}

	http.Redirect(c.Response().Writer, c.Request(), "/", http.StatusSeeOther)
	return nil
}

func (s *Server) rate(c echo.Context) error {
	r, err := s.rateService.Rate(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, r.Rate)
}

func (s *Server) apiSubscribe(c echo.Context) error {
	email := extractEmail(c)

	err := s.subscribeService.Subscribe(email)
	if err != nil {
		if errors.Is(err, services.ErrAlreadySubscribed) {
			return c.JSON(http.StatusConflict, err.Error())
		}
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err = s.emailService.SendEmails(c.Request().Context(), []string{email}, greetingsMessage)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, email)
}

func (s *Server) sendEmails(c echo.Context) error {
	r, err := s.rateService.Rate(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err = s.emailService.SendEmails(c.Request().Context(), s.subscribeService.EmailList(), models.NewMessage(r))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, "Emails sent")
}

func extractEmail(c echo.Context) string {
	email := c.FormValue("email")

	return strings.TrimSpace(email)
}

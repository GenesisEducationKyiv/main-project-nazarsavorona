package http

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type service interface {
	Subscribe(email string) error
	SendEmails(ctx context.Context) error
	Rate(ctx context.Context) (float64, error)
	EmailList() []string
}

type Server struct {
	router  *echo.Echo
	service service

	template *template.Template
}

func NewServer(s service) *Server {
	functionMap := template.FuncMap{"add": func(x, y int) int { return x + y }}

	e := echo.New()
	e.HideBanner = true

	server := &Server{
		router:  e,
		service: s,

		template: template.Must(template.New("").Funcs(functionMap).ParseGlob("./templates/*.gohtml")),
	}

	e.Use(middleware.Recover(), middleware.Logger())

	server.routes()

	return server
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

func (s *Server) rate(c echo.Context) error {
	rate, err := s.service.Rate(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, rate)
}

func (s *Server) apiSubscribe(c echo.Context) error {
	email := getEmail(c)

	err := s.service.Subscribe(email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, email)
}

func (s *Server) sendEmails(c echo.Context) error {
	_ = s.service.SendEmails(c.Request().Context())

	return c.JSON(http.StatusOK, "Emails sent")
}

func (s *Server) Start(address string) error {
	return s.router.Start(address)
}

func (s *Server) index(c echo.Context) error {
	emails := s.service.EmailList()
	sort.Strings(emails)

	rate, err := s.service.Rate(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	indexData := struct {
		Rate   string
		Emails []string
	}{fmt.Sprintf("%.2f", rate), emails}

	err = s.template.ExecuteTemplate(c.Response().Writer, "index.gohtml", indexData)
	if err != nil {
		http.Redirect(c.Response().Writer, c.Request(), "/", http.StatusInternalServerError)
	}

	return nil
}

func (s *Server) conflict(c echo.Context) error {
	return s.template.ExecuteTemplate(c.Response().Writer, "conflict.gohtml", nil)
}

func (s *Server) webSubscribe(c echo.Context) error {
	email := getEmail(c)

	err := s.service.Subscribe(email)
	if err != nil {
		http.Redirect(c.Response().Writer, c.Request(), "/conflict", http.StatusSeeOther)
	}

	http.Redirect(c.Response().Writer, c.Request(), "/", http.StatusSeeOther)
	return nil
}

func (s *Server) webSendEmails(c echo.Context) error {
	err := s.service.SendEmails(c.Request().Context())
	if err != nil {
		http.Redirect(c.Response().Writer, c.Request(), "/", http.StatusBadRequest)
		return err
	}

	http.Redirect(c.Response().Writer, c.Request(), "/", http.StatusSeeOther)
	return nil
}

func getEmail(c echo.Context) string {
	email := c.FormValue("email")

	return strings.TrimSpace(email)
}

package handlers

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strings"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/models"
	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/services"
	"github.com/labstack/echo/v4"
)

type WebHandlers struct {
	emailService     emailService
	rateService      rateService
	subscribeService subscribeService

	template *template.Template
}

func NewWebHandlers(emailService emailService,
	rateService rateService,
	subscribeService subscribeService) *WebHandlers {
	functionMap := template.FuncMap{"add": func(x, y int) int { return x + y }}

	server := &WebHandlers{
		emailService:     emailService,
		rateService:      rateService,
		subscribeService: subscribeService,

		template: template.Must(template.New("").Funcs(functionMap).ParseGlob("./templates/*.gohtml")),
	}

	return server
}

func (h *WebHandlers) Index(c echo.Context) error {
	emails := h.subscribeService.EmailList()
	sort.Strings(emails)

	r, err := h.rateService.Rate(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	indexData := struct {
		Rate   string
		Emails []string
	}{fmt.Sprintf("%.2f", r.Rate), emails}

	err = h.template.ExecuteTemplate(c.Response().Writer, "index.gohtml", indexData)
	if err != nil {
		http.Redirect(c.Response().Writer, c.Request(), "/", http.StatusInternalServerError)
	}

	return nil
}

func (h *WebHandlers) Conflict(c echo.Context) error {
	return h.template.ExecuteTemplate(c.Response().Writer, "conflict.gohtml", nil)
}

func (h *WebHandlers) Subscribe(c echo.Context) error {
	email := strings.TrimSpace(c.FormValue("email"))

	err := h.subscribeService.Subscribe(email)
	if err != nil {
		if errors.Is(err, services.ErrAlreadySubscribed) {
			http.Redirect(c.Response().Writer, c.Request(), "/conflict", http.StatusSeeOther)
			return nil
		}
		http.Redirect(c.Response().Writer, c.Request(), "/", http.StatusBadRequest)
	}

	err = h.emailService.SendEmails(c.Request().Context(), []string{email}, &models.Message{
		Subject: "Subscription",
		Body:    "You have successfully subscribed to the service",
	})
	if err != nil {
		http.Redirect(c.Response().Writer, c.Request(), "/", http.StatusBadRequest)
		return err
	}

	http.Redirect(c.Response().Writer, c.Request(), "/", http.StatusSeeOther)
	return nil
}

func (h *WebHandlers) SendEmails(c echo.Context) error {
	r, err := h.rateService.Rate(c.Request().Context())
	if err != nil {
		http.Redirect(c.Response().Writer, c.Request(), "/", http.StatusBadRequest)
		return err
	}

	err = h.emailService.SendEmails(c.Request().Context(), h.subscribeService.EmailList(), models.NewMessage(r))
	if err != nil {
		http.Redirect(c.Response().Writer, c.Request(), "/", http.StatusBadRequest)
		return err
	}

	http.Redirect(c.Response().Writer, c.Request(), "/", http.StatusSeeOther)
	return nil
}

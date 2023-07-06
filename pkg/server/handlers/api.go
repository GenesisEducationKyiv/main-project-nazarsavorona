package handlers

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/models"
	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/services"
	"github.com/labstack/echo/v4"
)

type (
	emailService interface {
		SendEmails(context.Context, []string, *models.Message) error
	}

	rateService interface {
		Rate(ctx context.Context) (*models.Rate, error)
	}

	subscribeService interface {
		Subscribe(email string) error
		EmailList() []string
	}

	APIHandlers struct {
		emailService     emailService
		rateService      rateService
		subscribeService subscribeService
	}
)

func NewAPIHandlers(emailService emailService,
	rateService rateService,
	subscribeService subscribeService) *APIHandlers {
	return &APIHandlers{
		emailService:     emailService,
		rateService:      rateService,
		subscribeService: subscribeService,
	}
}

func (h *APIHandlers) Rate(c echo.Context) error {
	r, err := h.rateService.Rate(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, r.Rate)
}

func (h *APIHandlers) Subscribe(c echo.Context) error {
	email := strings.TrimSpace(c.FormValue("email"))

	err := h.subscribeService.Subscribe(email)
	if err != nil {
		if errors.Is(err, services.ErrAlreadySubscribed) {
			return c.JSON(http.StatusConflict, err.Error())
		}
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err = h.emailService.SendEmails(c.Request().Context(), []string{email}, &models.Message{
		Subject: "Subscription",
		Body:    "You have successfully subscribed to the service",
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, email)
}

func (h *APIHandlers) SendEmails(c echo.Context) error {
	r, err := h.rateService.Rate(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err = h.emailService.SendEmails(c.Request().Context(), h.subscribeService.EmailList(), models.NewMessageFromRate(r))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, "Emails sent")
}

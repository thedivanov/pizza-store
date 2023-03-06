package handlers

import (
	"context"
	"delivery/models"
	"net/http"

	"github.com/labstack/echo"
)

type Handler struct {
	E        *echo.Echo
	rabbitMQ RabbitMQ
	db       DB
}

type RabbitMQ interface {
	SendDeliveredMessage(order models.Order) error
}

type DB interface {
	SetStatusDelivered(ctx context.Context, id int) (*models.Order, error)
}

func NewHandler(e *echo.Echo, rabbitMQ RabbitMQ, db DB) (*Handler, error) {
	return &Handler{
		E:        e,
		rabbitMQ: rabbitMQ,
		db:       db,
	}, nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.E.ServeHTTP(w, r)
}

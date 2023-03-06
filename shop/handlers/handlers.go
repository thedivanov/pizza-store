package handlers

import (
	"context"
	"net/http"

	"shop/models"

	"github.com/labstack/echo"
)

type Handler struct {
	E        *echo.Echo
	rabbitMQ RabbitMQ
	db       DB
}

type RabbitMQ interface {
	CreateOreder(order models.Order) error
}

type DB interface {
	GetOrders(ctx context.Context, offset, limit int, needTotal bool) ([]*models.Order, int64, error)
	CreateOreder(ctx context.Context, order *models.Order) error
	SetCancelOrder(ctx context.Context, id int) (*models.Order, error)
	SetConfirmOrder(ctx context.Context, id int) (*models.Order, error)
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

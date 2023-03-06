package handlers

import (
	"context"
	"net/http"

	"kitchen/models"

	"github.com/labstack/echo"
)

type Handler struct {
	E        *echo.Echo
	RabbitMQ RabbitMQ
	db       DB
}

type RabbitMQ interface {
	CreateHandoverOrder(order models.Order) error
}

type DB interface {
	GetOrders(ctx context.Context, limit, offset int, needTotal bool) ([]*models.KitchenOrder, int64, error)
	CreateOrder(ctx context.Context, order models.Order) error
	SetCookingOrder(ctx context.Context, id int) (*models.Order, error)
	SetCookedOrder(ctx context.Context, id int) (*models.Order, error)
	SetHandoverOrder(ctx context.Context, id int) (*models.Order, error)
}

func NewHandler(e *echo.Echo, rabbitMQ RabbitMQ, db DB) (*Handler, error) {
	return &Handler{
		E:        e,
		RabbitMQ: rabbitMQ,
		db:       db,
	}, nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.E.ServeHTTP(w, r)
}

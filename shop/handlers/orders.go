package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"shop/models"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func (h *Handler) createOreder(c echo.Context) error {
	order := models.Order{}

	err := json.NewDecoder(c.Request().Body).Decode(&order)
	if err != nil {
		log.WithError(err).Error("Decode body error")
		return echo.ErrInternalServerError
	}
	if err = validator.New().Struct(order); err != nil {
		log.WithError(err).Error("Validate body error")
		return echo.ErrBadRequest
	}
	err = h.db.CreateOreder(context.TODO(), &order)
	if err != nil {
		log.WithError(err).Error("Write to DB error")
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, order)
}

func (h *Handler) getOrders(c echo.Context) error {
	orders, _, err := h.db.GetOrders(context.TODO(), 0, 0, false)
	if err != nil {
		log.WithError(err).Error("Get orders from db error")
		return echo.ErrInternalServerError
	}

	ordersResponse := models.GetOrdersResponse{}
	ordersResponse.Orders = orders

	return c.JSON(http.StatusOK, ordersResponse)
}

func (h *Handler) adminGetOrders(c echo.Context) error {
	orders, count, err := h.db.GetOrders(context.TODO(), 0, 0, true)
	if err != nil {
		log.WithError(err).Error("Get orders from db error")
		return echo.ErrInternalServerError
	}

	ordersResponse := models.GetOrdersAdminResponse{}
	ordersResponse.Orders = orders
	ordersResponse.Total = count
	ordersResponse.Offset = 0

	return c.JSON(http.StatusOK, ordersResponse)
}

func (h *Handler) confirmOrder(c echo.Context) error {
	orderIDParam := c.Param("order_id")

	orderID, err := strconv.ParseInt(orderIDParam, 10, 64)
	if err != nil {
		log.WithError(err).Error("Parse order_id error")
		return echo.ErrInternalServerError
	}
	if orderID <= 0 {
		log.WithError(err).Error("Bad order_id error")
		return echo.ErrBadRequest
	}

	order, err := h.db.SetConfirmOrder(context.TODO(), int(orderID))
	if err != nil {
		log.WithError(err).Error("Confirm order error")
		return echo.ErrInternalServerError
	}

	err = h.rabbitMQ.CreateOreder(*order)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, nil)
}

func (h *Handler) cancelOrder(c echo.Context) error {
	orderIDParam := c.Param("order_id")

	orderID, err := strconv.ParseInt(orderIDParam, 10, 64)
	if err != nil {
		log.WithError(err).Error("Parse order_id error")
		return echo.ErrInternalServerError
	}
	if orderID <= 0 {
		log.WithError(err).Error("Bad order_id error")
		return echo.ErrBadRequest
	}

	_, err = h.db.SetCancelOrder(context.TODO(), int(orderID))
	if err != nil {
		log.WithError(err).Error("Set canceled status update error")
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, nil)
}

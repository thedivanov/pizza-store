package handlers

import (
	"context"
	"kitchen/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func (h *Handler) getOrders(c echo.Context) error {
	orders, count, err := h.db.GetOrders(context.TODO(), 0, 0, true)
	if err != nil {
		log.WithError(err).Error("Get orders from db error")
		return echo.ErrInternalServerError
	}

	ordersResponse := models.OrdersResponse{}
	ordersResponseOrder := []*models.OrdersResponseOrder{}
	for _, order := range orders {
		respOrder := &models.OrdersResponseOrder{}
		respOrder.ID = order.ID
		respOrder.Status = order.Status
		respOrder.Items = append(respOrder.Items, order.Order.Items...)

		ordersResponseOrder = append(ordersResponseOrder, respOrder)
	}
	ordersResponse.Orders = ordersResponseOrder
	ordersResponse.Total = count

	return c.JSON(http.StatusOK, ordersResponse)
}

func (h *Handler) cookingStart(c echo.Context) error {
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

	_, err = h.db.SetCookingOrder(context.TODO(), int(orderID))
	if err != nil {
		log.WithError(err).Error("Confirm order error")
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, nil)
}

func (h *Handler) cookingEnd(c echo.Context) error {
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

	_, err = h.db.SetCookedOrder(context.TODO(), int(orderID))
	if err != nil {
		log.WithError(err).Error("Confirm order error")
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, nil)
}

func (h *Handler) cookingHandover(c echo.Context) error {
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

	order, err := h.db.SetHandoverOrder(context.TODO(), int(orderID))
	if err != nil {
		log.WithError(err).Error("Confirm order error")
		return echo.ErrInternalServerError
	}

	err = h.RabbitMQ.CreateHandoverOrder(*order)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, nil)
}

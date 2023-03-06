package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func (h *Handler) deliver(c echo.Context) error {
	orderIDParam := c.Param("order_id")

	orderID, err := strconv.ParseInt(orderIDParam, 10, 64)
	if err != nil {
		return echo.ErrInternalServerError
	}
	if orderID <= 0 {
		log.Error("Bad order_id error")
		return echo.ErrBadRequest
	}

	order, err := h.db.SetStatusDelivered(context.TODO(), int(orderID))
	if err != nil {
		return echo.ErrInternalServerError
	}

	err = h.rabbitMQ.SendDeliveredMessage(*order)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}

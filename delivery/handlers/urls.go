package handlers

import "github.com/labstack/echo"

func (h *Handler) AddURLs() {
	v1 := h.E.Group("/rider/v1")
	h.groupV1Routes(v1)
}

func (h *Handler) groupV1Routes(v1Group *echo.Group) {
	v1Group.POST("/orders/:order_id/deliver", h.deliver)
}

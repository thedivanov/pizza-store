package handlers

import "github.com/labstack/echo"

func (h *Handler) AddURLs() {
	v1 := h.E.Group("/v1")
	h.groupV1Routes(v1)
}

func (h *Handler) groupV1Routes(v1Group *echo.Group) {
	v1Group.POST("/orders", h.createOreder)
	v1Group.GET("/orders", h.getOrders)
	v1Group.GET("/admin/orders", h.adminGetOrders)
	v1Group.POST("/admin/orders/:order_id/confirm", h.confirmOrder)
	v1Group.POST("/admin/orders/:order_id/cancel", h.cancelOrder)
}

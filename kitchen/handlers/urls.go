package handlers

import "github.com/labstack/echo"

func (h *Handler) AddURLs() {
	v1 := h.E.Group("/cook/v1")
	h.groupV1Routes(v1)
}

func (h *Handler) groupV1Routes(v1Group *echo.Group) {
	v1Group.GET("/orders", h.getOrders)
	v1Group.POST("/orders/:order_id/cooking/start", h.cookingStart)
	v1Group.POST("/orders/:order_id/cooking/end", h.cookingEnd)
	v1Group.POST("/orders/:order_id/handover", h.cookingHandover)
}

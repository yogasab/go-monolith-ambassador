package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/yogasab/go-monolith-ambassador/src/services"
)

type orderController struct {
	orderService services.OrderService
}

func NewrOrderController(orderService services.OrderService) *orderController {
	return &orderController{orderService: orderService}
}

func (h *orderController) GetOrders(ctx *fiber.Ctx) error {
	orders, err := h.orderService.GetOrders()
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).
			JSON(fiber.Map{
				"code":    http.StatusInternalServerError,
				"message": "internal server errors",
				"error":   err.Error(),
			})
	}

	for i, order := range orders {
		orders[i].Name = order.GetFullName()
		orders[i].Total = order.GetTotalPrice()
	}

	return ctx.
		Status(http.StatusOK).
		JSON(fiber.Map{
			"code":    http.StatusOK,
			"message": "orders fetched successfully",
			"data":    orders,
		})
}

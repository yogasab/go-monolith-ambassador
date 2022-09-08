package controllers

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/yogasab/go-monolith-ambassador/src/models/dto"
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

func (h *orderController) GetOrdersRankings(ctx *fiber.Ctx) error {
	var formatters []dto.OrderRankingDTO
	results, err := h.orderService.GetOrdersRankings()

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).
			JSON(fiber.Map{
				"code":    http.StatusInternalServerError,
				"message": "internal server errors",
				"error":   err.Error(),
			})
	}

	bytes, _ := json.Marshal(results)

	if err := json.Unmarshal(bytes, &formatters); err != nil {
		return ctx.Status(http.StatusInternalServerError).
			JSON(fiber.Map{
				"code":    http.StatusInternalServerError,
				"message": "internal server errors",
				"error":   err.Error(),
			})
	}

	if rankingParam := ctx.Query("rank"); rankingParam != "" {
		loweredParam := strings.ToLower(rankingParam)
		if loweredParam == "lowest" {
			sort.Slice(formatters, func(i, j int) bool {
				return formatters[i].Revenue < formatters[j].Revenue
			})
		} else {
			sort.Slice(formatters, func(i, j int) bool {
				return formatters[i].Revenue > formatters[j].Revenue
			})
		}
	}

	return ctx.
		Status(http.StatusOK).
		JSON(fiber.Map{
			"code":    http.StatusOK,
			"message": "order rankings fetched successfully",
			"data":    formatters,
		})
}

func (h *orderController) CreateOrder(ctx *fiber.Ctx) error {
	dto := new(dto.CreateOrderDTO)
	if err := ctx.BodyParser(dto); err != nil {
		return ctx.
			Status(http.StatusUnprocessableEntity).
			JSON(fiber.Map{
				"code":    http.StatusUnprocessableEntity,
				"message": "failed to process request",
				"error":   err,
			})
	}

	if errors := ValidateInput(*dto); errors != nil {
		return ctx.
			Status(http.StatusUnprocessableEntity).
			JSON(fiber.Map{
				"code":    http.StatusUnprocessableEntity,
				"message": "error validation request",
				"error":   errors,
			})
	}

	newOrder, err := h.orderService.CreateOrder(dto)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).
			JSON(fiber.Map{
				"code":    http.StatusInternalServerError,
				"message": "internal server errors",
				"error":   err.Error(),
			})
	}

	return ctx.
		Status(http.StatusCreated).
		JSON(fiber.Map{
			"code":    http.StatusCreated,
			"message": "order created successfully",
			"data":    newOrder,
		})
}

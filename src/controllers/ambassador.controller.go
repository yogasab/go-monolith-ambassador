package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/yogasab/go-monolith-ambassador/src/services"
)

type ambassadorController struct {
	ambassadorService services.AmbassadorService
}

func NewAmbassadorController(ambassadorService services.AmbassadorService) *ambassadorController {
	return &ambassadorController{ambassadorService: ambassadorService}
}

func (h *ambassadorController) GetAmbassadors(ctx *fiber.Ctx) error {
	ambassadors, err := h.ambassadorService.GetAmbassadors()
	if err != nil {
		return ctx.
			Status(http.StatusInternalServerError).
			JSON(fiber.Map{
				"code":    http.StatusInternalServerError,
				"message": "internal server error",
				"error":   err.Error(),
			})
	}
	return ctx.
		Status(http.StatusOK).
		JSON(fiber.Map{
			"code":    http.StatusOK,
			"message": "ambassadors fetched successfully",
			"data":    ambassadors,
		})
}

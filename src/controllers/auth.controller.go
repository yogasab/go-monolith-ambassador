package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/yogasab/go-monolith-ambassador/src/models/dto"
	"github.com/yogasab/go-monolith-ambassador/src/services"
)

type authController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) *authController {
	return &authController{authService: authService}
}

func (h *authController) RegisterAdmin(ctx *fiber.Ctx) error {
	dto := new(dto.RegisterDTO)
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

	newUser, err := h.authService.RegisterAdmin(dto)
	if err != nil {
		return ctx.
			Status(http.StatusInternalServerError).
			JSON(fiber.Map{
				"code":    http.StatusInternalServerError,
				"message": "internal server error",
				"error":   err,
			})
	}

	return ctx.
		Status(http.StatusCreated).
		JSON(fiber.Map{
			"code":    http.StatusCreated,
			"message": "user registered successfully",
			"data":    newUser,
		})
}

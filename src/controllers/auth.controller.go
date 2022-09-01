package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/yogasab/go-monolith-ambassador/src/models/dto"
	"github.com/yogasab/go-monolith-ambassador/src/services"
)

type authController struct {
	authService  services.AuthService
	tokenService services.TokenService
}

func NewAuthController(
	authService services.AuthService, tokenService services.TokenService) *authController {
	return &authController{
		authService:  authService,
		tokenService: tokenService,
	}
}

func (h *authController) Register(ctx *fiber.Ctx) error {
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

	newUser, err := h.authService.Register(dto)
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

func (h *authController) Login(ctx *fiber.Ctx) error {
	dto := new(dto.LoginDTO)
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

	registeredUser, err := h.authService.Login(dto)
	if err != nil {
		return ctx.
			Status(http.StatusInternalServerError).
			JSON(fiber.Map{
				"code":    http.StatusInternalServerError,
				"message": "internal server error",
				"error":   err.Error(),
			})
	}

	token, err := h.tokenService.GenerateToken(registeredUser.ID)
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
			"message": "user logged in successfully",
			"data":    token,
		})
}

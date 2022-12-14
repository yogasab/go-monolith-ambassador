package controllers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/yogasab/go-monolith-ambassador/src/middlewares"
	"github.com/yogasab/go-monolith-ambassador/src/models/dto"
	"github.com/yogasab/go-monolith-ambassador/src/services"
)

type authController struct {
	authService  services.AuthService
	tokenService services.TokenService
	orderService services.OrderService
}

func NewAuthController(
	authService services.AuthService,
	tokenService services.TokenService,
	orderService services.OrderService) *authController {
	return &authController{
		authService:  authService,
		tokenService: tokenService,
		orderService: orderService,
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

	dto.IsAmbassador = strings.Contains(ctx.Path(), "/api/ambassadors")

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
	scope := "admin"
	isAmbassador := strings.Contains(ctx.Path(), "/api/ambassadors")
	if isAmbassador {
		scope = "ambassadors"
	}

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

	// Prohibited if the ambassador access the admin endpoint
	if !isAmbassador && registeredUser.IsAmbassador {
		return ctx.
			Status(http.StatusForbidden).
			JSON(fiber.Map{
				"code":    http.StatusForbidden,
				"message": "prohibited to aceess this route",
				"errors":  "you are prohibited to acesss this route",
			})
	}

	token, err := h.tokenService.GenerateToken(registeredUser.ID, scope)
	if err != nil {
		return ctx.
			Status(http.StatusInternalServerError).
			JSON(fiber.Map{
				"code":    http.StatusInternalServerError,
				"message": "internal server error",
				"error":   err.Error(),
			})
	}

	loggedUser := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		HTTPOnly: true,
		Expires:  time.Now().Add(time.Hour * 24),
	}
	ctx.Cookie(&loggedUser)

	return ctx.
		Status(http.StatusOK).
		JSON(fiber.Map{
			"code":    http.StatusOK,
			"message": "user logged in successfully",
			"data":    token,
		})
}

func (h *authController) Profile(ctx *fiber.Ctx) error {
	ID, _ := middlewares.GetUserID(ctx)

	user, err := h.authService.GetProfile(ID)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"code":    http.StatusUnauthorized,
			"message": "invalid access token",
			"data":    nil,
		})
	}

	revenue, err := h.orderService.GetAmbassadorsRevenue(ID)
	if err != nil {
		return ctx.
			Status(http.StatusInternalServerError).
			JSON(fiber.Map{
				"code":    http.StatusInternalServerError,
				"message": "internal server error",
				"errors":  err,
			})
	}

	if !strings.Contains(ctx.Path(), "/api/ambassadors") {
		user.Revenue = 0
	}
	user.Revenue = revenue

	return ctx.
		Status(http.StatusOK).
		JSON(fiber.Map{
			"code":    http.StatusOK,
			"message": "profile fetched successfully",
			"data":    user,
		})
}

func (h *authController) Logout(ctx *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:    "jwt",
		Value:   "",
		Expires: time.Now().Add(-time.Hour),
	}
	ctx.Cookie(&cookie)
	return ctx.
		Status(http.StatusOK).
		JSON(fiber.Map{
			"code":    http.StatusOK,
			"message": "user logged out successfully",
			"data":    nil,
		})
}

func (h *authController) UpdateProfile(ctx *fiber.Ctx) error {
	dto := new(dto.UpdateProfileDTO)
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

	ID, _ := middlewares.GetUserID(ctx)
	dto.ID = ID
	updatedUser, err := h.authService.UpdateProfile(dto)
	if err != nil {
		return ctx.
			Status(http.StatusInternalServerError).
			JSON(fiber.Map{
				"code":    http.StatusInternalServerError,
				"message": "internal server errors",
				"error":   err,
			})
	}

	return ctx.
		Status(http.StatusOK).
		JSON(fiber.Map{
			"code":    http.StatusOK,
			"message": "user profile updated successfully",
			"data":    updatedUser,
		})
}

func (h *authController) UpdateProfilePassword(ctx *fiber.Ctx) error {
	dto := new(dto.UpdateProfilePassword)
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

	ID, _ := middlewares.GetUserID(ctx)
	dto.ID = ID
	isUpdated, err := h.authService.UpdateProfilePassword(dto)
	if err != nil {
		return ctx.
			Status(http.StatusInternalServerError).
			JSON(fiber.Map{
				"code":    http.StatusInternalServerError,
				"message": "internal server errors",
				"error":   err.Error(),
			})
	}

	return ctx.
		Status(http.StatusOK).
		JSON(fiber.Map{
			"code":       http.StatusOK,
			"message":    "user profile updated successfully",
			"is_updated": isUpdated,
		})
}

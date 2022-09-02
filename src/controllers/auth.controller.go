package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
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
	cookie := ctx.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte("123"), nil
	})

	if !token.Valid || err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"code":    http.StatusUnauthorized,
			"message": "invalid access token",
			"data":    nil,
		})
	}

	payload := token.Claims.(*jwt.StandardClaims)
	ID, _ := strconv.Atoi(payload.Subject)
	user, err := h.authService.GetProfile(ID)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"code":    http.StatusUnauthorized,
			"message": "invalid access token",
			"data":    nil,
		})
	}

	return ctx.
		Status(http.StatusOK).
		JSON(fiber.Map{
			"code":    http.StatusOK,
			"message": "profile fetched successfully",
			"data":    user,
		})
}

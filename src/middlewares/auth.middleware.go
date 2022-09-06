package middlewares

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

type ClaimsWithScope struct {
	jwt.StandardClaims
	Scope string
}

func IsAuthenticated(ctx *fiber.Ctx) error {
	cookie := ctx.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &ClaimsWithScope{}, func(t *jwt.Token) (interface{}, error) {
		return []byte("123"), nil
	})

	if err != nil || !token.Valid {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"code":    http.StatusUnauthorized,
			"message": "invalid access token",
			"data":    nil,
		})
	}

	payload := token.Claims.(*ClaimsWithScope)
	isAmbassador := strings.Contains(ctx.Path(), "/api/ambassadors")
	if (payload.Scope == "admin" && isAmbassador) || (payload.Scope == "ambassadors" && !isAmbassador) {
		return ctx.Status(http.StatusForbidden).JSON(fiber.Map{
			"code":    http.StatusForbidden,
			"message": "prohibited to aceess this route",
			"errors":  "you are prohibited to acesss this route",
		})
	}

	return ctx.Next()
}

func GetUserID(ctx *fiber.Ctx) (int, error) {
	cookie := ctx.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte("123"), nil
	})
	if err != nil {
		return 0, err
	}

	payload := token.Claims.(*jwt.StandardClaims)
	ID, _ := strconv.Atoi(payload.Subject)
	return ID, nil
}

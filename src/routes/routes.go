package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yogasab/go-monolith-ambassador/src/controllers"
	"github.com/yogasab/go-monolith-ambassador/src/database"
	"github.com/yogasab/go-monolith-ambassador/src/repositories"
	"github.com/yogasab/go-monolith-ambassador/src/services"
)

func Setup(app *fiber.App) {
	api := app.Group("/api")

	authController := controllers.NewAuthController(
		services.NewAuthService(repositories.NewUserRepository(database.DB)), services.NewJWTService())

	admin := api.Group("/admin")
	admin.Post("/auth/register", authController.Register)
	admin.Post("/auth/login", authController.Login)
}

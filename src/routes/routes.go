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

	admin := api.Group("/admin")
	authController := controllers.NewAuthController(
		services.NewAuthService(repositories.NewUserRepository(database.DB)))
	admin.Post("/auth/register", authController.RegisterAdmin)
}

package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yogasab/go-monolith-ambassador/src/controllers"
	"github.com/yogasab/go-monolith-ambassador/src/database"
	"github.com/yogasab/go-monolith-ambassador/src/middlewares"
	"github.com/yogasab/go-monolith-ambassador/src/repositories"
	"github.com/yogasab/go-monolith-ambassador/src/services"
)

func Setup(app *fiber.App) {
	api := app.Group("/api")

	authController := controllers.NewAuthController(
		services.NewAuthService(repositories.NewUserRepository(database.DB)),
		services.NewJWTService(),
	)
	ambassadorController := controllers.NewAmbassadorController(
		services.NewAmbassadorService(repositories.NewUserRepository(database.DB)),
	)
	productController := controllers.NewProductController(
		services.NewProductsService(repositories.NewProductRepository(database.DB)),
	)

	admin := api.Group("/admin")
	admin.Post("auth/register", authController.Register)
	admin.Post("auth/login", authController.Login)

	authenticatedAdmin := admin.Use(middlewares.IsAuthenticated)
	authenticatedAdmin.Get("profile", authController.Profile)
	authenticatedAdmin.Post("logout", authController.Logout)
	authenticatedAdmin.Put("profile/update", authController.UpdateProfile)
	authenticatedAdmin.Put("profile/password", authController.UpdateProfilePassword)
	authenticatedAdmin.Get("ambassadors", ambassadorController.GetAmbassadors)
	authenticatedAdmin.Get("products", productController.GetProducts)
	authenticatedAdmin.Post("products", productController.CreateProduct)
	authenticatedAdmin.Get("products/:id", productController.GetProduct)
	authenticatedAdmin.Put("products/:id", productController.UpdateProduct)
	authenticatedAdmin.Delete("products/:id", productController.DeleteProduct)
}

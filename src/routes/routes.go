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

	authController := controllers.NewAuthController(
		services.NewAuthService(repositories.NewUserRepository(database.DB)),
		services.NewJWTService(),
		services.NewOrderService(repositories.NewOrderRepository(database.DB)),
	)
	ambassadorController := controllers.NewAmbassadorController(
		services.NewAmbassadorService(repositories.NewUserRepository(database.DB)),
	)
	productController := controllers.NewProductController(
		services.NewProductsService(repositories.NewProductRepository(database.DB)),
		services.NewRedisService(repositories.NewRedisRepository(database.RedisClient)),
	)
	linkController := controllers.NewLinkController(
		services.NewLinkService(repositories.NewLinkRepository(database.DB)),
	)
	orderController := controllers.NewrOrderController(
		services.NewOrderService(repositories.NewOrderRepository(database.DB)),
	)

	api := app.Group("/api")

	admin := api.Group("/admin")
	admin.Post("auth/register", authController.Register)
	admin.Post("auth/login", authController.Login)

	authenticatedAdmin := admin.Use(middlewares.IsAuthenticated)
	authenticatedAdmin.Get("profile", authController.Profile)
	authenticatedAdmin.Post("logout", authController.Logout)
	authenticatedAdmin.Put("profile/update", authController.UpdateProfile)
	authenticatedAdmin.Put("profile/password", authController.UpdateProfilePassword)
	authenticatedAdmin.Get("profile/:id/links", linkController.GetUserLinks)
	authenticatedAdmin.Get("ambassadors", ambassadorController.GetAmbassadors)
	authenticatedAdmin.Get("orders", orderController.GetOrders)
	authenticatedAdmin.Get("products", productController.GetProducts)
	authenticatedAdmin.Post("products", productController.CreateProduct)
	authenticatedAdmin.Get("products/:id", productController.GetProduct)
	authenticatedAdmin.Put("products/:id", productController.UpdateProduct)
	authenticatedAdmin.Delete("products/:id", productController.DeleteProduct)

	ambassador := api.Group("/ambassadors")
	ambassador.Post("auth/register", authController.Register)
	ambassador.Post("auth/login", authController.Login)

	authenticatedAmbassador := ambassador.Use(middlewares.IsAuthenticated)
	authenticatedAmbassador.Get("profile", authController.Profile)
	authenticatedAmbassador.Post("logout", authController.Logout)
	authenticatedAmbassador.Put("profile/update", authController.UpdateProfile)
	authenticatedAmbassador.Put("profile/password", authController.UpdateProfilePassword)
	authenticatedAmbassador.Get("products/frontend", productController.GetProductsFrontend)
	authenticatedAmbassador.Get("products/backend", productController.GetProductsBackend)

}

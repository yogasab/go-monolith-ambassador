package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yogasab/go-monolith-ambassador/src/database"
	"github.com/yogasab/go-monolith-ambassador/src/routes"
)

func main() {
	database.Connect()
	database.AutoMigrate()

	app := fiber.New()
	routes.Setup(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Listen(":5000")
}

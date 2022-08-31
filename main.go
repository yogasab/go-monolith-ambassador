package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yogasab/go-monolith-ambassador/src/database"
)

func main() {
	database.Connect()
	database.AutoMigrate()

	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	app.Listen(":3000")
}

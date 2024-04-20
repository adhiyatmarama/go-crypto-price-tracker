package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	database "github.com/adhiyatmarama/go-crypto-price-tracker/database"
)

func main() {
	database.ConnectDatabase()

	app := fiber.New()
	port := "3000"

	app.Use(logger.New(logger.Config{
		TimeZone: "UTC",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("This is crypto price tracker application, written using go")
	})

	app.Listen(":" + port)
}

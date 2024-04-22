package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/adhiyatmarama/go-crypto-price-tracker/coin/coincontroller"
	database "github.com/adhiyatmarama/go-crypto-price-tracker/database"
	"github.com/adhiyatmarama/go-crypto-price-tracker/tracker/trackercontroller"
	usercontroller "github.com/adhiyatmarama/go-crypto-price-tracker/user/usercontroller"
)

func main() {
	database.ConnectDatabase()

	app := fiber.New()
	port := "3000"

	app.Use(logger.New(logger.Config{
		TimeZone: "UTC",
	}))
	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("This is crypto price tracker application, written using go")
	})
	app.Mount("/user", usercontroller.GetRoutes())
	app.Mount("/coin", coincontroller.GetRoutes())
	app.Mount("/tracker", trackercontroller.GetRoutes())

	app.Listen(":" + port)
}

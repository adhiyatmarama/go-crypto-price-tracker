package coincontroller

import (
	"github.com/adhiyatmarama/go-crypto-price-tracker/coin/coinservice"
	"github.com/adhiyatmarama/go-crypto-price-tracker/middlewares"
	"github.com/gofiber/fiber/v2"
)

func GetRoutes() *fiber.App {
	coinRoutes := fiber.New()

	coinRoutes.Use(middlewares.JWTMiddleware)
	coinRoutes.Get("/", GetCoins)
	coinRoutes.Get("/:id", GetCoinById)

	return coinRoutes
}

func GetCoins(c *fiber.Ctx) error {
	query := c.Queries()

	coins, err := coinservice.GetCoins(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": coins,
	})
}

func GetCoinById(c *fiber.Ctx) error {
	id := c.Params("id")

	coin, err := coinservice.GetCoinById(id)
	if err != nil {
		switch err.Error() {
		case "not found":
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Coin Not Found",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal server error",
				"error":   err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": coin,
	})
}

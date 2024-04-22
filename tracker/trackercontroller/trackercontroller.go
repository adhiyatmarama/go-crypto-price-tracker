package trackercontroller

import (
	"fmt"
	"log"

	"github.com/adhiyatmarama/go-crypto-price-tracker/database"
	"github.com/adhiyatmarama/go-crypto-price-tracker/middlewares"
	"github.com/adhiyatmarama/go-crypto-price-tracker/tracker/trackermodel"
	"github.com/adhiyatmarama/go-crypto-price-tracker/tracker/trackerservice"
	"github.com/gofiber/fiber/v2"
)

func GetRoutes() *fiber.App {
	trackerRoutes := fiber.New()

	trackerRoutes.Use(middlewares.JWTMiddleware)
	trackerRoutes.Get("/", GetTrackersByUser)
	trackerRoutes.Post("/", CreateTracker)
	trackerRoutes.Delete("/", RemoveTracker)

	return trackerRoutes
}

func CreateTracker(c *fiber.Ctx) error {
	userEmail := fmt.Sprintf("%v", c.Locals("user_email"))

	var tracker trackermodel.Tracker
	if err := c.BodyParser(&tracker); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "bad request",
			"error":   err.Error(),
		})
	}

	if tracker.CoinId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Coin id must not be empty",
		})
	}

	created, err := trackerservice.CreateTracker(userEmail, tracker)
	if err != nil {
		switch err.Error() {
		case "coin has been added to your tracker":
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": err.Error(),
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error when create tracker to DB",
				"error":   err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully create a tracker",
		"tracker": fiber.Map{
			"user_email": created.UserEmail,
			"coin_id":    created.CoinId,
		},
	})
}

func GetTrackerByUserEmailAndCoinId(userEmail string, coinId string) ([]trackermodel.Tracker, error) {
	var result []trackermodel.Tracker
	rows, err := database.DB.Query(
		fmt.Sprintf(
			"select user_email, coin_id from Trackers where user_email = '%s' and coin_id = '%s'",
			userEmail,
			coinId,
		),
	)
	if rows == nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user_email string
		var coin_id string
		err = rows.Scan(&user_email, &coin_id)
		if err != nil {
			return nil, err
		}
		result = append(result, trackermodel.Tracker{
			UserEmail: user_email,
			CoinId:    coin_id,
		})
	}

	return result, nil

}

func GetTrackersByUser(c *fiber.Ctx) error {
	userEmail := fmt.Sprintf("%v", c.Locals("user_email"))

	result, err := trackerservice.GetTrackersByUser(userEmail)
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error when get user tracker",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": result,
	})
}

func RemoveTracker(c *fiber.Ctx) error {
	userEmail := fmt.Sprintf("%v", c.Locals("user_email"))

	var tracker trackermodel.Tracker
	if err := c.BodyParser(&tracker); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "bad request",
			"error":   err.Error(),
		})
	}

	if tracker.CoinId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Coin id must not be empty",
		})
	}

	result, err := database.DB.Exec(
		fmt.Sprintf(
			"delete from Trackers where user_email = '%s' and coin_id = '%s'",
			userEmail,
			tracker.CoinId,
		),
	)
	if err != nil {
		log.Print(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Tracker not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully remove coin from tracker",
	})
}

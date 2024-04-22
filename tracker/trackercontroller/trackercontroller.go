package trackercontroller

import (
	"fmt"
	"sort"
	"sync"

	"github.com/adhiyatmarama/go-crypto-price-tracker/database"
	"github.com/adhiyatmarama/go-crypto-price-tracker/libs/libscoincap"
	"github.com/adhiyatmarama/go-crypto-price-tracker/middlewares"
	"github.com/adhiyatmarama/go-crypto-price-tracker/tracker/trackermodel"
	"github.com/gofiber/fiber/v2"
)

func GetRoutes() *fiber.App {
	trackerRoutes := fiber.New()

	trackerRoutes.Use(middlewares.JWTMiddleware)
	trackerRoutes.Get("/", GetTrackersByUser)
	trackerRoutes.Post("/", CreateTracker)

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

	existing, err := GetTrackerByUserEmailAndCoinId(userEmail, tracker.CoinId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Server Error",
			"error":   err.Error(),
		})
	}
	if len(existing) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Coin has been added to your tracker",
		})
	}

	// Add tracker to table
	_, err = database.DB.Exec(fmt.Sprintf("INSERT INTO Trackers(user_email, coin_id) VALUES('%s', '%s' )", userEmail, tracker.CoinId))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error when create tracker to DB",
			"error":   err.Error(),
		})
	}

	created, err := GetTrackerByUserEmailAndCoinId(userEmail, tracker.CoinId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error when create tracker to DB",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully create a tracker",
		"tracker": fiber.Map{
			"user_email": created[0].UserEmail,
			"coin_id":    created[0].CoinId,
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
	result := []*libscoincap.Coin{}

	// Get coin Ids tracked by user
	rows, err := database.DB.Query(
		fmt.Sprintf(
			"select coin_id from Trackers where user_email = '%s'",
			userEmail,
		),
	)
	if rows == nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"data": result,
		})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error when get tracker to DB",
			"error":   err.Error(),
		})
	}
	defer rows.Close()
	coinIds := []string{}
	for rows.Next() {
		var coin_id string
		err = rows.Scan(&coin_id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error when get tracker to DB",
				"error":   err.Error(),
			})
		}
		coinIds = append(coinIds, coin_id)
	}

	// Get detail of each coin by coin id concurrently using go routines,  wait group, and channels
	wg := sync.WaitGroup{}
	getCoinErrorCh := make(chan error, len(coinIds))
	coinCh := make(chan *libscoincap.Coin, len(coinIds))
	for _, coinId := range coinIds {
		wg.Add(1)
		go func() {
			defer wg.Done()
			coin, err := libscoincap.GetAssetById(coinId)
			if err != nil {
				fmt.Println(err, " === ada error")
				getCoinErrorCh <- err
				return
			}
			coinCh <- coin
		}()
	}
	wg.Wait()
	close(getCoinErrorCh)
	close(coinCh)

	// return error if got error
	if err := <-getCoinErrorCh; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error when get coin data",
		})
	}

	// get coin detail result and sort by coin's rank
	for coin := range coinCh {
		result = append(result, coin)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[j].Rank > result[i].Rank
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": result,
	})
}

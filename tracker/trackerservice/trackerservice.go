package trackerservice

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/adhiyatmarama/go-crypto-price-tracker/database"
	"github.com/adhiyatmarama/go-crypto-price-tracker/libs/libscoincap"
	"github.com/adhiyatmarama/go-crypto-price-tracker/libs/libscurrencyapi"
	"github.com/adhiyatmarama/go-crypto-price-tracker/tracker/trackermodel"
)

func CreateTracker(userEmail string, tracker trackermodel.Tracker) (*trackermodel.Tracker, error) {
	existing, err := GetTrackerByUserEmailAndCoinId(userEmail, tracker.CoinId)
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}

	if len(existing) > 0 {
		return nil, errors.New("coin has been added to your tracker")
	}

	// Add tracker to table
	_, err = database.DB.Exec(fmt.Sprintf("INSERT INTO Trackers(user_email, coin_id) VALUES('%s', '%s' )", userEmail, tracker.CoinId))
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}

	created, err := GetTrackerByUserEmailAndCoinId(userEmail, tracker.CoinId)
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}

	return &created[0], nil
}

func GetTrackersByUser(userEmail string) ([]libscoincap.Coin, error) {
	// Get coin Ids tracked by user
	rows, err := database.DB.Query(
		fmt.Sprintf(
			"select coin_id from Trackers where user_email = '%s'",
			userEmail,
		),
	)
	if rows == nil {
		return []libscoincap.Coin{}, nil
	}
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}
	defer rows.Close()
	coinIds := []string{}
	for rows.Next() {
		var coin_id string
		err = rows.Scan(&coin_id)
		if err != nil {
			log.Print(err.Error())
			return nil, err
		}
		coinIds = append(coinIds, coin_id)
	}

	// Get IDR value to USD
	idrVal, err := libscurrencyapi.GetLatestExchangeRate("USD", "IDR")
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}

	// Get detail of each coin by coin id concurrently using go routines,  wait group, and channels
	wg := sync.WaitGroup{}
	result := []libscoincap.Coin{}
	getCoinErrorCh := make(chan error, len(coinIds))
	coinCh := make(chan *libscoincap.Coin, len(coinIds))
	for _, coinId := range coinIds {
		wg.Add(1)
		go func(coinId string, idrVal interface{}) {
			defer wg.Done()
			coin, err := libscoincap.GetAssetById(coinId)
			if err != nil {
				getCoinErrorCh <- err
				return
			}
			usdVal, err := strconv.ParseFloat(coin.PriceUsd, 64)
			if err != nil {
				getCoinErrorCh <- err
				return
			}
			coin.PriceIdr = fmt.Sprintf("%.4f", usdVal*idrVal.(float64))
			coinCh <- coin
		}(coinId, idrVal)
	}
	wg.Wait()
	close(getCoinErrorCh)
	close(coinCh)

	// return error if got error
	if err := <-getCoinErrorCh; err != nil {
		log.Print(err.Error())
		return nil, err
	}

	// get coin detail result and sort by coin's rank
	for coin := range coinCh {
		result = append(result, *coin)
	}

	return result, nil
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
		log.Print(err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user_email string
		var coin_id string
		err = rows.Scan(&user_email, &coin_id)
		if err != nil {
			log.Print(err.Error())
			return nil, err
		}
		result = append(result, trackermodel.Tracker{
			UserEmail: user_email,
			CoinId:    coin_id,
		})
	}

	return result, nil

}

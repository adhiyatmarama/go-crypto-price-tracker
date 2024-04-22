package coinservice

import (
	"fmt"
	"strconv"

	"github.com/adhiyatmarama/go-crypto-price-tracker/libs/libscoincap"
	"github.com/adhiyatmarama/go-crypto-price-tracker/libs/libscurrencyapi"
)

func GetCoins(query map[string]string) ([]libscoincap.Coin, error) {
	coins, err := libscoincap.GetAssets(query)
	if err != nil {
		return nil, err
	}

	// Get IDR value to USD
	idrVal, err := libscurrencyapi.GetLatestExchangeRate("USD", "IDR")
	if err != nil {
		return nil, err
	}

	// Add IDR price
	for i := 0; i < len(coins); i++ {
		usdVal, _ := strconv.ParseFloat(coins[i].PriceUsd, 64)
		coins[i].PriceIdr = fmt.Sprintf("%.4f", usdVal*idrVal.(float64))
	}
	return coins, nil
}

func GetCoinById(coinId string) (*libscoincap.Coin, error) {
	coin, err := libscoincap.GetAssetById(coinId)
	if err != nil {
		return nil, err
	}

	// Get IDR value to USD
	idrVal, err := libscurrencyapi.GetLatestExchangeRate("USD", "IDR")
	if err != nil {
		return nil, err
	}

	// Add IDR Price
	usdVal, _ := strconv.ParseFloat(coin.PriceUsd, 64)
	coin.PriceIdr = fmt.Sprintf("%.4f", usdVal*idrVal.(float64))

	return coin, err
}

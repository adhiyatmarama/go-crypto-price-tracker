package coinservice

import (
	"github.com/adhiyatmarama/go-crypto-price-tracker/libs/libscoincap"
)

func GetCoins(query map[string]string) ([]libscoincap.Coin, error) {
	coins, err := libscoincap.GetAssets(query)
	if err != nil {
		return nil, err
	}

	return coins, nil
}

func GetCoinById(coinId string) (*libscoincap.Coin, error) {
	coin, err := libscoincap.GetAssetById(coinId)
	if err != nil {
		return nil, err
	}

	return coin, err
}

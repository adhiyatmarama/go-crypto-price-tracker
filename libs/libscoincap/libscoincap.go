package libscoincap

import (
	"encoding/json"
	"log"
	"net/http"
)

var BASE_URL = "https://api.coincap.io/v2"

type Coin struct {
	Id                string `json:"id"`
	Rank              string `json:"rank"`
	Symbol            string `json:"symbol"`
	Name              string `json:"name"`
	Supply            string `json:"supply"`
	MaxSupply         string `json:"maxSupply"`
	MarketCapUsd      string `json:"marketCapUsd"`
	VolumeUsd24Hr     string `json:"volumeUsd24Hr"`
	PriceUsd          string `json:"priceUsd"`
	ChangePercent24Hr string `json:"changePercent24Hr"`
	Vwap24Hr          string `json:"vwap24Hr"`
	Explorer          string `json:"explorer"`
}

type GetAssetByIdResponse struct {
	Data      Coin `json:"data"`
	Timestamp int  `json:"timestamp"`
}

func GetAssetById(coinId string) (*Coin, error) {
	var getResponse GetAssetByIdResponse

	response, err := http.Get(BASE_URL + "/assets/" + coinId)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&getResponse); err != nil {
		log.Print(err)
		return nil, err
	}

	return &getResponse.Data, nil
}

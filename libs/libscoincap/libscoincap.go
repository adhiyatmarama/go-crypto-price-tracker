package libscoincap

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
)

var BASE_URL = "https://api.coincap.io/v2"

type Coin struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	PriceUsd string `json:"priceUsd"`
	PriceIdr string `json:"priceIdr"`
}

type GetAssetByIdResponse struct {
	Data      Coin `json:"data"`
	Timestamp int  `json:"timestamp"`
}

type GetAssetsResponse struct {
	Data      []Coin `json:"data"`
	Timestamp int    `json:"timestamp"`
}

func GetAssetById(coinId string) (*Coin, error) {
	var getResponse GetAssetByIdResponse

	response, err := http.Get(BASE_URL + "/assets/" + coinId)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer response.Body.Close()

	resStatusCode := response.StatusCode
	if resStatusCode == 404 {
		return nil, errors.New("not found")
	}

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&getResponse); err != nil {
		log.Print(err)
		return nil, err
	}

	return &getResponse.Data, nil
}

func GetAssets(query map[string]string) ([]Coin, error) {
	var getResponse GetAssetsResponse

	baseUrl, _ := url.Parse(BASE_URL + "/assets")
	params := url.Values{}
	params.Add("limit", query["limit"])
	params.Add("offset", query["offset"])
	baseUrl.RawQuery = params.Encode()

	response, err := http.Get(baseUrl.String())
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

	return getResponse.Data, nil
}

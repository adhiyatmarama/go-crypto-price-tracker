package libscurrencyapi

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

type GetLatestExchangeRateResponse struct {
	Data map[string]map[string]interface{}
}

var BASE_URL = "https://api.currencyapi.com/v3"
var apiKey = "cur_live_SOBRQcke8yWWlvwrrxgETMCJgWuQQHBODbQLnwiU"

func GetLatestExchangeRate(baseCurrency string, currencies string) (interface{}, error) {
	var getResponse GetLatestExchangeRateResponse

	baseUrl, _ := url.Parse(BASE_URL + "/latest")
	params := url.Values{}
	params.Add("base_currency", baseCurrency)
	params.Add("currencies", currencies)
	baseUrl.RawQuery = params.Encode()

	client := &http.Client{}
	req, err := http.NewRequest("GET", baseUrl.String(), nil)
	if err != nil {
		log.Print("Error creating HTTP request:", err)
		return 0, err
	}
	req.Header.Add("apiKey", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		log.Print("Error sending HTTP request:", err)
		return 0, err
	}

	// Read the response body
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&getResponse); err != nil {
		log.Print(err)
		return nil, err
	}

	return getResponse.Data["IDR"]["value"], nil
}

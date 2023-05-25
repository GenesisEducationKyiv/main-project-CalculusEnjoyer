package providsers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const CoinGeckoUrl = "https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=%s"

type CoinGeckoRateProvider struct{}

type CoinGeckoResponse map[string]map[string]float64

func (*CoinGeckoRateProvider) GetExchangeRate(baseCurrency, targetCurrency string) (rate float64, err error) {
	url := fmt.Sprintf(CoinGeckoUrl, baseCurrency, targetCurrency)

	response, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to make request: %v", err)
	}
	defer response.Body.Close()

	var data CoinGeckoResponse
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return 0, fmt.Errorf("failed to decode response: %v", err)
	}

	rates := data[strings.ToLower(baseCurrency)]
	rate = rates[strings.ToLower(targetCurrency)]

	return rate, nil
}

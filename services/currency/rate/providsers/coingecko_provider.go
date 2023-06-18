package providsers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type CoinGeckoRateProvider struct{ coinGeckoURL string }

type CoinGeckoResponse map[string]map[string]float64

func NewCoinGeckoRateProvider() *CoinGeckoRateProvider {
	err := godotenv.Load()
	if err != err {
		log.Fatalf("Failed to load configs")
	}

	return &CoinGeckoRateProvider{
		os.Getenv("COINGEKCO_URL"),
	}
}

func (p *CoinGeckoRateProvider) GetExchangeRate(baseCurrency, targetCurrency string) (rate float64, err error) {
	url := fmt.Sprintf(p.coinGeckoURL, baseCurrency, targetCurrency)

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

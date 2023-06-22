package crypto

import (
	"currency/cerror"
	"currency/config"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type CoinGeckoRateProvider struct{ coinGeckoURL string }

type CoinGeckoResponse map[string]map[string]float64

func NewCoinGeckoRateProvider(conf config.Config) *CoinGeckoRateProvider {
	return &CoinGeckoRateProvider{
		conf.CoinGekcoURL,
	}
}

func (p *CoinGeckoRateProvider) GetExchangeRate(baseCurrency, targetCurrency string) (rate float64, errBodyClose error) {
	url := fmt.Sprintf(p.coinGeckoURL, baseCurrency, targetCurrency)

	response, err := http.Get(url)
	if err != nil || response.StatusCode != http.StatusOK {
		return cerror.ErrRateValue, cerror.ErrRate
	}
	defer response.Body.Close()

	var data CoinGeckoResponse
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return cerror.ErrRateValue, cerror.ErrDecode
	}

	rates := data[strings.ToLower(baseCurrency)]
	rate = rates[strings.ToLower(targetCurrency)]

	return rate, nil
}

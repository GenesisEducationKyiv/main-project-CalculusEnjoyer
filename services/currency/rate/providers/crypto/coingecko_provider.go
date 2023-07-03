package crypto

import (
	"currency/cerror"
	"currency/config"
	"currency/rate/messages"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type CoinGeckoRateProvider struct {
	coinGeckoURL        string
	supportedCurrencies map[messages.Currency]string
}

type CoinGeckoResponse map[string]map[string]float64

func NewCoinGeckoRateProvider(conf config.Config) *CoinGeckoRateProvider {
	currencies := map[messages.Currency]string{
		messages.BTC: "bitcoin",
		messages.UAH: "uah",
	}

	return &CoinGeckoRateProvider{
		conf.CoinGekcoURL,
		currencies,
	}
}

func (p *CoinGeckoRateProvider) GetExchangeRate(baseCurrency, targetCurrency messages.Currency) (float64, error) {
	convertedBase, err := p.currencyToString(baseCurrency)
	if err != nil {
		return cerror.ErrRateValue, err
	}

	convertedTarget, err := p.currencyToString(targetCurrency)
	if err != nil {
		return cerror.ErrRateValue, err
	}
	url := fmt.Sprintf(p.coinGeckoURL, convertedBase, convertedTarget)

	response, err := http.Get(url)
	if err != nil || response.StatusCode != http.StatusOK {
		return cerror.ErrRateValue, cerror.ErrRate
	}
	defer response.Body.Close()

	return decodeRateResponse(response, convertedBase, convertedTarget)
}

func decodeRateResponse(resp *http.Response, baseCurrencyName, targetCurrencyName string) (float64, error) {
	var data CoinGeckoResponse
	err := json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return cerror.ErrRateValue, err
	}

	rates := data[baseCurrencyName]
	rate, ok := rates[targetCurrencyName]
	if !ok {
		return cerror.ErrRateValue, cerror.ErrRate
	}

	log.Printf("Getting rate from COINGECKO: %f", rate)
	return rate, nil
}

func (p *CoinGeckoRateProvider) currencyToString(currency messages.Currency) (string, error) {
	result := p.supportedCurrencies[currency]
	if result == "" {
		return result, fmt.Errorf("%s is unsupported currency", string(currency))
	}
	return result, nil
}

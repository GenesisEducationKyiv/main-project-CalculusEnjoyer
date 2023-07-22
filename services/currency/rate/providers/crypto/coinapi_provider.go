package crypto

import (
	"currency/cerror"
	"currency/config"
	"currency/domain"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

type CoinAPIProvider struct {
	endpoint   string
	apiKey     string
	currencies map[domain.Currency]string
}

func NewCoinAPIProvider(conf config.Config) *CoinAPIProvider {
	currencies := map[domain.Currency]string{
		domain.BTC: "BTC",
		domain.UAH: "UAH",
	}

	return &CoinAPIProvider{
		endpoint:   conf.CoinApiURL,
		apiKey:     conf.CoinApiKey,
		currencies: currencies,
	}
}

func (p *CoinAPIProvider) GetExchangeRate(baseCurrency, targetCurrency domain.Currency) (float64, error) {
	request, err := p.generateHTTPRequest(baseCurrency, targetCurrency)
	if err != nil {
		return cerror.ErrRateValue, errors.Wrap(err, "can not generate http request for getting rate")
	}

	res, err := http.DefaultClient.Do(request)
	defer res.Body.Close()

	if err != nil || res.StatusCode != http.StatusOK {
		return cerror.ErrRateValue, errors.Wrap(cerror.ErrRate, "can not make request to the COIN API")
	}

	return p.extractRate(res)
}

func (p *CoinAPIProvider) generateHTTPRequest(baseCurrency, targetCurrency domain.Currency) (*http.Request, error) {
	endpoint, err := p.generateEndpoint(baseCurrency, targetCurrency)
	if err != nil {
		return nil, errors.Wrap(err, "can not generate request")
	}

	req, err := http.NewRequest(
		http.MethodGet,
		endpoint,
		nil,
	)

	if err != nil {
		return nil, errors.Wrap(cerror.ErrRate, "can not generate http request")
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-CoinAPI-Key", p.apiKey)

	return req, nil
}

func (p *CoinAPIProvider) Name() string {
	return "COINAPI"
}

func (p *CoinAPIProvider) generateEndpoint(baseCurrency, targetCurrency domain.Currency) (string, error) {
	convertedBase, err := p.currencyToString(baseCurrency)
	if err != nil {
		return "", err
	}

	convertedTarget, err := p.currencyToString(targetCurrency)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(p.endpoint, convertedBase, convertedTarget), nil
}

func (p *CoinAPIProvider) extractRate(response *http.Response) (float64, error) {
	type rateResponse struct {
		Rate float64 `json:"rate"`
	}
	var data rateResponse

	err := json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return cerror.ErrRateValue, err
	}

	log.Printf("Getting rate from COINAPI: %f", data.Rate)
	return data.Rate, nil
}

func (p *CoinAPIProvider) currencyToString(currency domain.Currency) (string, error) {
	result := p.currencies[currency]
	if result == "" {
		return result, fmt.Errorf("%s is unsupported currency", string(currency))
	}
	return result, nil
}

package crypto

import (
	"currency/cerror"
	"currency/config"
	"currency/domain"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

type KunaRateProvider struct {
	kunaURL             string
	supportedCurrencies map[domain.Currency]string
}

type KunaRateProviderResponse map[string][]map[string]interface{}

func NewKunaRateProvider(conf config.Config) *KunaRateProvider {
	currencies := map[domain.Currency]string{
		domain.BTC: "BTC",
		domain.UAH: "UAH",
		domain.ETH: "ETH",
		domain.USD: "USD",
		domain.XMR: "XMR",
	}

	return &KunaRateProvider{
		conf.KunaURL,
		currencies,
	}
}

func (p *KunaRateProvider) GetExchangeRate(baseCurrency, targetCurrency domain.Currency) (float64, error) {
	request, err := p.generateHTTPRequest(baseCurrency, targetCurrency)
	if err != nil {
		return cerror.ErrRateValue, errors.Wrap(err, "can not generate request to Kuna.io")
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil || response.StatusCode != http.StatusOK {
		return cerror.ErrRateValue, cerror.ErrRate
	}
	defer response.Body.Close()

	return p.extractRate(response)
}

func (p *KunaRateProvider) Name() string {
	return "KUNA"
}

func (p *KunaRateProvider) extractRate(response *http.Response) (float64, error) {
	var data KunaRateProviderResponse
	err := json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return cerror.ErrRateValue, cerror.ErrDecode
	}

	pair := data["data"]
	if len(pair) == 0 {
		return cerror.ErrRateValue, cerror.ErrDecode
	}

	rate, ok := pair[0]["price"].(string)
	if !ok {
		return cerror.ErrRateValue, cerror.ErrDecode
	}

	float, err := strconv.ParseFloat(rate, 64)
	if err != nil {
		return cerror.ErrRateValue, errors.Wrap(err, "failed to decode Kuna response to float")
	}

	return float, nil
}

func (p *KunaRateProvider) generateHTTPRequest(baseCurrency, targetCurrency domain.Currency) (*http.Request, error) {
	convertedBase, err := p.currencyToString(baseCurrency)
	if err != nil {
		return nil, err
	}

	convertedTarget, err := p.currencyToString(targetCurrency)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(p.kunaURL, convertedBase, convertedTarget)

	if err != nil {
		return nil, errors.Wrap(err, "can not generate request")
	}

	req, err := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(cerror.ErrRate, "can not generate http request")
	}

	req.Header.Add("Accept", "application/json")

	return req, nil
}

func (p *KunaRateProvider) currencyToString(currency domain.Currency) (string, error) {
	result := p.supportedCurrencies[currency]
	if result == "" {
		return result, fmt.Errorf("%s is unsupported currency", string(currency))
	}
	return result, nil
}

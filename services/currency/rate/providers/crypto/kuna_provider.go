package crypto

import (
	"currency/cerror"
	"currency/config"
	"currency/rate/messages"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

type KunaRateProvider struct {
	kunaURL             string
	supportedCurrencies map[messages.Currency]string
}

type KunaRateProviderResponse map[string][]map[string]interface{}

func NewKunaRateProvider(conf config.Config) *KunaRateProvider {
	currencies := map[messages.Currency]string{
		messages.BTC: "BTC",
		messages.UAH: "UAH",
	}

	return &KunaRateProvider{
		conf.KunaURL,
		currencies,
	}
}

func (p *KunaRateProvider) GetExchangeRate(baseCurrency, targetCurrency messages.Currency) (float64, error) {
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
	rate := pair[0]["price"].(string)

	float, err := strconv.ParseFloat(rate, 64)
	if err != nil {
		return cerror.ErrRateValue, errors.Wrap(err, "failed to decode Kuna response to float")
	}

	return float, nil
}

func (p *KunaRateProvider) currencyToString(currency messages.Currency) (string, error) {
	result := p.supportedCurrencies[currency]
	if result == "" {
		return result, fmt.Errorf("%s is unsupported currency", string(currency))
	}
	return result, nil
}

func (p *KunaRateProvider) generateHTTPRequest(baseCurrency, targetCurrency messages.Currency) (*http.Request, error) {
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

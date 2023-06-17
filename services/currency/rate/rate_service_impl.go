package rate

import (
	"currency/rate/messages"
	"currency/rate/providsers"
	"time"
)

type rateService struct {
	Provider providsers.RateProvider
}

func NewService(provider providsers.RateProvider) RateService { return &rateService{provider} }

func (r *rateService) GetRate(currencies messages.RateRequest) (rate messages.RateResult, err error) {
	btcRate, err := r.Provider.GetExchangeRate(currencies.BaseCurrency, currencies.TargetCurrency)

	return messages.RateResult{Rate: btcRate, Timestamp: time.Now().String()}, err
}

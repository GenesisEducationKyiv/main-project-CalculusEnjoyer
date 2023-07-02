package rate

import (
	"currency/rate/messages"
	"time"
)

type TimeProvider interface {
	Now() time.Time
}

type RateProvider interface {
	GetExchangeRate(baseCurrency, targetCurrency messages.Currency) (float64, error)
}

type RateService interface {
	GetRate(currencies messages.RateRequest) (messages.RateResult, error)
}

type rateService struct {
	timeProvider TimeProvider
	rateProvider RateProvider
}

func NewRateService(rateProvider RateProvider, timeProvider TimeProvider) RateService {
	return &rateService{timeProvider, rateProvider}
}

func (r *rateService) GetRate(currencies messages.RateRequest) (rate messages.RateResult, err error) {
	btcRate, err := r.rateProvider.GetExchangeRate(currencies.BaseCurrency, currencies.TargetCurrency)

	return messages.RateResult{Rate: btcRate, Timestamp: r.timeProvider.Now()}, err
}

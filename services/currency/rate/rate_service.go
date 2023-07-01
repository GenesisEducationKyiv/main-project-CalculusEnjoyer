package rate

import (
	"currency/rate/messages"
	"time"
)

type TimeProvider interface {
	Now() time.Time
}

type RateProvider interface {
	GetExchangeRate(baseCurrency, targetCurrency string) (rate float64, err error)
}

type BaseRateService interface {
	GetRate(currencies messages.RateRequest) (rate messages.RateResult, err error)
}

type rateService struct {
	timeProvider TimeProvider
	rateProvider RateProvider
}

func NewRateService(rateProvider RateProvider, timeProvider TimeProvider) BaseRateService {
	return &rateService{timeProvider, rateProvider}
}

func (r *rateService) GetRate(currencies messages.RateRequest) (rate messages.RateResult, err error) {
	btcRate, err := r.rateProvider.GetExchangeRate(currencies.BaseCurrency, currencies.TargetCurrency)

	return messages.RateResult{Rate: btcRate, Timestamp: r.timeProvider.Now()}, err
}

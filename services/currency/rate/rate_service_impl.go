package rate

import (
	"currency/rate/messages"
	"currency/rate/providsers/crypto"
	mytime "currency/rate/providsers/time"
)

type rateService struct {
	timeProvider mytime.TimeProvider
	rateProvider crypto.RateProvider
}

func NewService(rateProvider crypto.RateProvider, timeProvider mytime.TimeProvider) RateService {
	return &rateService{timeProvider, rateProvider}
}

func (r *rateService) GetRate(currencies messages.RateRequest) (rate messages.RateResult, err error) {
	btcRate, err := r.rateProvider.GetExchangeRate(currencies.BaseCurrency, currencies.TargetCurrency)

	return messages.RateResult{Rate: btcRate, Timestamp: r.timeProvider.Now()}, err
}

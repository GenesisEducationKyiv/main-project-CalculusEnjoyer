package rate

import (
	"currency/domain"
	"time"
)

type TimeProvider interface {
	Now() time.Time
}

type RateProvider interface {
	GetExchangeRate(baseCurrency, targetCurrency domain.Currency) (float64, error)
	Name() string
}

type RateService interface {
	GetRate(currencies domain.RateRequest) (domain.RateResult, error)
}

type rateService struct {
	timeProvider TimeProvider
	rateProvider RateProvider
}

func NewRateService(rateProvider RateProvider, timeProvider TimeProvider) RateService {
	return &rateService{timeProvider, rateProvider}
}

func (r *rateService) GetRate(currencies domain.RateRequest) (rate domain.RateResult, err error) {
	btcRate, err := r.rateProvider.GetExchangeRate(currencies.BaseCurrency, currencies.TargetCurrency)

	return domain.RateResult{Rate: btcRate, Timestamp: r.timeProvider.Now()}, err
}

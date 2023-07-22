package rate

import (
	"currency/cerror"
	"currency/config"
	"currency/domain"
	"fmt"
	"time"
)

type CachedProvider struct {
	provider     RateProvider
	conf         config.Config
	timeProvider TimeProvider
	lastRateTime time.Time
	rate         float64
}

func NewCachedProvider(provider RateProvider, conf config.Config, timeProvider TimeProvider) *CachedProvider {
	return &CachedProvider{
		provider:     provider,
		conf:         conf,
		timeProvider: timeProvider,
		lastRateTime: time.Unix(0, 0),
	}
}

func (r *CachedProvider) GetExchangeRate(baseCurrency, targetCurrency domain.Currency) (float64, error) {
	if r.timeProvider.Now().Sub(r.lastRateTime) <= r.conf.CacheValidTime {
		return r.rate, nil
	}

	rate, err := r.provider.GetExchangeRate(baseCurrency, targetCurrency)
	if err != nil {
		return cerror.ErrRateValue, err
	}

	r.rate = rate
	r.lastRateTime = r.timeProvider.Now()
	return rate, nil
}

func (r *CachedProvider) Name() string {
	return fmt.Sprintf("Cached provider powered by: %s", r.provider.Name())
}

package rate

import (
	"currency/cerror"
	"currency/config"
	"currency/domain"
	"currency/logger"
	"fmt"
	"github.com/patrickmn/go-cache"
)

type CachedProvider struct {
	provider RateProvider
	cache    *cache.Cache
}

func NewCachedProvider(provider RateProvider, conf config.Config) *CachedProvider {
	return &CachedProvider{
		provider: provider,
		cache:    cache.New(conf.CacheValidTime, conf.CacheValidTime),
	}
}

func (r *CachedProvider) GetExchangeRate(baseCurrency, targetCurrency domain.Currency) (float64, error) {
	cacheRate, found := r.cache.Get(makeKey(baseCurrency, targetCurrency))
	if found {
		logger.DefaultLog(logger.INFO, "getting rate from cache")
		return cacheRate.(float64), nil
	}

	rate, err := r.provider.GetExchangeRate(baseCurrency, targetCurrency)
	if err != nil {
		return cerror.ErrRateValue, err
	}

	if err = r.cache.Add(makeKey(baseCurrency, targetCurrency), rate, cache.DefaultExpiration); err != nil {
		logger.DefaultLog(logger.ERROR, "can not add rate to the cache")
	}

	return rate, nil
}

func (r *CachedProvider) Name() string {
	return fmt.Sprintf("Cached provider powered by: %s", r.provider.Name())
}

func makeKey(baseCurrency, targetCurrency domain.Currency) string {
	return string(baseCurrency) + string(targetCurrency)
}

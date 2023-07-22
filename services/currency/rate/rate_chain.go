package rate

import (
	"currency/domain"
)

type RateLink struct {
	next     *RateLink
	provider RateProvider
}

func NewRateLink(provider RateProvider) *RateLink {
	return &RateLink{
		provider: provider,
	}
}

func (r *RateLink) SetNextLink(next *RateLink) {
	r.next = next
}

func (r *RateLink) GetExchangeRate(baseCurrency, targetCurrency domain.Currency) (float64, error) {
	rate, err := r.provider.GetExchangeRate(baseCurrency, targetCurrency)
	if err != nil && r.next != nil {
		return r.next.GetExchangeRate(baseCurrency, targetCurrency)
	}
	return rate, err
}

func (r *RateLink) Name() string {
	return r.provider.Name()
}

package providsers

type RateProvider interface {
	GetExchangeRate(baseCurrency, targetCurrency string) (rate float64, err error)
}

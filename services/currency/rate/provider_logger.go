package rate

import (
	"currency/domain"
	"currency/logger"
	"fmt"
)

type RateLogger struct {
	provider RateProvider
	logger   logger.Logger
}

func NewRateLogger(provider RateProvider, logger logger.Logger) *RateLogger {
	return &RateLogger{
		provider: provider,
		logger:   logger,
	}
}

func (p *RateLogger) GetExchangeRate(baseCurrency, targetCurrency domain.Currency) (float64, error) {
	rate, err := p.provider.GetExchangeRate(baseCurrency, targetCurrency)
	if err != nil {
		p.logger.Log(logger.ERROR, fmt.Sprintf("Getting rate from %s FAILED: %v", p.provider.Name(), err))
	} else {
		p.logger.Log(logger.INFO, fmt.Sprintf("Getting rate from %s: %f", p.provider.Name(), rate))
	}

	return rate, err
}

func (p *RateLogger) Name() string {
	return p.provider.Name()
}

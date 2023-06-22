package ctrl

import (
	"api/models"
)

type CurrencyProvider interface {
	GetRate(request *models.RateRequest) (*models.RateResponse, error)
}

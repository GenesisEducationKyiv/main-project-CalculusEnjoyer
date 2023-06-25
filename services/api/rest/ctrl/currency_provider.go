package ctrl

import (
	"api/models"
	"context"
)

type CurrencyProvider interface {
	GetRate(request models.RateRequest, cnx context.Context) (*models.RateResponse, error)
}

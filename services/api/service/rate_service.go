package service

import (
	"api/domain"
	"context"
)

type RateProvider interface {
	GetRate(request domain.RateRequest, cnx context.Context) (*domain.RateResponse, error)
}

type RateService struct {
	rateProvider RateProvider
}

func NewRateServiece(provider RateProvider) *RateService {
	return &RateService{rateProvider: provider}
}

func (s *RateService) GetRate(request domain.RateRequest, cnx context.Context) (*domain.RateResponse, error) {
	return s.rateProvider.GetRate(request, cnx)
}

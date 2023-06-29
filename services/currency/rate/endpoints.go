package rate

import (
	"context"
	"currency/rate/messages"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	GetRateEndpoint endpoint.Endpoint
}

func NewEndpointSet(svc RateService) Endpoints {
	return Endpoints{
		GetRateEndpoint: MakeGetRateEndpoint(svc),
	}
}

func MakeGetRateEndpoint(svc RateService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(messages.RateRequest)
		result, err := svc.GetRate(req)
		return result, err
	}
}

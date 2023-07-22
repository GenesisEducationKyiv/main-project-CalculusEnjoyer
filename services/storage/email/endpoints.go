package email

import (
	context1 "context"
	"storage/domain"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	AddEmail     endpoint.Endpoint
	GetAllEmails endpoint.Endpoint
}

func NewEndpointSet(svc EmailRepository) Endpoints {
	return Endpoints{
		AddEmail:     MakeAddEmailEndpoint(svc),
		GetAllEmails: MakeGetAllEmailsEndpoint(svc),
	}
}

func MakeGetAllEmailsEndpoint(svc EmailRepository) endpoint.Endpoint {
	return func(ctx context1.Context, request interface{}) (interface{}, error) {
		emails, err := svc.GetAll()
		return emails, err
	}
}

func MakeAddEmailEndpoint(svc EmailRepository) endpoint.Endpoint {
	return func(ctx context1.Context, request interface{}) (interface{}, error) {
		req := request.(domain.Email)
		err := svc.Add(req)
		return nil, err
	}
}

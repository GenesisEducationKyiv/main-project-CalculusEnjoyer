package emails

import (
	context1 "context"
	"storage/emails/messages"

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
		req := request.(messages.Email)
		err := svc.Add(req)
		return nil, err
	}
}

package transport

import (
	"context"
	"storage/domain"
	"storage/email"
	proto2 "storage/transport/proto"

	"google.golang.org/grpc/status"

	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	getAllEmails grpctransport.Handler
	addEmail     grpctransport.Handler
}

func NewGRPCServer(ep email.Endpoints) proto2.StorageServiceServer {
	return &grpcServer{
		addEmail: grpctransport.NewServer(
			ep.AddEmail,
			decodeGRPCAddEmailRequest,
			decodeGRPCAddEmailResponse,
		),
		getAllEmails: grpctransport.NewServer(
			ep.GetAllEmails,
			decodeGRPCGetAllEmailsRequest,
			decodeGRPCGetAllEmailsResponse,
		),
	}
}

func (g grpcServer) AddEmail(ctx context.Context, request *proto2.AddEmailRequest) (*proto2.AddEmailResponse, error) {
	_, rep, err := g.addEmail.ServeGRPC(ctx, request)
	if err != nil {
		return &proto2.AddEmailResponse{}, err
	}
	response := rep.(proto2.AddEmailResponse)
	return &response, nil
}

func (g grpcServer) GetAllEmails(ctx context.Context, request *proto2.GetAllEmailsRequest) (*proto2.GetAllEmailsResponse, error) {
	_, rep, err := g.getAllEmails.ServeGRPC(ctx, request)
	if err != nil {
		return nil, status.Error(status.Code(err), err.Error())
	}
	response := rep.(proto2.GetAllEmailsResponse)
	return &response, nil
}

func decodeGRPCAddEmailRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto2.AddEmailRequest)
	return domain.Email{Value: req.Email}, nil
}

func decodeGRPCAddEmailResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	if grpcRes == nil {
		return proto2.AddEmailResponse{}, nil
	}
	req := grpcRes.(error)
	return nil, status.Error(status.Code(req), req.Error())
}

func decodeGRPCGetAllEmailsRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto2.GetAllEmailsRequest)

	return req, nil
}

func decodeGRPCGetAllEmailsResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.([]domain.Email)

	var emails []string
	for i := range res {
		emails = append(emails, res[i].Value)
	}

	return proto2.GetAllEmailsResponse{Email: emails}, nil
}

package transport

import (
	"context"
	"storage/domain"
	"storage/email"
	"storage/transport/proto"

	"google.golang.org/protobuf/types/known/emptypb"

	"google.golang.org/grpc/status"

	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	getAllEmails   grpctransport.Handler
	addEmailRevert grpctransport.Handler
	addEmail       grpctransport.Handler
}

func NewGRPCServer(ep email.Endpoints) proto.StorageServiceServer {
	return &grpcServer{
		addEmail: grpctransport.NewServer(
			ep.AddEmail,
			decodeGRPCAddEmailRequest,
			decodeGRPCAddEmailResponse,
		),
		addEmailRevert: grpctransport.NewServer(
			ep.AddEmailRevert,
			decodeGRPCAddEmailRequest,
			decodeGRPCAddEmailRevertResponse,
		),
		getAllEmails: grpctransport.NewServer(
			ep.GetAllEmails,
			decodeGRPCGetAllEmailsRequest,
			decodeGRPCGetAllEmailsResponse,
		),
	}
}

func (g grpcServer) AddEmail(ctx context.Context, request *proto.AddEmailRequest) (*proto.AddEmailResponse, error) {
	_, rep, err := g.addEmail.ServeGRPC(ctx, request)
	if err != nil {
		return &proto.AddEmailResponse{}, err
	}
	response := rep.(proto.AddEmailResponse)
	return &response, nil
}

func (g grpcServer) AddEmailRevert(ctx context.Context, request *proto.AddEmailRequest) (*emptypb.Empty, error) {
	_, _, err := g.addEmailRevert.ServeGRPC(ctx, request)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

func (g grpcServer) GetAllEmails(ctx context.Context, request *proto.GetAllEmailsRequest) (*proto.GetAllEmailsResponse, error) {
	_, rep, err := g.getAllEmails.ServeGRPC(ctx, request)
	if err != nil {
		return nil, status.Error(status.Code(err), err.Error())
	}
	response := rep.(proto.GetAllEmailsResponse)
	return &response, nil
}

func decodeGRPCAddEmailRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.AddEmailRequest)
	return domain.Email{Value: req.Email}, nil
}

func decodeGRPCAddEmailResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	if grpcRes == nil {
		return proto.AddEmailResponse{}, nil
	}
	req := grpcRes.(error)
	return nil, status.Error(status.Code(req), req.Error())
}

func decodeGRPCGetAllEmailsRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.GetAllEmailsRequest)

	return req, nil
}

func decodeGRPCGetAllEmailsResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.([]domain.Email)

	var emails []string
	for i := range res {
		emails = append(emails, res[i].Value)
	}

	return proto.GetAllEmailsResponse{Email: emails}, nil
}

func decodeGRPCAddEmailRevertResponse(_ context.Context, _ interface{}) (interface{}, error) {
	return emptypb.Empty{}, nil
}

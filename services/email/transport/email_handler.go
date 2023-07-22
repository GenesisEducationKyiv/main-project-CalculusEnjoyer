package transport

import (
	"context"
	"email/dispatcher"
	domain2 "email/domain"
	proto2 "email/transport/proto"

	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	sendEmail grpctransport.Handler
}

func NewGRPCServer(ep dispatcher.Endpoints) proto2.EmailServiceServer {
	return &grpcServer{
		sendEmail: grpctransport.NewServer(
			ep.SendEmail,
			decodeGRPCSendEmailRequest,
			decodeGRPCSendEmailResponse,
		),
	}
}

func (g grpcServer) SendEmail(ctx context.Context, request *proto2.SendEmailRequest) (*proto2.SendEmailResponse, error) {
	_, rep, err := g.sendEmail.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	response := rep.(proto2.SendEmailResponse)
	return &response, nil
}

func decodeGRPCSendEmailRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto2.SendEmailRequest)
	return domain2.SendEmailRequest{Content: domain2.EmailContent{Subject: req.Subject, Body: req.Body}, To: req.To}, nil
}

func decodeGRPCSendEmailResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	if grpcReply == nil {
		return proto2.SendEmailResponse{}, nil
	}
	err := grpcReply.(error)
	return nil, err
}

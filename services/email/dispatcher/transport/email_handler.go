package transport

import (
	"context"
	"email/dispatcher"
	"email/dispatcher/executor"
	"email/dispatcher/executor/templates"
	"email/dispatcher/messages/proto"

	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	sendEmail grpctransport.Handler
}

func NewGRPCServer(ep dispatcher.Endpoints) proto.EmailServiceServer {
	return &grpcServer{
		sendEmail: grpctransport.NewServer(
			ep.SendEmail,
			decodeGRPCSendEmailRequest,
			decodeGRPCSendEmailResponse,
		),
	}
}

func (g grpcServer) SendEmail(ctx context.Context, request *proto.SendEmailRequest) (*proto.SendEmailResponse, error) {
	_, rep, err := g.sendEmail.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	response := rep.(proto.SendEmailResponse)
	return &response, nil
}

func decodeGRPCSendEmailRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.SendEmailRequest)
	return executor.SendEmailRequest{Content: templates.EmailContent{Subject: req.Subject, Body: req.Body}, To: req.To}, nil
}

func decodeGRPCSendEmailResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	if grpcReply == nil {
		return proto.SendEmailResponse{}, nil
	}
	err := grpcReply.(error)
	return nil, err
}

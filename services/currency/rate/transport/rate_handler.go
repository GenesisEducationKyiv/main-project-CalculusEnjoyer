package transport

import (
	"context"
	"currency/rate"
	"currency/rate/messages"
	"currency/rate/messages/proto"

	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	getRate grpctransport.Handler
}

func NewGRPCServer(ep rate.Endpoints) proto.RateServiceServer {
	return &grpcServer{
		getRate: grpctransport.NewServer(
			ep.GetRateEndpoint,
			decodeGRPCRateRequest,
			decodeGRPCGetResponse,
		),
	}
}

func (g *grpcServer) GetRate(ctx context.Context, r *proto.RateRequest) (*proto.RateResponse, error) {
	_, rep, err := g.getRate.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	response := rep.(proto.RateResponse)
	return &response, err
}

func decodeGRPCRateRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.RateRequest)
	return messages.RateRequest{BaseCurrency: messages.Currency(req.BaseCurrency), TargetCurrency: messages.Currency(req.TargetCurrency)}, nil
}

func decodeGRPCGetResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(messages.RateResult)
	return proto.RateResponse{Rate: reply.Rate, Timestamp: reply.Timestamp.UTC().String()}, nil
}

package transport

import (
	"context"
	domain2 "currency/domain"
	"currency/logger"
	"currency/rate"
	proto2 "currency/transport/proto"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	getRate grpctransport.Handler
}

func NewGRPCServer(ep rate.Endpoints) proto2.RateServiceServer {
	return &grpcServer{
		getRate: grpctransport.NewServer(
			ep.GetRateEndpoint,
			decodeGRPCRateRequest,
			decodeGRPCGetResponse,
		),
	}
}

func (g *grpcServer) GetRate(ctx context.Context, r *proto2.RateRequest) (*proto2.RateResponse, error) {
	_, rep, err := g.getRate.ServeGRPC(ctx, r)
	logger.DefaultLog(logger.DEBUG, "receiving gRPC request")
	if err != nil {
		return nil, err
	}
	response := rep.(proto2.RateResponse)
	return &response, err
}

func decodeGRPCRateRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto2.RateRequest)
	logger.DefaultLog(logger.DEBUG, "decoding gRPC request")
	return domain2.RateRequest{BaseCurrency: domain2.Currency(req.BaseCurrency), TargetCurrency: domain2.Currency(req.TargetCurrency)}, nil
}

func decodeGRPCGetResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(domain2.RateResult)
	logger.DefaultLog(logger.DEBUG, "decoding gRPC response")
	return proto2.RateResponse{Rate: reply.Rate, Timestamp: reply.Timestamp.UTC().String()}, nil
}

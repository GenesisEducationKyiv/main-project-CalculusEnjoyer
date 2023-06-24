package main

import (
	"currency/rate"
	"currency/rate/messages/proto"
	"currency/rate/transport"
	"currency/test"
	"log"
	"net"
	"strconv"
	"time"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

func main() {
	runTest()
}

func runTest() {
	conf := test.LoadFromENV()

	service := rate.NewRateService(&stubRateProvider{TestRate: conf.TestRate}, &stubTimeProvider{TestTimeStamp: conf.TestUnixTimeStamp})
	eps := rate.NewEndpointSet(service)
	grpcServer := transport.NewGRPCServer(eps)
	baseServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
	proto.RegisterRateServiceServer(baseServer, grpcServer)
	lis, err := net.Listen(conf.Network, ":"+strconv.Itoa(conf.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	if err = baseServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

type stubRateProvider struct {
	TestRate float64
}

func (r *stubRateProvider) GetExchangeRate(baseCurrency, targetCurrency string) (rate float64, err error) {
	return r.TestRate, nil
}

type stubTimeProvider struct {
	TestTimeStamp time.Time
}

func (t *stubTimeProvider) Now() time.Time {
	return t.TestTimeStamp
}

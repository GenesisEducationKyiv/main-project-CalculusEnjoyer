package main

import (
	"currency/config"
	"currency/rate"
	"currency/rate/messages/proto"
	"currency/rate/providers/crypto"
	"currency/rate/providers/time"
	"currency/rate/transport"
	"log"
	"net"
	"strconv"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

func main() {
	run()
}

func run() {
	conf := config.LoadFromENV()

	service := rate.NewRateService(crypto.NewCoinGeckoRateProvider(conf), &time.SystemTime{})
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

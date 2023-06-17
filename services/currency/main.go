package main

import (
	"currency/rate"
	"currency/rate/messages/proto"
	"currency/rate/providsers"
	"currency/rate/transport"
	"log"
	"net"
	"os"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	run()
}

func run() {
	err := godotenv.Load()
	if err != err {
		log.Fatalf("Failed to load configs")
	}
	network := os.Getenv("NETWORK")
	port := os.Getenv("PORT")

	service := rate.NewService(&providsers.CoinGeckoRateProvider{})
	eps := rate.NewEndpointSet(service)
	grpcServer := transport.NewGRPCServer(eps)
	baseServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
	proto.RegisterRateServiceServer(baseServer, grpcServer)
	lis, err := net.Listen(network, ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	if err := baseServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

package main

import (
	"email/dispatcher"
	sender "email/dispatcher/executor"
	"email/dispatcher/messages/proto"
	"email/dispatcher/transport"
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
	if err != nil {
		log.Fatalf("Can not load .env config")
	}
	network := os.Getenv("NETWORK")
	port := os.Getenv("PORT")

	service := dispatcher.NewService(sender.NewGoSender())
	eps := dispatcher.NewEndpointSet(service)
	grpcServer := transport.NewGRPCServer(eps)
	baseServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
	proto.RegisterEmailServiceServer(baseServer, grpcServer)
	lis, err := net.Listen(network, ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	if err := baseServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

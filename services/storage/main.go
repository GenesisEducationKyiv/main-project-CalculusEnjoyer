package main

import (
	"log"
	"net"
	"os"
	"storage/emails"
	"storage/emails/messages/proto"
	"storage/emails/transport"

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
		log.Fatalf("Error getting the .env config")
	}
	network := os.Getenv("NETWORK")
	port := os.Getenv("PORT")

	service := emails.NewService()
	eps := emails.NewEndpointSet(service)
	grpcServer := transport.NewGRPCServer(eps)
	baseServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
	proto.RegisterStorageServiceServer(baseServer, grpcServer)
	lis, err := net.Listen(network, ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	if err := baseServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

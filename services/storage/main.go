package main

import (
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"storage/emails"
	"storage/emails/messages/proto"
	"storage/emails/transport"
)

var StorageFile *os.File

func main() {
	run()
}

func run() {
	err := godotenv.Load()
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

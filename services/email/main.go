package main

import (
	"email/config"
	"email/dispatcher"
	sender "email/dispatcher/executor"
	"email/dispatcher/messages/proto"
	"email/dispatcher/transport"
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

	service := dispatcher.NewService(sender.NewGoSender(conf))
	eps := dispatcher.NewEndpointSet(service)
	grpcServer := transport.NewGRPCServer(eps)
	baseServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
	proto.RegisterEmailServiceServer(baseServer, grpcServer)

	lis, err := net.Listen(conf.Network, ":"+strconv.Itoa(conf.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	if err := baseServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

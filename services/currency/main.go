package main

import (
	"currency/config"
	logger "currency/logger"
	"currency/rate"
	"currency/rate/providers/crypto"
	"currency/rate/providers/time"
	"currency/transport"
	"currency/transport/proto"
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

	loggerInstance := logger.NewBrokerLogger(&time.SystemTime{}, conf)
	loggerInstance.Init()
	logger.SetDefaultLogger(loggerInstance)

	service := rate.NewRateService(rate.NewCachedProvider(bootstrapRateProviders(conf, loggerInstance), conf), &time.SystemTime{})
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

func bootstrapRateProviders(conf config.Config, logger logger.Logger) *rate.RateLink {
	kunaLink := rate.NewRateLink(rate.NewRateLogger(crypto.NewKunaRateProvider(conf), logger))
	coinApiLink := rate.NewRateLink(rate.NewRateLogger(crypto.NewCoinAPIProvider(conf), logger))
	coinGeckoLink := rate.NewRateLink(rate.NewRateLogger(crypto.NewCoinGeckoRateProvider(conf), logger))

	kunaLink.SetNextLink(coinApiLink)
	coinApiLink.SetNextLink(coinGeckoLink)

	return kunaLink
}

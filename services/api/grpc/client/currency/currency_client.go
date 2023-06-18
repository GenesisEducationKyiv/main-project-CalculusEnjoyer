package currency

import (
	"context"
	"currency/rate/messages/proto"
	"log"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

type CurrencyGRPCClient struct {
	network string
	port    string
}

func NewCurrencyGRPCClient() *CurrencyGRPCClient {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Can not load .env config")
	}

	return &CurrencyGRPCClient{
		os.Getenv("CURRENCY_NETWORK"),
		os.Getenv("CURRENCY_SERVICE_PORT"),
	}
}

func (c *CurrencyGRPCClient) GetRate(request *proto.RateRequest) (*proto.RateResponse, error) {
	conn := c.getConnection()
	defer conn.Close()

	client := proto.NewRateServiceClient(conn)

	response, err := client.GetRate(context.Background(), request)
	if err != nil {
		log.Printf("Failed to call GetRate: %v", err)
		return nil, err
	}

	return response, err
}

func (c *CurrencyGRPCClient) getConnection() *grpc.ClientConn {
	conn, err := grpc.Dial(c.network+":"+c.port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	return conn
}

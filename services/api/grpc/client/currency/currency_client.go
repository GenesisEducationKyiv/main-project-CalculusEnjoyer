package currency

import (
	"api/config"
	"api/models"
	"context"
	"currency/rate/messages/proto"
	"log"
	"strconv"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type CurrencyGRPCClient struct {
	conf config.Config
}

func NewCurrencyGRPCClient(conf config.Config) *CurrencyGRPCClient {
	return &CurrencyGRPCClient{conf: conf}
}

func (c *CurrencyGRPCClient) GetRate(request *models.RateRequest) (*models.RateResponse, error) {
	conn := c.getConnection()
	defer conn.Close()

	client := proto.NewRateServiceClient(conn)

	response, err := client.GetRate(context.Background(), modelRateToProto(request))
	if err != nil {
		return nil, errors.Wrap(err, "can not get rate")
	}

	return protoRateToModel(response), err
}

func (c *CurrencyGRPCClient) getConnection() *grpc.ClientConn {
	conn, err := grpc.Dial(c.conf.CurrencyNetwork+":"+strconv.Itoa(c.conf.CurrencyPort), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	return conn
}

func protoRateToModel(response *proto.RateResponse) *models.RateResponse {
	return &models.RateResponse{
		Rate: response.Rate,
	}
}

func modelRateToProto(request *models.RateRequest) *proto.RateRequest {
	return &proto.RateRequest{
		BaseCurrency:   request.BaseCurrency,
		TargetCurrency: request.TargetCurrency,
	}
}

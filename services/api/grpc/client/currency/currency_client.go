package currency

import (
	"api/config"
	"api/models"
	"context"
	"currency/rate/messages/proto"
	"log"
	"strconv"

	"google.golang.org/grpc/connectivity"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type CurrencyGRPCClient struct {
	conf config.Config
	conn *grpc.ClientConn
}

func NewCurrencyGRPCClient(conf config.Config) *CurrencyGRPCClient {
	client := CurrencyGRPCClient{conf: conf}
	client.conn = client.getConnection()
	return &client
}

func (c *CurrencyGRPCClient) GetRate(request *models.RateRequest, cnx context.Context) (*models.RateResponse, error) {
	conn := c.connection()

	client := proto.NewRateServiceClient(conn)

	response, err := client.GetRate(cnx, modelRateToProto(request))
	if err != nil {
		return nil, errors.Wrap(err, "can not get rate")
	}

	return protoRateToModel(response), err
}

func (c *CurrencyGRPCClient) getConnection() *grpc.ClientConn {
	conn, err := grpc.Dial(c.conf.CurrencyNetwork+":"+strconv.Itoa(c.conf.CurrencyPort), grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to connect: %v", err)
	}

	return conn
}

func (c *CurrencyGRPCClient) connection() *grpc.ClientConn {
	if c.conn != nil && c.conn.GetState() == connectivity.Ready {
		return c.conn
	} else {
		c.conn = c.getConnection()
		return c.conn
	}
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

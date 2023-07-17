package currency

import (
	"api/config"
	"api/domain"
	"api/logger"
	"context"
	"currency/transport/proto"
	"strconv"

	"google.golang.org/grpc/credentials/insecure"

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
	client.conn, _ = openConnection(conf.CurrencyNetwork, conf.CurrencyPort)
	return &client
}

func (c *CurrencyGRPCClient) GetRate(request domain.RateRequest, cnx context.Context) (*domain.RateResponse, error) {
	conn, err := c.connection()
	logger.DefaultLog(logger.DEBUG, "trying to get rate from gRPC client")
	if err != nil {
		logger.DefaultLog(logger.ERROR, "fail to get connection")
		return nil, errors.Wrap(err, "fail to get connection")
	}

	client := proto.NewRateServiceClient(conn)

	proto, err := c.modelRateToProto(&request)
	if err != nil {
		logger.DefaultLog(logger.ERROR, "can not make proto request for getting rate")
		return nil, errors.Wrap(err, "can not make proto request for getting rate")
	}

	response, err := client.GetRate(cnx, proto)
	if err != nil {
		logger.DefaultLog(logger.ERROR, "can not get rate")
		return nil, errors.Wrap(err, "can not get rate")
	}

	return protoRateToModel(response), err
}

func (c *CurrencyGRPCClient) connection() (*grpc.ClientConn, error) {
	logger.DefaultLog(logger.DEBUG, "using connection to the currency service")
	if c.conn != nil && c.conn.GetState() == connectivity.Ready {
		return c.conn, nil
	}

	con, err := openConnection(c.conf.CurrencyNetwork, c.conf.CurrencyPort)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get conn")
	}

	c.conn = con
	return c.conn, nil
}

func openConnection(network string, port int) (*grpc.ClientConn, error) {
	insecureHack := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.Dial(network+":"+strconv.Itoa(port), insecureHack)

	return conn, errors.Wrap(err, "failed to grpc connect")
}

func protoRateToModel(response *proto.RateResponse) *domain.RateResponse {
	return &domain.RateResponse{
		Rate: response.Rate,
	}
}

func (c *CurrencyGRPCClient) modelRateToProto(request *domain.RateRequest) (*proto.RateRequest, error) {
	return &proto.RateRequest{
		BaseCurrency:   string(request.BaseCurrency),
		TargetCurrency: string(request.TargetCurrency),
	}, nil
}

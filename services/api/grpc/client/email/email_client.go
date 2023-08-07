package email

import (
	"api/config"
	"api/domain"
	"api/logger"
	"context"
	"email/transport/proto"
	"strconv"

	"github.com/pkg/errors"
	"google.golang.org/grpc/credentials/insecure"

	"google.golang.org/grpc/connectivity"

	"google.golang.org/grpc"
)

type EmailGRPCClient struct {
	conf config.Config
	conn *grpc.ClientConn
}

func NewEmailGRPCClient(conf config.Config) *EmailGRPCClient {
	client := EmailGRPCClient{conf: conf}
	client.conn, _ = openConnection(conf.EmailNetwork, conf.EmailPort)
	return &client
}

func (c *EmailGRPCClient) SendEmail(request domain.SendEmailsRequest, cnx context.Context) error {
	conn, err := c.connection()
	logger.DefaultLog(logger.DEBUG, "trying to send emails from gRPC client")
	if err != nil {
		logger.DefaultLog(logger.ERROR, "failed to send emails from gRPC client")
		return errors.Wrap(err, "can not get connection SendEmail")
	}

	client := proto.NewEmailServiceClient(conn)

	_, err = client.SendEmail(cnx, modelSendEmailsToProto(request))

	return err
}

func (c *EmailGRPCClient) connection() (*grpc.ClientConn, error) {
	if c.conn != nil && c.conn.GetState() == connectivity.Ready {
		return c.conn, nil
	}

	con, err := openConnection(c.conf.EmailNetwork, c.conf.EmailPort)
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

func modelSendEmailsToProto(request domain.SendEmailsRequest) *proto.SendEmailRequest {
	return &proto.SendEmailRequest{
		Subject: request.Template.Subject,
		Body:    request.Template.Body,
		To:      request.Interceptor.Value,
	}
}

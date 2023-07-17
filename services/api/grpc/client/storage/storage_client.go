package storage

import (
	"api/aerror"
	"api/config"
	"api/domain"
	"api/logger"
	"context"
	"storage/transport/proto"
	"strconv"

	"google.golang.org/grpc/credentials/insecure"

	"google.golang.org/grpc/connectivity"

	"github.com/pkg/errors"

	"google.golang.org/grpc"
)

type StorageGRPCClient struct {
	conf config.Config
	conn *grpc.ClientConn
}

func NewStorageGRPCClient(conf config.Config) *StorageGRPCClient {
	client := StorageGRPCClient{conf: conf}
	client.conn, _ = openConnection(conf.StorageNetwork, conf.StoragePort)
	return &client
}

func (c *StorageGRPCClient) AddEmail(request domain.AddEmailRequest, cnx context.Context) error {
	logger.DefaultLog(logger.DEBUG, "trying to add emails from gRPC client")
	conn, err := c.connection()
	if err != nil {
		return errors.Wrap(err, "failed to get connection")
	}

	client := proto.NewStorageServiceClient(conn)

	_, err = client.AddEmail(cnx, modelAddEmailToProto(request))
	return err
}

func (c *StorageGRPCClient) GetAllEmails(cnx context.Context) ([]domain.Email, error) {
	logger.DefaultLog(logger.DEBUG, "trying to get emails from gRPC client")
	conn, err := c.connection()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get connection()")
	}

	client := proto.NewStorageServiceClient(conn)

	response, err := client.GetAllEmails(cnx, &proto.GetAllEmailsRequest{})
	if err != nil {
		return nil, errors.Wrap(err, aerror.ErrGRPC.Error())
	}

	return protoEmailsToSlice(response), nil
}

func (c *StorageGRPCClient) connection() (*grpc.ClientConn, error) {
	if c.conn != nil && c.conn.GetState() == connectivity.Ready {
		return c.conn, nil
	}

	con, err := openConnection(c.conf.StorageNetwork, c.conf.StoragePort)
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

func modelAddEmailToProto(request domain.AddEmailRequest) *proto.AddEmailRequest {
	return &proto.AddEmailRequest{
		Email: request.Email.Value,
	}
}

func protoEmailsToSlice(response *proto.GetAllEmailsResponse) []domain.Email {
	emails := make([]domain.Email, len(response.Email))
	for i, email := range emails {
		email.Value = response.Email[i]
		emails[i] = email
	}

	return emails
}

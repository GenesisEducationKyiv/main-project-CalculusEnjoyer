package storage

import (
	"api/aerror"
	"api/config"
	"api/models"
	"context"
	"log"
	"storage/emails/messages/proto"
	"strconv"

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
	client.conn = client.getConnection()
	return &client
}

func (c *StorageGRPCClient) AddEmail(request models.AddEmailRequest) error {
	conn := c.connection()

	client := proto.NewStorageServiceClient(conn)

	_, err := client.AddEmail(context.Background(), modelAddEmailToProto(request))
	return err
}

func (c *StorageGRPCClient) GetAllEmails() ([]models.Email, error) {
	conn := c.connection()

	client := proto.NewStorageServiceClient(conn)

	response, err := client.GetAllEmails(context.Background(), &proto.GetAllEmailsRequest{})
	if err != nil {
		errors.Wrap(err, aerror.ErrGRPC.Error())
	}

	return protoEmailsToSlice(response), nil
}

func (c *StorageGRPCClient) connection() *grpc.ClientConn {
	if c.conn != nil && c.conn.GetState() == connectivity.Ready {
		return c.conn
	} else {
		c.conn = c.getConnection()
		return c.conn
	}
}

func (c *StorageGRPCClient) getConnection() *grpc.ClientConn {
	conn, err := grpc.Dial(c.conf.StorageNetwork+":"+strconv.Itoa(c.conf.StoragePort), grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to connect: %v", err)
	}

	return conn
}

func modelAddEmailToProto(request models.AddEmailRequest) *proto.AddEmailRequest {
	return &proto.AddEmailRequest{
		Email: request.Email.Value,
	}
}

func protoEmailsToSlice(response *proto.GetAllEmailsResponse) []models.Email {
	emails := make([]models.Email, len(response.Email))
	for i, email := range emails {
		email.Value = response.Email[i]
		emails[i] = email
	}

	return emails
}

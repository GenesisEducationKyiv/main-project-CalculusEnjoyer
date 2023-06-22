package storage

import (
	"api/aerror"
	"api/config"
	"api/models"
	"context"
	"storage/emails/messages/proto"
	"strconv"

	"github.com/pkg/errors"

	"google.golang.org/grpc"
)

type StorageGRPCClient struct {
	conf config.Config
}

func NewStorageGRPCClient(conf config.Config) *StorageGRPCClient {
	return &StorageGRPCClient{conf: conf}
}

func (c *StorageGRPCClient) AddEmail(request models.AddEmailRequest) error {
	conn := c.getConnection()
	defer conn.Close()

	client := proto.NewStorageServiceClient(conn)

	_, err := client.AddEmail(context.Background(), modelAddEmailToProto(request))
	return err
}

func (c *StorageGRPCClient) GetAllEmails() ([]models.Email, error) {
	conn := c.getConnection()
	defer conn.Close()

	client := proto.NewStorageServiceClient(conn)

	response, err := client.GetAllEmails(context.Background(), &proto.GetAllEmailsRequest{})
	if err != nil {
		errors.Wrap(err, aerror.ErrGRPC.Error())
	}

	return protoEmailsToSlice(response), nil
}

func (c *StorageGRPCClient) getConnection() *grpc.ClientConn {
	conn, err := grpc.Dial(c.conf.StorageNetwork+":"+strconv.Itoa(c.conf.StoragePort), grpc.WithInsecure())
	if err != nil {
		errors.Wrap(err, aerror.ErrGRPC.Error())
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

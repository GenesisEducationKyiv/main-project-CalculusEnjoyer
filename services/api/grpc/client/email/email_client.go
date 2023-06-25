package email

import (
	"api/config"
	"api/models"
	"context"
	"email/dispatcher/messages/proto"
	"log"
	"strconv"

	"google.golang.org/grpc/connectivity"

	"google.golang.org/grpc"
)

type EmailGRPCClient struct {
	conf config.Config
	conn *grpc.ClientConn
}

func NewEmailGRPCClient(conf config.Config) *EmailGRPCClient {
	client := EmailGRPCClient{conf: conf}
	client.conn = client.getConnection()
	return &client
}

func (c *EmailGRPCClient) SendEmail(request models.SendEmailsRequest) error {
	conn := c.connection()

	client := proto.NewEmailServiceClient(conn)

	_, err := client.SendEmail(context.Background(), modelSendEmailstoProto(request))

	return err
}

func (c *EmailGRPCClient) getConnection() *grpc.ClientConn {
	conn, err := grpc.Dial(c.conf.EmailNetwork+":"+strconv.Itoa(c.conf.EmailPort), grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to connect: %v", err)
	}

	return conn
}

func (c *EmailGRPCClient) connection() *grpc.ClientConn {
	if c.conn != nil && c.conn.GetState() == connectivity.Ready {
		return c.conn
	} else {
		c.conn = c.getConnection()
		return c.conn
	}
}

func modelSendEmailstoProto(request models.SendEmailsRequest) *proto.SendEmailRequest {
	return &proto.SendEmailRequest{
		Subject: request.Template.Subject,
		Body:    request.Template.Body,
		To:      request.Interceptor.Value,
	}
}

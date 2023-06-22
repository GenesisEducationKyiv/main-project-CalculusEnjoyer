package email

import (
	"api/config"
	"api/models"
	"context"
	"email/dispatcher/messages/proto"
	"log"
	"strconv"

	"google.golang.org/grpc"
)

type EmailGRPCClient struct {
	conf config.Config
}

func NewEmailGRPCClient(conf config.Config) *EmailGRPCClient {
	return &EmailGRPCClient{conf: conf}
}

func (c *EmailGRPCClient) SendEmail(request models.SendEmailsRequest) error {
	conn := c.getConnection()
	defer conn.Close()

	client := proto.NewEmailServiceClient(conn)

	_, err := client.SendEmail(context.Background(), modelSendEmailstoProto(request))

	return err
}

func (c *EmailGRPCClient) getConnection() *grpc.ClientConn {
	conn, err := grpc.Dial(c.conf.EmailNetwork+":"+strconv.Itoa(c.conf.EmailPort), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	return conn
}

func modelSendEmailstoProto(request models.SendEmailsRequest) *proto.SendEmailRequest {
	return &proto.SendEmailRequest{
		Subject: request.Template.Subject,
		Body:    request.Template.Body,
		To:      request.Interceptor.Value,
	}
}

package email

import (
	"context"
	"email/dispatcher/messages/proto"
	"log"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

type EmailGRPCClient struct {
	network string
	port    string
}

func NewEmailGRPCClient() *EmailGRPCClient {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Can not load .env config")
	}

	return &EmailGRPCClient{
		network: os.Getenv("EMAIL_NETWORK"),
		port:    os.Getenv("EMAIL_SERVICE_PORT"),
	}
}

func (c *EmailGRPCClient) SendEmail(request proto.SendEmailRequest) error {
	conn := c.getConnection()
	defer conn.Close()

	client := proto.NewEmailServiceClient(conn)

	_, err := client.SendEmail(context.Background(), &request)

	return err
}

func (c *EmailGRPCClient) getConnection() *grpc.ClientConn {
	conn, err := grpc.Dial(c.network+":"+c.port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	return conn
}

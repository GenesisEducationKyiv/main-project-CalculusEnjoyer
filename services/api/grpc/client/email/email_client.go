package email

import (
	"context"
	"email/dispatcher/messages/proto"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"os"
)

var network string
var port string
var client proto.EmailServiceClient
var connection *grpc.ClientConn

type EmailGRPCClient struct{}

func init() {
	_ = godotenv.Load()
	network = os.Getenv("EMAIL_NETWORK")
	port = os.Getenv("EMAIL_SERVICE_PORT")
}

func (c *EmailGRPCClient) SendEmail(request proto.SendEmailRequest) error {
	conn := c.getConnection()
	defer conn.Close()

	client = proto.NewEmailServiceClient(conn)

	_, err := client.SendEmail(context.Background(), &request)

	return err
}

func (c *EmailGRPCClient) getConnection() *grpc.ClientConn {
	conn, err := grpc.Dial(network+":"+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	return conn
}

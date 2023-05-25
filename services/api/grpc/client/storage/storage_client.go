package storage

import (
	"context"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"os"
	"storage/emails/messages/proto"
)

var storageNetwork string
var port string
var client proto.StorageServiceClient
var connection *grpc.ClientConn

type StorageGRPCClient struct{}

func init() {
	_ = godotenv.Load()
	storageNetwork = os.Getenv("STORAGE_NETWORK")
	port = os.Getenv("STORAGE_SERVICE_PORT")
}

func (c *StorageGRPCClient) AddEmail(request proto.AddEmailRequest) (proto.AddEmailResponse, error) {
	conn := c.getConnection()
	defer conn.Close()

	client = proto.NewStorageServiceClient(conn)

	response, err := client.AddEmail(context.Background(), &request)
	if response == nil {
		return proto.AddEmailResponse{}, err
	}
	return *response, err
}

func (c *StorageGRPCClient) GetAllEmails(request proto.GetAllEmailsRequest) proto.GetAllEmailsResponse {
	conn := c.getConnection()
	defer conn.Close()

	client = proto.NewStorageServiceClient(conn)

	response, err := client.GetAllEmails(context.Background(), &request)
	if err != nil {
		log.Fatalf("Failed to call GetRate: %v", err)
	}

	return *response
}

func (c *StorageGRPCClient) getConnection() *grpc.ClientConn {
	conn, err := grpc.Dial(storageNetwork+":"+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	return conn
}

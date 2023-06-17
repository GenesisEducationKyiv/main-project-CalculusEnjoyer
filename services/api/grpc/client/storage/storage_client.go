package storage

import (
	"context"
	"log"
	"os"
	"storage/emails/messages/proto"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

type StorageGRPCClient struct {
	storageNetwork string
	port           string
}

func NewStorageGRPCClient() *StorageGRPCClient {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Can not load .env config")
	}

	return &StorageGRPCClient{
		os.Getenv("STORAGE_NETWORK"),
		os.Getenv("STORAGE_SERVICE_PORT"),
	}
}

func (c *StorageGRPCClient) AddEmail(request proto.AddEmailRequest) (proto.AddEmailResponse, error) {
	conn := c.getConnection()
	defer conn.Close()

	client := proto.NewStorageServiceClient(conn)

	response, err := client.AddEmail(context.Background(), &request)
	if response == nil {
		return proto.AddEmailResponse{}, err
	}
	return *response, err
}

func (c *StorageGRPCClient) GetAllEmails(request proto.GetAllEmailsRequest) proto.GetAllEmailsResponse {
	conn := c.getConnection()
	defer conn.Close()

	client := proto.NewStorageServiceClient(conn)

	response, err := client.GetAllEmails(context.Background(), &request)
	if err != nil {
		log.Fatalf("Failed to call GetRate: %v", err)
	}

	return *response
}

func (c *StorageGRPCClient) getConnection() *grpc.ClientConn {
	conn, err := grpc.Dial(c.storageNetwork+":"+c.port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	return conn
}

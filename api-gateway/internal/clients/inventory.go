package clients

import (
	"api-gateway/internal/api/proto/inventory"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type InventoryClient struct {
	inventory.InventoryServiceClient
}

func NewInventoryClient(address string) *InventoryClient {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}

	return &InventoryClient{inventory.NewInventoryServiceClient(conn)}
}

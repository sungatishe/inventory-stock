package server

import (
	"fmt"
	"google.golang.org/grpc"
	"inventory-service/api/proto"
	"log"
	"net"
)

type InventoryServer struct {
	server           *grpc.Server
	inventoryService proto.InventoryServiceServer
}

func NewInventoryServer(inventoryService proto.InventoryServiceServer) *InventoryServer {
	grpcServer := grpc.NewServer()
	proto.RegisterInventoryServiceServer(grpcServer, inventoryService)

	return &InventoryServer{
		server:           grpcServer,
		inventoryService: inventoryService,
	}
}

func (s *InventoryServer) Run(port string) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	fmt.Printf("Inventory service is running on port %s\n", port)
	if err := s.server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

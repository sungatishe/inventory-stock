package server

import (
	"auth-service/api/proto"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

type AuthServer struct {
	server      *grpc.Server
	authService proto.AuthServiceServer
}

func NewAuthServer(authService proto.AuthServiceServer) *AuthServer {
	grpcServer := grpc.NewServer()
	proto.RegisterAuthServiceServer(grpcServer, authService)

	return &AuthServer{
		server:      grpcServer,
		authService: authService,
	}
}

func (s *AuthServer) Run(port string) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	fmt.Printf("Auth service is running on port %s\n", port)
	if err := s.server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

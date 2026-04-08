package server

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	Server *grpc.Server
}

func NewGrpcServer() *GRPCServer {
	server := grpc.NewServer()
	return &GRPCServer{
		Server: server,
	}
}

func (s *GRPCServer) Run(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	return s.Server.Serve(lis)
}

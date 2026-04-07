package server

import (
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

func (s *GRPCServer) Run(port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	return s.Server.Serve(lis)
}

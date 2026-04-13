package server

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	Server *grpc.Server
}

func NewGrpcServer() *GRPCServer {
	server := grpc.NewServer(grpc.UnaryInterceptor(
		func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
			return nil, status.Error(codes.Unavailable, "something went wrong")
		},
	))
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

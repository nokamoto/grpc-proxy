package main

import (
	"fmt"
	pb "github.com/nokamoto/grpc-proxy/examples/ping"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
)

type pingService struct{}

func (s *pingService) Send(_ context.Context, _ *pb.Ping) (*pb.Pong, error) {
	return &pb.Pong{}, nil
}

func (s *pingService) SendStreamC(_ pb.PingService_SendStreamCServer) error {
	return status.Error(codes.Unimplemented, "not implemented yet")
}

func (s *pingService) SendStreamS(_ *pb.Ping, _ pb.PingService_SendStreamSServer) error {
	return status.Error(codes.Unimplemented, "not implemented yet")
}

func (s *pingService) SendStreamB(_ pb.PingService_SendStreamBServer) error {
	return status.Error(codes.Unimplemented, "not implemented yet")
}

func withPingServer(f func() error) error {
	port := 9002

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	opts := []grpc.ServerOption{}
	srv := grpc.NewServer(opts...)
	svc := &pingService{}

	pb.RegisterPingServiceServer(srv, svc)

	go func() {
		srv.Serve(lis)
	}()
	defer srv.GracefulStop()

	return f()
}

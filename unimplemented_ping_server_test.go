package main

import (
	pb "github.com/nokamoto/grpc-proxy/examples/ping"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type unimplementedPingService struct{}

func (s *unimplementedPingService) Send(_ context.Context, _ *pb.Ping) (*pb.Pong, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "[ping] not implemented yet")
}

func (s *unimplementedPingService) SendStreamC(stream pb.PingService_SendStreamCServer) error {
	return grpc.Errorf(codes.Unimplemented, "[ping] not implemented yet")
}

func (s *unimplementedPingService) SendStreamS(_ *pb.Ping, stream pb.PingService_SendStreamSServer) error {
	return grpc.Errorf(codes.Unimplemented, "[ping] not implemented yet")
}

func (s *unimplementedPingService) SendStreamB(stream pb.PingService_SendStreamBServer) error {
	return grpc.Errorf(codes.Unimplemented, "[ping] not implemented yet")
}

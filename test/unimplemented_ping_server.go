package test

import (
	pb "github.com/nokamoto/grpc-proxy/examples/ping"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnimplementedPingServer implements pb.PingServiceServer only for test.
type UnimplementedPingServer struct{}

// Send returns Unimplemented.
func (*UnimplementedPingServer) Send(_ context.Context, m *pb.Ping) (*pb.Pong, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented yet")
}

// SendStreamC returns Unimplemented.
func (s *UnimplementedPingServer) SendStreamC(stream pb.PingService_SendStreamCServer) error {
	return status.Error(codes.Unimplemented, "not implemented yet")
}

// SendStreamS returns Unimplemented.
func (s *UnimplementedPingServer) SendStreamS(m *pb.Ping, stream pb.PingService_SendStreamSServer) error {
	return status.Error(codes.Unimplemented, "not implemented yet")
}

// SendStreamB returns Unimplemented.
func (s *UnimplementedPingServer) SendStreamB(stream pb.PingService_SendStreamBServer) error {
	return status.Error(codes.Unimplemented, "not implemented yet")
}

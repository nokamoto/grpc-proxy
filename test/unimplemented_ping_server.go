package test

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnimplementedPingServer implements pb.PingServiceServer only for test.
type UnimplementedPingServer struct{}

// Send returns Unimplemented.
func (*UnimplementedPingServer) Send(_ context.Context, m *Ping) (*Pong, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented yet")
}

// SendStreamC returns Unimplemented.
func (s *UnimplementedPingServer) SendStreamC(stream PingService_SendStreamCServer) error {
	return status.Error(codes.Unimplemented, "not implemented yet")
}

// SendStreamS returns Unimplemented.
func (s *UnimplementedPingServer) SendStreamS(m *Ping, stream PingService_SendStreamSServer) error {
	return status.Error(codes.Unimplemented, "not implemented yet")
}

// SendStreamB returns Unimplemented.
func (s *UnimplementedPingServer) SendStreamB(stream PingService_SendStreamBServer) error {
	return status.Error(codes.Unimplemented, "not implemented yet")
}

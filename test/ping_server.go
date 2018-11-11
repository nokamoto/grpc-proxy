package test

import (
	pb "github.com/nokamoto/grpc-proxy/examples/ping"
	"golang.org/x/net/context"
	"io"
	"time"
)

// PingServer implements pb.PingServiceServer only for test.
type PingServer struct{}

func (s *PingServer) pong(source []*pb.Ping) *pb.Pong {
	return &pb.Pong{Source: source, Ts: time.Now().Unix()}
}

// Send returns pb.Pong.
func (s *PingServer) Send(_ context.Context, m *pb.Ping) (*pb.Pong, error) {
	return s.pong([]*pb.Ping{m}), nil
}

// SendStreamC returns pb.Pong.
func (s *PingServer) SendStreamC(stream pb.PingService_SendStreamCServer) error {
	source := make([]*pb.Ping, 0)
	for {
		m, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(s.pong(source))
		}

		if err != nil {
			return err
		}

		source = append(source, m)
	}
}

// SendStreamS returns pb.Pong.
func (s *PingServer) SendStreamS(m *pb.Ping, stream pb.PingService_SendStreamSServer) error {
	for i := 0; i < 10; i++ {
		err := stream.Send(s.pong([]*pb.Ping{m}))
		if err != nil {
			return err
		}
	}
	return nil
}

// SendStreamB returns pb.Pong.
func (s *PingServer) SendStreamB(stream pb.PingService_SendStreamBServer) error {
	for {
		m, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		err = stream.Send(s.pong([]*pb.Ping{m}))
		if err != nil {
			return err
		}
	}
}

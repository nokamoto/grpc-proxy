package test

import (
	"golang.org/x/net/context"
	"io"
	"time"
)

// PingServer implements PingServiceServer only for test.
type PingServer struct{}

func (s *PingServer) pong(source []*Ping) *Pong {
	return &Pong{Source: source, Ts: time.Now().Unix()}
}

// Send returns Pong.
func (s *PingServer) Send(_ context.Context, m *Ping) (*Pong, error) {
	return s.pong([]*Ping{m}), nil
}

// SendStreamC returns Pong.
func (s *PingServer) SendStreamC(stream PingService_SendStreamCServer) error {
	source := make([]*Ping, 0)
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

// SendStreamS returns Pong.
func (s *PingServer) SendStreamS(m *Ping, stream PingService_SendStreamSServer) error {
	for i := 0; i < 10; i++ {
		err := stream.Send(s.pong([]*Ping{m}))
		if err != nil {
			return err
		}
	}
	return nil
}

// SendStreamB returns Pong.
func (s *PingServer) SendStreamB(stream PingService_SendStreamBServer) error {
	for {
		m, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		err = stream.Send(s.pong([]*Ping{m}))
		if err != nil {
			return err
		}
	}
}

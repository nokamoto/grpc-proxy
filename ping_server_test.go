package main

import (
	pb "github.com/nokamoto/grpc-proxy/examples/ping"
	"golang.org/x/net/context"
	"io"
)

type pingService struct{}

func (s *pingService) Send(_ context.Context, _ *pb.Ping) (*pb.Pong, error) {
	return &pb.Pong{}, nil
}

func (s *pingService) SendStreamC(stream pb.PingService_SendStreamCServer) error {
	for {
		_, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&pb.Pong{})
		}

		if err != nil {
			return err
		}
	}
}

func (s *pingService) SendStreamS(_ *pb.Ping, stream pb.PingService_SendStreamSServer) error {
	for i := 0; i < 10; i++ {
		err := stream.Send(&pb.Pong{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *pingService) SendStreamB(stream pb.PingService_SendStreamBServer) error {
	for {
		_, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		err = stream.Send(&pb.Pong{})
		if err != nil {
			return err
		}
	}
}

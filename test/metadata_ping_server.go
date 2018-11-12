package test

import (
	pb "github.com/nokamoto/grpc-proxy/examples/ping"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"io"
)

// MetadataPingServer implements pb.PingServiceServer only for test.
type MetadataPingServer struct {
	key     string
	header  func(string, []string) metadata.MD
	trailer func(string, []string) metadata.MD
}

// Send returns Metadata.
func (s *MetadataPingServer) Send(ctx context.Context, m *pb.Ping) (*pb.Pong, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, grpc.Errorf(codes.InvalidArgument, "metadata not found")
	}

	grpc.SetHeader(ctx, s.header(s.key, md.Get(s.key)))
	grpc.SetTrailer(ctx, s.trailer(s.key, md.Get(s.key)))

	return &pb.Pong{}, nil
}

// SendStreamC returns Metadata.
func (s *MetadataPingServer) SendStreamC(stream pb.PingService_SendStreamCServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return grpc.Errorf(codes.InvalidArgument, "metadata not found")
	}

	stream.SetHeader(s.header(s.key, md.Get(s.key)))
	stream.SetTrailer(s.trailer(s.key, md.Get(s.key)))

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

// SendStreamS returns Metadata.
func (s *MetadataPingServer) SendStreamS(m *pb.Ping, stream pb.PingService_SendStreamSServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return grpc.Errorf(codes.InvalidArgument, "metadata not found")
	}

	stream.SetHeader(s.header(s.key, md.Get(s.key)))
	stream.SetTrailer(s.trailer(s.key, md.Get(s.key)))

	for i := 0; i < 10; i++ {
		err := stream.Send(&pb.Pong{})
		if err != nil {
			return err
		}
	}
	return nil
}

// SendStreamB returns Metadata.
func (s *MetadataPingServer) SendStreamB(stream pb.PingService_SendStreamBServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return grpc.Errorf(codes.InvalidArgument, "metadata not found")
	}

	stream.SetHeader(s.header(s.key, md.Get(s.key)))
	stream.SetTrailer(s.trailer(s.key, md.Get(s.key)))

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

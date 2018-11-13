package test

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"io"
)

// MetadataPingServer implements PingServiceServer only for test.
type MetadataPingServer struct {
	key     string
	header  func(string, []string) metadata.MD
	trailer func(string, []string) metadata.MD
}

// Send returns Metadata.
func (s *MetadataPingServer) Send(ctx context.Context, m *Ping) (*Pong, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, grpc.Errorf(codes.InvalidArgument, "metadata not found")
	}

	grpc.SetHeader(ctx, s.header(s.key, md.Get(s.key)))
	grpc.SetTrailer(ctx, s.trailer(s.key, md.Get(s.key)))

	return &Pong{}, nil
}

// SendStreamC returns Metadata.
func (s *MetadataPingServer) SendStreamC(stream PingService_SendStreamCServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return grpc.Errorf(codes.InvalidArgument, "metadata not found")
	}

	stream.SetHeader(s.header(s.key, md.Get(s.key)))
	stream.SetTrailer(s.trailer(s.key, md.Get(s.key)))

	for {
		_, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&Pong{})
		}

		if err != nil {
			return err
		}
	}
}

// SendStreamS returns Metadata.
func (s *MetadataPingServer) SendStreamS(m *Ping, stream PingService_SendStreamSServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return grpc.Errorf(codes.InvalidArgument, "metadata not found")
	}

	stream.SetHeader(s.header(s.key, md.Get(s.key)))
	stream.SetTrailer(s.trailer(s.key, md.Get(s.key)))

	for i := 0; i < 10; i++ {
		err := stream.Send(&Pong{})
		if err != nil {
			return err
		}
	}
	return nil
}

// SendStreamB returns Metadata.
func (s *MetadataPingServer) SendStreamB(stream PingService_SendStreamBServer) error {
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

		err = stream.Send(&Pong{})
		if err != nil {
			return err
		}
	}
}

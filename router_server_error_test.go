package main

import (
	ping "github.com/nokamoto/grpc-proxy/examples/ping"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func Test_router_server_error_unary(t *testing.T) {
	withServer(t, &unimplementedPingService{}, func(ctx context.Context, cc *grpc.ClientConn) {
		c := ping.NewPingServiceClient(cc)

		_, err := c.Send(ctx, &ping.Ping{})
		s, _ := status.FromError(err)
		if s.Code() != codes.Unimplemented {
			t.Errorf("%v != %v %s", s.Code(), codes.Unimplemented, s.Message())
		}
	})
}

func Test_router_server_error_streamC(t *testing.T) {
	withServer(t, &unimplementedPingService{}, func(ctx context.Context, cc *grpc.ClientConn) {
		c := ping.NewPingServiceClient(cc)

		stream, err := c.SendStreamC(ctx)
		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 10; i++ {
			err = stream.Send(&ping.Ping{})
			if err != nil {
				t.Fatal(err)
			}
		}

		_, err = stream.CloseAndRecv()

		s, _ := status.FromError(err)
		if s.Code() != codes.Unimplemented {
			t.Errorf("%v != %v %s", s.Code(), codes.Unimplemented, s.Message())
		}
	})
}

func Test_router_server_error_streamS(t *testing.T) {
	withServer(t, &unimplementedPingService{}, func(ctx context.Context, cc *grpc.ClientConn) {
		c := ping.NewPingServiceClient(cc)

		stream, err := c.SendStreamS(ctx, &ping.Ping{})
		if err != nil {
			t.Fatal(err)
		}

		_, err = stream.Recv()

		s, _ := status.FromError(err)
		if s.Code() != codes.Unimplemented {
			t.Errorf("%v != %v %s", s.Code(), codes.Unimplemented, s.Message())
		}
	})
}

func Test_router_server_error_streamB(t *testing.T) {
	withServer(t, &unimplementedPingService{}, func(ctx context.Context, cc *grpc.ClientConn) {
		c := ping.NewPingServiceClient(cc)

		stream, err := c.SendStreamB(ctx)
		if err != nil {
			t.Fatal(err)
		}

		err = stream.Send(&ping.Ping{})
		if err != nil {
			t.Fatal(err)
		}

		_, err = stream.Recv()

		s, _ := status.FromError(err)
		if s.Code() != codes.Unimplemented {
			t.Errorf("%v != %v %s", s.Code(), codes.Unimplemented, s.Message())
		}
	})
}

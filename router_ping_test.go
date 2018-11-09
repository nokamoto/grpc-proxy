package main

import (
	ping "github.com/nokamoto/grpc-proxy/examples/ping"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"testing"
)

func Test_router_ping_unary(t *testing.T) {
	testWithPingServer(t, &pingService{}, func(ctx context.Context, cc *grpc.ClientConn) {
		c := ping.NewPingServiceClient(cc)

		_, err := c.Send(ctx, &ping.Ping{})
		s, _ := status.FromError(err)
		if s.Code() != codes.OK {
			t.Errorf("%v != %v %s", s.Code(), codes.OK, s.Message())
		}
	})
}

func Test_router_ping_streamC(t *testing.T) {
	testWithPingServer(t, &pingService{}, func(ctx context.Context, cc *grpc.ClientConn) {
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
		if s.Code() != codes.OK {
			t.Errorf("%v != %v %s", s.Code(), codes.OK, s.Message())
		}
	})
}

func Test_router_ping_streamS(t *testing.T) {
	testWithPingServer(t, &pingService{}, func(ctx context.Context, cc *grpc.ClientConn) {
		c := ping.NewPingServiceClient(cc)

		stream, err := c.SendStreamS(ctx, &ping.Ping{})
		if err != nil {
			t.Fatal(err)
		}

		i := 0
		for {
			_, err = stream.Recv()
			if err == io.EOF {
				break
			}

			i++

			s, _ := status.FromError(err)
			if s.Code() != codes.OK {
				t.Errorf("%v != %v %s", s.Code(), codes.OK, s.Message())
			}
		}

		if i != 10 {
			t.Errorf("%d != 10", i)
		}
	})
}

func Test_router_ping_streamB(t *testing.T) {
	testWithPingServer(t, &pingService{}, func(ctx context.Context, cc *grpc.ClientConn) {
		c := ping.NewPingServiceClient(cc)

		stream, err := c.SendStreamB(ctx)
		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 10; i++ {
			err = stream.Send(&ping.Ping{})
			if err != nil {
				t.Fatal(err)
			}

			_, err = stream.Recv()

			s, _ := status.FromError(err)
			if s.Code() != codes.OK {
				t.Errorf("%v != %v %s", s.Code(), codes.OK, s.Message())
			}
		}

		err = stream.CloseSend()
		if err != nil {
			t.Error(err)
		}
	})
}

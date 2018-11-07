package main

import (
	"fmt"
	ping "github.com/nokamoto/grpc-proxy/examples/ping"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
	"testing"
	"time"
)

func withServer(t *testing.T, f func(context.Context, *grpc.ClientConn)) {
	port := 9000

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		t.Fatal(err)
	}

	server := newGrpcServer()

	desc, err := newDescriptor("examples/ping/example.pb")
	if err != nil {
		t.Fatal(err)
	}

	for _, sd := range desc.serviceDescriptors() {
		server.RegisterService(sd, &grpcProxyServer{})
	}

	go func() {
		server.Serve(lis)
	}()
	defer server.GracefulStop()

	cc, err := grpc.Dial(fmt.Sprintf("%s:%d", "localhost", port), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer cc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	f(ctx, cc)
}

func Test_proxy_server_ping_unary(t *testing.T) {
	withServer(t, func(ctx context.Context, cc *grpc.ClientConn) {
		c := ping.NewPingServiceClient(cc)

		_, err := c.Send(ctx, &ping.Ping{})
		s, _ := status.FromError(err)
		if s.Code() != codes.Unimplemented {
			t.Errorf("%v != %v %s", s.Code(), codes.Unimplemented, s.Message())
		}
	})
}

func Test_proxy_server_ping_streamC(t *testing.T) {
	withServer(t, func(ctx context.Context, cc *grpc.ClientConn) {
		c := ping.NewPingServiceClient(cc)

		stream, err := c.SendStreamC(ctx)
		if err != nil {
			t.Fatal(err)
		}

		err = stream.Send(&ping.Ping{})
		if err != nil {
			t.Fatal(err)
		}

		_, err = stream.CloseAndRecv()

		s, _ := status.FromError(err)
		if s.Code() != codes.Unimplemented {
			t.Errorf("%v != %v %s", s.Code(), codes.Unimplemented, s.Message())
		}
	})
}

func Test_proxy_server_ping_streamS(t *testing.T) {
	withServer(t, func(ctx context.Context, cc *grpc.ClientConn) {
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

func Test_proxy_server_ping_streamB(t *testing.T) {
	withServer(t, func(ctx context.Context, cc *grpc.ClientConn) {
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

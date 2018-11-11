package test

import (
	"fmt"
	pb "github.com/nokamoto/grpc-proxy/examples/ping"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"net"
	"testing"
	"time"
)

func TestPingServer_Send(t *testing.T) {
	ctx, c, afterEach := beforeEachPing(t, 9000)
	defer afterEach()

	ping := pb.Ping{Ts: time.Now().Unix()}

	pong, err := c.Send(ctx, &ping)
	if err != nil {
		t.Fatal(err)
	}

	if pong.Source[0].Ts != ping.Ts {
		t.Errorf("%d != %d", pong.Source[0].Ts, ping.Ts)
	}
}

func TestPingServer_SendStreamC(t *testing.T) {
	ctx, c, afterEach := beforeEachPing(t, 9000)
	defer afterEach()

	stream, err := c.SendStreamC(ctx)
	if err != nil {
		t.Fatal(err)
	}

	source := make([]*pb.Ping, 0)
	for i := 0; i < 10; i++ {
		ping := &pb.Ping{Ts: time.Now().Unix()}

		err = stream.Send(ping)
		if err != nil {
			t.Fatal(err)
		}

		source = append(source, ping)
	}

	pong, err := stream.CloseAndRecv()
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		if pong.Source[i].Ts != source[i].Ts {
			t.Errorf("%d != %d", pong.Source[i].Ts, source[i].Ts)
		}
	}
}

func TestPingServer_SendStreamS(t *testing.T) {
	ctx, c, afterEach := beforeEachPing(t, 9000)
	defer afterEach()

	ping := &pb.Ping{Ts: time.Now().Unix()}

	stream, err := c.SendStreamS(ctx, ping)
	if err != nil {
		t.Fatal(err)
	}

	for {
		pong, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			t.Fatal(err)
		}

		if pong.Source[0].Ts != ping.Ts {
			t.Errorf("%d != %d", pong.Source[0].Ts, ping.Ts)
		}
	}
}

func TestPingServer_SendStreamB(t *testing.T) {
	ctx, c, afterEach := beforeEachPing(t, 9000)
	defer afterEach()

	stream, err := c.SendStreamB(ctx)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		ping := &pb.Ping{Ts: time.Now().Unix()}

		err := stream.Send(ping)
		if err != nil {
			t.Fatal(err)
		}

		pong, err := stream.Recv()
		if err != nil {
			t.Fatal(err)
		}

		if pong.Source[0].Ts != ping.Ts {
			t.Errorf("%d != %d", pong.Source[0].Ts, ping.Ts)
		}
	}

	err = stream.CloseSend()
	if err != nil {
		t.Fatal(err)
	}
}

func beforeEachPing(t *testing.T, port int) (context.Context, pb.PingServiceClient, func()) {
	afterEachServer := beforeEachPingServer(t, 9002)
	afterEachGrpcProxy := beforeEachGrpcProxy(t, port, "../examples/ping/example.pb", "../examples/ping/example.yaml")

	cc, err := grpc.Dial(fmt.Sprintf("%s:%d", "localhost", port), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	return ctx, pb.NewPingServiceClient(cc), func() {
		cc.Close()
		cancel()

		afterEachGrpcProxy()

		afterEachServer()
	}
}

func beforeEachPingServer(t *testing.T, port int) func() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		t.Fatal(err)
	}

	opts := []grpc.ServerOption{}
	srv := grpc.NewServer(opts...)

	pb.RegisterPingServiceServer(srv, &PingServer{})

	go func() {
		srv.Serve(lis)
	}()

	return func() {
		srv.GracefulStop()
	}
}

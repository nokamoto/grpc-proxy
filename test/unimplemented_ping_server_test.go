package test

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
	"testing"
	"time"
)

func TestUnimplementedPingServer_Send(t *testing.T) {
	ctx, c, afterEach := beforeEachUnimplementedPing(t, 9000)
	defer afterEach()

	ping := Ping{Ts: time.Now().Unix()}

	_, err := c.Send(ctx, &ping)
	if code := status.Code(err); code != codes.Unimplemented {
		t.Errorf("%v != %v", code, codes.Unimplemented)
	}
}

func TestUnimplementedPingServer_SendStreamC(t *testing.T) {
	ctx, c, afterEach := beforeEachUnimplementedPing(t, 9000)
	defer afterEach()

	stream, err := c.SendStreamC(ctx)
	if err != nil {
		t.Fatal(err)
	}

	source := make([]*Ping, 0)
	for i := 0; i < 10; i++ {
		ping := &Ping{Ts: time.Now().Unix()}

		err = stream.Send(ping)
		if err != nil {
			t.Fatal(err)
		}

		source = append(source, ping)
	}

	_, err = stream.CloseAndRecv()
	if code := status.Code(err); code != codes.Unimplemented {
		t.Errorf("%v != %v", code, codes.Unimplemented)
	}
}

func TestUnimplementedPingServer_SendStreamS(t *testing.T) {
	ctx, c, afterEach := beforeEachUnimplementedPing(t, 9000)
	defer afterEach()

	ping := &Ping{Ts: time.Now().Unix()}

	stream, err := c.SendStreamS(ctx, ping)
	if err != nil {
		t.Fatal(err)
	}

	_, err = stream.Recv()
	if code := status.Code(err); code != codes.Unimplemented {
		t.Errorf("%v != %v", code, codes.Unimplemented)
	}
}

func TestUnimplementedPingServer_SendStreamB(t *testing.T) {
	ctx, c, afterEach := beforeEachUnimplementedPing(t, 9000)
	defer afterEach()

	stream, err := c.SendStreamB(ctx)
	if err != nil {
		t.Fatal(err)
	}

	ping := &Ping{Ts: time.Now().Unix()}

	err = stream.Send(ping)
	if err != nil {
		t.Fatal(err)
	}

	_, err = stream.Recv()
	if code := status.Code(err); code != codes.Unimplemented {
		t.Errorf("%v != %v", code, codes.Unimplemented)
	}
}

func beforeEachUnimplementedPing(t *testing.T, port int) (context.Context, PingServiceClient, func()) {
	afterEachServer := beforeEachUnimplementedPingServer(t, 9002)
	afterEachGrpcProxy := beforeEachGrpcProxy(t, port, "../testdata/protobuf/ping/ping.pb", "../testdata/yaml/ping.yaml")

	cc, err := grpc.Dial(fmt.Sprintf("%s:%d", "localhost", port), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	return ctx, NewPingServiceClient(cc), func() {
		cc.Close()
		cancel()

		afterEachGrpcProxy()

		afterEachServer()
	}
}

func beforeEachUnimplementedPingServer(t *testing.T, port int) func() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		t.Fatal(err)
	}

	opts := []grpc.ServerOption{}
	srv := grpc.NewServer(opts...)

	RegisterPingServiceServer(srv, &UnimplementedPingServer{})

	go func() {
		srv.Serve(lis)
	}()

	return func() {
		srv.GracefulStop()
	}
}

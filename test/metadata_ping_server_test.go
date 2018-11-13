package test

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
	"net"
	"strings"
	"testing"
	"time"
)

var (
	metadataKey = "metadatakey"

	doubleMetadataKey = strings.Repeat(metadataKey, 2)

	tripleMetadataKey = strings.Repeat(metadataKey, 3)

	doubleHeader = func(key string, values []string) metadata.MD {
		if len(values) == 0 {
			return metadata.Pairs()
		}
		return metadata.Pairs(strings.Repeat(key, 2), strings.Repeat(values[0], 2))
	}

	tripleTrailer = func(key string, values []string) metadata.MD {
		if len(values) == 0 {
			return metadata.Pairs()
		}
		return metadata.Pairs(strings.Repeat(key, 3), strings.Repeat(values[0], 3))
	}
)

func TestMetadataPingServer_Send(t *testing.T) {
	ctx, c, afterEach := beforeEachMetadataPing(t, 9000)
	defer afterEach()

	ctx = metadata.AppendToOutgoingContext(ctx, metadataKey, "v")

	var header, trailer metadata.MD

	_, err := c.Send(ctx, &Ping{}, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		t.Fatal(err)
	}

	testMetadata(t, header, doubleMetadataKey, "vv")
	testMetadata(t, trailer, tripleMetadataKey, "vvv")
}

func TestMetadataPingServer_SendStreamC(t *testing.T) {
	ctx, c, afterEach := beforeEachMetadataPing(t, 9000)
	defer afterEach()

	ctx = metadata.AppendToOutgoingContext(ctx, metadataKey, "v")

	var header, trailer metadata.MD

	stream, err := c.SendStreamC(ctx, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		t.Fatal(err)
	}

	err = stream.Send(&Ping{})
	if err != nil {
		t.Fatal(err)
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		t.Fatal(err)
	}

	testMetadata(t, header, doubleMetadataKey, "vv")
	testMetadata(t, trailer, tripleMetadataKey, "vvv")
}

func TestMetadataPingServer_SendStreamS(t *testing.T) {
	ctx, c, afterEach := beforeEachMetadataPing(t, 9000)
	defer afterEach()

	ctx = metadata.AppendToOutgoingContext(ctx, metadataKey, "v")

	var header, trailer metadata.MD

	stream, err := c.SendStreamS(ctx, &Ping{}, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		t.Fatal(err)
	}

	for {
		_, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			t.Fatal(err)
		}
	}

	testMetadata(t, header, doubleMetadataKey, "vv")
	testMetadata(t, trailer, tripleMetadataKey, "vvv")
}

func TestMetadataPingServer_SendStreamB(t *testing.T) {
	ctx, c, afterEach := beforeEachMetadataPing(t, 9000)
	defer afterEach()

	ctx = metadata.AppendToOutgoingContext(ctx, metadataKey, "v")

	var header, trailer metadata.MD

	stream, err := c.SendStreamB(ctx, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		err := stream.Send(&Ping{})
		if err != nil {
			t.Fatal(err)
		}

		_, err = stream.Recv()
		if err != nil {
			t.Fatal(err)
		}
	}

	err = stream.CloseSend()
	if err != nil {
		t.Fatal(err)
	}

	t.Skip("bidirectional stream metadata proxy impl.")
}

func testMetadata(t *testing.T, md metadata.MD, k, v string) {
	t.Helper()
	if len(md.Get(k)) != 1 {
		t.Fatalf("len(%v) != 1", md)
	}
	if s := md.Get(k)[0]; s != v {
		t.Errorf("%s != %s", s, v)
	}
}

func beforeEachMetadataPing(t *testing.T, port int) (context.Context, PingServiceClient, func()) {
	afterEachServer := beforeEachMetadataPingServer(t, 9002)
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

func beforeEachMetadataPingServer(t *testing.T, port int) func() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		t.Fatal(err)
	}

	opts := []grpc.ServerOption{}
	srv := grpc.NewServer(opts...)

	RegisterPingServiceServer(srv, &MetadataPingServer{key: metadataKey, header: doubleHeader, trailer: tripleTrailer})

	go func() {
		srv.Serve(lis)
	}()

	return func() {
		srv.GracefulStop()
	}
}

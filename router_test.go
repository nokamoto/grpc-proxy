package main

import (
	"fmt"
	"github.com/nokamoto/grpc-proxy/descriptor"
	empty "github.com/nokamoto/grpc-proxy/examples/empty-package"
	ping "github.com/nokamoto/grpc-proxy/examples/ping"
	"github.com/nokamoto/grpc-proxy/yaml"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
	"testing"
	"time"
)

func withPingServer(svc ping.PingServiceServer, f func()) error {
	port := 9002

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	opts := []grpc.ServerOption{}
	srv := grpc.NewServer(opts...)

	ping.RegisterPingServiceServer(srv, svc)

	go func() {
		srv.Serve(lis)
	}()
	defer srv.GracefulStop()

	f()

	return nil
}

func withEmptyServer(svc empty.ServiceServer, f func()) error {
	port := 9001

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	opts := []grpc.ServerOption{}
	srv := grpc.NewServer(opts...)

	empty.RegisterServiceServer(srv, svc)

	go func() {
		srv.Serve(lis)
	}()
	defer srv.GracefulStop()

	f()

	return nil
}

func withProxyServer(t *testing.T, pb string, yml string, f func(context.Context, *grpc.ClientConn)) func() {
	return func() {
		port := 9000

		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			t.Fatal(err)
		}

		server := newGrpcServer()

		desc, err := descriptor.NewDescriptor(pb)
		if err != nil {
			t.Fatal(err)
		}

		routes, clusters, err := yaml.NewYaml(yml)
		if err != nil {
			panic(err)
		}

		router, err := newRouter(desc.FileDescriptorSet, routes, clusters)
		if err != nil {
			panic(err)
		}

		for _, sd := range desc.ServiceDescs() {
			server.RegisterService(sd, router)
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
}

func testWithPingServer(t *testing.T, svc ping.PingServiceServer, f func(context.Context, *grpc.ClientConn)) {
	err := withPingServer(svc, withProxyServer(t, "examples/ping/example.pb", "examples/ping/example.yaml", f))

	if err != nil {
		t.Error(err)
	}
}

func testWithEmptyServer(t *testing.T, svc empty.ServiceServer, f func(context.Context, *grpc.ClientConn)) {
	err := withEmptyServer(svc, withProxyServer(t, "examples/empty-package/example.pb", "examples/empty-package/example.yaml", f))

	if err != nil {
		t.Error(err)
	}
}

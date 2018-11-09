package main

import (
	"fmt"
	ping "github.com/nokamoto/grpc-proxy/examples/ping"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
	"testing"
	"time"
)

func withPingServer(svc ping.PingServiceServer, f func() error) error {
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

	return f()
}

func withServer(t *testing.T, svc ping.PingServiceServer, f func(context.Context, *grpc.ClientConn)) {
	err := withPingServer(svc, func() error {
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

		routes, clusters, err := newYaml("examples/ping/example.yaml")
		if err != nil {
			panic(err)
		}

		router, err := newRouter(desc.FileDescriptorSet, routes, clusters)
		if err != nil {
			panic(err)
		}

		for _, sd := range desc.serviceDescriptors() {
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
		return nil
	})

	if err != nil {
		t.Error(err)
	}
}

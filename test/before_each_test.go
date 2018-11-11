package test

import (
	"fmt"
	"github.com/nokamoto/grpc-proxy/codec"
	"github.com/nokamoto/grpc-proxy/descriptor"
	"github.com/nokamoto/grpc-proxy/route"
	"github.com/nokamoto/grpc-proxy/yaml"
	"google.golang.org/grpc"
	"net"
	"testing"
)

func beforeEachGrpcProxy(t *testing.T, port int, pb, yml string) func() {
	opts := []grpc.ServerOption{grpc.CustomCodec(codec.RawCodec{})}
	srv := grpc.NewServer(opts...)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		t.Fatal(err)
	}

	desc, err := descriptor.NewDescriptor(pb)
	if err != nil {
		t.Fatal(err)
	}

	routes, clusters, err := yaml.NewYaml(yml)
	if err != nil {
		panic(err)
	}

	router, err := route.NewRoutes(desc, routes, clusters)
	if err != nil {
		panic(err)
	}

	for _, sd := range descriptor.ServiceDescs(desc) {
		srv.RegisterService(sd, router)
	}

	go func() {
		srv.Serve(lis)
	}()

	return func() {
		srv.GracefulStop()
	}
}

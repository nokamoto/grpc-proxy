package main

import (
	"flag"
	"fmt"
	"github.com/nokamoto/grpc-proxy/codec"
	"github.com/nokamoto/grpc-proxy/descriptor"
	"github.com/nokamoto/grpc-proxy/route"
	"github.com/nokamoto/grpc-proxy/yaml"
	"google.golang.org/grpc"
	"net"
)

func newGrpcServer() *grpc.Server {
	opts := []grpc.ServerOption{grpc.CustomCodec(codec.RawCodec{})}
	return grpc.NewServer(opts...)
}

func main() {
	var (
		port = flag.Int("p", 9000, "gRPC server port")
		pb   = flag.String("pb", "", "file descriptor protocol buffers filepath")
		y    = flag.String("yaml", "", "yaml configuration filepath")
	)

	flag.Parse()

	routes, clusters, err := yaml.NewYaml(*y)
	if err != nil {
		panic(err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		panic(err)
	}

	server := newGrpcServer()

	desc, err := descriptor.NewDescriptor(*pb)
	if err != nil {
		panic(err)
	}

	router, err := route.NewRoutes(desc, routes, clusters)
	if err != nil {
		panic(err)
	}

	for _, sd := range descriptor.ServiceDescs(desc) {
		server.RegisterService(sd, router)
	}

	server.Serve(lis)
}

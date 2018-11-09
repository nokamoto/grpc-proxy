package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

func newGrpcServer() *grpc.Server {
	opts := []grpc.ServerOption{grpc.CustomCodec(codec{})}
	return grpc.NewServer(opts...)
}

func main() {
	var (
		port = flag.Int("p", 9000, "gRPC server port")
		pb   = flag.String("pb", "", "file descriptor protocol buffers filepath")
		y    = flag.String("yaml", "", "yaml configuration filepath")
	)

	flag.Parse()

	routes, clusters, err := newYaml(*y)
	if err != nil {
		panic(err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		panic(err)
	}

	server := newGrpcServer()

	desc, err := newDescriptor(*pb)
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

	server.Serve(lis)
}

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
	)

	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		panic(err)
	}

	server := newGrpcServer()

	desc, err := newDescriptor(*pb)
	if err != nil {
		panic(err)
	}

	for _, sd := range desc.serviceDescriptors() {
		server.RegisterService(sd, &grpcProxyServer{})
	}

	server.Serve(lis)
}

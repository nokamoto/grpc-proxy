package main

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type cluster interface {
	invokeUnary(context.Context, *message, string) (*message, error)
	invokeStreamC(proxyStreamCServer, *grpc.StreamDesc, string) error
	invokeStreamS(proxyStreamSServer, *grpc.StreamDesc, string) error
}

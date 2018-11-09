package main

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type proxyServer interface {
	unary(context.Context, *message, string) (*message, error)
	streamC(proxyStreamCServer) error
	streamS(proxyStreamSServer) error
	streamB(grpc.ServerStream) error
}

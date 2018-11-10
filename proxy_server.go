package main

import (
	"github.com/nokamoto/grpc-proxy/codec"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type proxyServer interface {
	unary(context.Context, *codec.RawMessage, string) (*codec.RawMessage, error)
	streamC(proxyStreamCServer, *grpc.StreamDesc, string) error
	streamS(proxyStreamSServer, *grpc.StreamDesc, string) error
	streamB(proxyStreamBServer, *grpc.StreamDesc, string) error
}

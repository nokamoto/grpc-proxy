package main

import (
	"github.com/nokamoto/grpc-proxy/codec"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type cluster interface {
	invokeUnary(context.Context, *codec.RawMessage, string) (*codec.RawMessage, error)
	invokeStreamC(proxyStreamCServer, *grpc.StreamDesc, string) error
	invokeStreamS(proxyStreamSServer, *grpc.StreamDesc, string) error
	invokeStreamB(proxyStreamBServer, *grpc.StreamDesc, string) error
}

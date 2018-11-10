package main

import (
	"github.com/nokamoto/grpc-proxy/codec"
	"github.com/nokamoto/grpc-proxy/server"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type cluster interface {
	invokeUnary(context.Context, *codec.RawMessage, string) (*codec.RawMessage, error)
	invokeStreamC(server.RawServerStreamC, *grpc.StreamDesc, string) error
	invokeStreamS(server.RawServerStreamS, *grpc.StreamDesc, string) error
	invokeStreamB(server.RawServerStreamB, *grpc.StreamDesc, string) error
}

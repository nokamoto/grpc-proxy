package server

import (
	"github.com/nokamoto/grpc-proxy/codec"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Server represents a gRPC server which sends/receives codec.RawMessage.
type Server interface {
	Unary(context.Context, *codec.RawMessage, string) (*codec.RawMessage, error)
	StreamC(RawServerStreamC, *grpc.StreamDesc, string) error
	StreamS(RawServerStreamS, *grpc.StreamDesc, string) error
	StreamB(RawServerStreamB, *grpc.StreamDesc, string) error
}

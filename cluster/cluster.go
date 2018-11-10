package cluster

import (
	"github.com/nokamoto/grpc-proxy/codec"
	"github.com/nokamoto/grpc-proxy/server"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Cluster interface {
	InvokeUnary(context.Context, *codec.RawMessage, string) (*codec.RawMessage, error)
	InvokeStreamC(server.RawServerStreamC, *grpc.StreamDesc, string) error
	InvokeStreamS(server.RawServerStreamS, *grpc.StreamDesc, string) error
	InvokeStreamB(server.RawServerStreamB, *grpc.StreamDesc, string) error
}

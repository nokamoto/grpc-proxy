package route

import (
	"github.com/nokamoto/grpc-proxy/cluster"
	"github.com/nokamoto/grpc-proxy/codec"
	"github.com/nokamoto/grpc-proxy/observe"
	"github.com/nokamoto/grpc-proxy/server"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type route struct {
	cluster cluster.Cluster
	log     observe.Log
}

func (r *route) unary(ctx context.Context, m *codec.RawMessage, method string) (*codec.RawMessage, error) {
	return r.cluster.InvokeUnary(ctx, m, method)
}

func (r *route) streamC(stream server.RawServerStreamC, desc *grpc.StreamDesc, method string) error {
	return r.cluster.InvokeStreamC(stream, desc, method)
}

func (r *route) streamS(stream server.RawServerStreamS, desc *grpc.StreamDesc, method string) error {
	return r.cluster.InvokeStreamS(stream, desc, method)
}

func (r *route) streamB(stream server.RawServerStreamB, desc *grpc.StreamDesc, method string) error {
	return r.cluster.InvokeStreamB(stream, desc, method)
}

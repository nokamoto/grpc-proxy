package route

import (
	"fmt"
	"github.com/nokamoto/grpc-proxy/cluster"
	"github.com/nokamoto/grpc-proxy/codec"
	"github.com/nokamoto/grpc-proxy/observe"
	"github.com/nokamoto/grpc-proxy/server"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"os"
	"time"
)

type route struct {
	cluster cluster.Cluster
	log     observe.Log
}

func (r *route) unary(ctx context.Context, m *codec.RawMessage, method string) (*codec.RawMessage, error) {
	start := time.Now()

	res, err := r.cluster.InvokeUnary(ctx, m, method)

	d := time.Since(start)

	_, e := r.log.Write(method, status.Code(err), m.Size(), res.Size(), d)
	if e != nil {
		fmt.Fprintf(os.Stderr, "%s: unary access log error: %s\n", method, e)
	}

	return res, err
}

func (r *route) streamC(stream server.RawServerStreamC, desc *grpc.StreamDesc, method string) error {
	start := time.Now()

	err := r.cluster.InvokeStreamC(stream, desc, method)

	d := time.Since(start)

	_, e := r.log.Write(method, status.Code(err), -1, -1, d)
	if e != nil {
		fmt.Fprintf(os.Stderr, "%s: stream c access log error: %s\n", method, e)
	}

	return err
}

func (r *route) streamS(stream server.RawServerStreamS, desc *grpc.StreamDesc, method string) error {
	start := time.Now()

	err := r.cluster.InvokeStreamS(stream, desc, method)

	d := time.Since(start)

	_, e := r.log.Write(method, status.Code(err), -1, -1, d)
	if e != nil {
		fmt.Fprintf(os.Stderr, "%s: stream s access log error: %s\n", method, e)
	}

	return err
}

func (r *route) streamB(stream server.RawServerStreamB, desc *grpc.StreamDesc, method string) error {
	start := time.Now()

	err := r.cluster.InvokeStreamB(stream, desc, method)

	d := time.Since(start)

	_, e := r.log.Write(method, status.Code(err), -1, -1, d)
	if e != nil {
		fmt.Fprintf(os.Stderr, "%s: stream b access log error: %s\n", method, e)
	}

	return err
}

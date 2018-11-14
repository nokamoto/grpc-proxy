package route

import (
	"fmt"
	"github.com/nokamoto/grpc-proxy/cluster"
	"github.com/nokamoto/grpc-proxy/codec"
	"github.com/nokamoto/grpc-proxy/observe"
	"github.com/nokamoto/grpc-proxy/server"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
	"time"
)

type route struct {
	cluster cluster.Cluster
	log     *observe.Log
	prom    *observe.Prom
}

func (r *route) observe(method string, code codes.Code, req, res int, d time.Duration) {
	if r.log != nil {
		_, e := (*r.log).Write(method, code, req, res, d)
		if e != nil {
			fmt.Fprintf(os.Stderr, "%s: access log error: %s\n", method, e)
		}
	}
	if r.prom != nil {
		e := (*r.prom).Observe(method, code, req, res, d)
		if e != nil {
			fmt.Fprintf(os.Stderr, "%s: unary prom: %s\n", method, e)
		}
	}
}

func (r *route) unary(ctx context.Context, m *codec.RawMessage, method string) (*codec.RawMessage, error) {
	start := time.Now()

	res, err := r.cluster.InvokeUnary(ctx, m, method)

	d := time.Since(start)

	code := status.Code(err)

	r.observe(method, code, m.Size(), res.Size(), d)

	return res, err
}

func (r *route) streamC(stream server.RawServerStreamC, desc *grpc.StreamDesc, method string) error {
	start := time.Now()

	err := r.cluster.InvokeStreamC(stream, desc, method)

	d := time.Since(start)

	code := status.Code(err)

	r.observe(method, code, -1, -1, d)

	return err
}

func (r *route) streamS(stream server.RawServerStreamS, desc *grpc.StreamDesc, method string) error {
	start := time.Now()

	err := r.cluster.InvokeStreamS(stream, desc, method)

	d := time.Since(start)

	code := status.Code(err)

	r.observe(method, code, -1, -1, d)

	return err
}

func (r *route) streamB(stream server.RawServerStreamB, desc *grpc.StreamDesc, method string) error {
	start := time.Now()

	err := r.cluster.InvokeStreamB(stream, desc, method)

	d := time.Since(start)

	code := status.Code(err)

	r.observe(method, code, -1, -1, d)

	return err
}

func (r *route) destroy() {
	if r.prom != nil {
		(*r.prom).Destroy()
	}
}

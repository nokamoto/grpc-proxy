package route

import (
	"fmt"
	pb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/nokamoto/grpc-proxy/cluster"
	"github.com/nokamoto/grpc-proxy/codec"
	"github.com/nokamoto/grpc-proxy/descriptor"
	obs "github.com/nokamoto/grpc-proxy/observe"
	"github.com/nokamoto/grpc-proxy/server"
	"github.com/nokamoto/grpc-proxy/yaml"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// Routes implements server.Server.
type Routes struct {
	routes map[string]*route
}

// NewRoutes returns Routes from the yaml configurations.
func NewRoutes(fds *pb.FileDescriptorSet, yml *yaml.Yaml) (*Routes, error) {
	r := &Routes{
		routes: make(map[string]*route),
	}

	cs := make(map[string]cluster.Cluster)

	for _, yc := range yml.Clusters {
		c, err := cluster.NewRoundRobin(yc)
		if err != nil {
			return nil, err
		}

		cs[yc.Name] = c
	}

	ls := make(map[string]obs.Log)

	for _, yl := range yml.Observe.Logs {
		l, err := obs.NewLog(yl)
		if err != nil {
			return nil, err
		}

		ls[yl.Name] = l
	}

	ps := make(map[string]obs.Prom)

	for _, yp := range yml.Observe.Prom {
		p, err := obs.NewProm(yp)
		if err != nil {
			return nil, err
		}

		ps[yp.Name] = p
	}

	for _, fd := range fds.File {
		for _, sd := range fd.GetService() {
			for _, md := range sd.GetMethod() {
				full := descriptor.FullMethod(fd, sd, md)
				yr := yml.FindByFullMethod(full)

				if len(yr) == 0 {
					return nil, fmt.Errorf("%s has no route", full)
				} else if len(yr) > 1 {
					return nil, fmt.Errorf("%s has ambiguous routes: %v", full, yr)
				}

				head := yr[0]

				cluster, ok := cs[head.Cluster.Name]
				if !ok {
					return nil, fmt.Errorf("cluster %s is undefined", head.Cluster.Name)
				}

				log, ok := ls[head.Observe.Log.Name]
				if !ok {
					return nil, fmt.Errorf("log %s is undefined", head.Observe.Log.Name)
				}

				prom, ok := ps[head.Observe.Prom.Name]
				if !ok {
					return nil, fmt.Errorf("prom %s is undefined", head.Observe.Prom.Name)
				}

				r.routes[full] = &route{cluster: cluster, log: log, prom: prom}
			}
		}
	}

	return r, nil
}

// Unary routes codec.RawMessage to a selected cluster.
func (r *Routes) Unary(ctx context.Context, m *codec.RawMessage, method string) (*codec.RawMessage, error) {
	c, ok := r.routes[method]
	if !ok {
		return nil, grpc.Errorf(codes.Unknown, "[grpc-proxy] unknown")
	}
	return c.unary(ctx, m, method)
}

// StreamC routes the client side stream to a selected cluster.
func (r *Routes) StreamC(stream server.RawServerStreamC, desc *grpc.StreamDesc, method string) error {
	c, ok := r.routes[method]
	if !ok {
		return grpc.Errorf(codes.Unknown, "[grpc-proxy] unknown")
	}
	return c.streamC(stream, desc, method)
}

// StreamS routes the server side stream to a selected cluster.
func (r *Routes) StreamS(stream server.RawServerStreamS, desc *grpc.StreamDesc, method string) error {
	c, ok := r.routes[method]
	if !ok {
		return grpc.Errorf(codes.Unknown, "[grpc-proxy] unknown")
	}
	return c.streamS(stream, desc, method)
}

// StreamB routes the bidirectional stream to a selected cluster.
func (r *Routes) StreamB(stream server.RawServerStreamB, desc *grpc.StreamDesc, method string) error {
	c, ok := r.routes[method]
	if !ok {
		return grpc.Errorf(codes.Unknown, "[grpc-proxy] unknown")
	}
	return c.streamB(stream, desc, method)
}

func (r *Routes) Destroy() {
	for _, r := range r.routes {
		r.destroy()
	}
}

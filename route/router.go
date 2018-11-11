package route

import (
	"fmt"
	pb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/nokamoto/grpc-proxy/cluster"
	"github.com/nokamoto/grpc-proxy/codec"
	"github.com/nokamoto/grpc-proxy/descriptor"
	"github.com/nokamoto/grpc-proxy/server"
	"github.com/nokamoto/grpc-proxy/yaml"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// Routes implements server.Server.
type Routes struct {
	clusters map[string]cluster.Cluster
}

// NewRoutes returns Routes from the yaml configurations.
func NewRoutes(fds *pb.FileDescriptorSet, routes *yaml.Routes, clusters *yaml.Clusters) (*Routes, error) {
	r := &Routes{
		clusters: make(map[string]cluster.Cluster),
	}

	cs := make(map[string]cluster.Cluster)

	for _, yc := range clusters.Clusters {
		c, err := cluster.NewRoundRobin(yc)
		if err != nil {
			return nil, err
		}

		cs[yc.Name] = c
	}

	for _, fd := range fds.File {
		for _, sd := range fd.GetService() {
			for _, md := range sd.GetMethod() {
				full := descriptor.FullMethod(fd, sd, md)
				route := routes.FindByFullMethod(full)

				if len(route) == 0 {
					return nil, fmt.Errorf("%s has no route", full)
				} else if len(route) > 1 {
					return nil, fmt.Errorf("%s has ambiguous routes: %v", full, route)
				}

				cluster, ok := cs[route[0].Cluster.Name]
				if !ok {
					return nil, fmt.Errorf("cluster %s is undefined", route[0].Cluster.Name)
				}

				r.clusters[full] = cluster
			}
		}
	}

	return r, nil
}

// Unary routes codec.RawMessage to a selected cluster.
func (r *Routes) Unary(ctx context.Context, m *codec.RawMessage, method string) (*codec.RawMessage, error) {
	c, ok := r.clusters[method]
	if !ok {
		return nil, grpc.Errorf(codes.Unknown, "[grpc-proxy] unknown")
	}
	return c.InvokeUnary(ctx, m, method)
}

// StreamC routes the client side stream to a selected cluster.
func (r *Routes) StreamC(stream server.RawServerStreamC, desc *grpc.StreamDesc, method string) error {
	c, ok := r.clusters[method]
	if !ok {
		return grpc.Errorf(codes.Unknown, "[grpc-proxy] unknown")
	}
	return c.InvokeStreamC(stream, desc, method)
}

// StreamS routes the server side stream to a selected cluster.
func (r *Routes) StreamS(stream server.RawServerStreamS, desc *grpc.StreamDesc, method string) error {
	c, ok := r.clusters[method]
	if !ok {
		return grpc.Errorf(codes.Unknown, "[grpc-proxy] unknown")
	}
	return c.InvokeStreamS(stream, desc, method)
}

// StreamB routes the bidirectional stream to a selected cluster.
func (r *Routes) StreamB(stream server.RawServerStreamB, desc *grpc.StreamDesc, method string) error {
	c, ok := r.clusters[method]
	if !ok {
		return grpc.Errorf(codes.Unknown, "[grpc-proxy] unknown")
	}
	return c.InvokeStreamB(stream, desc, method)
}

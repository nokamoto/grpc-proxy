package main

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

type router struct {
	clusters map[string]cluster.Cluster
}

func newRouter(fds *pb.FileDescriptorSet, routes *yaml.Routes, clusters *yaml.Clusters) (*router, error) {
	r := &router{
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

func (r *router) Unary(ctx context.Context, m *codec.RawMessage, method string) (*codec.RawMessage, error) {
	c, ok := r.clusters[method]
	if !ok {
		return nil, grpc.Errorf(codes.Unknown, "[grpc-proxy] unknown")
	}
	return c.InvokeUnary(ctx, m, method)
}

func (r *router) StreamC(stream server.RawServerStreamC, desc *grpc.StreamDesc, method string) error {
	c, ok := r.clusters[method]
	if !ok {
		return grpc.Errorf(codes.Unknown, "[grpc-proxy] unknown")
	}
	return c.InvokeStreamC(stream, desc, method)
}

func (r *router) StreamS(stream server.RawServerStreamS, desc *grpc.StreamDesc, method string) error {
	c, ok := r.clusters[method]
	if !ok {
		return grpc.Errorf(codes.Unknown, "[grpc-proxy] unknown")
	}
	return c.InvokeStreamS(stream, desc, method)
}

func (r *router) StreamB(stream server.RawServerStreamB, desc *grpc.StreamDesc, method string) error {
	c, ok := r.clusters[method]
	if !ok {
		return grpc.Errorf(codes.Unknown, "[grpc-proxy] unknown")
	}
	return c.InvokeStreamB(stream, desc, method)
}

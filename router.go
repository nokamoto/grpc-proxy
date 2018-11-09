package main

import (
	"fmt"
	pb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type router struct {
	clusters map[string]cluster
}

func newRouter(fds *pb.FileDescriptorSet, routes *yamlRoutes, clusters *yamlClusters) (*router, error) {
	r := &router{
		clusters: make(map[string]cluster),
	}

	cs := make(map[string]cluster)

	for _, cluster := range clusters.Clusters {
		c, err := newClusterRoundRobin(cluster)
		if err != nil {
			return nil, err
		}

		cs[cluster.Name] = c
	}

	for _, fd := range fds.File {
		for _, sd := range fd.GetService() {
			for _, md := range sd.GetMethod() {
				full := fullMethod(fd, sd, md)
				route := routes.findByFullMethod(full)

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

func (r *router) unary(ctx context.Context, m *message, method string) (*message, error) {
	cluster, ok := r.clusters[method]
	if !ok {
		return nil, grpc.Errorf(codes.Unknown, "[grpc-proxy] unknown")
	}
	return cluster.invokeUnary(ctx, m, method)
}

func (r *router) streamC(stream proxyStreamCServer, desc *grpc.StreamDesc, method string) error {
	cluster, ok := r.clusters[method]
	if !ok {
		return grpc.Errorf(codes.Unknown, "[grpc-proxy] unknown")
	}
	return cluster.invokeStreamC(stream, desc, method)
}

func (r *router) streamS(stream proxyStreamSServer, desc *grpc.StreamDesc, method string) error {
	cluster, ok := r.clusters[method]
	if !ok {
		return grpc.Errorf(codes.Unknown, "[grpc-proxy] unknown")
	}
	return cluster.invokeStreamS(stream, desc, method)
}

func (r *router) streamB(stream proxyStreamBServer, desc *grpc.StreamDesc, method string) error {
	cluster, ok := r.clusters[method]
	if !ok {
		return grpc.Errorf(codes.Unknown, "[grpc-proxy] unknown")
	}
	return cluster.invokeStreamB(stream, desc, method)
}

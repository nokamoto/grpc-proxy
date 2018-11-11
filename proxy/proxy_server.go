package proxy

import (
	"fmt"
	"github.com/nokamoto/grpc-proxy/codec"
	"github.com/nokamoto/grpc-proxy/descriptor"
	"github.com/nokamoto/grpc-proxy/route"
	"github.com/nokamoto/grpc-proxy/yaml"
	"google.golang.org/grpc"
	"net"
)

// Server represents a gRPC server combined all routes.
type Server struct {
	srv *grpc.Server
	lis net.Listener
}

// NewServer returns a gRPC server from the gRPC proxy server port, the file descriptor protocol buffers filepath,
// and the yaml configuration filepath.NewProxyServer
func NewServer(port int, pb, yml string) (*Server, error) {
	routes, clusters, err := yaml.NewYaml(yml)
	if err != nil {
		return nil, err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	opts := []grpc.ServerOption{grpc.CustomCodec(codec.RawCodec{})}

	srv := grpc.NewServer(opts...)

	desc, err := descriptor.NewDescriptor(pb)
	if err != nil {
		return nil, err
	}

	router, err := route.NewRoutes(desc, routes, clusters)
	if err != nil {
		return nil, err
	}

	for _, sd := range descriptor.ServiceDescs(desc) {
		srv.RegisterService(sd, router)
	}

	return &Server{srv: srv, lis: lis}, nil
}

// Serve starts the gRPC proxy server.
func (s *Server) Serve() error {
	return s.srv.Serve(s.lis)
}

// GracefulStop gracefully stops the gRPC proxy server.
func (s *Server) GracefulStop() {
	s.srv.GracefulStop()
}

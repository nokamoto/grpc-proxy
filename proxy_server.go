package main

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type proxyServer interface {
	unary(context.Context, *message) (*message, error)
	streamC(proxyStreamCServer) error
	streamS(proxyStreamSServer) error
	streamB(grpc.ServerStream) error
}

type grpcProxyServer struct{}

func (s *grpcProxyServer) unary(context.Context, *message) (*message, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "[grpc-proxy] unimplemented")
}

func (s *grpcProxyServer) streamC(proxyStreamCServer) error {
	return grpc.Errorf(codes.Unimplemented, "[grpc-proxy] unimplemented")
}

func (s *grpcProxyServer) streamS(proxyStreamSServer) error {
	return grpc.Errorf(codes.Unimplemented, "[grpc-proxy] unimplemented")
}

func (s *grpcProxyServer) streamB(grpc.ServerStream) error {
	return grpc.Errorf(codes.Unimplemented, "[grpc-proxy] unimplemented")
}

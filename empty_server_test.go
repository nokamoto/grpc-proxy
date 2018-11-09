package main

import (
	pb "github.com/nokamoto/grpc-proxy/examples/empty-package"
	"golang.org/x/net/context"
)

type emptyServer struct{}

func (s *emptyServer) Reverse(_ context.Context, a *pb.A) (*pb.B, error) {
	b := ""

	for _, v := range a.A {
		b = string(v) + b
	}

	return &pb.B{B: b}, nil
}

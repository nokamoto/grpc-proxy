package main

import (
	"github.com/nokamoto/grpc-proxy/codec"
	"google.golang.org/grpc"
)

type streamSServer interface {
	send(*codec.RawMessage) error
	grpc.ServerStream
}

type proxyStreamSServer struct {
	grpc.ServerStream
}

func (x *proxyStreamSServer) send(m *codec.RawMessage) error {
	return x.ServerStream.SendMsg(m)
}

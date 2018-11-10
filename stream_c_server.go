package main

import (
	"github.com/nokamoto/grpc-proxy/codec"
	"google.golang.org/grpc"
)

type streamCServer interface {
	sendAndClose(*codec.RawMessage) error
	recv() (*codec.RawMessage, error)
	grpc.ServerStream
}

type proxyStreamCServer struct {
	grpc.ServerStream
}

func (x *proxyStreamCServer) sendAndClose(m *codec.RawMessage) error {
	return x.ServerStream.SendMsg(m)
}

func (x *proxyStreamCServer) Recv() (*codec.RawMessage, error) {
	m := new(codec.RawMessage)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

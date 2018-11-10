package main

import (
	"github.com/nokamoto/grpc-proxy/codec"
	"google.golang.org/grpc"
)

type streamBServer interface {
	send(*codec.RawMessage) error
	recv() (*codec.RawMessage, error)
	grpc.ServerStream
}

type proxyStreamBServer struct {
	grpc.ServerStream
}

func (x *proxyStreamBServer) send(m *codec.RawMessage) error {
	return x.ServerStream.SendMsg(m)
}

func (x *proxyStreamBServer) recv() (*codec.RawMessage, error) {
	m := new(codec.RawMessage)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

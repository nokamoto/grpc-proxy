package main

import (
	"google.golang.org/grpc"
)

type streamCServer interface {
	sendAndClose(*message) error
	recv() (*message, error)
	grpc.ServerStream
}

type proxyStreamCServer struct {
	grpc.ServerStream
}

func (x *proxyStreamCServer) sendAndClose(m *message) error {
	return x.ServerStream.SendMsg(m)
}

func (x *proxyStreamCServer) Recv() (*message, error) {
	m := new(message)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

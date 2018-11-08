package main

import (
	"google.golang.org/grpc"
)

type streamBServer interface {
	send(*message) error
	recv() (*message, error)
	grpc.ServerStream
}

type proxyStreamBServer struct {
	grpc.ServerStream
}

func (x *proxyStreamBServer) send(m *message) error {
	return x.ServerStream.SendMsg(m)
}

func (x *proxyStreamBServer) recv() (*message, error) {
	m := new(message)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

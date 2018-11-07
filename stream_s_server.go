package main

import (
	"google.golang.org/grpc"
)

type streamSServer interface {
	send(*message) error
	grpc.ServerStream
}

type proxyStreamSServer struct {
	grpc.ServerStream
}

func (x *proxyStreamSServer) send(m *message) error {
	return x.ServerStream.SendMsg(m)
}

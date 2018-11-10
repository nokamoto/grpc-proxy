package server

import (
	"github.com/nokamoto/grpc-proxy/codec"
	"google.golang.org/grpc"
)

type RawServerStreamS struct {
	grpc.ServerStream
}

func (x *RawServerStreamS) Send(m *codec.RawMessage) error {
	return x.ServerStream.SendMsg(m)
}

func (x *RawServerStreamS) Recv() (*codec.RawMessage, error) {
	m := new(codec.RawMessage)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

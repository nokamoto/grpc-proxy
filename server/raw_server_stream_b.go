package server

import (
	"github.com/nokamoto/grpc-proxy/codec"
	"google.golang.org/grpc"
)

type RawServerStreamB struct {
	grpc.ServerStream
}

func (x *RawServerStreamB) Send(m *codec.RawMessage) error {
	return x.ServerStream.SendMsg(m)
}

func (x *RawServerStreamB) Recv() (*codec.RawMessage, error) {
	m := new(codec.RawMessage)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

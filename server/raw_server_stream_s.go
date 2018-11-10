package server

import (
	"github.com/nokamoto/grpc-proxy/codec"
	"google.golang.org/grpc"
)

// RawServerStreamS implements for server side grpc.ServerStream.
type RawServerStreamS struct {
	grpc.ServerStream
}

// Send sends codec.RawMessage to the downstream.
func (x *RawServerStreamS) Send(m *codec.RawMessage) error {
	return x.ServerStream.SendMsg(m)
}

// Recv receives codec.RawMessage from the downstream.
func (x *RawServerStreamS) Recv() (*codec.RawMessage, error) {
	m := new(codec.RawMessage)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

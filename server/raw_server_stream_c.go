package server

import (
	"github.com/nokamoto/grpc-proxy/codec"
	"google.golang.org/grpc"
)

// RawServerStreamC implements for client side grpc.ServerStream.
type RawServerStreamC struct {
	grpc.ServerStream
}

// SendAndClose sends codec.RawMessage to the downstream.
func (x *RawServerStreamC) SendAndClose(m *codec.RawMessage) error {
	return x.ServerStream.SendMsg(m)
}

// Recv receives codec.RawMessage from the downstream.
func (x *RawServerStreamC) Recv() (*codec.RawMessage, error) {
	m := new(codec.RawMessage)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

package server

import (
	"github.com/nokamoto/grpc-proxy/codec"
	"google.golang.org/grpc"
)

// RawServerStreamB implements for bidirectional grpc.ServerStream.
type RawServerStreamB struct {
	grpc.ServerStream
}

// Send sends codec.RawMessage to the downstream.
func (x *RawServerStreamB) Send(m *codec.RawMessage) error {
	return x.ServerStream.SendMsg(m)
}

// Recv receives codec.RawMessage from the downstream.
func (x *RawServerStreamB) Recv() (*codec.RawMessage, error) {
	m := new(codec.RawMessage)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

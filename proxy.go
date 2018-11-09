package main

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type proxy struct {
	con *grpc.ClientConn
}

func newProxy(address string) (*proxy, error) {
	opts := []grpc.DialOption{grpc.WithInsecure()}

	con, err := grpc.Dial(address, opts...)
	if err != nil {
		return nil, err
	}

	return &proxy{con}, err
}

func (p *proxy) invokeUnary(ctx context.Context, m *message, method string) (*message, error) {
	rep := new(message)
	err := p.con.Invoke(ctx, method, m, rep, grpc.CallCustomCodec(codec{}))
	return rep, err
}

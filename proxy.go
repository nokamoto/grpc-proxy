package main

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"io"
	"time"
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

func (p *proxy) invokeStreamC(downstream proxyStreamCServer, desc *grpc.StreamDesc, method string) error {
	// todo: timeout configuration
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	upstream, err := p.con.NewStream(ctx, desc, method, grpc.CallCustomCodec(codec{}))
	if err != nil {
		return grpc.Errorf(codes.Internal, "[grpc-proxy] stream c error: %s", err)
	}

	for {
		m, err := downstream.Recv()

		if err == io.EOF {
			res := new(message)

			err = upstream.CloseSend()
			if err != nil {
				return err
			}

			err = upstream.RecvMsg(res)
			if err != nil {
				return err
			}

			return downstream.sendAndClose(res)
		}

		if err != nil {
			return err
		}

		err = upstream.SendMsg(m)
		if err != nil {
			return err
		}
	}
}

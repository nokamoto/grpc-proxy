package main

import (
	"io"
	"time"

	"github.com/nokamoto/grpc-proxy/codec"
	"github.com/nokamoto/grpc-proxy/server"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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

func (p *proxy) invokeUnary(ctx context.Context, m *codec.RawMessage, method string) (*codec.RawMessage, error) {
	rep := new(codec.RawMessage)
	err := p.con.Invoke(ctx, method, m, rep, grpc.CallCustomCodec(codec.RawCodec{}))
	return rep, err
}

func (p *proxy) invokeStreamC(downstream server.RawServerStreamC, desc *grpc.StreamDesc, method string) error {
	// todo: timeout configuration.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	upstream, err := p.con.NewStream(ctx, desc, method, grpc.CallCustomCodec(codec.RawCodec{}))
	if err != nil {
		return grpc.Errorf(codes.Internal, "[grpc-proxy] stream c error: %s", err)
	}

	for {
		m, err := downstream.Recv()

		if err == io.EOF {
			res := new(codec.RawMessage)

			err = upstream.CloseSend()
			if err != nil {
				return err
			}

			err = upstream.RecvMsg(res)
			if err != nil {
				return err
			}

			return downstream.SendAndClose(res)
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

func (p *proxy) invokeStreamS(downstream server.RawServerStreamS, desc *grpc.StreamDesc, method string) error {
	// todo: timeout configuration.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	upstream, err := p.con.NewStream(ctx, desc, method, grpc.CallCustomCodec(codec.RawCodec{}))
	if err != nil {
		return grpc.Errorf(codes.Internal, "[grpc-proxy] stream s error: %s", err)
	}

	req, err := downstream.Recv()
	if err != nil {
		return err
	}

	err = upstream.SendMsg(req)
	if err != nil {
		return err
	}

	for {
		m := new(codec.RawMessage)
		err = upstream.RecvMsg(m)
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		err = downstream.Send(m)
		if err != nil {
			return err
		}
	}
}

func (p *proxy) invokeStreamB(downstream server.RawServerStreamB, desc *grpc.StreamDesc, method string) error {
	// todo: timeout configuration.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	upstream, err := p.con.NewStream(ctx, desc, method, grpc.CallCustomCodec(codec.RawCodec{}))
	if err != nil {
		return grpc.Errorf(codes.Internal, "[grpc-proxy] stream b error: %s", err)
	}

	chanDownstream := make(chan error)

	go func() {
		for {
			m := new(codec.RawMessage)

			err := downstream.RecvMsg(m)
			if err == io.EOF {
				chanDownstream <- nil
				break
			}

			if err != nil {
				chanDownstream <- err
				break
			}

			err = upstream.SendMsg(m)
			if err != nil {
				chanDownstream <- err
				break
			}
		}
	}()

	go func() {
		for {
			m := new(codec.RawMessage)

			err := upstream.RecvMsg(m)
			if err == io.EOF {
				chanDownstream <- nil
				break
			}

			if err != nil {
				chanDownstream <- err
				break
			}

			err = downstream.SendMsg(m)
			if err != nil {
				chanDownstream <- err
				break
			}
		}
	}()

	err = <-chanDownstream

	// close anyway.
	upstream.CloseSend()

	return err
}

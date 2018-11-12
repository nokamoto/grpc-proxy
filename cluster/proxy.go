package cluster

import (
	"fmt"
	"github.com/nokamoto/grpc-proxy/codec"
	"github.com/nokamoto/grpc-proxy/server"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"io"
	"os"
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
	octx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs())

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		octx = metadata.NewOutgoingContext(octx, md)
	}

	fmt.Println(md)

	var header, trailer metadata.MD

	rep := new(codec.RawMessage)

	err := p.con.Invoke(octx, method, m, rep, grpc.CallCustomCodec(codec.RawCodec{}), grpc.Header(&header), grpc.Trailer(&trailer))

	grpc.SetHeader(ctx, header)
	grpc.SetTrailer(ctx, trailer)

	return rep, err
}

func (p *proxy) invokeStreamC(downstream server.RawServerStreamC, desc *grpc.StreamDesc, method string) error {
	octx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs())

	md, ok := metadata.FromIncomingContext(downstream.Context())
	if ok {
		octx = metadata.NewOutgoingContext(octx, md)
	}

	fmt.Println(md)

	var header, trailer metadata.MD

	upstream, err := p.con.NewStream(octx, desc, method, grpc.CallCustomCodec(codec.RawCodec{}), grpc.Header(&header), grpc.Trailer(&trailer))
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

			downstream.SetHeader(header)
			downstream.SetTrailer(trailer)

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
	octx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs())

	md, ok := metadata.FromIncomingContext(downstream.Context())
	if ok {
		octx = metadata.NewOutgoingContext(octx, md)
	}

	fmt.Println(md)

	upstream, err := p.con.NewStream(octx, desc, method, grpc.CallCustomCodec(codec.RawCodec{}))
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

	headerCheck := true

	for {
		m := new(codec.RawMessage)
		err = upstream.RecvMsg(m)
		if err == io.EOF {
			downstream.SetTrailer(upstream.Trailer())
			return nil
		}

		if err != nil {
			return err
		}

		if headerCheck {
			h, err := upstream.Header()
			if err != nil {
				fmt.Fprintf(os.Stderr, "stream s upstream header error: %s", err)
			}
			if err == nil {
				downstream.SetHeader(h)
			}
			headerCheck = false
		}

		err = downstream.Send(m)
		if err != nil {
			return err
		}
	}
}

// TODO: matadata proxy
func (p *proxy) invokeStreamB(downstream server.RawServerStreamB, desc *grpc.StreamDesc, method string) error {
	octx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs())

	md, ok := metadata.FromIncomingContext(downstream.Context())
	if ok {
		octx = metadata.NewOutgoingContext(octx, md)
	}

	fmt.Println(md)

	upstream, err := p.con.NewStream(octx, desc, method, grpc.CallCustomCodec(codec.RawCodec{}))
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

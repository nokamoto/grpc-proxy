package server

import (
	"github.com/nokamoto/grpc-proxy/codec"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// RawUnaryHandler returns grpc.methodHandler for codec.RawMessage.
func RawUnaryHandler(fullMethod string) func(interface{}, context.Context, func(interface{}) error, grpc.UnaryServerInterceptor) (interface{}, error) {
	return func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
		in := new(codec.RawMessage)
		if err := dec(in); err != nil {
			return nil, err
		}
		if interceptor == nil {
			return srv.(Server).Unary(ctx, in, fullMethod)
		}
		info := &grpc.UnaryServerInfo{
			Server:     srv,
			FullMethod: fullMethod,
		}
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.(Server).Unary(ctx, req.(*codec.RawMessage), fullMethod)
		}
		return interceptor(ctx, in, info, handler)
	}
}

// RawServerStreamCHandler returns grpc.StreamHandler for the client side codec.RawMessage stream.
func RawServerStreamCHandler(fullMethod string, desc *grpc.StreamDesc) grpc.StreamHandler {
	return func(srv interface{}, stream grpc.ServerStream) error {
		return srv.(Server).StreamC(RawServerStreamC{stream}, desc, fullMethod)
	}
}

// RawServerStreamSHandler returns grpc.StreamHandler for the server side codec.RawMessage stream.
func RawServerStreamSHandler(fullMethod string, desc *grpc.StreamDesc) grpc.StreamHandler {
	return func(srv interface{}, stream grpc.ServerStream) error {
		return srv.(Server).StreamS(RawServerStreamS{stream}, desc, fullMethod)
	}
}

// RawServerStreamBHandler returns grpc.StreamHandler for the bidirectional codec.RawMessage stream.
func RawServerStreamBHandler(fullMethod string, desc *grpc.StreamDesc) grpc.StreamHandler {
	return func(srv interface{}, stream grpc.ServerStream) error {
		return srv.(Server).StreamB(RawServerStreamB{stream}, desc, fullMethod)
	}
}

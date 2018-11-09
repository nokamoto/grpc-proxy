package main

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func unaryProxyHandler(fullMethod string) func(interface{}, context.Context, func(interface{}) error, grpc.UnaryServerInterceptor) (interface{}, error) {
	return func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
		in := new(message)
		if err := dec(in); err != nil {
			return nil, err
		}
		if interceptor == nil {
			return srv.(proxyServer).unary(ctx, in, fullMethod)
		}
		info := &grpc.UnaryServerInfo{
			Server:     srv,
			FullMethod: fullMethod,
		}
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.(proxyServer).unary(ctx, req.(*message), fullMethod)
		}
		return interceptor(ctx, in, info, handler)
	}
}

func streamCProxyHandler(fullMethod string, desc *grpc.StreamDesc) func(interface{}, grpc.ServerStream) error {
	return func(srv interface{}, stream grpc.ServerStream) error {
		return srv.(proxyServer).streamC(proxyStreamCServer{stream}, desc, fullMethod)
	}
}

func streamSProxyHandler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(proxyServer).streamS(proxyStreamSServer{stream})
}

func streamBProxyHandler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(proxyServer).streamB(stream)
}

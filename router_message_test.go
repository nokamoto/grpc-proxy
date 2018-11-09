package main

import (
	empty "github.com/nokamoto/grpc-proxy/examples/empty-package"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"testing"
)

func Test_router_message(t *testing.T) {
	testWithEmptyServer(t, &emptyServer{}, func(ctx context.Context, cc *grpc.ClientConn) {
		c := empty.NewServiceClient(cc)

		b, err := c.Reverse(ctx, &empty.A{A: "abcdefg"})
		if err != nil {
			t.Error(err)
		}

		if b.B != "gfedcba" {
			t.Errorf("%s != gfedcba", b.B)
		}
	})
}

package test

import (
	"github.com/nokamoto/grpc-proxy/proxy"
	"testing"
)

func beforeEachGrpcProxy(t *testing.T, port int, pb, yml string) func() {
	t.Helper()
	srv, err := proxy.NewServer(port, pb, yml)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		srv.Serve()
	}()

	return func() {
		srv.GracefulStop()
	}
}

package main

import (
	"flag"
	"github.com/nokamoto/grpc-proxy/proxy"
)

func main() {
	var (
		port = flag.Int("p", 9000, "gRPC server port")
		pb   = flag.String("pb", "", "file descriptor protocol buffers filepath")
		yml  = flag.String("yaml", "", "yaml configuration filepath")
	)

	flag.Parse()

	srv, err := proxy.NewServer(*port, *pb, *yml)
	if err != nil {
		panic(err)
	}

	srv.Serve()
}

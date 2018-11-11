package main

import (
	"flag"
	"fmt"
	"github.com/nokamoto/grpc-proxy/proxy"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func main() {
	var (
		port = flag.Int("p", 9000, "gRPC server port")
		pb   = flag.String("pb", "", "file descriptor protocol buffers filepath")
		yml  = flag.String("yaml", "", "yaml configuration filepath")
		prom = flag.Int("metrics", 9001, "Prometheus exporter port")
	)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		fmt.Println(http.ListenAndServe(fmt.Sprintf(":%d", prom), nil))
	}()

	flag.Parse()

	srv, err := proxy.NewServer(*port, *pb, *yml)
	if err != nil {
		panic(err)
	}

	srv.Serve()
}

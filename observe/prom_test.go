package observe

import (
	"fmt"
	"github.com/nokamoto/grpc-proxy/yaml"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc/codes"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestProm_Observe(t *testing.T) {
	port := 9000

	afterEachProm := beforeEachProm(t, port)
	defer afterEachProm()

	y, err := yaml.NewYaml("../examples/empty-package/example.yaml")
	if err != nil {
		t.Fatal(err)
	}

	p, err := NewProm(y.Observe.Prom[0])
	if err != nil {
		t.Fatal(err)
	}

	write := func(method string, code codes.Code, req int, res int, nanos time.Duration) {
		err = p.Observe(method, code, req, res, nanos)
		if err != nil {
			t.Fatal(err)
		}
	}

	write("x", codes.OK, 32, 64, 1500 * time.Millisecond)
	write("x", codes.OK, 64, 128, 750 * time.Millisecond)
	write("y", codes.OK, 128, 256, 250 * time.Millisecond)

	retriveProm(t, `request_count{method="x",status="OK"} 2`, port)
	retriveProm(t, `request_count{method="y",status="OK"} 1`, port)

	retriveProm(t, `latency_seconds_bucket{method="x",status="OK",le="+Inf"} 2`, port)
	retriveProm(t, `latency_seconds_bucket{method="x",status="OK",le="1"} 2`, port)
	retriveProm(t, `latency_seconds_bucket{method="x",status="OK",le="0.5"} 1`, port)
	retriveProm(t, `latency_seconds_bucket{method="y",status="OK",le="0.5"} 1`, port)
}

func retriveProm(t *testing.T, expected string, port int) {
	res, err := http.Get(fmt.Sprintf("http://localhost:%d/metrics", port))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	for _, line := range strings.Split(string(bytes), "\n") {
		if expected == line {
			return
		}
	}

	t.Fatalf("%s not found: %s", expected, string(bytes))
}

func beforeEachProm(t *testing.T, port int) func() {
	http.Handle("/metrics", promhttp.Handler())
	srv := http.Server{Addr: fmt.Sprintf(":%d", port), Handler: nil}

	go func() {
		srv.ListenAndServe()
	}()

	return func() {
		srv.Close()
	}
}

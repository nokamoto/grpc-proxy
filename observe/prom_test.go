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

	write("x", codes.OK, 5, 10, 3*time.Millisecond)
	write("x", codes.OK, 10, 15, 4*time.Millisecond)
	write("y", codes.OK, 15, 20, 5*time.Millisecond)

	counter := retriveProm(t, `request_count{method="x",status="OK"}`, port)
	expected := `request_count{method="x",status="OK"} 2`
	if counter != expected {
		t.Errorf("%s != %s", counter, expected)
	}

	counter = retriveProm(t, `request_count{method="y",status="OK"}`, port)
	expected = `request_count{method="y",status="OK"} 1`
	if counter != expected {
		t.Errorf("%s != %s", counter, expected)
	}
}

func retriveProm(t *testing.T, name string, port int) string {
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
		if strings.HasPrefix(line, name) {
			return line
		}
	}

	t.Fatalf("%s not found: %s", name, string(bytes))

	return ""
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

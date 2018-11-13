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

func TestProm_NewProm(t *testing.T) {
	yml, err := yaml.NewYaml("../testdata/yaml/prom_new.yaml")
	if err != nil {
		t.Fatal(err)
	}

	_, err = NewProm(yml.Observe.Prom[0])
	if err != nil {
		t.Fatal(err)
	}

	_, err = NewProm(yml.Observe.Prom[1])
	if err != nil {
		t.Fatal(err)
	}
}

func TestProm_NewProm_duplicated(t *testing.T) {
	yml, err := yaml.NewYaml("../testdata/yaml/prom_duplicated.yaml")
	if err != nil {
		t.Fatal(err)
	}

	_, err = NewProm(yml.Observe.Prom[0])
	if err != nil {
		t.Fatal(err)
	}

	_, err = NewProm(yml.Observe.Prom[1])
	if err == nil {
		t.Fatal()
	}
}

func TestProm_Observe(t *testing.T) {
	port := 9000

	afterEachProm := beforeEachProm(t, port)
	defer afterEachProm()

	yml, err := yaml.NewYaml("../testdata/yaml/prom.yaml")
	if err != nil {
		t.Fatal(err)
	}

	p, err := NewProm(yml.Observe.Prom[0])
	if err != nil {
		t.Fatal(err)
	}

	write := func(method string, code codes.Code, req int, res int, nanos time.Duration) {
		err = p.Observe(method, code, req, res, nanos)
		if err != nil {
			t.Fatal(err)
		}
	}

	x := "x"
	y := "y"

	write(x, codes.OK, 127, 255, 1500*time.Millisecond)
	write(x, codes.OK, 255, 127, 750*time.Millisecond)
	write(y, codes.OK, 511, 63, 250*time.Millisecond)

	counter := func(method string, code codes.Code, n int) {
		retriveProm(t, fmt.Sprintf(`default_request_count{method="%s",status="%s"} %d`, method, code, n), port)
	}

	counter(x, codes.OK, 2)
	counter(y, codes.OK, 1)

	hist := func(bucket string, method string, code codes.Code, le string, n int) {
		retriveProm(t, fmt.Sprintf(`default_%s{method="%s",status="%s",le="%s"} %d`, bucket, method, code, le, n), port)
	}

	latency := "latency_seconds_bucket"
	hist(latency, x, codes.OK, "+Inf", 2)
	hist(latency, x, codes.OK, "1", 1)
	hist(latency, x, codes.OK, "0.5", 0)

	hist(latency, y, codes.OK, "+Inf", 1)
	hist(latency, y, codes.OK, "1", 1)
	hist(latency, y, codes.OK, "0.5", 1)

	request := "request_bytes_bucket"
	hist(request, x, codes.OK, "+Inf", 2)
	hist(request, x, codes.OK, "256", 2)
	hist(request, x, codes.OK, "128", 1)

	hist(request, y, codes.OK, "+Inf", 1)
	hist(request, y, codes.OK, "256", 0)
	hist(request, y, codes.OK, "128", 0)

	response := "response_bytes_bucket"
	hist(response, x, codes.OK, "+Inf", 2)
	hist(response, x, codes.OK, "128", 1)
	hist(response, x, codes.OK, "64", 0)

	hist(response, y, codes.OK, "+Inf", 1)
	hist(response, y, codes.OK, "128", 1)
	hist(response, y, codes.OK, "64", 1)
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

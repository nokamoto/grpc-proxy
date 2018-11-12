package observe

import (
	"fmt"
	"github.com/nokamoto/grpc-proxy/yaml"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/codes"
	"sort"
	"time"
)

type Prom interface {
	Observe(string, codes.Code, int, int, time.Duration) error
	Destroy()
}

func NewProm(c yaml.Prom) (Prom, error) {
	labels := []string{"method", "status"}

	counter := prometheus.NewCounterVec(prometheus.CounterOpts{Name: fmt.Sprintf("%s_request_count", c.Name)}, labels)

	err := prometheus.Register(counter)
	if err != nil {
		return nil, err
	}

	sorted := func(f []float64) []float64 {
		sort.Sort(sort.Float64Slice(f))
		return f
	}

	hist := func(name string, buckets []float64) (*prometheus.HistogramVec, error) {
		h := prometheus.NewHistogramVec(prometheus.HistogramOpts{Name: fmt.Sprintf("%s_%s", c.Name, name), Buckets: sorted(buckets)}, labels)
		err := prometheus.Register(h)
		return h, err
	}

	latency, err := hist("latency_seconds", c.Buckets.LatencySeconds)
	if err != nil {
		return nil, err
	}

	req, err := hist("request_bytes", c.Buckets.RequestBytes)
	if err != nil {
		return nil, err
	}

	res, err := hist("response_bytes", c.Buckets.ResponseBytes)
	if err != nil {
		return nil, err
	}

	return &prom{counter: counter, hist: hists{req: req, res: res, latency: latency}}, nil
}

type prom struct {
	counter *prometheus.CounterVec
	hist    hists
}

type hists struct {
	req     *prometheus.HistogramVec
	res     *prometheus.HistogramVec
	latency *prometheus.HistogramVec
}

func (p *prom) Observe(method string, code codes.Code, req int, res int, nanos time.Duration) error {
	labels := []string{method, code.String()}

	c, err := p.counter.GetMetricWithLabelValues(labels...)
	if err != nil {
		return err
	}

	hreq, err := p.hist.req.GetMetricWithLabelValues(labels...)
	if err != nil {
		return err
	}

	hres, err := p.hist.res.GetMetricWithLabelValues(labels...)
	if err != nil {
		return err
	}

	hlatency, err := p.hist.latency.GetMetricWithLabelValues(labels...)
	if err != nil {
		return err
	}

	c.Inc()
	hreq.Observe(float64(req))
	hres.Observe(float64(res))
	hlatency.Observe(float64(nanos) / (1000 * 1000 * 1000))

	return nil
}

func (p *prom) Destroy() {
	prometheus.Unregister(p.counter)
	prometheus.Unregister(p.hist.latency)
	prometheus.Unregister(p.hist.req)
	prometheus.Unregister(p.hist.res)
}

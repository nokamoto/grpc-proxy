package observe

import (
	"github.com/nokamoto/grpc-proxy/yaml"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/codes"
	"sort"
	"time"
)

type Prom interface {
	Observe(string, codes.Code, int, int, time.Duration) error
}

func NewProm(c yaml.Prom) (Prom, error) {
	labels := []string{"method", "status"}
	counter := prometheus.NewCounterVec(prometheus.CounterOpts{Name: "request_count"}, labels)

	err := prometheus.Register(counter)
	if err != nil {
		return nil, err
	}

	sorted := func(f []float64) []float64 {
		sort.Sort(sort.Float64Slice(f))
		return f
	}

	labels = []string{"method"}
	latency := prometheus.NewHistogramVec(prometheus.HistogramOpts{Name: "latency_seconds", Buckets: sorted(c.Buckets.LatencySeconds)}, labels)
	req := prometheus.NewHistogramVec(prometheus.HistogramOpts{Name: "request_bytes", Buckets: sorted(c.Buckets.RequestBytes)}, labels)
	res := prometheus.NewHistogramVec(prometheus.HistogramOpts{Name: "response_bytes", Buckets: sorted(c.Buckets.ResponseBytes)}, labels)

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
	c, err := p.counter.GetMetricWithLabelValues(method, code.String())
	if err != nil {
		return err
	}

	hreq, err := p.hist.req.GetMetricWithLabelValues(method)
	if err != nil {
		return err
	}

	hres, err := p.hist.res.GetMetricWithLabelValues(method)
	if err != nil {
		return err
	}

	hlatency, err := p.hist.latency.GetMetricWithLabelValues(method)
	if err != nil {
		return err
	}

	c.Inc()
	hreq.Observe(float64(req))
	hres.Observe(float64(res))
	hlatency.Observe(float64(nanos / (1000 * 1000 * 1000)))

	return nil
}

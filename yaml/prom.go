package yaml

// Prom represents a configuration of Prometheus gRPC proxy metrics.
type Prom struct {
	Name    string
	Buckets struct {
		LatencySeconds []float64 `yaml:"latency-seconds"`
		RequestBytes   []float64 `yaml:"request-bytes"`
		ResponseBytes  []float64 `yaml:"response-bytes"`
	}
}

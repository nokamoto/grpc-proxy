package yaml

// Observe represents a configuration of gRPC observability.
type Observe struct {
	Observe struct {
		Logs []Log
	}
}

// Log represents a configuration of gRPC access logging.
type Log struct {
	Name string
	File string
}

package yaml

type observe struct {
	Logs []Log
	Prom []Prom
}

// Log represents a configuration of gRPC access logging.
type Log struct {
	Name string
	File string
}

package yaml

// Cluster represent a configuration of a single upstream cluster.
type Cluster struct {
	Name       string
	RoundRobin []string `yaml:"round_robin"`
}

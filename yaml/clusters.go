package yaml

// Clusters represents a configuration of upstream clusters.
type Clusters struct {
	Clusters []Cluster
}

type Cluster struct {
	Name       string
	RoundRobin []string `yaml:"round_robin"`
}

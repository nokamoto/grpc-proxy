package yaml

// Clusters represents a configuration of upstream clusters.
type Clusters struct {
	Clusters []Cluster
}

// Cluster represent a configuration of a single upstream cluster.
type Cluster struct {
	Name       string
	RoundRobin []string `yaml:"round_robin"`
}

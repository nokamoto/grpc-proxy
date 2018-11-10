package yaml

type Clusters struct {
	Clusters []Cluster
}

type Cluster struct {
	Name       string
	RoundRobin []string `yaml:"round_robin"`
}

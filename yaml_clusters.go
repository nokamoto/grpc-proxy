package main

type yamlClusters struct {
	Clusters []yamlCluster
}

type yamlCluster struct {
	Name       string
	RoundRobin []string `yaml:"round_robin"`
}

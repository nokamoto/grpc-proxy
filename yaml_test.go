package main

import (
	"testing"
)

func Test_newYaml_empty_package(t *testing.T) {
	routes, clusters, err := newYaml("examples/empty-package/example.yaml")
	if err != nil {
		t.Fatal(err)
	}

	if s := routes.Routes[0].Method.Prefix; s != "/" {
		t.Errorf("%s != /", s)
	}
	if s := routes.Routes[0].Cluster.Name; s != "local" {
		t.Errorf("%s != local", s)
	}

	if s := clusters.Clusters[0].Name; s != "local" {
		t.Errorf("%s != local", s)
	}
	if s := clusters.Clusters[0].RoundRobin[0]; s != "localhost:9001" {
		t.Errorf("%s != localhost:9001", s)
	}
}

func Test_newYaml_ping(t *testing.T) {
	routes, clusters, err := newYaml("examples/ping/example.yaml")
	if err != nil {
		t.Fatal(err)
	}

	if s := routes.Routes[0].Method.Prefix; s != "/" {
		t.Errorf("%s != /", s)
	}
	if s := routes.Routes[0].Cluster.Name; s != "local" {
		t.Errorf("%s != local", s)
	}

	if s := clusters.Clusters[0].Name; s != "local" {
		t.Errorf("%s != local", s)
	}
	if s := clusters.Clusters[0].RoundRobin[0]; s != "localhost:9002" {
		t.Errorf("%s != localhost:9002", s)
	}
}

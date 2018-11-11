package yaml

import (
	"testing"
)

func Test_NewYaml_empty_package(t *testing.T) {
	routes, clusters, observe, err := NewYaml("../examples/empty-package/example.yaml")
	if err != nil {
		t.Fatal(err)
	}

	if s := routes.Routes[0].Method.Prefix; s != "/" {
		t.Errorf("%s != /", s)
	}
	if s := routes.Routes[0].Cluster.Name; s != "local" {
		t.Errorf("%s != local", s)
	}
	if s := routes.Routes[0].Observe.Log.Name; s != "stdout" {
		t.Errorf("%s != stdout", s)
	}

	if s := clusters.Clusters[0].Name; s != "local" {
		t.Errorf("%s != local", s)
	}
	if s := clusters.Clusters[0].RoundRobin[0]; s != "localhost:9001" {
		t.Errorf("%s != localhost:9001", s)
	}

	if s := observe.Observe.Logs[0].Name; s != "stdout" {
		t.Errorf("%s != stdout", s)
	}
	if s := observe.Observe.Logs[0].File; s != "/dev/stdout" {
		t.Errorf("%s != /dev/stdout", s)
	}
}

func Test_NewYaml_ping(t *testing.T) {
	routes, clusters, observe, err := NewYaml("../examples/ping/example.yaml")
	if err != nil {
		t.Fatal(err)
	}

	if s := routes.Routes[0].Method.Prefix; s != "/" {
		t.Errorf("%s != /", s)
	}
	if s := routes.Routes[0].Cluster.Name; s != "local" {
		t.Errorf("%s != local", s)
	}
	if s := routes.Routes[0].Observe.Log.Name; s != "stdout" {
		t.Errorf("%s != stdout", s)
	}

	if s := clusters.Clusters[0].Name; s != "local" {
		t.Errorf("%s != local", s)
	}
	if s := clusters.Clusters[0].RoundRobin[0]; s != "localhost:9002" {
		t.Errorf("%s != localhost:9002", s)
	}

	if s := observe.Observe.Logs[0].Name; s != "stdout" {
		t.Errorf("%s != stdout", s)
	}
	if s := observe.Observe.Logs[0].File; s != "/dev/stdout" {
		t.Errorf("%s != /dev/stdout", s)
	}
}

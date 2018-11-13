package yaml

import (
	"testing"
)

func Test_NewYaml_ping(t *testing.T) {
	yaml, err := NewYaml("../testdata/yaml/ping.yaml")
	if err != nil {
		t.Fatal(err)
	}

	if s := yaml.Routes[0].Method.Prefix; s != "/" {
		t.Errorf("%s != /", s)
	}
	if s := yaml.Routes[0].Cluster.Name; s != "local" {
		t.Errorf("%s != local", s)
	}
	if s := yaml.Routes[0].Observe.Log.Name; s != "stdout" {
		t.Errorf("%s != stdout", s)
	}

	if s := yaml.Clusters[0].Name; s != "local" {
		t.Errorf("%s != local", s)
	}
	if s := yaml.Clusters[0].RoundRobin[0]; s != "localhost:9002" {
		t.Errorf("%s != localhost:9002", s)
	}

	if s := yaml.Observe.Logs[0].Name; s != "stdout" {
		t.Errorf("%s != stdout", s)
	}
	if s := yaml.Observe.Logs[0].File; s != "/dev/stdout" {
		t.Errorf("%s != /dev/stdout", s)
	}
}

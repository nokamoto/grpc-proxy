package yaml

import (
	"reflect"
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
	if s := yaml.Routes[0].Observe.Log.Name; s != nil {
		t.Fatal(*s)
	}
	if s := yaml.Routes[0].Observe.Prom.Name; s != nil {
		t.Fatal(*s)
	}

	if s := yaml.Clusters[0].Name; s != "local" {
		t.Errorf("%s != local", s)
	}
	if s := yaml.Clusters[0].RoundRobin[0]; s != "localhost:9002" {
		t.Errorf("%s != localhost:9002", s)
	}

	if l := len(yaml.Observe.Logs); l != 0 {
		t.Errorf("%d != 0", l)
	}
	if l := len(yaml.Observe.Prom); l != 0 {
		t.Errorf("%d != 0", l)
	}
}

func Test_NewYaml_ping_log(t *testing.T) {
	yaml, err := NewYaml("../testdata/yaml/ping_log.yaml")
	if err != nil {
		t.Fatal(err)
	}

	if s := yaml.Routes[0].Observe.Log.Name; *s != "stdout" {
		t.Errorf("%s != stdout", *s)
	}

	if yaml.Observe.Logs == nil {
		t.Fatal()
	}

	logs := yaml.Observe.Logs

	if s := logs[0].Name; s != "stdout" {
		t.Errorf("%s != stdout", s)
	}
	if s := logs[0].File; s != "/dev/stdout" {
		t.Errorf("%s != /dev/stdout", s)
	}
}

func Test_NewYaml_ping_prom(t *testing.T) {
	yaml, err := NewYaml("../testdata/yaml/ping_prom.yaml")
	if err != nil {
		t.Fatal(err)
	}

	if s := yaml.Routes[0].Observe.Prom.Name; *s != "default" {
		t.Errorf("%s != stdout", *s)
	}

	if yaml.Observe.Prom == nil {
		t.Fatal()
	}

	prom := yaml.Observe.Prom

	if s := prom[0].Name; s != "default" {
		t.Errorf("%s != default", s)
	}
	if x, y := prom[0].Buckets.LatencySeconds, []float64{1.0, 0.5}; !reflect.DeepEqual(x, y) {
		t.Errorf("%v != %v", x, y)
	}
	if x, y := prom[0].Buckets.RequestBytes, []float64{256.0, 128.0}; !reflect.DeepEqual(x, y) {
		t.Errorf("%v != %v", x, y)
	}
	if x, y := prom[0].Buckets.ResponseBytes, []float64{128.0, 64.0}; !reflect.DeepEqual(x, y) {
		t.Errorf("%v != %v", x, y)
	}
}

package yaml

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"reflect"
	"testing"
)

func Test_NewYaml_ping(t *testing.T) {
	expected := &Yaml{
		Routes: []Route{
			prefix("/", routeObserve{}, routeAuth{}),
		},
		Clusters: defaultClusters(),
	}

	test(t, "ping.yaml", expected)
}

func Test_NewYaml_ping_method_equal(t *testing.T) {
	expected := &Yaml{
		Routes: []Route{
			eq("/ping.PingService/Send", routeObserve{}, routeAuth{}),
			eq("/ping.PingService/SendStreamC", routeObserve{}, routeAuth{}),
			eq("/ping.PingService/SendStreamS", routeObserve{}, routeAuth{}),
			eq("/ping.PingService/SendStreamB", routeObserve{}, routeAuth{}),
		},
		Clusters: defaultClusters(),
	}

	test(t, "ping_method_equal.yaml", expected)
}

func Test_NewYaml_ping_log(t *testing.T) {
	expected := &Yaml{
		Routes: []Route{
			prefix(
				"/",
				routeObserve{
					Log: routeObserveLog{
						Name: ref("stdout"),
					},
				},
				routeAuth{},
			),
		},
		Clusters: defaultClusters(),
		Observe: observe{
			Logs: []Log{
				Log{
					Name: "stdout",
					File: "/dev/stdout",
				},
			},
		},
	}

	test(t, "ping_log.yaml", expected)
}

func Test_NewYaml_ping_prom(t *testing.T) {
	expected := &Yaml{
		Routes: []Route{
			prefix(
				"/",
				routeObserve{
					Prom: routeObserveProm{
						Name: ref("default"),
					},
				},
				routeAuth{},
			),
		},
		Clusters: defaultClusters(),
		Observe: observe{
			Prom: []Prom{
				Prom{
					Name: "default",
					Buckets: promBuckets{
						LatencySeconds: []float64{1.0, 0.5},
						RequestBytes:   []float64{256.0, 128.0},
						ResponseBytes:  []float64{128.0, 64.0},
					},
				},
			},
		},
	}

	test(t, "ping_prom.yaml", expected)
}

func Test_NewYaml_ping_anonymous(t *testing.T) {
	expected := &Yaml{
		Routes: []Route{
			prefix(
				"/",
				routeObserve{},
				routeAuth{
					KeyAuth: &routeKeyAuth{
						Metadata:        "x-apikey",
						Anonymous:       ref("anonymous_users"),
						HideCredentials: true,
					},
				},
			),
		},
		Clusters: defaultClusters(),
		Clients: []Client{
			client("anonymous_users"),
			client("admin_users", "ure3Wee2", "shae5Aig"),
			client("developers", "Pae4shua"),
		},
	}

	test(t, "ping_anonymous.yaml", expected)
}

func ref(s string) *string {
	return &s
}

func Test_NewYaml_errors(t *testing.T) {
	check := func(y string) {
		yaml, err := NewYaml(y)
		if err == nil {
			t.Fatalf("%v", yaml)
		}
	}

	check("../testdata/yaml/yaml_ambiguous_method.yaml")
	check("../testdata/yaml/yaml_no_method.yaml")
}

func defaultClusters() []Cluster {
	return []Cluster{
		Cluster{
			Name:       "local",
			RoundRobin: []string{"localhost:9002"},
		},
	}
}

func client(name string, keys ...string) Client {
	return Client{
		Name: name,
		Keys: keys,
	}
}

func prefix(s string, observe routeObserve, auth routeAuth) Route {
	return Route{
		Method:  routeMethod{Prefix: &s},
		Cluster: routeCluster{Name: "local"},
		Observe: observe,
		Auth:    auth,
	}
}

func eq(s string, observe routeObserve, auth routeAuth) Route {
	return Route{
		Method:  routeMethod{Equal: &s},
		Cluster: routeCluster{Name: "local"},
		Observe: observe,
		Auth:    auth,
	}
}

func test(t *testing.T, file string, expected *Yaml) {
	t.Helper()

	actual, err := NewYaml(fmt.Sprintf("../testdata/yaml/%s", file))
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		l, err := yaml.Marshal(actual)
		if err != nil {
			t.Fatal(err)
		}

		r, err := yaml.Marshal(expected)
		if err != nil {
			t.Fatal(err)
		}

		t.Errorf("%s != %s", string(l), string(r))
	}
}

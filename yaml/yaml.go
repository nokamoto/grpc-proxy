package yaml

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

// Yaml represents a configuration of a gRPC proxy.
type Yaml struct {
	Routes   []Route
	Clusters []Cluster
	Observe  struct {
		Logs []Log
	}
}

// NewYaml returns routes and clusters configurations read from the yaml file.
func NewYaml(y string) (*Yaml, error) {
	bytes, err := ioutil.ReadFile(y)
	if err != nil {
		return nil, err
	}

	yml := &Yaml{}
	if err := yaml.Unmarshal(bytes, yml); err != nil {
		return nil, err
	}

	return yml, nil
}

// FindByFullMethod returns all routes match fully qualified the gRPC service method name.
func (y *Yaml) FindByFullMethod(name string) []Route {
	routes := make([]Route, 0)

	for _, route := range y.Routes {
		if strings.HasPrefix(name, route.Method.Prefix) {
			routes = append(routes, route)
		}
	}

	return routes
}

package yaml

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

// Yaml represents a configuration of a gRPC proxy.
type Yaml struct {
	Routes   []Route
	Clusters []Cluster
	Observe  observe
	Clients  []Client
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

	return yml, yml.validate()
}

// FindByFullMethod returns all routes match fully qualified the gRPC service method name.
func (y *Yaml) FindByFullMethod(name string) []Route {
	routes := make([]Route, 0)

	for _, route := range y.Routes {
		if equal := route.Method.Equal; equal != nil {
			if *equal == name {
				routes = append(routes, route)
			}
		}
		if prefix := route.Method.Prefix; prefix != nil {
			if strings.HasPrefix(name, *prefix) {
				routes = append(routes, route)
			}
		}
	}

	return routes
}

func (y *Yaml) validate() error {
	errors := make([]error, 0)
	for _, route := range y.Routes {
		err := route.validate()
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) != 0 {
		s := ""
		for _, err := range errors {
			if len(s) != 0 {
				s += " ,"
			}
			s += fmt.Sprintf("%s", err.Error())
		}
		return fmt.Errorf("yaml validation %d error(s): %s", len(errors), s)
	}
	return nil
}

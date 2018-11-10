package yaml

import (
	"strings"
)

type Routes struct {
	Routes []Route
}

type Route struct {
	Method  Method
	Cluster RouteCluster
}

type Method struct {
	Prefix string
}

type RouteCluster struct {
	Name string
}

func (r *Routes) FindByFullMethod(name string) []Route {
	routes := make([]Route, 0)

	for _, route := range r.Routes {
		if strings.HasPrefix(name, route.Method.Prefix) {
			routes = append(routes, route)
		}
	}

	return routes
}

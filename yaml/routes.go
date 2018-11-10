package yaml

import (
	"strings"
)

// Routes represents a configuration of gRPC routing roules.
type Routes struct {
	Routes []Route
}

// Route represents a configuration of a single gRPC routing roule.
type Route struct {
	Method  Method
	Cluster RouteCluster
}

// Method represents a gRPC service method matching rule.
type Method struct {
	Prefix string
}

// RouteCluster represents a configuration of a binding between the method and the cluster.
type RouteCluster struct {
	Name string
}

// FindByFullMethod returns all routes match fully qualified the gRPC service method name.
func (r *Routes) FindByFullMethod(name string) []Route {
	routes := make([]Route, 0)

	for _, route := range r.Routes {
		if strings.HasPrefix(name, route.Method.Prefix) {
			routes = append(routes, route)
		}
	}

	return routes
}

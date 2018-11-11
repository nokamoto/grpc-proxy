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
	Method struct {
		Prefix string
	}
	Cluster struct {
		Name string
	}
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

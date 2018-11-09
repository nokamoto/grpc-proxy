package main

import (
	"strings"
)

type yamlRoutes struct {
	Routes []yamlRoute
}

type yamlRoute struct {
	Method  yamlMethod
	Cluster yamlRouteCluster
}

type yamlMethod struct {
	Prefix string
}

type yamlRouteCluster struct {
	Name string
}

func (r *yamlRoutes) findByFullMethod(name string) []yamlRoute {
	routes := make([]yamlRoute, 0)

	for _, route := range r.Routes {
		if strings.HasPrefix(name, route.Method.Prefix) {
			routes = append(routes, route)
		}
	}

	return routes
}

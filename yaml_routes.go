package main

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

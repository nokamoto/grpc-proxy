package yaml

import (
	"fmt"
)

// Route represents a configuration of a single gRPC routing roule.
type Route struct {
	Method  routeMethod
	Cluster routeCluster
	Observe routeObserve
}

type routeMethod struct {
	Prefix *string `yaml:",omitempty"`
	Equal  *string `yaml:",omitempty"`
}

type routeCluster struct {
	Name string
}

type routeObserve struct {
	Log  routeObserveLog
	Prom routeObserveProm
}

type routeObserveLog struct {
	Name *string
}

type routeObserveProm struct {
	Name *string
}

func (r Route) validate() error {
	if r.Method.Equal == nil && r.Method.Prefix == nil {
		return fmt.Errorf("routes.method: equal or prefix must be defined")
	}
	if r.Method.Equal != nil && r.Method.Prefix != nil {
		return fmt.Errorf("routes.method: equal and prefix must not be defined: %s %s", *r.Method.Equal, *r.Method.Prefix)
	}
	return nil
}

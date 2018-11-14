package yaml

import (
	"fmt"
)

// Route represents a configuration of a single gRPC routing roule.
type Route struct {
	Method struct {
		Prefix *string `yaml:",omitempty"`
		Equal  *string `yaml:",omitempty"`
	}
	Cluster struct {
		Name string
	}
	Observe struct {
		Log struct {
			Name *string
		}
		Prom struct {
			Name *string
		}
	}
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

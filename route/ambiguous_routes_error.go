package route

import (
	"fmt"
	"github.com/nokamoto/grpc-proxy/yaml"
	yml "gopkg.in/yaml.v2"
)

type ambiguousRoutesError struct {
	method     string
	candidates []yaml.Route
}

func (e *ambiguousRoutesError) Error() string {
	b, err := yml.Marshal(e.candidates)
	if err != nil {
		return fmt.Sprintf("%s has ambiguous routes", e.method)
	}
	return fmt.Sprintf("%s has ambiguous routes:\n%v", e.method, string(b))
}

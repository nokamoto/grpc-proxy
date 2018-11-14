package route

import (
	"fmt"
)

type missingRoutesError struct {
	method string
}

func (e *missingRoutesError) Error() string {
	return fmt.Sprintf("%s has no routes", e.method)
}

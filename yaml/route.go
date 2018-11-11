package yaml

// Route represents a configuration of a single gRPC routing roule.
type Route struct {
	Method struct {
		Prefix string
	}
	Cluster struct {
		Name string
	}
	Observe struct {
		Log struct {
			Name string
		}
	}
}

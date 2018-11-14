package route

import (
	"github.com/nokamoto/grpc-proxy/descriptor"
	"github.com/nokamoto/grpc-proxy/yaml"
	"testing"
)

func TestNewRoutes_ping_method_prefix(t *testing.T) {
	_, afterEach, err := testRoutes(t, "../testdata/yaml/ping.yaml")
	defer afterEach()

	if err != nil {
		t.Fatal(err)
	}
}

func TestNewRoutes_ping_method_equal(t *testing.T) {
	_, afterEach, err := testRoutes(t, "../testdata/yaml/ping_method_equal.yaml")
	defer afterEach()

	if err != nil {
		t.Fatal(err)
	}
}

func TestNewRoutes_ping_method_missing(t *testing.T) {
	_, afterEach, err := testRoutes(t, "../testdata/yaml/ping_method_missing.yaml")
	defer afterEach()

	if err == nil {
		t.Fatal()
	}

	_, ok := err.(*missingRoutesError)
	if !ok {
		t.Fatal(err)
	}
}

func TestNewRoutes_ping_method_ambiguous(t *testing.T) {
	_, afterEach, err := testRoutes(t, "../testdata/yaml/ping_method_ambiguous.yaml")
	defer afterEach()

	if err == nil {
		t.Fatal()
	}

	_, ok := err.(*ambiguousRoutesError)
	if !ok {
		t.Fatal(err)
	}
}

func testRoutes(t *testing.T, y string) (*Routes, func(), error) {
	t.Helper()

	pb, err := descriptor.NewDescriptor("../testdata/protobuf/ping/ping.pb")
	if err != nil {
		t.Fatal(err)
	}

	yml, err := yaml.NewYaml(y)
	if err != nil {
		t.Fatal(err)
	}

	r, err := NewRoutes(pb, yml)

	f := func() {
		if r != nil {
			r.Destroy()
		}
	}

	return r, f, err
}

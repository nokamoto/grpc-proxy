package yaml

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// NewYaml returns routes and clusters configurations read from the yaml file.
func NewYaml(y string) (*Routes, *Clusters, *Observe, error) {
	bytes, err := ioutil.ReadFile(y)
	if err != nil {
		return nil, nil, nil, err
	}

	r := &Routes{}
	if err := yaml.Unmarshal(bytes, r); err != nil {
		return nil, nil, nil, err
	}

	c := &Clusters{}
	if err := yaml.Unmarshal(bytes, c); err != nil {
		return nil, nil, nil, err
	}

	o := &Observe{}
	if err := yaml.Unmarshal(bytes, o); err != nil {
		return nil, nil, nil, err
	}

	return r, c, o, nil
}

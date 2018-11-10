package yaml

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func NewYaml(y string) (*Routes, *Clusters, error) {
	bytes, err := ioutil.ReadFile(y)
	if err != nil {
		return nil, nil, err
	}

	r := &Routes{}
	if err := yaml.Unmarshal(bytes, r); err != nil {
		return nil, nil, err
	}

	c := &Clusters{}
	if err := yaml.Unmarshal(bytes, c); err != nil {
		return nil, nil, err
	}

	return r, c, nil
}

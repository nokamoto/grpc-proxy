package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func newYaml(y string) (*yamlRoutes, *yamlClusters, error) {
	bytes, err := ioutil.ReadFile(y)
	if err != nil {
		return nil, nil, err
	}

	r := &yamlRoutes{}
	if err := yaml.Unmarshal(bytes, r); err != nil {
		return nil, nil, err
	}

	c := &yamlClusters{}
	if err := yaml.Unmarshal(bytes, c); err != nil {
		return nil, nil, err
	}

	return r, c, nil
}

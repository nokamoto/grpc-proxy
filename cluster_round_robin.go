package main

import (
	"fmt"
	"golang.org/x/net/context"
	"sync"
)

type clusterRoundRobin struct {
	mu      sync.Mutex
	proxies []*proxy
	next    int
}

func newClusterRoundRobin(c yamlCluster) (*clusterRoundRobin, error) {
	proxies := make([]*proxy, 0)

	for _, address := range c.RoundRobin {
		proxy, err := newProxy(address)
		if err != nil {
			return nil, err
		}

		proxies = append(proxies, proxy)
	}

	if len(proxies) == 0 {
		return nil, fmt.Errorf("cluster %s empty round robin", c.Name)
	}

	return &clusterRoundRobin{proxies: proxies}, nil
}

func (c *clusterRoundRobin) invokeUnary(ctx context.Context, m *message, method string) (*message, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.next = (c.next + 1) % len(c.proxies)

	return c.proxies[c.next].invokeUnary(ctx, m, method)
}

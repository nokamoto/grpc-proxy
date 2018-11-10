package cluster

import (
	"fmt"
	"github.com/nokamoto/grpc-proxy/codec"
	"github.com/nokamoto/grpc-proxy/server"
	"github.com/nokamoto/grpc-proxy/yaml"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"sync"
)

type roundRobin struct {
	mu      sync.Mutex
	proxies []*proxy
	next    int
}

// NewRoundRobin returns Cluster with round robin load balancing.
func NewRoundRobin(c yaml.Cluster) (Cluster, error) {
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

	return &roundRobin{proxies: proxies}, nil
}

func (c *roundRobin) nextProxy() *proxy {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.next = (c.next + 1) % len(c.proxies)

	return c.proxies[c.next]
}

func (c *roundRobin) InvokeUnary(ctx context.Context, m *codec.RawMessage, method string) (*codec.RawMessage, error) {
	return c.nextProxy().invokeUnary(ctx, m, method)
}

func (c *roundRobin) InvokeStreamC(stream server.RawServerStreamC, desc *grpc.StreamDesc, method string) error {
	return c.nextProxy().invokeStreamC(stream, desc, method)
}

func (c *roundRobin) InvokeStreamS(stream server.RawServerStreamS, desc *grpc.StreamDesc, method string) error {
	return c.nextProxy().invokeStreamS(stream, desc, method)
}

func (c *roundRobin) InvokeStreamB(stream server.RawServerStreamB, desc *grpc.StreamDesc, method string) error {
	return c.nextProxy().invokeStreamB(stream, desc, method)
}

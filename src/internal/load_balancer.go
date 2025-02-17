package loadbalancer

import (
	"sync"
	"time"
)

// represents a backend server
type Server struct {
	URL   string
	Alive bool
	mux   sync.RWMutex
	LastChecked time.Time
	FailCount int
}

// configs for load balancer
type loadBalancer struct {
	servers []*Server
	mux sync.RWMutex
	// roundrobin counter
	current int
}

// init new load balancer
func NewLoadBalancer() *loadBalancer {
	return &loadBalancer{
		servers: make([]*Server, 0),
		current: 0,
	}
}

// used to add a new server to the load balancer
func (lb *loadBalancer) AddServer(url string) {
	lb.mux.Lock()
	defer lb.mux.Unlock()

	server := &Server{
		URL: url,
		Alive: true,
	}
	lb.servers = append(lb.servers, server)
}

// using a round-robin algorithm for nextServer which returns the next available server
func (lb *loadBalancer) NextServer() *Server {
	lb.mux.Lock()
	defer lb.mux.Unlock()

	// loopin through the servers to find a healthy one
	for i := 0; i < len(lb.servers); i++ {
		lb.current = (lb.current + 1) % len(lb.servers)
		if lb.servers[lb.current].Alive {
			return lb.servers[lb.current]
		}
	}
	return nil
}

func (lb *loadBalancer) GetServers() []*Server {
	lb.mux.RLock()
	defer lb.mux.RUnlock()

	return lb.servers
}
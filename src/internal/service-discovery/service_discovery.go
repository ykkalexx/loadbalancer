package discovery

import (
	"sync"
	"time"
)

type ServiceRegistry struct {
    services map[string][]string  // service name -> list of endpoints
    mu       sync.RWMutex
    ttl      time.Duration
}

func NewServiceRegistry() *ServiceRegistry {
    return &ServiceRegistry{
        services: make(map[string][]string),
        ttl:      time.Minute * 5,
    }
}
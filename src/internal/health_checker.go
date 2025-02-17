package loadbalancer

import (
	"net/http"
	"time"
)

const (
	healthCheckInterval = 20 * time.Second
	maxFailCount = 3
)

// i use this to check the health of a server
// if a server is not healthy, it will not be used for load balancing
func (s *Server) CheckHealth() {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(s.URL)
	s.mux.Lock()
	defer s.mux.Unlock()

	if err != nil {
		s.FailCount++
		if s.FailCount >= maxFailCount {
			s.Alive = false
		}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		s.Alive = true
		s.FailCount = 0
	} else {
		s.FailCount++
		if s.FailCount >= maxFailCount {
			s.Alive = false
		}
	}
	s.LastChecked = time.Now()
}

// start health check for all servers
func (lb *loadBalancer) StartHealthCheck() {
    ticker := time.NewTicker(healthCheckInterval)
    go func() {
        for {
            select {
            case <-ticker.C:
                lb.mux.RLock()
                for _, server := range lb.servers {
                    go server.CheckHealth()
                }
                lb.mux.RUnlock()
            }
        }
    }()
}

func (lb *loadBalancer) GetHealthyServerCount() int {
	lb.mux.RLock()
	defer lb.mux.RUnlock()

	count := 0
	for _, server := range lb.servers {
		if server.Alive {
			count++
		}
	}
	return count
}
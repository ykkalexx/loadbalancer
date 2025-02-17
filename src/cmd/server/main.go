package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	loadbalancer "github.com/ykkalexx/load-balancer/internal"
	"github.com/ykkalexx/load-balancer/internal/config"
	"github.com/ykkalexx/load-balancer/internal/ratelimiter"
)

func main() {
	config, err := config.LoadConfig("C:\\Projects\\DistributedLoadBalancer\\config.json")
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}


	// init load balancer
    lb := loadbalancer.NewLoadBalancer()
    metrics := loadbalancer.NewMetrics()
    ratelimiter := ratelimiter.NewRateLimiter(
        config.RateLimit.RequestsPerSecond,
        time.Second,
    )

	// adding servers 
    for _, server := range config.Servers {
        lb.AddServerWithWeight(server.URL, server.Weight)
    }


	// start health checks on the servers
	lb.StartHealthCheck()

    // creating a proxy handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// rate limiter
		if !ratelimiter.Allow(r.RemoteAddr) {
            metrics.RecordRequest(0, false)
            http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
            return
        }

		server := lb.NextServer()
		start := time.Now()
		if server == nil {
			metrics.RecordRequest(time.Since(start), false)
			log.Printf("No available servers")
			http.Error(w, "No available servers", http.StatusServiceUnavailable)
			return
		}
	
		log.Printf("Routing request to: %s", server.URL)
		targetURL, err := url.Parse(server.URL)
		if err != nil {
			log.Printf("Error parsing URL: %v", err)
			http.Error(w, "Invalid backend server URL", http.StatusInternalServerError)
			return
		}
	
		metrics.RecordRequest(time.Since(start), true)
		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		proxy.ServeHTTP(w, r)
	})

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
        data := map[string]interface{}{
            "total_requests":       metrics.TotalRequests,
            "successful_requests": metrics.SuccesfullRequests,
            "failed_requests":     metrics.FailedRequests,
            "server_count":        len(lb.GetServers()),
            "healthy_servers":     lb.GetHealthyServerCount(),
        }
        json.NewEncoder(w).Encode(data)
    })


	// starting the load balancer
	log.Printf("Load balancer started at :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
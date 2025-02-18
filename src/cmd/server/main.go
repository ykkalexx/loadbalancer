package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	loadbalancer "github.com/ykkalexx/load-balancer/internal"
	circuitbreaker "github.com/ykkalexx/load-balancer/internal/circuit-breaker"
	"github.com/ykkalexx/load-balancer/internal/cluster"
	"github.com/ykkalexx/load-balancer/internal/config"
	"github.com/ykkalexx/load-balancer/internal/ratelimiter"
	retry "github.com/ykkalexx/load-balancer/internal/retryMechanicsm"
)

func main() {
    configPath := os.Args[2] 
    config, err := config.LoadConfig(configPath)
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
    cb := circuitbreaker.NewCircuitBreaker(5, 10*time.Second)
    retryPolicy := retry.NewRetryPolicy()
    clusterManager := cluster.NewClusterManager(
        config.Cluster.NodeID,
        config.Cluster.IsPrimary,
        config.Cluster.PeerNodes,
    )
    clusterManager.Start()

    // adding servers 
    for _, server := range config.Servers {
        lb.AddServerWithWeight(server.URL, server.Weight)
    }

    // start health checks on the servers
    lb.StartHealthCheck()

    // Create a new ServeMux for routing
    mux := http.NewServeMux()

    // Add metrics endpoint
    mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        data := map[string]interface{}{
            "total_requests":      metrics.TotalRequests,
            "successful_requests": metrics.SuccesfullRequests,
            "failed_requests":     metrics.FailedRequests,
            "server_count":        len(lb.GetServers()),
            "healthy_servers":     lb.GetHealthyServerCount(),
            "cluster_info": map[string]interface{}{
                "node_id":    config.Cluster.NodeID,
                "is_primary": config.Cluster.IsPrimary,
                "healthy_nodes": clusterManager.GetHealthyNodes(),
            },
        }
        json.NewEncoder(w).Encode(data)
    })

    // Create the main proxy handler
    proxyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // rate limiter
        if !ratelimiter.Allow(r.RemoteAddr) {
            metrics.RecordRequest(0, false)
            http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
            return
        }

        // circuit breaker check
        if !cb.IsAllowed() {
            metrics.RecordRequest(0, false)
            http.Error(w, "Service temporarily unavailable", http.StatusServiceUnavailable)
            return
        }

        var lastErr error
        start := time.Now()

        // Retry logic
        for attempt := 0; attempt < retryPolicy.MaxAttempts; attempt++ {
            server := lb.NextServer()
            if server == nil {
                continue
            }

            targetURL, err := url.Parse(server.URL)
            if err != nil {
                lastErr = err
                continue
            }

			proxy := httputil.NewSingleHostReverseProxy(targetURL)
			proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
				lastErr = err
				cb.RecordFailure()
				metrics.RecordRequest(time.Since(start), false)
			}			

			log.Printf("Routing request to: %s", server.URL)
			proxy.ServeHTTP(w, r)
			cb.RecordSuccess() 
			metrics.RecordRequest(time.Since(start), true)
			return
        }

        // if retries failed
        metrics.RecordRequest(time.Since(start), false)
        errMsg := "All attempts failed"
        if lastErr != nil {
            errMsg = fmt.Sprintf("All attempts failed: %v", lastErr)
        }
        http.Error(w, errMsg, http.StatusServiceUnavailable)
    })

    // the main handler
    mux.Handle("/", proxyHandler)

    // Add a health check endpoint
    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]interface{}{
            "status":    "healthy",
            "node_id":   config.Cluster.NodeID,
            "is_primary": config.Cluster.IsPrimary,
            "timestamp": time.Now(),
        })
    })

    // starting the load balancer
    addr := fmt.Sprintf(":%d", config.Port)
    log.Printf("Load balancer node %s starting at %s", config.Cluster.NodeID, addr)
    if err := http.ListenAndServe(addr, mux); err != nil {
        log.Fatal(err)
    }
}
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	loadbalancer "github.com/ykkalexx/load-balancer/internal"
	circuitbreaker "github.com/ykkalexx/load-balancer/internal/circuit-breaker"
	"github.com/ykkalexx/load-balancer/internal/config"
	"github.com/ykkalexx/load-balancer/internal/ratelimiter"
	retry "github.com/ykkalexx/load-balancer/internal/retryMechanicsm"
)

func main() {
    config, err := config.LoadConfig("C:/Projects/DistributedLoadBalancer/config.json")
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

    mux.HandleFunc("/cluster/join", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }
        // Handle node join
    })

    mux.HandleFunc("/cluster/heartbeat", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }
        // Handle heartbeat
    })

    // starting the load balancer
    log.Printf("Load balancer started at :8080")
    if err := http.ListenAndServe(":8080", mux); err != nil {
        log.Fatal(err)
    }
}
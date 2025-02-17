package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	loadbalancer "github.com/ykkalexx/load-balancer/internal"
)

func main() {
	// init load balancer
	lb := loadbalancer.NewLoadBalancer()

	// adding servers 
	lb.AddServer("http://localhost:5001")
	lb.AddServer("http://localhost:5002")
	lb.AddServer("http://localhost:5003")

	// start health checks on the servers
	lb.StartHealthCheck()

    // creating a proxy handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server := lb.NextServer()
		if server == nil {
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
	
		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		proxy.ServeHTTP(w, r)
	})

	// starting the load balancer
	log.Printf("Load balancer started at :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
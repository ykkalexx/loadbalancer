package ratelimiter

import (
	"sync"
	"time"
)

// this implements a sliding window rate limiting
type RateLimiter struct {
    requests map[string][]time.Time
    limit    int           // maximum requests per window
    window   time.Duration // time window for rate limiting
    mu       sync.RWMutex
}

// creates a new rate limiter instance
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
    return &RateLimiter{
        requests: make(map[string][]time.Time),
        limit:    limit,
        window:   window,
    }
}

// checks if a request from an IP should be allowed
func (rl *RateLimiter) Allow(ip string) bool {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    now := time.Now()
    windowStart := now.Add(-rl.window)

    // Clean up old requests
    if requests, exists := rl.requests[ip]; exists {
        var recent []time.Time
        for _, t := range requests {
            if t.After(windowStart) {
                recent = append(recent, t)
            }
        }
        rl.requests[ip] = recent
    }

    // Check if limit is exceeded
    if len(rl.requests[ip]) >= rl.limit {
        return false
    }

    // Add new request
    rl.requests[ip] = append(rl.requests[ip], now)
    return true
}

// returns the number of requests for an IP in the current window
func (rl *RateLimiter) GetCurrentRequests(ip string) int {
    rl.mu.RLock()
    defer rl.mu.RUnlock()
    return len(rl.requests[ip])
}
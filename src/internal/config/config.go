package config

import (
	"encoding/json"
	"os"
	"time"
)

type ServerConfig struct {
    URL    string `json:"url"`
    Weight int    `json:"weight"`
}

type RateLimitConfig struct {
    RequestsPerSecond int `json:"requests_per_second"`
    BurstSize        int `json:"burst_size"`
}

type HealthCheckConfig struct {
    IntervalSeconds int `json:"interval_seconds"`
    TimeoutSeconds  int `json:"timeout_seconds"`
    MaxFailures     int `json:"max_failures"`
}

type ClusterConfig struct {
    NodeID      string   `json:"node_id"`
    Port        int      `json:"port"`
    PeerNodes   []string `json:"peer_nodes"`
    IsPrimary   bool     `json:"is_primary"`
}



type Config struct {
    Port        int              `json:"port"`
    Servers     []ServerConfig   `json:"servers"`
    RateLimit   RateLimitConfig  `json:"rate_limit"`
    HealthCheck HealthCheckConfig `json:"health_check"`
    Cluster ClusterConfig `json:"cluster"`
}

func LoadConfig(filename string) (*Config, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var config Config
    if err := json.NewDecoder(file).Decode(&config); err != nil {
        return nil, err
    }

    // Set defaults if not specified
    if config.Port == 0 {
        config.Port = 8080
    }
    if config.HealthCheck.IntervalSeconds == 0 {
        config.HealthCheck.IntervalSeconds = 20
    }
    if config.HealthCheck.TimeoutSeconds == 0 {
        config.HealthCheck.TimeoutSeconds = 5
    }
    if config.HealthCheck.MaxFailures == 0 {
        config.HealthCheck.MaxFailures = 3
    }

    return &config, nil
}

// GetHealthCheckDuration returns the health check interval as time.Duration
func (c *Config) GetHealthCheckDuration() time.Duration {
    return time.Duration(c.HealthCheck.IntervalSeconds) * time.Second
}

// GetTimeoutDuration returns the timeout as time.Duration
func (c *Config) GetTimeoutDuration() time.Duration {
    return time.Duration(c.HealthCheck.TimeoutSeconds) * time.Second
}
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
	BurstSize         int `json:"burst_size"`
}

type HealthCheckConfig struct {
	Interval    time.Duration `json:"interval"`
	Timeout     time.Duration `json:"timeout"`
	MaxFailures int           `json:"max_failures"`
}

type Config struct {
	Port        int               `json:"port"`
	Servers     []ServerConfig    `json:"servers"`
	RateLimit   RateLimitConfig   `json:"rate_limit"`
	HealthCheck HealthCheckConfig `json:"health_check"`
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
    if config.RateLimit.RequestsPerSecond == 0 {
        config.RateLimit.RequestsPerSecond = 100
    }
    if config.RateLimit.BurstSize == 0 {
        config.RateLimit.BurstSize = 20
    }
    if config.HealthCheck.Interval == 0 {
        config.HealthCheck.Interval = 20 * time.Second
    }
    if config.HealthCheck.Timeout == 0 {
        config.HealthCheck.Timeout = 5 * time.Second
    }
    if config.HealthCheck.MaxFailures == 0 {
        config.HealthCheck.MaxFailures = 3
    }

    return &config, nil
}
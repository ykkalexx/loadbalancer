package retry

import (
	"math"
	"time"
)

type RetryPolicy struct {
    MaxAttempts     int
    InitialInterval time.Duration
    MaxInterval     time.Duration
    Multiplier      float64
}

func NewRetryPolicy() *RetryPolicy {
    return &RetryPolicy{
        MaxAttempts:     3,
        InitialInterval: 100 * time.Millisecond,
        MaxInterval:     2 * time.Second,
        Multiplier:      2.0,
    }
}

func (rp *RetryPolicy) GetNextInterval(attempt int) time.Duration {
    if attempt <= 0 {
        return rp.InitialInterval
    }

    interval := float64(rp.InitialInterval) * math.Pow(rp.Multiplier, float64(attempt-1))
    if interval > float64(rp.MaxInterval) {
        return rp.MaxInterval
    }
    return time.Duration(interval)
}
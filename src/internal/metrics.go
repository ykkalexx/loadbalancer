package loadbalancer

import (
	"sync"
	"time"
)

type Metrics struct {
	TotalRequests      int64
	SuccesfullRequests int64
	FailedRequests     int64
	ResponseTimes      []time.Duration
	mux sync.RWMutex	
}

// NewMetrics creates a new Metrics object
func NewMetrics() *Metrics {
	return &Metrics{
		ResponseTimes: make([]time.Duration, 0),
	}
}

// function used to record the request metrics
func (m *Metrics) RecordRequest(duration time.Duration, success bool) {
	m.mux.Lock()
	defer m.mux.Unlock()

	m.TotalRequests++
	if success {
		m.SuccesfullRequests++
	} else {
		m.FailedRequests++
	}
	m.ResponseTimes = append(m.ResponseTimes, duration)
}


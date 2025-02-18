package cluster

import (
	"sync"
	"time"
)

type Node struct {
    ID        string
    Address   string
    LastSeen  time.Time
    IsLeader  bool
}

type Cluster struct {
    nodes    map[string]*Node
    mu       sync.RWMutex
    selfID   string
}
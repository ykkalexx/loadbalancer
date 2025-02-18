package cluster

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type NodeState int

const (
    StateHealthy NodeState = iota
    StateUnhealthy
    StateDown
)

type Node struct {
    ID       string
    Address  string
    State    NodeState
    LastSeen time.Time
}

type ClusterManager struct {
    nodeID    string
    nodes     map[string]*Node
    mu        sync.RWMutex
    isPrimary bool
}

func NewClusterManager(nodeID string, isPrimary bool, peerAddresses []string) *ClusterManager {
    cm := &ClusterManager{
        nodeID:    nodeID,
        nodes:     make(map[string]*Node),
        isPrimary: isPrimary,
    }

    // Register peer nodes
    for _, addr := range peerAddresses {
        cm.nodes[addr] = &Node{
            Address:  addr,
            State:   StateHealthy,
            LastSeen: time.Now(),
        }
    }

    return cm
}

func (cm *ClusterManager) Start() {
    go cm.healthCheckLoop()
    log.Printf("Starting cluster manager for node %s (Primary: %v)", cm.nodeID, cm.isPrimary)
}

func (cm *ClusterManager) healthCheckLoop() {
    ticker := time.NewTicker(5 * time.Second)
    for range ticker.C {
        cm.checkPeerHealth()
    }
}

func (cm *ClusterManager) checkPeerHealth() {
    cm.mu.Lock()
    defer cm.mu.Unlock()

    for _, node := range cm.nodes {
        resp, err := http.Get(fmt.Sprintf("%s/health", node.Address))
        if err != nil || resp.StatusCode != http.StatusOK {
            node.State = StateUnhealthy
            log.Printf("Node %s is unhealthy", node.Address)
        } else {
            node.State = StateHealthy
            node.LastSeen = time.Now()
        }
    }
}

func (cm *ClusterManager) GetHealthyNodes() []string {
    cm.mu.RLock()
    defer cm.mu.RUnlock()

    var healthy []string
    for _, node := range cm.nodes {
        if node.State == StateHealthy {
            healthy = append(healthy, node.Address)
        }
    }
    return healthy
}

func (cm *ClusterManager) IsPrimary() bool {
    return cm.isPrimary
}
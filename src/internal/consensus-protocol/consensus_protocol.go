package consensus

import (
	"sync"
)

type NodeState int

const (
    Follower NodeState = iota
    Candidate
    Leader
)

type RaftNode struct {
    state       NodeState
    currentTerm int
    votedFor    string
    mu          sync.Mutex
    peers       []string
    leaderId    string
}
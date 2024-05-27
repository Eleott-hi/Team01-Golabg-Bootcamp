package raft

import (
	"math/rand"
	"sync"
	"time"
)

type State int

const (
	Follower State = iota
	Candidate
	Leader
)

type LogEntry struct {
	Term    int
	Command interface{}
}

type Raft struct {
	mu          sync.Mutex
	peers       []string
	me          int
	state       State
	currentTerm int
	votedFor    int
	log         []LogEntry

	commitIndex int
	lastApplied int

	nextIndex  []int
	matchIndex []int

	electionTimeout  time.Duration
	heartbeatTimeout time.Duration

	electionTimer  *time.Timer
	heartbeatTimer *time.Timer
}

func NewRaft(peers []string, me int) *Raft {
	r := &Raft{
		peers:            peers,
		me:               me,
		state:            Follower,
		votedFor:         -1,
		electionTimeout:  time.Duration(rand.Intn(150)+150) * time.Millisecond,
		heartbeatTimeout: 50 * time.Millisecond,
	}

	r.electionTimer = time.NewTimer(r.electionTimeout)
	r.heartbeatTimer = time.NewTimer(r.heartbeatTimeout)

	return r
}

func (r *Raft) resetElectionTimer() {
	r.electionTimer.Stop()
	r.electionTimer.Reset(r.electionTimeout)
}

func (r *Raft) resetHeartbeatTimer() {
	r.heartbeatTimer.Stop()
	r.heartbeatTimer.Reset(r.heartbeatTimeout)
}

func (r *Raft) startElection() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.state = Candidate
	r.currentTerm++
	r.votedFor = r.me

	votesRecieved := 1

	for _, peer := range r.peers {
		if peer == r.peers[r.me] {
			continue
		}

		go func(peer string) {
			voteGranted := r.requestVote(peer)

			if voteGranted {
				r.mu.Lock()
				votesRecieved++
				if votesRecieved > len(r.peers)/2 && r.state == Candidate {
					r.state = Leader
					r.resetHeartbeatTimer()
				}
				r.mu.Unlock()
			}
		}(peer)
	}

	r.resetElectionTimer()
}

func (r *Raft) requestVote(peer string) bool {
	// TODO: request
	return true
}

func (r *Raft) run() {
	for {
		switch r.state {
		case Follower:
			select {
			case <-r.electionTimer.C:
				r.startElection()
			}
		case Candidate:
			select {
			case <-r.electionTimer.C:
				r.startElection()
			}
		case Leader:
			select {
			case <-r.heartbeatTimer.C:
				r.sendHeartbeats()
				r.resetHeartbeatTimer()
			}
		}
	}
}

func (r *Raft) sendHeartbeats() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, peer := range r.peers {
		if peer == r.peers[r.me] {
			continue
		}

		go func(peer string) {
			r.appendEntries(peer)
		}(peer)
	}
}

func (r *Raft) appendEntries(peer string) {
	// TODO
}

func manamain() {
	peers := []string{"server1", "server2", "server3"}
	r := NewRaft(peers, 0)
	go r.run()

	select {}
}

package main

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
	id          int
	state       State
	currentTerm int
	votedFor    int
	log         []LogEntry
	commitIndex int
	lastApplied int

	// Cluster Information
	peers            []string
	applyCh          chan ApplyMsg
	electionTimeout  time.Duration
	heartbeatTimeout time.Duration

	// Volatile state on leaders
	nextIndex  []int
	matchIndex []int
}

type ApplyMsg struct {
	Index   int
	Command interface{}
}

type RequestVoteArgs struct {
	Term         int
	CandidateId  int
	LastLogIndex int
	LastLogTerm  int
}

type RequestVoteReply struct {
	Term        int
	VoteGranted bool
}

type AppendEntriesArgs struct {
	Term         int
	LeaderId     int
	PrevLogIndex int
	PrevLogTerm  int
	Entries      []LogEntry
	LeaderCommit int
}

type AppendEntriesReply struct {
	Term    int
	Success bool
}

func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	if args.Term > rf.currentTerm {
		rf.currentTerm = args.Term
		rf.votedFor = -1
		rf.state = Follower
	}

	if rf.votedFor == -1 || rf.votedFor == args.CandidateId {
		if args.LastLogTerm > rf.lastLogTerm() || (args.LastLogTerm == rf.lastLogTerm() && args.LastLogIndex >= len(rf.log)-1) {
			reply.VoteGranted = true
			rf.votedFor = args.CandidateId
			rf.electionTimeout = rf.electionTimeoutDuration()
		}
	}
	reply.Term = rf.currentTerm
}

func (rf *Raft) lastLogTerm() int {
	if len(rf.log) == 0 {
		return -1
	}
	return rf.log[len(rf.log)-1].Term
}

func (rf *Raft) run() {
	for {
		switch rf.state {
		case Follower:
			rf.runFollower()
		case Candidate:
			rf.runCandidate()
		case Leader:
			rf.runLeader()
		}
	}
}

func (rf *Raft) runFollower() {
	select {
	case <-time.After(rf.electionTimeoutDuration()):
		rf.state = Candidate
	case <-rf.applyCh:
		// Handle apply channel message
	}
}

func (rf *Raft) electionTimeoutDuration() time.Duration {
	return time.Duration(150+rand.Intn(150)) * time.Millisecond
}

func (rf *Raft) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	if args.Term < rf.currentTerm {
		reply.Term = rf.currentTerm
		reply.Success = false
		return
	}

	if len(rf.log) <= args.PrevLogIndex || rf.log[args.PrevLogIndex].Term != args.PrevLogTerm {
		reply.Success = false
		return
	}

	rf.log = append(rf.log[:args.PrevLogIndex+1], args.Entries...)
	if args.LeaderCommit > rf.commitIndex {
		rf.commitIndex = min(args.LeaderCommit, len(rf.log)-1)
	}

	reply.Success = true
	reply.Term = rf.currentTerm
}

func (rf *Raft) runLeader() {
	for {
		for _, peer := range rf.peers {
			go func(peer string) {
				args := &AppendEntriesArgs{
					Term:         rf.currentTerm,
					LeaderId:     rf.id,
					PrevLogIndex: len(rf.log) - 1,
					PrevLogTerm:  rf.lastLogTerm(),
					Entries:      []LogEntry{},
					LeaderCommit: rf.commitIndex,
				}
				reply := &AppendEntriesReply{}
				rf.sendAppendEntries(peer, args, reply)
			}(peer)
		}
		time.Sleep(rf.heartbeatTimeout)
	}
}

func (rf *Raft) runCandidate() {
	rf.mu.Lock()
	rf.currentTerm++
	rf.votedFor = rf.id
	rf.state = Candidate
	rf.mu.Unlock()

	votes := 1
	voteCh := make(chan bool)
	resetCh := make(chan bool)

	for _, peer := range rf.peers {
		if peer != rf.peers[rf.id] {
			go func(peer string) {
				args := &RequestVoteArgs{
					Term:         rf.currentTerm,
					CandidateId:  rf.id,
					LastLogIndex: len(rf.log) - 1,
					LastLogTerm:  rf.lastLogTerm(),
				}
				reply := &RequestVoteReply{}
				rf.sendRequestVote(peer, args, reply)

				rf.mu.Lock()
				defer rf.mu.Unlock()
				if reply.Term > rf.currentTerm {
					rf.currentTerm = reply.Term
					rf.state = Follower
					rf.votedFor = -1
					resetCh <- true
					return
				}

				if reply.VoteGranted {
					votes++
					if votes > len(rf.peers)/2 {
						rf.state = Leader
						voteCh <- true
						return
					}
				}
			}(peer)
		}
	}

	timeout := time.After(rf.electionTimeoutDuration())

	for {
		select {
		case <-timeout:
			rf.mu.Lock()
			rf.state = Follower
			rf.mu.Unlock()
			return
		case <-voteCh:
			return
		case <-resetCh:
			return
		}
	}
}

func (rf *Raft) sendRequestVote(peer string, args *RequestVoteArgs, reply *RequestVoteReply) {
	// Simulate RPC call to peer (this should actually be an RPC call in a real implementation)
	// For simplicity, let's assume the function directly modifies the reply
	// In practice, you would use a library such as gRPC or net/rpc to handle this.
}

func (rf *Raft) sendAppendEntries(peer string, args *AppendEntriesArgs, reply *AppendEntriesReply) {
	// Simulate RPC call to peer
}

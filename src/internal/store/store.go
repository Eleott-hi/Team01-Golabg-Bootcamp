package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/hashicorp/raft"
)

const (
	retainSnapshotCount = 2
	raftTimeout         = 10 * time.Second
)

type command struct {
	Op    string `json:"op,omitempty"`
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type Store struct {
	RaftDir  string
	RaftBind string

	mu      sync.Mutex
	storage map[string]string
	raft    *raft.Raft
	logger  *log.Logger
}

func New(raftDir, raftAddr string) *Store {
	return &Store{
		RaftDir:  raftDir,
		RaftBind: raftAddr,
		storage:  make(map[string]string),
		logger:   log.New(os.Stderr, "[store] ", log.LstdFlags),
	}
}

func (s *Store) Open(enableSingle bool, localID string) error {
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(localID)

	addr, err := net.ResolveTCPAddr("tcp", s.RaftBind)
	if err != nil {
		return err
	}

	transport, err := raft.NewTCPTransport(s.RaftBind, addr, 3, 10*time.Second, os.Stderr)
	if err != nil {
		return err
	}

	snapshots, err := raft.NewFileSnapshotStore(s.RaftDir, retainSnapshotCount, os.Stderr)
	if err != nil {
		return err
	}

	logStore := raft.NewInmemStore()
	stableStore := raft.NewInmemStore()

	r, err := raft.NewRaft(config, (*fsm)(s), logStore, stableStore, snapshots, transport)
	if err != nil {
		return err
	}
	s.raft = r

	if enableSingle {
		configuration := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      config.LocalID,
					Address: transport.LocalAddr(),
				},
			},
		}
		r.BootstrapCluster(configuration)
	}

	return nil
}

func (s *Store) Get(key string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.storage[key]
}

func (s *Store) Set(key, value string) error {
	if s.raft.State() != raft.Leader {
		return errors.ErrUnsupported
	}

	c := &command{Op: "set", Key: key, Value: value}
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	return s.raft.Apply(b, raftTimeout).Error()
}

func (s *Store) Delete(key string) error {
	if s.raft.State() != raft.Leader {
		return errors.ErrUnsupported
	}

	c := &command{Op: "delete", Key: key}
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	return s.raft.Apply(b, raftTimeout).Error()
}

func (s *Store) Join(nodeID, addr string) error {
	s.logger.Printf("received join request for remote node %s at %s", nodeID, addr)

	configFuture := s.raft.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		s.logger.Printf("failed to get raft configuration: %v", err)
		return err
	}

	for _, srv := range configFuture.Configuration().Servers {
		if srv.ID == raft.ServerID(nodeID) || srv.Address == raft.ServerAddress(addr) {
			if srv.ID == raft.ServerID(nodeID) && srv.Address == raft.ServerAddress(addr) {
				s.logger.Printf("node %s at %s already member of cluster", nodeID, addr)
				return nil
			}

			future := s.raft.RemoveServer(srv.ID, 0, 0)
			if err := future.Error(); err != nil {
				return fmt.Errorf("error removing existing node %s at %s: %s", nodeID, addr, err)
			}
		}
	}

	f := s.raft.AddVoter(raft.ServerID(nodeID), raft.ServerAddress(addr), 0, 0)
	if f.Error() != nil {
		return f.Error()
	}
	s.logger.Printf("node %s at %s joined successfully", nodeID, addr)
	return nil
}

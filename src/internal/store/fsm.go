package store

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/hashicorp/raft"
)

type fsm Store

// Apply implements raft.FSM.
func (f *fsm) Apply(l *raft.Log) interface{} {
	var c command
	if err := json.Unmarshal(l.Data, &c); err != nil {
		panic("failed to unmarshal command: " + err.Error())
	}

	switch c.Op {
	case "set":
		return f.applySet(c.Key, c.Value)
	case "delete":
		return f.applyDelete(c.Key)
	default:
		return errors.ErrUnsupported
	}
}

// Restore implements raft.FSM.
func (f *fsm) Restore(snapshot io.ReadCloser) error {
	o := make(map[string]string)
	if err := json.NewDecoder(snapshot).Decode(&o); err != nil {
		return err
	}

	f.storage = o
	return nil
}

// Snapshot implements raft.FSM.
func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	o := make(map[string]string)
	for key, value := range f.storage {
		o[key] = value
	}
	return &fsmSnapshot{store: o}, nil
}

func (f *fsm) applySet(key, value string) interface{} {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.storage[key] = value
	return nil
}

func (f *fsm) applyDelete(key string) interface{} {
	f.mu.Lock()
	defer f.mu.Unlock()

	delete(f.storage, key)
	return nil
}

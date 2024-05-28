package warehouse

import (
	"encoding/json"

	"github.com/hashicorp/raft"
)

type fsmSnapshot struct {
	store map[string]string
}

// Persist implements raft.FSMSnapshot.
func (f *fsmSnapshot) Persist(sink raft.SnapshotSink) error {
	err := func() error {
		b, err := json.Marshal(f.store)
		if err != nil {
			return err
		}
		if _, err := sink.Write(b); err != nil {
			return err
		}
		return sink.Close()
	}()
	if err != nil {
		sink.Cancel()
	}
	return err
}

// Release implements raft.FSMSnapshot.
func (f *fsmSnapshot) Release() {}

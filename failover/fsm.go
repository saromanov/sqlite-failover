package failover

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/hashicorp/raft"
)

// FSM implementation of the Raft FSM
type FSM struct {
	mu      *sync.Mutex
	masters []string
	m       map[string]string
}

type command struct {
	Op    string `json:"op,omitempty"`
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type masterSnapshot struct {
	masters []string
	store   map[string]string
}

func (m *masterSnapshot) Persist(sink raft.SnapshotSink) error {
	err := func() error {
		b, err := json.Marshal(m.store)
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

func (m *masterSnapshot) Release() {
}

func (f *FSM) Apply(l *raft.Log) interface{} {
	var c command
	if err := json.Unmarshal(l.Data, &c); err != nil {
		panic(fmt.Sprintf("failed to unmarshal command: %s", err.Error()))
	}

	f.handleAction(&c)
	return nil
}

func (f *FSM) handleAction(c *command) {
	switch c.Op {
	case "set":
		f.handleSet(c)
	case "del":
		f.handleDelete(c)
	default:
		panic("unable to find command")
	}
}

func (f *FSM) handleSet(c *command) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.m[c.Key] = c.Value
	return nil
}

func (f *FSM) handleDelete(c *command) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	for k, v := range f.m {
		if k == c.Key {
			delete(f.m, k)
		}
	}
	return nil
}

func (f *FSM) Snapshot() (raft.FSMSnapshot, error) {
	snap := new(masterSnapshot)
	snap.masters = make([]string, 0, len(f.masters))

	f.mu.Lock()
	defer f.mu.Unlock()
	for _, master := range f.masters {
		snap.masters = append(snap.masters, master)
	}
	return snap, nil
}

func (f *FSM) Restore(io.ReadCloser) error {
	return nil
}

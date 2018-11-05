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
	*sync.Mutex
	masters []string
}

type command struct {
	Op    string `json:"op,omitempty"`
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

func (f *FSM) Apply(l *raft.Log) interface{} {
	var c command
	if err := json.Unmarshal(l.Data, &c); err != nil {
		panic(fmt.Sprintf("failed to unmarshal command: %s", err.Error()))
	}

	f.handleAction(&c)
	return nil
}

func (f *FSM) handleAction(c*Command){

}

func (f *FSM) Snapshot() (raft.FSMSnapshot, error) {
	snap := new(masterSnapshot)
	snap.masters = make([]string, 0, len(f.masters))

	f.Lock()
	defer f.Unlock()
	for master := range fsm.masters {
		snap.masters = append(snap.masters, master)
	}
	return snap, nil
}

func (f *FSM) Restore(io.ReadCloser) error {
	return nil
}

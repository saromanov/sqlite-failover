package failover

import (
	"io"

	"github.com/hashicorp/raft"
)

// FSM implementation of the Raft FSM
type FSM struct {
}

func (f *FSM) Apply(*raft.Log) interface{} {
	return nil
}

func (f *FSM) Snapshot() (raft.FSMSnapshot, error) {
	return nil, nil
}

func (f *FSM) Restore(io.ReadCloser) error {
	return nil
}

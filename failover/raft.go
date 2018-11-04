package failover

import (
	"os"
	"time"

	"github.com/hashicorp/raft"
	"github.com/hashicorp/raft-boltdb"
)

// Failover defines main structure
type Failover struct {
	raft      *raft.Raft
	dbStore   *raftboltdb.BoltStore
	raftAddr  string
	transport *raft.NetworkTransport
}

// New creates a new failover
func New(c *Config) (*Failover, error) {
	conf := raft.DefaultConfig()
	fileStore, err := raft.NewFileSnapshotStore(c.RaftDir, 1, os.Stderr)
	if err != nil {
		return nil, err
	}
	trans, err := raft.NewTCPTransport(c.RaftAddr, nil, 3, 3*time.Second, os.Stderr)
	if err != nil {
		return nil, err
	}
	dbStore, err := raftboltdb.NewBoltStore(c.RaftDBPath)
	if err != nil {
		return nil, err
	}
	logStore := dbStore
	fsm := &FSM{}
	r, err := raft.NewRaft(conf, fsm, logStore, dbStore, fileStore, trans)
	if err != nil {
		return nil, err
	}
	return &Failover{
		raft:      r,
		transport: trans,
	}, nil
}

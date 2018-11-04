package failover

import (
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
	fileStore, err := raft.NewFileSnapshotStore(c.RaftDir, 1, log)
	if err != nil {
		return nil, err
	}
	trans, err := raft.NewTCPTransport(c.RaftAddr, nil, 3, 3*time.Second, log)
	if err != nil {
		return nil, err
	}
	dbStore, err := raftboltdb.NewBoltStore(c.RaftDBPath)
	if err != nil {
		return nil, err
	}
	peers := []string{}
	r, err := raft.NewRaft(c, fsm, dbStore, dbStore, fileStore, peers, trans)
	return &Failover{
		raft:      r,
		transport: trans,
	}, nil
}

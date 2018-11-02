package failover

import (
	"github.com/hashicorp/raft"
	"github.com/hashicorp/raft-boltdb"
)

// Failover defines main structure
type Failover struct {
	raft     *raft.Raft
	dbStore  *raftboltdb.BoltStore
	raftAddr string
}

// New creates a new failover
func New(c *Config) *Failover {
	conf := raft.DefaultConfig()
	r, err := raft.NewRaft(config, fsm, dbStore, dbStore, fileStore, r.peerStore, nil)
	return &Failover{
		raft:r,
	}
}

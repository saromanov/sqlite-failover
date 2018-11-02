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
func New() *Failover {
	return &Failover{}
}

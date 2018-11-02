package failover

import (
	"github.com/hashicorp/raft"
)

// Failover defines main structure
type Failover struct {
	raft *raft.Raft
}

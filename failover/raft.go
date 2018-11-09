package failover

import (
	"errors"
	"os"
	"time"

	"github.com/hashicorp/raft"
	"github.com/hashicorp/raft-boltdb"
)

var errNotLeader = errors.New("node not a leader")

// Failover defines main structure
type Failover struct {
	raft       *raft.Raft
	dbStore    *raftboltdb.BoltStore
	raftAddr   string
	transport  *raft.NetworkTransport
	raftConfig *raft.Config
}

// NewRaft creates a new raft failover
func NewRaft(c *Config) (*Failover, error) {
	conf := raft.DefaultConfig()
	conf.LocalID = raft.ServerID(c.LocalID)
	fileStore, err := raft.NewFileSnapshotStore(c.RaftDir, 1, os.Stderr)
	if err != nil {
		return nil, err
	}
	trans, err := raft.NewTCPTransport(c.RaftAddr, nil, 3, 3*time.Second, os.Stderr)
	if err != nil {
		return nil, err
	}
	logStore, dbStore, err := initDBStore(c)
	if err != nil {
		return nil, err
	}
	fsm := &FSM{}
	r, err := raft.NewRaft(conf, fsm, logStore, dbStore, fileStore, trans)
	if err != nil {
		return nil, err
	}
	return &Failover{
		raft:       r,
		transport:  trans,
		raftConfig: conf,
	}, nil
}

// initDBStore provides init of the store
func initDBStore(c *Config) (raft.LogStore, raft.StableStore, error) {
	if c.InMemoryStore {
		return raft.NewInmemStore(), raft.NewInmemStore(), nil
	}
	dbStore, err := raftboltdb.NewBoltStore(c.RaftDBPath)
	if err != nil {
		return nil, nil, err
	}
	return dbStore, dbStore, nil
}

// ISLeader retruns true if node is leader
func (f *Failover) IsLeader() bool {
	addr := f.raft.Leader()
	if addr == "" {
		return false
	}
	return string(addr) == f.raftAddr
}

// Set provides setting of the key value to raft store
func (f *Failover) Set(key, value string) error {
	if f.raft.State() != raft.Leader {
		return errNotLeader
	}

	return nil
}

// Get provides getting of the key
func (f *Failover) Get(key string) string {
	return ""
}

// Run provides starting of the application
func (f *Failover) Run() error {
	configuration := raft.Configuration{
		Servers: []raft.Server{
			{
				ID:      f.raftConfig.LocalID,
				Address: f.transport.LocalAddr(),
			},
		},
	}

	f.raft.BootstrapCluster(configuration)
	return nil
}

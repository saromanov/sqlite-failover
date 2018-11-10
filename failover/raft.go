package failover

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/raft"
	"github.com/hashicorp/raft-boltdb"
)

var errNotLeader = errors.New("node not a leader")

// Failover defines main structure
type Failover struct {
	raft        *raft.Raft
	dbStore     *raftboltdb.BoltStore
	raftAddr    string
	transport   *raft.NetworkTransport
	raftConfig  *raft.Config
	masters     FSM
	mastersAddr []string
}

// NewRaft creates a new raft failover
func NewRaft(c *Config) (*Failover, error) {
	conf := raft.DefaultConfig()
	c.LocalID = "machine-2"
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
		raft:        r,
		transport:   trans,
		raftConfig:  conf,
		masters:     FSM{},
		mastersAddr: []string{},
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

// IsLeader retruns true if node is leader
func (f *Failover) IsLeader() bool {
	addr := f.raft.Leader()
	if addr == "" {
		return false
	}
	return string(addr) == f.raftAddr
}

// AddVoter provides joining node to the cluster
func (f *Failover) AddVoter(addr []string) error {
	for _, a := range addr {
		fmt.Println(a)
		resp := f.raft.AddVoter(raft.ServerID(a), raft.ServerAddress(a), 0, 0)
		if resp.Error() != nil {
			fmt.Println("Error:", resp.Error())
			continue
		}
	}
	return nil
}

// Set provides setting of the key value to raft store
func (f *Failover) Set(key, value string) error {
	if f.raft.State() != raft.Leader {
		return errNotLeader
	}

	c := command{
		Op:    "set",
		Key:   key,
		Value: value,
	}
	d, err := json.Marshal(c)
	if err != nil {
		return err
	}

	return f.raft.Apply(d, 5*time.Second).Error()
}

// AddMasters provides append of new masters to the cluster
func (f *Failover) AddMasters(addr []string) error {
	if f.raft.State() != raft.Leader {
		return errNotLeader
	}

	c := command{
		Op:      "add_masters",
		Masters: addr,
	}
	d, err := json.Marshal(c)
	if err != nil {
		return err
	}

	return f.raft.Apply(d, 5*time.Second).Error()
}

// Get provides getting of the key
func (f *Failover) Get(key string) string {
	return ""
}

// GetMasters provides getting of masters
func (f *Failover) GetMasters() []string {
	return f.mastersAddr
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

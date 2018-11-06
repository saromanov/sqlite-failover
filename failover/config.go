package failover

// Config defines configuration for sqlite-failover
type Config struct {
	RaftAddr      string `yaml:"raft_addr"`
	RaftDir       string `yaml:"raft_dir"`
	RaftDBPath    string `yaml:"raft_db_path"`
	LocalID       string `yaml:"local_id"`
	InMemoryStore bool   `yaml:"in_memory_store"`
	Interval      int64  `yaml:"interval"`
	Addr          string `yaml:"addr"`
}

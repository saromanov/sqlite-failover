package failover

// Config defines configuration for sqlite-failover
type Config struct {
	RaftAddr string `yaml:"raft_addr"`
}

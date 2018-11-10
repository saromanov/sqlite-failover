package main

import (
	"github.com/saromanov/sqlite-failover/failover"
)

func main() {
	cfg := &failover.Config{
		RaftDir:    "$HOME/raftdir",
		RaftDBPath: "$HOME/raftdbdir",
		RaftAddr:   "127.0.0.1:9000",
		LocalID:    "machine-1",
	}
	f := failover.New(cfg)
	f.Start()
}

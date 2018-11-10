package main

import (
	"flag"
	"fmt"

	"github.com/saromanov/sqlite-failover/failover"
)

var (
	addr     *string
	raftAddr *string
)

func parse() {
	addr = flag.String("addr", "Address", "a string")
	raftAddr = flag.String("raft-addr", "Raft Address", "a string")
	flag.Parse()
}
func main() {
	parse()
	fmt.Println(*addr)
	cfg := &failover.Config{
		RaftDir:    "$HOME/raftdir",
		RaftDBPath: "$HOME/raftdbdir",
		RaftAddr:   *raftAddr,
		LocalID:    "machine-1",
		Addr:       *addr,
	}
	f := failover.New(cfg)
	f.Start()
}

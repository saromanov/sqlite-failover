package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/saromanov/sqlite-failover/failover"
)

var (
	addr     *string
	raftAddr *string
	// join defines list of the addresses
	// for joining to the cluster
	join *string
)

func parse() {
	addr = flag.String("addr", "", "a string")
	raftAddr = flag.String("raft-addr", "", "a string")
	join = flag.String("join", "List of the addrsses for join to the cluster", "addr")
	flag.Parse()
}
func main() {
	parse()
	fmt.Println(*addr)
	if *addr == "" {
		panic("address is not defined")
	}
	if *raftAddr == "" {
		panic("raft address is not defined")
	}
	cfg := &failover.Config{
		RaftDir:       "$HOME/raftdir",
		RaftDBPath:    "$HOME/raftdbdir",
		RaftAddr:      *raftAddr,
		LocalID:       *addr,
		Addr:          *addr,
		InMemoryStore: true,
	}
	f := failover.New(cfg)
	time.Sleep(5 * time.Second)
	if *join != "" {
		splitter := strings.Split(*join, ",")
		f.Join(splitter)
	}
	f.Start()
}

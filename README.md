# sqlite-failover

Sqlite with Failover based on Raft

Starting of the App
`go run main.go -addr=127.0.0.1:8080 -raft-addr=127.0.0.1:9001`

Joinign node to the cluster
`go run main.go -addr=127.0.0.1:8081 -raft-addr=127.0.0.1:9000 -join=127.0.0.1:9001`
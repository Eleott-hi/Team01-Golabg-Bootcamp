package main

import "flag"

const (
	DefaultHTTPAddr = "localhost:11000"
	DefaultRaftAddr = "localhost:12000"
)

var (
	httpAddr string
	raftAddr string
	joinAddr string
	nodeID   string
)

func init() {
	flag.StringVar(&httpAddr, "h", DefaultHTTPAddr, "HTTP bind address")
	flag.StringVar(&raftAddr, "r", DefaultRaftAddr, "Raft bind address")
	flag.StringVar(&joinAddr, "j", "", "Join address (if any)")
	flag.StringVar(&nodeID, "n", "", "node id. if not set, same as raft bind address")
}

func main() {
	flag.Parse()

	if nodeID == "" {
		nodeID = raftAddr
	}

}

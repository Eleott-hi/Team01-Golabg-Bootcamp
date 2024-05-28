package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"team01/internal/store"
	"team01/internal/warehouse"
)

const (
	DefaultHTTPAddr = "localhost:11000"
	DefaultRaftAddr = "localhost:12000"
)

var (
	httpAddr string
	raftAddr string
	raftDir  string
	joinAddr string
	nodeID   string
)

func init() {
	flag.StringVar(&httpAddr, "h", DefaultHTTPAddr, "HTTP bind address")
	flag.StringVar(&raftAddr, "r", DefaultRaftAddr, "Raft bind address")
	flag.StringVar(&raftDir, "o", "", "Raft directory")
	flag.StringVar(&joinAddr, "j", "", "Join address (if any)")
	flag.StringVar(&nodeID, "n", "", "node id. if not set, same as raft bind address")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s [options] <raft-data-path>\n", os.Args[0])
	}
}

func main() {
	flag.Parse()

	if raftDir == "" {
		flag.PrintDefaults()
		log.Fatal("no raft storage dir specified")
	}

	if nodeID == "" {
		nodeID = raftAddr
	}

	if err := os.MkdirAll(raftDir, 0700); err != nil {
		log.Fatalf("failed to create path for Raft storage: %v", err)
	}

	s := store.New(raftDir, raftAddr)
	if err := s.Open(joinAddr == "", nodeID); err != nil {
		log.Fatalf("failed to open store: %v", err)
	}

	wh := warehouse.New(httpAddr, s)
	if err := wh.Start(); err != nil {
		log.Fatalf("failed to start warehouse service: %v", err)
	}

	if joinAddr != "" {
		if err := join(joinAddr, raftDir, nodeID); err != nil {
			log.Fatalf("failed to join node at %s: %v", joinAddr, err)
		}
	}

	log.Printf("warehouse listening on %s", httpAddr)

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt)
	<-terminate
}

func join(joinAddr, raftAddr, nodeID string) error {
	b, err := json.Marshal(map[string]string{"addr": raftAddr, "id": nodeID})
	if err != nil {
		return err
	}
	res, err := http.Post(fmt.Sprintf("http://%s/join", joinAddr), "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	res.Body.Close()
	return nil
}

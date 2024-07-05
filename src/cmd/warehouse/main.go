package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"team01/internal/store"
	"team01/internal/warehouse"
)

const (
	DefaultHTTPAddr = "localhost:11000"
	DefaultRaftAddr = "localhost:12000"
)

type statemachine struct {
	db     *sync.Map
	server int
}

type commandKind uint8

const (
	setCommand commandKind = iota
	getCommand
)

type command struct {
	kind  commandKind
	key   string
	value string
}

func (s *statemachine) Apply(cmd []byte) ([]byte, error) {
	c := decodeCommand(cmd)

	switch c.kind {
	case setCommand:
		s.db.Store(c.key, c.value)
	case getCommand:
		value, ok := s.db.Load(c.key)
		if !ok {
			return nil, nil
		}
		return []byte(value.(string)), nil
	default:
		return nil, fmt.Errorf("unknown command: %x", cmd)
	}

	return nil, nil
}

func encodeCommand(c command) []byte {
	msg := bytes.NewBuffer(nil)
	err := msg.WriteByte(uint8(c.kind))
	if err != nil {
		panic(err)
	}

	err = binary.Write(msg, binary.LittleEndian, uint64(len(c.key)))
	if err != nil {
		panic(err)
	}

	msg.WriteString(c.key)

	err = binary.Write(msg, binary.LittleEndian, uint64(len(c.value)))
	if err != nil {
		panic(err)
	}

	msg.WriteString(c.value)

	return msg.Bytes()
}

func decodeCommand(msg []byte) command {
	var c command
	c.kind = commandKind(msg[0])

	keyLen := binary.LittleEndian.Uint64(msg[1:9])
	c.key = string(msg[9 : 9+keyLen])

	if c.kind == setCommand {
		valLen := binary.LittleEndian.Uint64(msg[9+keyLen : 9+keyLen+8])
		c.value = string(msg[9+keyLen+8 : 9+keyLen+valLen])
	}

	return c
}

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

	if nodeID == "" {
		nodeID = raftAddr
	}
	if raftDir == "" {
		raftDir = nodeID
	}

	if err := os.MkdirAll(raftDir, 0o700); err != nil {
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
		if err := join(joinAddr, raftAddr, nodeID); err != nil {
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

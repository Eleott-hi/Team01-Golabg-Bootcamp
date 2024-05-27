package heartbeat

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

type Node struct {
	Address       string    `json:"address"`
	IsLeader      bool      `json:"is_leader"`
	LastHeartbeat time.Time `json:"-"`
}

type Cluster struct {
	Nodes  map[string]*Node `json:"nodes"`
	Leader string           `json:"leader"`
	mu     sync.Mutex
}

var cluster Cluster

func heartbeatHandler(w http.ResponseWriter, r *http.Request) {
	var node Node
	if err := json.NewDecoder(r.Body).Decode(&node); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	cluster.mu.Lock()
	defer cluster.mu.Unlock()

	if existingNode, ok := cluster.Nodes[node.Address]; ok {
		existingNode.LastHeartbeat = time.Now()
		existingNode.IsLeader = node.IsLeader
	} else {
		node.LastHeartbeat = time.Now()
		cluster.Nodes[node.Address] = &node
	}

	if node.IsLeader {
		cluster.Leader = node.Address
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cluster)
}

func sendHeartbeat(address, myAddress string, isLeader bool) {
	node := Node{
		Address:  myAddress,
		IsLeader: isLeader,
	}
	data, _ := json.Marshal(node)
	if _, err := http.Post("http://"+address+"/heartbeat", "application/json", bytes.NewBuffer(data)); err != nil {
		log.Printf("falied to send heartbeat to %s: %v", address, err)
	}

}

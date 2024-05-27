package warehouse

import (
	context "context"
	"encoding/json"
	"errors"
	"net/http"
	sync "sync"

	"google.golang.org/grpc"
)

func StartServer(grpcServer *grpc.Server) {
	RegisterWareHouseServer(grpcServer, &wareHouseServer{
		storage: make(map[string]string),
	})
}

type wareHouseServer struct {
	UnimplementedWareHouseServer
	storage map[string]string
	mu      sync.RWMutex
}

func (s *wareHouseServer) SetValue(ctx context.Context, pair *Pair) (*Empty, error) {
	tx := map[string]string{
		"key":   pair.GetKey(),
		"value": pair.GetValue(),
	}
	txBytes, err := json.Marshal(tx)
	if err != nil {
		return nil, err
	}

	res, err := http.Post("http://127.0.0.1:26657/broadcast_tx_commit?tx=\""+string(txBytes)+"\"", "application/json", nil)
	if err != nil || res.StatusCode != http.StatusOK {
		return nil, err
	}

	return &Empty{}, nil
}

func (s *wareHouseServer) GetValue(ctx context.Context, key *Key) (*Result, error) {
	res, err := http.Get("http://127.0.0.1:26657/abci_query?data=\"" + key.GetKey() + "\"")
	if err != nil || res.StatusCode != http.StatusOK {
		return nil, err
	}

	var queryRes map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&queryRes); err != nil {
		return nil, err
	}

	value := queryRes["result"].(map[string]interface{})["response"].(map[string]interface{})["value"]
	if value == nil {
		return nil, errors.New("not found")
	}

	return &Result{Message: value.(string)}, nil
}

func (s *wareHouseServer) DeleteValue(ctx context.Context, key *Key) (*Empty, error) {
	tx := map[string]string{"key": key.GetKey(), "delete": "true"}
	txBytes, err := json.Marshal(tx)
	if err != nil {
		return nil, err
	}

	res, err := http.Post("http://127.0.0.1:26657/broadcast_tx_commit?tx=\""+string(txBytes)+"\"", "application/json", nil)
	if err != nil || res.StatusCode != http.StatusOK {
		return nil, err
	}

	return &Empty{}, nil
}

// func (s *wareHouseServer) SetValue(ctx context.Context, pair *Pair) (*Empty, error) {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	s.storage[pair.GetKey()] = pair.GetValue()

// 	return &Empty{}, nil
// }

// func (s *wareHouseServer) GetValue(ctx context.Context, key *Key) (*Result, error) {
// 	var value string
// 	s.mu.RLock()
// 	defer s.mu.RUnlock()

// 	value = s.storage[key.GetKey()]

// 	return &Result{Message: value}, nil
// }

// func (s *wareHouseServer) DeleteValue(ctx context.Context, key *Key) (*Empty, error) {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	delete(s.storage, key.GetKey())

// 	return &Empty{}, nil
// }

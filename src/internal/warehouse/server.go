package warehouse

import (
	context "context"
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
	s.mu.Lock()
	defer s.mu.Unlock()

	s.storage[pair.GetKey()] = pair.GetValue()

	return &Empty{}, nil
}

func (s *wareHouseServer) GetValue(ctx context.Context, key *Key) (*Result, error) {
	var value string
	s.mu.RLock()
	defer s.mu.RUnlock()

	value = s.storage[key.GetKey()]

	return &Result{Message: value}, nil
}

func (s *wareHouseServer) DeleteValue(ctx context.Context, key *Key) (*Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.storage, key.GetKey())

	return &Empty{}, nil
}

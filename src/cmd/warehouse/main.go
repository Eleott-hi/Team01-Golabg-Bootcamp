package main

import (
	"context"
	"log"
	"net"
	"sync"
	pb "team01/internal/warehouse"

	"google.golang.org/grpc"
)

type wareHouseServer struct {
	pb.UnimplementedWareHouseServer
	storage map[string]string
	mu      sync.RWMutex
}

func newServer() *wareHouseServer {
	return &wareHouseServer{
		storage: make(map[string]string),
	}
}

func (s *wareHouseServer) SetValue(ctx context.Context, pair *pb.Pair) (*pb.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.storage[pair.GetKey()] = pair.GetValue()

	return &pb.Empty{}, nil
}

func (s *wareHouseServer) GetValue(ctx context.Context, key *pb.Key) (*pb.Result, error) {
	var value string
	s.mu.RLock()
	defer s.mu.RUnlock()

	value = s.storage[key.GetKey()]

	return &pb.Result{Message: value}, nil
}

func (s *wareHouseServer) DeleteValue(ctx context.Context, key *pb.Key) (*pb.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.storage, key.GetKey())

	return &pb.Empty{}, nil
}

func main() {
	// replicationFactor := flag.Int("r", 2, "replication factor")
	// flag.Parse()

	lis, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterWareHouseServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}

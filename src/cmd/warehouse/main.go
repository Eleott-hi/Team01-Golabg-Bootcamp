package main

import (
	"log"
	"net"
	"team01/internal/tendermmint"
	"team01/internal/warehouse"

	"google.golang.org/grpc"
)

func main() {
	// replicationFactor := flag.Int("r", 2, "replication factor")
	// flag.Parse()

	app := &tendermmint.Application{}
	go app.StartABCI("127.0.0.1:26658")

	lis, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	warehouse.StartServer(grpcServer)
	grpcServer.Serve(lis)
}

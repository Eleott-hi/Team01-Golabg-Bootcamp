run-cli:
	go run cmd/warehouse-cli/main.go

run-warehouse:
	go run cmd/warehouse/main.go

build:
	go build -o warehouse cmd/warehouse/main.go

start-leader:
	./warehouse -h :11000 -r :12000 -n node0 node0

start-follower1:
	./warehouse -h :11001 -r :12001 -j :11000 -n node1 node1

start-follower2:
	./warehouse -h :11002 -r :12002 -j :11000 -n node2 node2

proto:
	protoc \
		--proto_path=api \
		--go_out . \
		--go-grpc_out . \
		api/warehouse.proto

tidy:
	go mod tidy

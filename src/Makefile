run-cli:
	go run cmd/warehouse-cli/main.go

run-warehouse:
	go run cmd/warehouse/main.go

build:
	go build -o warehouse cmd/warehouse/main.go

start-leader:
	./warehouse \
		-h localhost:11000 \
		-r localhost:12000 \
		-n node0

start-follower1:
	./warehouse \
		-h localhost:11001 \
		-r localhost:12001 \
		-j localhost:11000 \
		-n node1

start-follower2:
	./warehouse \
		-h localhost:11002 \
		-r localhost:12002 \
		-j localhost:11000 \
		-n node2

proto:
	protoc \
		--proto_path=api \
		--go_out . \
		--go-grpc_out . \
		api/warehouse.proto

tidy:
	go mod tidy

run-cli:
	go run cmd/warehouse-cli/main.go

run-warehouse:
	go run cmd/warehouse/main.go

proto:
	protoc \
		--proto_path=api \
		--go_out . \
		--go-grpc_out . \
		api/warehouse.proto

tidy:
	go mod tidy

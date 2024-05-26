all:
	@echo "No target specified for \"all\""

run-server: go.mod tidy
	go run ./cmd/database

run-client: go.mod tidy
	go run ./cmd/client

go.mod:
	go mod init team01

tidy:
	go mod tidy

build: go.mod tidy
	go build -o build/database ./cmd/database
	go build -o build/client ./cmd/client

start-server: build
	./build/database

start-client: build
	./build/client

.PHONY: test
test:
	go test -v -coverpkg=./... ./test/...
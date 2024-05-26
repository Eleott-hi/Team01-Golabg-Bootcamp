

all: run

run: go.mod tidy
	go run ./cmd/database

go.mod:
	go mod init team01

tidy:
	go mod tidy

build: go.mod tidy
	go build -o build/database ./cmd/database

start: build
	./build/database

.PHONY: test
test:
	go test -v -coverpkg=./... ./test
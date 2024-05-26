

all: run

run: go.mod tidy
	go run .

go.mod:
	go mod init team01

tidy:
	go mod tidy

build: go.mod tidy
	go build -o build/main .

start: build
	./build/main

.PHONY: test
test:
	go test -v ./test
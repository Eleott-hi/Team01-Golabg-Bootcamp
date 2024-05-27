package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

var storage = make(map[string]string)

func main() {
	// replicationFactor := flag.Int("r", 2, "replication factor")
	// flag.Parse()

	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatalf("Error listening: %v", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting request: ", err)
			return
		}
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	request, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading request: ", err)
		return
	}
	command := strings.TrimSpace(request)
	response := executeCommand(command)
	conn.Write([]byte(response + "\n"))
}

func executeCommand(command string) string {
	fmt.Println("Command: ", command)

	parts := strings.Fields(command)
	if len(parts) < 2 {
		return "Error: imvalid command"
	}

	action := parts[0]
	key := parts[1]

	switch action {
	case "SET":
		if len(parts) < 3 {
			return "Error: invalid command"
		}
		value := parts[2]
		storage[key] = value
		return fmt.Sprintf("Created: [%s]=%s", key, value)
	case "GET":
		if value, exists := storage[key]; exists {
			return value
		} else {
			return "Not found"
		}
	case "DELETE":
		delete(storage, key)
		return "Deleted"
	default:
		return "Error: unknown command"
	}
}

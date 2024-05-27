package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

var (
	host string
	port string
)

func main() {
	flag.StringVar(&host, "H", "127.0.0.1", "database host")
	flag.StringVar(&port, "P", "8080", "database port")
	flag.Parse()

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		log.Fatalf("Cannot dial: %v", err)
	}
	defer conn.Close()
	fmt.Printf("Connected to a database of Warehouse {NUM} at %s:%s\n", host, port)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err == io.EOF {
			continue
		}
		text = strings.TrimSpace(text)

		if len(text) > 0 {
			handleCommand(conn, text)
		}
	}
}

func handleCommand(conn net.Conn, command string) {
	fmt.Println("Command: ", command)
	fmt.Fprintf(conn, command+"\n")
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err == io.EOF {
		fmt.Println("Empty response")
		return
	}
	if err != nil {
		fmt.Println("Error reading response: ", err)
		return
	}
	fmt.Println("Response: ", response)
}

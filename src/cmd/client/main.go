package main

import (
	"log"
	"os"
	
	cli "team01/pkg/client-cli"
)

func main() {
	client := cli.New()
	if err := client.Run(); err != nil {
		log.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

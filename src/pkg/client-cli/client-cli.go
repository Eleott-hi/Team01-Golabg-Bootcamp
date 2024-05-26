package clientcli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type ClientCli struct{}

func New() *ClientCli {
	return &ClientCli{}
}

func (c *ClientCli) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nReceived interrupt signal. Shutting down gracefully...")
		cancel()
	}()

	inputChan := make(chan string)
	go c.readInput(ctx, inputChan)

	for {
		fmt.Printf("Name, Age: ")
		
		select {
		case <-ctx.Done():
			fmt.Println("Context cancelled, exiting loop.")
			return nil
		case input := <-inputChan:

			if input == "" {
				continue
			}
			var name string
			var age int
			if _, err := fmt.Sscanf(input, "%s %d", &name, &age); err != nil {
				fmt.Println("Invalid input. Please enter your name and age.")
			} else {
				fmt.Printf("Name: %s, Age: %d\n", name, age)
			}
		}
	}
}

func (c *ClientCli) readInput(ctx context.Context, inputChan chan<- string) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		select {
		case <-ctx.Done():
			close(inputChan)
			return
		default:
			if scanner.Scan() {
				inputChan <- scanner.Text()
			} else {
				close(inputChan)
				return
			}
		}
	}
}

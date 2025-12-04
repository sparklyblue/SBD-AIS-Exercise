package main

import (
	"exc8/client"
	"exc8/server"
	"log"
	"time"
)

func main() {
	// Start server in a goroutine
	go func() {
		if err := server.StartGrpcServer(); err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	// Wait for server to start
	time.Sleep(500 * time.Millisecond)

	// Run client
	c, err := client.NewGrpcClient()
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	if err := c.Run(); err != nil {
		log.Fatalf("client error: %v", err)
	}
}

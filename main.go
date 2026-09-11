package main

import (
	"log"

	"github.com/Phran6ix/redis-kv/cmd/server"
)

func main() {
	listener, err := server.Listener(":6979")
	if err != nil {
		log.Fatalf("Failed to start the listener, %v", err)
	}

	log.Printf("now receiveing connection on port 6979")

	err = server.Start(listener)
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

package main

import (
	"fmt"
	"log"
)

const (
	AppName    = "Stumpf.Works NAS"
	AppVersion = "0.1.0-alpha"
)

func main() {
	fmt.Printf("%s v%s\n", AppName, AppVersion)
	fmt.Println("Starting server...")

	// TODO: Initialize configuration
	// TODO: Initialize database
	// TODO: Initialize logger
	// TODO: Initialize API server
	// TODO: Initialize WebSocket server
	// TODO: Load plugins
	// TODO: Start HTTP server

	log.Println("Server will start here (Phase 2)")
}

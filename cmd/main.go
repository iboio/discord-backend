package main

import (
	"discord/internal/handler"
)

func main() {
	// Start all services (gRPC server, REST API, etc.)
	handler.Handler()
}

package main

import (
	"log"
	"pongServer/pkg/router"
)

func main() {
	//initialize router
	r := router.SetupRouter()
	// start the server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server Failed to start %v", err)
	}
}

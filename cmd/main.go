package main

import (
	"github.com/radean0909/guild-chat/api"
)

func main() {
	// Create a new API service
	svc := api.New()

	svc.Start(":8000")
}

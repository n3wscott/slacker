package main

import (
	"log"

	"github.com/n3wscott/slacker/pkg/commands"
)

func main() {
	if err := commands.New().Execute(); err != nil {
		log.Fatalf("error during command execution: %v", err)
	}
}

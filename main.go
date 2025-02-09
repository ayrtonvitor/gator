package main

import (
	"log"
	"os"

	"github.com/ayrtonvitor/gator/internal/command"
	_ "github.com/lib/pq"
)

func main() {
	state, err := setup()
	if err != nil {
		log.Fatal("Could not initialize app: %w", err)
	}

	commands := command.GetCommandList()
	err = commands.TryRunInputCommand(os.Args, state)
	if err != nil {
		log.Fatalf("Could not execute command: %v", err)
	}
}

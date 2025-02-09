package main

import (
	"log"
	"os"

	"github.com/ayrtonvitor/gator/internal/command"
	"github.com/ayrtonvitor/gator/internal/config"
	"github.com/ayrtonvitor/gator/internal/state"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		log.Fatalf("Could not read %v", err)
	}
	state := &state.State{
		Config: &conf,
	}
	commands := command.GetCommandList()
	err = commands.TryRunInputCommand(os.Args, state)
	if err != nil {
		log.Fatalf("Could not execute command: %v", err)
	}
}

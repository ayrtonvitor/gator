package command

import (
	"github.com/ayrtonvitor/gator/internal/state"
)

type command struct {
	Name string
	Args []string
}

type commands struct {
	Handlers map[string]func(*state.State, command) error
}

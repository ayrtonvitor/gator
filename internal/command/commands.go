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

func (c *commands) register(name string, f func(*state.State, command) error) {
	c.Handlers[name] = f
}

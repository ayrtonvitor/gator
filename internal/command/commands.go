package command

import (
	"fmt"

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

func (c *commands) run(s *state.State, cmd command) error {
	handler, ok := c.Handlers[cmd.Name]
	if !ok {
		return fmt.Errorf("%s is not a valid command.", cmd.Name)
	}
	return handler(s, cmd)
}

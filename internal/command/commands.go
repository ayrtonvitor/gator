package command

import (
	"errors"
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

func GetCommandList() commands {
	cmds := commands{
		Handlers: make(map[string]func(*state.State, command) error),
	}
	cmds.register("login", login)
	cmds.register("register", register)
	cmds.register("reset", reset)
	cmds.register("users", listUsers)
	cmds.register("agg", aggregate)
	cmds.register("addfeed", middlewareLoggedIn(addFeed))
	cmds.register("feeds", listFeeds)
	cmds.register("follow", middlewareLoggedIn(follow))
	cmds.register("following", following)
	cmds.register("unfollow", middlewareLoggedIn(unfollow))
	return cmds
}

func (c *commands) TryRunInputCommand(cmdLineargs []string, state *state.State) error {
	if len(cmdLineargs) < 2 {
		return errors.New("Did not pass a command")
	}
	cmd := command{
		Name: cmdLineargs[1],
		Args: cmdLineargs[2:],
	}
	return c.run(state, cmd)
}

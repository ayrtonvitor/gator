package command

import (
	"errors"
	"fmt"

	"github.com/ayrtonvitor/gator/internal/state"
)

func login(s *state.State, c command) error {
	if len(c.Args) != 1 {
		return errors.New("Login expects a single argument (`user name`)")
	}
	usr := c.Args[0]
	err := s.Config.SetUser(usr)
	if err != nil {
		return fmt.Errorf("Could not login: %w", err)
	}
	fmt.Printf("Logged in as %s\n", usr)
	return nil
}

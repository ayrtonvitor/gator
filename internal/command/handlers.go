package command

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ayrtonvitor/gator/internal/database"
	"github.com/ayrtonvitor/gator/internal/state"
	"github.com/google/uuid"
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

func register(s *state.State, c command) error {
	if len(c.Args) != 1 {
		return errors.New("Register expects a single argument (`user name`)")
	}

	_, err := s.Db.GetUser(context.Background(), c.Args[0])
	if err == nil {
		return fmt.Errorf("User %s already exists", c.Args[0])
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("Could not register %s: %w", c.Args[0], err)
	}

	usr := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      c.Args[0],
	}
	userModel, err := s.Db.CreateUser(context.Background(), usr)
	if err != nil {
		return fmt.Errorf("Could not register %s: %w", c.Args[0], err)
	}
	return s.Config.SetUser(userModel.Name)

}

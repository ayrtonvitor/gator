package command

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ayrtonvitor/gator/internal/database"
	"github.com/ayrtonvitor/gator/internal/rss"
	"github.com/ayrtonvitor/gator/internal/state"
	"github.com/google/uuid"
)

func login(s *state.State, c command) error {
	if len(c.Args) != 1 {
		return errors.New("Login expects a single argument (`user name`)")
	}
	usrName := c.Args[0]
	if usrName == s.Config.CurrentUserName {
		return errors.New("Can not change login to the same user")
	}

	usr, err := s.Db.GetUser(context.Background(), usrName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("User %s does not exist. Register a new user with `register`", usrName)
		}
		return fmt.Errorf("Could not login with user %s", usrName)
	}

	err = s.Config.SetUser(usr.Name)
	if err != nil {
		return fmt.Errorf("Could not login: %w", err)
	}
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

func reset(s *state.State, _ command) error {
	err := s.Db.Reset(context.Background())
	if err != nil {
		return fmt.Errorf("Unsuccessful reset: %w", err)
	}
	fmt.Println("Successful reset")
	return nil
}

func listUsers(s *state.State, _ command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Could not get users from db: %w", err)
	}
	for _, user := range users {
		if user.Name == s.Config.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
			continue
		}

		fmt.Printf("* %s\n", user.Name)
	}
	return nil
}

func aggregate(s *state.State, _ command) error {
	feed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml", *http.DefaultClient)

	if err != nil {
		return err
	}
	fmt.Println(feed)
	return nil
}

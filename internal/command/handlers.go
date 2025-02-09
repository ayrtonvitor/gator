package command

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
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
	ctx := context.Background()
	deletes := make([]func(*database.Queries, context.Context) error, 0)
	deletes = append(deletes, (*database.Queries).DeleteUsers, (*database.Queries).DeleteFeeds)
	for _, delFunc := range deletes {
		err := delFunc(s.Db, ctx)
		if err != nil {
			return fmt.Errorf("Unsuccessful reset. Database might be in invalid state: %w", err)
		}
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

func addFeed(s *state.State, c command) error {
	if len(c.Args) != 2 {
		return errors.New("Command addFeed gets exactly two arguments `name` and `url` in this order")
	}

	newFeed := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      c.Args[0],
		Url:       c.Args[1],
	}
	feed, err := s.Db.CreateFeed(context.Background(), newFeed)
	if err != nil {
		return fmt.Errorf("Could add new feed %s: %w", c.Args[0], err)
	}
	fmt.Printf("New feed %s added: %v\n", feed.Name, feed)
	return nil
}

func listFeeds(s *state.State, _ command) error {
	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Could not get feeds from db: %w", err)
	}
	for _, feed := range feeds {
		identedRowItem, err := json.MarshalIndent(struct {
			Name     string
			Url      string
			UserName string
		}{
			Name:     feed.Name,
			Url:      feed.Url,
			UserName: feed.UserName.String,
		}, "", "  ")
		if err != nil {
			fmt.Printf("%v\n", feed)
			continue
		}
		replacer := strings.NewReplacer("\"", "", "{\n", "", "\n}", "")
		fmt.Printf("\n%s\n", replacer.Replace(string(identedRowItem)))
	}
	return nil
}

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

	err = s.Config.SetUser(usr.Name, usr.ID)
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
	return s.Config.SetUser(userModel.Name, userModel.ID)

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

func aggregate(s *state.State, c command) error {
	if len(c.Args) != 1 {
		return errors.New("Aggregate command takes a single required argument `delay`")
	}

	client := &http.Client{Timeout: time.Duration(10 * time.Second)}

	delay, err := getScrapringInterval(c.Args[0])
	if err != nil {
		return fmt.Errorf("Could not start aggregating feed %w", err)
	}
	ticker := time.NewTicker(delay)
	for ; ; <-ticker.C {
		err := scrapeFeed(s, client)
		if err != nil {
			fmt.Printf("Error while scraping feed: %v", errors.Unwrap(err))
		}
	}
}

func addFeed(s *state.State, c command, user database.User) error {
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
	return follow(s, command{Name: "follow", Args: []string{c.Args[1]}}, user)
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

func follow(s *state.State, c command, user database.User) error {
	if len(c.Args) != 1 {
		return errors.New("Command follow gets exactly one argument `feed url`")
	}
	feed, err := s.Db.GetFeedByUrl(context.Background(), c.Args[0])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("Can not follow feed %s because it is not registered", c.Args[0])
		}
		return fmt.Errorf("Could not get feed to follow feed: %w", err)
	}

	feedFollow, err := s.Db.CreateFeedFollow(context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    feed.ID,
		})
	if err != nil {
		return fmt.Errorf("Could not follow feed: %w", err)
	}
	fmt.Printf("%s is now following %s\n", s.Config.CurrentUserName, feedFollow.FeedName)
	return nil
}

func following(s *state.State, c command) error {
	if len(c.Args) > 0 {
		return errors.New("Command following gets no arguments")
	}
	feeds, err := s.Db.GetFeedFollowsForUser(context.Background(), s.Config.CurrentUserName)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("Could not get the feeds being followed\n")
	}
	fmt.Println(s.Config.CurrentUserName + " follows:")
	for _, feed := range feeds {
		fmt.Println(feed.FeedName)
	}
	return nil
}

func unfollow(s *state.State, c command, user database.User) error {
	if len(c.Args) != 1 {
		return errors.New("Command unfollow get exactly one argument `feed url`")
	}

	feed, err := s.Db.GetFeedByUrl(context.Background(), c.Args[0])
	if err != nil {
		return fmt.Errorf("Could not get feed %s to unfollow: %w", c.Args[0], err)
	}

	err = s.Db.UnfollowFeed(context.Background(), database.UnfollowFeedParams{UserID: user.ID, FeedID: feed.ID})
	if err != nil {
		return fmt.Errorf("Could not unfollow %s", c.Args[0])
	}
	fmt.Printf("Successfully unfollowed %s\n", feed.Name)
	return nil
}

func middlewareLoggedIn(handler func(s *state.State, c command, user database.User) error) func(*state.State, command) error {
	return func(s *state.State, c command) error {
		usr, err := s.Db.GetUser(context.Background(), s.Config.CurrentUserName)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("User %s does not exist. Register a new user with `register`", s.Config.CurrentUserName)
			}
			return fmt.Errorf("Could not login with user %s: %w", s.Config.CurrentUserName, err)
		}
		return handler(s, c, usr)
	}
}

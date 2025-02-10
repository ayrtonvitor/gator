package command

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ayrtonvitor/gator/internal/database"
	"github.com/ayrtonvitor/gator/internal/rss"
	"github.com/ayrtonvitor/gator/internal/state"
	"github.com/google/uuid"
)

func scrapeFeed(s *state.State, client *http.Client) error {
	next, err := s.Db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("Could not get next feed to fetch: %w", err)
	}

	feed, err := rss.FetchFeed(context.Background(), next.Url, client)
	if err != nil {
		return fmt.Errorf("Could not fetch feed %s: %w", next.Name, err)
	}
	s.Db.MarkFeedFetched(context.Background(), next.ID)

	err = errors.Join()
	for _, item := range feed.Channel.Item {
		hasTitle := len(item.Title) != 0
		hasDescription := len(item.Description) != 0
		publicationTime := parsePublTime(item.PubDate)
		_, dbErr := s.Db.CreatePost(context.Background(),
			database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Title:       sql.NullString{String: item.Title, Valid: hasTitle},
				Url:         item.Link,
				Description: sql.NullString{String: item.Description, Valid: hasDescription},
				PublishedAt: publicationTime,
				FeedID:      next.ID,
			})
		if dbErr != nil {
			err = errors.Join(err, dbErr)
		}
		fmt.Println(sql.NullString{String: item.Description, Valid: hasDescription})
	}
	return err
}

func getScrapringInterval(arg string) (time.Duration, error) {
	const errorMsg string = "in the format #d where `d` is the time interval multiplier (`s`, `m`, `h`, `d` for seconds,\n" +
		"minutes, hours and days, respectively and `#` is the number fo such intervals"
	if len(arg) < 2 {
		return 0, fmt.Errorf(errorMsg)
	}
	s, err := getScrapingIntervalMult(arg[len(arg)-1:])
	if err != nil {
		return 0, fmt.Errorf(errorMsg)
	}
	n, err := strconv.Atoi(arg[:len(arg)-1])
	if err != nil {
		return 0, fmt.Errorf(errorMsg)
	}

	return time.Duration(n) * s, nil
}

func getScrapingIntervalMult(c string) (time.Duration, error) {
	c = strings.ToLower(c)
	switch c {
	case "s":
		return time.Second, nil
	case "m":
		return time.Minute, nil
	case "h":
		return time.Hour, nil
	default:
		return 0, fmt.Errorf("Invalid suffix")
	}
}

func parsePublTime(timeAsString string) sql.NullTime {
	parsed, err := time.Parse(time.RFC3339, timeAsString)
	if err != nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{
		Time:  parsed,
		Valid: true,
	}
}

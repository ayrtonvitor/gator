package command

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ayrtonvitor/gator/internal/rss"
	"github.com/ayrtonvitor/gator/internal/state"
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

	feed.PrintFeed()
	return nil
}

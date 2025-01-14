package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Arguments) < 1 {
		return fmt.Errorf("usage: %s <name> <time_between_requests>", cmd.Name)
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.Arguments[0])
	if err != nil {
		return errors.New("error parsing time duration")
	}

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return errors.New("error fetching next feed from database")
	}
	err = s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return errors.New("error marking feed as fetched")
	}

	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("error trying to fetch feed: %w", err)
	}
	for _, item := range rssFeed.Channel.Item {
		fmt.Printf(" * %s\n", item.Title)
	}
	return nil
}

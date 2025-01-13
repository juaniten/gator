package main

import (
	"context"
	"fmt"
)

func handlerAgg(s *state, cmd command) error {
	feedUrl := "https://www.wagslane.dev/index.xml"
	rssFeed, err := fetchFeed(context.Background(), feedUrl)
	if err != nil {
		return fmt.Errorf("error trying to fetch feed: %w", err)
	}
	fmt.Printf("Requested feed: %+v\n", rssFeed)
	return nil
}

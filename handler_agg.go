package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/juaniten/gator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Arguments) < 1 {
		return fmt.Errorf("usage: %s <time_between_requests>", cmd.Name)
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("invalid time duration string: %w", err)
	}

	fmt.Printf("Collecting feeds every %s...\n", timeBetweenRequests)
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
	log.Println("Found next feed to fetch.")

	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("error trying to fetch feed: %w", err)
	}
	err = s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return fmt.Errorf("error marking feed as fetched: %w", err)
	}

	for _, post := range rssFeed.Channel.Item {
		savePost(s, feed.ID, post)
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))

	return nil
}

func savePost(s *state, feedId uuid.UUID, post RSSItem) {
	pubDate, err := time.Parse(time.RFC1123Z, post.PubDate)
	valid := err != nil

	_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Title:     post.Title,
		Url:       post.Link,
		Description: sql.NullString{
			String: post.Description,
			Valid:  true,
		},
		PublishedAt: sql.NullTime{
			Time:  pubDate,
			Valid: valid,
		},
		FeedID: feedId,
	})
	if err != nil {
		if !errorIsUrlCollision(err) {
			fmt.Println(fmt.Errorf("error creating post: %w", err).Error())
		}
	}
}

func errorIsUrlCollision(err error) bool {
	return strings.Contains(err.Error(), `duplicate key value violates unique constraint "posts_url_key"`)
}

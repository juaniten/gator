package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/juaniten/gator/internal/database"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Arguments) < 1 {
		return fmt.Errorf("usage: %s <feed_url>", cmd.Name)
	}
	url := cmd.Arguments[0]

	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("couldn't get feed: %w", err)
	}

	args := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	followRow, err := s.db.CreateFeedFollow(context.Background(), args)
	if err != nil {
		return fmt.Errorf("couldn't create feed follow: %w", err)
	}

	fmt.Println("Feed Follow created:")
	printFeedFollow(followRow.UserName, followRow.FeedName)
	return nil

}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Arguments) < 1 {
		return fmt.Errorf("usage: %s <feed_url>", cmd.Name)
	}
	url := cmd.Arguments[0]

	args := database.DeleteFeedFollowParams{
		Name: user.Name,
		Url:  url,
	}
	err := s.db.DeleteFeedFollow(context.Background(), args)
	if err != nil {
		return fmt.Errorf("couldn't delete feed follow: %w", err)
	}

	fmt.Printf("Feed unfollowed successfully.\n")
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return fmt.Errorf("couldn't get feed follows: %w", err)
	}

	if len(feedFollows) == 0 {
		fmt.Printf("No feed follows found for user %s.\n", user.Name)
		return nil
	}
	fmt.Printf("User %s is following these feeds: \n", user.Name)
	for _, feedFollow := range feedFollows {
		fmt.Printf(" * %s\n", feedFollow.FeedName)
	}
	return nil
}

func printFeedFollow(username, feedname string) {
	fmt.Printf("* User:          %s\n", username)
	fmt.Printf("* Feed:          %s\n", feedname)
}

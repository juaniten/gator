package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/juaniten/gator/internal/database"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) < 2 {
		log.Fatal("command called without enough arguments")
	}
	name := cmd.arguments[0]
	url := cmd.arguments[1]

	args := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	}
	feed, err := s.db.CreateFeed(context.Background(), args)
	if err != nil {
		log.Fatal("error creating feed: ", err)
	}
	fmt.Printf("%+v\n", feed)

	argsFeedFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	_, err = s.db.CreateFeedFollow(context.Background(), argsFeedFollow)
	if err != nil {
		log.Fatal("error creating feed follow: ", err)
	}

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("Feeds: %+v\n", feeds)
	for _, feed := range feeds {
		fmt.Printf("\nFeed: %+v\n", feed)
		username, err := s.db.GetUserName(context.Background(), feed.UserID)
		if err != nil {
			fmt.Printf("error retrieving username: %v\n", err)
		}
		fmt.Println("Name of the feed's user: ", username)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) < 1 {
		log.Fatal("command follow called without enough arguments")
	}
	url := cmd.arguments[0]

	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		log.Fatal("error getting feed by url: ", err)
	}
	fmt.Printf("Found feed: %v\n", feed.ID)

	args := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	followRow, err := s.db.CreateFeedFollow(context.Background(), args)
	if err != nil {
		log.Fatal("error getting feed follow: ", err)
	}

	fmt.Printf("Feed Follow: %+v\n", followRow)
	return nil

}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) < 1 {
		log.Fatal("command unfollow called without enough arguments")
	}
	url := cmd.arguments[0]

	args := database.DeleteFeedFollowParams{
		Name: user.Name,
		Url:  url,
	}
	err := s.db.DeleteFeedFollow(context.Background(), args)
	if err != nil {
		log.Fatal("error deleting feed by url: ", err)
	}

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		log.Fatal("error getting feed follows: ", err)
	}
	fmt.Printf("User %s is following these feeds: \n", user.Name)
	for _, feedFollow := range feedFollows {
		fmt.Println(feedFollow.FeedName)
	}
	return nil
}
